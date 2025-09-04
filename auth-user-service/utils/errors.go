package utils

import "errors"

// Errores personalizados del sistema
var (
	ErrEmailExists      = errors.New("el email ya está registrado")
	ErrInvalidLogin     = errors.New("credenciales inválidas")
	ErrInvalidToken     = errors.New("token inválido o expirado")
	ErrUserNotFound     = errors.New("usuario no encontrado")
	ErrUnauthorized     = errors.New("no autorizado")
	ErrInvalidPassword  = errors.New("contraseña actual incorrecta")
	ErrDatabaseError    = errors.New("error de base de datos")
	ErrInvalidFileType  = errors.New("tipo de archivo no permitido")
	ErrFileTooLarge     = errors.New("archivo demasiado grande")
	ErrEmailSendFailed  = errors.New("error al enviar email")
	ErrTokenExpired     = errors.New("token expirado")
	ErrInvalidRequest   = errors.New("solicitud inválida")
)