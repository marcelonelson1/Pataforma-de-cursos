package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	DBHost   string
	DBPort   string
	DBUser   string
	DBPass   string
	DBName   string
	JWTSecret string
	AppEnv   string
	CORSOrigins string
	
	// URLs de otros servicios
	AuthServiceURL      string
	PaymentServiceURL   string
	CourseServiceURL    string
	ContactServiceURL   string
	PortfolioServiceURL string
	HomeServiceURL      string
	FrontendURL         string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	return &Config{
		Port:     getEnv("PORT", "8007"),
		DBHost:   getEnv("DB_HOST", "localhost"),
		DBPort:   getEnv("DB_PORT", "3306"),
		DBUser:   getEnv("DB_USER", "root"),
		DBPass:   getEnv("DB_PASSWORD", "root"),
		DBName:   getEnv("DB_NAME", "analytics_db"),
		JWTSecret: getEnv("JWT_SECRET", "yoyo"),
		AppEnv:   getEnv("APP_ENV", "development"),
		CORSOrigins: getEnv("CORS_ORIGINS", "http://localhost:3000"),
		
		AuthServiceURL:      getEnv("AUTH_SERVICE_URL", "http://localhost:8001"),
		PaymentServiceURL:   getEnv("PAYMENT_SERVICE_URL", "http://localhost:8002"),
		CourseServiceURL:    getEnv("COURSE_SERVICE_URL", "http://localhost:8003"),
		ContactServiceURL:   getEnv("CONTACT_SERVICE_URL", "http://localhost:8004"),
		PortfolioServiceURL: getEnv("PORTFOLIO_SERVICE_URL", "http://localhost:8005"),
		HomeServiceURL:      getEnv("HOME_SERVICE_URL", "http://localhost:8006"),
		FrontendURL:         getEnv("FRONTEND_URL", "http://localhost:3000"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}