// services/chapter_service.go
package services

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"

	"course-service/config"
	"course-service/models"
)

type ChapterService struct {
	db     *gorm.DB
	config *config.Config
}

func NewChapterService(db *gorm.DB, cfg *config.Config) *ChapterService {
	return &ChapterService{
		db:     db,
		config: cfg,
	}
}

// GetChaptersByCourse obtiene todos los capítulos de un curso
func (cs *ChapterService) GetChaptersByCourse(courseID uint) ([]models.Chapter, error) {
	var chapters []models.Chapter
	if err := cs.db.Where("curso_id = ?", courseID).Order("orden ASC").Find(&chapters).Error; err != nil {
		return nil, fmt.Errorf("error al obtener capítulos: %v", err)
	}
	return chapters, nil
}

// GetChapterByID obtiene un capítulo específico por ID
func (cs *ChapterService) GetChapterByID(chapterID uint) (*models.Chapter, error) {
	var chapter models.Chapter
	if err := cs.db.First(&chapter, chapterID).Error; err != nil {
		return nil, fmt.Errorf("capítulo no encontrado")
	}
	return &chapter, nil
}

// CreateChapter crea un nuevo capítulo
func (cs *ChapterService) CreateChapter(req *models.ChapterRequest) (*models.Chapter, error) {
	// Validaciones
	if err := cs.validateChapterRequest(req); err != nil {
		return nil, err
	}

	// Extraer nombre del archivo del video URL si existe
	videoNombre := ""
	if req.VideoURL != "" {
		parts := strings.Split(req.VideoURL, "/")
		if len(parts) > 0 {
			videoNombre = parts[len(parts)-1]
		}
	}

	// Establecer valores por defecto
	if req.TipoContenido == "" {
		req.TipoContenido = models.ContentTypeVideo
	}

	chapter := &models.Chapter{
		CursoID:       req.CursoID,
		Titulo:        req.Titulo,
		Descripcion:   req.Descripcion,
		Duracion:      req.Duracion,
		VideoURL:      req.VideoURL,
		VideoNombre:   videoNombre,
		Publicado:     req.Publicado,
		Orden:         req.Orden,
		TipoContenido: req.TipoContenido,
		Gratis:        req.Gratis,
	}

	if err := cs.db.Create(chapter).Error; err != nil {
		return nil, fmt.Errorf("error al crear capítulo: %v", err)
	}

	log.Printf("Capítulo creado exitosamente: ID %v", chapter.ID)
	return chapter, nil
}

// UpdateChapter actualiza un capítulo existente
func (cs *ChapterService) UpdateChapter(chapterID uint, req *models.ChapterRequest) (*models.Chapter, error) {
	// Validaciones
	if err := cs.validateChapterRequest(req); err != nil {
		return nil, err
	}

	// Obtener capítulo existente
	var chapter models.Chapter
	if err := cs.db.First(&chapter, chapterID).Error; err != nil {
		return nil, fmt.Errorf("capítulo no encontrado")
	}

	// Extraer nombre del archivo si cambió el video URL
	videoNombre := chapter.VideoNombre
	if req.VideoURL != chapter.VideoURL {
		videoNombre = ""
		if req.VideoURL != "" {
			parts := strings.Split(req.VideoURL, "/")
			if len(parts) > 0 {
				videoNombre = parts[len(parts)-1]
			}
		}
	}

	// Actualizar campos
	chapter.CursoID = req.CursoID
	chapter.Titulo = req.Titulo
	chapter.Descripcion = req.Descripcion
	chapter.Duracion = req.Duracion
	chapter.VideoURL = req.VideoURL
	chapter.VideoNombre = videoNombre
	chapter.Publicado = req.Publicado
	chapter.Orden = req.Orden
	chapter.TipoContenido = req.TipoContenido
	chapter.Gratis = req.Gratis

	if err := cs.db.Save(&chapter).Error; err != nil {
		return nil, fmt.Errorf("error al actualizar capítulo: %v", err)
	}

	log.Printf("Capítulo actualizado exitosamente: ID %v", chapter.ID)
	return &chapter, nil
}

// DeleteChapter elimina un capítulo y sus archivos asociados
func (cs *ChapterService) DeleteChapter(chapterID uint) error {
	// Obtener capítulo
	var chapter models.Chapter
	if err := cs.db.First(&chapter, chapterID).Error; err != nil {
		return fmt.Errorf("capítulo no encontrado")
	}

	// Eliminar archivo de video si existe
	if chapter.VideoNombre != "" && chapter.CursoID > 0 {
		fileService := NewFileService(cs.db, cs.config)
		courseIDStr := fmt.Sprintf("%d", chapter.CursoID)
		if err := fileService.DeleteFile(courseIDStr, chapter.VideoNombre); err != nil {
			log.Printf("Advertencia: No se pudo eliminar archivo de video: %v", err)
		}
	}

	// Eliminar progreso asociado al capítulo
	if err := cs.db.Where("capitulo_id = ?", chapterID).Delete(&models.ChapterProgress{}).Error; err != nil {
		log.Printf("Error al eliminar progreso del capítulo: %v", err)
	}

	// Eliminar capítulo
	if err := cs.db.Delete(&chapter).Error; err != nil {
		return fmt.Errorf("error al eliminar capítulo: %v", err)
	}

	log.Printf("Capítulo ID %d eliminado exitosamente", chapterID)
	return nil
}

