package controllers

import (
	"curso-platform/middleware"
	"curso-platform/services"
	"curso-platform/utils"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PortfolioController gestiona las operaciones relacionadas con el portafolio de proyectos
type PortfolioController struct {
	portfolioService *services.PortfolioService
}

// NewPortfolioController crea una nueva instancia del controlador de portafolio
func NewPortfolioController(portfolioService *services.PortfolioService) *PortfolioController {
	return &PortfolioController{
		portfolioService: portfolioService,
	}
}

// GetAllProjects obtiene todos los proyectos del portfolio
func (c *PortfolioController) GetAllProjects(ctx *gin.Context) {
	// Paginación opcional
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "50"))
	
	// Filtro por categoría opcional
	category := ctx.Query("category")
	
	// Obtener proyectos
	projects, total, err := c.portfolioService.GetAllProjects(page, limit, category)
	if err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}
	
	// Si se solicita paginación explícitamente
	if ctx.Query("page") != "" {
		utils.SendSuccessResponse(ctx, gin.H{
			"projects": projects,
			"pagination": gin.H{
				"total": total,
				"page":  page,
				"limit": limit,
				"pages": (total + int64(limit) - 1) / int64(limit),
			},
		})
		return
	}
	
	// Respuesta estándar sin paginación
	utils.SendSuccessResponse(ctx, projects)
}

// GetProjectById obtiene un proyecto específico por ID
func (c *PortfolioController) GetProjectById(ctx *gin.Context) {
	id := ctx.Param("id")
	projectID, err := strconv.Atoi(id)
	if err != nil {
		utils.SendErrorResponse(ctx, utils.ErrBadRequest, http.StatusBadRequest)
		return
	}
	
	project, err := c.portfolioService.GetProjectByID(uint(projectID))
	if err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusNotFound)
		return
	}
	
	utils.SendSuccessResponse(ctx, project)
}

// GetProjectsByCategory obtiene proyectos filtrados por categoría
func (c *PortfolioController) GetProjectsByCategory(ctx *gin.Context) {
	category := ctx.Param("category")
	
	// Paginación opcional
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "50"))
	
	projects, total, err := c.portfolioService.GetProjectsByCategory(category, page, limit)
	if err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}
	
	// Si se solicita paginación explícitamente
	if ctx.Query("page") != "" {
		utils.SendSuccessResponse(ctx, gin.H{
			"projects": projects,
			"category": category,
			"pagination": gin.H{
				"total": total,
				"page":  page,
				"limit": limit,
				"pages": (total + int64(limit) - 1) / int64(limit),
			},
		})
		return
	}
	
	// Respuesta estándar sin paginación
	utils.SendSuccessResponse(ctx, gin.H{
		"projects": projects,
		"category": category,
	})
}

// CreateProject crea un nuevo proyecto en el portfolio
func (c *PortfolioController) CreateProject(ctx *gin.Context) {
	// Obtener datos del formulario
	title := ctx.PostForm("title")
	category := ctx.PostForm("category")
	description := ctx.PostForm("description")
	
	// Validar campos obligatorios
	if title == "" || category == "" {
		utils.SendErrorResponse(ctx, utils.ErrMissingFields, http.StatusBadRequest)
		return
	}
	
	// Manejar la carga de la imagen
	file, err := ctx.FormFile("image")
	if err != nil {
		utils.SendErrorResponse(ctx, utils.ErrBadRequest, http.StatusBadRequest)
		return
	}
	
	// Crear archivo temporal para procesar
	tempFile, err := os.CreateTemp("", "portfolio-*.tmp")
	if err != nil {
		utils.SendErrorResponse(ctx, utils.ErrServerError, http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()
	
	// Guardar FormFile en archivo temporal
	if err := ctx.SaveUploadedFile(file, tempFile.Name()); err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}
	
	// Crear proyecto
	project, err := c.portfolioService.CreateProject(title, category, description, tempFile, file.Filename)
	if err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}
	
	utils.SendSuccessResponse(ctx, gin.H{
		"message": "Proyecto creado con éxito",
		"project": project,
	})
}

