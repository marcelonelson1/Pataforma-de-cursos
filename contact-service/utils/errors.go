package utils

import (
	"errors"
	"os"
)

// Errores personalizados del sistema
var (
	ErrDatabaseError     = errors.New("error de base de datos")
	ErrResourceNotFound  = errors.New("recurso no encontrado")
	ErrInvalidRequest    = errors.New("solicitud invalida")
	ErrEmailSendError    = errors.New("error al enviar email")
	ErrUnauthorized      = errors.New("no autorizado")
	ErrForbidden         = errors.New("acceso denegado")
	ErrValidationFailed  = errors.New("error de validacion")
)

// GetEnv obtiene una variable de entorno con valor por defecto
func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}