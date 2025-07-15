package controllers

import (
	"curso-platform/middleware"
	"curso-platform/models"
	"curso-platform/services"
	"curso-platform/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ProgresoController gestiona las operaciones relacionadas con el progreso del usuario
type ProgresoController struct {
	progresoService *services.ProgresoService
}

// NewProgresoController crea una nueva instancia del controlador de progreso
func NewProgresoController(progresoService *services.ProgresoService) *ProgresoController {
	return &ProgresoController{
		progresoService: progresoService,
	}
}

// GetProgresoUsuario obtiene el progreso de un usuario en un curso específico
func (c *ProgresoController) GetProgresoUsuario(ctx *gin.Context) {
	// Obtener el usuario autenticado
	userValue, exists := ctx.Get("user")
	if !exists {
		utils.SendErrorResponse(ctx, utils.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	usuario, ok := userValue.(models.Usuario)
	if !ok {
		utils.SendErrorResponse(ctx, utils.ErrServerError, http.StatusInternalServerError)
		return
	}

	// Obtener ID del curso
	cursoID, err := strconv.ParseUint(ctx.Param("cursoId"), 10, 64)
	if err != nil {
		utils.SendErrorResponse(ctx, utils.ErrBadRequest, http.StatusBadRequest)
		return
	}

	// Obtener el progreso
	progreso, err := c.progresoService.GetProgresoUsuario(usuario.ID, uint(cursoID))
	if err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, progreso)
}

// MarcarCapituloCompletado marca un capítulo como completado/incompleto
func (c *ProgresoController) MarcarCapituloCompletado(ctx *gin.Context) {
	// Obtener el usuario autenticado
	userValue, exists := ctx.Get("user")
	if !exists {
		utils.SendErrorResponse(ctx, utils.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	usuario, ok := userValue.(models.Usuario)
	if !ok {
		utils.SendErrorResponse(ctx, utils.ErrServerError, http.StatusInternalServerError)
		return
	}

	// Obtener datos de la solicitud
	var req struct {
		CursoID    uint    `json:"curso_id" binding:"required"`
		CapituloID uint    `json:"capitulo_id" binding:"required"`
		Completado bool    `json:"completado"`
		Progreso   float64 `json:"progreso"`
	}
	
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	// Marcar capítulo
	progresoCapitulo, err := c.progresoService.MarcarCapituloCompletado(usuario.ID, req)
	if err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "Progreso actualizado correctamente",
		"progreso": progresoCapitulo,
	})
}

// GuardarUltimoCapitulo guarda el último capítulo visto
func (c *ProgresoController) GuardarUltimoCapitulo(ctx *gin.Context) {
	// Obtener el usuario autenticado
	userValue, exists := ctx.Get("user")
	if !exists {
		utils.SendErrorResponse(ctx, utils.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	usuario, ok := userValue.(models.Usuario)
	if !ok {
		utils.SendErrorResponse(ctx, utils.ErrServerError, http.StatusInternalServerError)
		return
	}

	// Obtener datos de la solicitud
	var req struct {
		CursoID    uint `json:"curso_id" binding:"required"`
		CapituloID uint `json:"capitulo_id" binding:"required"`
	}
	
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := c.progresoService.GuardarUltimoCapitulo(usuario.ID, req); err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Último capítulo actualizado correctamente",
	})
}

// RegisterRoutes registra todas las rutas relacionadas con el progreso
func (c *ProgresoController) RegisterRoutes(router *gin.Engine) {
	progreso := router.Group("/api/progreso")
	{
		progreso.Use(middleware.AuthMiddleware())
		progreso.GET("/curso/:cursoId", c.GetProgresoUsuario)
		progreso.POST("/capitulo/completado", c.MarcarCapituloCompletado)
		progreso.POST("/ultimo-capitulo", c.GuardarUltimoCapitulo)
	}
}