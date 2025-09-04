// main.go
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"course-service/config"
	"course-service/routes"
	"course-service/utils"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env, usando variables del sistema")
	}

	// Cargar configuración
	cfg := config.Load()

	// Conectar a la base de datos
	db := utils.ConnectDatabase(cfg)
	
	// Migrar modelos
	utils.MigrateDatabase(db)

	// Inicializar directorios estáticos
	utils.InitStaticDirs(cfg)

	// Inicializar proveedores de pago (para compatibilidad)
	//utils.InitPaymentProviders()

	// Configurar router
	router := setupRouter(cfg)

	// Configurar rutas
	routes.SetupCourseRoutes(router, db, cfg)

	// Iniciar servidor
	startServer(router, cfg)
}

func setupRouter(cfg *config.Config) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Configurar CORS dinámico
	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	if cfg.AllowAllOrigins {
		corsConfig.AllowAllOrigins = true
		log.Printf("🌐 CORS: Permitiendo TODOS los orígenes (desarrollo)")
	} else {
		corsConfig.AllowOrigins = append(cfg.AllowedOrigins, cfg.FrontendURL)
		log.Printf("🌐 CORS: Orígenes permitidos: %v", corsConfig.AllowOrigins)
	}

	router.Use(cors.New(corsConfig))

	// Servir archivos estáticos ANTES de otros middlewares
	router.Static("/static", cfg.UploadPath)
	log.Printf("📁 Archivos estáticos disponibles en: /static -> %s", cfg.UploadPath)

	// Middlewares solo para rutas API
	router.Use(utils.JSONMiddleware())
	router.Use(utils.RecoveryHandler())
	router.Use(utils.ErrorHandler())
	router.Use(utils.LoggerMiddleware())

	// Rate limiting en producción
	if cfg.AppEnv == "production" {
		router.Use(utils.RateLimitMiddleware(100, time.Minute))
	}

	return router
}

func startServer(router *gin.Engine, cfg *config.Config) {
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  30 * time.Second,  // Aumentado para videos grandes
		WriteTimeout: 120 * time.Second, // Aumentado significativamente para streaming
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("🚀 Course Service iniciado en puerto %s", cfg.Port)
	log.Printf("📊 Ambiente: %s", cfg.AppEnv)
	log.Printf("📁 Upload Path: %s", cfg.UploadPath)
	log.Printf("💾 Cache: %v (TTL: %ds)", cfg.Cache.Enabled, cfg.Cache.TTL)
	
	if cfg.CDN.Enabled {
		log.Printf("🌐 CDN habilitado: %s", cfg.CDN.URL)
	}
	
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ Error al iniciar servidor: %v", err)
	}
}