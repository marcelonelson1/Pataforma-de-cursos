package models

import "time"

// NotificationSettings configuraciones de notificaciones del usuario
type NotificationSettings struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	UserID            uint      `gorm:"not null;uniqueIndex" json:"user_id"`
	EmailNotifications bool     `gorm:"default:true" json:"email_notifications"`
	PushNotifications  bool     `gorm:"default:true" json:"push_notifications"`
	CourseUpdates      bool     `gorm:"default:true" json:"course_updates"`
	PaymentAlerts      bool     `gorm:"default:true" json:"payment_alerts"`
	MarketingEmails    bool     `gorm:"default:false" json:"marketing_emails"`
	WeeklyDigest       bool     `gorm:"default:true" json:"weekly_digest"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// UpdateNotificationSettingsRequest estructura para actualizar configuraciones
type UpdateNotificationSettingsRequest struct {
	EmailNotifications bool `json:"email_notifications"`
	PushNotifications  bool `json:"push_notifications"`
	CourseUpdates      bool `json:"course_updates"`
	PaymentAlerts      bool `json:"payment_alerts"`
	MarketingEmails    bool `json:"marketing_emails"`
	WeeklyDigest       bool `json:"weekly_digest"`
}