package models

import "time"

// PasswordReset modelo para almacenar los tokens de restablecimiento
type PasswordReset struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"size:100;not null;index" json:"email"`
	Token     string    `gorm:"size:100;not null;uniqueIndex" json:"token"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `gorm:"default:false" json:"used"`
}