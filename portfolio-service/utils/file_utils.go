package utils

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// AllowedImageTypes tipos de imagen permitidos
var AllowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// MaxFileSize tamano maximo de archivo (5MB)
const MaxFileSize = 5 * 1024 * 1024

// SaveUploadedFile guarda un archivo subido
func SaveUploadedFile(file *multipart.FileHeader, destinationDir string) (string, error) {
	log.Printf("🔥 DEBUG SaveUploadedFile: Starting with file: %s, size: %d", file.Filename, file.Size)
	
	// Validar tamano del archivo
	if file.Size > MaxFileSize {
		log.Printf("🔥 DEBUG SaveUploadedFile: File too large: %d > %d", file.Size, MaxFileSize)
		return "", fmt.Errorf("archivo demasiado grande. Maximo permitido: %d bytes", MaxFileSize)
	}
	log.Printf("🔥 DEBUG SaveUploadedFile: File size OK: %d", file.Size)

	// Validar tipo de archivo
	contentType := file.Header.Get("Content-Type")
	log.Printf("🔥 DEBUG SaveUploadedFile: Content-Type: %s", contentType)
	if !AllowedImageTypes[contentType] {
		log.Printf("🔥 DEBUG SaveUploadedFile: Invalid content type: %s", contentType)
		return "", ErrInvalidFileType
	}
	log.Printf("🔥 DEBUG SaveUploadedFile: Content-Type OK")

	// Crear directorio si no existe
	log.Printf("🔥 DEBUG SaveUploadedFile: Creating directory: %s", destinationDir)
	if err := os.MkdirAll(destinationDir, 0755); err != nil {
		log.Printf("🔥 DEBUG SaveUploadedFile: Failed to create directory: %v", err)
		return "", fmt.Errorf("error creando directorio: %v", err)
	}
	log.Printf("🔥 DEBUG SaveUploadedFile: Directory created successfully")

	// Generar nombre unico para el archivo
	ext := filepath.Ext(file.Filename)
	uniqueFilename := fmt.Sprintf("%s-%d%s", uuid.New().String(), time.Now().Unix(), ext)
	log.Printf("🔥 DEBUG SaveUploadedFile: Generated filename: %s", uniqueFilename)
	
	// Ruta completa del archivo
	filePath := filepath.Join(destinationDir, uniqueFilename)
	log.Printf("🔥 DEBUG SaveUploadedFile: Full file path: %s", filePath)

	// Abrir archivo subido
	log.Printf("🔥 DEBUG SaveUploadedFile: Opening source file...")
	src, err := file.Open()
	if err != nil {
		log.Printf("🔥 DEBUG SaveUploadedFile: Failed to open source file: %v", err)
		return "", fmt.Errorf("error abriendo archivo: %v", err)
	}
	defer src.Close()
	log.Printf("🔥 DEBUG SaveUploadedFile: Source file opened successfully")

	// Crear archivo destino
	log.Printf("🔥 DEBUG SaveUploadedFile: Creating destination file: %s", filePath)
	dst, err := os.Create(filePath)
	if err != nil {
		log.Printf("🔥 DEBUG SaveUploadedFile: Failed to create destination file: %v", err)
		return "", fmt.Errorf("error creando archivo: %v", err)
	}
	defer dst.Close()
	log.Printf("🔥 DEBUG SaveUploadedFile: Destination file created successfully")

	// Copiar contenido
	log.Printf("🔥 DEBUG SaveUploadedFile: Starting file copy...")
	bytes, err := io.Copy(dst, src)
	if err != nil {
		log.Printf("🔥 DEBUG SaveUploadedFile: Failed to copy file content: %v", err)
		return "", fmt.Errorf("error copiando archivo: %v", err)
	}
	log.Printf("🔥 DEBUG SaveUploadedFile: File copied successfully, %d bytes written", bytes)

	log.Printf("🔥 DEBUG SaveUploadedFile: SUCCESS! File saved as: %s", uniqueFilename)
	return uniqueFilename, nil
}

// DeleteFile elimina un archivo del sistema
func DeleteFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // El archivo no existe, no hay error
	}

	return os.Remove(filePath)
}

// GetFileURL construye la URL publica del archivo
func GetFileURL(filename string, baseURL string) string {
	if filename == "" {
		return ""
	}

	// Si ya es una URL completa, retornarla tal como esta
	if strings.HasPrefix(filename, "http://") || strings.HasPrefix(filename, "https://") {
		return filename
	}

	// Construir URL relativa
	return fmt.Sprintf("%s/static/portfolio/%s", baseURL, filename)
}

// ValidateImageFile valida que el archivo sea una imagen valida
func ValidateImageFile(file *multipart.FileHeader) error {
	// Validar tamano
	if file.Size > MaxFileSize {
		return fmt.Errorf("archivo demasiado grande. Maximo: %d bytes", MaxFileSize)
	}

	// Validar tipo MIME
	contentType := file.Header.Get("Content-Type")
	if !AllowedImageTypes[contentType] {
		return ErrInvalidFileType
	}

	// Validar extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	validExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	if !validExtensions[ext] {
		return ErrInvalidFileType
	}

	return nil
}

// GetImagePath retorna la ruta completa de una imagen
func GetImagePath(filename string) string {
	if filename == "" {
		return ""
	}
	return filepath.Join("./static/portfolio", filename)
}

// EnsureStaticDirs crea los directorios estaticos necesarios
func EnsureStaticDirs() error {
	dirs := []string{
		"./static",
		"./static/portfolio",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("error creando directorio %s: %v", dir, err)
		}
	}

	return nil
}