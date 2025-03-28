package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidToken     = errors.New("token inválido o expirado")
	ErrUserNotFound     = errors.New("usuario no encontrado")
	ErrEmailExists      = errors.New("email ya registrado")
	ErrInvalidLogin     = errors.New("credenciales inválidas")
	ErrPasswordMismatch = errors.New("las contraseñas no coinciden")
	ErrServerError      = errors.New("error interno del servidor")
	ErrDatabaseError    = errors.New("error de base de datos")
	ErrEmailSendError   = errors.New("error al enviar correo")
)

func SendErrorResponse(c *gin.Context, err error, status int) {
	c.JSON(status, gin.H{
		"success": false,
		"error":   err.Error(),
	})
	c.Abort()
}

func SendSuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}