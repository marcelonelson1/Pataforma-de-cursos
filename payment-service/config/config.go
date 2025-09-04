// config/config.go
package config

import (
	"log"
	"os"
)

type Config struct {
	// Server
	Port      string
	AppEnv    string
	JWTSecret string

	// Database
	Database DatabaseConfig

	// External Services
	UserServiceURL   string
	CourseServiceURL string
	FrontendURL      string
	BaseURL          string

	// Payment Providers
	PayPal      PayPalConfig
	Coinbase    CoinbaseConfig
	MercadoPago MercadoPagoConfig // 🔥 AGREGADO
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type PayPalConfig struct {
	ClientID string
	Secret   string
	Env      string // sandbox or live
}

type CoinbaseConfig struct {
	APIKey string
}

// 🔥 NUEVA ESTRUCTURA para Mercado Pago
type MercadoPagoConfig struct {
	AccessToken string
	Environment string // sandbox o production
	AcceptUSD   bool   // si acepta pagos en USD
}

func Load() *Config {
	cfg := &Config{
		// Server
		Port:      getEnv("PORT", "8002"),
		AppEnv:    getEnv("APP_ENV", "development"),
		JWTSecret: getEnv("JWT_SECRET", "mi_clave_secreta_muy_segura"),

		// Database
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "payment_db"),
		},

		// External Services
		UserServiceURL:   getEnv("USER_SERVICE_URL", "http://localhost:8004"),
		CourseServiceURL: getEnv("COURSE_SERVICE_URL", "http://localhost:8003"),
		FrontendURL:      getEnv("FRONTEND_URL", "http://localhost:3000"),
		BaseURL:          getEnv("BASE_URL", "http://localhost:8002"),

		// Payment Providers - CREDENCIALES ACTUALIZADAS
		PayPal: PayPalConfig{
			ClientID: getEnv("PAYPAL_CLIENT_ID", "ASYN839bjb4gjMr6nRCc-7YYR8HutdM48kFMWhq-Sxp-PgB5c5R38yGiLBEPwDBIptFj8IJ71OPVXVUt"),
			Secret:   getEnv("PAYPAL_SECRET", "EHGs6eflLFMvOWTrhWjCWmwBmbXIL8he0dM6bIbcVDFhwGStuz3PGFp_nreODGiJueoNyjxfZG1Hqi0-"),
			Env:      getEnv("PAYPAL_ENV", "sandbox"),
		},

		Coinbase: CoinbaseConfig{
			APIKey: getEnv("COINBASE_COMMERCE_API_KEY", ""),
		},

		// 🔥 NUEVA CONFIGURACIÓN para Mercado Pago
		MercadoPago: MercadoPagoConfig{
			AccessToken: getEnv("MERCADOPAGO_ACCESS_TOKEN", ""),
			Environment: getEnv("MERCADOPAGO_ENV", "sandbox"),
			AcceptUSD:   getEnv("MERCADOPAGO_ACCEPT_USD", "false") == "true",
		},
	}

	// 🔍 Debug: Verificar que las credenciales se cargaron correctamente
	log.Printf("🔧 PayPal Client ID: %s... (length: %d)",
		maskString(cfg.PayPal.ClientID), len(cfg.PayPal.ClientID))
	log.Printf("🔧 PayPal Secret: %s... (length: %d)",
		maskString(cfg.PayPal.Secret), len(cfg.PayPal.Secret))
	log.Printf("🔧 PayPal Environment: %s", cfg.PayPal.Env)

	// 🔥 NUEVO: Debug para Mercado Pago
	log.Printf("🔧 MercadoPago Access Token: %s... (length: %d)",
		maskString(cfg.MercadoPago.AccessToken), len(cfg.MercadoPago.AccessToken))
	log.Printf("🔧 MercadoPago Environment: %s", cfg.MercadoPago.Environment)
	log.Printf("🔧 MercadoPago Accept USD: %t", cfg.MercadoPago.AcceptUSD)

	// Validar que las credenciales no estén vacías
	if cfg.PayPal.ClientID == "" {
		log.Printf("⚠️  PayPal Client ID está vacío")
	}
	if cfg.PayPal.Secret == "" {
		log.Printf("⚠️  PayPal Secret está vacío")
	}

	// 🔥 NUEVA VALIDACIÓN para Mercado Pago
	if cfg.MercadoPago.AccessToken == "" {
		log.Printf("⚠️  Mercado Pago Access Token está vacío")
	}

	return cfg
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

// Helper function para ocultar credenciales sensibles en logs
func maskString(s string) string {
	if len(s) <= 8 {
		return "****"
	}
	return s[:4] + "****"
}
