package services

import (
	"auth-user-service/models"
	"auth-user-service/utils"
	"log"
	"time"

	"gorm.io/gorm"
)

// Aliases para evitar repetir models. todo el tiempo
type (
	Usuario = models.Usuario
	NotificationSettings = models.NotificationSettings
	UpdateNotificationSettingsRequest = models.UpdateNotificationSettingsRequest
)

type UserService struct {
	db *gorm.DB
}

func NewUserService() *UserService {
	return &UserService{
		db: utils.GetDB(),
	}
}

// GetUserByID obtiene un usuario por ID
func (u *UserService) GetUserByID(userID uint) (*models.Usuario, error) {
	var user models.Usuario
	if result := u.db.First(&user, userID); result.Error != nil {
		return nil, utils.ErrUserNotFound
	}
	
	// No devolver la contraseña
	user.Password = ""
	return &user, nil
}

// GetUserProfile obtiene el perfil completo del usuario
func (u *UserService) GetUserProfile(userID uint) (*models.UserProfile, error) {
	var user models.Usuario
	if result := u.db.First(&user, userID); result.Error != nil {
		return nil, utils.ErrUserNotFound
	}

	// Obtener configuraciones de notificación
	var notificationSettings NotificationSettings
	if result := u.db.Where("user_id = ?", userID).First(&notificationSettings); result.Error != nil {
		// Si no existen, crear configuraciones por defecto
		notificationSettings = NotificationSettings{
			UserID:            userID,
			EmailNotifications: true,
			PushNotifications:  true,
			CourseUpdates:      true,
			PaymentAlerts:      true,
			MarketingEmails:    false,
			WeeklyDigest:       true,
		}
		u.db.Create(&notificationSettings)
	}

	profile := &models.UserProfile{
		UserID:              user.ID,
		Nombre:              user.Nombre,
		Email:               user.Email,
		Phone:               user.Phone,
		ImageURL:            user.ImageURL,
		Role:                user.Role,
		LastLogin:           user.LastLogin,
		CreatedAt:           user.CreatedAt,
		NotificationSettings: &notificationSettings,
	}

	return profile, nil
}

// UpdateProfile actualiza el perfil del usuario
func (u *UserService) UpdateProfile(userID uint, req *models.UpdateProfileRequest) error {
	var user models.Usuario
	if result := u.db.First(&user, userID); result.Error != nil {
		return utils.ErrUserNotFound
	}

	// Actualizar campos
	updates := make(map[string]interface{})
	if req.Nombre != "" {
		updates["nombre"] = req.Nombre
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}

	if len(updates) > 0 {
		if result := u.db.Model(&user).Updates(updates); result.Error != nil {
			return result.Error
		}

		// Log de actividad
		u.LogActivity(userID, "update_profile", "Perfil actualizado", "", "")
	}

	return nil
}

// UpdateProfileImage actualiza la imagen de perfil
func (u *UserService) UpdateProfileImage(userID uint, imageURL string) error {
	var user models.Usuario
	if result := u.db.First(&user, userID); result.Error != nil {
		return utils.ErrUserNotFound
	}

	if result := u.db.Model(&user).Update("image_url", imageURL); result.Error != nil {
		return result.Error
	}

	// Log de actividad
	u.LogActivity(userID, "update_profile_image", "Imagen de perfil actualizada", "", "")

	return nil
}

// GetNotificationSettings obtiene las configuraciones de notificación
func (u *UserService) GetNotificationSettings(userID uint) (*NotificationSettings, error) {
	var settings NotificationSettings
	if result := u.db.Where("user_id = ?", userID).First(&settings); result.Error != nil {
		// Si no existen, crear por defecto
		settings = NotificationSettings{
			UserID:            userID,
			EmailNotifications: true,
			PushNotifications:  true,
			CourseUpdates:      true,
			PaymentAlerts:      true,
			MarketingEmails:    false,
			WeeklyDigest:       true,
		}
		u.db.Create(&settings)
	}

	return &settings, nil
}

