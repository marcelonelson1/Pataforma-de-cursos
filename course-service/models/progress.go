// models/progress.go
package models

import (
	"time"
	"gorm.io/gorm"
)

// UserProgress representa el progreso de un usuario en un curso (migrado de ProgresoUsuario)
type UserProgress struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	UsuarioID        uint           `gorm:"not null;index:idx_usuario_curso" json:"usuario_id"`
	CursoID          uint           `gorm:"not null;index:idx_usuario_curso" json:"curso_id"`
	PorcentajeTotal  float64        `gorm:"type:decimal(5,2);default:0" json:"porcentaje_total"`
	UltimoCapitulo   uint           `gorm:"default:0" json:"ultimo_capitulo"`
	FechaInicio      *time.Time     `json:"fecha_inicio"`
	FechaCompletado  *time.Time     `json:"fecha_completado"`
	TiempoTotal      int            `gorm:"default:0" json:"tiempo_total"` // en segundos
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName especifica el nombre de la tabla para UserProgress
func (UserProgress) TableName() string {
	return "progreso_usuarios"
}

// ChapterProgress representa el progreso de un usuario en un capítulo específico (migrado de ProgresoCapitulo)
type ChapterProgress struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	UsuarioID       uint           `gorm:"not null;index:idx_usuario_capitulo" json:"usuario_id"`
	CursoID         uint           `gorm:"not null;index:idx_usuario_capitulo" json:"curso_id"`
	CapituloID      uint           `gorm:"not null;index:idx_usuario_capitulo" json:"capitulo_id"`
	Completado      bool           `gorm:"default:false" json:"completado"`
	Progreso        float64        `gorm:"type:decimal(5,2);default:0" json:"progreso"` // 0-100
	TiempoVisto     int            `gorm:"default:0" json:"tiempo_visto"` // en segundos
	FechaCompletado *time.Time     `json:"fecha_completado"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName especifica el nombre de la tabla para ChapterProgress
func (ChapterProgress) TableName() string {
	return "progreso_capitulos"
}

// ProgressResponse estructura para respuestas de progreso (migrado de ProgresoResponse)
type ProgressResponse struct {
	ProgresoTotal     float64                   `json:"progreso_total"`
	UltimoCapitulo    uint                      `json:"ultimo_capitulo"`
	TiempoTotal       int                       `json:"tiempo_total"`
	FechaInicio       *time.Time                `json:"fecha_inicio"`
	FechaCompletado   *time.Time                `json:"fecha_completado"`
	CapitulosProgreso map[uint]ChapterProgress  `json:"capitulos_progreso"`
	EstadisticasCurso CourseStatistics          `json:"estadisticas_curso"`
}

// CourseStatistics estadísticas del curso
type CourseStatistics struct {
	TotalCapitulos      int `json:"total_capitulos"`
	CapitulosCompletados int `json:"capitulos_completados"`
	CapitulosPendientes int `json:"capitulos_pendientes"`
}

// ChapterProgressRequest estructura para actualizar progreso de capítulo (migrado de ProgresoCapituloRequest)
type ChapterProgressRequest struct {
	CursoID    uint    `json:"curso_id" binding:"required"`
	CapituloID uint    `json:"capitulo_id" binding:"required"`
	Completado bool    `json:"completado"`
	Progreso   float64 `json:"progreso"`
	TiempoVisto int    `json:"tiempo_visto"`
}

// LastChapterRequest estructura para actualizar último capítulo visto (migrado de UltimoCapituloRequest)
type LastChapterRequest struct {
	CursoID    uint `json:"curso_id" binding:"required"`
	CapituloID uint `json:"capitulo_id" binding:"required"`
}

// UserProgressSummary resumen del progreso del usuario
type UserProgressSummary struct {
	TotalCursos          int                `json:"total_cursos"`
	CursosCompletados    int                `json:"cursos_completados"`
	CursosEnProgreso     int                `json:"cursos_en_progreso"`
	TiempoTotalEstudio   int                `json:"tiempo_total_estudio"`
	PromedioCompletitud  float64            `json:"promedio_completitud"`
	CursosRecientes      []CourseProgress   `json:"cursos_recientes"`
}

// CourseProgress progreso en un curso específico
type CourseProgress struct {
	CursoID         uint      `json:"curso_id"`
	TituloCurso     string    `json:"titulo_curso"`
	ImagenURL       string    `json:"imagen_url"`
	Progreso        float64   `json:"progreso"`
	UltimoCapitulo  uint      `json:"ultimo_capitulo"`
	FechaUltimaVez  time.Time `json:"fecha_ultima_vez"`
}

// IsCompleted verifica si el progreso está completado
func (up *UserProgress) IsCompleted() bool {
	return up.PorcentajeTotal >= 100.0
}

// MarkAsStarted marca el curso como iniciado
func (up *UserProgress) MarkAsStarted() {
	if up.FechaInicio == nil {
		now := time.Now()
		up.FechaInicio = &now
	}
}

// MarkAsCompleted marca el curso como completado
func (up *UserProgress) MarkAsCompleted() {
	if up.PorcentajeTotal >= 100.0 && up.FechaCompletado == nil {
		now := time.Now()
		up.FechaCompletado = &now
	}
}

// UpdateTotalProgress actualiza el progreso total basado en capítulos
func (up *UserProgress) UpdateTotalProgress(totalChapters, completedChapters int) {
	if totalChapters > 0 {
		up.PorcentajeTotal = (float64(completedChapters) / float64(totalChapters)) * 100
		
		// Marcar como completado si llegó al 100%
		if up.PorcentajeTotal >= 100.0 {
			up.MarkAsCompleted()
		}
	}
}

// MarkAsCompleted marca el capítulo como completado
func (cp *ChapterProgress) MarkAsCompleted() {
	cp.Completado = true
	cp.Progreso = 100.0
	now := time.Now()
	cp.FechaCompletado = &now
}

// UpdateProgress actualiza el progreso del capítulo
func (cp *ChapterProgress) UpdateProgress(progress float64, timeWatched int) {
	cp.Progreso = progress
	cp.TiempoVisto = timeWatched
	
	// Si llegó al 100%, marcarlo como completado
	if progress >= 100.0 {
		cp.MarkAsCompleted()
	}
}

// GetProgressPercentage calcula el porcentaje basado en capítulos completados
func GetProgressPercentage(totalChapters, completedChapters int) float64 {
	if totalChapters == 0 {
		return 0.0
	}
	return (float64(completedChapters) / float64(totalChapters)) * 100
}