package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	AppEnv     string
	
	// Email configuration
	EmailFrom     string
	EmailPassword string
	SMTPHost      string
	SMTPPort      string
	ContactEmail  string
}

var AppConfig *Config

// LoadConfig carga la configuracion desde variables de entorno
func LoadConfig() {
	// Cargar archivo .env si existe
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontro archivo .env, usando variables de entorno del sistema")
	}

	AppConfig = &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "cursos_db"),
		JWTSecret:  getEnv("JWT_SECRET", "your-secret-key"),
		AppEnv:     getEnv("APP_ENV", "development"),
		
		// Email config
		EmailFrom:     getEnv("EMAIL_FROM", "noreply@example.com"),
		EmailPassword: getEnv("EMAIL_PASSWORD", ""),
		SMTPHost:      getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:      getEnv("SMTP_PORT", "587"),
		ContactEmail:  getEnv("CONTACT_EMAIL", "marcelinho.nelson@gmail.com"),
	}

	log.Printf("Configuracion cargada - Entorno: %s, DB: %s:%s", 
		AppConfig.AppEnv, AppConfig.DBHost, AppConfig.DBPort)
}

// getEnv obtiene una variable de entorno con valor por defecto
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}