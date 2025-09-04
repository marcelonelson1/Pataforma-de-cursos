// routes/routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"payment-service/config"
	"payment-service/controllers"
	"payment-service/middleware"
)

// SetupPaymentRoutes configura todas las rutas del Payment Service
func SetupPaymentRoutes(router *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// Inicializar controladores
	paymentController := controllers.NewPaymentController(db, cfg)
	webhookController := controllers.NewWebhookController(db, cfg)
	paymentReturnController := controllers.NewPaymentReturnController(db, cfg) // 🆕 Controlador genérico de retornos

	// Grupo principal de API
	api := router.Group("/api")

	// Rutas de pagos (requieren autenticación)
	pagos := api.Group("/pagos")
	pagos.Use(middleware.AuthMiddleware(cfg))
	{
		pagos.POST("", paymentController.CreatePayment)               // POST /api/pagos
		pagos.GET("/:id", paymentController.GetPaymentStatus)        // GET /api/pagos/:id
		pagos.GET("/user/history", paymentController.GetUserPayments) // GET /api/pagos/user/history
		pagos.PUT("/:id/status", paymentController.UpdatePaymentStatus) // PUT /api/pagos/:id/status
		pagos.POST("/:id/refresh", paymentController.RefreshPaymentStatus) // POST /api/pagos/:id/refresh - Actualizar estado manualmente
	}

	// Rutas públicas para comunicación entre microservicios
	publicPagos := api.Group("/pagos")
	{
		publicPagos.GET("/verify-course-access", paymentController.VerifyCourseAccess)      // GET /api/pagos/verify-course-access
		publicPagos.GET("/course-access-info", paymentController.GetCourseAccessInfo)      // GET /api/pagos/course-access-info
	}

	// Rutas de webhooks (NO requieren autenticación)
	webhooks := api.Group("/pagos")
	{
		// Webhook genérico
		webhooks.POST("/webhook", webhookController.GenericWebhook) // POST /api/pagos/webhook
		
		// Webhooks específicos por proveedor
		webhooks.POST("/paypal/webhook", webhookController.PayPalWebhook)     // POST /api/pagos/paypal/webhook
		webhooks.POST("/coinbase/webhook", webhookController.CoinbaseWebhook) // POST /api/pagos/coinbase/webhook
		webhooks.POST("/mercadopago/webhook", webhookController.HandleMercadoPagoWebhookNew) // POST /api/pagos/mercadopago/webhook
		
		// 🆕 NUEVO: Controlador de retorno genérico (compatible con todos los proveedores)
		webhooks.GET("/return", paymentReturnController.HandleGenericReturn) // GET /api/pagos/return
		
		// ✅ EXISTENTES: Callbacks específicos (mantener compatibilidad)
		webhooks.GET("/paypal/callback", webhookController.PayPalCallback)    // GET /api/pagos/paypal/callback
		webhooks.GET("/mercadopago/return", paymentController.HandleMercadoPagoReturn) // GET /api/pagos/mercadopago/return
	}

	// Rutas de administración (requieren autenticación + rol admin)
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(cfg), middleware.AdminMiddleware())
	{
		// Rutas completas de administración
		admin.GET("/payments", paymentController.GetAllPayments)       // GET /api/admin/payments
		admin.GET("/payments/stats", paymentController.GetPaymentStats) // GET /api/admin/payments/stats
		
		// Testing de proveedores de pago
		admin.POST("/test/mercadopago", paymentController.TestMercadoPago) // POST /api/admin/test/mercadopago
		
		// 🔥 NUEVO: Diagnósticos avanzados (solo admin)
		admin.GET("/diagnostic/mercadopago", paymentController.DiagnosticMercadoPago) // GET /api/admin/diagnostic/mercadopago
		admin.GET("/diagnostic/mercadopago/deep", paymentController.DeepDiagnosticMercadoPago) // GET /api/admin/diagnostic/mercadopago/deep
		
		// 🔥 NUEVO: Gestión de expiración de pagos (solo admin)
		admin.POST("/cleanup/expired", paymentController.CleanupExpiredPayments)      // POST /api/admin/cleanup/expired
		admin.GET("/cleanup/stats", paymentController.GetCleanupStats)                // GET /api/admin/cleanup/stats
		admin.POST("/cleanup/start", paymentController.StartCleanupScheduler)         // POST /api/admin/cleanup/start
		admin.POST("/cleanup/stop", paymentController.StopCleanupScheduler)           // POST /api/admin/cleanup/stop
		admin.POST("/payments/:id/expire", paymentController.ExpirePayment)           // POST /api/admin/payments/:id/expire
		admin.POST("/payments/:id/extend", paymentController.ExtendPaymentExpiration) // POST /api/admin/payments/:id/extend
	}

	// Rutas de salud y métricas (públicas)
	health := api.Group("/health")
	{
		health.GET("", paymentController.HealthCheck)                // GET /api/health
		health.GET("/paypal", paymentController.PayPalHealthCheck)   // GET /api/health/paypal
		health.GET("/mercadopago", paymentController.MercadoPagoHealthCheck) // GET /api/health/mercadopago
	}

	// 🔥 NUEVO: Rutas de debug para desarrollo (públicas en dev)
	if cfg.AppEnv == "development" {
		debug := api.Group("/debug")
		{
			debug.GET("/paypal", paymentController.PayPalHealthCheck)
			debug.GET("/mercadopago", paymentController.MercadoPagoHealthCheck)
			debug.GET("/mercadopago/diagnostic", paymentController.DiagnosticMercadoPago) // Versión pública para dev
			debug.GET("/mercadopago/deep", paymentController.DeepDiagnosticMercadoPago) // Versión pública para dev
		}
	}

	// Ruta raíz para verificar que el servicio está funcionando
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "payment-service",
			"status":  "running",
			"version": "1.0.0",
		})
	})

	// Ruta para listar todas las rutas disponibles (útil para desarrollo)
	if cfg.AppEnv == "development" {
		router.GET("/routes", func(c *gin.Context) {
			routes := []string{
				"GET  /",
				"GET  /routes                                    [DEV ONLY]",
				"",
				"--- HEALTH CHECKS ---",
				"GET  /api/health                                [PUBLIC]",
				"GET  /api/health/paypal                         [PUBLIC]",
				"GET  /api/health/mercadopago                    [PUBLIC]",
				"",
				"--- PAYMENT OPERATIONS ---",
				"POST /api/pagos                                 [AUTH REQUIRED]",
				"GET  /api/pagos/:id                             [AUTH REQUIRED]",
				"GET  /api/pagos/user/history                    [AUTH REQUIRED]",
				"PUT  /api/pagos/:id/status                      [AUTH REQUIRED]",
				"POST /api/pagos/:id/refresh                     [AUTH REQUIRED] - Refresh payment status",
				"",
				"--- MICROSERVICE COMMUNICATION ---",
				"GET  /api/pagos/verify-course-access            [PUBLIC] - Check user course access",
				"GET  /api/pagos/course-access-info              [PUBLIC] - Get detailed course access info",
				"",
				"--- WEBHOOKS ---",
				"POST /api/pagos/webhook                         [PUBLIC]",
				"POST /api/pagos/paypal/webhook                  [PUBLIC]",
				"POST /api/pagos/coinbase/webhook                [PUBLIC]",
				"POST /api/pagos/mercadopago/webhook             [PUBLIC]",
				"GET  /api/pagos/paypal/callback                 [PUBLIC]",
				"",
				"--- ADMIN OPERATIONS ---",
				"GET  /api/admin/payments                        [ADMIN REQUIRED]",
				"GET  /api/admin/payments/stats                  [ADMIN REQUIRED]",
				"POST /api/admin/test/mercadopago                [ADMIN REQUIRED]",
				"GET  /api/admin/diagnostic/mercadopago          [ADMIN REQUIRED]",
				"GET  /api/admin/diagnostic/mercadopago/deep     [ADMIN REQUIRED]",
				"",
				"--- DEBUG (DEV ONLY) ---",
				"GET  /api/debug/paypal                          [DEV ONLY]",
				"GET  /api/debug/mercadopago                     [DEV ONLY]",
				"GET  /api/debug/mercadopago/diagnostic          [DEV ONLY]",
				"GET  /api/debug/mercadopago/deep                [DEV ONLY]",
			}
			c.JSON(200, gin.H{
				"service": "payment-service",
				"routes":  routes,
				"auth_middleware": "ENABLED",
				"jwt_secret_configured": cfg.JWTSecret != "",
				"payment_methods_supported": []string{
					"dev",
					"paypal", 
					"coinbase", 
					"mercadopago",
					"stripe", 
					"card", 
					"transfer",
				},
				"debug_endpoints": map[string]string{
					"simple_diagnostic":    "/api/debug/mercadopago/diagnostic",
					"deep_diagnostic":      "/api/debug/mercadopago/deep",
					"mercadopago_health":   "/api/health/mercadopago",
					"all_routes":           "/routes",
				},
			})
		})
	}
}