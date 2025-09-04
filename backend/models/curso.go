package models

import (
	"time"
)

// Curso representa un curso en la plataforma
type Curso struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Titulo      string     `gorm:"size:200;not null" json:"titulo"`
	Descripcion string     `gorm:"size:500" json:"descripcion"`
	Contenido   string     `gorm:"type:text" json:"contenido"`
	Precio      float64    `gorm:"type:decimal(10,2);default:29.99" json:"precio"`
	Estado      string     `gorm:"size:20;default:'Borrador'" json:"estado"`
	ImagenURL   string     `gorm:"size:255" json:"imagen_url"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Capitulos   []Capitulo `gorm:"foreignKey:CursoID" json:"capitulos,omitempty"`
}

// TableName especifica el nombre de la tabla para el modelo Curso
func (Curso) TableName() string {
	return "cursos"
}