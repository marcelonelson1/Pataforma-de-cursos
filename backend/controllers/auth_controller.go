// controllers/auth_controller.go
package controllers

import (
	"curso-platform/middleware"
	"curso-platform/models"
	"curso-platform/services"
	"curso-platform/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthController gestiona las operaciones relacionadas con la autenticación
type AuthController struct {
	authService          *services.AuthService
	passwordResetService *services.PasswordResetService
}

// NewAuthController crea una nueva instancia del controlador de autenticación
func NewAuthController(authService *services.AuthService, passwordResetService *services.PasswordResetService) *AuthController {
	return &AuthController{
		authService:          authService,
		passwordResetService: passwordResetService,
	}
}

// Register maneja el registro de nuevos usuarios
func (c *AuthController) Register(ctx *gin.Context) {
	var req models.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := c.authService.Register(req); err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Usuario registrado correctamente"})
}

// Login maneja el inicio de sesión de usuarios
func (c *AuthController) Login(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	response, err := c.authService.Login(req)
	if err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusUnauthorized)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// RefreshToken renueva el token de un usuario
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	// Obtener el usuario del contexto (establecido por AuthMiddleware)
	userValue, exists := ctx.Get("user")
	if !exists {
		utils.SendErrorResponse(ctx, utils.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	user, ok := userValue.(models.Usuario)
	if !ok {
		utils.SendErrorResponse(ctx, utils.ErrServerError, http.StatusInternalServerError)
		return
	}

	// Generar nuevo token
	token, err := c.authService.RefreshToken(user.ID, user.Role)
	if err != nil {
		utils.SendErrorResponse(ctx, utils.ErrServerError, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"token":   token,
	})
}

// CheckAdmin verifica si el usuario actual es administrador
func (c *AuthController) CheckAdmin(ctx *gin.Context) {
	// Obtener el usuario del contexto (establecido por AuthMiddleware)
	userValue, exists := ctx.Get("user")
	if !exists {
		utils.SendErrorResponse(ctx, utils.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	user, ok := userValue.(models.Usuario)
	if !ok {
		utils.SendErrorResponse(ctx, utils.ErrServerError, http.StatusInternalServerError)
		return
	}

	// Verificar si el usuario tiene rol de administrador
	isAdmin := user.Role == "admin"

	// Si no es admin, registrar el intento para auditoría
	if !isAdmin {
		log.Printf("Verificación de admin fallida. Usuario ID: %d, Email: %s, Nombre: %s %s, Rol: %s",
			user.ID, user.Email, user.Nombre, user.Apellido, user.Role)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"isAdmin": isAdmin,
	})
}

// ForgotPassword maneja la solicitud de restablecimiento de contraseña
func (c *AuthController) ForgotPassword(ctx *gin.Context) {
	var req models.ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	log.Printf("Solicitud de recuperación recibida para: %s", req.Email)

	// Respuesta base segura (no revelar si el email existe)
	response := gin.H{
		"message": "Si el email existe en nuestra base de datos, recibirás un enlace para restablecer tu contraseña.",
	}

	// Usar el servicio para generar el token
	reset, err := c.passwordResetService.GenerateResetToken(req.Email)
	if err != nil {
		// No revelamos problemas específicos por seguridad
		log.Printf("Error al generar token: %v", err)
		ctx.JSON(http.StatusOK, response)
		return
	}

	// Construir el enlace de restablecimiento
	frontendURL := utils.GetEnv("FRONTEND_URL", "http://localhost:3000")
	resetLink := frontendURL + "/reset-password/" + reset.Token
	log.Printf("Enlace de recuperación generado: %s", resetLink)

	// En una aplicación real, buscaríamos los detalles del usuario y su nombre completo
	var user models.Usuario
	if result := c.passwordResetService.GetDB().Where("email = ?", req.Email).First(&user); result.Error != nil {
		log.Printf("Error al buscar usuario para email: %v", result.Error)
		ctx.JSON(http.StatusOK, response)
		return
	}

	// Usar el nombre completo del usuario
	nombreCompleto := user.GetFullName()

	// Enviar correo electrónico usando el servicio
	emailError := c.passwordResetService.SendPasswordResetEmail(req.Email, nombreCompleto, resetLink)

	// En entorno de producción, mostrar también el enlace en los logs
	// Esto es útil para depurar problemas de envío de correo en producción
	// sin afectar la seguridad de la aplicación (los logs son solo visibles para administradores)
	log.Printf("INFO IMPORTANTE: Token generado: %s", reset.Token)
	log.Printf("INFO IMPORTANTE: Enlace de recuperación: %s", resetLink)
	log.Printf("Si el usuario no recibe el correo, proporciona este enlace manualmente")

	if emailError != nil {
		// En producción, registrar el error pero aún así dar respuesta de éxito al usuario
		log.Printf("Error al enviar correo: %v", emailError)
		log.Printf("Proporciona manualmente el enlace: %s", resetLink)
	} else {
		log.Printf("Correo enviado correctamente a: %s (%s)", req.Email, nombreCompleto)
	}

	ctx.JSON(http.StatusOK, response)
}

// ValidateResetToken valida un token de restablecimiento
func (c *AuthController) ValidateResetToken(ctx *gin.Context) {
	token := ctx.Param("token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "token requerido", "valid": false})
		return
	}

	log.Printf("Validando token de recuperación: %s", token)

	// Validar el token usando el servicio
	reset, err := c.passwordResetService.ValidateResetToken(token)
	if err != nil {
		log.Printf("Token inválido: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "token inválido o expirado", "valid": false})
		return
	}

	log.Printf("Token válido para email: %s", reset.Email)
	ctx.JSON(http.StatusOK, gin.H{"valid": true, "email": reset.Email})
}

// ResetPassword restablece la contraseña de un usuario
func (c *AuthController) ResetPassword(ctx *gin.Context) {
	var req models.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	log.Printf("Solicitud de restablecimiento de contraseña recibida para token: %s", req.Token)

	// Usar el servicio para restablecer la contraseña
	err := c.passwordResetService.ResetPassword(req.Token, req.Password)
	if err != nil {
		log.Printf("Error al restablecer contraseña: %v", err)
		utils.SendErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	log.Println("Contraseña restablecida correctamente")
	ctx.JSON(http.StatusOK, gin.H{"message": "Contraseña actualizada correctamente"})
}

// RegisterRoutes registra todas las rutas relacionadas con la autenticación
func (c *AuthController) RegisterRoutes(router *gin.Engine) {
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", c.Register)
		auth.POST("/login", c.Login)
		auth.POST("/forgot-password", c.ForgotPassword)
		auth.GET("/reset-password/:token/validate", c.ValidateResetToken)
		auth.POST("/reset-password", c.ResetPassword)
		auth.GET("/check-admin", middleware.AuthMiddleware(), c.CheckAdmin)

		profile := auth.Group("")
		profile.Use(middleware.AuthMiddleware())
		{
			profile.POST("/refresh-token", c.RefreshToken)
		}
	}
}
