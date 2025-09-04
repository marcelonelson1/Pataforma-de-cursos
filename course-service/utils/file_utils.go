// utils/file_utils.go
package utils

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"course-service/config"
)

// CreateDirIfNotExists crea un directorio si no existe
func CreateDirIfNotExists(dirPath string) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		log.Printf("Creando directorio: %s", dirPath)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			log.Printf("Error al crear directorio %s: %v", dirPath, err)
		}
	}
}

// SaveImageFile guarda un archivo de imagen desde un formulario
func SaveImageFile(c *gin.Context, fieldName string, cfg *config.Config) (string, error) {
	file, header, err := c.Request.FormFile(fieldName)
	if err != nil {
		// Si no hay archivo, no es un error
		if err == http.ErrMissingFile {
			return "", nil
		}
		return "", err
	}
	defer file.Close()

	// Validar tamaño del archivo
	if header.Size > cfg.MaxImageSize {
		return "", fmt.Errorf("la imagen no debe superar los %d MB", cfg.MaxImageSize/(1024*1024))
	}

	// Validar tipo de archivo
	fileExt := strings.ToLower(filepath.Ext(header.Filename))
	if fileExt == "" {
		return "", fmt.Errorf("archivo sin extensión")
	}

	// Quitar el punto de la extensión para comparar
	ext := fileExt[1:]
	validExt := false
	for _, allowedExt := range cfg.AllowedImageFormats {
		if ext == strings.ToLower(allowedExt) {
			validExt = true
			break
		}
	}

	if !validExt {
		return "", fmt.Errorf("tipo de archivo no permitido. Formatos válidos: %s", 
			strings.Join(cfg.AllowedImageFormats, ", "))
	}

	// Generar nombre de archivo único
	uniqueID := uuid.New().String()
	filename := fmt.Sprintf("%s%s", uniqueID, fileExt)

	// Crear directorio de imágenes si no existe
	imageDir := filepath.Join(cfg.UploadPath, "images")
	CreateDirIfNotExists(imageDir)

	// Ruta completa del archivo
	filePath := filepath.Join(imageDir, filename)

	// Crear el archivo en el servidor
	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("error al crear archivo: %v", err)
	}
	defer out.Close()

	// Copiar el contenido del archivo
	if _, err = io.Copy(out, file); err != nil {
		return "", fmt.Errorf("error al copiar archivo: %v", err)
	}

	// URL donde se puede acceder a la imagen
	imageURL := fmt.Sprintf("/static/images/%s", filename)

	log.Printf("Imagen guardada exitosamente: %s", imageURL)
	return imageURL, nil
}

// ValidateVideoFile valida un archivo de video
func ValidateVideoFile(header *multipart.FileHeader, cfg *config.Config) error {
	// Validar tamaño
	if header.Size > cfg.MaxVideoSize {
		return fmt.Errorf("el archivo no debe superar los %d MB", cfg.MaxVideoSize/(1024*1024))
	}

	// Validar extensión
	fileExt := strings.ToLower(filepath.Ext(header.Filename))
	if fileExt == "" {
		return fmt.Errorf("archivo sin extensión")
	}

	// Quitar el punto de la extensión para comparar
	ext := fileExt[1:]
	validExt := false
	for _, allowedExt := range cfg.AllowedVideoFormats {
		if ext == strings.ToLower(allowedExt) {
			validExt = true
			break
		}
	}

	if !validExt {
		return fmt.Errorf("tipo de archivo no permitido. Formatos válidos: %s", 
			strings.Join(cfg.AllowedVideoFormats, ", "))
	}

	return nil
}

