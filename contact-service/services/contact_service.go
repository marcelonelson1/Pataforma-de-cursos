package services

import (
	"contact-service/models"
	"contact-service/utils"
	"log"

	"gorm.io/gorm"
)

type ContactService struct {
	db           *gorm.DB
	emailService *EmailService
}

// NewContactService crea una nueva instancia de ContactService
func NewContactService() *ContactService {
	return &ContactService{
		db:           utils.GetDB(),
		emailService: NewEmailService(),
	}
}

// CreateMessage crea un nuevo mensaje de contacto
func (cs *ContactService) CreateMessage(req *models.ContactRequest) (*models.ContactMessage, error) {
	message := &models.ContactMessage{
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
		Message: req.Message,
		Read:    false,
		Starred: false,
	}

	// Guardar en base de datos
	if err := cs.db.Create(message).Error; err != nil {
		log.Printf("Error guardando mensaje: %v", err)
		return nil, utils.ErrDatabaseError
	}

	// Enviar notificacion por email (asincrono)
	go func() {
		if err := cs.emailService.SendContactNotification(req); err != nil {
			log.Printf("Error enviando notificacion de contacto: %v", err)
		}
	}()

	return message, nil
}

// GetAllMessages obtiene todos los mensajes de contacto
func (cs *ContactService) GetAllMessages() ([]models.ContactMessage, error) {
	var messages []models.ContactMessage

	if err := cs.db.Order("created_at DESC").Find(&messages).Error; err != nil {
		return nil, utils.ErrDatabaseError
	}

	return messages, nil
}

// GetMessageByID obtiene un mensaje por ID
func (cs *ContactService) GetMessageByID(id uint) (*models.ContactMessage, error) {
	var message models.ContactMessage

	if err := cs.db.First(&message, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrResourceNotFound
		}
		return nil, utils.ErrDatabaseError
	}

	return &message, nil
}

// MarkAsRead marca un mensaje como leido
func (cs *ContactService) MarkAsRead(id uint) error {
	message, err := cs.GetMessageByID(id)
	if err != nil {
		return err
	}

	if !message.Read {
		message.MarkAsRead()
		if err := cs.db.Save(message).Error; err != nil {
			return utils.ErrDatabaseError
		}
	}

	return nil
}

// UpdateMessageStatus actualiza el estado de un mensaje (read/star)
func (cs *ContactService) UpdateMessageStatus(id uint, action string) (*models.ContactMessage, error) {
	message, err := cs.GetMessageByID(id)
	if err != nil {
		return nil, err
	}

	switch action {
	case "read":
		message.ToggleRead()
	case "star":
		message.ToggleStar()
	default:
		return nil, utils.ErrInvalidRequest
	}

	if err := cs.db.Save(message).Error; err != nil {
		return nil, utils.ErrDatabaseError
	}

	return message, nil
}

// DeleteMessage elimina un mensaje
func (cs *ContactService) DeleteMessage(id uint) error {
	if err := cs.db.Delete(&models.ContactMessage{}, id).Error; err != nil {
		return utils.ErrDatabaseError
	}

	return nil
}

// ReplyToMessage responde a un mensaje de contacto
func (cs *ContactService) ReplyToMessage(id uint, replyText string) error {
	message, err := cs.GetMessageByID(id)
	if err != nil {
		return err
	}

	// Enviar email de respuesta
	if err := cs.emailService.SendReplyEmail(message, replyText); err != nil {
		log.Printf("Error enviando respuesta: %v", err)
		return utils.ErrEmailSendError
	}

	// Marcar como leido si no lo estaba
	if !message.Read {
		message.MarkAsRead()
		cs.db.Save(message)
	}

	log.Printf("Respuesta enviada a %s (%s): %s", message.Name, message.Email, replyText)
	return nil
}

// GetUnreadCount obtiene el numero de mensajes no leidos
func (cs *ContactService) GetUnreadCount() (int64, error) {
	var count int64
	
	if err := cs.db.Model(&models.ContactMessage{}).Where("read = ?", false).Count(&count).Error; err != nil {
		return 0, utils.ErrDatabaseError
	}

	return count, nil
}

// GetStarredMessages obtiene mensajes marcados con estrella
func (cs *ContactService) GetStarredMessages() ([]models.ContactMessage, error) {
	var messages []models.ContactMessage

	if err := cs.db.Where("starred = ?", true).Order("created_at DESC").Find(&messages).Error; err != nil {
		return nil, utils.ErrDatabaseError
	}

	return messages, nil
}