// UpdateProject actualiza un proyecto existente
func (c *PortfolioController) UpdateProject(ctx *gin.Context) {
	id := ctx.Param("id")
	projectID, err := strconv.Atoi(id)
	if err != nil {
		utils.SendErrorResponse(ctx, utils.ErrBadRequest, http.StatusBadRequest)
		return
	}
	
	// Obtener datos del formulario
	title := ctx.PostForm("title")
	category := ctx.PostForm("category")
	description := ctx.PostForm("description")
	
	// Manejar la carga de la imagen si hay una nueva
	var imageFile *os.File
	file, err := ctx.FormFile("image")
	if err == nil {
		// Si hay un archivo de imagen, crear un archivo temporal
		tempFile, err := os.CreateTemp("", "portfolio-*.tmp")
		if err != nil {
			utils.SendErrorResponse(ctx, utils.ErrServerError, http.StatusInternalServerError)
			return
		}
		defer os.Remove(tempFile.Name())
		defer tempFile.Close()
		
		// Guardar FormFile en archivo temporal
		if err := ctx.SaveUploadedFile(file, tempFile.Name()); err != nil {
			utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
			return
		}
		
		imageFile = tempFile
	}
	
	// Actualizar proyecto
	project, err := c.portfolioService.UpdateProject(uint(projectID), title, category, description, imageFile, file.Filename)
	if err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}
	
	utils.SendSuccessResponse(ctx, gin.H{
		"message": "Proyecto actualizado con éxito",
		"project": project,
	})
}

// DeleteProject elimina un proyecto
func (c *PortfolioController) DeleteProject(ctx *gin.Context) {
	id := ctx.Param("id")
	projectID, err := strconv.Atoi(id)
	if err != nil {
		utils.SendErrorResponse(ctx, utils.ErrBadRequest, http.StatusBadRequest)
		return
	}
	
	if err := c.portfolioService.DeleteProject(uint(projectID)); err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}
	
	utils.SendSuccessResponse(ctx, gin.H{
		"message": "Proyecto eliminado con éxito",
	})
}

// ReorderProjects actualiza el orden de los proyectos
func (c *PortfolioController) ReorderProjects(ctx *gin.Context) {
	var requestData struct {
		ProjectIds []uint `json:"projectIds" binding:"required"`
	}
	
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		utils.SendErrorResponse(ctx, utils.ErrInvalidJSON, http.StatusBadRequest)
		return
	}
	
	if err := c.portfolioService.ReorderProjects(requestData.ProjectIds); err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}
	
	utils.SendSuccessResponse(ctx, gin.H{
		"message": "Orden actualizado con éxito",
	})
}

// GetPortfolioStats obtiene estadísticas del portfolio
func (c *PortfolioController) GetPortfolioStats(ctx *gin.Context) {
	stats, err := c.portfolioService.GetPortfolioStats()
	if err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}
	
	utils.SendSuccessResponse(ctx, stats)
}

// RegisterRoutes registra todas las rutas relacionadas con el portfolio
func (c *PortfolioController) RegisterRoutes(router *gin.Engine) {
	portfolio := router.Group("/api/portfolio")
	{
		portfolio.GET("", c.GetAllProjects)
		portfolio.GET("/:id", c.GetProjectById)
		portfolio.GET("/category/:category", c.GetProjectsByCategory)
		
		// Rutas protegidas para administración (solo admin)
		admin := portfolio.Group("")
		admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			admin.POST("", c.CreateProject)
			admin.PUT("/:id", c.UpdateProject)
			admin.DELETE("/:id", c.DeleteProject)
			admin.POST("/reorder", c.ReorderProjects)
			admin.GET("/stats", c.GetPortfolioStats)
		}
	}
}