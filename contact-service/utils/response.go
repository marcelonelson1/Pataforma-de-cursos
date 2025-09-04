package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// SendSuccessResponse envia una respuesta exitosa
func SendSuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

// SendSuccessMessage envia un mensaje exitoso con datos opcionales
func SendSuccessMessage(c *gin.Context, message string, data interface{}) {
	response := gin.H{
		"success": true,
		"message": message,
	}
	
	if data != nil {
		response["data"] = data
	}
	
	c.JSON(http.StatusOK, response)
}

// SendErrorResponse envia una respuesta de error
func SendErrorResponse(c *gin.Context, err error, statusCode int) {
	c.JSON(statusCode, gin.H{
		"success": false,
		"error":   err.Error(),
	})
}

// SendErrorMessage envia un mensaje de error personalizado
func SendErrorMessage(c *gin.Context, message string, statusCode int) {
	c.JSON(statusCode, gin.H{
		"success": false,
		"error":   message,
	})
}

// SendValidationError envia errores de validacion
func SendValidationError(c *gin.Context, err error) {
	var validationErrors []string
	
	if validationErr, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErr {
			validationErrors = append(validationErrors, fieldError.Error())
		}
	} else {
		validationErrors = append(validationErrors, err.Error())
	}
	
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"error":   "Error de validacion",
		"details": validationErrors,
	})
}