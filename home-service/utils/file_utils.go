package utils

import (
	"fmt"
	"io"
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

// MaxFileSize tamano maximo de archivo (10MB para imagenes del home)
const MaxFileSize = 10 * 1024 * 1024

// SaveUploadedFile guarda un archivo subido
func SaveUploadedFile(file *multipart.FileHeader, destinationDir string) (string, error) {
	// Validar tamano del archivo
	if file.Size > MaxFileSize {
		return "", fmt.Errorf("archivo demasiado grande. Maximo permitido: %d bytes", MaxFileSize)
	}

	// Validar tipo de archivo
	if !AllowedImageTypes[file.Header.Get("Content-Type")] {
		return "", ErrInvalidFileType
	}

	// Crear directorio si no existe
	if err := os.MkdirAll(destinationDir, 0755); err != nil {
		return "", fmt.Errorf("error creando directorio: %v", err)
	}

	// Generar nombre unico para el archivo
	ext := filepath.Ext(file.Filename)
	uniqueFilename := fmt.Sprintf("%s-%d%s", uuid.New().String(), time.Now().Unix(), ext)
	
	// Ruta completa del archivo
	filePath := filepath.Join(destinationDir, uniqueFilename)

	// Abrir archivo subido
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("error abriendo archivo: %v", err)
	}
	defer src.Close()

	// Crear archivo destino
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("error creando archivo: %v", err)
	}
	defer dst.Close()

	// Copiar contenido
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("error copiando archivo: %v", err)
	}

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
	return fmt.Sprintf("%s/static/home/%s", baseURL, filename)
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
	return filepath.Join("./static/home", filename)
}

// EnsureStaticDirs crea los directorios estaticos necesarios
func EnsureStaticDirs() error {
	dirs := []string{
		"./static",
		"./static/home",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("error creando directorio %s: %v", dir, err)
		}
	}

	return nil
}