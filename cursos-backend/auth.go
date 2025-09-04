package main

import (
	"errors"
	"net/http"
	"strings"
	"time"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Estructuras para solicitudes y respuestas
type RegisterRequest struct {
	Nombre   string `json:"nombre" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"` // Opcional, por defecto será "user"
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string  `json:"token"`
	User  Usuario `json:"user"`
}

// JWT Claims personalizado
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"` // Añadimos el rol a los claims
	jwt.RegisteredClaims
}

// Secret para JWT
var jwtSecret = []byte(getEnv("JWT_SECRET", "mi_clave_secreta_muy_segura"))

// Controlador para registro de usuarios
func register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	// Verificar si el email ya existe
	var existingUser Usuario
	if result := db.Where("email = ?", req.Email).First(&existingUser); result.Error == nil {
		SendErrorResponse(c, ErrEmailExists, http.StatusBadRequest)
		return
	}

	// Hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		SendErrorResponse(c, errors.New("error al procesar la contraseña"), http.StatusInternalServerError)
		return
	}

	// Establecer rol por defecto si no se proporciona
	role := req.Role
	if role == "" {
		role = "user" // Por defecto, todos los usuarios son "user"
	}

	// Crear nuevo usuario
	user := Usuario{
		Nombre:   req.Nombre,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     role,
	}

	if result := db.Create(&user); result.Error != nil {
		SendErrorResponse(c, errors.New("error al crear usuario"), http.StatusInternalServerError)
		return
	}

	// Usar directamente c.JSON como en el código original que funcionaba
	c.JSON(http.StatusCreated, gin.H{"message": "Usuario registrado correctamente"})
}

// Controlador para login de usuarios
func login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	// Buscar usuario por email
	var user Usuario
	if result := db.Where("email = ?", req.Email).First(&user); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			SendErrorResponse(c, ErrInvalidLogin, http.StatusUnauthorized)
			return
		}
		SendErrorResponse(c, errors.New("error al buscar usuario"), http.StatusInternalServerError)
		return
	}

	// Verificar contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		SendErrorResponse(c, ErrInvalidLogin, http.StatusUnauthorized)
		return
	}

	// Actualizar última conexión
	now := time.Now()
	if err := db.Model(&user).Update("last_login", now).Error; err != nil {
		// Solo registramos el error, no detenemos el proceso de login
		log.Printf("Error al actualizar última conexión: %v", err)
	}

	// Generar token JWT
	token, err := generateToken(user.ID, user.Role)
	if err != nil {
		SendErrorResponse(c, errors.New("error al generar token"), http.StatusInternalServerError)
		return
	}

	// No enviar la contraseña en la respuesta
	user.Password = ""

	// Actualizar last_login en la instancia local para devolverla en la respuesta
	user.LastLogin = now

	// Usar el formato original de AuthResponse
	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User:  user,
	})
}

// Función para generar token JWT
func generateToken(userID uint, role string) (string, error) {
	// Aumentamos el tiempo de expiración para evitar desconexiones frecuentes
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		Role:   role, // Incluir el rol en el token
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Middleware de autenticación
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el token del encabezado de autorización
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			SendErrorResponse(c, errors.New("token de autorización requerido"), http.StatusUnauthorized)
			c.Abort()
			return
		}

		// Verificar que el token tenga el formato correcto
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			SendErrorResponse(c, errors.New("formato de token inválido"), http.StatusUnauthorized)
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Validar el token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			SendErrorResponse(c, ErrInvalidToken, http.StatusUnauthorized)
			c.Abort()
			return
		}

		// Verificar si el usuario existe en la base de datos
		var user Usuario
		if result := db.First(&user, claims.UserID); result.Error != nil {
			SendErrorResponse(c, ErrUserNotFound, http.StatusUnauthorized)
			c.Abort()
			return
		}

		// Verificar que el rol en el token coincida con el de la base de datos
		if claims.Role != user.Role {
			// Si han cambiado los roles, actualizamos la información del usuario
			log.Printf("Diferencia en roles: Token (%s) vs. DB (%s) para usuario ID: %d", claims.Role, user.Role, user.ID)
		}

		// Añadir el usuario al contexto para que los controladores puedan acceder a él
		c.Set("user", user)

		c.Next()
	}
}

// adminMiddleware verifica si el usuario autenticado tiene rol de administrador
func adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el usuario del contexto (establecido por authMiddleware)
		userValue, exists := c.Get("user")
		if !exists {
			SendErrorResponse(c, ErrUnauthorized, http.StatusUnauthorized)
			c.Abort()
			return
		}

		user, ok := userValue.(Usuario)
		if !ok {
			SendErrorResponse(c, errors.New("error al obtener información del usuario"), http.StatusInternalServerError)
			c.Abort()
			return
		}

		// Verificar si el usuario tiene rol de administrador
		if user.Role != "admin" {
			// Registrar intento de acceso no autorizado para auditoría de seguridad
			log.Printf("Intento de acceso a área administrativa por usuario sin permisos. ID: %d, Email: %s, Rol: %s", 
				user.ID, user.Email, user.Role)
				
			SendErrorResponse(c, errors.New("acceso denegado: se requiere rol de administrador"), http.StatusForbidden)
			c.Abort()
			return
		}

		// Verificar adicionalmente en la base de datos para asegurar que el rol no ha cambiado
		var dbUser Usuario
		if result := db.Select("role").First(&dbUser, user.ID); result.Error != nil {
			log.Printf("Error al verificar rol en DB para usuario ID: %d: %v", user.ID, result.Error)
			SendErrorResponse(c, errors.New("error al verificar permisos de administrador"), http.StatusInternalServerError)
			c.Abort()
			return
		}

		// Comprobar que el rol en la base de datos sigue siendo admin
		if dbUser.Role != "admin" {
			log.Printf("Discrepancia de roles detectada. Token: admin, DB: %s para usuario ID: %d", dbUser.Role, user.ID)
			SendErrorResponse(c, errors.New("acceso denegado: usuario ya no tiene permisos de administrador"), http.StatusForbidden)
			c.Abort()
			return
		}

		// Registrar acceso administrativo exitoso para auditoría
		log.Printf("Acceso administrativo exitoso. Usuario ID: %d, Email: %s", user.ID, user.Email)

		c.Next()
	}
}

// checkAdmin verifica y responde si el usuario tiene rol admin
func checkAdmin(c *gin.Context) {
	// Obtener el usuario del contexto (establecido por authMiddleware)
	userValue, exists := c.Get("user")
	if !exists {
		SendErrorResponse(c, ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	user, ok := userValue.(Usuario)
	if !ok {
		SendErrorResponse(c, errors.New("error al obtener información del usuario"), http.StatusInternalServerError)
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

// refreshToken extiende la sesión del usuario generando un nuevo token
func refreshToken(c *gin.Context) {
	// Obtener el usuario del contexto (establecido por authMiddleware)
	userValue, exists := c.Get("user")
	if !exists {
		SendErrorResponse(c, ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	user, ok := userValue.(Usuario)
	if !ok {
		SendErrorResponse(c, errors.New("error al obtener información del usuario"), http.StatusInternalServerError)
		return
	}

	// Generar nuevo token
	token, err := generateToken(user.ID, user.Role)
	if err != nil {
		SendErrorResponse(c, errors.New("error al renovar el token"), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"token": token,
	})
}