package controllers

import (
	"auth-user-service/services"
	"auth-user-service/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: services.NewUserService(),
	}
}

// GetUserByID obtiene un usuario por ID (para otros servicios - requiere autenticación)
func (u *UserController) GetUserByID(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.SendErrorMessage(c, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	user, err := u.userService.GetUserByID(uint(userID))
	if err != nil {
		if err == utils.ErrUserNotFound {
			utils.SendErrorResponse(c, err, http.StatusNotFound)
		} else {
			utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		}
		return
	}

	utils.SendSuccessResponse(c, user)
}

// GetUserByIDInternal obtiene un usuario por ID para comunicación entre servicios
// NO requiere autenticación JWT - específicamente para payment-service
func (u *UserController) GetUserByIDInternal(c *gin.Context) {
	// Validar origen del servicio (opcional pero recomendado)
	serviceHeader := c.GetHeader("X-Service-Name")
	clientIP := c.ClientIP()

	if serviceHeader == "" {
		log.Printf("Advertencia: Petición a endpoint interno sin identificación de servicio desde IP: %s", clientIP)
		// En producción podrías agregar validación más estricta aquí
	} else {
		log.Printf("Petición interna recibida del servicio: %s desde IP: %s", serviceHeader, clientIP)
	}

	// Parsear ID del usuario
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.SendErrorMessage(c, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	// Obtener usuario
	user, err := u.userService.GetUserByID(uint(userID))
	if err != nil {
		if err == utils.ErrUserNotFound {
			utils.SendErrorResponse(c, err, http.StatusNotFound)
		} else {
			utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		}
		return
	}

	// Estructura de respuesta específica que espera payment-service
	response := gin.H{
		"success": true,
		"data": gin.H{
			"id":     user.ID,
			"nombre": user.Nombre,
			"email":  user.Email,
			"role":   user.Role,
		},
		"message": "Usuario encontrado",
	}

	log.Printf("Usuario encontrado para servicio interno: ID=%d, Email=%s, Role=%s",
		user.ID, user.Email, user.Role)

	c.JSON(http.StatusOK, response)
}

// GetUsers obtiene lista de usuarios con paginación (para otros servicios)
func (u *UserController) GetUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	users, total, err := u.userService.ListUsers(page, limit)
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
		},
	}

	utils.SendSuccessResponse(c, response)
}
