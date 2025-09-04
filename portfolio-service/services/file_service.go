package services

import (
	"portfolio-service/utils"
	"mime/multipart"
	"fmt"
	"log"
	"path/filepath"
)

type FileService struct {
	portfolioService *PortfolioService
}

// NewFileService crea una nueva instancia de FileService
func NewFileService() *FileService {
	return &FileService{
		portfolioService: NewPortfolioService(),
	}
}

// UploadProjectImage sube una imagen para un proyecto
func (fs *FileService) UploadProjectImage(projectID uint, file *multipart.FileHeader) (string, error) {
	log.Printf("DEBUG FileService: UploadProjectImage called for project %d", projectID)
	log.Printf("DEBUG FileService: File info - Name: %s, Size: %d", file.Filename, file.Size)
	
	// Verificar que el proyecto existe
	_, err := fs.portfolioService.GetProjectByID(projectID)
	if err != nil {
		log.Printf("DEBUG FileService: Project %d not found: %v", projectID, err)
		return "", err
	}
	log.Printf("DEBUG FileService: Project %d exists", projectID)

	// Validar archivo
	log.Printf("DEBUG FileService: Validating image file...")
	if err := utils.ValidateImageFile(file); err != nil {
		log.Printf("DEBUG FileService: File validation failed: %v", err)
		return "", err
	}
	log.Printf("DEBUG FileService: File validation passed")

	// Definir directorio de destino
	destinationDir := "./static/portfolio"
	log.Printf("DEBUG FileService: Destination directory: %s", destinationDir)

	// Guardar archivo
	log.Printf("DEBUG FileService: Attempting to save file...")
	filename, err := utils.SaveUploadedFile(file, destinationDir)
	if err != nil {
		log.Printf("Error guardando archivo: %v", err)
		return "", utils.ErrFileUploadFailed
	}
	log.Printf("DEBUG FileService: File saved successfully as: %s", filename)

	// Actualizar proyecto con nueva imagen
	_, err = fs.portfolioService.UpdateProjectImage(projectID, filename)
	if err != nil {
		// Si falla la actualizacion, eliminar archivo subido
		filePath := filepath.Join(destinationDir, filename)
		utils.DeleteFile(filePath)
		return "", err
	}

	log.Printf("Imagen subida exitosamente para proyecto %d: %s", projectID, filename)
	
	// Retornar URL relativa de la imagen
	return fmt.Sprintf("/static/portfolio/%s", filename), nil
}

// DeleteProjectImage elimina la imagen de un proyecto
func (fs *FileService) DeleteProjectImage(projectID uint) error {
	// Obtener proyecto
	project, err := fs.portfolioService.GetProjectByID(projectID)
	if err != nil {
		return err
	}

	// Verificar que tiene imagen
	if project.ImageURL == "" {
		return fmt.Errorf("el proyecto no tiene imagen")
	}

	// Eliminar archivo del disco
	imagePath := utils.GetImagePath(project.ImageURL)
	if err := utils.DeleteFile(imagePath); err != nil {
		log.Printf("Error eliminando archivo de imagen: %v", err)
	}

	// Actualizar proyecto removiendo la URL de la imagen
	_, err = fs.portfolioService.UpdateProjectImage(projectID, "")
	if err != nil {
		return err
	}

	log.Printf("Imagen eliminada exitosamente para proyecto %d", projectID)
	return nil
}

// GetImageURL construye la URL completa de una imagen
func (fs *FileService) GetImageURL(filename string, baseURL string) string {
	return utils.GetFileURL(filename, baseURL)
}

// ValidateAndPrepareUpload valida un archivo antes de subirlo
func (fs *FileService) ValidateAndPrepareUpload(file *multipart.FileHeader, projectID uint) error {
	// Verificar que el proyecto existe
	_, err := fs.portfolioService.GetProjectByID(projectID)
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