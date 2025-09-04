package models

import (
	

	"gorm.io/gorm"
)

// ContactMessage representa un mensaje de contacto enviado por un usuario
type ContactMessage struct {
	gorm.Model
	Name    string `gorm:"size:100;not null" json:"name"`
	Email   string `gorm:"size:100;not null" json:"email"`
	Phone   string `gorm:"size:20" json:"phone"`
	Message string `gorm:"type:text;not null" json:"message"`
	Read    bool   `gorm:"default:false" json:"read"`
	Starred bool   `gorm:"default:false" json:"starred"`
}

// TableName especifica el nombre de la tabla para el modelo ContactMessage
func (ContactMessage) TableName() string {
	return "contact_messages"
}

// ContactRequest representa una solicitud de contacto enviada por el cliente
type ContactRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Phone   string `json:"phone"`
	Message string `json:"message" binding:"required"`
}