package services

import (
	"home-service/models"
	"home-service/utils"
	"log"

	"gorm.io/gorm"
)

type HomeService struct {
	db *gorm.DB
}

// NewHomeService crea una nueva instancia de HomeService
func NewHomeService() *HomeService {
	return &HomeService{
		db: utils.GetDB(),
	}
}

// GetAllImages obtiene todas las imagenes activas (publico)
func (hs *HomeService) GetAllImages() ([]models.HomeImage, error) {
	var images []models.HomeImage

	if err := hs.db.Where("is_active = ?", true).Order("image_order ASC, created_at DESC").Find(&images).Error; err != nil {
		return nil, utils.ErrDatabaseError
	}

	return images, nil
}

// GetAllImagesAdmin obtiene todas las imagenes para admin
func (hs *HomeService) GetAllImagesAdmin() ([]models.HomeImage, error) {
	var images []models.HomeImage

	if err := hs.db.Order("image_order ASC, created_at DESC").Find(&images).Error; err != nil {
		return nil, utils.ErrDatabaseError
	}

	return images, nil
}

// GetImageByID obtiene una imagen por ID
func (hs *HomeService) GetImageByID(id uint) (*models.HomeImage, error) {
	var image models.HomeImage

	if err := hs.db.First(&image, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrResourceNotFound
		}
		return nil, utils.ErrDatabaseError
	}

	return &image, nil
}

// CreateImage crea una nueva imagen del home
func (hs *HomeService) CreateImage(req *models.HomeImageRequest) (*models.HomeImage, error) {
	// Si no se especifica orden, obtener el siguiente disponible
	if req.Order == 0 {
		var maxOrder int
		hs.db.Model(&models.HomeImage{}).Select("COALESCE(MAX(image_order), 0)").Scan(&maxOrder)
		req.Order = maxOrder + 1
	}

	image := &models.HomeImage{
		Title:    req.Title,
		Subtitle: req.Subtitle,
		Order:    req.Order,
		IsActive: req.IsActive,
	}

	if err := hs.db.Create(image).Error; err != nil {
		log.Printf("Error creando imagen: %v", err)
		return nil, utils.ErrDatabaseError
	}

	return image, nil
}

// UpdateImage actualiza una imagen existente
func (hs *HomeService) UpdateImage(id uint, req *models.HomeImageRequest) (*models.HomeImage, error) {
	image, err := hs.GetImageByID(id)
	if err != nil {
		return nil, err
	}

	// Actualizar campos
	image.UpdateFromRequest(req)

	if err := hs.db.Save(image).Error; err != nil {
		log.Printf("Error actualizando imagen: %v", err)
		return nil, utils.ErrDatabaseError
	}

	return image, nil
}

// DeleteImage elimina una imagen
func (hs *HomeService) DeleteImage(id uint) error {
	image, err := hs.GetImageByID(id)
	if err != nil {
		return err
	}

	// Eliminar archivo asociado si existe
	if image.ImageURL != "" {
		imagePath := utils.GetImagePath(image.ImageURL)
		if err := utils.DeleteFile(imagePath); err != nil {
			log.Printf("Error eliminando imagen %s: %v", imagePath, err)
		}
	}

	if err := hs.db.Delete(image).Error; err != nil {
		return utils.ErrDatabaseError
	}

	// Reordenar imagenes restantes
	hs.reorderImagesAfterDeletion(image.Order)

	return nil
}

// UpdateImageFile actualiza el archivo de una imagen
func (hs *HomeService) UpdateImageFile(id uint, imageURL string) (*models.HomeImage, error) {
	image, err := hs.GetImageByID(id)
	if err != nil {
		return nil, err
	}

	// Eliminar imagen anterior si existe
	if image.ImageURL != "" && image.ImageURL != imageURL {
		oldImagePath := utils.GetImagePath(image.ImageURL)
		if err := utils.DeleteFile(oldImagePath); err != nil {
			log.Printf("Error eliminando imagen anterior %s: %v", oldImagePath, err)
		}
	}

	image.SetImageURL(imageURL)

	if err := hs.db.Save(image).Error; err != nil {
		return nil, utils.ErrDatabaseError
	}

	return image, nil
}

// ReorderImages reordena las imagenes segun una lista de IDs
func (hs *HomeService) ReorderImages(imageIDs []uint) error {
	// Validar que todos los IDs existan
	var existingCount int64
	hs.db.Model(&models.HomeImage{}).Where("id IN ?", imageIDs).Count(&existingCount)
	
	if int(existingCount) != len(imageIDs) {
		return utils.ErrResourceNotFound
	}

	// Actualizar orden en transaccion
	return hs.db.Transaction(func(tx *gorm.DB) error {
		for i, imageID := range imageIDs {
			if err := tx.Model(&models.HomeImage{}).
				Where("id = ?", imageID).
				Update("image_order", i+1).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetHomeStats obtiene estadisticas de las imagenes del home
func (hs *HomeService) GetHomeStats() (*models.HomeStats, error) {
	var totalImages, activeImages, inactiveImages int64

	// Contar imagenes totales
	hs.db.Model(&models.HomeImage{}).Count(&totalImages)
	
	// Contar imagenes activas
	hs.db.Model(&models.HomeImage{}).Where("is_active = ?", true).Count(&activeImages)
	
	// Contar imagenes inactivas
	inactiveImages = totalImages - activeImages

	// Imagenes recientes
	var recentImages []models.HomeImage
	hs.db.Where("is_active = ?", true).
		Order("created_at DESC").
		Limit(5).
		Find(&recentImages)

	// Convertir a respuestas
	var recentResponses []models.HomeImageResponse
	for _, image := range recentImages {
		recentResponses = append(recentResponses, image.ToResponse())
	}

	stats := &models.HomeStats{
		TotalImages:    int(totalImages),
		ActiveImages:   int(activeImages),
		InactiveImages: int(inactiveImages),
		RecentImages:   recentResponses,
	}

	return stats, nil
}

// ToggleImageStatus cambia el estado activo/inactivo de una imagen
func (hs *HomeService) ToggleImageStatus(id uint) (*models.HomeImage, error) {
	image, err := hs.GetImageByID(id)
	if err != nil {
		return nil, err
	}

	if image.IsActive {
		image.Deactivate()
	} else {
		image.Activate()
	}

	if err := hs.db.Save(image).Error; err != nil {
		return nil, utils.ErrDatabaseError
	}

	return image, nil
}

// reorderImagesAfterDeletion reordena imagenes despues de eliminar una
func (hs *HomeService) reorderImagesAfterDeletion(deletedOrder int) {
	// Decrementar el orden de todas las imagenes que tenian un orden mayor
	hs.db.Model(&models.HomeImage{}).
		Where("image_order > ?", deletedOrder).
		UpdateColumn("image_order", gorm.Expr("image_order - 1"))
}