package services

import (
	"curso-platform/config"
	"curso-platform/models"
	"curso-platform/utils"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// PortfolioService maneja la lógica relacionada con el portafolio de proyectos
type PortfolioService struct{}

// NewPortfolioService crea una nueva instancia de PortfolioService
func NewPortfolioService() *PortfolioService {
	return &PortfolioService{}
}

// GetAllProjects obtiene todos los proyectos del portfolio
func (s *PortfolioService) GetAllProjects(page, limit int, category string) ([]models.ProjectPortfolio, int64, error) {
	var projects []models.ProjectPortfolio
	query := config.DB.Model(&models.ProjectPortfolio{}).Order("order_index asc")
	
	// Aplicar filtro si existe
	if category != "" {
		query = query.Where("category = ?", category)
	}
	
	// Total para paginación
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, utils.ErrDatabaseError
	}
	
	// Calcular offset para paginación
	offset := (page - 1) * limit
	
	// Obtener proyectos con paginación
	if err := query.Limit(limit).Offset(offset).Find(&projects).Error; err != nil {
		return nil, 0, utils.ErrDatabaseError
	}

	return projects, total, nil
}

// GetProjectByID obtiene un proyecto específico por ID
func (s *PortfolioService) GetProjectByID(id uint) (*models.ProjectPortfolio, error) {
	var project models.ProjectPortfolio

	if err := config.DB.First(&project, id).Error; err != nil {
		return nil, utils.ErrResourceNotFound
	}

	return &project, nil
}

// GetProjectsByCategory obtiene proyectos filtrados por categoría
func (s *PortfolioService) GetProjectsByCategory(category string, page, limit int) ([]models.ProjectPortfolio, int64, error) {
	var projects []models.ProjectPortfolio
	
	// Query con filtro por categoría
	query := config.DB.Where("category = ?", category).Order("order_index asc")
	
	// Total para paginación
	var total int64
	if err := query.Model(&models.ProjectPortfolio{}).Count(&total).Error; err != nil {
		return nil, 0, utils.ErrDatabaseError
	}
	
	// Calcular offset para paginación
	offset := (page - 1) * limit
	
	// Obtener proyectos con paginación
	if err := query.Limit(limit).Offset(offset).Find(&projects).Error; err != nil {
		return nil, 0, utils.ErrDatabaseError
	}

	return projects, total, nil
}

// CreateProject crea un nuevo proyecto en el portfolio
func (s *PortfolioService) CreateProject(title, category, description string, imageFile *os.File, fileName string) (*models.ProjectPortfolio, error) {
	// Validar campos obligatorios
	if title == "" || category == "" {
		return nil, utils.ErrMissingFields
	}

	// Crear directorio si no existe
	portfolioDir := "./static/portfolio"
	if err := utils.CreateDirIfNotExists(portfolioDir); err != nil {
		return nil, err
	}

	// Generar nombre único para el archivo
	filename := utils.GenerateUniqueFileName(fileName)
	filePath := filepath.Join(portfolioDir, filename)

	// Leer contenido del archivo original
	fileData, err := os.ReadFile(imageFile.Name())
	if err != nil {
		return nil, fmt.Errorf("error al leer archivo imagen: %v", err)
	}

	// Escribir a nuevo archivo
	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		return nil, fmt.Errorf("error al guardar imagen: %v", err)
	}

	// Obtener el índice máximo actual para ordenar
	var maxOrderIndex int
	config.DB.Model(&models.ProjectPortfolio{}).Select("COALESCE(MAX(order_index), 0)").Scan(&maxOrderIndex)

	// Crear el nuevo proyecto con URL relativa para el frontend
	project := models.ProjectPortfolio{
		Title:       title,
		Category:    category,
		Description: description,
		ImageURL:    "/static/portfolio/" + filename,
		OrderIndex:  maxOrderIndex + 1,
	}

	// Guardar en la base de datos
	if err := config.DB.Create(&project).Error; err != nil {
		// Si hay error, eliminar la imagen subida
		os.Remove(filePath)
		return nil, utils.ErrDatabaseError
	}

	return &project, nil
}

