package services

import (
	"curso-platform/config"
	"curso-platform/models"
	"curso-platform/utils"
	"fmt"
	"log"

	"gorm.io/gorm"
)

// ContactService maneja la lógica relacionada con los mensajes de contacto
type ContactService struct{}

// NewContactService crea una nueva instancia de ContactService
func NewContactService() *ContactService {
	return &ContactService{}
}

// SendContactMessage guarda un mensaje de contacto y envía email de notificación
func (s *ContactService) SendContactMessage(req models.ContactRequest) error {
	// Guardar en la base de datos
	message := models.ContactMessage{
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
		Message: req.Message,
	}

	if err := config.DB.Create(&message).Error; err != nil {
		log.Printf("Error guardando mensaje: %v", err)
		return utils.ErrDatabaseError
	}

	// Enviar por email
	if err := s.sendContactEmail(req); err != nil {
		log.Printf("Error enviando email: %v", err)
		return utils.ErrEmailSendError
	}

	return nil
}

// GetAllMessages obtiene todos los mensajes de contacto
func (s *ContactService) GetAllMessages() ([]models.ContactMessage, error) {
	var messages []models.ContactMessage

	if err := config.DB.Order("created_at DESC").Find(&messages).Error; err != nil {
		return nil, utils.ErrDatabaseError
	}

	return messages, nil
}

// GetMessageByID obtiene un mensaje de contacto por su ID
func (s *ContactService) GetMessageByID(id uint) (*models.ContactMessage, error) {
	var message models.ContactMessage
	if err := config.DB.First(&message, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrResourceNotFound
		}
		return nil, utils.ErrDatabaseError
	}

	// Marcar como leído si no lo estaba
	if !message.Read {
		config.DB.Model(&message).Update("read", true)
	}

	return &message, nil
}

// UpdateMessageStatus actualiza el estado (leído/destacado) de un mensaje
func (s *ContactService) UpdateMessageStatus(id uint, action string) (*models.ContactMessage, error) {
	var message models.ContactMessage
	if err := config.DB.First(&message, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrResourceNotFound
		}
		return nil, utils.ErrDatabaseError
	}

	switch action {
	case "read":
		config.DB.Model(&message).Update("read", !message.Read)
	case "star":
		config.DB.Model(&message).Update("starred", !message.Starred)
	default:
		return nil, utils.ErrInvalidRequest
	}

	// Recuperar el mensaje actualizado
	if err := config.DB.First(&message, id).Error; err != nil {
		return nil, utils.ErrDatabaseError
	}

	return &message, nil
}

// DeleteMessage elimina un mensaje de contacto
func (s *ContactService) DeleteMessage(id uint) error {
	if err := config.DB.Delete(&models.ContactMessage{}, id).Error; err != nil {
		return utils.ErrDatabaseError
	}
	return nil
}

// ReplyToMessage envía una respuesta a un mensaje de contacto
func (s *ContactService) ReplyToMessage(id uint, replyText string) error {
	var message models.ContactMessage
	if err := config.DB.First(&message, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrResourceNotFound
		}
		return utils.ErrDatabaseError
	}

	// Enviar email de respuesta
	if err := s.sendReplyEmail(message, replyText); err != nil {
		return utils.ErrEmailSendError
	}

	log.Printf("Respuesta enviada a %s (%s): %s", message.Name, message.Email, replyText)
	return nil
}

// Enviar email al recibir un contacto
func (s *ContactService) sendContactEmail(contact models.ContactRequest) error {
	to := utils.GetEnv("CONTACT_EMAIL", "admin@example.com")
	subject := "Nuevo mensaje de contacto - " + contact.Name
	
	// Datos para la plantilla (en una implementación real, usaríamos una plantilla HTML)
	body := fmt.Sprintf(`
Nuevo mensaje de contacto recibido:

Nombre: %s
Email: %s
Teléfono: %s
Mensaje:
%s
	`, contact.Name, contact.Email, contact.Phone, contact.Message)
	
	// Enviar email
	return utils.SendEmail(to, subject, body, "text/plain; charset=\"UTF-8\"")
}

// Enviar email de respuesta a quien contactó
func (s *ContactService) sendReplyEmail(contact models.ContactMessage, replyText string) error {
	// Preparar datos para el template
	data := struct {
		Name        string
		Email       string
		OriginalMsg string
		ReplyMsg    string
	}{
		Name:        contact.Name,
		Email:       contact.Email,
		OriginalMsg: contact.Message,
		ReplyMsg:    replyText,
	}
	
	// Renderizar la plantilla
	emailBody, err := utils.RenderEmailTemplate("contact_reply", data)
	if err != nil {
		return err
	}
	
	// Enviar email
	subject := "Respuesta a tu mensaje"
	return utils.SendEmail(contact.Email, subject, emailBody, "text/html; charset=\"UTF-8\"")
}