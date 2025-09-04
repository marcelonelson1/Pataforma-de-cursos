package controllers

import (
	"portfolio-service/models"
	"portfolio-service/services"
	"portfolio-service/utils"
	"net/http"
	"strconv"
	"strings"
	"log"

	"github.com/gin-gonic/gin"
)

type PortfolioController struct {
	portfolioService *services.PortfolioService
	fileService      *services.FileService
}

// NewPortfolioController crea una nueva instancia de PortfolioController
func NewPortfolioController() *PortfolioController {
	return &PortfolioController{
		portfolioService: services.NewPortfolioService(),
		fileService:      services.NewFileService(),
	}
}

// GetAllProjects obtiene todos los proyectos activos (publico)
func (pc *PortfolioController) GetAllProjects(c *gin.Context) {
	projects, err := pc.portfolioService.GetAllProjects()
	if err != nil {
		utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Convertir a respuestas
	var responses []models.ProjectResponse
	for _, project := range projects {
		responses = append(responses, project.ToResponse())
	}

	utils.SendSuccessResponse(c, gin.H{
		"projects": responses,
		"total":    len(responses),
	})
}

// GetAllProjectsAdmin obtiene todos los proyectos para admin
func (pc *PortfolioController) GetAllProjectsAdmin(c *gin.Context) {
	projects, err := pc.portfolioService.GetAllProjectsAdmin()
	if err != nil {
		utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Convertir a respuestas
	var responses []models.ProjectResponse
	for _, project := range projects {
		responses = append(responses, project.ToResponse())
	}

	utils.SendSuccessResponse(c, gin.H{
		"projects": responses,
		"total":    len(responses),
	})
}

// GetProjectByID obtiene un proyecto especifico
func (pc *PortfolioController) GetProjectByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendErrorMessage(c, "ID invalido", http.StatusBadRequest)
		return
	}

	project, err := pc.portfolioService.GetProjectByID(uint(id))
	if err != nil {
		if err == utils.ErrResourceNotFound {
			utils.SendErrorResponse(c, err, http.StatusNotFound)
		} else {
			utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		}
		return
	}

	utils.SendSuccessResponse(c, project.ToResponse())
}

// GetProjectsByCategory obtiene proyectos por categoria
func (pc *PortfolioController) GetProjectsByCategory(c *gin.Context) {
	category := c.Param("category")

	projects, err := pc.portfolioService.GetProjectsByCategory(category)
	if err != nil {
		if err == utils.ErrInvalidCategory {
			utils.SendErrorResponse(c, err, http.StatusBadRequest)
		} else {
			utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		}
		return
	}

	// Convertir a respuestas
	var responses []models.ProjectResponse
	for _, project := range projects {
		responses = append(responses, project.ToResponse())
	}

	utils.SendSuccessResponse(c, gin.H{
		"category": category,
		"projects": responses,
		"total":    len(responses),
	})
}

// CreateProject crea un nuevo proyecto (admin)
func (pc *PortfolioController) CreateProject(c *gin.Context) {
	contentType := c.GetHeader("Content-Type")
	log.Printf("🔥 DEBUG CREATE: Content-Type: %s", contentType)
	log.Printf("🔥 DEBUG CREATE: CreateProject called")
	
	var req models.ProjectRequest
	var hasImage bool

	// Detectar si es JSON o Multipart
	if strings.Contains(contentType, "multipart/form-data") {
		log.Printf("🔥 DEBUG CREATE: Processing multipart form data")
		
		// Obtener datos del formulario
		req = models.ProjectRequest{
			Title:       c.PostForm("title"),
			Category:    c.PostForm("category"),
			Description: c.PostForm("description"),
			IsActive:    true,
		}
		
		// Verificar si hay imagen
		_, err := c.FormFile("image")
		hasImage = (err == nil)
		log.Printf("🔥 DEBUG CREATE: Has image: %v", hasImage)
		
	} else {
		log.Printf("🔥 DEBUG CREATE: Processing JSON data")
		
		// Obtener datos del JSON
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("🔥 DEBUG CREATE: Error binding JSON: %v", err)
			utils.SendValidationError(c, err)
			return
		}
		hasImage = false
	}

	log.Printf("🔥 DEBUG CREATE: Title: %s, Category: %s", req.Title, req.Category)

	// Validar datos requeridos
	if req.Title == "" {
		utils.SendErrorMessage(c, "title es requerido", http.StatusBadRequest)
		return
	}
	if req.Category == "" {
		utils.SendErrorMessage(c, "category es requerido", http.StatusBadRequest)
		return
	}

	// Crear proyecto básico
	project, err := pc.portfolioService.CreateProject(&req)
	if err != nil {
		log.Printf("🔥 DEBUG CREATE: Error creating project: %v", err)
		if err == utils.ErrInvalidCategory {
			utils.SendErrorResponse(c, err, http.StatusBadRequest)
		} else {
			utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		}
		return
	}

	// Si hay imagen, procesarla
	if hasImage {
		log.Printf("🔥 DEBUG CREATE: Processing image upload...")
		file, err := c.FormFile("image")
		if err == nil && file != nil {
			log.Printf("🔥 DEBUG CREATE: Found image file: %s, Size: %d", file.Filename, file.Size)
			imageURL, uploadErr := pc.fileService.UploadProjectImage(project.ID, file)
			if uploadErr != nil {
				log.Printf("🔥 DEBUG CREATE: Error uploading image: %v", uploadErr)
			} else {
				log.Printf("🔥 DEBUG CREATE: Image uploaded successfully: %s", imageURL)
				// Actualizar proyecto con URL de imagen
				project, _ = pc.portfolioService.UpdateProjectImage(project.ID, imageURL)
			}
		}
	}

	log.Printf("🔥 DEBUG CREATE: Project created successfully")
	utils.SendSuccessMessage(c, "Proyecto creado exitosamente", project.ToResponse())
}

