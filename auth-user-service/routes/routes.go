package routes

import (
	"auth-user-service/controllers"
	"auth-user-service/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura todas las rutas de la aplicación
func SetupRoutes() *gin.Engine {
	router := gin.Default()

	// Configurar middlewares globales
	router.Use(middleware.SetupCORS())
	router.Use(middleware.JSONMiddleware())
	router.Use(middleware.RecoveryMiddleware())

	// Inicializar controladores
	authController := controllers.NewAuthController()
	profileController := controllers.NewProfileController()
	userController := controllers.NewUserController()
	adminController := controllers.NewAdminController()
	fileController := controllers.NewFileController()

	// Health check
	router.GET("/api/health", healthCheck)

	// Rutas de autenticación (públicas)
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
		auth.POST("/forgot-password", authController.ForgotPassword)
		auth.GET("/reset-password/:token/validate", authController.ValidateResetToken)
		auth.POST("/reset-password", authController.ResetPassword)
		auth.POST("/validate-token", authController.ValidateToken) // Para otros servicios

		// Rutas protegidas de autenticación
		authProtected := auth.Group("")
		authProtected.Use(middleware.AuthMiddleware())
		{
			authProtected.GET("/check-admin", authController.CheckAdmin)
			authProtected.POST("/refresh-token", authController.RefreshToken)
			authProtected.POST("/change-password", authController.ChangePassword)
		}
	}

	// Rutas de usuario/perfil (requieren autenticación)
	users := router.Group("/api/users")
	users.Use(middleware.AuthMiddleware())
	{
		users.GET("/profile", profileController.GetProfile)
		users.PUT("/profile", profileController.UpdateProfile)
		users.POST("/profile/image", fileController.UploadProfileImage)
		users.GET("/notification-settings", profileController.GetNotificationSettings)
		users.PUT("/notification-settings", profileController.UpdateNotificationSettings)

		// Endpoint para otros servicios (obtener usuario por ID)
		users.GET("/:id", userController.GetUserByID)
		users.GET("", userController.GetUsers) // Lista con paginación
	}

	// Rutas de administración (requieren autenticación + rol admin)
	admin := router.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		admin.GET("/stats", adminController.GetAdminStats)
		admin.GET("/users", adminController.ListUsers)
		admin.GET("/users/:id", adminController.GetUserByID)
		admin.PUT("/users/:id", adminController.UpdateUser)
		admin.DELETE("/users/:id", adminController.DeleteUser)
		admin.PUT("/users/:id/role", adminController.ChangeUserRole)
		admin.GET("/activity-log", adminController.GetActivityLog)
	}

	// *** NUEVAS RUTAS INTERNAS PARA COMUNICACIÓN ENTRE SERVICIOS ***
	// Estas rutas NO requieren autenticación JWT para permitir comunicación entre servicios
	internal := router.Group("/api/internal")
	{
		internal.GET("/users/:id", userController.GetUserByIDInternal)
		internal.GET("/health", internalHealthCheck)
	}

	// Rutas estáticas para servir imágenes
	router.Static("/static/profiles", "./static/profiles")

	return router
}

// healthCheck endpoint de health check
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"service":   "auth-user-service",
		"version":   "1.0.0",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// internalHealthCheck endpoint de health check para servicios internos
func internalHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"service":   "auth-user-service-internal",
		"version":   "1.0.0",
		"endpoint":  "internal",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}