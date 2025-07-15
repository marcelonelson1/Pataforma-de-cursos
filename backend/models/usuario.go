package models

import (
	"time"
)

// Usuario representa un usuario en el sistema
type Usuario struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Nombre    string    `gorm:"size:50;not null" json:"nombre"`
	Apellido  string    `gorm:"size:50;not null" json:"apellido"`
	Email     string    `gorm:"size:100;not null;uniqueIndex" json:"email"`
	Password  string    `gorm:"size:100;not null" json:"-"`
	Role      string    `gorm:"size:20;default:'user'" json:"role"`
	Phone     string    `gorm:"size:20" json:"phone"`
	ImageURL  string    `gorm:"size:255" json:"image_url"`
	LastLogin time.Time `json:"last_login"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName especifica el nombre de la tabla para el modelo Usuario
func (Usuario) TableName() string {
	return "usuarios"
}

// GetFullName devuelve el nombre completo del usuario
func (u Usuario) GetFullName() string {
	return u.Nombre + " " + u.Apellido
}
