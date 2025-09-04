package main

import (
	"analytics-service/config"
	"analytics-service/controllers"
	"analytics-service/middleware"
	"analytics-service/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	
	log.Println("🚀 Starting Analytics Service...")
	log.Printf("📊 Environment: %s", cfg.AppEnv)
	log.Printf("🔧 Port: %s", cfg.Port)
	
	// Initialize database
	err := utils.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	
	// Run migrations
	err = utils.RunMigrations()
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	
	// Set Gin mode
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	
	// Create router
	router := gin.Default()
	
	// Add middleware
	router.Use(middleware.CORSMiddleware(cfg))
	
	// Initialize controller
	analyticsController := controllers.NewAnalyticsController(cfg)
	
	// Setup routes
	analyticsController.SetupRoutes(router)
	
	// Start server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	log.Printf("✅ Analytics Service started on port %s", cfg.Port)
	log.Printf("🔗 Health Check: http://localhost:%s/api/health", cfg.Port)
	log.Printf("📊 Dashboard: http://localhost:%s/api/admin/dashboard", cfg.Port)
	log.Printf("📈 Sales Stats: http://localhost:%s/api/admin/sales-stats", cfg.Port)
	log.Printf("👤 Profile: http://localhost:%s/api/auth/profile", cfg.Port)
	log.Printf("🔐 Admin Check: http://localhost:%s/api/auth/check-admin", cfg.Port)
	
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}