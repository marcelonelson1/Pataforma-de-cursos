package middleware

import (
	"curso-platform/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware maneja panics y asegura respuestas JSON consistentes
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				switch e := err.(type) {
				case error:
					utils.SendErrorResponse(c, e, http.StatusInternalServerError)
				case string:
					utils.SendErrorResponse(c, errors.New(e), http.StatusInternalServerError)
				default:
					utils.SendErrorResponse(c, utils.ErrServerError, http.StatusInternalServerError)
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
			if c.Writer.Size() == 0 {
				lastError := c.Errors.Last()
				var status int
				
				switch lastError.Err {
				case utils.ErrUnauthorized:
					status = http.StatusUnauthorized
				case utils.ErrResourceNotFound:
					status = http.StatusNotFound
				case utils.ErrInvalidRequest, utils.ErrValidationFailed, utils.ErrMissingFields:
					status = http.StatusBadRequest
				default:
					status = http.StatusInternalServerError
				}

				utils.SendErrorResponse(c, lastError.Err, status)
			}
		}
	}
}