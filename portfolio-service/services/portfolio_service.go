package services

import (
	"portfolio-service/models"
	"portfolio-service/utils"
	"log"

	"gorm.io/gorm"
)

type PortfolioService struct {
	db *gorm.DB
}

// NewPortfolioService crea una nueva instancia de PortfolioService
func NewPortfolioService() *PortfolioService {
	return &PortfolioService{
		db: utils.GetDB(),
	}
}

// GetAllProjects obtiene todos los proyectos activos (publico)
func (ps *PortfolioService) GetAllProjects() ([]models.ProjectPortfolio, error) {
	var projects []models.ProjectPortfolio

	if err := ps.db.Where("is_active = ?", true).Order("order_index ASC, created_at DESC").Find(&projects).Error; err != nil {
		return nil, utils.ErrDatabaseError
	}

	return projects, nil
}

// GetAllProjectsAdmin obtiene todos los proyectos para admin
func (ps *PortfolioService) GetAllProjectsAdmin() ([]models.ProjectPortfolio, error) {
	var projects []models.ProjectPortfolio

	if err := ps.db.Order("order_index ASC, created_at DESC").Find(&projects).Error; err != nil {
		return nil, utils.ErrDatabaseError
	}

	return projects, nil
}

// GetProjectByID obtiene un proyecto por ID
func (ps *PortfolioService) GetProjectByID(id uint) (*models.ProjectPortfolio, error) {
	var project models.ProjectPortfolio

	if err := ps.db.First(&project, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrResourceNotFound
		}
		return nil, utils.ErrDatabaseError
	}

	return &project, nil
}

// GetProjectsByCategory obtiene proyectos por categoria
func (ps *PortfolioService) GetProjectsByCategory(category string) ([]models.ProjectPortfolio, error) {
	// Nota: Se permite cualquier categoría personalizada

	var projects []models.ProjectPortfolio

	if err := ps.db.Where("category = ? AND is_active = ?", category, true).
		Order("order_index ASC, created_at DESC").Find(&projects).Error; err != nil {
		return nil, utils.ErrDatabaseError
	}

	return projects, nil
}

// CreateProject crea un nuevo proyecto
func (ps *PortfolioService) CreateProject(req *models.ProjectRequest) (*models.ProjectPortfolio, error) {
	// Nota: Se permite cualquier categoría personalizada

	// Si no se especifica orden, obtener el siguiente disponible
	if req.OrderIndex == 0 {
		var maxOrder int
		ps.db.Model(&models.ProjectPortfolio{}).Select("COALESCE(MAX(order_index), 0)").Scan(&maxOrder)
		req.OrderIndex = maxOrder + 1
	}

	project := &models.ProjectPortfolio{
		Title:       req.Title,
		Category:    req.Category,
		Description: req.Description,
		OrderIndex:  req.OrderIndex,
		IsActive:    req.IsActive,
	}

	if err := ps.db.Create(project).Error; err != nil {
		log.Printf("Error creando proyecto: %v", err)
		return nil, utils.ErrDatabaseError
	}

	return project, nil
}

// UpdateProject actualiza un proyecto existente
func (ps *PortfolioService) UpdateProject(id uint, req *models.ProjectRequest) (*models.ProjectPortfolio, error) {
	project, err := ps.GetProjectByID(id)
	if err != nil {
		return nil, err
	}

	// Nota: Se permite cualquier categoría personalizada

	// Actualizar campos
	project.UpdateFromRequest(req)

	if err := ps.db.Save(project).Error; err != nil {
		log.Printf("Error actualizando proyecto: %v", err)
		return nil, utils.ErrDatabaseError
	}

	return project, nil
}

// DeleteProject elimina un proyecto
func (ps *PortfolioService) DeleteProject(id uint) error {
	project, err := ps.GetProjectByID(id)
	if err != nil {
		return err
	}

	// Eliminar imagen asociada si existe
	if project.ImageURL != "" {
		imagePath := utils.GetImagePath(project.ImageURL)
		if err := utils.DeleteFile(imagePath); err != nil {
			log.Printf("Error eliminando imagen %s: %v", imagePath, err)
		}
	}

	if err := ps.db.Delete(project).Error; err != nil {
		return utils.ErrDatabaseError
	}

	// Reordenar proyectos restantes
	ps.reorderProjectsAfterDeletion(project.OrderIndex)

	return nil
}

