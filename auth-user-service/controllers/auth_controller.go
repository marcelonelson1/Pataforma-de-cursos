package controllers

import (
	"auth-user-service/config"
	"auth-user-service/models"
	"auth-user-service/services"
	"auth-user-service/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService  *services.AuthService
	emailService *services.EmailService
}

func NewAuthController() *AuthController {
	return &AuthController{
		authService:  services.NewAuthService(),
		emailService: services.NewEmailService(),
	}
}

// Register registra un nuevo usuario
func (a *AuthController) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, err)
		return
	}

	if err := a.authService.Register(&req); err != nil {
		if err == utils.ErrEmailExists {
			utils.SendErrorResponse(c, err, http.StatusBadRequest)
		} else {
			utils.SendErrorMessage(c, "error al crear usuario", http.StatusInternalServerError)
		}
		return
	}

	// Usar directamente c.JSON como en el código original que funcionaba
	c.JSON(http.StatusCreated, gin.H{"message": "Usuario registrado correctamente"})
}

// Login autentica un usuario
func (a *AuthController) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, err)
		return
	}

	authResponse, err := a.authService.Login(&req)
	if err != nil {
		if err == utils.ErrInvalidLogin {
			utils.SendErrorResponse(c, err, http.StatusUnauthorized)
		} else {
			utils.SendErrorMessage(c, "error al procesar login", http.StatusInternalServerError)
		}
		return
	}

	// Usar el formato original de AuthResponse
	c.JSON(http.StatusOK, *authResponse)
}

// ForgotPassword solicita restablecimiento de contraseña
func (a *AuthController) ForgotPassword(c *gin.Context) {
	var req models.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, err)
		return
	}

	// Intentar generar token de reset
	reset, err := a.authService.ForgotPassword(req.Email)
	if err != nil {
		// No revelamos si el email existe o no por seguridad
		c.JSON(http.StatusOK, gin.H{"message": "Si el email existe en nuestra base de datos, recibirás un enlace para restablecer tu contraseña."})
		log.Printf("Email no encontrado o error de DB: %v", err)
		return
	}

	// Buscar el usuario para obtener el nombre
	var user models.Usuario
	if result := utils.GetDB().Where("email = ?", req.Email).First(&user); result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Si el email existe en nuestra base de datos, recibirás un enlace para restablecer tu contraseña."})
		return
	}

	// Construir el enlace de restablecimiento
	resetLink := config.AppConfig.FrontendURL + "/reset-password/" + reset.Token

	// Enviar correo electrónico
	emailError := a.emailService.SendPasswordResetEmail(user.Email, user.Nombre, resetLink)
	if emailError != nil {
		log.Printf("Error al enviar correo: %v", emailError)
		// En producción, esto podría ser un problema crítico
		if config.AppConfig.AppEnv == "production" {
			log.Printf("ERROR CRÍTICO: Fallo al enviar email en producción: %v", emailError)
		}
	}

	// Respuesta base
	response := gin.H{
		"message": "Si el email existe en nuestra base de datos, recibirás un enlace para restablecer tu contraseña.",
	}

	// Solo enviar el token y enlace en desarrollo
	if config.AppConfig.AppEnv == "development" {
		response["resetToken"] = reset.Token
		response["resetLink"] = resetLink
		// En desarrollo, podemos informar sobre el error de email si ocurrió
		if emailError != nil {
			response["emailError"] = emailError.Error()
		}
	}

	// Usar el formato original que funcionaba
	c.JSON(http.StatusOK, response)
}

// ValidateResetToken valida el token de restablecimiento
func (a *AuthController) ValidateResetToken(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token requerido", "valid": false})
		return
	}

	// Validar el token
	reset, err := a.authService.ValidateResetToken(token)
	if err != nil {
		log.Printf("Token inválido: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "token inválido o expirado", "valid": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": true, "email": reset.Email})
}

// ResetPassword restablece la contraseña
func (a *AuthController) ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, err)
		return
	}

	if err := a.authService.ResetPassword(req.Token, req.Password); err != nil {
		if err == utils.ErrInvalidToken {
			utils.SendErrorMessage(c, "token inválido o expirado", http.StatusBadRequest)
		} else {
			log.Printf("Error al restablecer contraseña: %v", err)
			utils.SendErrorMessage(c, "error al actualizar la contraseña", http.StatusInternalServerError)
		}
		return
	}

	// Devolver un JSON válido usando exactamente el mismo formato original
	c.JSON(http.StatusOK, gin.H{"message": "Contraseña actualizada correctamente"})
}

// ChangePassword cambia la contraseña del usuario
func (a *AuthController) ChangePassword(c *gin.Context) {
	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, err)
		return
	}

	// Obtener el usuario del contexto
	userValue, exists := c.Get("user")
	if !exists {
		utils.SendErrorResponse(c, utils.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	user, ok := userValue.(*models.Usuario)
	if !ok {
		utils.SendErrorMessage(c, "error al obtener información del usuario", http.StatusInternalServerError)
		return
	}

	if err := a.authService.ChangePassword(user.ID, req.CurrentPassword, req.NewPassword); err != nil {
		if err == utils.ErrInvalidPassword {
			utils.SendErrorResponse(c, err, http.StatusBadRequest)
		} else {
			utils.SendErrorMessage(c, "error al cambiar contraseña", http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contraseña cambiada correctamente"})
}

// RefreshToken extiende la sesión del usuario generando un nuevo token
func (a *AuthController) RefreshToken(c *gin.Context) {
	// Obtener el usuario del contexto
	userValue, exists := c.Get("user")
	if !exists {
		utils.SendErrorResponse(c, utils.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	user, ok := userValue.(*models.Usuario)
	if !ok {
		utils.SendErrorMessage(c, "error al obtener información del usuario", http.StatusInternalServerError)
		return
	}

	// Generar nuevo token
	token, err := a.authService.RefreshToken(user.ID)
	if err != nil {
		utils.SendErrorMessage(c, "error al renovar el token", http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"token":   token,
	})
}

// CheckAdmin verifica y responde si el usuario tiene rol admin
func (a *AuthController) CheckAdmin(c *gin.Context) {
	// Obtener el usuario del contexto
	userValue, exists := c.Get("user")
	if !exists {
		utils.SendErrorResponse(c, utils.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	user, ok := userValue.(*models.Usuario)
	if !ok {
		utils.SendErrorMessage(c, "error al obtener información del usuario", http.StatusInternalServerError)
		return
	}

	// Verificar si el usuario tiene rol de administrador
	isAdmin := user.Role == "admin"

	// Si no es admin, registrar el intento para auditoría
	if !isAdmin {
		log.Printf("Verificación de admin fallida. Usuario ID: %d, Email: %s, Rol: %s", 
			user.ID, user.Email, user.Role)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"isAdmin": isAdmin,
	})
}

// ValidateToken endpoint para validar tokens (para otros servicios)
func (a *AuthController) ValidateToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"valid": false,
			"error": "token requerido",
		})
		return
	}

	// Extraer token
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		c.JSON(http.StatusBadRequest, gin.H{
			"valid": false,
			"error": "formato de token inválido",
		})
		return
	}

	validationResponse, err := a.authService.ValidateToken(tokenParts[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"valid": false,
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, validationResponse)
}