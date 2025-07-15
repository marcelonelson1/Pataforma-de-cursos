package models

import (
	"time"
)

// ActivityLog representa un registro de actividad en el sistema
type ActivityLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Action    string    `gorm:"size:50;not null" json:"action"`
	Details   string    `gorm:"size:500" json:"details"`
	IP        string    `gorm:"size:50" json:"ip"`
	UserAgent string    `gorm:"size:255" json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName especifica el nombre de la tabla para el modelo ActivityLog
func (ActivityLog) TableName() string {
	return "activity_logs"
}

// ActivityLogData representa la respuesta con información detallada de un registro de actividad
type ActivityLogData struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	UserName    string    `json:"user_name"`
	UserSurname string    `json:"user_surname"`
	FullName    string    `json:"full_name"`
	Action      string    `json:"action"`
	Details     string    `json:"details"`
	IP          string    `json:"ip"`
	UserAgent   string    `json:"user_agent"`
	CreatedAt   time.Time `json:"created_at"`
}
