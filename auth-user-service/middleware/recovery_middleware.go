package middleware

import (
	"auth-user-service/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware middleware para manejar panics
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Printf("Panic recuperado: %v", recovered)
		
		utils.SendErrorMessage(c, "Error interno del servidor", http.StatusInternalServerError)
		c.Abort()
	})
}

// JSONMiddleware asegura que todas las respuestas sean JSON
func JSONMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Next()
	}
}