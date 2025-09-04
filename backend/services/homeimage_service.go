package services

import (
	"curso-platform/config"
	"curso-platform/models"
	"curso-platform/utils"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// HomeImageService maneja la lógica relacionada con las imágenes de la página de inicio
type HomeImageService struct{}

// NewHomeImageService crea una nueva instancia de HomeImageService
func NewHomeImageService() *HomeImageService {
	return &HomeImageService{}
}

// GetHomeImages obtiene todas las imágenes de la página de inicio
func (s *HomeImageService) GetHomeImages() ([]models.HomeImage, error) {
	var images []models.HomeImage
	if err := config.DB.Order("image_order ASC").Find(&images).Error; err != nil {
		return nil, utils.ErrDatabaseError
	}
	return images, nil
}

// GetPublicHomeImages obtiene solo las imágenes activas para uso público
func (s *HomeImageService) GetPublicHomeImages() ([]models.HomeImage, error) {
	var images []models.HomeImage
	if err := config.DB.Where("is_active = ?", true).Order("image_order ASC").Find(&images).Error; err != nil {
		return nil, utils.ErrDatabaseError
	}
	return images, nil
}

// UploadHomeImage sube una nueva imagen para la página de inicio
func (s *HomeImageService) UploadHomeImage(title, subtitle string, imageFile *os.File, fileName string) (*models.HomeImage, error) {
	// Crear directorio si no existe
	if err := utils.CreateDirIfNotExists("./static/home"); err != nil {
		return nil, utils.ErrServerError
	}

	// Generar nombre único para el archivo
	ext := filepath.Ext(fileName)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := fmt.Sprintf("./static/home/%s", filename)

	// Leer contenido del archivo original
	fileData, err := os.ReadFile(imageFile.Name())
	if err != nil {
		return nil, fmt.Errorf("error al leer archivo imagen: %v", err)
	}

	// Escribir a nuevo archivo
	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		return nil, fmt.Errorf("error al guardar imagen: %v", err)
	}

	// Obtener el último orden para añadir la nueva imagen al final
	var lastOrder int
	config.DB.Model(&models.HomeImage{}).Select("COALESCE(MAX(image_order), -1)").Scan(&lastOrder)
	newOrder := lastOrder + 1

	image := models.HomeImage{
		ImageURL:  fmt.Sprintf("/static/home/%s", filename),
		Title:     title,
		Subtitle:  subtitle,
		Order:     newOrder,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := config.DB.Create(&image).Error; err != nil {
		// Eliminar archivo si falla la BD
		os.Remove(filePath)
		return nil, utils.ErrDatabaseError
	}

	return &image, nil
}

// UpdateHomeImage actualiza una imagen existente
func (s *HomeImageService) UpdateHomeImage(id uint, title, subtitle string, isActive *bool) (*models.HomeImage, error) {
	var image models.HomeImage
	if err := config.DB.First(&image, id).Error; err != nil {
		return nil, utils.ErrResourceNotFound
	}

	// Actualizar solo los campos proporcionados
	updates := make(map[string]interface{})
	
	if title != "" {
		updates["title"] = title
	}
	if subtitle != "" {
		updates["subtitle"] = subtitle
	}
	if isActive != nil {
		updates["is_active"] = *isActive
	}
	updates["updated_at"] = time.Now()

	if err := config.DB.Model(&image).Updates(updates).Error; err != nil {
		return nil, utils.ErrDatabaseError
	}

	// Obtener la imagen actualizada
	if err := config.DB.First(&image, id).Error; err != nil {
		return nil, utils.ErrDatabaseError
	}

	return &image, nil
}

// DeleteHomeImage elimina una imagen
func (s *HomeImageService) DeleteHomeImage(id uint) error {
	var image models.HomeImage
	if err := config.DB.First(&image, id).Error; err != nil {
		return utils.ErrResourceNotFound
	}

	// Eliminar archivo físico
	filePath := fmt.Sprintf(".%s", image.ImageURL)
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		// Solo registramos el error, no interrumpimos el proceso
		fmt.Printf("Advertencia al eliminar archivo de imagen: %v", err)
	}

	if err := config.DB.Delete(&image).Error; err != nil {
		return utils.ErrDatabaseError
	}

	return nil
}

// ReorderHomeImages reordena las imágenes de la página de inicio
func (s *HomeImageService) ReorderHomeImages(newOrder []struct {
	ID    uint `json:"id"`
	Order int  `json:"order"`
}) error {
	tx := config.DB.Begin()
	
	for _, item := range newOrder {
		if err := tx.Model(&models.HomeImage{}).Where("id = ?", item.ID).Update("image_order", item.Order).Error; err != nil {
			tx.Rollback()
			return utils.ErrDatabaseError
		}
	}
	
	if err := tx.Commit().Error; err != nil {
		return utils.ErrDatabaseError
	}
	
	return nil
}