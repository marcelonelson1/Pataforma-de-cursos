package main

import (
	"contact-service/config"
	"contact-service/routes"
	"contact-service/utils"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Configurar logging
	setupLogging()

	// Cargar configuracion
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
		log.Println("No se pudo abrir el archivo de log, usando salida estandar")
	}
}

// initStaticDirs crea los directorios estaticos necesarios
func initStaticDirs() {
	createDirIfNotExists("./static")
	createDirIfNotExists("./templates")
}

// createDirIfNotExists crea un directorio si no existe
func createDirIfNotExists(dirPath string) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		log.Printf("Creando directorio: %s", dirPath)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			log.Printf("Error al crear directorio %s: %v", dirPath, err)
		}
	}
}

// startServer inicia el servidor HTTP
func startServer(router *gin.Engine) {
	port := utils.GetEnv("PORT", "8004")
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Contact Service iniciando en puerto %s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error al iniciar servidor: %v", err)
	}
}