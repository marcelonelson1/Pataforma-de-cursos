package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Errores comunes del sistema
var (
	ErrInvalidToken     = errors.New("token inválido o expirado")
	ErrUserNotFound     = errors.New("usuario no encontrado")
	ErrEmailExists      = errors.New("este email ya está registrado")
	ErrInvalidLogin     = errors.New("credenciales inválidas")
	ErrPasswordMismatch = errors.New("las contraseñas no coinciden")
	ErrServerError      = errors.New("error interno del servidor")
	ErrDatabaseError    = errors.New("error de conexión a la base de datos")
	ErrEmailSendError   = errors.New("error al enviar correo electrónico")
)

// HandleError función para manejar errores de forma centralizada
func HandleError(c *gin.Context, err error, status int) {
	c.JSON(status, gin.H{"error": err.Error()})
}

// SendErrorResponse envía una respuesta de error estandarizada
func SendErrorResponse(c *gin.Context, err error, status int) {
	c.JSON(status, gin.H{
		"status": "error",
		"error":  err.Error(),
	})
}

// SendSuccessResponse envía una respuesta de éxito estandarizada
func SendSuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   data,
	})
}