// UpdateProject actualiza un proyecto (admin) - Maneja JSON y Multipart
func (pc *PortfolioController) UpdateProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendErrorMessage(c, "ID invalido", http.StatusBadRequest)
		return
	}

	contentType := c.GetHeader("Content-Type")
	log.Printf("🔥 DEBUG UPDATE: Content-Type: %s", contentType)

	var req models.ProjectRequest
	var hasImage bool

	// Detectar si es JSON o Multipart
	if strings.Contains(contentType, "multipart/form-data") {
		log.Printf("🔥 DEBUG UPDATE: Processing multipart form data")
		
		// Obtener datos del formulario
		req = models.ProjectRequest{
			Title:       c.PostForm("title"),
			Category:    c.PostForm("category"),
			Description: c.PostForm("description"),
			IsActive:    true,
		}
		
		// Verificar si hay imagen
		_, err := c.FormFile("image")
		hasImage = (err == nil)
		log.Printf("🔥 DEBUG UPDATE: Has image: %v", hasImage)
		
	} else {
		log.Printf("🔥 DEBUG UPDATE: Processing JSON data")
		
		// Obtener datos del JSON
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("🔥 DEBUG UPDATE: Error binding JSON: %v", err)
			utils.SendValidationError(c, err)
			return
		}
		hasImage = false
	}

	log.Printf("🔥 DEBUG UPDATE: Title: %s, Category: %s", req.Title, req.Category)

	// Validar datos requeridos
	if req.Title == "" {
		utils.SendErrorMessage(c, "title es requerido", http.StatusBadRequest)
		return
	}
	if req.Category == "" {
		utils.SendErrorMessage(c, "category es requerido", http.StatusBadRequest)
		return
	}

	// Actualizar proyecto básico
	project, err := pc.portfolioService.UpdateProject(uint(id), &req)
	if err != nil {
		log.Printf("🔥 DEBUG UPDATE: Error updating project: %v", err)
		if err == utils.ErrResourceNotFound {
			utils.SendErrorResponse(c, err, http.StatusNotFound)
		} else if err == utils.ErrInvalidCategory {
			utils.SendErrorResponse(c, err, http.StatusBadRequest)
		} else {
			utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		}
		return
	}

	// Si hay imagen, procesarla
	if hasImage {
		log.Printf("🔥 DEBUG UPDATE: Processing image upload...")
		file, err := c.FormFile("image")
		if err == nil && file != nil {
			log.Printf("🔥 DEBUG UPDATE: Found image file: %s, Size: %d", file.Filename, file.Size)
			imageURL, uploadErr := pc.fileService.UploadProjectImage(uint(id), file)
			if uploadErr != nil {
				log.Printf("🔥 DEBUG UPDATE: Error uploading image: %v", uploadErr)
			} else {
				log.Printf("🔥 DEBUG UPDATE: Image uploaded successfully: %s", imageURL)
				// Actualizar proyecto con URL de imagen
				project, _ = pc.portfolioService.UpdateProjectImage(uint(id), imageURL)
			}
		}
	}

	log.Printf("🔥 DEBUG UPDATE: Project updated successfully")
	utils.SendSuccessMessage(c, "Proyecto actualizado exitosamente", project.ToResponse())
}