// UpdateProjectImage actualiza la imagen de un proyecto
func (ps *PortfolioService) UpdateProjectImage(id uint, imageURL string) (*models.ProjectPortfolio, error) {
	project, err := ps.GetProjectByID(id)
	if err != nil {
		return nil, err
	}

	// Eliminar imagen anterior si existe
	if project.ImageURL != "" && project.ImageURL != imageURL {
		oldImagePath := utils.GetImagePath(project.ImageURL)
		if err := utils.DeleteFile(oldImagePath); err != nil {
			log.Printf("Error eliminando imagen anterior %s: %v", oldImagePath, err)
		}
	}

	project.SetImageURL(imageURL)

	if err := ps.db.Save(project).Error; err != nil {
		return nil, utils.ErrDatabaseError
	}

	return project, nil
}

// ReorderProjects reordena los proyectos segun una lista de IDs
func (ps *PortfolioService) ReorderProjects(projectIDs []uint) error {
	// Validar que todos los IDs existan
	var existingCount int64
	ps.db.Model(&models.ProjectPortfolio{}).Where("id IN ?", projectIDs).Count(&existingCount)
	
	if int(existingCount) != len(projectIDs) {
		return utils.ErrResourceNotFound
	}

	// Actualizar orden en transaccion
	return ps.db.Transaction(func(tx *gorm.DB) error {
		for i, projectID := range projectIDs {
			if err := tx.Model(&models.ProjectPortfolio{}).
				Where("id = ?", projectID).
				Update("order_index", i+1).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetPortfolioStats obtiene estadisticas del portafolio
func (ps *PortfolioService) GetPortfolioStats() (*models.PortfolioStats, error) {
	var totalProjects, activeProjects, inactiveProjects int64

	// Contar proyectos totales
	ps.db.Model(&models.ProjectPortfolio{}).Count(&totalProjects)
	
	// Contar proyectos activos
	ps.db.Model(&models.ProjectPortfolio{}).Where("is_active = ?", true).Count(&activeProjects)
	
	// Contar proyectos inactivos
	inactiveProjects = totalProjects - activeProjects

	// Estadisticas por categoria
	var categoryStats []models.CategoryStats
	ps.db.Model(&models.ProjectPortfolio{}).
		Select("category, COUNT(*) as count").
		Where("is_active = ?", true).
		Group("category").
		Find(&categoryStats)

	// Proyectos recientes
	var recentProjects []models.ProjectPortfolio
	ps.db.Where("is_active = ?", true).
		Order("created_at DESC").
		Limit(5).
		Find(&recentProjects)

	// Convertir a respuestas
	var recentResponses []models.ProjectResponse
	for _, project := range recentProjects {
		recentResponses = append(recentResponses, project.ToResponse())
	}

	stats := &models.PortfolioStats{
		TotalProjects:    int(totalProjects),
		ActiveProjects:   int(activeProjects),
		InactiveProjects: int(inactiveProjects),
		CategoriesStats:  categoryStats,
		RecentProjects:   recentResponses,
	}

	return stats, nil
}

// ToggleProjectStatus cambia el estado activo/inactivo de un proyecto
func (ps *PortfolioService) ToggleProjectStatus(id uint) (*models.ProjectPortfolio, error) {
	project, err := ps.GetProjectByID(id)
	if err != nil {
		return nil, err
	}

	if project.IsActive {
		project.Deactivate()
	} else {
		project.Activate()
	}

	if err := ps.db.Save(project).Error; err != nil {
		return nil, utils.ErrDatabaseError
	}

	return project, nil
}

// reorderProjectsAfterDeletion reordena proyectos despues de eliminar uno
func (ps *PortfolioService) reorderProjectsAfterDeletion(deletedOrder int) {
	// Decrementar el orden de todos los proyectos que tenian un orden mayor
	ps.db.Model(&models.ProjectPortfolio{}).
		Where("order_index > ?", deletedOrder).
		UpdateColumn("order_index", gorm.Expr("order_index - 1"))
}

// GetUniqueCategories obtiene todas las categorias unicas
func (ps *PortfolioService) GetUniqueCategories() ([]string, error) {
	var categories []string

	if err := ps.db.Model(&models.ProjectPortfolio{}).
		Select("DISTINCT category").
		Where("is_active = ?", true).
		Pluck("category", &categories).Error; err != nil {
		return nil, utils.ErrDatabaseError
	}

	return categories, nil
}