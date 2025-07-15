package services

import (
	"curso-platform/models"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// EnrollmentService gestiona la lógica de negocio para inscripciones
type EnrollmentService struct {
	db *gorm.DB
}

// NewEnrollmentService crea una nueva instancia del servicio de inscripciones
func NewEnrollmentService(db *gorm.DB) *EnrollmentService {
	return &EnrollmentService{
		db: db,
	}
}

// Enroll inscribe a un usuario en una actividad
func (s *EnrollmentService) Enroll(userID uint, activityID uint) (*models.Enrollment, error) {
	var activity models.Activity
	
	// Verificar si la actividad existe
	if err := s.db.First(&activity, activityID).Error; err != nil {
		return nil, fmt.Errorf("actividad no encontrada: %w", err)
	}

	// Verificar si hay cupo disponible
	if activity.Enrolled >= activity.Capacity {
		return nil, errors.New("no hay cupo disponible para esta actividad")
	}

	// Verificar si el usuario ya está inscrito
	var existingCount int64
	s.db.Model(&models.Enrollment{}).
		Where("user_id = ? AND activity_id = ? AND status = 'active'", userID, activityID).
		Count(&existingCount)
	
	if existingCount > 0 {
		return nil, errors.New("ya estás inscrito en esta actividad")
	}

	// Iniciar transacción
	tx := s.db.Begin()

	// Crear la inscripción
	enrollment := models.Enrollment{
		UserID:     userID,
		ActivityID: activityID,
		EnrolledAt: time.Now(),
		Status:     "active",
	}

	// Guardar la inscripción
	if err := tx.Create(&enrollment).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error al crear la inscripción: %w", err)
	}

	// Incrementar el contador de inscripciones
	if err := tx.Model(&models.Activity{}).
		Where("id = ?", activityID).
		Update("enrolled", gorm.Expr("enrolled + ?", 1)).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error al actualizar contador de inscripciones: %w", err)
	}

	// Confirmar transacción
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error al confirmar la inscripción: %w", err)
	}

	return &enrollment, nil
}

// GetUserActivities obtiene las actividades en las que está inscrito el usuario
func (s *EnrollmentService) GetUserActivities(userID uint) ([]models.Activity, error) {
	var activities []models.Activity
	
	if err := s.db.Table("activities").
		Joins("JOIN enrollments ON activities.id = enrollments.activity_id").
		Where("enrollments.user_id = ? AND enrollments.status = 'active'", userID).
		Find(&activities).Error; err != nil {
		return nil, fmt.Errorf("error al obtener actividades del usuario: %w", err)
	}
	
	return activities, nil
}