package middleware

import (
	"auth-user-service/services"
	"auth-user-service/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware middleware de autenticación JWT
func AuthMiddleware() gin.HandlerFunc {
	authService := services.NewAuthService()

	return func(c *gin.Context) {
		// Obtener el token del encabezado de autorización
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.SendErrorMessage(c, "token de autorización requerido", http.StatusUnauthorized)
			c.Abort()
			return
		}

		// Verificar que el token tenga el formato correcto
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			utils.SendErrorMessage(c, "formato de token inválido", http.StatusUnauthorized)
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Validar el token
		validationResponse, err := authService.ValidateToken(tokenString)
		if err != nil || !validationResponse.Valid {
			utils.SendErrorResponse(c, utils.ErrInvalidToken, http.StatusUnauthorized)
			c.Abort()
			return
		}

		// Añadir el usuario al contexto para que los controladores puedan acceder a él
		c.Set("user", validationResponse.User)
		c.Set("user_id", validationResponse.UserID)
		c.Set("user_role", validationResponse.Role)

		c.Next()
	}
}

// OptionalAuthMiddleware middleware de autenticación opcional (no falla si no hay token)
func OptionalAuthMiddleware() gin.HandlerFunc {
	authService := services.NewAuthService()

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := tokenParts[1]
		validationResponse, err := authService.ValidateToken(tokenString)
		if err == nil && validationResponse.Valid {
			c.Set("user", validationResponse.User)
			c.Set("user_id", validationResponse.UserID)
			c.Set("user_role", validationResponse.Role)
		}

		c.Next()
	}
}