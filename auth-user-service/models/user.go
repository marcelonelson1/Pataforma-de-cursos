package models

import "time"

// Usuario modelo principal de usuario
type Usuario struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Nombre    string    `gorm:"size:100;not null" json:"nombre"`
	Email     string    `gorm:"size:100;not null;uniqueIndex" json:"email"`
	Password  string    `gorm:"size:100;not null" json:"-"`
	Role      string    `gorm:"size:20;default:'user'" json:"role"`
	Phone     string    `gorm:"size:20" json:"phone"`
	ImageURL  string    `gorm:"size:255" json:"image_url"`
	LastLogin time.Time `json:"last_login"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserProfile información extendida del perfil
type UserProfile struct {
	UserID            uint   `json:"user_id"`
	Nombre            string `json:"nombre"`
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	ImageURL          string `json:"image_url"`
	Role              string `json:"role"`
	LastLogin         time.Time `json:"last_login"`
	CreatedAt         time.Time `json:"created_at"`
	NotificationSettings *NotificationSettings `json:"notification_settings,omitempty"`
}