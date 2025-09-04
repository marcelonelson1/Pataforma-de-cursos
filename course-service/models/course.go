// models/course.go
package models

import (
	"time"
	"gorm.io/gorm"
)

// Course representa un curso en el sistema (migrado del modelo Curso original)
type Course struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Titulo         string         `gorm:"size:200;not null" json:"titulo"`
	Descripcion    string         `gorm:"size:500" json:"descripcion"`
	Contenido      string         `gorm:"type:text" json:"contenido"`
	Precio         float64        `gorm:"type:decimal(10,2);default:29.99" json:"precio"`
	Estado         string         `gorm:"size:20;default:'Borrador'" json:"estado"`
	ImagenURL      string         `gorm:"size:255" json:"imagen_url"`
	InstructorID   uint           `gorm:"not null;index" json:"instructor_id"`
	DuracionTotal  string         `gorm:"size:20" json:"duracion_total"`
	Nivel          string         `gorm:"size:20;default:'Principiante'" json:"nivel"`
	CategoriaID    uint           `gorm:"index" json:"categoria_id"`
	Destacado      bool           `gorm:"default:false" json:"destacado"`
	Gratuito       bool           `gorm:"default:false" json:"gratuito"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	
	// Relaciones
	Capitulos []Chapter `gorm:"foreignKey:CursoID" json:"capitulos,omitempty"`
	Categoria *Category `gorm:"foreignKey:CategoriaID" json:"categoria,omitempty"`
}

// TableName especifica el nombre de la tabla para Course
func (Course) TableName() string {
	return "cursos"
}

// Chapter representa un capítulo de un curso (migrado del modelo Capitulo original)
type Chapter struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	CursoID       uint           `gorm:"not null;index" json:"curso_id"`
	Titulo        string         `gorm:"size:200;not null" json:"titulo"`
	Descripcion   string         `gorm:"size:500" json:"descripcion"`
	Duracion      string         `gorm:"size:10" json:"duracion"`
	VideoURL      string         `gorm:"size:255" json:"video_url"`
	VideoNombre   string         `gorm:"size:255" json:"video_nombre"`
	Orden         int            `gorm:"default:0" json:"orden"`
	Publicado     bool           `gorm:"default:false" json:"publicado"`
	TipoContenido string         `gorm:"size:50;default:'video'" json:"tipo_contenido"`
	Gratis        bool           `gorm:"default:false" json:"gratis"` // Preview gratuito
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	
	// Relación
	Curso *Course `gorm:"foreignKey:CursoID" json:"-"`
}

// TableName especifica el nombre de la tabla para Chapter
func (Chapter) TableName() string {
	return "capitulos"
}

// Category representa una categoría de cursos
type Category struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Nombre      string         `gorm:"size:100;not null" json:"nombre"`
	Descripcion string         `gorm:"size:255" json:"descripcion"`
	Slug        string         `gorm:"size:100;uniqueIndex" json:"slug"`
	PadreID     *uint          `gorm:"index" json:"padre_id"`
	Activa      bool           `gorm:"default:true" json:"activa"`
	Orden       int            `gorm:"default:0" json:"orden"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	
	// Relaciones
	Padre      *Category  `gorm:"foreignKey:PadreID" json:"padre,omitempty"`
	Subcategorias []Category `gorm:"foreignKey:PadreID" json:"subcategorias,omitempty"`
	Cursos     []Course   `gorm:"foreignKey:CategoriaID" json:"-"`
}

// CourseRequest estructura para crear/actualizar cursos
type CourseRequest struct {
	Titulo       string  `json:"titulo" binding:"required"`
	Descripcion  string  `json:"descripcion" binding:"required"`
	Contenido    string  `json:"contenido" binding:"required"`
	Precio       float64 `json:"precio"`
	Estado       string  `json:"estado"`
	ImagenURL    string  `json:"imagen_url"`
	Nivel        string  `json:"nivel"`
	CategoriaID  uint    `json:"categoria_id"`
	Destacado    bool    `json:"destacado"`
	Gratuito     bool    `json:"gratuito"`
}

