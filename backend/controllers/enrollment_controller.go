package controllers

import (
	"curso-platform/middleware"
	"curso-platform/models"
	"curso-platform/services"
	"curso-platform/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// EnrollmentController gestiona las operaciones de inscripción a actividades
type EnrollmentController struct {
	enrollmentService *services.EnrollmentService
}

// NewEnrollmentController crea una nueva instancia del controlador de inscripciones
func NewEnrollmentController(enrollmentService *services.EnrollmentService) *EnrollmentController {
	return &EnrollmentController{
		enrollmentService: enrollmentService,
	}
}

// Enroll inscribe a un usuario en una actividad
func (c *EnrollmentController) Enroll(ctx *gin.Context) {
	// Obtener usuario del contexto
	userValue, exists := ctx.Get("user")
	if !exists {
		utils.SendErrorResponse(ctx, utils.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	user, ok := userValue.(models.Usuario)
	if !ok {
		utils.SendErrorResponse(ctx, utils.ErrServerError, http.StatusInternalServerError)
		return
	}

	// Obtener ID de la actividad
	id := ctx.Param("id")
	activityID, err := strconv.Atoi(id)
	if err != nil {
		utils.SendErrorResponse(ctx, utils.ErrBadRequest, http.StatusBadRequest)
		return
	}

	// Realizar la inscripción
	enrollment, err := c.enrollmentService.Enroll(user.ID, uint(activityID))
	if err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	// Registrar actividad en el log
	middleware.LogActivity(ctx, user.ID, "enrollment", 
		fmt.Sprintf("Usuario se inscribió en la actividad ID: %d", activityID))

	utils.SendSuccessResponse(ctx, gin.H{
		"message":    "Inscripción exitosa",
		"enrollment": enrollment,
	})
}

// GetUserActivities obtiene las actividades en las que está inscrito el usuario
func (c *EnrollmentController) GetUserActivities(ctx *gin.Context) {
	// Obtener usuario del contexto
	userValue, exists := ctx.Get("user")
	if !exists {
		utils.SendErrorResponse(ctx, utils.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	user, ok := userValue.(models.Usuario)
	if !ok {
		utils.SendErrorResponse(ctx, utils.ErrServerError, http.StatusInternalServerError)
		return
	}

	// Obtener actividades del usuario
	activities, err := c.enrollmentService.GetUserActivities(user.ID)
	if err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(ctx, activities)
}

// RegisterRoutes registra las rutas para el controlador de inscripciones
func (c *EnrollmentController) RegisterRoutes(router *gin.Engine) {
	enrollment := router.Group("/api/enrollments")
	enrollment.Use(middleware.AuthMiddleware())
	{
		enrollment.POST("/:id", c.Enroll)
		enrollment.GET("/my-activities", c.GetUserActivities)
	}
}