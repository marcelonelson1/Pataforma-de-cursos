package middleware

import (
	"auth-user-service/config"
	"log"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupCORS configura el middleware CORS
func SetupCORS() gin.HandlerFunc {
	// Obtener orígenes permitidos de la configuración
	allowedOrigins := strings.Split(config.AppConfig.CORSOrigins, ",")
	
	// Limpiar espacios en blanco
	for i, origin := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(origin)
	}

	log.Printf("Configurando CORS para permitir orígenes: %v", allowedOrigins)

	corsConfig := cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	return cors.New(corsConfig)
}