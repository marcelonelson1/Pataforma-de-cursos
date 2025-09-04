package services

import (
	"auth-user-service/models"
	"auth-user-service/utils"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

type AuthService struct {
	db             *gorm.DB
	jwtService     *JWTService
	passwordService *PasswordService
	userService    *UserService
}

func NewAuthService() *AuthService {
	return &AuthService{
		db:             utils.GetDB(),
		jwtService:     NewJWTService(),
		passwordService: NewPasswordService(),
		userService:    NewUserService(),
	}
}

// Register registra un nuevo usuario
func (a *AuthService) Register(req *models.RegisterRequest) error {
	// Verificar si el email ya existe
	var existingUser models.Usuario
	if result := a.db.Where("email = ?", req.Email).First(&existingUser); result.Error == nil {
		return utils.ErrEmailExists
	}

	// Hash de la contraseña
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	// Establecer rol por defecto si no se proporciona
	role := req.Role
	if role == "" {
		role = "user"
	}

	// Crear nuevo usuario
	user := models.Usuario{
		Nombre:   req.Nombre,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     role,
	}

	// Iniciar transacción
	tx := a.db.Begin()

	// Crear usuario
	if result := tx.Create(&user); result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	// Crear configuraciones de notificación por defecto
	notificationSettings := models.NotificationSettings{
		UserID:            user.ID,
		EmailNotifications: true,
		PushNotifications:  true,
		CourseUpdates:      true,
		PaymentAlerts:      true,
	}

	if result := tx.Create(&notificationSettings); result.Error != nil {
		tx.Rollback()
		log.Printf("Error al crear configuraciones de notificación: %v", result.Error)
		return result.Error
	}

	// Confirmar transacción
	tx.Commit()

	// Log de actividad
	a.userService.LogActivity(user.ID, "register", "Usuario registrado", "", "")

	return nil
}

// Login autentica un usuario
func (a *AuthService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
	// Buscar usuario por email
	var user models.Usuario
	if result := a.db.Where("email = ?", req.Email).First(&user); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, utils.ErrInvalidLogin
		}
		return nil, result.Error
	}

	// Verificar contraseña
	if err := utils.CheckPassword(user.Password, req.Password); err != nil {
		return nil, utils.ErrInvalidLogin
	}

	// Actualizar última conexión
	now := time.Now()
	if err := a.db.Model(&user).Update("last_login", now).Error; err != nil {
		log.Printf("Error al actualizar última conexión: %v", err)
	}

	// Generar token JWT
	token, err := a.jwtService.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	// No enviar la contraseña en la respuesta
	user.Password = ""
	user.LastLogin = now

	// Log de actividad
	a.userService.LogActivity(user.ID, "login", "Usuario inició sesión", "", "")

	return &models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

// ValidateToken valida un token JWT
func (a *AuthService) ValidateToken(tokenString string) (*models.TokenValidationResponse, error) {
	claims, err := a.jwtService.ValidateToken(tokenString)
	if err != nil {
		return &models.TokenValidationResponse{Valid: false}, utils.ErrInvalidToken
	}

	// Verificar si el usuario existe en la base de datos
	var user models.Usuario
	if result := a.db.First(&user, claims.UserID); result.Error != nil {
		return &models.TokenValidationResponse{Valid: false}, utils.ErrUserNotFound
	}

	// Verificar que el rol en el token coincida con el de la base de datos
	if claims.Role != user.Role {
		log.Printf("Diferencia en roles: Token (%s) vs. DB (%s) para usuario ID: %d", claims.Role, user.Role, user.ID)
	}

	// No enviar la contraseña
	user.Password = ""

	return &models.TokenValidationResponse{
		Valid:  true,
		UserID: user.ID,
		Role:   user.Role,
		User:   &user,
	}, nil
}

// RefreshToken genera un nuevo token
func (a *AuthService) RefreshToken(userID uint) (string, error) {
	// Buscar el usuario para obtener el rol actual
	var user models.Usuario
	if result := a.db.First(&user, userID); result.Error != nil {
		return "", utils.ErrUserNotFound
	}

	return a.jwtService.RefreshToken(user.ID, user.Role)
}

// CheckAdmin verifica si un usuario es administrador
func (a *AuthService) CheckAdmin(userID uint) (bool, error) {
	var user models.Usuario
	if result := a.db.Select("role").First(&user, userID); result.Error != nil {
		return false, utils.ErrUserNotFound
	}

	return user.Role == "admin", nil
}

// ForgotPassword inicia el proceso de recuperación de contraseña
func (a *AuthService) ForgotPassword(email string) (*models.PasswordReset, error) {
	// Verificar si el email existe
	var user models.Usuario
	result := a.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		// Por seguridad, no revelamos si el email existe
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, utils.ErrUserNotFound
		}
		return nil, result.Error
	}

	// Generar token de restablecimiento
	return a.passwordService.GenerateResetToken(email)
}

// ChangePassword cambia la contraseña del usuario
func (a *AuthService) ChangePassword(userID uint, currentPassword, newPassword string) error {
	err := a.passwordService.ChangePassword(userID, currentPassword, newPassword)
	if err == nil {
		// Log de actividad
		a.userService.LogActivity(userID, "change_password", "Contraseña cambiada", "", "")
	}
	return err
}

// ResetPassword restablece la contraseña usando un token
func (a *AuthService) ResetPassword(token, newPassword string) error {
	err := a.passwordService.ResetPassword(token, newPassword)
	if err == nil {
		// Buscar el usuario por el token para hacer log
		reset, _ := a.passwordService.ValidateResetToken(token)
		if reset != nil {
			var user models.Usuario
			if result := a.db.Where("email = ?", reset.Email).First(&user); result.Error == nil {
				a.userService.LogActivity(user.ID, "reset_password", "Contraseña restablecida", "", "")
			}
		}
	}
	return err
}

// ValidateResetToken valida un token de reset
func (a *AuthService) ValidateResetToken(token string) (*models.PasswordReset, error) {
	return a.passwordService.ValidateResetToken(token)
}