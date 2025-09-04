package services

import (
	"auth-user-service/models"
	"auth-user-service/utils"
	"log"
	"time"

	"gorm.io/gorm"
)

type ActivityService struct {
	db *gorm.DB
}

func NewActivityService() *ActivityService {
	return &ActivityService{
		db: utils.GetDB(),
	}
}

// LogActivity registra una nueva actividad
func (a *ActivityService) LogActivity(req *models.CreateActivityLogRequest) error {
	activity := models.ActivityLog{
		UserID:    req.UserID,
		Action:    req.Action,
		Details:   req.Details,
		IP:        req.IP,
		UserAgent: req.UserAgent,
		CreatedAt: time.Now(),
	}

	if result := a.db.Create(&activity); result.Error != nil {
		log.Printf("Error al registrar actividad: %v", result.Error)
		return result.Error
	}

	return nil
}

// GetUserActivity obtiene la actividad de un usuario específico
func (a *ActivityService) GetUserActivity(userID uint, page, limit int) ([]models.ActivityLog, int64, error) {
	var activities []models.ActivityLog
	var total int64

	// Contar total
	a.db.Model(&models.ActivityLog{}).Where("user_id = ?", userID).Count(&total)

	// Obtener actividades paginadas
	offset := (page - 1) * limit
	if result := a.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&activities); result.Error != nil {
		return nil, 0, result.Error
	}

	return activities, total, nil
}

// GetAllActivity obtiene toda la actividad del sistema (para admin)
func (a *ActivityService) GetAllActivity(page, limit int, userID *uint) ([]models.ActivityLog, int64, error) {
	var activities []models.ActivityLog
	var total int64

	query := a.db.Model(&models.ActivityLog{})
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	// Contar total
	query.Count(&total)

	// Obtener actividades paginadas
	offset := (page - 1) * limit
	if result := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&activities); result.Error != nil {
		return nil, 0, result.Error
	}

	return activities, total, nil
}

// DeleteOldActivity elimina logs de actividad antiguos (tarea de mantenimiento)
func (a *ActivityService) DeleteOldActivity(daysOld int) error {
	cutoffDate := time.Now().AddDate(0, 0, -daysOld)
	
	result := a.db.Where("created_at < ?", cutoffDate).Delete(&models.ActivityLog{})
	if result.Error != nil {
		return result.Error
	}

	log.Printf("Eliminados %d registros de actividad antiguos", result.RowsAffected)
	return nil
}

// GetTopActions obtiene las acciones más frecuentes (para estadísticas)
func (a *ActivityService) GetTopActions(limit int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	
	if err := a.db.Model(&models.ActivityLog{}).
		Select("action, COUNT(*) as count").
		Group("action").
		Order("count DESC").
		Limit(limit).
		Find(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}