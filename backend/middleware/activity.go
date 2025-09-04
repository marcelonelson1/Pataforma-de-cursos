package middleware

import (
	"curso-platform/config"
	"curso-platform/models"
	"curso-platform/utils"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LogActivity registra una actividad en el sistema
func LogActivity(c *gin.Context, userID uint, action, details string) {
	// Obtener IP y User-Agent
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	
	activityLog := models.ActivityLog{
		UserID:    userID,
		Action:    action,
		Details:   details,
		IP:        ip,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
	}
	
	if err := config.DB.Create(&activityLog).Error; err != nil {
		log.Printf("Error al registrar actividad: %v", err)
	}
}

// ActivityLogger registra ciertas actividades automáticamente
func ActivityLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ejecutar el resto de los handlers primero
		c.Next()

		// Solo registrar actividad para ciertas rutas
		path := c.Request.URL.Path
		method := c.Request.Method

		// Lista de patrones de rutas que queremos monitorear
		monitoredPaths := map[string]string{
			"/api/auth/login":    "login",
			"/api/auth/register": "register",
			"/api/admin":         "admin_access",
			"/api/cursos":        "course_management",
		}

		// Verificar si la ruta actual está en la lista de monitoreo
		for pattern, actionType := range monitoredPaths {
			if utils.PathMatchesPattern(path, pattern) {
				// Solo registrar si hay un usuario en el contexto
				if userValue, exists := c.Get("user"); exists {
					user, ok := userValue.(models.Usuario)
					if ok {
						// Personalizar los detalles según la ruta
						details := generateActivityDetails(c, path, method, actionType)
						LogActivity(c, user.ID, actionType, details)
					}
				}
				break
			}
		}
	}
}

// Genera detalles personalizados para el registro de actividad
func generateActivityDetails(c *gin.Context, path, method, actionType string) string {
	switch actionType {
	case "login":
		return "Inicio de sesión exitoso"
	case "register":
		return "Registro de nuevo usuario"
	case "admin_access":
		return "Acceso a panel de administración: " + path
	case "course_management":
		operation := "acceso"
		if method == "POST" {
			operation = "creación"
		} else if method == "PUT" {
			operation = "actualización"
		} else if method == "DELETE" {
			operation = "eliminación"
		}
		return "Gestión de cursos: " + operation
	default:
		return "Acceso a: " + path
	}
}