// services/file_service.go
package services

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"course-service/config"
	
)

type FileService struct {
	db     *gorm.DB
	config *config.Config
}

type VideoResponse struct {
	VideoURL string `json:"video_url"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

type ImageResponse struct {
	ImageURL string `json:"image_url"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

func NewFileService(db *gorm.DB, cfg *config.Config) *FileService {
	return &FileService{
		db:     db,
		config: cfg,
	}
}

// UploadVideo maneja la subida de videos (migrado de uploadVideo)
func (fs *FileService) UploadVideo(file multipart.File, header *multipart.FileHeader, courseID uint, chapterID string) (*VideoResponse, error) {
	// Generar nombre único
	uniqueID := uuid.New().String()
	fileExt := strings.ToLower(filepath.Ext(header.Filename))
	filename := fmt.Sprintf("%s-%s-%s%s", fmt.Sprintf("%d", courseID), chapterID, uniqueID, fileExt)

	// Crear directorio específico para el curso
	courseDir := filepath.Join(fs.config.UploadPath, "videos", fmt.Sprintf("%d", courseID))
	if err := os.MkdirAll(courseDir, 0755); err != nil {
		return nil, fmt.Errorf("error al crear directorio para el curso: %v", err)
	}

	// Ruta completa del archivo
	filePath := filepath.Join(courseDir, filename)

	// Crear archivo en el servidor
	out, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("error al crear archivo: %v", err)
	}
	defer out.Close()

	// Copiar contenido
	if _, err = io.Copy(out, file); err != nil {
		return nil, fmt.Errorf("error al copiar archivo: %v", err)
	}

	// URL de acceso
	videoURL := fmt.Sprintf("/static/videos/%d/%s", courseID, filename)

	log.Printf("Video subido exitosamente: %s", videoURL)

	return &VideoResponse{
		VideoURL: videoURL,
		Filename: filename,
		Size:     header.Size,
	}, nil
}

// UploadImage maneja la subida de imágenes
func (fs *FileService) UploadImage(file multipart.File, header *multipart.FileHeader) (*ImageResponse, error) {
	// Generar nombre único
	uniqueID := uuid.New().String()
	fileExt := strings.ToLower(filepath.Ext(header.Filename))
	filename := uniqueID + fileExt

	// Crear directorio si no existe
	imageDir := filepath.Join(fs.config.UploadPath, "images")
	if err := os.MkdirAll(imageDir, 0755); err != nil {
		return nil, fmt.Errorf("error al crear directorio de imágenes: %v", err)
	}

	// Ruta completa del archivo
	filePath := filepath.Join(imageDir, filename)

	// Crear archivo en el servidor
	out, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("error al crear archivo: %v", err)
	}
	defer out.Close()

	// Copiar contenido
	if _, err = io.Copy(out, file); err != nil {
		return nil, fmt.Errorf("error al copiar archivo: %v", err)
	}

	// URL de acceso
	imageURL := fmt.Sprintf("/static/images/%s", filename)

	log.Printf("Imagen subida exitosamente: %s", imageURL)

	return &ImageResponse{
		ImageURL: imageURL,
		Filename: filename,
		Size:     header.Size,
	}, nil
}

// DeleteFile elimina un archivo específico (migrado de deleteVideo)
func (fs *FileService) DeleteFile(courseID, filename string) error {
	filePath := filepath.Join(fs.config.UploadPath, "videos", courseID, filename)

	// Verificar si el archivo existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("Archivo no encontrado para eliminar: %s", filePath)
		return fmt.Errorf("archivo no encontrado")
	}

	// Eliminar archivo
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("error al eliminar archivo: %v", err)
	}

	// Eliminar miniatura si existe
	thumbnailPath := filepath.Join(fs.config.UploadPath, "thumbnails", courseID, filename)
	if _, err := os.Stat(thumbnailPath); err == nil {
		os.Remove(thumbnailPath)
	}

	log.Printf("Archivo eliminado exitosamente: %s", filePath)
	return nil
}

// DeleteImageFile elimina un archivo de imagen
func (fs *FileService) DeleteImageFile(imageURL string) error {
	// Extraer nombre del archivo de la URL
	filename := filepath.Base(imageURL)
	filePath := filepath.Join(fs.config.UploadPath, "images", filename)

	// Verificar si el archivo existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("Imagen no encontrada para eliminar: %s", filePath)
		return nil // No es error crítico
	}

	// Eliminar archivo
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("error al eliminar imagen: %v", err)
	}

	log.Printf("Imagen eliminada exitosamente: %s", filePath)
	return nil
}

// DeleteCourseDirectory elimina todo el directorio de un curso
func (fs *FileService) DeleteCourseDirectory(courseID string) error {
	// Eliminar directorio de videos
	videosDir := filepath.Join(fs.config.UploadPath, "videos", courseID)
	if _, err := os.Stat(videosDir); err == nil {
		if err := os.RemoveAll(videosDir); err != nil {
			log.Printf("Error al eliminar directorio de videos: %v", err)
		} else {
			log.Printf("Directorio de videos eliminado: %s", videosDir)
		}
	}

	// Eliminar directorio de miniaturas
	thumbnailsDir := filepath.Join(fs.config.UploadPath, "thumbnails", courseID)
	if _, err := os.Stat(thumbnailsDir); err == nil {
		if err := os.RemoveAll(thumbnailsDir); err != nil {
			log.Printf("Error al eliminar directorio de miniaturas: %v", err)
		} else {
			log.Printf("Directorio de miniaturas eliminado: %s", thumbnailsDir)
		}
	}

	return nil
}

