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
	ErrUnauthorized     = errors.New("no autorizado")
	ErrServerError      = errors.New("error interno del servidor")
	ErrDatabaseError    = errors.New("error de base de datos")
	ErrEmailSendError   = errors.New("error al enviar correo")
	ErrInvalidJSON      = errors.New("JSON inválido")
	ErrInvalidRequest   = errors.New("solicitud inválida")
	ErrResourceNotFound = errors.New("recurso no encontrado")
	ErrBadRequest       = errors.New("solicitud mal formada")
	ErrValidationFailed = errors.New("validación fallida")
	ErrMissingFields    = errors.New("campos requeridos faltantes")
	ErrPaymentFailed    = errors.New("pago fallido")
	ErrPaymentExists    = errors.New("pago ya existente")
	ErrPaymentRejected  = errors.New("pago rechazado")
	ErrPaymentNotFound  = errors.New("pago no encontrado")
)

// SendErrorResponse envía una respuesta de error estandarizada
func SendErrorResponse(c *gin.Context, err error, status int) {
	c.JSON(status, gin.H{
		"success": false,
		"error":   err.Error(),
	})
	c.Abort()
}

// SendValidationErrorResponse envía una respuesta de error de validación con detalles
func SendValidationErrorResponse(c *gin.Context, err error, details interface{}) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"error":   err.Error(),
		"details": details,
	})
	c.Abort()
}

// SendSuccessResponse envía una respuesta exitosa estandarizada
func SendSuccessResponse(c *gin.Context, data interface{}) {
	response := gin.H{"success": true}
	if data != nil {
		response["data"] = data
	}
	c.JSON(http.StatusOK, response)
}

// RecoveryMiddleware maneja panics y asegura respuestas JSON consistentes
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				switch e := err.(type) {
				case error:
					SendErrorResponse(c, e, http.StatusInternalServerError)
				case string:
					SendErrorResponse(c, errors.New(e), http.StatusInternalServerError)
				default:
					SendErrorResponse(c, ErrServerError, http.StatusInternalServerError)
				}
			}
		}()
		c.Next()
	}
}

// JSONMiddleware asegura que todas las respuestas sean JSON
func JSONMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	}
}

// ErrorHandler maneja errores durante el procesamiento de solicitudes
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			// Si el cuerpo aún no se ha escrito
			if c.Writer.Size() == 0 {
				// Obtener el último error
				lastError := c.Errors.Last()
				var status int
				
				switch lastError.Err {
				case ErrUnauthorized:
					status = http.StatusUnauthorized
				case ErrResourceNotFound:
					status = http.StatusNotFound
				case ErrInvalidRequest, ErrValidationFailed, ErrMissingFields:
					status = http.StatusBadRequest
				default:
					status = http.StatusInternalServerError
				}

				SendErrorResponse(c, lastError, status)
			}
		}
	}
}