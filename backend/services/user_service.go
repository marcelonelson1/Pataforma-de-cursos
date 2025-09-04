package services

import (
	"curso-platform/config"
	"curso-platform/models"
	"curso-platform/utils"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/gorm"
)

// UserService maneja la lógica relacionada con los usuarios
type UserService struct{}

// NewUserService crea una nueva instancia de UserService
func NewUserService() *UserService {
	return &UserService{}
}

// GetUserByID obtiene un usuario por su ID
func (s *UserService) GetUserByID(id uint) (*models.Usuario, error) {
	var user models.Usuario
	if result := config.DB.First(&user, id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, utils.ErrResourceNotFound
		}
		return nil, utils.ErrDatabaseError
	}

	// No enviar la contraseña
	user.Password = ""

	return &user, nil
}

// UpdateProfile actualiza el perfil de un usuario - CORREGIDO: Agregado campo Apellido
func (s *UserService) UpdateProfile(userID uint, req models.UpdateProfileRequest) (*models.Usuario, error) {
	// Obtener usuario actual
	var currentUser models.Usuario
	if result := config.DB.First(&currentUser, userID); result.Error != nil {
		return nil, utils.ErrResourceNotFound
	}

	// Verificar si el nuevo email ya existe (si se cambia)
	if req.Email != "" && req.Email != currentUser.Email {
		var existingUser models.Usuario
		if result := config.DB.Where("email = ? AND id != ?", req.Email, userID).First(&existingUser); result.Error == nil {
			return nil, utils.ErrEmailExists
		}
	}

	// Preparar updates - CORREGIDO: Agregado campo Apellido
	updates := map[string]interface{}{}

	if req.Nombre != "" {
		updates["nombre"] = req.Nombre
	}

	// ✅ AGREGADO: Campo apellido
	if req.Apellido != "" {
		updates["apellido"] = req.Apellido
	}

	if req.Email != "" {
		updates["email"] = req.Email
	}

	if req.Phone != "" {
		updates["phone"] = req.Phone
	}

	// Actualizar usuario en la BD
	if result := config.DB.Model(&currentUser).Updates(updates); result.Error != nil {
		return nil, utils.ErrDatabaseError
	}

	// Obtener usuario actualizado
	var updatedUser models.Usuario
	if err := config.DB.First(&updatedUser, userID).Error; err != nil {
		return nil, utils.ErrDatabaseError
	}

	// No enviar contraseña en la respuesta
	updatedUser.Password = ""

	return &updatedUser, nil
}

// UploadProfileImage sube una nueva imagen de perfil para el usuario
func (s *UserService) UploadProfileImage(userID uint, file *os.File, filename string, fileSize int64) (string, error) {
	// Validar tamaño (máx 2MB)
	if fileSize > 2*1024*1024 {
		return "", errors.New("la imagen no debe superar los 2MB")
	}

	// Validar tipo de archivo
	extension := filepath.Ext(filename)
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	if !allowedExtensions[extension] {
		return "", errors.New("tipo de archivo no permitido")
	}

	// Crear directorio para imágenes de usuario si no existe
	uploadDir := "./static/profiles"
	if err := utils.CreateDirIfNotExists(uploadDir); err != nil {
		return "", err
	}

	// Nombre de archivo: user_ID.extension
	newFilename := fmt.Sprintf("user_%d%s", userID, extension)
	destinationPath := filepath.Join(uploadDir, newFilename)

	// Si ya existe un archivo anterior con el mismo nombre, eliminarlo
	_, err := os.Stat(destinationPath)
	if err == nil {
		if err := os.Remove(destinationPath); err != nil {
			return "", err
		}
	}

	// Crear nuevo archivo
	destination, err := os.Create(destinationPath)
	if err != nil {
		return "", err
	}
	defer destination.Close()

	// Leer el archivo original
	fileData, err := os.ReadFile(file.Name())
	if err != nil {
		return "", err
	}

	// Escribir en el nuevo archivo
	if _, err := destination.Write(fileData); err != nil {
		return "", err
	}

	// URL relativa para acceder a la imagen
	imageURL := "/static/profiles/" + newFilename

	// Actualizar URL de imagen en la BD
	if err := config.DB.Model(&models.Usuario{}).Where("id = ?", userID).Update("image_url", imageURL).Error; err != nil {
		return "", utils.ErrDatabaseError
	}

	return imageURL, nil
}

// ListUsers obtiene una lista paginada de usuarios - CORREGIDO: Búsqueda incluye apellido
func (s *UserService) ListUsers(page, limit int, role, search string) ([]models.Usuario, int64, error) {
	var users []models.Usuario
	query := config.DB.Order("created_at desc")

	// Calcular offset para paginación
	offset := (page - 1) * limit

	// Filtrado por rol
	if role != "" {
		query = query.Where("role = ?", role)
	}

	// Búsqueda por nombre, apellido o email - CORREGIDO: Agregada búsqueda por apellido
	if search != "" {
		query = query.Where("nombre LIKE ? OR apellido LIKE ? OR email LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Contar total de registros para paginación
	var total int64
	if err := query.Model(&models.Usuario{}).Count(&total).Error; err != nil {
		return nil, 0, utils.ErrDatabaseError
	}

	// Obtener usuarios con paginación
	if err := query.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, utils.ErrDatabaseError
	}

	// No enviar contraseñas
	for i := range users {
		users[i].Password = ""
	}

	return users, total, nil
}

// DeleteUser elimina un usuario
func (s *UserService) DeleteUser(id uint, adminID uint) error {
	// Verificar si el usuario existe
	var user models.Usuario
	if err := config.DB.First(&user, id).Error; err != nil {
		return utils.ErrResourceNotFound
	}

	// Evitar que un admin elimine su propia cuenta
	if adminID == user.ID {
		return errors.New("no puedes eliminar tu propia cuenta")
	}

	// Eliminar usuario (borrado en cascada configurado en la base de datos)
	if err := config.DB.Delete(&user).Error; err != nil {
		return utils.ErrDatabaseError
	}

	return nil
}
