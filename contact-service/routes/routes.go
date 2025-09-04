package routes

import (
	"contact-service/controllers"
	"contact-service/middleware"
	"contact-service/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura todas las rutas del Contact Service
func SetupRoutes() *gin.Engine {
	router := gin.Default()

	// Configurar middlewares globales
	router.Use(middleware.SetupCORS())
	router.Use(middleware.JSONMiddleware())
	router.Use(middleware.RecoveryMiddleware())

	// Inicializar controladores
	contactController := controllers.NewContactController()

	// Grupo principal de API
	api := router.Group("/api")

	// Rutas publicas de contacto
	contact := api.Group("/contact")
	{
		contact.POST("", contactController.CreateMessage) // POST /api/contact
	}

	// Rutas de administracion (requieren autenticacion + rol admin)
	admin := api.Group("/admin")
	admin.Use(utils.AuthMiddleware(), utils.AdminMiddleware())
	{
		// Gestion de mensajes
		admin.GET("/messages", contactController.GetAllMessages)                    // GET /api/admin/messages
		admin.GET("/messages/:id", contactController.GetMessageByID)               // GET /api/admin/messages/:id
		admin.PATCH("/messages/:id/:action", contactController.UpdateMessageStatus) // PATCH /api/admin/messages/:id/read|star
		admin.DELETE("/messages/:id", contactController.DeleteMessage)             // DELETE /api/admin/messages/:id
		admin.POST("/messages/:id/reply", contactController.ReplyToMessage)        // POST /api/admin/messages/:id/reply
		
		// Estadisticas y filtros
		admin.GET("/messages/stats", contactController.GetContactStats)      // GET /api/admin/messages/stats
		admin.GET("/messages/starred", contactController.GetStarredMessages) // GET /api/admin/messages/starred
	}

	// Health check
	api.GET("/health", contactController.HealthCheck) // GET /api/health

	// Ruta raiz para verificar que el servicio esta funcionando
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "contact-service",
			"status":  "running",
			"version": "1.0.0",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	return router
}