// ChapterRequest estructura para crear/actualizar capítulos
type ChapterRequest struct {
	CursoID       uint   `json:"curso_id" binding:"required"`
	Titulo        string `json:"titulo" binding:"required"`
	Descripcion   string `json:"descripcion"`
	Duracion      string `json:"duracion"`
	VideoURL      string `json:"video_url"`
	VideoNombre   string `json:"video_nombre"`
	Publicado     bool   `json:"publicado"`
	Orden         int    `json:"orden"`
	TipoContenido string `json:"tipo_contenido"`
	Gratis        bool   `json:"gratis"`
}

// CourseResponse estructura para respuestas de curso
type CourseResponse struct {
	ID            uint      `json:"id"`
	Titulo        string    `json:"titulo"`
	Descripcion   string    `json:"descripcion"`
	Precio        float64   `json:"precio"`
	Estado        string    `json:"estado"`
	ImagenURL     string    `json:"imagen_url"`
	Nivel         string    `json:"nivel"`
	Gratuito      bool      `json:"gratuito"`
	TieneAcceso   bool      `json:"tiene_acceso"`
	TotalChapters int       `json:"total_capitulos"`
	CreatedAt     time.Time `json:"created_at"`
}

// ChapterResponse estructura para respuestas de capítulo
type ChapterResponse struct {
	ID            uint      `json:"id"`
	Titulo        string    `json:"titulo"`
	Descripcion   string    `json:"descripcion"`
	Duracion      string    `json:"duracion"`
	VideoURL      string    `json:"video_url"`
	Orden         int       `json:"orden"`
	Publicado     bool      `json:"publicado"`
	TipoContenido string    `json:"tipo_contenido"`
	Gratis        bool      `json:"gratis"`
	TieneAcceso   bool      `json:"tiene_acceso"`
	CreatedAt     time.Time `json:"created_at"`
}

// Constantes para estados de curso
const (
	CourseStatusDraft     = "Borrador"
	CourseStatusPublished = "Publicado"
	CourseStatusArchived  = "Archivado"
	CourseStatusReview    = "En Revisión"
)

// Constantes para niveles de curso
const (
	CourseLevelBeginner     = "Principiante"
	CourseLevelIntermediate = "Intermedio"
	CourseLevelAdvanced     = "Avanzado"
	CourseLevelExpert       = "Experto"
)

// Constantes para tipos de contenido
const (
	ContentTypeVideo = "video"
	ContentTypeText  = "texto"
	ContentTypeQuiz  = "quiz"
	ContentTypeFile  = "archivo"
)

// ValidCourseStates mapa de estados válidos
var ValidCourseStates = map[string]bool{
	CourseStatusDraft:     true,
	CourseStatusPublished: true,
	CourseStatusArchived:  true,
	CourseStatusReview:    true,
}

// ValidCourseLevels mapa de niveles válidos
var ValidCourseLevels = map[string]bool{
	CourseLevelBeginner:     true,
	CourseLevelIntermediate: true,
	CourseLevelAdvanced:     true,
	CourseLevelExpert:       true,
}

// IsValidCourseState valida si un estado de curso es válido
func IsValidCourseState(state string) bool {
	return ValidCourseStates[state]
}

// IsValidCourseLevel valida si un nivel de curso es válido
func IsValidCourseLevel(level string) bool {
	return ValidCourseLevels[level]
}

// IsPublished verifica si el curso está publicado
func (c *Course) IsPublished() bool {
	return c.Estado == CourseStatusPublished
}

// IsFree verifica si el curso es gratuito
func (c *Course) IsFree() bool {
	return c.Gratuito || c.Precio == 0
}

// GetPublishedChapters obtiene solo los capítulos publicados
func (c *Course) GetPublishedChapters() []Chapter {
	var publishedChapters []Chapter
	for _, chapter := range c.Capitulos {
		if chapter.Publicado {
			publishedChapters = append(publishedChapters, chapter)
		}
	}
	return publishedChapters
}

// CanUserAccess verifica si un usuario puede acceder al capítulo
func (ch *Chapter) CanUserAccess(hasCoursePaid bool) bool {
	// Si el capítulo es gratis, siempre puede acceder
	if ch.Gratis {
		return true
	}
	
	// Si el capítulo no está publicado, no puede acceder
	if !ch.Publicado {
		return false
	}
	
	// Necesita haber pagado el curso
	return hasCoursePaid
}