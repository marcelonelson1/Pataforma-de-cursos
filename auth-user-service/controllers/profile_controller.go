package controllers

import (
	"auth-user-service/models"
	"auth-user-service/services"
	"auth-user-service/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProfileController struct {
	userService *services.UserService
}

func NewProfileController() *ProfileController {
	return &ProfileController{
		userService: services.NewUserService(),
	}
}

// GetProfile obtiene el perfil del usuario
func (p *ProfileController) GetProfile(c *gin.Context) {
	// Obtener el usuario del contexto
	userValue, exists := c.Get("user")
	if !exists {
		utils.SendErrorResponse(c, utils.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	user, ok := userValue.(*models.Usuario)
	if !ok {
		utils.SendErrorMessage(c, "error al obtener información del usuario", http.StatusInternalServerError)
		return
	}

	profile, err := p.userService.GetUserProfile(user.ID)
	if err != nil {
		utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, profile)
}

// UpdateProfile actualiza el perfil del usuario
func (p *ProfileController) UpdateProfile(c *gin.Context) {
	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, err)
		return
	}

	// Obtener el usuario del contexto
	userValue, exists := c.Get("user")
	if !exists {
		utils.SendErrorResponse(c, utils.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	user, ok := userValue.(*models.Usuario)
	if !ok {
		utils.SendErrorMessage(c, "error al obtener información del usuario", http.StatusInternalServerError)
		return
	}

	if err := p.userService.UpdateProfile(user.ID, &req); err != nil {
		utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	utils.SendSuccessMessage(c, "Perfil actualizado correctamente", nil)
}

// GetNotificationSettings obtiene las configuraciones de notificación
func (p *ProfileController) GetNotificationSettings(c *gin.Context) {
	// Obtener el usuario del contexto
	userValue, exists := c.Get("user")
	if !exists {
		utils.SendErrorResponse(c, utils.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	user, ok := userValue.(*models.Usuario)
	if !ok {
		utils.SendErrorMessage(c, "error al obtener información del usuario", http.StatusInternalServerError)
		return
	}

	settings, err := p.userService.GetNotificationSettings(user.ID)
	if err != nil {
		utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, settings)
}

// UpdateNotificationSettings actualiza las configuraciones de notificación
func (p *ProfileController) UpdateNotificationSettings(c *gin.Context) {
	var req models.UpdateNotificationSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, err)
		return
	}

	// Obtener el usuario del contexto
	userValue, exists := c.Get("user")
	if !exists {
		utils.SendErrorResponse(c, utils.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	user, ok := userValue.(*models.Usuario)
	if !ok {
		utils.SendErrorMessage(c, "error al obtener información del usuario", http.StatusInternalServerError)
		return
	}

	if err := p.userService.UpdateNotificationSettings(user.ID, &req); err != nil {
		utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	utils.SendSuccessMessage(c, "Configuraciones actualizadas correctamente", nil)
}