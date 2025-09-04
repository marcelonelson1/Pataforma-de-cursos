package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"log"
	"github.com/gin-gonic/gin"
)

type HomeImage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ImageURL  string    `gorm:"size:255" json:"image_url"`
	Title     string    `gorm:"size:100" json:"title"`
	Subtitle  string    `gorm:"size:200" json:"subtitle"`
	Order     int       `gorm:"column:image_order;default:0" json:"order"` // Renombrado para evitar palabra reservada
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName sobrescribe el nombre de la tabla predeterminado
func (HomeImage) TableName() string {
	return "home_images"
}

func initHomeImagesDirs() {
	createDirIfNotExists("./static/home")
}

func registerHomeImageRoutes(router *gin.Engine) {
	// Rutas protegidas que requieren autenticación de administrador
	homeImages := router.Group("/api/home-images")
	homeImages.Use(authMiddleware())
	homeImages.Use(adminMiddleware())
	{
		homeImages.GET("", getHomeImages)
		homeImages.POST("", uploadHomeImage)
		homeImages.PUT("/:id", updateHomeImage)
		homeImages.DELETE("/:id", deleteHomeImage)
		homeImages.PATCH("/reorder", reorderHomeImages)
	}

	// Endpoint público para obtener imágenes activas sin autenticación
	router.GET("/api/home-images/public", getPublicHomeImages)
}

// getPublicHomeImages retorna solo las imágenes activas para uso público
func getPublicHomeImages(c *gin.Context) {
	var images []HomeImage
	if err := db.Where("is_active = ?", true).Order("image_order ASC").Find(&images).Error; err != nil {
		log.Printf("Error al obtener imágenes públicas: %v", err)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}
	SendSuccessResponse(c, images)
}

func getHomeImages(c *gin.Context) {
	var images []HomeImage
	if err := db.Order("image_order ASC").Find(&images).Error; err != nil {
		log.Printf("Error al obtener todas las imágenes: %v", err)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}
	SendSuccessResponse(c, images)
}

func uploadHomeImage(c *gin.Context) {
	// Verificar autenticación
	userValue, exists := c.Get("user")
	if !exists {
		log.Println("Error en uploadHomeImage: Usuario no encontrado en el contexto")
		SendErrorResponse(c, errors.New("usuario no autorizado"), http.StatusUnauthorized)
		return
	}
	
	user, ok := userValue.(Usuario)
	if !ok || user.Role != "admin" {
		log.Printf("Error en uploadHomeImage: Usuario no es admin, role: %s", user.Role)
		SendErrorResponse(c, errors.New("permisos insuficientes"), http.StatusForbidden)
		return
	}

	// Procesar archivo de imagen
	file, err := c.FormFile("image")
	if err != nil {
		log.Printf("Error al leer el archivo de la solicitud: %v", err)
		SendErrorResponse(c, errors.New("no se pudo leer el archivo de imagen"), http.StatusBadRequest)
		return
	}

	// Validar tipo de archivo
	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		SendErrorResponse(c, errors.New("solo se permiten archivos de imagen"), http.StatusBadRequest)
		return
	}

	// Crear directorio si no existe
	if err := os.MkdirAll("./static/home", 0755); err != nil {
		log.Printf("Error al crear directorio: %v", err)
		SendErrorResponse(c, ErrServerError, http.StatusInternalServerError)
		return
	}

	// Generar nombre único para el archivo
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := fmt.Sprintf("./static/home/%s", filename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		log.Printf("Error al guardar el archivo: %v", err)
		SendErrorResponse(c, ErrServerError, http.StatusInternalServerError)
		return
	}

	title := c.PostForm("title")
	subtitle := c.PostForm("subtitle")

	// Obtener el último orden para añadir la nueva imagen al final
	var lastOrder int
	db.Model(&HomeImage{}).Select("COALESCE(MAX(image_order), -1)").Scan(&lastOrder)
	newOrder := lastOrder + 1

	image := HomeImage{
		ImageURL: fmt.Sprintf("/static/home/%s", filename),
		Title:    title,
		Subtitle: subtitle,
		Order:    newOrder,
		IsActive: true,
	}

	if err := db.Create(&image).Error; err != nil {
		log.Printf("Error al crear registro en DB: %v", err)
		os.Remove(filePath) // Eliminar archivo si falla la BD
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}

	SendSuccessResponse(c, image)
}

func updateHomeImage(c *gin.Context) {
	id := c.Param("id")
	var image HomeImage

	if err := db.First(&image, id).Error; err != nil {
		log.Printf("Error al buscar imagen %s: %v", id, err)
		SendErrorResponse(c, ErrResourceNotFound, http.StatusNotFound)
		return
	}

	var updateData struct {
		Title    string `json:"title"`
		Subtitle string `json:"subtitle"`
		IsActive *bool  `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		log.Printf("Error al procesar JSON: %v", err)
		SendErrorResponse(c, ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	if updateData.Title != "" {
		image.Title = updateData.Title
	}
	if updateData.Subtitle != "" {
		image.Subtitle = updateData.Subtitle
	}
	if updateData.IsActive != nil {
		image.IsActive = *updateData.IsActive
	}

	if err := db.Save(&image).Error; err != nil {
		log.Printf("Error al actualizar imagen: %v", err)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}

	SendSuccessResponse(c, image)
}

func deleteHomeImage(c *gin.Context) {
	id := c.Param("id")
	var image HomeImage

	if err := db.First(&image, id).Error; err != nil {
		log.Printf("Error al buscar imagen %s: %v", id, err)
		SendErrorResponse(c, ErrResourceNotFound, http.StatusNotFound)
		return
	}

	// Eliminar archivo físico
	filePath := fmt.Sprintf(".%s", image.ImageURL)
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		log.Printf("Advertencia al eliminar archivo de imagen: %v", err)
	}

	if err := db.Delete(&image).Error; err != nil {
		log.Printf("Error al eliminar imagen de la BD: %v", err)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}

	SendSuccessResponse(c, gin.H{"message": "Imagen eliminada correctamente"})
}

func reorderHomeImages(c *gin.Context) {
	var newOrder []struct {
		ID    uint `json:"id"`
		Order int  `json:"order"`
	}

	if err := c.ShouldBindJSON(&newOrder); err != nil {
		log.Printf("Error al procesar JSON: %v", err)
		SendErrorResponse(c, ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	tx := db.Begin()
	for _, item := range newOrder {
		if err := tx.Model(&HomeImage{}).Where("id = ?", item.ID).Update("image_order", item.Order).Error; err != nil {
			tx.Rollback()
			log.Printf("Error al actualizar orden de imagen %d: %v", item.ID, err)
			SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
			return
		}
	}
	tx.Commit()

	SendSuccessResponse(c, gin.H{"message": "Orden actualizado correctamente"})
}