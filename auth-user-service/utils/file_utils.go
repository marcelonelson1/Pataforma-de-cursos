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

// SaveUploadedFile guarda un archivo subido
func SaveUploadedFile(file *multipart.FileHeader, destinationPath string) (string, error) {
	// Abrir el archivo
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Crear directorio si no existe
	if err := os.MkdirAll(filepath.Dir(destinationPath), 0755); err != nil {
		return "", err
	}

	// Crear archivo de destino
	dst, err := os.Create(destinationPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copiar contenido
	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return destinationPath, nil
}

// GenerateUniqueFileName genera un nombre único para el archivo
func GenerateUniqueFileName(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	name := strings.TrimSuffix(originalFilename, ext)
	timestamp := time.Now().Unix()
	uuid := uuid.New().String()[:8]
	
	return fmt.Sprintf("%s_%d_%s%s", name, timestamp, uuid, ext)
}

// IsAllowedImageType verifica si el tipo de archivo está permitido
func IsAllowedImageType(filename string) bool {
	allowedTypes := []string{".jpg", ".jpeg", ".png", ".webp"}
	ext := strings.ToLower(filepath.Ext(filename))
	
	for _, allowedType := range allowedTypes {
		if ext == allowedType {
			return true
		}
	}
	return false
}

// GetFileSize obtiene el tamaño del archivo
func GetFileSize(file *multipart.FileHeader) int64 {
	return file.Size
}

// CreateDirIfNotExists crea un directorio si no existe
func CreateDirIfNotExists(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return os.MkdirAll(dirPath, 0755)
	}
	return nil
}

// DeleteFile elimina un archivo si existe
func DeleteFile(filePath string) error {
	if _, err := os.Stat(filePath); err == nil {
		return os.Remove(filePath)
	}
	return nil
}