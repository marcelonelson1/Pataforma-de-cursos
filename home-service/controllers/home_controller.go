package controllers

import (
	"home-service/models"
	"home-service/services"
	"home-service/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HomeController struct {
	homeService *services.HomeService
	fileService *services.FileService
}

// NewHomeController crea una nueva instancia de HomeController
func NewHomeController() *HomeController {
	return &HomeController{
		homeService: services.NewHomeService(),
		fileService: services.NewFileService(),
	}
}

// GetAllImages obtiene todas las imagenes activas (publico)
func (hc *HomeController) GetAllImages(c *gin.Context) {
	images, err := hc.homeService.GetAllImages()
	if err != nil {
		utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Convertir a respuestas
	var responses []models.HomeImageResponse
	for _, image := range images {
		responses = append(responses, image.ToResponse())
	}

	utils.SendSuccessResponse(c, gin.H{
		"images": responses,
		"total":  len(responses),
	})
}

// GetAllImagesAdmin obtiene todas las imagenes para admin
func (hc *HomeController) GetAllImagesAdmin(c *gin.Context) {
	images, err := hc.homeService.GetAllImagesAdmin()
	if err != nil {
		utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Convertir a respuestas
	var responses []models.HomeImageResponse
	for _, image := range images {
		responses = append(responses, image.ToResponse())
	}

	utils.SendSuccessResponse(c, gin.H{
		"images": responses,
		"total":  len(responses),
	})
}

// GetImageByID obtiene una imagen especifica
func (hc *HomeController) GetImageByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendErrorMessage(c, "ID invalido", http.StatusBadRequest)
		return
	}

	image, err := hc.homeService.GetImageByID(uint(id))
	if err != nil {
		if err == utils.ErrResourceNotFound {
			utils.SendErrorResponse(c, err, http.StatusNotFound)
		} else {
			utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		}
		return
	}

	utils.SendSuccessResponse(c, image.ToResponse())
}

// CreateImage crea una nueva imagen del home (admin)
func (hc *HomeController) CreateImage(c *gin.Context) {
	var req models.HomeImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, err)
		return
	}

	image, err := hc.homeService.CreateImage(&req)
	if err != nil {
		utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	utils.SendSuccessMessage(c, "Imagen del home creada exitosamente", image.ToResponse())
}

// UpdateImage actualiza una imagen del home (admin)
func (hc *HomeController) UpdateImage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendErrorMessage(c, "ID invalido", http.StatusBadRequest)
		return
	}

	var req models.HomeImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, err)
		return
	}

	image, err := hc.homeService.UpdateImage(uint(id), &req)
	if err != nil {
		if err == utils.ErrResourceNotFound {
			utils.SendErrorResponse(c, err, http.StatusNotFound)
		} else {
			utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		}
		return
	}

	utils.SendSuccessMessage(c, "Imagen del home actualizada exitosamente", image.ToResponse())
}

// DeleteImage elimina una imagen del home (admin)
func (hc *HomeController) DeleteImage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendErrorMessage(c, "ID invalido", http.StatusBadRequest)
		return
	}

	if err := hc.homeService.DeleteImage(uint(id)); err != nil {
		if err == utils.ErrResourceNotFound {
			utils.SendErrorResponse(c, err, http.StatusNotFound)
		} else {
			utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		}
		return
	}

	utils.SendSuccessMessage(c, "Imagen del home eliminada exitosamente", nil)
}

// UploadImage sube archivo para una imagen del home (admin)
func (hc *HomeController) UploadImage(c *gin.Context) {
	// Obtener imagen ID
	imageIDStr := c.PostForm("image_id")
	if imageIDStr == "" {
		utils.SendErrorMessage(c, "image_id requerido", http.StatusBadRequest)
		return
	}

	imageID, err := strconv.ParseUint(imageIDStr, 10, 32)
	if err != nil {
		utils.SendErrorMessage(c, "image_id invalido", http.StatusBadRequest)
		return
	}

	// Obtener archivo
	file, err := c.FormFile("image")
	if err != nil {
		utils.SendErrorMessage(c, "archivo imagen requerido", http.StatusBadRequest)
		return
	}

	// Subir imagen
	imageURL, err := hc.fileService.UploadHomeImage(uint(imageID), file)
	if err != nil {
		if err == utils.ErrResourceNotFound {
			utils.SendErrorResponse(c, err, http.StatusNotFound)
		} else if err == utils.ErrInvalidFileType {
			utils.SendErrorResponse(c, err, http.StatusBadRequest)
		} else {
			utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		}
		return
	}

	utils.SendSuccessMessage(c, "Imagen subida exitosamente", gin.H{
		"image_url": imageURL,
		"image_id":  imageID,
	})
}

// DeleteImageFile elimina archivo de una imagen del home (admin)
func (hc *HomeController) DeleteImageFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendErrorMessage(c, "ID invalido", http.StatusBadRequest)
		return
	}

	if err := hc.fileService.DeleteHomeImage(uint(id)); err != nil {
		if err == utils.ErrResourceNotFound {
			utils.SendErrorResponse(c, err, http.StatusNotFound)
		} else {
			utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		}
		return
	}

	utils.SendSuccessMessage(c, "Archivo de imagen eliminado exitosamente", nil)
}

// ReorderImages reordena imagenes del home (admin)
func (hc *HomeController) ReorderImages(c *gin.Context) {
	var req models.ReorderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, err)
		return
	}

	if err := hc.homeService.ReorderImages(req.ImageIDs); err != nil {
		if err == utils.ErrResourceNotFound {
			utils.SendErrorResponse(c, err, http.StatusNotFound)
		} else {
			utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		}
		return
	}

	utils.SendSuccessMessage(c, "Imagenes reordenadas exitosamente", nil)
}

// GetHomeStats obtiene estadisticas de las imagenes del home (admin)
func (hc *HomeController) GetHomeStats(c *gin.Context) {
	stats, err := hc.homeService.GetHomeStats()
	if err != nil {
		utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, stats)
}

// ToggleImageStatus cambia estado activo/inactivo (admin)
func (hc *HomeController) ToggleImageStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendErrorMessage(c, "ID invalido", http.StatusBadRequest)
		return
	}

	image, err := hc.homeService.ToggleImageStatus(uint(id))
	if err != nil {
		if err == utils.ErrResourceNotFound {
			utils.SendErrorResponse(c, err, http.StatusNotFound)
		} else {
			utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		}
		return
	}

	status := "inactiva"
	if image.IsActive {
		status = "activa"
	}

	utils.SendSuccessMessage(c, "Estado de la imagen actualizado a: "+status, image.ToResponse())
}

// HealthCheck endpoint de verificacion de salud
func (hc *HomeController) HealthCheck(c *gin.Context) {
	utils.SendSuccessResponse(c, gin.H{
		"status":  "ok",
		"service": "home-service",
		"version": "1.0.0",
	})
}