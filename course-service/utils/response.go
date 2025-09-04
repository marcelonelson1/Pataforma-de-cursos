// utils/response.go
package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response estructura estándar para respuestas de la API
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    int         `json:"code,omitempty"`
}

// ErrorResponse estructura para respuestas de error
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Details string `json:"details,omitempty"`
}

// PaginatedResponse estructura para respuestas paginadas
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// Pagination información de paginación
type Pagination struct {
	Total       int64 `json:"total"`
	Page        int   `json:"page"`
	Limit       int   `json:"limit"`
	Pages       int64 `json:"pages"`
	HasNext     bool  `json:"has_next"`
	HasPrevious bool  `json:"has_previous"`
}

// SendSuccessResponse envía una respuesta exitosa
func SendSuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// SendCreatedResponse envía una respuesta de recurso creado
func SendCreatedResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Recurso creado exitosamente",
		Data:    data,
	})
}

// SendErrorResponse envía una respuesta de error
func SendErrorResponse(c *gin.Context, message string, statusCode int) {
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Error:   message,
		Code:    statusCode,
	})
}

// SendDetailedErrorResponse envía una respuesta de error con detalles
func SendDetailedErrorResponse(c *gin.Context, message string, details string, statusCode int) {
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Error:   message,
		Code:    statusCode,
		Details: details,
	})
}

// SendValidationErrorResponse envía errores de validación
func SendValidationErrorResponse(c *gin.Context, errors []string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"error":   "Errores de validación",
		"code":    http.StatusBadRequest,
		"details": errors,
	})
}

// SendPaginatedResponse envía una respuesta paginada
func SendPaginatedResponse(c *gin.Context, data interface{}, total int64, page, limit int) {
	pages := (total + int64(limit) - 1) / int64(limit)
	
	pagination := Pagination{
		Total:       total,
		Page:        page,
		Limit:       limit,
		Pages:       pages,
		HasNext:     int64(page) < pages,
		HasPrevious: page > 1,
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Success:    true,
		Data:       data,
		Pagination: pagination,
	})
}

// SendNotFoundResponse envía respuesta de recurso no encontrado
func SendNotFoundResponse(c *gin.Context, resource string) {
	SendErrorResponse(c, resource+" no encontrado", http.StatusNotFound)
}

// SendUnauthorizedResponse envía respuesta de no autorizado
func SendUnauthorizedResponse(c *gin.Context) {
	SendErrorResponse(c, "No autorizado", http.StatusUnauthorized)
}

// SendForbiddenResponse envía respuesta de acceso prohibido
func SendForbiddenResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Acceso prohibido"
	}
	SendErrorResponse(c, message, http.StatusForbidden)
}

// SendInternalServerErrorResponse envía error interno del servidor
func SendInternalServerErrorResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Error interno del servidor"
	}
	SendErrorResponse(c, message, http.StatusInternalServerError)
}

// SendBadRequestResponse envía respuesta de solicitud incorrecta
func SendBadRequestResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Solicitud incorrecta"
	}
	SendErrorResponse(c, message, http.StatusBadRequest)
}

// SendCustomResponse envía una respuesta personalizada
func SendCustomResponse(c *gin.Context, statusCode int, success bool, message string, data interface{}) {
	response := Response{
		Success: success,
		Data:    data,
	}
	
	if message != "" {
		if success {
			response.Message = message
		} else {
			response.Error = message
		}
	}
	
	c.JSON(statusCode, response)
}

// HealthCheckResponse estructura para health check
type HealthCheckResponse struct {
	Status    string                 `json:"status"`
	Service   string                 `json:"service"`
	Version   string                 `json:"version"`
	Timestamp string                 `json:"timestamp"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

// SendHealthCheckResponse envía respuesta de health check
func SendHealthCheckResponse(c *gin.Context, status string, details map[string]interface{}) {
	statusCode := http.StatusOK
	if status != "ok" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, HealthCheckResponse{
		Status:    status,
		Service:   "course-service",
		Version:   "1.0.0",
		Timestamp: GetCurrentTimestamp(),
		Details:   details,
	})
}