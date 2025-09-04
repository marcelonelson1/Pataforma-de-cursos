package middleware

import (
	"portfolio-service/config"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupCORS configura el middleware de CORS
func SetupCORS() gin.HandlerFunc {
	frontendURL := "http://localhost:3000"
	if config.AppConfig != nil {
		// Puedes agregar una variable de entorno para la URL del frontend
		if envURL := config.AppConfig.AppEnv; envURL == "production" {
			frontendURL = "https://yourproductiondomain.com"
		}
	}

	log.Printf("Configurando CORS para permitir origen: %s", frontendURL)

	return cors.New(cors.Config{
		AllowOrigins:     []string{frontendURL, "http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

// JSONMiddleware fuerza respuestas JSON
func JSONMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Next()
	})
}

// RecoveryMiddleware maneja panics y errores inesperados
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Printf("Panic recuperado: %v", recovered)
		
		c.JSON(500, gin.H{
			"success": false,
			"error":   "Error interno del servidor",
		})
		
		c.Abort()
	})
}