package controllers

import (
	"analytics-service/config"
	"analytics-service/middleware"
	"analytics-service/services"
	"analytics-service/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AnalyticsController struct {
	analyticsService *services.AnalyticsService
	profileService   *services.ProfileService
	config          *config.Config
}

func NewAnalyticsController(cfg *config.Config) *AnalyticsController {
	analyticsService := services.NewAnalyticsService(
		cfg.AuthServiceURL,
		cfg.PaymentServiceURL,
		cfg.CourseServiceURL,
		cfg.ContactServiceURL,
		cfg.PortfolioServiceURL,
		cfg.HomeServiceURL,
	)
	
	profileService := services.NewProfileService(cfg.AuthServiceURL)
	
	return &AnalyticsController{
		analyticsService: analyticsService,
		profileService:   profileService,
		config:          cfg,
	}
}

// Health Check
func (ac *AnalyticsController) HealthCheck(c *gin.Context) {
	utils.SuccessResponse(c, gin.H{
		"service": "analytics-service",
		"status":  "healthy",
		"port":    ac.config.Port,
	}, "Analytics service is running")
}

// Dashboard metrics
func (ac *AnalyticsController) GetDashboard(c *gin.Context) {
	metrics, err := ac.analyticsService.GetDashboardMetrics()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Error al obtener métricas del dashboard")
		return
	}
	
	utils.SuccessResponse(c, metrics, "Dashboard metrics retrieved successfully")
}

// Sales statistics
func (ac *AnalyticsController) GetSalesStats(c *gin.Context) {
	period := c.DefaultQuery("period", "month")
	
	stats, err := ac.analyticsService.GetSalesStats(period)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Error al obtener estadísticas de ventas")
		return
	}
	
	utils.SuccessResponse(c, stats, "Sales statistics retrieved successfully")
}

// Activity logs
func (ac *AnalyticsController) GetActivityLog(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "50")
	
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}
	
	logs, err := ac.analyticsService.GetActivityLog(page, limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Error al obtener logs de actividad")
		return
	}
	
	utils.SuccessResponse(c, logs, "Activity logs retrieved successfully")
}

// Admin statistics
func (ac *AnalyticsController) GetAdminStats(c *gin.Context) {
	period := c.DefaultQuery("period", "month")
	
	stats, err := ac.analyticsService.GetSalesStats(period)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Error al obtener estadísticas")
		return
	}
	
	utils.SuccessResponse(c, stats, "Admin statistics retrieved successfully")
}

// User profile endpoints
func (ac *AnalyticsController) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User ID not found")
		return
	}
	
	profile, err := ac.profileService.GetUserProfile(userID.(uint))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Error al obtener perfil")
		return
	}
	
	utils.SuccessResponse(c, profile, "Profile retrieved successfully")
}

func (ac *AnalyticsController) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User ID not found")
		return
	}
	
	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.BadRequestResponse(c, "Datos inválidos")
		return
	}
	
	profile, err := ac.profileService.UpdateUserProfile(userID.(uint), updateData)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Error al actualizar perfil")
		return
	}
	
	utils.SuccessResponse(c, profile, "Profile updated successfully")
}

func (ac *AnalyticsController) GetNotificationSettings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User ID not found")
		return
	}
	
	settings, err := ac.profileService.GetNotificationSettings(userID.(uint))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Error al obtener configuraciones")
		return
	}
	
	utils.SuccessResponse(c, settings, "Notification settings retrieved successfully")
}

func (ac *AnalyticsController) UpdateNotificationSettings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User ID not found")
		return
	}
	
	var settings map[string]interface{}
	if err := c.ShouldBindJSON(&settings); err != nil {
		utils.BadRequestResponse(c, "Datos inválidos")
		return
	}
	
	updatedSettings, err := ac.profileService.UpdateNotificationSettings(userID.(uint), settings)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Error al actualizar configuraciones")
		return
	}
	
	utils.SuccessResponse(c, updatedSettings, "Notification settings updated successfully")
}

// Check admin status
func (ac *AnalyticsController) CheckAdmin(c *gin.Context) {
	userRole, exists := c.Get("user_role")
	if !exists {
		utils.UnauthorizedResponse(c, "User role not found")
		return
	}
	
	isAdmin := userRole.(string) == "admin"
	
	utils.SuccessResponse(c, gin.H{
		"isAdmin": isAdmin,
		"role":    userRole,
	}, "Admin status checked successfully")
}

// Setup routes
func (ac *AnalyticsController) SetupRoutes(router *gin.Engine) {
	// Public routes
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "analytics-service",
			"version": "1.0.0",
			"status":  "running",
		})
	})
	
	router.GET("/api/health", ac.HealthCheck)
	
	// Authenticated routes
	authGroup := router.Group("/api/auth")
	authGroup.Use(middleware.AuthMiddleware(ac.config))
	{
		authGroup.GET("/profile", ac.GetProfile)
		authGroup.PUT("/profile", ac.UpdateProfile)
		authGroup.GET("/notification-settings", ac.GetNotificationSettings)
		authGroup.PUT("/notification-settings", ac.UpdateNotificationSettings)
		authGroup.GET("/check-admin", ac.CheckAdmin)
	}
	
	// Admin routes
	adminGroup := router.Group("/api/admin")
	adminGroup.Use(middleware.AuthMiddleware(ac.config))
	adminGroup.Use(middleware.AdminMiddleware(ac.config))
	{
		adminGroup.GET("/dashboard", ac.GetDashboard)
		adminGroup.GET("/stats", ac.GetAdminStats)
		adminGroup.GET("/sales-stats", ac.GetSalesStats)
		adminGroup.GET("/activity-log", ac.GetActivityLog)
	}
}