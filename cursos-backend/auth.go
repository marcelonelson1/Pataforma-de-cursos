package main

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Estructuras para solicitudes y respuestas
type RegisterRequest struct {
	Nombre    string `json:"nombre" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
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
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// Secret para JWT - en producción usar variable de entorno
var jwtSecret = []byte(getEnv("JWT_SECRET", "mi_clave_secreta_muy_segura"))

// Controlador para registro de usuarios
func register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar si el email ya existe
	var existingUser Usuario
	if result := db.Where("email = ?", req.Email).First(&existingUser); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Este email ya está registrado"})
		return
	}

	// Hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar la contraseña"})
		return
	}

	// Crear nuevo usuario
	user := Usuario{
		Nombre:   req.Nombre,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if result := db.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear usuario"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Usuario registrado correctamente"})
}

// Controlador para login de usuarios
func login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Buscar usuario por email
	var user Usuario
	if result := db.Where("email = ?", req.Email).First(&user); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar usuario"})
		return
	}

	// Verificar contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
		return
	}

	// Generar token JWT
	token, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar token"})
		return
	}

	// No enviar la contraseña en la respuesta
	user.Password = ""

	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User:  user,
	})
}

// Middleware de autenticación
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Autorización requerida"})
			return
		}

		// Extraer token del header
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// Validar token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			return
		}

		// Verificar si el usuario existe
		var user Usuario
		if result := db.First(&user, claims.UserID); result.Error != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Usuario no encontrado"})
			return
		}

		// Guardar ID del usuario en el contexto
		c.Set("userID", claims.UserID)
		c.Next()
	}
}

// Función para generar token JWT
func generateToken(userID uint) (string, error) {
	// Definir tiempo de expiración (24 horas)
	expirationTime := time.Now().Add(24 * time.Hour)

	// Crear claims
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Crear token con claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar token y retornarlo
	return token.SignedString(jwtSecret)
}