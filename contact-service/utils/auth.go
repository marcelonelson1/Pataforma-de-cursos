package utils

import (
	"contact-service/config"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// JWTClaims estructura para los claims del JWT
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware middleware para validar JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			SendErrorResponse(c, ErrUnauthorized, http.StatusUnauthorized)
			c.Abort()
			return
		}

		// Verificar formato "Bearer token"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			SendErrorMessage(c, "formato de token invalido", http.StatusUnauthorized)
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validar token
		claims := &JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("metodo de firma inesperado: %v", token.Header["alg"])
			}
			return []byte(config.AppConfig.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			SendErrorMessage(c, "token invalido", http.StatusUnauthorized)
			c.Abort()
			return
		}

		// Guardar informacion del usuario en el contexto
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// AdminMiddleware middleware para validar rol de admin
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			SendErrorResponse(c, ErrUnauthorized, http.StatusUnauthorized)
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok || role != "admin" {
			SendErrorMessage(c, "acceso denegado: se requiere rol de administrador", http.StatusForbidden)
			c.Abort()
			return
		}

		c.Next()
	}
}