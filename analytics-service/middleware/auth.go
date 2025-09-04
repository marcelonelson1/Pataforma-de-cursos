package middleware

import (
	"analytics-service/config"
	"analytics-service/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := utils.ExtractTokenFromHeader(c)
		if err != nil {
			utils.UnauthorizedResponse(c, "Token requerido")
			c.Abort()
			return
		}

		claims, err := utils.ValidateJWT(tokenString, cfg.JWTSecret)
		if err != nil {
			utils.UnauthorizedResponse(c, "Token inválido")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Set("claims", claims)
		c.Next()
	}
}

func AdminMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			utils.UnauthorizedResponse(c, "Token requerido")
			c.Abort()
			return
		}

		userClaims, ok := claims.(*utils.Claims)
		if !ok {
			utils.UnauthorizedResponse(c, "Token inválido")
			c.Abort()
			return
		}

		if !utils.IsAdmin(userClaims) {
			utils.ForbiddenResponse(c, "Acceso denegado: Se requieren permisos de administrador")
			c.Abort()
			return
		}

		c.Next()
	}
}