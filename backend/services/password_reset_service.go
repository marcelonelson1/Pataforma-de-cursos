// services/password_reset_service.go
package services

import (
	"curso-platform/config"
	"curso-platform/models"
	"curso-platform/utils"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// PasswordResetService estructura para manejar la lógica de restablecimiento de contraseñas
type PasswordResetService struct {
	db *gorm.DB
}

// NewPasswordResetService crea una nueva instancia del servicio de restablecimiento de contraseñas
func NewPasswordResetService(db *gorm.DB) *PasswordResetService {
	return &PasswordResetService{db: db}
}

// GetDB devuelve la conexión de base de datos
func (s *PasswordResetService) GetDB() *gorm.DB {
	return s.db
}

// GenerateResetToken genera un token de restablecimiento para un usuario
func (s *PasswordResetService) GenerateResetToken(email string) (*models.PasswordReset, error) {
	log.Printf("Generando token de recuperación para email: %s", email)
	
	// Verificar si el email existe
	var user models.Usuario
	result := config.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		log.Printf("Email no encontrado o error de DB: %v", result.Error)
		return nil, errors.New("email no encontrado")
	}
	
	log.Printf("Usuario encontrado con ID: %d, Nombre: %s", user.ID, user.Nombre)

	// Generar token único
	tokenStr := uuid.New().String()
	log.Printf("Token generado: %s", tokenStr)
	
	// Crear el objeto de reset de contraseña
	reset := &models.PasswordReset{
		Email:     email,
		Token:     tokenStr,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour), // Aumentar a 24 horas para dar más tiempo
		Used:      false,
	}

	// Eliminar tokens anteriores para el mismo email
	if err := config.DB.Where("email = ?", email).Delete(&models.PasswordReset{}).Error; err != nil {
		log.Printf("Error al eliminar tokens anteriores: %v", err)
	} else {
		log.Println("Tokens anteriores eliminados correctamente")
	}

	// Guardar el nuevo token
	if result := config.DB.Create(reset); result.Error != nil {
		log.Printf("Error al crear token en la base de datos: %v", result.Error)
		return nil, result.Error
	}
	
	log.Printf("Token guardado correctamente en la base de datos con ID: %d", reset.ID)
	
	// Verificar que el token se guardó correctamente
	var savedReset models.PasswordReset
	if err := config.DB.Where("token = ?", reset.Token).First(&savedReset).Error; err != nil {
		log.Printf("ADVERTENCIA: No se pudo verificar el token guardado: %v", err)
		
		// Intentar recuperar el token con una consulta más general para diagnóstico
		var anyTokens []models.PasswordReset
		if err := config.DB.Where("email = ?", email).Find(&anyTokens).Error; err != nil {
			log.Printf("No se encontraron tokens para el email: %s", email)
		} else {
			log.Printf("Tokens encontrados para el email %s: %d", email, len(anyTokens))
			for i, t := range anyTokens {
				log.Printf("Token %d: ID=%d, Token=%s, Expiración=%v, Usado=%v", 
					i+1, t.ID, t.Token, t.ExpiresAt, t.Used)
			}
		}
	} else {
		log.Printf("Token verificado en la base de datos. ID: %d, Token: %s, Email: %s, Expira: %v",
			savedReset.ID, savedReset.Token, savedReset.Email, savedReset.ExpiresAt)
	}

	return reset, nil
}