// ReorderChapters reordena los capítulos de un curso
func (cs *ChapterService) ReorderChapters(courseID uint, chapterIDs []uint) error {
	// Validar que todos los IDs pertenecen al curso
	var count int64
	if err := cs.db.Model(&models.Chapter{}).
		Where("curso_id = ? AND id IN ?", courseID, chapterIDs).
		Count(&count).Error; err != nil {
		return fmt.Errorf("error al validar capítulos: %v", err)
	}

	if int(count) != len(chapterIDs) {
		return fmt.Errorf("algunos capítulos no pertenecen al curso o no existen")
	}

	// Iniciar transacción
	tx := cs.db.Begin()

	// Actualizar orden de cada capítulo
	for i, chapterID := range chapterIDs {
		if err := tx.Model(&models.Chapter{}).
			Where("id = ?", chapterID).
			Update("orden", i+1).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("error al actualizar orden del capítulo %d: %v", chapterID, err)
		}
	}

	// Confirmar transacción
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("error al confirmar reordenamiento: %v", err)
	}

	log.Printf("Capítulos reordenados correctamente para el curso %d", courseID)
	return nil
}

// UpdateChapterPublishedStatus actualiza solo el estado de publicación
func (cs *ChapterService) UpdateChapterPublishedStatus(chapterID uint, published bool) error {
	if err := cs.db.Model(&models.Chapter{}).
		Where("id = ?", chapterID).
		Update("publicado", published).Error; err != nil {
		return fmt.Errorf("error al actualizar estado de publicación: %v", err)
	}

	log.Printf("Estado de publicación del capítulo %d actualizado a: %v", chapterID, published)
	return nil
}

// PrepareChaptersResponse prepara la respuesta de capítulos con información de acceso
func (cs *ChapterService) PrepareChaptersResponse(chapters []models.Chapter, hasAccess bool) []models.ChapterResponse {
	var responses []models.ChapterResponse

	for _, chapter := range chapters {
		response := cs.PrepareChapterResponse(&chapter, hasAccess)
		responses = append(responses, response)
	}

	return responses
}

// PrepareChapterResponse prepara la respuesta de un capítulo específico
func (cs *ChapterService) PrepareChapterResponse(chapter *models.Chapter, hasAccess bool) models.ChapterResponse {
	// Determinar si el usuario puede acceder a este capítulo
	canAccess := chapter.CanUserAccess(hasAccess)

	response := models.ChapterResponse{
		ID:            chapter.ID,
		Titulo:        chapter.Titulo,
		Descripcion:   chapter.Descripcion,
		Duracion:      chapter.Duracion,
		Orden:         chapter.Orden,
		Publicado:     chapter.Publicado,
		TipoContenido: chapter.TipoContenido,
		Gratis:        chapter.Gratis,
		TieneAcceso:   canAccess,
		CreatedAt:     chapter.CreatedAt,
	}

	// Solo incluir URL del video si tiene acceso
	if canAccess {
		response.VideoURL = chapter.VideoURL
	}

	return response
}

// GetPublishedChaptersByCourse obtiene solo capítulos publicados
func (cs *ChapterService) GetPublishedChaptersByCourse(courseID uint) ([]models.Chapter, error) {
	var chapters []models.Chapter
	if err := cs.db.Where("curso_id = ? AND publicado = true", courseID).
		Order("orden ASC").Find(&chapters).Error; err != nil {
		return nil, fmt.Errorf("error al obtener capítulos publicados: %v", err)
	}
	return chapters, nil
}

// GetChapterStatistics obtiene estadísticas de capítulos para un curso
func (cs *ChapterService) GetChapterStatistics(courseID uint) (map[string]int, error) {
	stats := make(map[string]int)

	// Total de capítulos
	var total int64
	if err := cs.db.Model(&models.Chapter{}).Where("curso_id = ?", courseID).Count(&total).Error; err != nil {
		return nil, err
	}
	stats["total"] = int(total)

	// Capítulos publicados
	var published int64
	if err := cs.db.Model(&models.Chapter{}).
		Where("curso_id = ? AND publicado = true", courseID).Count(&published).Error; err != nil {
		return nil, err
	}
	stats["published"] = int(published)

	// Capítulos borradores
	stats["draft"] = int(total) - int(published)

	// Capítulos gratuitos
	var free int64
	if err := cs.db.Model(&models.Chapter{}).
		Where("curso_id = ? AND gratis = true", courseID).Count(&free).Error; err != nil {
		return nil, err
	}
	stats["free"] = int(free)

	return stats, nil
}

// validateChapterRequest valida los datos de un capítulo
func (cs *ChapterService) validateChapterRequest(req *models.ChapterRequest) error {
	if req.Titulo == "" {
		return fmt.Errorf("título es requerido")
	}

	if req.CursoID == 0 {
		return fmt.Errorf("ID de curso es requerido")
	}

	if req.Orden < 0 {
		return fmt.Errorf("orden no puede ser negativo")
	}

	return nil
}