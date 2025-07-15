package models

import (
	"time"
)

// Capitulo representa un capítulo dentro de un curso
type Capitulo struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CursoID     uint      `gorm:"not null;index" json:"curso_id"`
	Titulo      string    `gorm:"size:200;not null" json:"titulo"`
	Descripcion string    `gorm:"size:500" json:"descripcion"`
	Duracion    string    `gorm:"size:10" json:"duracion"`
	VideoURL    string    `gorm:"size:255" json:"video_url"`
	VideoNombre string    `gorm:"size:255" json:"video_nombre"`
	Orden       int       `gorm:"default:0" json:"orden"`
	Publicado   bool      `gorm:"default:false" json:"publicado"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Curso       *Curso    `gorm:"foreignKey:CursoID" json:"-"`
}

// TableName especifica el nombre de la tabla para el modelo Capitulo
func (Capitulo) TableName() string {
	return "capitulos"
}