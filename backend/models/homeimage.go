package models

import (
	"time"
)

// HomeImage representa una imagen para la página de inicio
type HomeImage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ImageURL  string    `gorm:"size:255" json:"image_url"`
	Title     string    `gorm:"size:100" json:"title"`
	Subtitle  string    `gorm:"size:200" json:"subtitle"`
	Order     int       `gorm:"column:image_order;default:0" json:"order"` // Renombrado para evitar palabra reservada
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName sobrescribe el nombre de la tabla predeterminado
func (HomeImage) TableName() string {
	return "home_images"
}