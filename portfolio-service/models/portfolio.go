package models

import (
	"time"
	"gorm.io/gorm"
)

// ProjectPortfolio representa un proyecto en el portafolio
type ProjectPortfolio struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"size:200;not null" json:"title"`
	Category    string         `gorm:"size:50;not null" json:"category"`
	Description string         `gorm:"size:500" json:"description"`
	ImageURL    string         `gorm:"size:255" json:"image_url"`
	OrderIndex  int            `gorm:"default:0" json:"order_index"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName especifica el nombre de la tabla
func (ProjectPortfolio) TableName() string {
	return "project_portfolios"
}

// ProjectRequest estructura para crear/actualizar un proyecto
type ProjectRequest struct {
	Title       string `json:"title" binding:"required"`
	Category    string `json:"category" binding:"required"`
	Description string `json:"description"`
	OrderIndex  int    `json:"order_index"`
	IsActive    bool   `json:"is_active"`
}

// ProjectResponse estructura para respuestas de proyectos
type ProjectResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	OrderIndex  int       `json:"order_index"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ReorderRequest estructura para reordenar proyectos
type ReorderRequest struct {
	ProjectIDs []uint `json:"project_ids" binding:"required"`
}

// CategoryStats estructura para estadisticas por categoria
type CategoryStats struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
}

// PortfolioStats estructura para estadisticas del portafolio
type PortfolioStats struct {
	TotalProjects    int             `json:"total_projects"`
	ActiveProjects   int             `json:"active_projects"`
	InactiveProjects int             `json:"inactive_projects"`
	CategoriesStats  []CategoryStats `json:"categories_stats"`
	RecentProjects   []ProjectResponse `json:"recent_projects"`
}

// FileUploadRequest estructura para subida de archivos
type FileUploadRequest struct {
	ProjectID uint   `form:"project_id" binding:"required"`
	Filename  string `form:"filename"`
}

// Constantes de categorias validas
const (
	CategoryWeb       = "web"
	CategoryMobile    = "mobile"
	CategoryDesktop   = "desktop"
	CategoryDesign    = "design"
	CategoryMarketing = "marketing"
	CategoryOther     = "other"
)

// ValidCategories lista de categorias validas
var ValidCategories = []string{
	CategoryWeb,
	CategoryMobile,
	CategoryDesktop,
	CategoryDesign,
	CategoryMarketing,
	CategoryOther,
}

// IsValidCategory verifica si una categoria es valida
func IsValidCategory(category string) bool {
	for _, validCategory := range ValidCategories {
		if category == validCategory {
			return true
		}
	}
	return false
}

// ToResponse convierte el modelo a estructura de respuesta
func (p *ProjectPortfolio) ToResponse() ProjectResponse {
	return ProjectResponse{
		ID:          p.ID,
		Title:       p.Title,
		Category:    p.Category,
		Description: p.Description,
		ImageURL:    p.ImageURL,
		OrderIndex:  p.OrderIndex,
		IsActive:    p.IsActive,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// UpdateFromRequest actualiza el proyecto desde una request
func (p *ProjectPortfolio) UpdateFromRequest(req *ProjectRequest) {
	p.Title = req.Title
	p.Category = req.Category
	p.Description = req.Description
	p.OrderIndex = req.OrderIndex
	p.IsActive = req.IsActive
}

// SetImageURL establece la URL de la imagen
func (p *ProjectPortfolio) SetImageURL(imageURL string) {
	p.ImageURL = imageURL
}

// Activate activa el proyecto
func (p *ProjectPortfolio) Activate() {
	p.IsActive = true
}

// Deactivate desactiva el proyecto
func (p *ProjectPortfolio) Deactivate() {
	p.IsActive = false
}

// SetOrder establece el orden del proyecto
func (p *ProjectPortfolio) SetOrder(order int) {
	p.OrderIndex = order
}