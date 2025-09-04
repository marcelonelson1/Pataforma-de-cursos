// config/config.go
package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	// Server
	Port      string
	AppEnv    string
	JWTSecret string

	// Database
	Database DatabaseConfig

	// External Services
	PaymentServiceURL string
	UserServiceURL    string
	FrontendURL       string
	BaseURL           string

	// CORS
	AllowedOrigins []string
	AllowAllOrigins bool

	// File Storage
	UploadPath           string
	MaxVideoSize         int64
	MaxImageSize         int64
	AllowedVideoFormats  []string
	AllowedImageFormats  []string

	// Cache
	Cache CacheConfig

	// CDN
	CDN CDNConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type CacheConfig struct {
	Enabled bool
	TTL     int
}

type CDNConfig struct {
	URL     string
	Enabled bool
}

func Load() *Config {
	return &Config{
		// Server
		Port:      getEnv("PORT", "8003"),
		AppEnv:    getEnv("APP_ENV", "development"),
		JWTSecret: getEnv("JWT_SECRET", "mi_clave_secreta_muy_segura"),

		// Database
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "course_user"),
			Password: getEnv("DB_PASSWORD", "course_password"),
			Name:     getEnv("DB_NAME", "course_db"),
		},

		// External Services
		PaymentServiceURL: getEnv("PAYMENT_SERVICE_URL", "http://localhost:8002"),
		UserServiceURL:    getEnv("USER_SERVICE_URL", "http://localhost:8004"),
		FrontendURL:       getEnv("FRONTEND_URL", "http://localhost:3000"),
		BaseURL:           getEnv("BASE_URL", "http://localhost:8003"),

		// CORS
		AllowedOrigins:  getEnvSlice("ALLOWED_ORIGINS", "http://localhost:3000"),
		AllowAllOrigins: getEnvBool("ALLOW_ALL_ORIGINS", false),

		// File Storage
		UploadPath:          getEnv("UPLOAD_PATH", "./static"),
		MaxVideoSize:        getEnvInt64("MAX_VIDEO_SIZE", 104857600), // 100MB
		MaxImageSize:        getEnvInt64("MAX_IMAGE_SIZE", 5242880),   // 5MB
		AllowedVideoFormats: getEnvSlice("ALLOWED_VIDEO_FORMATS", "mp4,webm,ogg"),
		AllowedImageFormats: getEnvSlice("ALLOWED_IMAGE_FORMATS", "jpg,jpeg,png,webp,gif"),

		// Cache
		Cache: CacheConfig{
			Enabled: getEnvBool("ENABLE_CACHE", true),
			TTL:     getEnvInt("CACHE_TTL", 300),
		},

		// CDN
		CDN: CDNConfig{
			URL:     getEnv("CDN_URL", ""),
			Enabled: getEnvBool("CDN_ENABLED", false),
		},
	}
}

func (d DatabaseConfig) DSN() string {
	return d.User + ":" + d.Password + "@tcp(" + d.Host + ":" + d.Port + ")/" + d.Name + "?charset=utf8mb4&parseTime=True&loc=Local"
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func getEnvSlice(key, defaultValue string) []string {
	value := getEnv(key, defaultValue)
	return strings.Split(value, ",")
}