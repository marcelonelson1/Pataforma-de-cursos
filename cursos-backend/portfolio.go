package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ProjectPortfolio representa un proyecto en el portafolio
type ProjectPortfolio struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"size:200;not null" json:"title"`
	Category    string    `gorm:"size:50;not null" json:"category"`
	Description string    `gorm:"size:500" json:"description"`
	ImageURL    string    `gorm:"size:255" json:"image_url"`
	OrderIndex  int       `gorm:"default:0" json:"order_index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// getAllProjects obtiene todos los proyectos del portfolio
func getAllProjects(c *gin.Context) {
	var projects []ProjectPortfolio

	// Paginación opcional
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset := (page - 1) * limit

	// Filtro por categoría opcional
	categoryFilter := c.Query("category")
	
	// Construir query
	query := db.Model(&ProjectPortfolio{}).Order("order_index asc")
	
	// Aplicar filtro si existe
	if categoryFilter != "" {
		query = query.Where("category = ?", categoryFilter)
	}
	
	// Total para paginación
	var total int64
	if err := query.Count(&total).Error; err != nil {
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}
	
	// Obtener proyectos con paginación
	if err := query.Limit(limit).Offset(offset).Find(&projects).Error; err != nil {
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}

	// Registrar actividad si hay un usuario autenticado
	if userValue, exists := c.Get("user"); exists {
		user := userValue.(Usuario)
		logActivity(c, user.ID, "view_portfolio", fmt.Sprintf("Visualización de portfolio (total: %d proyectos)", len(projects)))
	}

	// Si se solicita paginación explícitamente
	if c.Query("page") != "" {
		SendSuccessResponse(c, gin.H{
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
	SendSuccessResponse(c, projects)
}

// getProjectById obtiene un proyecto específico por ID
func getProjectById(c *gin.Context) {
	id := c.Param("id")
	var project ProjectPortfolio

	result := db.First(&project, id)
	if result.Error != nil {
		SendErrorResponse(c, ErrResourceNotFound, http.StatusNotFound)
		return
	}

	// Registrar actividad si hay un usuario autenticado
	if userValue, exists := c.Get("user"); exists {
		user := userValue.(Usuario)
		logActivity(c, user.ID, "view_project_detail", fmt.Sprintf("Visualización de detalle de proyecto ID: %s, Título: %s", id, project.Title))
	}

	SendSuccessResponse(c, project)
}

// getProjectsByCategory obtiene proyectos filtrados por categoría
func getProjectsByCategory(c *gin.Context) {
	category := c.Param("category")
	var projects []ProjectPortfolio

	// Paginación opcional
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset := (page - 1) * limit

	// Query con filtro por categoría
	query := db.Where("category = ?", category).Order("order_index asc")
	
	// Total para paginación
	var total int64
	if err := query.Model(&ProjectPortfolio{}).Count(&total).Error; err != nil {
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}
	
	// Obtener proyectos con paginación
	if err := query.Limit(limit).Offset(offset).Find(&projects).Error; err != nil {
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}

	// Registrar actividad si hay un usuario autenticado
	if userValue, exists := c.Get("user"); exists {
		user := userValue.(Usuario)
		logActivity(c, user.ID, "filter_portfolio", fmt.Sprintf("Filtrado de portfolio por categoría: %s (total: %d proyectos)", category, len(projects)))
	}

	// Si se solicita paginación explícitamente
	if c.Query("page") != "" {
		SendSuccessResponse(c, gin.H{
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
	SendSuccessResponse(c, gin.H{
		"projects": projects,
		"category": category,
	})
}

// createProject crea un nuevo proyecto en el portfolio
func createProject(c *gin.Context) {
	// Solo administradores pueden crear proyectos
	userValue, _ := c.Get("user")
	user := userValue.(Usuario)
	
	if user.Role != "admin" {
		SendErrorResponse(c, ErrUnauthorized, http.StatusForbidden)
		return
	}

	// Obtener datos del formulario
	title := c.PostForm("title")
	category := c.PostForm("category")
	description := c.PostForm("description")

	// Validar campos obligatorios
	if title == "" || category == "" {
		SendErrorResponse(c, ErrMissingFields, http.StatusBadRequest)
		return
	}

	// Manejar la carga de la imagen
	file, err := c.FormFile("image")
	if err != nil {
		SendErrorResponse(c, ErrBadRequest, http.StatusBadRequest)
		return
	}

	// Asegurar que el directorio de portfolio existe
	portfolioDir := "./static/portfolio"
	if _, err := os.Stat(portfolioDir); os.IsNotExist(err) {
		if err := os.MkdirAll(portfolioDir, 0755); err != nil {
			log.Printf("Error al crear directorio de portfolio: %v", err)
			SendErrorResponse(c, ErrServerError, http.StatusInternalServerError)
			return
		}
	}

	// Crear nombre de archivo único
	ext := filepath.Ext(file.Filename)
	fileName := uuid.New().String() + ext
	filePath := filepath.Join(portfolioDir, fileName)

	// Guardar archivo en el servidor
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		log.Printf("Error al guardar imagen: %v", err)
		SendErrorResponse(c, ErrServerError, http.StatusInternalServerError)
		return
	}

	// Obtener el índice máximo actual para ordenar
	var maxOrderIndex int
	db.Model(&ProjectPortfolio{}).Select("COALESCE(MAX(order_index), 0)").Scan(&maxOrderIndex)

	// Crear el nuevo proyecto con URL relativa para el frontend
	project := ProjectPortfolio{
		Title:       title,
		Category:    category,
		Description: description,
		ImageURL:    "/static/portfolio/" + fileName, // URL para acceso desde frontend
		OrderIndex:  maxOrderIndex + 1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Guardar en la base de datos
	result := db.Create(&project)
	if result.Error != nil {
		// Si hay error, eliminar la imagen subida
		os.Remove(filePath)
		log.Printf("Error al crear proyecto en BD: %v", result.Error)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}

	// Registrar actividad
	logActivity(c, user.ID, "create_portfolio", fmt.Sprintf("Creado proyecto '%s' en categoría '%s'", title, category))

	SendSuccessResponse(c, gin.H{
		"message": "Proyecto creado con éxito",
		"project": project,
	})
}

// updateProject actualiza un proyecto existente
func updateProject(c *gin.Context) {
	// Solo administradores pueden actualizar proyectos
	userValue, _ := c.Get("user")
	user := userValue.(Usuario)
	
	if user.Role != "admin" {
		SendErrorResponse(c, ErrUnauthorized, http.StatusForbidden)
		return
	}

	id := c.Param("id")
	var project ProjectPortfolio

	// Verificar si el proyecto existe
	result := db.First(&project, id)
	if result.Error != nil {
		SendErrorResponse(c, ErrResourceNotFound, http.StatusNotFound)
		return
	}

	// Guardar datos anteriores para el log
	oldTitle := project.Title
	oldCategory := project.Category

	// Obtener datos del formulario
	title := c.PostForm("title")
	category := c.PostForm("category")
	description := c.PostForm("description")

	// Actualizar campos si se proporcionan
	if title != "" {
		project.Title = title
	}
	if category != "" {
		project.Category = category
	}
	project.Description = description // Permitir actualizar a vacío

	// Verificar si se está subiendo una nueva imagen
	file, err := c.FormFile("image")
	if err == nil {
		// Asegurar que el directorio de portfolio existe
		portfolioDir := "./static/portfolio"
		if _, err := os.Stat(portfolioDir); os.IsNotExist(err) {
			if err := os.MkdirAll(portfolioDir, 0755); err != nil {
				log.Printf("Error al crear directorio de portfolio: %v", err)
				SendErrorResponse(c, ErrServerError, http.StatusInternalServerError)
				return
			}
		}

		// Si se proporcionó una nueva imagen, eliminar la anterior si existe
		if project.ImageURL != "" {
			oldImagePath := "." + project.ImageURL
			// Solo intentar eliminar si el archivo existe
			if _, err := os.Stat(oldImagePath); err == nil {
				os.Remove(oldImagePath)
			}
		}

		// Crear nombre de archivo único
		ext := filepath.Ext(file.Filename)
		fileName := uuid.New().String() + ext
		filePath := filepath.Join(portfolioDir, fileName)

		// Guardar archivo en el servidor
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			log.Printf("Error al guardar nueva imagen: %v", err)
			SendErrorResponse(c, ErrServerError, http.StatusInternalServerError)
			return
		}

		// Actualizar URL de la imagen
		project.ImageURL = "/static/portfolio/" + fileName
	}

	// Actualizar fecha de modificación
	project.UpdatedAt = time.Now()

	// Guardar cambios en la base de datos
	if err := db.Save(&project).Error; err != nil {
		log.Printf("Error al actualizar proyecto en BD: %v", err)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}

	// Registrar actividad
	logActivity(c, user.ID, "update_portfolio", fmt.Sprintf("Actualizado proyecto ID: %s, Título: '%s' -> '%s', Categoría: '%s' -> '%s'", 
		id, oldTitle, project.Title, oldCategory, project.Category))

	SendSuccessResponse(c, gin.H{
		"message": "Proyecto actualizado con éxito",
		"project": project,
	})
}

// deleteProject elimina un proyecto del portfolio
func deleteProject(c *gin.Context) {
	// Solo administradores pueden eliminar proyectos
	userValue, _ := c.Get("user")
	user := userValue.(Usuario)
	
	if user.Role != "admin" {
		SendErrorResponse(c, ErrUnauthorized, http.StatusForbidden)
		return
	}

	id := c.Param("id")
	var project ProjectPortfolio

	// Verificar si el proyecto existe
	result := db.First(&project, id)
	if result.Error != nil {
		SendErrorResponse(c, ErrResourceNotFound, http.StatusNotFound)
		return
	}

	// Guardar datos para el log
	projectTitle := project.Title
	projectCategory := project.Category

	// Eliminar la imagen asociada si existe
	if project.ImageURL != "" {
		imagePath := "." + project.ImageURL
		// Solo intentar eliminar si el archivo existe
		if _, err := os.Stat(imagePath); err == nil {
			os.Remove(imagePath)
		}
	}

	// Eliminar el proyecto de la base de datos
	if err := db.Delete(&project).Error; err != nil {
		log.Printf("Error al eliminar proyecto de BD: %v", err)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}

	// Reordenar índices de otros proyectos
	if err := db.Exec("UPDATE project_portfolios SET order_index = order_index - 1 WHERE order_index > ?", project.OrderIndex).Error; err != nil {
		// Solo loggeamos el error, no detenemos la ejecución
		log.Printf("Error al reordenar proyectos: %v", err)
	}

	// Registrar actividad
	logActivity(c, user.ID, "delete_portfolio", fmt.Sprintf("Eliminado proyecto ID: %s, Título: '%s', Categoría: '%s'", 
		id, projectTitle, projectCategory))

	SendSuccessResponse(c, gin.H{
		"message": "Proyecto eliminado con éxito",
	})
}

// reorderProjects actualiza el orden de los proyectos
func reorderProjects(c *gin.Context) {
	// Solo administradores pueden reordenar proyectos
	userValue, _ := c.Get("user")
	user := userValue.(Usuario)
	
	if user.Role != "admin" {
		SendErrorResponse(c, ErrUnauthorized, http.StatusForbidden)
		return
	}

	var requestData struct {
		ProjectIds []uint `json:"projectIds" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		SendErrorResponse(c, ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	// Validar que todos los IDs existan
	for _, id := range requestData.ProjectIds {
		var project ProjectPortfolio
		if err := db.First(&project, id).Error; err != nil {
			SendErrorResponse(c, ErrResourceNotFound, http.StatusBadRequest)
			return
		}
	}

	// Iniciar transacción para actualizar todos los índices
	tx := db.Begin()

	for i, projectID := range requestData.ProjectIds {
		if err := tx.Model(&ProjectPortfolio{}).Where("id = ?", projectID).Update("order_index", i+1).Error; err != nil {
			tx.Rollback()
			log.Printf("Error al actualizar orden de proyecto ID %d: %v", projectID, err)
			SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
			return
		}
	}

	// Confirmar transacción
	if err := tx.Commit().Error; err != nil {
		log.Printf("Error al confirmar transacción de reordenamiento: %v", err)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}

	// Registrar actividad
	logActivity(c, user.ID, "reorder_portfolio", fmt.Sprintf("Reordenados %d proyectos del portfolio", len(requestData.ProjectIds)))

	SendSuccessResponse(c, gin.H{
		"message": "Orden actualizado con éxito",
	})
}

// getPortfolioStats obtiene estadísticas del portfolio para el panel de administración
func getPortfolioStats(c *gin.Context) {
	// Solo administradores pueden ver estadísticas
	userValue, _ := c.Get("user")
	user := userValue.(Usuario)
	
	if user.Role != "admin" {
		SendErrorResponse(c, ErrUnauthorized, http.StatusForbidden)
		return
	}
	
	// Registrar actividad
	logActivity(c, user.ID, "view_portfolio_stats", "Visualización de estadísticas del portfolio")
	
	// Total de proyectos
	var totalProjects int64
	if err := db.Model(&ProjectPortfolio{}).Count(&totalProjects).Error; err != nil {
		log.Printf("Error al contar total de proyectos: %v", err)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}
	
	// Proyectos por categoría
	type CategoryCount struct {
		Category string `json:"category"`
		Count    int    `json:"count"`
	}
	
	var categoryCounts []CategoryCount
	if err := db.Model(&ProjectPortfolio{}).
		Select("category, COUNT(*) as count").
		Group("category").
		Order("count DESC").
		Scan(&categoryCounts).Error; err != nil {
		log.Printf("Error al obtener conteo por categorías: %v", err)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}
	
	// Proyectos recientes (últimos 5)
	var recentProjects []ProjectPortfolio
	if err := db.Order("created_at DESC").Limit(5).Find(&recentProjects).Error; err != nil {
		log.Printf("Error al obtener proyectos recientes: %v", err)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}
	
	// Proyecto más reciente
	var latestProject ProjectPortfolio
	hasLatest := false
	if err := db.Order("created_at DESC").First(&latestProject).Error; err == nil {
		hasLatest = true
	}
	
	// Respuesta con todas las estadísticas
	response := gin.H{
		"totalProjects": totalProjects,
		"byCategory":    categoryCounts,
		"recent":        recentProjects,
	}
	
	if hasLatest {
		response["latest"] = latestProject
	}
	
	SendSuccessResponse(c, gin.H{
		"stats": response,
	})
}