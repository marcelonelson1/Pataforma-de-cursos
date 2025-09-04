package models

import (
	"time"
)

// ProgresoUsuario representa el progreso general de un usuario en un curso
type ProgresoUsuario struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UsuarioID       uint      `gorm:"not null;index:idx_usuario_curso" json:"usuario_id"`
	CursoID         uint      `gorm:"not null;index:idx_usuario_curso" json:"curso_id"`
	PorcentajeTotal float64   `gorm:"type:decimal(5,2);default:0" json:"porcentaje_total"`
	UltimoCapitulo  uint      `gorm:"default:0" json:"ultimo_capitulo"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// TableName especifica el nombre de la tabla para el modelo ProgresoUsuario
func (ProgresoUsuario) TableName() string {
	return "progreso_usuarios"
}

// ProgresoCapitulo representa el progreso de un usuario en un capítulo específico
type ProgresoCapitulo struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UsuarioID  uint      `gorm:"not null;index:idx_usuario_capitulo" json:"usuario_id"`
	CursoID    uint      `gorm:"not null;index:idx_usuario_capitulo" json:"curso_id"`
	CapituloID uint      `gorm:"not null;index:idx_usuario_capitulo" json:"capitulo_id"`
	Completado bool      `gorm:"default:false" json:"completado"`
	Progreso   float64   `gorm:"type:decimal(5,2);default:0" json:"progreso"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TableName especifica el nombre de la tabla para el modelo ProgresoCapitulo
func (ProgresoCapitulo) TableName() string {
	return "progreso_capitulos"
}

// ProgresoResponse representa la respuesta con información de progreso para el frontend
type ProgresoResponse struct {
	ProgresoTotal     float64                   `json:"progreso_total"`
	UltimoCapitulo    uint                      `json:"ultimo_capitulo"`
	CapitulosProgreso map[uint]ProgresoCapitulo `json:"capitulos_progreso"`
}