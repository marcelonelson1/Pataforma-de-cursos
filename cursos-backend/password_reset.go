package main

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

// GenerateResetToken genera un nuevo token de restablecimiento
func GenerateResetToken(email string) (*PasswordReset, error) {
	tokenStr := uuid.New().String()
	reset := &PasswordReset{
		Email:     email,
		Token:     tokenStr,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      false,
	}

	// Eliminar tokens anteriores para el mismo email
	if err := db.Where("email = ?", email).Delete(&PasswordReset{}).Error; err != nil {
		log.Printf("Error al eliminar tokens anteriores: %v", err)
	}

	// Guardar el nuevo token
	if result := db.Create(&reset); result.Error != nil {
		log.Printf("Error al crear token: %v", result.Error)
		return nil, result.Error
	}

	return reset, nil
}

// ValidateResetToken verifica si un token es válido
func ValidateResetToken(token string) (*PasswordReset, error) {
	var reset PasswordReset
	result := db.Where("token = ? AND used = ? AND expires_at > ?", token, false, time.Now()).First(&reset)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidToken
		}
		log.Printf("Error de base de datos al validar token: %v", result.Error)
		return nil, result.Error
	}
	return &reset, nil
}

// MarkTokenAsUsed marca un token como utilizado
func MarkTokenAsUsed(token string) error {
	result := db.Model(&PasswordReset{}).Where("token = ?", token).Update("used", true)
	if result.Error != nil {
		log.Printf("Error al marcar token como usado: %v", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrInvalidToken
	}
	return nil
}