// ValidateImageFile valida un archivo de imagen
func ValidateImageFile(header *multipart.FileHeader, cfg *config.Config) error {
	// Validar tamaño
	if header.Size > cfg.MaxImageSize {
		return fmt.Errorf("la imagen no debe superar los %d MB", cfg.MaxImageSize/(1024*1024))
	}

	// Validar extensión
	fileExt := strings.ToLower(filepath.Ext(header.Filename))
	if fileExt == "" {
		return fmt.Errorf("archivo sin extensión")
	}

	// Quitar el punto de la extensión para comparar
	ext := fileExt[1:]
	validExt := false
	for _, allowedExt := range cfg.AllowedImageFormats {
		if ext == strings.ToLower(allowedExt) {
			validExt = true
			break
		}
	}

	if !validExt {
		return fmt.Errorf("tipo de archivo no permitido. Formatos válidos: %s", 
			strings.Join(cfg.AllowedImageFormats, ", "))
	}

	return nil
}

// GetFileExtension obtiene la extensión de un archivo
func GetFileExtension(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}

// GetContentType determina el tipo de contenido basado en la extensión
func GetContentType(filename string) string {
	ext := GetFileExtension(filename)
	
	switch ext {
	case ".mp4":
		return "video/mp4"
	case ".webm":
		return "video/webm"
	case ".ogg":
		return "video/ogg"
	case ".avi":
		return "video/avi"
	case ".mov":
		return "video/quicktime"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}

// DeleteFile elimina un archivo del sistema de archivos
func DeleteFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("archivo no encontrado: %s", filePath)
	}

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("error al eliminar archivo: %v", err)
	}

	log.Printf("Archivo eliminado: %s", filePath)
	return nil
}

// DeleteDirectory elimina un directorio y todo su contenido
func DeleteDirectory(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return nil // El directorio no existe, no es un error
	}

	if err := os.RemoveAll(dirPath); err != nil {
		return fmt.Errorf("error al eliminar directorio: %v", err)
	}

	log.Printf("Directorio eliminado: %s", dirPath)
	return nil
}

// GetFileSize obtiene el tamaño de un archivo
func GetFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

// FileExists verifica si un archivo existe
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// GetCurrentTimestamp retorna el timestamp actual en formato RFC3339
func GetCurrentTimestamp() string {
	return time.Now().Format(time.RFC3339)
}

// GenerateUniqueFilename genera un nombre de archivo único
func GenerateUniqueFilename(originalName string) string {
	ext := filepath.Ext(originalName)
	uniqueID := uuid.New().String()
	return fmt.Sprintf("%s%s", uniqueID, ext)
}

// GenerateUniqueFilenameWithPrefix genera un nombre de archivo único con prefijo
func GenerateUniqueFilenameWithPrefix(prefix, originalName string) string {
	ext := filepath.Ext(originalName)
	uniqueID := uuid.New().String()
	return fmt.Sprintf("%s-%s%s", prefix, uniqueID, ext)
}

// CleanupOldFiles limpia archivos antiguos de un directorio
func CleanupOldFiles(dirPath string, maxAge time.Duration) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && time.Since(info.ModTime()) > maxAge {
			log.Printf("Eliminando archivo antiguo: %s", path)
			if err := os.Remove(path); err != nil {
				log.Printf("Error al eliminar archivo antiguo %s: %v", path, err)
			}
		}

		return nil
	})
}

// GetDirectorySize calcula el tamaño total de un directorio
func GetDirectorySize(dirPath string) (int64, error) {
	var size int64
	
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	return size, err
}

// CopyFile copia un archivo de origen a destino
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return destinationFile.Sync()
}

// SetVideoStreamingHeaders establece headers apropiados para streaming de video
func SetVideoStreamingHeaders(c *gin.Context, filename string) {
	contentType := GetContentType(filename)
	
	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", filename))
	c.Header("Cache-Control", "public, max-age=31536000") // 1 año
	c.Header("Accept-Ranges", "bytes") // Permitir rangos para seeking
}

// SetImageHeaders establece headers para imágenes
func SetImageHeaders(c *gin.Context, filename string) {
	contentType := GetContentType(filename)
	
	c.Header("Content-Type", contentType)
	c.Header("Cache-Control", "public, max-age=31536000") // 1 año
}