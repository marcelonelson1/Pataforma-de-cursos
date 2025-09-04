// controllers/file_controller.go
package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"course-service/config"
	"course-service/services"
	"course-service/utils"
)

type FileController struct {
	db          *gorm.DB
	config      *config.Config
	fileService *services.FileService
}

func NewFileController(db *gorm.DB, cfg *config.Config) *FileController {
	return &FileController{
		db:          db,
		config:      cfg,
		fileService: services.NewFileService(db, cfg),
	}
}

// UploadVideo maneja la subida de videos (migrado de uploadVideo)
func (fc *FileController) UploadVideo(c *gin.Context) {
	// Verificar autenticación
	userID, exists := c.Get("user_id")
	if !exists {
		utils.SendErrorResponse(c, "Usuario no autenticado", http.StatusUnauthorized)
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		utils.SendErrorResponse(c, "Error al obtener información del usuario", http.StatusInternalServerError)
		return
	}

	log.Printf("Usuario %d iniciando subida de video", userIDUint)

	// Obtener parámetros del formulario
	cursoID := c.PostForm("curso_id")
	capituloID := c.PostForm("capitulo_id")

	if cursoID == "" {
		utils.SendErrorResponse(c, "ID del curso requerido", http.StatusBadRequest)
		return
	}

	// Verificar que el curso existe
	cursoIDInt, err := strconv.Atoi(cursoID)
	if err != nil {
		utils.SendErrorResponse(c, "ID de curso inválido", http.StatusBadRequest)
		return
	}

	// Obtener archivo
	file, header, err := c.Request.FormFile("video")
	if err != nil {
		log.Printf("Error al obtener el archivo: %v", err)
		utils.SendErrorResponse(c, "Archivo de video no proporcionado", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validar archivo
	if err := fc.fileService.ValidateVideoFile(header); err != nil {
		log.Printf("Validación de archivo falló: %v", err)
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	// Subir video
	videoResponse, err := fc.fileService.UploadVideo(file, header, uint(cursoIDInt), capituloID)
	if err != nil {
		log.Printf("Error al subir video: %v", err)
		utils.SendErrorResponse(c, "Error al guardar el video: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Video subido exitosamente: %s", videoResponse.VideoURL)
	utils.SendSuccessResponse(c, videoResponse)
}

// UploadImage maneja la subida de imágenes
func (fc *FileController) UploadImage(c *gin.Context) {
	// Verificar autenticación
	userID, exists := c.Get("user_id")
	if !exists {
		utils.SendErrorResponse(c, "Usuario no autenticado", http.StatusUnauthorized)
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		utils.SendErrorResponse(c, "Error al obtener información del usuario", http.StatusInternalServerError)
		return
	}

	log.Printf("Usuario %d iniciando subida de imagen", userIDUint)

	// Obtener archivo
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		log.Printf("Error al obtener el archivo: %v", err)
		utils.SendErrorResponse(c, "Archivo de imagen no proporcionado", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validar archivo
	if err := fc.fileService.ValidateImageFile(header); err != nil {
		log.Printf("Validación de archivo falló: %v", err)
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	// Subir imagen
	imageResponse, err := fc.fileService.UploadImage(file, header)
	if err != nil {
		log.Printf("Error al subir imagen: %v", err)
		utils.SendErrorResponse(c, "Error al guardar la imagen: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Imagen subida exitosamente: %s", imageResponse.ImageURL)
	utils.SendSuccessResponse(c, imageResponse)
}

// GetVideo maneja el streaming de videos (migrado de getVideo)
func (fc *FileController) GetVideo(c *gin.Context) {
	cursoID := c.Param("courseId")
	filename := c.Param("filename")

	if cursoID == "" || filename == "" {
		utils.SendErrorResponse(c, "Parámetros incompletos", http.StatusBadRequest)
		return
	}

	// Verificar que el usuario tiene acceso al curso (si está autenticado)
	userID, exists := c.Get("user_id")
	if exists {
		userIDUint, ok := userID.(uint)
		if ok {
			cursoIDUint, err := strconv.ParseUint(cursoID, 10, 32)
			if err == nil {
				// Verificar acceso solo si pudimos parsear el ID del curso
				courseService := services.NewCourseService(fc.db, fc.config)
				hasAccess := courseService.CheckUserAccessToCourse(uint(cursoIDUint), userIDUint)
				if !hasAccess {
					utils.SendErrorResponse(c, "No tienes acceso a este contenido", http.StatusForbidden)
					return
				}
			}
		}
	}

	// Construir ruta del archivo
	filePath := filepath.Join(fc.config.UploadPath, "videos", cursoID, filename)

	// Verificar si el archivo existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("Video no encontrado: %s", filePath)
		utils.SendErrorResponse(c, "Video no encontrado", http.StatusNotFound)
		return
	}

	// Establecer headers para streaming
	fc.setVideoHeaders(c, filename)

	// Servir el archivo
	c.File(filePath)
}

// GetImage maneja el servicio de imágenes
func (fc *FileController) GetImage(c *gin.Context) {
	filename := c.Param("filename")

	if filename == "" {
		utils.SendErrorResponse(c, "Nombre de archivo requerido", http.StatusBadRequest)
		return
	}

	// Construir ruta del archivo
	filePath := filepath.Join(fc.config.UploadPath, "images", filename)

	// Verificar si el archivo existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("Imagen no encontrada: %s", filePath)
		utils.SendErrorResponse(c, "Imagen no encontrada", http.StatusNotFound)
		return
	}

	// Establecer headers para caché
	c.Header("Cache-Control", "public, max-age=31536000") // 1 año
	c.Header("Content-Type", fc.getImageContentType(filename))

	// Servir el archivo
	c.File(filePath)
}

// DeleteFile elimina un archivo (migrado de deleteVideo)
func (fc *FileController) DeleteFile(c *gin.Context) {
	cursoID := c.Param("courseId")
	filename := c.Param("filename")

	if cursoID == "" || filename == "" {
		utils.SendErrorResponse(c, "Parámetros incompletos", http.StatusBadRequest)
		return
	}

	// Eliminar archivo
	err := fc.fileService.DeleteFile(cursoID, filename)
	if err != nil {
		log.Printf("Error al eliminar archivo: %v", err)
		utils.SendErrorResponse(c, "Error al eliminar archivo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Archivo eliminado exitosamente: %s/%s", cursoID, filename)
	utils.SendSuccessResponse(c, gin.H{"message": "Archivo eliminado correctamente"})
}

// GetThumbnail obtiene la miniatura de un video
func (fc *FileController) GetThumbnail(c *gin.Context) {
	cursoID := c.Param("courseId")
	filename := c.Param("filename")

	if cursoID == "" || filename == "" {
		utils.SendErrorResponse(c, "Parámetros incompletos", http.StatusBadRequest)
		return
	}

	// Construir ruta de la miniatura
	thumbnailPath := filepath.Join(fc.config.UploadPath, "thumbnails", cursoID, filename)

	// Verificar si existe la miniatura
	if _, err := os.Stat(thumbnailPath); os.IsNotExist(err) {
		// Si no existe, intentar generar una
		videoPath := filepath.Join(fc.config.UploadPath, "videos", cursoID, filename)
		if _, err := os.Stat(videoPath); err == nil {
			// Generar miniatura (implementación básica)
			err := fc.fileService.GenerateThumbnail(videoPath, thumbnailPath)
			if err != nil {
				log.Printf("Error al generar miniatura: %v", err)
				utils.SendErrorResponse(c, "Miniatura no disponible", http.StatusNotFound)
				return
			}
		} else {
			utils.SendErrorResponse(c, "Miniatura no encontrada", http.StatusNotFound)
			return
		}
	}

	// Establecer headers
	c.Header("Cache-Control", "public, max-age=86400") // 1 día
	c.Header("Content-Type", "image/jpeg")

	// Servir miniatura
	c.File(thumbnailPath)
}

// setVideoHeaders establece los headers apropiados para streaming de video
func (fc *FileController) setVideoHeaders(c *gin.Context, filename string) {
	ext := filepath.Ext(filename)
	
	var contentType string
	switch ext {
	case ".mp4":
		contentType = "video/mp4"
	case ".webm":
		contentType = "video/webm"
	case ".ogg":
		contentType = "video/ogg"
	default:
		contentType = "video/mp4"
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", filename))
	c.Header("Cache-Control", "public, max-age=31536000") // 1 año
	c.Header("Accept-Ranges", "bytes") // Permitir rangos para seeking
}

// getImageContentType determina el tipo de contenido para imágenes
func (fc *FileController) getImageContentType(filename string) string {
	ext := filepath.Ext(filename)
	
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "image/jpeg"
	}
}