// UpdateNotificationSettings actualiza las configuraciones de notificación
func (u *UserService) UpdateNotificationSettings(userID uint, req *UpdateNotificationSettingsRequest) error {
	var settings NotificationSettings
	result := u.db.Where("user_id = ?", userID).First(&settings)
	
	if result.Error != nil {
		// Si no existen, crear nuevas
		settings = NotificationSettings{
			UserID:            userID,
			EmailNotifications: req.EmailNotifications,
			PushNotifications:  req.PushNotifications,
			CourseUpdates:      req.CourseUpdates,
			PaymentAlerts:      req.PaymentAlerts,
			MarketingEmails:    req.MarketingEmails,
			WeeklyDigest:       req.WeeklyDigest,
		}
		return u.db.Create(&settings).Error
	} else {
		// Actualizar existentes
		updates := map[string]interface{}{
			"email_notifications": req.EmailNotifications,
			"push_notifications":  req.PushNotifications,
			"course_updates":      req.CourseUpdates,
			"payment_alerts":      req.PaymentAlerts,
			"marketing_emails":    req.MarketingEmails,
			"weekly_digest":       req.WeeklyDigest,
		}
		return u.db.Model(&settings).Updates(updates).Error
	}
}

// ListUsers lista usuarios (para admin)
func (u *UserService) ListUsers(page, limit int) ([]models.Usuario, int64, error) {
	var users []models.Usuario
	var total int64

	// Contar total
	u.db.Model(&models.Usuario{}).Count(&total)

	// Obtener usuarios paginados
	offset := (page - 1) * limit
	if result := u.db.Select("id, nombre, email, role, phone, image_url, last_login, created_at, updated_at").
		Offset(offset).Limit(limit).Find(&users); result.Error != nil {
		return nil, 0, result.Error
	}

	return users, total, nil
}

// UpdateUser actualiza un usuario (para admin)
func (u *UserService) UpdateUser(userID uint, updates map[string]interface{}) error {
	var user models.Usuario
	if result := u.db.First(&user, userID); result.Error != nil {
		return utils.ErrUserNotFound
	}

	if result := u.db.Model(&user).Updates(updates); result.Error != nil {
		return result.Error
	}

	// Log de actividad
	u.LogActivity(userID, "admin_update", "Usuario actualizado por admin", "", "")

	return nil
}

// DeleteUser elimina un usuario (para admin)
func (u *UserService) DeleteUser(userID uint) error {
	var user models.Usuario
	if result := u.db.First(&user, userID); result.Error != nil {
		return utils.ErrUserNotFound
	}

	// Eliminar en transacción
	tx := u.db.Begin()

	// Eliminar configuraciones de notificación
	tx.Where("user_id = ?", userID).Delete(&models.NotificationSettings{})

	// Eliminar logs de actividad
	tx.Where("user_id = ?", userID).Delete(&models.ActivityLog{})

	// Eliminar tokens de reset de contraseña
	tx.Where("email = ?", user.Email).Delete(&models.PasswordReset{})

	// Eliminar usuario
	if result := tx.Delete(&user); result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	tx.Commit()

	return nil
}

// ChangeUserRole cambia el rol de un usuario (para admin)
func (u *UserService) ChangeUserRole(userID uint, newRole string) error {
	var user models.Usuario
	if result := u.db.First(&user, userID); result.Error != nil {
		return utils.ErrUserNotFound
	}

	if result := u.db.Model(&user).Update("role", newRole); result.Error != nil {
		return result.Error
	}

	// Log de actividad
	u.LogActivity(userID, "role_change", "Rol cambiado a: "+newRole, "", "")

	return nil
}

// LogActivity registra una actividad del usuario
func (u *UserService) LogActivity(userID uint, action, details, ip, userAgent string) {
	activity := models.ActivityLog{
		UserID:    userID,
		Action:    action,
		Details:   details,
		IP:        ip,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
	}

	if result := u.db.Create(&activity); result.Error != nil {
		log.Printf("Error al registrar actividad: %v", result.Error)
	}
}

// GetActivityLog obtiene el log de actividades (para admin)
func (u *UserService) GetActivityLog(page, limit int, userID *uint) ([]models.ActivityLog, int64, error) {
	var activities []models.ActivityLog
	var total int64

	query := u.db.Model(&models.ActivityLog{})
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	// Contar total
	query.Count(&total)

	// Obtener actividades paginadas
	offset := (page - 1) * limit
	if result := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&activities); result.Error != nil {
		return nil, 0, result.Error
	}

	return activities, total, nil
}