package main

import (
	"curso-platform/config"
	"curso-platform/routes"
	"curso-platform/utils"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Configuración inicial
	setupLogging()
	loadEnv()
	utils.InitRand()
	utils.InitJWT()

	// Inicializar directorios estáticos
	utils.InitStaticDirs()

	// Configurar base de datos
	db := config.SetupDatabase() // Ahora solo devuelve un valor

	// Forzar ejecución de migraciones independientemente de la variable de entorno
	if err := config.RunMigrations(db); err != nil {
		log.Fatalf("Error al ejecutar migraciones: %v", err)
	}

	// Configurar router
	router := routes.SetupRouter(db)

	// Iniciar servidor
	startServer(router)
}

// setupLogging configura el sistema de logging
func setupLogging() {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		// Configura el logger para escribir tanto en archivo como en consola
		log.SetOutput(logFile)
	} else {
		log.Println("No se pudo abrir el archivo de log, usando salida estándar")
	}
}

// loadEnv carga las variables de entorno desde el archivo .env
func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env, usando variables de entorno del sistema")
	}
}

// startServer inicia el servidor HTTP
func startServer(router http.Handler) {
	port := utils.GetEnv("PORT", "5000")
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Iniciando servidor en puerto %s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error al iniciar servidor: %v", err)
	}
}