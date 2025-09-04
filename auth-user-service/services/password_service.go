package services

import (
	"auth-user-service/models"
	"auth-user-service/utils"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PasswordService struct {
	db *gorm.DB
}

func NewPasswordService() *PasswordService {
	return &PasswordService{
		db: utils.GetDB(),
	}
}

// GenerateResetToken genera un nuevo token de restablecimiento
func (p *PasswordService) GenerateResetToken(email string) (*models.PasswordReset, error) {
	tokenStr := uuid.New().String()
	reset := &models.PasswordReset{
		Email:     email,
		Token:     tokenStr,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      false,
	}

	// Eliminar tokens anteriores para el mismo email
	if err := p.db.Where("email = ?", email).Delete(&models.PasswordReset{}).Error; err != nil {
		log.Printf("Error al eliminar tokens anteriores: %v", err)
	}

	// Guardar el nuevo token
	if result := p.db.Create(&reset); result.Error != nil {
		log.Printf("Error al crear token: %v", result.Error)
		return nil, result.Error
	}

	return reset, nil
}

// ValidateResetToken verifica si un token es válido
func (p *PasswordService) ValidateResetToken(token string) (*models.PasswordReset, error) {
	var reset models.PasswordReset
	result := p.db.Where("token = ? AND used = ? AND expires_at > ?", token, false, time.Now()).First(&reset)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, utils.ErrInvalidToken
		}
		log.Printf("Error de base de datos al validar token: %v", result.Error)
		return nil, result.Error
	}
	return &reset, nil
}

// MarkTokenAsUsed marca un token como utilizado
func (p *PasswordService) MarkTokenAsUsed(token string) error {
	result := p.db.Model(&models.PasswordReset{}).Where("token = ?", token).Update("used", true)
	if result.Error != nil {
		log.Printf("Error al marcar token como usado: %v", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return utils.ErrInvalidToken
	}
	return nil
}

// ResetPassword cambia la contraseña usando un token de reset
func (p *PasswordService) ResetPassword(token, newPassword string) error {
	// Validar el token primero
	reset, err := p.ValidateResetToken(token)
	if err != nil {
		return err
	}

	// Hash de la nueva contraseña
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Actualizar la contraseña del usuario
	if result := p.db.Model(&models.Usuario{}).Where("email = ?", reset.Email).Update("password", hashedPassword); result.Error != nil {
		log.Printf("Error al actualizar contraseña: %v", result.Error)
		return result.Error
	}

	// Marcar el token como usado
	return p.MarkTokenAsUsed(token)
}

// ChangePassword cambia la contraseña del usuario verificando la actual
func (p *PasswordService) ChangePassword(userID uint, currentPassword, newPassword string) error {
	// Buscar el usuario
	var user models.Usuario
	if result := p.db.First(&user, userID); result.Error != nil {
		return utils.ErrUserNotFound
	}

	// Verificar contraseña actual
	if err := utils.CheckPassword(user.Password, currentPassword); err != nil {
		return utils.ErrInvalidPassword
	}

	// Hash de la nueva contraseña
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Actualizar la contraseña
	if result := p.db.Model(&user).Update("password", hashedPassword); result.Error != nil {
		return result.Error
	}

	return nil
}