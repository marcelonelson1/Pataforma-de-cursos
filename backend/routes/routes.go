package routes

import (
	"curso-platform/controllers"
	"curso-platform/middleware"
	"curso-platform/services"
	"curso-platform/utils"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter configura el router con todas las rutas y middleware
func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// Configurar CORS
	setupCORS(router)

	// Middleware global
	router.Use(middleware.JSONMiddleware())
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.ActivityLogger())

	// Iniciar servicios
	// Nota: Ahora cada servicio recibe directamente la conexión de base de datos
	authService := services.NewAuthService()
	passwordResetService := services.NewPasswordResetService(db) // Nuevo servicio MVC
	userService := services.NewUserService()
	statsService := services.NewStatsService()
	cursoService := services.NewCursoService()
	capituloService := services.NewCapituloService()
	progresoService := services.NewProgresoService()
	contactService := services.NewContactService()
	pagoService := services.NewPagoService()
	portfolioService := services.NewPortfolioService()
	homeImageService := services.NewHomeImageService()
	activityService := services.NewActivityService(db)
	enrollmentService := services.NewEnrollmentService(db)

	// Iniciar controladores
	// Modificamos el constructor de authController para pasar el passwordResetService
	authController := controllers.NewAuthController(authService, passwordResetService)
	profileController := controllers.NewProfileController(userService, authService)
	adminController := controllers.NewAdminController(userService, statsService, authService)
	cursoController := controllers.NewCursoController(cursoService)
	capituloController := controllers.NewCapituloController(capituloService)
	progresoController := controllers.NewProgresoController(progresoService)
	contactController := controllers.NewContactController(contactService)
	pagoController := controllers.NewPagoController(pagoService)
	portfolioController := controllers.NewPortfolioController(portfolioService)
	homeImageController := controllers.NewHomeImageController(homeImageService)
	activityController := controllers.NewActivityController(activityService)
	enrollmentController := controllers.NewEnrollmentController(enrollmentService)

	// Registrar rutas
	authController.RegisterRoutes(router)
	profileController.RegisterRoutes(router)
	adminController.RegisterRoutes(router)
	cursoController.RegisterRoutes(router)
	capituloController.RegisterRoutes(router)
	progresoController.RegisterRoutes(router)
	contactController.RegisterRoutes(router)
	pagoController.RegisterRoutes(router)
	portfolioController.RegisterRoutes(router)
	homeImageController.RegisterRoutes(router)
	activityController.RegisterRoutes(router)
	enrollmentController.RegisterRoutes(router)

	// Configurar rutas para archivos estáticos
	setupStaticRoutes(router)

	// Ruta de verificación de salud
	router.GET("/api/health", healthCheck)

	return router
}

// setupCORS configura la política CORS para el router
func setupCORS(router *gin.Engine) {
	frontendURL := utils.GetEnv("FRONTEND_URL", "http://localhost:3000")
	log.Printf("Configurando CORS para permitir origen: %s", frontendURL)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}

// setupStaticRoutes configura las rutas para servir archivos estáticos
func setupStaticRoutes(router *gin.Engine) {
	// Servir archivos estáticos de imágenes
	router.GET("/static/images/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		filepath := "./static/images/" + filename
		
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Imagen no encontrada"})
			return
		}
		
		c.Header("Cache-Control", "public, max-age=31536000")
		c.File(filepath)
	})
	
	// Servir archivos estáticos de perfiles
	router.GET("/static/profiles/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		c.File("./static/profiles/" + filename)
	})
	
	// Servir archivos estáticos de portafolio
	router.GET("/static/portfolio/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		c.File("./static/portfolio/" + filename)
	})

	// Servir archivos estáticos de home
	router.GET("/static/home/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		c.File("./static/home/" + filename)
	})

	// Servir archivos estáticos de actividades
	router.GET("/static/activities/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		filepath := "./static/activities/" + filename
		
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Imagen de actividad no encontrada"})
			return
		}
		
		c.Header("Cache-Control", "public, max-age=31536000")
		c.File(filepath)
	})
}

// healthCheck es un endpoint para verificar el estado de la aplicación
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"version": "1.0.0",
		"time":    time.Now().Format(time.RFC3339),
	})
}