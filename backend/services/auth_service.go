package services

import (
	"curso-platform/config"
	"curso-platform/models"
	"curso-platform/utils"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// AuthService maneja la lógica relacionada con la autenticación
type AuthService struct{}

// NewAuthService crea una nueva instancia de AuthService
func NewAuthService() *AuthService {
	return &AuthService{}
}

// Register registra un nuevo usuario
func (s *AuthService) Register(req models.RegisterRequest) error {
	// Verificar si el email ya existe
	var existingUser models.Usuario
	if result := config.DB.Where("email = ?", req.Email).First(&existingUser); result.Error == nil {
		return utils.ErrEmailExists
	}

	// Hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("error al procesar la contraseña")
	}

	// Establecer rol por defecto si no se proporciona
	role := req.Role
	if role == "" {
		role = "user" // Por defecto, todos los usuarios son "user"
	}

	// Crear nuevo usuario - CORREGIDO: Agregado campo Apellido
	user := models.Usuario{
		Nombre:   req.Nombre,
		Apellido: req.Apellido, // ✅ AGREGADO: Campo apellido
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     role,
	}

	if result := config.DB.Create(&user); result.Error != nil {
		return errors.New("error al crear usuario")
	}

	return nil
}

// Login inicia sesión de un usuario
func (s *AuthService) Login(req models.LoginRequest) (*models.AuthResponse, error) {
	// Buscar usuario por email
	var user models.Usuario
	if result := config.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		return nil, utils.ErrInvalidLogin
	}

	// Verificar contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, utils.ErrInvalidLogin
	}

	// Actualizar última conexión
	now := time.Now()
	if err := config.DB.Model(&user).Update("last_login", now).Error; err != nil {
		// Solo registramos el error, no detenemos el proceso de login
		log.Printf("Error al actualizar última conexión: %v", err)
	}

	// Generar token JWT
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, errors.New("error al generar token")
	}

	// No enviar la contraseña en la respuesta
	user.Password = ""

	// Actualizar last_login en la instancia local para devolverla en la respuesta
	user.LastLogin = now

	// Preparar respuesta
	response := &models.AuthResponse{
		Token: token,
		User:  user,
	}

	return response, nil
}

// RefreshToken renueva el token JWT de un usuario
func (s *AuthService) RefreshToken(userID uint, role string) (string, error) {
	return utils.GenerateToken(userID, role)
}

// ChangeUserRole cambia el rol de un usuario
func (s *AuthService) ChangeUserRole(userID uint, newRole string) (*models.Usuario, error) {
	// Verificar que el usuario existe
	var user models.Usuario
	if result := config.DB.First(&user, userID); result.Error != nil {
		return nil, utils.ErrResourceNotFound
	}

	// Validar que el rol sea válido
	if newRole != "admin" && newRole != "user" {
		return nil, errors.New("rol no válido")
	}

	// Actualizar rol del usuario
	if err := config.DB.Model(&user).Update("role", newRole).Error; err != nil {
		return nil, utils.ErrDatabaseError
	}

	// Obtener usuario actualizado
	if err := config.DB.First(&user, userID).Error; err != nil {
		return nil, utils.ErrDatabaseError
	}

	// No enviar contraseña
	user.Password = ""

	return &user, nil
}

// ChangePassword cambia la contraseña de un usuario
func (s *AuthService) ChangePassword(userID uint, currentPassword, newPassword string) error {
	// Buscar usuario
	var user models.Usuario
	if result := config.DB.First(&user, userID); result.Error != nil {
		return utils.ErrUserNotFound
	}

	// Verificar contraseña actual
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return errors.New("contraseña actual incorrecta")
	}

	// Hash de la nueva contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("error al procesar la nueva contraseña")
	}

	// Actualizar contraseña en la BD
	if result := config.DB.Model(&user).Update("password", string(hashedPassword)); result.Error != nil {
		return utils.ErrDatabaseError
	}

	return nil
}
