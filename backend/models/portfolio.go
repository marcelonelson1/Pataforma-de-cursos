package models

import (
	"time"
)

// ProjectPortfolio representa un proyecto en el portafolio
type ProjectPortfolio struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"size:200;not null" json:"title"`
	Category    string    `gorm:"size:50;not null" json:"category"`
	Description string    `gorm:"size:500" json:"description"`
	ImageURL    string    `gorm:"size:255" json:"image_url"`
	OrderIndex  int       `gorm:"default:0" json:"order_index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName especifica el nombre de la tabla para el modelo ProjectPortfolio
func (ProjectPortfolio) TableName() string {
	return "project_portfolios"
}