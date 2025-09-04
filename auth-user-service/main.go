package main

import (
	"auth-user-service/config"
	"auth-user-service/routes"
	"auth-user-service/utils"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// Configurar logging
	setupLogging()

	// Cargar configuración
	config.LoadConfig()

	// Configurar base de datos
	if err := utils.SetupDatabase(); err != nil {
		log.Fatalf("Error al configurar la base de datos: %v", err)
	}

	// Crear directorios necesarios
	initStaticDirs()

	// Configurar rutas
	router := routes.SetupRoutes()

	// Iniciar servidor
	startServer(router)
}

// setupLogging configura el sistema de logging
func setupLogging() {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(logFile)
	} else {
		log.Println("No se pudo abrir el archivo de log, usando salida estándar")
	}
}

// initStaticDirs crea los directorios estáticos necesarios
func initStaticDirs() {
	dirs := []string{
		"./static",
		"./static/profiles",
		"./templates",
	}

	for _, dir := range dirs {
		if err := utils.CreateDirIfNotExists(dir); err != nil {
			log.Printf("Error al crear directorio %s: %v", dir, err)
		} else {
			log.Printf("Directorio creado/verificado: %s", dir)
		}
	}
}

// startServer inicia el servidor HTTP
func startServer(router http.Handler) {
	server := &http.Server{
		Addr:         ":" + config.AppConfig.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("🚀 Auth-User Service iniciando en puerto %s", config.AppConfig.Port)
	log.Printf("🌍 Entorno: %s", config.AppConfig.AppEnv)
	log.Printf("🔗 Endpoints disponibles:")
	log.Printf("   - Health Check: http://localhost:%s/api/health", config.AppConfig.Port)
	log.Printf("   - Auth: http://localhost:%s/api/auth/*", config.AppConfig.Port)
	log.Printf("   - Users: http://localhost:%s/api/users/*", config.AppConfig.Port)
	log.Printf("   - Admin: http://localhost:%s/api/admin/*", config.AppConfig.Port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ Error al iniciar servidor: %v", err)
	}
}