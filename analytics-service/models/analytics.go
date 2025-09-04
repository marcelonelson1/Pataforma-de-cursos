package models

import (
	"time"
	"gorm.io/gorm"
)

// User representa el perfil de usuario para analytics
type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Nombre   string `json:"nombre" gorm:"size:100;not null"`
	Email    string `json:"email" gorm:"size:100;not null;unique"`
	Role     string `json:"role" gorm:"size:20;default:user"`
	ImageURL string `json:"image_url" gorm:"size:255"`
	
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// SalesStats representa estadísticas de ventas
type SalesStats struct {
	ID             uint    `json:"id" gorm:"primaryKey"`
	Period         string  `json:"period" gorm:"size:20;not null"` // day, week, month, year
	Date           string  `json:"date" gorm:"size:20;not null"`   // YYYY-MM-DD format
	TotalSales     float64 `json:"total_sales" gorm:"default:0"`
	TotalOrders    int     `json:"total_orders" gorm:"default:0"`
	TotalUsers     int     `json:"total_users" gorm:"default:0"`
	TotalCourses   int     `json:"total_courses" gorm:"default:0"`
	TotalContacts  int     `json:"total_contacts" gorm:"default:0"`
	TotalProjects  int     `json:"total_projects" gorm:"default:0"`
	
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// DashboardMetrics representa métricas del dashboard
type DashboardMetrics struct {
	ID               uint    `json:"id" gorm:"primaryKey"`
	Date             string  `json:"date" gorm:"size:20;not null"`
	ActiveUsers      int     `json:"active_users" gorm:"default:0"`
	NewRegistrations int     `json:"new_registrations" gorm:"default:0"`
	Revenue          float64 `json:"revenue" gorm:"default:0"`
	ConversionRate   float64 `json:"conversion_rate" gorm:"default:0"`
	
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// ActivityLog representa el log de actividades
type ActivityLog struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	UserID      uint   `json:"user_id"`
	Action      string `json:"action" gorm:"size:100;not null"`
	Description string `json:"description" gorm:"size:500"`
	IPAddress   string `json:"ip_address" gorm:"size:45"`
	UserAgent   string `json:"user_agent" gorm:"size:255"`
	
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}