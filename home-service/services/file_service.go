package services

import (
	"home-service/utils"
	"mime/multipart"
	"fmt"
	"log"
	"path/filepath"
)

type FileService struct {
	homeService *HomeService
}

// NewFileService crea una nueva instancia de FileService
func NewFileService() *FileService {
	return &FileService{
		homeService: NewHomeService(),
	}
}

// UploadHomeImage sube una imagen para el home
func (fs *FileService) UploadHomeImage(imageID uint, file *multipart.FileHeader) (string, error) {
	// Verificar que la imagen existe
	_, err := fs.homeService.GetImageByID(imageID)
	if err != nil {
		return "", err
	}

	// Validar archivo
	if err := utils.ValidateImageFile(file); err != nil {
		return "", err
	}

	// Definir directorio de destino
	destinationDir := "./static/home"

	// Guardar archivo
	filename, err := utils.SaveUploadedFile(file, destinationDir)
	if err != nil {
		log.Printf("Error guardando archivo: %v", err)
		return "", utils.ErrFileUploadFailed
	}

	// Actualizar imagen con nuevo archivo
	_, err = fs.homeService.UpdateImageFile(imageID, filename)
	if err != nil {
		// Si falla la actualizacion, eliminar archivo subido
		filePath := filepath.Join(destinationDir, filename)
		utils.DeleteFile(filePath)
		return "", err
	}

	log.Printf("Imagen subida exitosamente para home ID %d: %s", imageID, filename)
	
	// Retornar URL relativa de la imagen
	return fmt.Sprintf("/static/home/%s", filename), nil
}

// DeleteHomeImage elimina la imagen de un elemento del home
func (fs *FileService) DeleteHomeImage(imageID uint) error {
	// Obtener imagen
	image, err := fs.homeService.GetImageByID(imageID)
	if err != nil {
		return err
	}

	// Verificar que tiene imagen
	if image.ImageURL == "" {
		return fmt.Errorf("la imagen del home no tiene archivo")
	}

	// Eliminar archivo del disco
	imagePath := utils.GetImagePath(image.ImageURL)
	if err := utils.DeleteFile(imagePath); err != nil {
		log.Printf("Error eliminando archivo de imagen: %v", err)
	}

	// Actualizar imagen removiendo la URL del archivo
	_, err = fs.homeService.UpdateImageFile(imageID, "")
	if err != nil {
		return err
	}

	log.Printf("Imagen eliminada exitosamente para home ID %d", imageID)
	return nil
}

// GetImageURL construye la URL completa de una imagen
func (fs *FileService) GetImageURL(filename string, baseURL string) string {
	return utils.GetFileURL(filename, baseURL)
}

// ValidateAndPrepareUpload valida un archivo antes de subirlo
func (fs *FileService) ValidateAndPrepareUpload(file *multipart.FileHeader, imageID uint) error {
	// Verificar que la imagen existe
	_, err := fs.homeService.GetImageByID(imageID)
	if err != nil {
		return err
	}

	// Validar archivo
	if err := utils.ValidateImageFile(file); err != nil {
		return err
	}

	// Asegurar que existen los directorios necesarios
	if err := utils.EnsureStaticDirs(); err != nil {
		return fmt.Errorf("error preparando directorios: %v", err)
	}

	return nil
}

// GetAllowedFileTypes retorna los tipos de archivo permitidos
func (fs *FileService) GetAllowedFileTypes() map[string]bool {
	return utils.AllowedImageTypes
}

// GetMaxFileSize retorna el tamano maximo permitido
func (fs *FileService) GetMaxFileSize() int64 {
	return utils.MaxFileSize
}