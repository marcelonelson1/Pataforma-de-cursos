// utils/errors.go
package utils

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"course-service/config"
)

// Constantes de errores
const (
	ErrDatabaseError     = "Error de base de datos"
	ErrValidationError   = "Error de validación"
	ErrNotFound          = "Recurso no encontrado"
	ErrUnauthorized      = "No autorizado"
	ErrForbidden         = "Acceso prohibido"
	ErrInternalServer    = "Error interno del servidor"
	ErrBadRequest        = "Solicitud incorrecta"
	ErrFileUpload        = "Error al subir archivo"
	ErrFileNotFound      = "Archivo no encontrado"
	ErrInvalidToken      = "Token inválido"
	ErrTokenExpired      = "Token expirado"
	ErrPermissionDenied  = "Permisos insuficientes"
)

// CourseError estructura personalizada de errores
type CourseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *CourseError) Error() string {
	return e.Message
}

// NewCourseError crea un nuevo error personalizado
func NewCourseError(code int, message, details string) *CourseError {
	return &CourseError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// ErrorHandler middleware para manejo centralizado de errores
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Procesar errores si los hay
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			
			switch e := err.Err.(type) {
			case *CourseError:
				c.JSON(e.Code, gin.H{
					"success": false,
					"error":   e.Message,
					"details": e.Details,
					"code":    e.Code,
				})
			default:
				// Error genérico
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"error":   ErrInternalServer,
					"code":    http.StatusInternalServerError,
				})
			}
		}
	}
}

// RecoveryHandler middleware para recuperación de panics
func RecoveryHandler() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Printf("Panic recovered: %v\n%s", recovered, debug.Stack())
		
		SendInternalServerErrorResponse(c, "Ha ocurrido un error inesperado")
		c.Abort()
	})
}

// AuthMiddleware middleware de autenticación JWT
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			SendUnauthorizedResponse(c)
			c.Abort()
			return
		}

		// Verificar formato Bearer
		if !strings.HasPrefix(authHeader, "Bearer ") {
			SendErrorResponse(c, "Formato de token inválido", http.StatusUnauthorized)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parsear y validar token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verificar método de firma
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, NewCourseError(http.StatusUnauthorized, ErrInvalidToken, "Método de firma inválido")
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil {
			log.Printf("Error al parsear token: %v", err)
			SendErrorResponse(c, ErrInvalidToken, http.StatusUnauthorized)
			c.Abort()
			return
		}

		if !token.Valid {
			SendErrorResponse(c, ErrTokenExpired, http.StatusUnauthorized)
			c.Abort()
			return
		}

		// Extraer claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userID, exists := claims["user_id"]
			if !exists {
				SendErrorResponse(c, "Token inválido: falta user_id", http.StatusUnauthorized)
				c.Abort()
				return
			}

			// Convertir userID a uint
			var userIDUint uint
			switch v := userID.(type) {
			case float64:
				userIDUint = uint(v)
			case int:
				userIDUint = uint(v)
			case uint:
				userIDUint = v
			default:
				SendErrorResponse(c, "Token inválido: user_id no válido", http.StatusUnauthorized)
				c.Abort()
				return
			}

			// Guardar información del usuario en el contexto
			c.Set("user_id", userIDUint)
			
			// Extraer role si existe
			if role, exists := claims["role"]; exists {
				c.Set("user_role", role)
			}

			// Extraer email si existe
			if email, exists := claims["email"]; exists {
				c.Set("user_email", email)
			}
		} else {
			SendErrorResponse(c, "Claims del token inválidos", http.StatusUnauthorized)
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminMiddleware middleware para verificar permisos de administrador
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			SendForbiddenResponse(c, ErrPermissionDenied)
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok || (roleStr != "admin" && roleStr != "instructor") {
			SendForbiddenResponse(c, "Requiere permisos de administrador")
			c.Abort()
			return
		}

		c.Next()
	}
}

// InstructorMiddleware middleware para verificar permisos de instructor
func InstructorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			SendForbiddenResponse(c, ErrPermissionDenied)
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok || (roleStr != "admin" && roleStr != "instructor") {
			SendForbiddenResponse(c, "Requiere permisos de instructor")
			c.Abort()
			return
		}

		c.Next()
	}
}

// JSONMiddleware middleware para establecer headers JSON
func JSONMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	}
}

// CORSMiddleware middleware para manejo de CORS (ya incluido en main.go, pero por completitud)
func CORSMiddleware(allowOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Verificar si el origen está permitido
		allowed := false
		for _, allowedOrigin := range allowOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
		}

		c.Next()
	}
}

// ValidateJSONMiddleware middleware para validar que el contenido sea JSON válido
func ValidateJSONMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			contentType := c.GetHeader("Content-Type")
			if strings.Contains(contentType, "application/json") {
				// Gin ya maneja la validación JSON automáticamente
				// Este middleware puede extenderse para validaciones adicionales
			}
		}
		c.Next()
	}
}

// LoggerMiddleware middleware personalizado para logging
func LoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] %s %s %d %s %s\n",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.ClientIP,
		)
	})
}

// RateLimitMiddleware middleware básico de rate limiting
func RateLimitMiddleware(maxRequests int, window time.Duration) gin.HandlerFunc {
	clients := make(map[string][]time.Time)
	mutex := sync.RWMutex{}

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		now := time.Now()

		mutex.Lock()
		defer mutex.Unlock()

		// Limpiar requests antiguos
		if requests, exists := clients[clientIP]; exists {
			validRequests := []time.Time{}
			for _, reqTime := range requests {
				if now.Sub(reqTime) < window {
					validRequests = append(validRequests, reqTime)
				}
			}
			clients[clientIP] = validRequests
		}

		// Verificar límite
		if len(clients[clientIP]) >= maxRequests {
			SendErrorResponse(c, "Demasiadas solicitudes", http.StatusTooManyRequests)
			c.Abort()
			return
		}

		// Agregar request actual
		clients[clientIP] = append(clients[clientIP], now)
		c.Next()
	}
	}
	