package routes

import (
	"home-service/controllers"
	"home-service/middleware"
	"home-service/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura todas las rutas del Home Service
func SetupRoutes() *gin.Engine {
	router := gin.Default()

	// Configurar middlewares globales
	router.Use(middleware.SetupCORS())
	router.Use(middleware.JSONMiddleware())
	router.Use(middleware.RecoveryMiddleware())

	// Inicializar controladores
	homeController := controllers.NewHomeController()

	// Grupo principal de API
	api := router.Group("/api")

	// Rutas publicas del home
	home := api.Group("/home")
	{
		home.GET("/images", homeController.GetAllImages)        // GET /api/home/images
		home.GET("/images/:id", homeController.GetImageByID)    // GET /api/home/images/:id
	}

	// Rutas de administracion (requieren autenticacion + rol admin)
	admin := api.Group("/admin")
	admin.Use(utils.AuthMiddleware(), utils.AdminMiddleware())
	{
		// CRUD de imagenes del home
		admin.GET("/home/images", homeController.GetAllImagesAdmin)           // GET /api/admin/home/images
		admin.POST("/home/images", homeController.CreateImage)               // POST /api/admin/home/images
		admin.PUT("/home/images/:id", homeController.UpdateImage)            // PUT /api/admin/home/images/:id
		admin.DELETE("/home/images/:id", homeController.DeleteImage)         // DELETE /api/admin/home/images/:id
		
		// Gestion de archivos de imagenes
		admin.POST("/home/upload-image", homeController.UploadImage)         // POST /api/admin/home/upload-image
		admin.DELETE("/home/images/:id/file", homeController.DeleteImageFile) // DELETE /api/admin/home/images/:id/file
		
		// Organizacion y estado
		admin.POST("/home/images/reorder", homeController.ReorderImages)     // POST /api/admin/home/images/reorder
		admin.PATCH("/home/images/:id/toggle", homeController.ToggleImageStatus) // PATCH /api/admin/home/images/:id/toggle
		
		// Estadisticas
		admin.GET("/home/stats", homeController.GetHomeStats)                // GET /api/admin/home/stats
	}

	// Rutas estaticas para imagenes
	router.Static("/static/home", "./static/home")

	// Health check
	api.GET("/health", homeController.HealthCheck) // GET /api/health

	// Ruta raiz para verificar que el servicio esta funcionando
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "home-service",
			"status":  "running",
			"version": "1.0.0",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	return router
}