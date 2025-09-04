package models

import (
	"time"

	"gorm.io/gorm"
)

// Enrollment representa una inscripción de un usuario a una actividad
type Enrollment struct {
	gorm.Model
	UserID     uint      `json:"userId" gorm:"not null"`
	ActivityID uint      `json:"activityId" gorm:"not null"`
	EnrolledAt time.Time `json:"enrolledAt" gorm:"not null"`
	Status     string    `json:"status" gorm:"not null;default:'active'"` // active, cancelled, completed
}

// EnrollmentResponse es el modelo de respuesta para una inscripción que incluye detalles de la actividad
type EnrollmentResponse struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"userId"`
	ActivityID uint      `json:"activityId"`
	EnrolledAt time.Time `json:"enrolledAt"`
	Status     string    `json:"status"`
	Activity   Activity  `json:"activity"`
}