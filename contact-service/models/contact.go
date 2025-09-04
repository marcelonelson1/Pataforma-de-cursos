package models

import (
	"time"
	"gorm.io/gorm"
)

// ContactMessage representa un mensaje de contacto
type ContactMessage struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Email     string         `gorm:"size:100;not null" json:"email"`
	Phone     string         `gorm:"size:20" json:"phone"`
	Message   string         `gorm:"type:text;not null" json:"message"`
	Read      bool           `gorm:"default:false" json:"read"`
	Starred   bool           `gorm:"default:false" json:"starred"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName especifica el nombre de la tabla
func (ContactMessage) TableName() string {
	return "contact_messages"
}

// ContactRequest estructura para crear un mensaje de contacto
type ContactRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Phone   string `json:"phone"`
	Message string `json:"message" binding:"required"`
}

// ContactResponse estructura para respuestas de mensajes
type ContactResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Message   string    `json:"message"`
	Read      bool      `json:"read"`
	Starred   bool      `json:"starred"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ReplyRequest estructura para responder a un mensaje
type ReplyRequest struct {
	Message string `json:"message" binding:"required"`
}

// EmailTemplateData estructura para datos del template de email
type EmailTemplateData struct {
	Name        string
	Email       string
	Subject     string
	OriginalMsg string
	ReplyMsg    string
}

// MarkAsRead marca el mensaje como leido
func (cm *ContactMessage) MarkAsRead() {
	cm.Read = true
}

// ToggleStar cambia el estado de estrella del mensaje
func (cm *ContactMessage) ToggleStar() {
	cm.Starred = !cm.Starred
}

// ToggleRead cambia el estado de lectura del mensaje
func (cm *ContactMessage) ToggleRead() {
	cm.Read = !cm.Read
}

// ToResponse convierte el modelo a estructura de respuesta
func (cm *ContactMessage) ToResponse() ContactResponse {
	return ContactResponse{
		ID:        cm.ID,
		Name:      cm.Name,
		Email:     cm.Email,
		Phone:     cm.Phone,
		Message:   cm.Message,
		Read:      cm.Read,
		Starred:   cm.Starred,
		CreatedAt: cm.CreatedAt,
		UpdatedAt: cm.UpdatedAt,
	}
}