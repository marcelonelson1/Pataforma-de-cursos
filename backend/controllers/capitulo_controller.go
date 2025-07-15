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

// CapituloController gestiona las operaciones relacionadas con los capítulos
type CapituloController struct {
	capituloService *services.CapituloService
}

// NewCapituloController crea una nueva instancia del controlador de capítulos
func NewCapituloController(capituloService *services.CapituloService) *CapituloController {
	return &CapituloController{
		capituloService: capituloService,
	}
}

// GetCapitulosByCurso obtiene todos los capítulos de un curso
func (c *CapituloController) GetCapitulosByCurso(ctx *gin.Context) {
	cursoId := ctx.Param("cursoId")
	cursoIDInt, err := strconv.Atoi(cursoId)
	if err != nil {
		utils.SendErrorResponse(ctx, utils.ErrBadRequest, http.StatusBadRequest)
		return
	}

	capitulos, err := c.capituloService.GetByCursoID(uint(cursoIDInt))
	if err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, capitulos)
}

// CreateCapitulo crea un nuevo capítulo
func (c *CapituloController) CreateCapitulo(ctx *gin.Context) {
	var req models.Capitulo
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	capitulo, err := c.capituloService.Create(req)
	if err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, capitulo)
}

// UpdateCapitulo actualiza un capítulo existente
func (c *CapituloController) UpdateCapitulo(ctx *gin.Context) {
	id := ctx.Param("id")
	capituloID, err := strconv.Atoi(id)
	if err != nil {
		utils.SendErrorResponse(ctx, utils.ErrBadRequest, http.StatusBadRequest)
		return
	}

	var req models.Capitulo
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	capitulo, err := c.capituloService.Update(uint(capituloID), req)
	if err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, capitulo)
}

// DeleteCapitulo elimina un capítulo
func (c *CapituloController) DeleteCapitulo(ctx *gin.Context) {
	id := ctx.Param("id")
	capituloID, err := strconv.Atoi(id)
	if err != nil {
		utils.SendErrorResponse(ctx, utils.ErrBadRequest, http.StatusBadRequest)
		return
	}

	if err := c.capituloService.Delete(uint(capituloID)); err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Capítulo eliminado correctamente"})
}

// RegisterRoutes registra todas las rutas relacionadas con los capítulos
func (c *CapituloController) RegisterRoutes(router *gin.Engine) {
	capitulos := router.Group("/api/capitulos")
	{
		capitulos.Use(middleware.AuthMiddleware())
		capitulos.GET("/curso/:cursoId", c.GetCapitulosByCurso)
		capitulos.POST("", c.CreateCapitulo)
		capitulos.PUT("/:id", c.UpdateCapitulo)
		capitulos.DELETE("/:id", c.DeleteCapitulo)
	}
}
