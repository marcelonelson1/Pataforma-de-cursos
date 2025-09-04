package utils

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// Errores predefinidos del sistema de pagos
var (
	ErrInvalidToken       = errors.New("token inválido o expirado")
	ErrUserNotFound       = errors.New("usuario no encontrado")
	ErrUnauthorized       = errors.New("no autorizado")
	ErrPaymentNotFound    = errors.New("pago no encontrado")
	ErrCourseNotFound     = errors.New("curso no encontrado")
	ErrInvalidPaymentData = errors.New("datos de pago inválidos")
	ErrPaymentExists      = errors.New("pago ya existe para este curso")
	ErrDatabaseError      = errors.New("error de base de datos")
	ErrExternalService    = errors.New("error en servicio externo")
	ErrInvalidRequest     = errors.New("solicitud inválida")
)

// Claims estructura para JWT
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware middleware de autenticación JWT
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener token del header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			SendErrorResponse(c, "Token de autorización requerido", http.StatusUnauthorized)
			return
		}

		// Verificar formato Bearer
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			SendErrorResponse(c, "Formato de token inválido", http.StatusUnauthorized)
			return
		}

		tokenString := tokenParts[1]

		// Validar token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			log.Printf("Token inválido: %v", err)
			SendErrorResponse(c, "Token inválido o expirado", http.StatusUnauthorized)
			return
		}

		// Agregar información del usuario al contexto
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// AdminMiddleware middleware para verificar permisos de administrador
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			SendErrorResponse(c, "Información de usuario no encontrada", http.StatusUnauthorized)
			return
		}

		role, ok := userRole.(string)
		if !ok || role != "admin" {
			SendErrorResponse(c, "Se requieren permisos de administrador", http.StatusForbidden)
			return
		}

		c.Next()
	}
}

// RecoveryMiddleware maneja panics y asegura respuestas JSON consistentes
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recuperado: %v", err)
				SendErrorResponse(c, "Error interno del servidor", http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

// LoggingMiddleware registra todas las requests
func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format("02/Jan/2006:15:04:05 -0700"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// CORSMiddleware maneja CORS
func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Verificar si el origin está permitido
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin,Content-Type,Accept,Authorization,X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// ValidatePaymentRequest valida los datos de una solicitud de pago
func ValidatePaymentRequest(req interface{}) error {
	// Aquí se pueden agregar validaciones específicas de negocio
	// Por ejemplo, validar montos mínimos, métodos permitidos, etc.
	return nil
}