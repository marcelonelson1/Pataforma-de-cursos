package models

import (
	"time"
	"gorm.io/gorm"
)

// HomeImage representa una imagen del home/carousel
type HomeImage struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ImageURL  string         `gorm:"size:255" json:"image_url"`
	Title     string         `gorm:"size:100" json:"title"`
	Subtitle  string         `gorm:"size:200" json:"subtitle"`
	Order     int            `gorm:"column:image_order;default:0" json:"order"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName especifica el nombre de la tabla
func (HomeImage) TableName() string {
	return "home_images"
}

// HomeImageRequest estructura para crear/actualizar una imagen del home
type HomeImageRequest struct {
	Title    string `json:"title" binding:"required"`
	Subtitle string `json:"subtitle"`
	Order    int    `json:"order"`
	IsActive bool   `json:"is_active"`
}

// HomeImageResponse estructura para respuestas de imagenes del home
type HomeImageResponse struct {
	ID        uint      `json:"id"`
	ImageURL  string    `json:"image_url"`
	Title     string    `json:"title"`
	Subtitle  string    `json:"subtitle"`
	Order     int       `json:"order"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ReorderRequest estructura para reordenar imagenes
type ReorderRequest struct {
	ImageIDs []uint `json:"image_ids" binding:"required"`
}

// HomeStats estructura para estadisticas del home
type HomeStats struct {
	TotalImages    int                 `json:"total_images"`
	ActiveImages   int                 `json:"active_images"`
	InactiveImages int                 `json:"inactive_images"`
	RecentImages   []HomeImageResponse `json:"recent_images"`
}

// FileUploadRequest estructura para subida de archivos
type FileUploadRequest struct {
	ImageID  uint   `form:"image_id"`
	Filename string `form:"filename"`
}

// ToResponse convierte el modelo a estructura de respuesta
func (h *HomeImage) ToResponse() HomeImageResponse {
	return HomeImageResponse{
		ID:        h.ID,
		ImageURL:  h.ImageURL,
		Title:     h.Title,
		Subtitle:  h.Subtitle,
		Order:     h.Order,
		IsActive:  h.IsActive,
		CreatedAt: h.CreatedAt,
		UpdatedAt: h.UpdatedAt,
	}
}

// UpdateFromRequest actualiza la imagen desde una request
func (h *HomeImage) UpdateFromRequest(req *HomeImageRequest) {
	h.Title = req.Title
	h.Subtitle = req.Subtitle
	h.Order = req.Order
	h.IsActive = req.IsActive
}

// SetImageURL establece la URL de la imagen
func (h *HomeImage) SetImageURL(imageURL string) {
	h.ImageURL = imageURL
}

// Activate activa la imagen
func (h *HomeImage) Activate() {
	h.IsActive = true
}

// Deactivate desactiva la imagen
func (h *HomeImage) Deactivate() {
	h.IsActive = false
}

// SetOrder establece el orden de la imagen
func (h *HomeImage) SetOrder(order int) {
	h.Order = order
}