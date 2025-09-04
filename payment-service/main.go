// main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"payment-service/config"
	"payment-service/routes"
	"payment-service/services"
	"payment-service/utils"
)

func main() {
	log.Printf("🚀 [STARTUP] Iniciando Payment Service...")

	// Cargar variables de entorno desde .env
	log.Printf("🔧 [CONFIG] Cargando variables de entorno...")
	if err := godotenv.Load(); err != nil {
		log.Printf("⚠️ [CONFIG] No se encontró archivo .env: %v", err)
		log.Println("Usando variables del sistema")
	} else {
		log.Printf("✅ [CONFIG] Archivo .env cargado exitosamente")
	}

	// Cargar configuración
	log.Printf("🔧 [CONFIG] Cargando configuración...")
	cfg := config.Load()
	
	// Debug: Verificar configuración crítica
	log.Printf("🔍 [CONFIG] Verificando configuración:")
	log.Printf("   Port: %s", cfg.Port)
	log.Printf("   App Environment: %s", cfg.AppEnv)
	log.Printf("   Database Host: %s", cfg.Database.Host)
	log.Printf("   JWT Secret configurado: %t (longitud: %d)", cfg.JWTSecret != "", len(cfg.JWTSecret))
	log.Printf("   PayPal Client ID: %s... (length: %d)", maskString(cfg.PayPal.ClientID), len(cfg.PayPal.ClientID))
	log.Printf("   PayPal Secret configurado: %t", cfg.PayPal.Secret != "")
	log.Printf("   PayPal Environment: %s", cfg.PayPal.Env)
	
	// 🔥 NUEVO: Debug de Mercado Pago
	log.Printf("   MercadoPago Access Token: %s... (length: %d)", maskString(cfg.MercadoPago.AccessToken), len(cfg.MercadoPago.AccessToken))
	log.Printf("   MercadoPago Environment: %s", cfg.MercadoPago.Environment)
	log.Printf("   MercadoPago Accept USD: %t", cfg.MercadoPago.AcceptUSD)
	
	log.Printf("   User Service URL: %s", cfg.UserServiceURL)
	log.Printf("   Course Service URL: %s", cfg.CourseServiceURL)
	log.Printf("   Frontend URL: %s", cfg.FrontendURL)
	log.Printf("   Base URL: %s", cfg.BaseURL)

	// 🔥 VERIFICAR JWT SECRET CRÍTICO
	if cfg.JWTSecret == "" {
		log.Fatal("❌ ERROR CRÍTICO: JWT_SECRET no está configurado en .env")
	}
	log.Printf("✅ [CONFIG] JWT Secret configurado correctamente")

	// Verificar credenciales PayPal
	if cfg.PayPal.ClientID == "" || cfg.PayPal.Secret == "" {
		log.Printf("❌ [ERROR] Credenciales PayPal no configuradas")
		log.Printf("💡 [HELP] Por favor configura PAYPAL_CLIENT_ID y PAYPAL_SECRET en tu archivo .env")
		os.Exit(1)
	}
	log.Printf("✅ [CONFIG] Credenciales PayPal configuradas correctamente")

	// 🔥 NUEVO: Verificar credenciales de Mercado Pago
	if cfg.MercadoPago.AccessToken == "" {
		log.Printf("⚠️ [WARNING] Credenciales de Mercado Pago no configuradas")
		log.Printf("💡 [HELP] Para habilitar Mercado Pago, configura MERCADOPAGO_ACCESS_TOKEN en tu archivo .env")
		log.Printf("💡 [HELP] Obtén tu Access Token en: https://www.mercadopago.com.ar/developers/panel")
	} else {
		log.Printf("✅ [CONFIG] Credenciales de Mercado Pago configuradas correctamente")
		
		// Test de conectividad con Mercado Pago
		if err := services.TestMercadoPagoConnection(cfg); err != nil {
			log.Printf("⚠️ [WARNING] Problema con conectividad de Mercado Pago: %v", err)
		} else {
			log.Printf("🎉 [SUCCESS] Mercado Pago conectado y funcionando correctamente")
		}
	}

	// Conectar a la base de datos
	log.Printf("🔧 [DATABASE] Conectando a la base de datos...")
	db := utils.ConnectDatabase(cfg)
	log.Printf("✅ [DATABASE] Conexión a base de datos establecida")

	// Ejecutar migraciones
	log.Printf("🔧 [DATABASE] Ejecutando migraciones...")
	utils.MigrateDatabase(db)
	log.Printf("✅ [DATABASE] Migraciones completadas exitosamente")

	// Inicializar proveedores de pago
	log.Printf("🔧 [PROVIDERS] Inicializando proveedores de pago...")
	services.InitPaymentProviders(cfg)
	log.Printf("✅ [PROVIDERS] Proveedores de pago inicializados")

	// Configurar router
	log.Printf("🔧 [ROUTER] Configurando router...")
	router := setupRouter(cfg)
	log.Printf("✅ [ROUTER] Router configurado")

	// Configurar rutas CON MIDDLEWARE DE AUTENTICACIÓN
	log.Printf("🔧 [ROUTES] Configurando rutas...")
	routes.SetupPaymentRoutes(router, db, cfg)
	log.Printf("✅ [ROUTES] Rutas configuradas con middleware de autenticación")

	// 🔥 NUEVO: Inicializar auto-updater automático
	if cfg.MercadoPago.AccessToken != "" {
		log.Printf("🤖 [AUTO_UPDATER] Iniciando actualizador automático de pagos...")
		autoUpdater := services.NewAutoUpdater(db, cfg)
		autoUpdater.Start()
		log.Printf("✅ [AUTO_UPDATER] Actualizador automático iniciado")
	}

	// 🔥 NUEVO: Inicializar servicio de limpieza automática de pagos expirados
	log.Printf("🧹 [CLEANUP] Iniciando servicio de limpieza automática...")
	cleanupService := services.NewCleanupService(db, cfg)
	cleanupService.StartCleanupScheduler()
	log.Printf("✅ [CLEANUP] Servicio de limpieza automática iniciado")

	// Iniciar servidor
	startServer(router, cfg)
}