// DeleteProject elimina un proyecto (admin)
func (pc *PortfolioController) DeleteProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendErrorMessage(c, "ID invalido", http.StatusBadRequest)
		return
	}

	if err := pc.portfolioService.DeleteProject(uint(id)); err != nil {
		if err == utils.ErrResourceNotFound {
			utils.SendErrorResponse(c, err, http.StatusNotFound)
		} else {
			utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		}
		return
	}

	utils.SendSuccessMessage(c, "Proyecto eliminado exitosamente", nil)
}

// UploadProjectImage sube imagen para un proyecto (admin)
func (pc *PortfolioController) UploadProjectImage(c *gin.Context) {
	// Obtener proyecto ID
	projectIDStr := c.PostForm("project_id")
	if projectIDStr == "" {
		utils.SendErrorMessage(c, "project_id requerido", http.StatusBadRequest)
		return
	}

	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.SendErrorMessage(c, "project_id invalido", http.StatusBadRequest)
		return
	}

	// Obtener archivo
	file, err := c.FormFile("image")
	if err != nil {
		utils.SendErrorMessage(c, "archivo imagen requerido", http.StatusBadRequest)
		return
	}

	// Subir imagen
	imageURL, err := pc.fileService.UploadProjectImage(uint(projectID), file)
	if err != nil {
		if err == utils.ErrResourceNotFound {
			utils.SendErrorResponse(c, err, http.StatusNotFound)
		} else if err == utils.ErrInvalidFileType {
			utils.SendErrorResponse(c, err, http.StatusBadRequest)
		} else {
			utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		}
		return
	}

	utils.SendSuccessMessage(c, "Imagen subida exitosamente", gin.H{
		"image_url": imageURL,
		"project_id": projectID,
	})
}

// DeleteProjectImage elimina imagen de un proyecto (admin)
func (pc *PortfolioController) DeleteProjectImage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendErrorMessage(c, "ID invalido", http.StatusBadRequest)
		return
	}

	if err := pc.fileService.DeleteProjectImage(uint(id)); err != nil {
		if err == utils.ErrResourceNotFound {
			utils.SendErrorResponse(c, err, http.StatusNotFound)
		} else {
			utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		}
		return
	}

	utils.SendSuccessMessage(c, "Imagen eliminada exitosamente", nil)
}

// ReorderProjects reordena proyectos (admin)
func (pc *PortfolioController) ReorderProjects(c *gin.Context) {
	var req models.ReorderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, err)
		return
	}

	if err := pc.portfolioService.ReorderProjects(req.ProjectIDs); err != nil {
		if err == utils.ErrResourceNotFound {
			utils.SendErrorResponse(c, err, http.StatusNotFound)
		} else {
			utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		}
		return
	}

	utils.SendSuccessMessage(c, "Proyectos reordenados exitosamente", nil)
}

// GetPortfolioStats obtiene estadisticas del portafolio (admin)
func (pc *PortfolioController) GetPortfolioStats(c *gin.Context) {
	stats, err := pc.portfolioService.GetPortfolioStats()
	if err != nil {
		utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, stats)
}

// ToggleProjectStatus cambia estado activo/inactivo (admin)
func (pc *PortfolioController) ToggleProjectStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendErrorMessage(c, "ID invalido", http.StatusBadRequest)
		return
	}

	project, err := pc.portfolioService.ToggleProjectStatus(uint(id))
	if err != nil {
		if err == utils.ErrResourceNotFound {
			utils.SendErrorResponse(c, err, http.StatusNotFound)
		} else {
			utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		}
		return
	}

	status := "inactivo"
	if project.IsActive {
		status = "activo"
	}

	utils.SendSuccessMessage(c, "Estado del proyecto actualizado a: "+status, project.ToResponse())
}

// GetCategories obtiene todas las categorias disponibles
func (pc *PortfolioController) GetCategories(c *gin.Context) {
	categories, err := pc.portfolioService.GetUniqueCategories()
	if err != nil {
		utils.SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, gin.H{
		"categories": categories,
		"valid_categories": models.ValidCategories,
	})
}

// HealthCheck endpoint de verificacion de salud
func (pc *PortfolioController) HealthCheck(c *gin.Context) {
	utils.SendSuccessResponse(c, gin.H{
		"status":  "ok",
		"service": "portfolio-service",
		"version": "1.0.0",
	})
}