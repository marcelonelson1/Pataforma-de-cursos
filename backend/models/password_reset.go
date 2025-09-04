package models

import (
	"time"
)

// PasswordReset modelo para almacenar los tokens de restablecimiento
type PasswordReset struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"size:100;not null;index" json:"email"`
	Token     string    `gorm:"size:100;not null;uniqueIndex" json:"token"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `gorm:"default:false" json:"used"`
}

// TableName especifica el nombre de la tabla para el modelo PasswordReset
func (PasswordReset) TableName() string {
	return "password_resets"
}

// ForgotPasswordRequest estructura para la solicitud de recuperación
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest estructura para restablecer contraseña
type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}