package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Server
	Port   string
	AppEnv string

	// JWT
	JWTSecret string

	// External Services
	PaymentServiceURL string
	CourseServiceURL  string
	FrontendURL       string

	// File Storage
	UploadPath           string
	MaxImageSize         string
	AllowedImageFormats  string

	// Email
	EmailFrom    string
	EmailPassword string
	SMTPHost     string
	SMTPPort     string
	MockEmail    string

	// CORS
	CORSOrigins string
}

var AppConfig *Config

func LoadConfig() {
	// Cargar .env si existe
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env, usando variables de entorno del sistema")
	}

	AppConfig = &Config{
		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "auth_user"),
		DBPassword: getEnv("DB_PASSWORD", "auth_password"),
		DBName:     getEnv("DB_NAME", "auth_user_db"),

		// Server
		Port:   getEnv("PORT", "8001"),
		AppEnv: getEnv("APP_ENV", "development"),

		// JWT
		JWTSecret: getEnv("JWT_SECRET", "mi_clave_secreta_muy_segura"),

		// External Services
		PaymentServiceURL: getEnv("PAYMENT_SERVICE_URL", "http://localhost:8002"),
		CourseServiceURL:  getEnv("COURSE_SERVICE_URL", "http://localhost:8003"),
		FrontendURL:       getEnv("FRONTEND_URL", "http://localhost:3000"),

		// File Storage
		UploadPath:          getEnv("UPLOAD_PATH", "./static"),
		MaxImageSize:        getEnv("MAX_IMAGE_SIZE", "5MB"),
		AllowedImageFormats: getEnv("ALLOWED_IMAGE_FORMATS", "jpg,jpeg,png,webp"),

		// Email
		EmailFrom:     getEnv("EMAIL_FROM", "noreply@example.com"),
		EmailPassword: getEnv("EMAIL_PASSWORD", ""),
		SMTPHost:      getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:      getEnv("SMTP_PORT", "587"),
		MockEmail:     getEnv("MOCK_EMAIL", "true"),

		// CORS
		CORSOrigins: getEnv("CORS_ORIGINS", "http://localhost:3000"),
	}

	log.Printf("Configuración cargada - Puerto: %s, Entorno: %s", AppConfig.Port, AppConfig.AppEnv)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}