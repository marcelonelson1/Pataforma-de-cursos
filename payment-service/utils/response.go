// utils/response.go
package utils

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// SendSuccessResponse envía una respuesta exitosa estandarizada
func SendSuccessResponse(c *gin.Context, data interface{}) {
	response := gin.H{"success": true}
	if data != nil {
		response["data"] = data
	}
	c.JSON(http.StatusOK, response)
}

// SendErrorResponse envía una respuesta de error estandarizada
func SendErrorResponse(c *gin.Context, message string, status int) {
	c.JSON(status, gin.H{
		"success": false,
		"error":   message,
	})
	c.Abort()
}

// SendValidationErrorResponse envía una respuesta de error de validación con detalles
func SendValidationErrorResponse(c *gin.Context, message string, details interface{}) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"error":   message,
		"details": details,
	})
	c.Abort()
}

// SendCreatedResponse envía una respuesta de creación exitosa
func SendCreatedResponse(c *gin.Context, data interface{}) {
	response := gin.H{"success": true}
	if data != nil {
		response["data"] = data
	}
	c.JSON(http.StatusCreated, response)
}

// SendNoContentResponse envía una respuesta sin contenido
func SendNoContentResponse(c *gin.Context) {
	c.JSON(http.StatusNoContent, gin.H{
		"success": true,
	})
}

// utils/errors.go
