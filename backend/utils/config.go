package utils

import (
	"log"
	"os"
	"path/filepath"
)

// GetEnv recupera una variable de entorno o devuelve un valor predeterminado
func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// CreateDirIfNotExists crea un directorio si no existe
func CreateDirIfNotExists(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		log.Printf("Creando directorio: %s", dirPath)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			log.Printf("Error al crear directorio %s: %v", dirPath, err)
			return err
		}
	}
	return nil
}

// InitStaticDirs inicializa todos los directorios estáticos necesarios
func InitStaticDirs() {
	dirs := []string{
		"./static",
		"./static/videos",
		"./static/images",
		"./static/profiles",
		"./static/portfolio",
		"./static/home",
	}

	for _, dir := range dirs {
		CreateDirIfNotExists(dir)
	}
}

// EnsureDir asegura que un directorio específico exista
func EnsureDir(path string) error {
	dirPath := filepath.Dir(path)
	return CreateDirIfNotExists(dirPath)
}