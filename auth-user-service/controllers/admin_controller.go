package controllers

import (
	"auth-user-service/models"
	"auth-user-service/services"
	"auth-user-service/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	userService *services.UserService
}

func NewAdminController() *AdminController {
	return &AdminController{
		userService: services.NewUserService(),
	}
}

// ListUsers lista todos los usuarios (para admin)
func (a *AdminController) ListUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	users, total, err := a.userService.ListUsers(page, limit)
	if err != nil {
		utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	response := gin.H{
		"users": users,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (int(total) + limit - 1) / limit,
		},
	}

	utils.SendSuccessResponse(c, response)
}

// GetUserByID obtiene un usuario específico (para admin)
func (a *AdminController) GetUserByID(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.SendErrorMessage(c, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	// Para admin, incluir el perfil completo directamente
	profile, err := a.userService.GetUserProfile(uint(userID))
	if err != nil {
		if err == utils.ErrUserNotFound {
			utils.SendErrorResponse(c, err, http.StatusNotFound)
		} else {
			utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		}
		return
	}

	utils.SendSuccessResponse(c, profile)
}

// UpdateUser actualiza un usuario (para admin)
func (a *AdminController) UpdateUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.SendErrorMessage(c, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.SendValidationError(c, err)
		return
	}

	// No permitir actualizar ciertos campos críticos
	delete(updates, "id")
	delete(updates, "password") // La contraseña se cambia por otros endpoints
	delete(updates, "created_at")

	if len(updates) == 0 {
		utils.SendErrorMessage(c, "No hay campos para actualizar", http.StatusBadRequest)
		return
	}

	if err := a.userService.UpdateUser(uint(userID), updates); err != nil {
		if err == utils.ErrUserNotFound {
			utils.SendErrorResponse(c, err, http.StatusNotFound)
		} else {
			utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		}
		return
	}

	utils.SendSuccessMessage(c, "Usuario actualizado correctamente", nil)
}

// DeleteUser elimina un usuario (para admin)
func (a *AdminController) DeleteUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.SendErrorMessage(c, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	// Obtener información del admin que está eliminando
	adminValue, exists := c.Get("user")
	if !exists {
		utils.SendErrorResponse(c, utils.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	admin, ok := adminValue.(*models.Usuario)
	if !ok {
		utils.SendErrorResponse(c, utils.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	// Prevenir que un admin se elimine a sí mismo
	if admin.ID == uint(userID) {
		utils.SendErrorMessage(c, "No puedes eliminarte a ti mismo", http.StatusBadRequest)
		return
	}

	if err := a.userService.DeleteUser(uint(userID)); err != nil {
		if err == utils.ErrUserNotFound {
			utils.SendErrorResponse(c, err, http.StatusNotFound)
		} else {
			utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		}
		return
	}

	// Log de la eliminación por parte del admin
	a.userService.LogActivity(admin.ID, "admin_delete_user", 
		fmt.Sprintf("Admin eliminó usuario ID: %d", userID), 
		c.ClientIP(), c.GetHeader("User-Agent"))

	utils.SendSuccessMessage(c, "Usuario eliminado correctamente", nil)
}

// ChangeUserRole cambia el rol de un usuario (para admin)
func (a *AdminController) ChangeUserRole(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.SendErrorMessage(c, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	var request struct {
		Role string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.SendValidationError(c, err)
		return
	}

	// Validar rol
	validRoles := []string{"user", "admin"}
	isValidRole := false
	for _, validRole := range validRoles {
		if request.Role == validRole {
			isValidRole = true
			break
		}
	}

	if !isValidRole {
		utils.SendErrorMessage(c, "Rol inválido. Roles permitidos: user, admin", http.StatusBadRequest)
		return
	}

	// Obtener información del admin que está cambiando el rol
	adminValue, exists := c.Get("user")
	if !exists {
		utils.SendErrorResponse(c, utils.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	admin := adminValue.(*models.Usuario)

	// Prevenir que un admin se quite sus propios permisos
	if admin.ID == uint(userID) && request.Role != "admin" {
		utils.SendErrorMessage(c, "No puedes cambiar tu propio rol de admin", http.StatusBadRequest)
		return
	}

	if err := a.userService.ChangeUserRole(uint(userID), request.Role); err != nil {
		if err == utils.ErrUserNotFound {
			utils.SendErrorResponse(c, err, http.StatusNotFound)
		} else {
			utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		}
		return
	}

	utils.SendSuccessMessage(c, "Rol actualizado correctamente", gin.H{
		"user_id": userID,
		"new_role": request.Role,
	})
}

// GetActivityLog obtiene el log de actividades (para admin)
func (a *AdminController) GetActivityLog(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "50")
	userIDStr := c.Query("user_id")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 50
	}

	var userID *uint
	if userIDStr != "" {
		if parsedUserID, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
			userIDUint := uint(parsedUserID)
			userID = &userIDUint
		}
	}

	activities, total, err := a.userService.GetActivityLog(page, limit, userID)
	if err != nil {
		utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	response := gin.H{
		"activities": activities,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (int(total) + limit - 1) / limit,
		},
	}

	utils.SendSuccessResponse(c, response)
}

// GetAdminStats obtiene estadísticas para el dashboard de admin
func (a *AdminController) GetAdminStats(c *gin.Context) {
	// Contar usuarios por rol
	var userCount, adminCount int64
	db := utils.GetDB()
	db.Model(&models.Usuario{}).Where("role = ?", "user").Count(&userCount)
	db.Model(&models.Usuario{}).Where("role = ?", "admin").Count(&adminCount)

	// Actividad reciente
	recentActivities, _, err := a.userService.GetActivityLog(1, 10, nil)
	if err != nil {
		utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	stats := gin.H{
		"total_users":       userCount + adminCount,
		"regular_users":     userCount,
		"admin_users":       adminCount,
		"recent_activities": recentActivities,
	}

	utils.SendSuccessResponse(c, stats)
}