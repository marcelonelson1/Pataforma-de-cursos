package middleware

import (
	"auth-user-service/models"
	"auth-user-service/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware verifica si el usuario autenticado tiene rol de administrador
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el usuario del contexto (establecido por AuthMiddleware)
		userValue, exists := c.Get("user")
		if !exists {
			utils.SendErrorResponse(c, utils.ErrUnauthorized, http.StatusUnauthorized)
			c.Abort()
			return
		}

		user, ok := userValue.(*models.Usuario)
		if !ok {
			utils.SendErrorMessage(c, "error al obtener información del usuario", http.StatusInternalServerError)
			c.Abort()
			return
		}

		// Verificar si el usuario tiene rol de administrador
		if user.Role != "admin" {
			// Registrar intento de acceso no autorizado para auditoría de seguridad
			log.Printf("Intento de acceso a área administrativa por usuario sin permisos. ID: %d, Email: %s, Rol: %s", 
				user.ID, user.Email, user.Role)
				
			utils.SendErrorMessage(c, "acceso denegado: se requiere rol de administrador", http.StatusForbidden)
			c.Abort()
			return
		}

		// Registrar acceso administrativo exitoso para auditoría
		log.Printf("Acceso administrativo exitoso. Usuario ID: %d, Email: %s", user.ID, user.Email)

		c.Next()
	}
}