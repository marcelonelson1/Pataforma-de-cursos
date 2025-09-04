// middleware/auth.go - VERSIÓN CORREGIDA
package middleware

import (
	"log"
	"net/http"
	"strings"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"payment-service/config"
	"payment-service/utils"
)

// Claims estructura para JWT
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware middleware de autenticación JWT para payment-service
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener token del header Authorization
		authHeader := c.GetHeader("Authorization")
		log.Printf("🔍 AUTH HEADER RECIBIDO: '%s'", authHeader)
		
		if authHeader == "" {
			log.Printf("❌ Error: No se proporcionó header Authorization")
			utils.SendErrorResponse(c, "Token de autorización requerido", http.StatusUnauthorized)
			return
		}

		// Verificar formato Bearer
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			log.Printf("❌ Error: Formato de token inválido: %s", authHeader)
			utils.SendErrorResponse(c, "Formato de token inválido", http.StatusUnauthorized)
			return
		}

		tokenString := tokenParts[1]
		
		// 🔥 CORRECCIÓN: Verificar longitud del token antes de hacer slice
		if len(tokenString) > 20 {
			log.Printf("🔍 TOKEN EXTRAÍDO: %s...", tokenString[:20])
		} else {
			log.Printf("🔍 TOKEN EXTRAÍDO: %s", tokenString)
		}

		// Validar token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// 🔥 VERIFICACIÓN ADICIONAL: Comprobar método de firma
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			log.Printf("❌ Error al validar token: %v", err)
			utils.SendErrorResponse(c, "Token inválido o expirado", http.StatusUnauthorized)
			return
		}

		log.Printf("✅ Token válido para usuario ID: %d, Role: %s", claims.UserID, claims.Role)

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
			log.Printf("❌ Error: Información de usuario no encontrada en contexto")
			utils.SendErrorResponse(c, "Información de usuario no encontrada", http.StatusUnauthorized)
			return
		}

		role, ok := userRole.(string)
		if !ok || role != "admin" {
			log.Printf("❌ Error: Se requieren permisos de admin. Role actual: %s", role)
			utils.SendErrorResponse(c, "Se requieren permisos de administrador", http.StatusForbidden)
			return
		}

		log.Printf("✅ Usuario admin verificado")
		c.Next()
	}
}

// CORSMiddleware maneja CORS para payment-service
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Lista de orígenes permitidos (actualizar según tu configuración)
		allowedOrigins := []string{
			"http://localhost:3000",
			"http://localhost:3001", 
			"https://tu-dominio.com",
		}

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

// 🔥 MEJORADO: RecoveryMiddleware con mejor logging
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("❌ PANIC RECUPERADO en payment-service:")
				log.Printf("   Error: %v", err)
				log.Printf("   Path: %s %s", c.Request.Method, c.Request.URL.Path)
				log.Printf("   IP: %s", c.ClientIP())
				log.Printf("   User-Agent: %s", c.Request.UserAgent())
				
				// Verificar si es un error de slice bounds específico
				if strings.Contains(fmt.Sprintf("%v", err), "slice bounds out of range") {
					log.Printf("   🔍 CAUSA: Error de slice bounds - probablemente en middleware de auth")
				}
				
				utils.SendErrorResponse(c, "Error interno del servidor", http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

// LoggingMiddleware registra todas las requests
func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("🌐 %s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
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

// 🔥 NUEVA FUNCIÓN: Middleware opcional para debug (solo desarrollo)
func DebugMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/api/debug/mercadopago/fix" || 
		   c.Request.URL.Path == "/api/debug/mercadopago/quick-test" ||
		   strings.HasPrefix(c.Request.URL.Path, "/api/debug/") {
			log.Printf("🔧 [DEBUG] Endpoint de debug accedido: %s", c.Request.URL.Path)
		}
		c.Next()
	}
}