// UpdateProject actualiza un proyecto existente
func (s *PortfolioService) UpdateProject(id uint, title, category, description string, imageFile *os.File, fileName string) (*models.ProjectPortfolio, error) {
	// Verificar si el proyecto existe
	var project models.ProjectPortfolio
	if err := config.DB.First(&project, id).Error; err != nil {
		return nil, utils.ErrResourceNotFound
	}

	// Actualizar campos si se proporcionan
	if title != "" {
		project.Title = title
	}
	if category != "" {
		project.Category = category
	}
	project.Description = description // Permitir actualizar a vacío

	// Verificar si se está subiendo una nueva imagen
	if imageFile != nil {
		// Si se proporcionó una nueva imagen, eliminar la anterior si existe
		if project.ImageURL != "" {
			oldImagePath := "." + project.ImageURL
			// Solo intentar eliminar si el archivo existe
			if _, err := os.Stat(oldImagePath); err == nil {
				os.Remove(oldImagePath)
			}
		}

		// Crear directorio si no existe
		portfolioDir := "./static/portfolio"
		if err := utils.CreateDirIfNotExists(portfolioDir); err != nil {
			return nil, err
		}

		// Generar nombre único para el archivo
		filename := utils.GenerateUniqueFileName(fileName)
		filePath := filepath.Join(portfolioDir, filename)

		// Leer contenido del archivo original
		fileData, err := os.ReadFile(imageFile.Name())
		if err != nil {
			return nil, fmt.Errorf("error al leer archivo imagen: %v", err)
		}

		// Escribir a nuevo archivo
		if err := os.WriteFile(filePath, fileData, 0644); err != nil {
			return nil, fmt.Errorf("error al guardar imagen: %v", err)
		}

		// Actualizar URL de la imagen
		project.ImageURL = "/static/portfolio/" + filename
	}

	// Guardar cambios en la base de datos
	if err := config.DB.Save(&project).Error; err != nil {
		return nil, utils.ErrDatabaseError
	}

	return &project, nil
}

// DeleteProject elimina un proyecto del portfolio
func (s *PortfolioService) DeleteProject(id uint) error {
	var project models.ProjectPortfolio
	if err := config.DB.First(&project, id).Error; err != nil {
		return utils.ErrResourceNotFound
	}

	// Eliminar la imagen asociada si existe
	if project.ImageURL != "" {
		imagePath := "." + project.ImageURL
		// Solo intentar eliminar si el archivo existe
		if _, err := os.Stat(imagePath); err == nil {
			os.Remove(imagePath)
		}
	}

	// Eliminar el proyecto de la base de datos
	if err := config.DB.Delete(&project).Error; err != nil {
		return utils.ErrDatabaseError
	}

	// Reordenar índices de otros proyectos
	if err := config.DB.Exec("UPDATE project_portfolios SET order_index = order_index - 1 WHERE order_index > ?", project.OrderIndex).Error; err != nil {
		// Solo loggeamos el error, no detenemos la ejecución
		log.Printf("Error al reordenar proyectos: %v", err)
	}

	return nil
}

// ReorderProjects actualiza el orden de los proyectos
func (s *PortfolioService) ReorderProjects(projectIds []uint) error {
	// Validar que todos los IDs existan
	for _, id := range projectIds {
		var project models.ProjectPortfolio
		if err := config.DB.First(&project, id).Error; err != nil {
			return utils.ErrResourceNotFound
		}
	}

	// Iniciar transacción para actualizar todos los índices
	tx := config.DB.Begin()

	for i, projectID := range projectIds {
		if err := tx.Model(&models.ProjectPortfolio{}).Where("id = ?", projectID).Update("order_index", i+1).Error; err != nil {
			tx.Rollback()
			return utils.ErrDatabaseError
		}
	}

	// Confirmar transacción
	if err := tx.Commit().Error; err != nil {
		return utils.ErrDatabaseError
	}

	return nil
}

// GetPortfolioStats obtiene estadísticas del portfolio para el panel de administración
func (s *PortfolioService) GetPortfolioStats() (map[string]interface{}, error) {
	// Total de proyectos
	var totalProjects int64
	if err := config.DB.Model(&models.ProjectPortfolio{}).Count(&totalProjects).Error; err != nil {
		return nil, utils.ErrDatabaseError
	}
	
	// Proyectos por categoría
	type CategoryCount struct {
		Category string `json:"category"`
		Count    int    `json:"count"`
	}
	
	var categoryCounts []CategoryCount
	if err := config.DB.Model(&models.ProjectPortfolio{}).
		Select("category, COUNT(*) as count").
		Group("category").
		Order("count DESC").
		Scan(&categoryCounts).Error; err != nil {
		return nil, utils.ErrDatabaseError
	}
	
	// Proyectos recientes (últimos 5)
	var recentProjects []models.ProjectPortfolio
	if err := config.DB.Order("created_at DESC").Limit(5).Find(&recentProjects).Error; err != nil {
		return nil, utils.ErrDatabaseError
	}
	
	// Proyecto más reciente
	var latestProject models.ProjectPortfolio
	hasLatest := config.DB.Order("created_at DESC").First(&latestProject).Error == nil
	
	// Respuesta con todas las estadísticas
	stats := map[string]interface{}{
		"totalProjects": totalProjects,
		"byCategory":    categoryCounts,
		"recent":        recentProjects,
	}
	
	if hasLatest {
		stats["latest"] = latestProject
	}
	
	return stats, nil
}