// ValidateVideoFile valida un archivo de video
func (fs *FileService) ValidateVideoFile(header *multipart.FileHeader) error {
	// Validar tamaño
	if header.Size > fs.config.MaxVideoSize {
		return fmt.Errorf("el archivo no debe superar los %d MB", fs.config.MaxVideoSize/(1024*1024))
	}

	// Validar extensión
	fileExt := strings.ToLower(filepath.Ext(header.Filename))
	if fileExt == "" {
		return fmt.Errorf("archivo sin extensión")
	}

	// Quitar el punto de la extensión para comparar
	ext := fileExt[1:]
	validExt := false
	for _, allowedExt := range fs.config.AllowedVideoFormats {
		if ext == strings.ToLower(allowedExt) {
			validExt = true
			break
		}
	}

	if !validExt {
		return fmt.Errorf("tipo de archivo no permitido. Formatos válidos: %s", 
			strings.Join(fs.config.AllowedVideoFormats, ", "))
	}

	return nil
}

// ValidateImageFile valida un archivo de imagen
func (fs *FileService) ValidateImageFile(header *multipart.FileHeader) error {
	// Validar tamaño
	if header.Size > fs.config.MaxImageSize {
		return fmt.Errorf("la imagen no debe superar los %d MB", fs.config.MaxImageSize/(1024*1024))
	}

	// Validar extensión
	fileExt := strings.ToLower(filepath.Ext(header.Filename))
	if fileExt == "" {
		return fmt.Errorf("archivo sin extensión")
	}

	// Quitar el punto de la extensión para comparar
	ext := fileExt[1:]
	validExt := false
	for _, allowedExt := range fs.config.AllowedImageFormats {
		if ext == strings.ToLower(allowedExt) {
			validExt = true
			break
		}
	}

	if !validExt {
		return fmt.Errorf("tipo de archivo no permitido. Formatos válidos: %s", 
			strings.Join(fs.config.AllowedImageFormats, ", "))
	}

	return nil
}

// GenerateThumbnail genera una miniatura para un video (implementación básica)
func (fs *FileService) GenerateThumbnail(videoPath, thumbnailPath string) error {
	// Crear directorio de miniaturas si no existe
	thumbnailDir := filepath.Dir(thumbnailPath)
	if err := os.MkdirAll(thumbnailDir, 0755); err != nil {
		return fmt.Errorf("error al crear directorio de miniaturas: %v", err)
	}

	// Esta es una implementación básica
	// En producción, usarías ffmpeg o similar para generar miniaturas reales
	log.Printf("Generando miniatura para %s -> %s", videoPath, thumbnailPath)
	
	// Por ahora, crear un archivo vacío como placeholder
	file, err := os.Create(thumbnailPath)
	if err != nil {
		return fmt.Errorf("error al crear miniatura: %v", err)
	}
	defer file.Close()

	// Escribir un placeholder (en producción sería la imagen real)
	file.WriteString("thumbnail placeholder")

	return nil
}

// GetFileInfo obtiene información de un archivo
func (fs *FileService) GetFileInfo(coursePath, filename string) (os.FileInfo, error) {
	filePath := filepath.Join(fs.config.UploadPath, coursePath, filename)
	return os.Stat(filePath)
}

// GetDiskUsage obtiene el uso de disco por curso
func (fs *FileService) GetDiskUsage(courseID string) (int64, error) {
	courseDir := filepath.Join(fs.config.UploadPath, "videos", courseID)
	
	var totalSize int64
	err := filepath.Walk(courseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("error al calcular uso de disco: %v", err)
	}

	return totalSize, nil
}

// CleanupOrphanedFiles limpia archivos huérfanos (sin registro en BD)
func (fs *FileService) CleanupOrphanedFiles() error {
	videosDir := filepath.Join(fs.config.UploadPath, "videos")
	
	// Obtener todos los archivos de video registrados en BD
	var chapters []struct {
		CursoID     uint   `json:"curso_id"`
		VideoNombre string `json:"video_nombre"`
	}
	
	if err := fs.db.Model(&struct {
		CursoID     uint   `gorm:"column:curso_id"`
		VideoNombre string `gorm:"column:video_nombre"`
	}{}).
		Select("curso_id, video_nombre").
		Where("video_nombre != ''").
		Find(&chapters).Error; err != nil {
		return fmt.Errorf("error al obtener archivos registrados: %v", err)
	}

	// Crear mapa de archivos válidos
	validFiles := make(map[string]bool)
	for _, chapter := range chapters {
		path := fmt.Sprintf("%d/%s", chapter.CursoID, chapter.VideoNombre)
		validFiles[path] = true
	}

	// Recorrer directorio y eliminar archivos huérfanos
	return filepath.Walk(videosDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		// Obtener ruta relativa
		relPath, err := filepath.Rel(videosDir, path)
		if err != nil {
			return err
		}

		// Verificar si el archivo está registrado
		if !validFiles[relPath] {
			log.Printf("Eliminando archivo huérfano: %s", relPath)
			if err := os.Remove(path); err != nil {
				log.Printf("Error al eliminar archivo huérfano %s: %v", relPath, err)
			}
		}

		return nil
	})
}