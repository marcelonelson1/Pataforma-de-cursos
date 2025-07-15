package controllers

import (
	"curso-platform/middleware"
	"curso-platform/models"
	"curso-platform/services"
	"curso-platform/utils"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// HomeImageController gestiona las operaciones relacionadas con las imágenes de la página de inicio
type HomeImageController struct {
	homeImageService *services.HomeImageService
}

// NewHomeImageController crea una nueva instancia del controlador de imágenes de inicio
func NewHomeImageController(homeImageService *services.HomeImageService) *HomeImageController {
	return &HomeImageController{
		homeImageService: homeImageService,
	}
}

// GetHomeImages obtiene todas las imágenes de la página de inicio
func (c *HomeImageController) GetHomeImages(ctx *gin.Context) {
	images, err := c.homeImageService.GetHomeImages()
	if err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}
	
	utils.SendSuccessResponse(ctx, images)
}

// GetPublicHomeImages obtiene solo las imágenes activas para uso público
func (c *HomeImageController) GetPublicHomeImages(ctx *gin.Context) {
	images, err := c.homeImageService.GetPublicHomeImages()
	if err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}
	
	utils.SendSuccessResponse(ctx, images)
}

// UploadHomeImage sube una nueva imagen para la página de inicio
func (c *HomeImageController) UploadHomeImage(ctx *gin.Context) {
	// Verificar autenticación
	userValue, exists := ctx.Get("user")
	if !exists {
		utils.SendErrorResponse(ctx, utils.ErrUnauthorized, http.StatusUnauthorized)
		return
	}
	
	user, ok := userValue.(models.Usuario)
	if !ok || user.Role != "admin" {
		log.Printf("Error en uploadHomeImage: Usuario no es admin, role: %s", user.Role)
		utils.SendErrorResponse(ctx, utils.ErrUnauthorized, http.StatusForbidden)
		return
	}

	// Procesar archivo de imagen
	file, err := ctx.FormFile("image")
	if err != nil {
		utils.SendErrorResponse(ctx, utils.ErrBadRequest, http.StatusBadRequest)
		return
	}

	// Crear archivo temporal para procesar
	tempFile, err := os.CreateTemp("", "homeimage-*.tmp")
	if err != nil {
		utils.SendErrorResponse(ctx, utils.ErrServerError, http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()
	
	// Guardar FormFile en archivo temporal
	if err := ctx.SaveUploadedFile(file, tempFile.Name()); err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	title := ctx.PostForm("title")
	subtitle := ctx.PostForm("subtitle")

	// Subir imagen
	image, err := c.homeImageService.UploadHomeImage(title, subtitle, tempFile, file.Filename)
	if err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(ctx, image)
}

// UpdateHomeImage actualiza una imagen existente
func (c *HomeImageController) UpdateHomeImage(ctx *gin.Context) {
	id := ctx.Param("id")
	imageID, err := strconv.Atoi(id)
	if err != nil {
		utils.SendErrorResponse(ctx, utils.ErrBadRequest, http.StatusBadRequest)
		return
	}

	var updateData struct {
		Title    string `json:"title"`
		Subtitle string `json:"subtitle"`
		IsActive *bool  `json:"is_active"`
	}

	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		utils.SendErrorResponse(ctx, utils.ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	image, err := c.homeImageService.UpdateHomeImage(uint(imageID), updateData.Title, updateData.Subtitle, updateData.IsActive)
	if err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(ctx, image)
}

// DeleteHomeImage elimina una imagen
func (c *HomeImageController) DeleteHomeImage(ctx *gin.Context) {
	id := ctx.Param("id")
	imageID, err := strconv.Atoi(id)
	if err != nil {
		utils.SendErrorResponse(ctx, utils.ErrBadRequest, http.StatusBadRequest)
		return
	}

	if err := c.homeImageService.DeleteHomeImage(uint(imageID)); err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(ctx, gin.H{"message": "Imagen eliminada correctamente"})
}

// ReorderHomeImages reordena las imágenes de la página de inicio
func (c *HomeImageController) ReorderHomeImages(ctx *gin.Context) {
	var newOrder []struct {
		ID    uint `json:"id"`
		Order int  `json:"order"`
	}

	if err := ctx.ShouldBindJSON(&newOrder); err != nil {
		utils.SendErrorResponse(ctx, utils.ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	if err := c.homeImageService.ReorderHomeImages(newOrder); err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(ctx, gin.H{"message": "Orden actualizado correctamente"})
}

// RegisterRoutes registra todas las rutas relacionadas con las imágenes de inicio
func (c *HomeImageController) RegisterRoutes(router *gin.Engine) {
	// Endpoint público para obtener imágenes activas sin autenticación
	router.GET("/api/home-images/public", c.GetPublicHomeImages)

	// Rutas protegidas que requieren autenticación de administrador
	homeImages := router.Group("/api/home-images")
	homeImages.Use(middleware.AuthMiddleware())
	homeImages.Use(middleware.AdminMiddleware())
	{
		homeImages.GET("", c.GetHomeImages)
		homeImages.POST("", c.UploadHomeImage)
		homeImages.PUT("/:id", c.UpdateHomeImage)
		homeImages.DELETE("/:id", c.DeleteHomeImage)
		homeImages.PATCH("/reorder", c.ReorderHomeImages)
	}
}