func setupRouter(cfg *config.Config) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
		log.Printf("🔧 [GIN] Modo establecido: Release")
	} else {
		log.Printf("🔧 [GIN] Modo establecido: Debug")
	}

	router := gin.Default()

	// Configurar CORS mejorado
	log.Printf("🔧 [CORS] Configurando CORS...")
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.FrontendURL, "http://localhost:3000", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Middleware de JSON
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	})

	// Middleware de recovery
	router.Use(gin.Recovery())

	// Middleware de logging mejorado
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("🌐 %s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format("02/Jan/2006:15:04:05 -0700"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	log.Printf("✅ [CORS] CORS y middlewares configurados")
	return router
}

// Helper function para ocultar credenciales sensibles en logs
func maskString(s string) string {
	if len(s) <= 8 {
		return "****"
	}
	return s[:4] + "****"
}

func startServer(router *gin.Engine, cfg *config.Config) {
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("🚀 Payment Service iniciado en puerto %s", cfg.Port)
	log.Printf("📊 Ambiente: %s", cfg.AppEnv)
	log.Printf("💳 PayPal Env: %s", cfg.PayPal.Env)
	
	// 🔥 NUEVO: Info de Mercado Pago
	if cfg.MercadoPago.AccessToken != "" {
		log.Printf("💰 MercadoPago Env: %s", cfg.MercadoPago.Environment)
		log.Printf("💵 MercadoPago USD: %t", cfg.MercadoPago.AcceptUSD)
	}
	
	log.Printf("🔐 Autenticación JWT: HABILITADA")
	log.Printf("🔗 Endpoints principales:")
	log.Printf("   • Health Check: http://localhost:%s/api/health", cfg.Port)
	log.Printf("   • Crear Pago: POST http://localhost:%s/api/pagos [AUTH REQUIRED]", cfg.Port)
	log.Printf("   • Verificar Pago: GET http://localhost:%s/api/pagos/:id [AUTH REQUIRED]", cfg.Port)
	log.Printf("   • Webhook MercadoPago: POST http://localhost:%s/api/pagos/mercadopago/webhook [PUBLIC]", cfg.Port)
	
	if cfg.AppEnv == "development" {
		log.Printf("🔧 [DEV] Endpoints de diagnóstico:")
		log.Printf("   • Ver Rutas: http://localhost:%s/routes", cfg.Port)
		log.Printf("   • PayPal Debug: http://localhost:%s/api/debug/paypal", cfg.Port)
		log.Printf("   • MercadoPago Health: http://localhost:%s/api/debug/mercadopago", cfg.Port)
		log.Printf("   • MercadoPago Diagnostic: http://localhost:%s/api/debug/mercadopago/diagnostic", cfg.Port)
		log.Printf("   • MercadoPago Deep Diagnostic: http://localhost:%s/api/debug/mercadopago/deep", cfg.Port)
		log.Printf("")
		log.Printf("🔍 [SOLUCIÓN] Para diagnosticar tu problema:")
		log.Printf("   1. Ejecuta: curl http://localhost:%s/api/debug/mercadopago/deep", cfg.Port)
		log.Printf("   2. O visita en tu navegador: http://localhost:%s/api/debug/mercadopago/deep", cfg.Port)
		log.Printf("   3. Revisa el campo 'recommendations' en la respuesta")
		log.Printf("")
	}
	
	log.Printf("💡 Métodos de pago disponibles:")
	log.Printf("   • PayPal: %s", getStatus(cfg.PayPal.ClientID != ""))
	log.Printf("   • Coinbase: %s", getStatus(cfg.Coinbase.APIKey != ""))
	log.Printf("   • MercadoPago: %s", getStatus(cfg.MercadoPago.AccessToken != ""))
	log.Printf("   • Desarrollo: ✅ ACTIVO")
	
	if cfg.MercadoPago.AccessToken != "" {
		log.Printf("")
		log.Printf("🎯 [MERCADOPAGO] Para resolver el error 'Una de las partes es de prueba':")
		log.Printf("   1. Ve a: https://www.mercadopago.com.ar/developers/panel")
		log.Printf("   2. Asegúrate de usar cuentas de PRUEBA para pagar")
		log.Printf("   3. Usa la cuenta COMPRADOR de prueba, no tu cuenta real")
		log.Printf("   4. Usa tarjetas de prueba: 4509 9535 6623 3704")
		log.Printf("   5. Ejecuta el diagnóstico completo en: /api/debug/mercadopago/deep")
		log.Printf("")
	}
	
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ Error al iniciar servidor: %v", err)
	}
}

// Helper function para mostrar estado de configuración
func getStatus(configured bool) string {
	if configured {
		return "✅ CONFIGURADO"
	}
	return "❌ NO CONFIGURADO"
}