// ValidateResetToken verifica si un token es válido
func (s *PasswordResetService) ValidateResetToken(token string) (*models.PasswordReset, error) {
	log.Printf("Validando token: %s", token)
	
	// Validación estándar del token
	var reset models.PasswordReset
	result := config.DB.Where("token = ? AND used = ?", token, false).First(&reset)
	
	// Si encontramos el token pero está expirado, dar un mensaje específico
	if result.Error == nil && reset.ExpiresAt.Before(time.Now()) {
		log.Printf("Token expirado: %s, expiró el: %v", token, reset.ExpiresAt)
		return nil, errors.New("el enlace de recuperación ha expirado")
	}
	
	// Si hay otro error (token no encontrado o error de DB)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("Token no encontrado: %s", token)
			
			// Para diagnóstico, buscar cualquier token similar
			var similarTokens []models.PasswordReset
			config.DB.Where("token LIKE ?", token[:8]+"%").Find(&similarTokens)
			if len(similarTokens) > 0 {
				log.Printf("Tokens similares encontrados: %d", len(similarTokens))
				for i, t := range similarTokens {
					log.Printf("Token similar %d: %s (para email: %s, usado: %v)", 
						i+1, t.Token, t.Email, t.Used)
				}
			}
			
			return nil, utils.ErrInvalidToken
		}
		log.Printf("Error de base de datos al validar token: %v", result.Error)
		return nil, result.Error
	}
	
	log.Printf("Token válido para email: %s, expira en: %v", reset.Email, reset.ExpiresAt)
	return &reset, nil
}

// ResetPassword restablece la contraseña del usuario
func (s *PasswordResetService) ResetPassword(token string, newPassword string) error {
	log.Printf("Intentando restablecer contraseña con token: %s", token)
	
	// Validar el token primero
	reset, err := s.ValidateResetToken(token)
	if err != nil {
		log.Printf("Error al validar token: %v", err)
		return err
	}
	
	log.Printf("Token validado correctamente para email: %s", reset.Email)

	// Buscar el usuario por email
	var user models.Usuario
	if result := config.DB.Where("email = ?", reset.Email).First(&user); result.Error != nil {
		log.Printf("Usuario no encontrado: %v", result.Error)
		return utils.ErrUserNotFound
	}
	
	log.Printf("Usuario encontrado: ID=%d, Nombre=%s", user.ID, user.Nombre)

	// Hash de la nueva contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error al generar hash de la contraseña: %v", err)
		return errors.New("error al procesar la contraseña")
	}
	
	log.Println("Hash de contraseña generado correctamente")

	// Actualizar la contraseña del usuario
	if result := config.DB.Model(&user).Update("password", string(hashedPassword)); result.Error != nil {
		log.Printf("Error al actualizar contraseña en la base de datos: %v", result.Error)
		return utils.ErrDatabaseError
	}
	
	log.Printf("Contraseña actualizada correctamente para usuario ID: %d", user.ID)

	// Marcar el token como usado - Usar la instancia real de la base de datos
	if result := config.DB.Model(&models.PasswordReset{}).Where("token = ?", token).Update("used", true); result.Error != nil {
		log.Printf("Error al marcar token como usado: %v", result.Error)
		return utils.ErrDatabaseError
	}
	
	log.Printf("Token marcado como usado correctamente: %s", token)

	return nil
}

// SendPasswordResetEmail envía un correo electrónico con el enlace de recuperación
func (s *PasswordResetService) SendPasswordResetEmail(email, name, resetLink string) error {
	log.Printf("Preparando envío de email de recuperación a: %s", email)
	log.Printf("Enlace de recuperación: %s", resetLink)
	
	// Preparar datos para la plantilla
	data := struct {
		Name      string
		ResetLink string
	}{
		Name:      name,
		ResetLink: resetLink,
	}

	// Renderizar la plantilla de correo
	htmlContent, err := utils.RenderEmailTemplate("reset_password", data)
	if err != nil {
		log.Printf("Error al renderizar plantilla de correo: %v", err)
		return err
	}
	
	log.Println("Plantilla de correo renderizada correctamente")

	// Configuración de envío de correo
	smtpHost := utils.GetEnv("SMTP_HOST", "smtp.gmail.com")
	smtpPort := utils.GetEnv("SMTP_PORT", "587")
	emailFrom := utils.GetEnv("EMAIL_FROM", "")
	
	
	log.Printf("Configuración SMTP: Host=%s, Puerto=%s, From=%s", 
		smtpHost, smtpPort, emailFrom)
	
	// Enviar el correo
	subject := "Recuperación de contraseña"
	return utils.SendEmail(email, subject, htmlContent, "text/html; charset=\"UTF-8\"")
}