package middleware

import (
	"curso-platform/config"
	"curso-platform/models"
	"curso-platform/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware verifica si el usuario está autenticado mediante un token JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el token del encabezado de autorización
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.SendErrorResponse(c, utils.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		// Verificar que el token tenga el formato correcto
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			utils.SendErrorResponse(c, utils.ErrInvalidToken, http.StatusUnauthorized)
			return
		}

		tokenString := tokenParts[1]

		// Validar el token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			utils.SendErrorResponse(c, utils.ErrInvalidToken, http.StatusUnauthorized)
			return
		}

		// Verificar si el usuario existe en la base de datos
		var user models.Usuario
		if result := config.DB.First(&user, claims.UserID); result.Error != nil {
			utils.SendErrorResponse(c, utils.ErrUserNotFound, http.StatusUnauthorized)
			return
		}

		// Verificar que el rol en el token coincida con el de la base de datos
		if claims.Role != user.Role {
			// Si han cambiado los roles, registramos la diferencia
			log.Printf("Diferencia en roles: Token (%s) vs. DB (%s) para usuario ID: %d (%s %s)",
				claims.Role, user.Role, user.ID, user.Nombre, user.Apellido)
		}

		// Añadir el usuario al contexto para que los controladores puedan acceder a él
		c.Set("user", user)

		c.Next()
	}
}

// AdminMiddleware verifica si el usuario autenticado tiene rol de administrador
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el usuario del contexto (establecido por AuthMiddleware)
		userValue, exists := c.Get("user")
		if !exists {
			utils.SendErrorResponse(c, utils.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		user, ok := userValue.(models.Usuario)
		if !ok {
			utils.SendErrorResponse(c, utils.ErrServerError, http.StatusInternalServerError)
			return
		}

		// Verificar si el usuario tiene rol de administrador
		if user.Role != "admin" {
			// Registrar intento de acceso no autorizado para auditoría de seguridad
			log.Printf("Intento de acceso a área administrativa por usuario sin permisos. ID: %d, Email: %s, Nombre: %s %s, Rol: %s",
				user.ID, user.Email, user.Nombre, user.Apellido, user.Role)

			utils.SendErrorResponse(c, utils.ErrUnauthorized, http.StatusForbidden)
			return
		}

		// Verificar adicionalmente en la base de datos para asegurar que el rol no ha cambiado
		var dbUser models.Usuario
		if result := config.DB.Select("role").First(&dbUser, user.ID); result.Error != nil {
			log.Printf("Error al verificar rol en DB para usuario ID: %d (%s %s): %v",
				user.ID, user.Nombre, user.Apellido, result.Error)
			utils.SendErrorResponse(c, utils.ErrServerError, http.StatusInternalServerError)
			return
		}

		// Comprobar que el rol en la base de datos sigue siendo admin
		if dbUser.Role != "admin" {
			log.Printf("Discrepancia de roles detectada. Token: admin, DB: %s para usuario ID: %d (%s %s)",
				dbUser.Role, user.ID, user.Nombre, user.Apellido)
			utils.SendErrorResponse(c, utils.ErrUnauthorized, http.StatusForbidden)
			return
		}

		// Registrar acceso administrativo exitoso para auditoría
		log.Printf("Acceso administrativo exitoso. Usuario ID: %d, Email: %s, Nombre: %s %s",
			user.ID, user.Email, user.Nombre, user.Apellido)

		c.Next()
	}
}
