package middleware

import (
	"curso-platform/utils"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CorsMiddleware configura los encabezados CORS para permitir solicitudes de orígenes específicos
func CorsMiddleware() gin.HandlerFunc {
	frontendURL := utils.GetEnv("FRONTEND_URL", "http://localhost:3000")
	log.Printf("Configurando CORS para permitir origen: %s", frontendURL)

	return cors.New(cors.Config{
		AllowOrigins:     []string{frontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}