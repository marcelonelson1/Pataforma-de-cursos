package models

import "time"

// ActivityLog modelo para logs de actividad del usuario
type ActivityLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Action    string    `gorm:"size:50;not null" json:"action"`
	Details   string    `gorm:"size:500" json:"details"`
	IP        string    `gorm:"size:50" json:"ip"`
	UserAgent string    `gorm:"size:255" json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateActivityLogRequest estructura para crear logs
type CreateActivityLogRequest struct {
	UserID    uint   `json:"user_id"`
	Action    string `json:"action"`
	Details   string `json:"details"`
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
}