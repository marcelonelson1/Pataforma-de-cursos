package models

import (
	"github.com/golang-jwt/jwt/v4"
)

// RegisterRequest estructura para registro de usuarios
type RegisterRequest struct {
	Nombre   string `json:"nombre" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"` // Opcional, por defecto será "user"
}

// LoginRequest estructura para login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse respuesta de autenticación
type AuthResponse struct {
	Token string  `json:"token"`
	User  Usuario `json:"user"`
}

// ChangePasswordRequest estructura para cambio de contraseña
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

// ForgotPasswordRequest estructura para solicitud de recuperación
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest estructura para restablecer contraseña
type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// UpdateProfileRequest estructura para actualizar perfil
type UpdateProfileRequest struct {
	Nombre string `json:"nombre"`
	Phone  string `json:"phone"`
}

// Claims JWT personalizado
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// TokenValidationResponse respuesta de validación de token
type TokenValidationResponse struct {
	Valid  bool     `json:"valid"`
	UserID uint     `json:"user_id,omitempty"`
	Role   string   `json:"role,omitempty"`
	User   *Usuario `json:"user,omitempty"`
}