package routes

import (
	"portfolio-service/controllers"
	"portfolio-service/middleware"
	"portfolio-service/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura todas las rutas del Portfolio Service
func SetupRoutes() *gin.Engine {
	router := gin.Default()

	// Configurar middlewares globales
	router.Use(middleware.SetupCORS())
	router.Use(middleware.JSONMiddleware())
	router.Use(middleware.RecoveryMiddleware())

	// Inicializar controladores
	portfolioController := controllers.NewPortfolioController()

	// Grupo principal de API
	api := router.Group("/api")

	// Rutas publicas del portafolio
	portfolio := api.Group("/portfolio")
	{
		portfolio.GET("", portfolioController.GetAllProjects)                        // GET /api/portfolio
		portfolio.GET("/:id", portfolioController.GetProjectByID)                   // GET /api/portfolio/:id
		portfolio.GET("/category/:category", portfolioController.GetProjectsByCategory) // GET /api/portfolio/category/:category
		portfolio.GET("/categories", portfolioController.GetCategories)             // GET /api/portfolio/categories
	}

	// Rutas de administracion (requieren autenticacion + rol admin)
	admin := api.Group("/admin")
	admin.Use(utils.AuthMiddleware(), utils.AdminMiddleware())
	{
		// CRUD de proyectos
		admin.GET("/portfolio", portfolioController.GetAllProjectsAdmin)           // GET /api/admin/portfolio
		admin.POST("/portfolio", portfolioController.CreateProject)               // POST /api/admin/portfolio
		admin.PUT("/portfolio/:id", portfolioController.UpdateProject)            // PUT /api/admin/portfolio/:id
		admin.DELETE("/portfolio/:id", portfolioController.DeleteProject)         // DELETE /api/admin/portfolio/:id
		
		// Gestion de imagenes
		admin.POST("/portfolio/upload-image", portfolioController.UploadProjectImage)   // POST /api/admin/portfolio/upload-image
		admin.DELETE("/portfolio/:id/image", portfolioController.DeleteProjectImage)   // DELETE /api/admin/portfolio/:id/image
		
		// Organizacion y estado
		admin.POST("/portfolio/reorder", portfolioController.ReorderProjects)     // POST /api/admin/portfolio/reorder
		admin.PATCH("/portfolio/:id/toggle", portfolioController.ToggleProjectStatus) // PATCH /api/admin/portfolio/:id/toggle
		
		// Estadisticas
		admin.GET("/portfolio/stats", portfolioController.GetPortfolioStats)      // GET /api/admin/portfolio/stats
	}

	// Rutas estaticas para imagenes con CORS
	staticGroup := router.Group("/static")
	staticGroup.Use(middleware.SetupCORS())
	staticGroup.Static("/portfolio", "./static/portfolio")

	// Health check
	api.GET("/health", portfolioController.HealthCheck) // GET /api/health

	// Ruta raiz para verificar que el servicio esta funcionando
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "portfolio-service",
			"status":  "running",
			"version": "1.0.0",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	return router
}