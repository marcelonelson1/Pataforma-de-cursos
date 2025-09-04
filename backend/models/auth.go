package models

// RegisterRequest representa una solicitud de registro de usuario
type RegisterRequest struct {
	Nombre   string `json:"nombre" binding:"required,min=2,max=50"`
	Apellido string `json:"apellido" binding:"required,min=2,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"` // Opcional, por defecto será "user"
}

// LoginRequest representa una solicitud de inicio de sesión
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse representa la respuesta a una solicitud de autenticación exitosa
type AuthResponse struct {
	Token string  `json:"token"`
	User  Usuario `json:"user"`
}

// ChangeRoleRequest representa una solicitud para cambiar el rol de un usuario
type ChangeRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

// UpdateProfileRequest representa una solicitud para actualizar el perfil de usuario
type UpdateProfileRequest struct {
	Nombre   string `json:"nombre" binding:"required,min=2,max=50"`
	Apellido string `json:"apellido" binding:"required,min=2,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone"`
}

// ChangePasswordRequest representa una solicitud para cambiar la contraseña
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=8"`
}

// NotificationSettingsRequest representa una solicitud para actualizar las preferencias de notificación
type NotificationSettingsRequest struct {
	EmailNotifications bool `json:"emailNotifications"`
	NewMessages        bool `json:"newMessages"`
	NewStudents        bool `json:"newStudents"`
	SalesReports       bool `json:"salesReports"`
	SystemUpdates      bool `json:"systemUpdates"`
}
