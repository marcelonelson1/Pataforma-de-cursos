package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config contiene todas las configuraciones de la aplicación
type Config struct {
	Port        int
	DatabaseURL string
	JWTSecret   string
	JWTExpires  int
}

// Load carga la configuración desde variables de entorno
func Load() (*Config, error) {
	// Intenta cargar variables desde .env si existe
	_ = godotenv.Load()

	port, err := strconv.Atoi(getEnv("PORT", "5000"))
	if err != nil {
		return nil, fmt.Errorf("invalid port: %v", err)
	}

	jwtExpires, err := strconv.Atoi(getEnv("JWT_EXPIRES", "24"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT expiration: %v", err)
	}

	return &Config{
		Port: port,
		DatabaseURL: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			getEnv("DB_USER", "root"),
			getEnv("DB_PASSWORD", ""),
			getEnv("DB_HOST", "localhost"),
			getEnv("DB_PORT", "3306"),
			getEnv("DB_NAME", "plataforma_cursos"),
		),
		JWTSecret:  getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpires: jwtExpires,
	}, nil
}

// getEnv obtiene una variable de entorno o usa un valor predeterminado
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}