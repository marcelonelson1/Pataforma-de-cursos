package services

import (
	"auth-user-service/config"
	"auth-user-service/utils"
	"fmt"
	"mime/multipart"
	"path/filepath"
)

type FileService struct {
	userService *UserService
}

func NewFileService() *FileService {
	return &FileService{
		userService: NewUserService(),
	}
}

// UploadProfileImage sube una imagen de perfil
func (f *FileService) UploadProfileImage(userID uint, file *multipart.FileHeader) (string, error) {
	// Validar tipo de archivo
	if !utils.IsAllowedImageType(file.Filename) {
		return "", utils.ErrInvalidFileType
	}

	// Validar tamaño (5MB máximo)
	maxSize := int64(5 * 1024 * 1024) // 5MB
	if utils.GetFileSize(file) > maxSize {
		return "", utils.ErrFileTooLarge
	}

	// Generar nombre único
	uniqueFilename := utils.GenerateUniqueFileName(file.Filename)
	
	// Crear directorio si no existe
	profilesDir := filepath.Join(config.AppConfig.UploadPath, "profiles")
	if err := utils.CreateDirIfNotExists(profilesDir); err != nil {
		return "", err
	}

	// Ruta completa del archivo
	filePath := filepath.Join(profilesDir, uniqueFilename)

	// Guardar archivo
	if _, err := utils.SaveUploadedFile(file, filePath); err != nil {
		return "", err
	}

	// URL relativa para devolver
	imageURL := fmt.Sprintf("/static/profiles/%s", uniqueFilename)

	// Actualizar en la base de datos
	if err := f.userService.UpdateProfileImage(userID, imageURL); err != nil {
		// Si falla la BD, eliminar archivo
		utils.DeleteFile(filePath)
		return "", err
	}

	return imageURL, nil
}

// DeleteProfileImage elimina una imagen de perfil
func (f *FileService) DeleteProfileImage(imageURL string) error {
	if imageURL == "" {
		return nil
	}

	// Construir ruta completa
	filename := filepath.Base(imageURL)
	filePath := filepath.Join(config.AppConfig.UploadPath, "profiles", filename)

	return utils.DeleteFile(filePath)
}