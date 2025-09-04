package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// StandardResponse estructura estándar de respuesta
type StandardResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SendSuccessResponse envía una respuesta exitosa
func SendSuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, StandardResponse{
		Success: true,
		Data:    data,
	})
}

// SendSuccessMessage envía una respuesta exitosa con mensaje
func SendSuccessMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SendErrorResponse envía una respuesta de error
func SendErrorResponse(c *gin.Context, err error, statusCode int) {
	c.JSON(statusCode, StandardResponse{
		Success: false,
		Error:   err.Error(),
	})
}

// SendErrorMessage envía una respuesta de error con mensaje personalizado
func SendErrorMessage(c *gin.Context, message string, statusCode int) {
	c.JSON(statusCode, StandardResponse{
		Success: false,
		Error:   message,
	})
}

// SendValidationError envía errores de validación
func SendValidationError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, StandardResponse{
		Success: false,
		Error:   "Error de validación: " + err.Error(),
	})
}