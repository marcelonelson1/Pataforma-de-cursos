package services

import (
	"curso-platform/config"
	"curso-platform/models"
	
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// CapituloService maneja la lógica relacionada con los capítulos
type CapituloService struct{}

// NewCapituloService crea una nueva instancia de CapituloService
func NewCapituloService() *CapituloService {
	return &CapituloService{}
}

// GetByCursoID obtiene todos los capítulos de un curso
func (s *CapituloService) GetByCursoID(cursoID uint) ([]models.Capitulo, error) {
	var capitulos []models.Capitulo
	if err := config.DB.Where("curso_id = ?", cursoID).Order("orden ASC").Find(&capitulos).Error; err != nil {
		return nil, fmt.Errorf("error al obtener capítulos: %v", err)
	}

	return capitulos, nil
}

// Create crea un nuevo capítulo
func (s *CapituloService) Create(req models.Capitulo) (*models.Capitulo, error) {
	// Validamos que el curso exista
	var curso models.Curso
	if err := config.DB.First(&curso, req.CursoID).Error; err != nil {
		return nil, fmt.Errorf("curso no encontrado")
	}

	// Extraer el nombre del archivo desde la URL si existe
	videoNombre := ""
	if req.VideoURL != "" {
		parts := strings.Split(req.VideoURL, "/")
		if len(parts) > 0 {
			videoNombre = parts[len(parts)-1]
		}
	}

	capitulo := models.Capitulo{
		CursoID:     req.CursoID,
		Titulo:      req.Titulo,
		Descripcion: req.Descripcion,
		Duracion:    req.Duracion,
		VideoURL:    req.VideoURL,
		VideoNombre: videoNombre,
		Publicado:   req.Publicado,
		Orden:       req.Orden,
	}

	if err := config.DB.Create(&capitulo).Error; err != nil {
		return nil, fmt.Errorf("error al crear capítulo: %v", err)
	}

	return &capitulo, nil
}

// Update actualiza un capítulo existente
func (s *CapituloService) Update(id uint, req models.Capitulo) (*models.Capitulo, error) {
	var capitulo models.Capitulo
	if err := config.DB.First(&capitulo, id).Error; err != nil {
		return nil, fmt.Errorf("capítulo no encontrado")
	}

	// Verificamos que el curso exista (seguridad adicional)
	var curso models.Curso
	if err := config.DB.First(&curso, req.CursoID).Error; err != nil {
		return nil, fmt.Errorf("curso no encontrado")
	}

	// Extraer el nombre del archivo desde la URL si cambió
	videoNombre := capitulo.VideoNombre
	if req.VideoURL != capitulo.VideoURL {
		videoNombre = ""
		if req.VideoURL != "" {
			parts := strings.Split(req.VideoURL, "/")
			if len(parts) > 0 {
				videoNombre = parts[len(parts)-1]
			}
		}
	}

	capitulo.CursoID = req.CursoID
	capitulo.Titulo = req.Titulo
	capitulo.Descripcion = req.Descripcion
	capitulo.Duracion = req.Duracion
	capitulo.VideoURL = req.VideoURL
	capitulo.VideoNombre = videoNombre
	capitulo.Publicado = req.Publicado
	capitulo.Orden = req.Orden

	if err := config.DB.Save(&capitulo).Error; err != nil {
		return nil, fmt.Errorf("error al actualizar capítulo: %v", err)
	}

	return &capitulo, nil
}

// Delete elimina un capítulo
func (s *CapituloService) Delete(id uint) error {
	var capitulo models.Capitulo
	if err := config.DB.First(&capitulo, id).Error; err != nil {
		return fmt.Errorf("capítulo no encontrado")
	}

	// Si hay un video asociado, eliminarlo
	if capitulo.VideoNombre != "" && capitulo.CursoID > 0 {
		cursoID := fmt.Sprintf("%d", capitulo.CursoID)
		filePath := filepath.Join("./static/videos", cursoID, capitulo.VideoNombre)

		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			log.Printf("Eliminando archivo de video: %s", filePath)
			if err := os.Remove(filePath); err != nil {
				log.Printf("Advertencia: No se pudo eliminar el archivo de video: %v", err)
			}
		}
	}

	// Eliminar progreso asociado al capítulo
	if err := config.DB.Where("capitulo_id = ?", id).Delete(&models.ProgresoCapitulo{}).Error; err != nil {
		log.Printf("Error al eliminar progreso del capítulo: %v", err)
	}

	if err := config.DB.Delete(&capitulo).Error; err != nil {
		return fmt.Errorf("error al eliminar capítulo: %v", err)
	}

	return nil
}

// GetCapitulo obtiene un capítulo por su ID
func (s *CapituloService) GetCapitulo(id uint) (*models.Capitulo, error) {
	var capitulo models.Capitulo
	if err := config.DB.First(&capitulo, id).Error; err != nil {
		return nil, fmt.Errorf("capítulo no encontrado")
	}
	return &capitulo, nil
}