package controllers

import (
	"contact-service/models"
	"contact-service/services"
	"contact-service/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ContactController struct {
	contactService *services.ContactService
}

// NewContactController crea una nueva instancia de ContactController
func NewContactController() *ContactController {
	return &ContactController{
		contactService: services.NewContactService(),
	}
}

// CreateMessage crea un nuevo mensaje de contacto (publico)
func (cc *ContactController) CreateMessage(c *gin.Context) {
	var req models.ContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, err)
		return
	}

	message, err := cc.contactService.CreateMessage(&req)
	if err != nil {
		utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	utils.SendSuccessMessage(c, "Mensaje enviado exitosamente", message.ToResponse())
}

// GetAllMessages obtiene todos los mensajes (admin)
func (cc *ContactController) GetAllMessages(c *gin.Context) {
	messages, err := cc.contactService.GetAllMessages()
	if err != nil {
		utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Convertir a respuestas
	var responses []models.ContactResponse
	for _, msg := range messages {
		responses = append(responses, msg.ToResponse())
	}

	utils.SendSuccessResponse(c, gin.H{
		"messages": responses,
		"total":    len(responses),
	})
}

// GetMessageByID obtiene un mensaje especifico (admin)
func (cc *ContactController) GetMessageByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendErrorMessage(c, "ID invalido", http.StatusBadRequest)
		return
	}

	message, err := cc.contactService.GetMessageByID(uint(id))
	if err != nil {
		if err == utils.ErrResourceNotFound {
			utils.SendErrorResponse(c, err, http.StatusNotFound)
		} else {
			utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		}
		return
	}

	// Marcar como leido automaticamente
	cc.contactService.MarkAsRead(uint(id))

	utils.SendSuccessResponse(c, message.ToResponse())
}

// UpdateMessageStatus actualiza el estado de un mensaje (admin)
func (cc *ContactController) UpdateMessageStatus(c *gin.Context) {
	idStr := c.Param("id")
	action := c.Param("action")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendErrorMessage(c, "ID invalido", http.StatusBadRequest)
		return
	}

	// Validar acciones permitidas
	if action != "read" && action != "star" {
		utils.SendErrorMessage(c, "Accion invalida. Permitidas: read, star", http.StatusBadRequest)
		return
	}

	message, err := cc.contactService.UpdateMessageStatus(uint(id), action)
	if err != nil {
		if err == utils.ErrResourceNotFound {
			utils.SendErrorResponse(c, err, http.StatusNotFound)
		} else {
			utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		}
		return
	}

	utils.SendSuccessMessage(c, "Estado actualizado correctamente", message.ToResponse())
}

// DeleteMessage elimina un mensaje (admin)
func (cc *ContactController) DeleteMessage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendErrorMessage(c, "ID invalido", http.StatusBadRequest)
		return
	}

	if err := cc.contactService.DeleteMessage(uint(id)); err != nil {
		utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	utils.SendSuccessMessage(c, "Mensaje eliminado correctamente", nil)
}

// ReplyToMessage responde a un mensaje de contacto (admin)
func (cc *ContactController) ReplyToMessage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendErrorMessage(c, "ID invalido", http.StatusBadRequest)
		return
	}

	var req models.ReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, err)
		return
	}

	if err := cc.contactService.ReplyToMessage(uint(id), req.Message); err != nil {
		if err == utils.ErrResourceNotFound {
			utils.SendErrorResponse(c, err, http.StatusNotFound)
		} else {
			utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		}
		return
	}

	utils.SendSuccessMessage(c, "Respuesta enviada correctamente", nil)
}

// GetContactStats obtiene estadisticas de mensajes (admin)
func (cc *ContactController) GetContactStats(c *gin.Context) {
	allMessages, err := cc.contactService.GetAllMessages()
	if err != nil {
		utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	unreadCount, err := cc.contactService.GetUnreadCount()
	if err != nil {
		utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	starredMessages, err := cc.contactService.GetStarredMessages()
	if err != nil {
		utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	stats := gin.H{
		"total_messages":   len(allMessages),
		"unread_messages":  unreadCount,
		"starred_messages": len(starredMessages),
		"read_messages":    len(allMessages) - int(unreadCount),
	}

	utils.SendSuccessResponse(c, stats)
}

// GetStarredMessages obtiene mensajes con estrella (admin)
func (cc *ContactController) GetStarredMessages(c *gin.Context) {
	messages, err := cc.contactService.GetStarredMessages()
	if err != nil {
		utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Convertir a respuestas
	var responses []models.ContactResponse
	for _, msg := range messages {
		responses = append(responses, msg.ToResponse())
	}

	utils.SendSuccessResponse(c, gin.H{
		"messages": responses,
		"total":    len(responses),
	})
}

// HealthCheck endpoint de verificacion de salud
func (cc *ContactController) HealthCheck(c *gin.Context) {
	utils.SendSuccessResponse(c, gin.H{
		"status":  "ok",
		"service": "contact-service",
		"version": "1.0.0",
	})
}