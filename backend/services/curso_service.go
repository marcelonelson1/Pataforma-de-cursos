package services

import (
	"curso-platform/config"
	"curso-platform/models"
	"curso-platform/utils"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// CursoService maneja la lógica relacionada con los cursos
type CursoService struct{}

// NewCursoService crea una nueva instancia de CursoService
func NewCursoService() *CursoService {
	return &CursoService{}
}

// GetAll obtiene todos los cursos
func (s *CursoService) GetAll() ([]models.Curso, error) {
	var cursos []models.Curso

	// Utilizamos Preload para cargar también los capítulos de cada curso
	if err := config.DB.Preload("Capitulos").Find(&cursos).Error; err != nil {
		return nil, fmt.Errorf("error al obtener cursos: %v", err)
	}

	return cursos, nil
}

// GetByID obtiene un curso por su ID
func (s *CursoService) GetByID(id uint) (*models.Curso, error) {
	var curso models.Curso
	if err := config.DB.Preload("Capitulos").First(&curso, id).Error; err != nil {
		return nil, fmt.Errorf("curso no encontrado: %v", err)
	}

	return &curso, nil
}

// Create crea un nuevo curso
func (s *CursoService) Create(titulo, descripcion, contenido, estado string, precio float64, imagenURL string, imagenFile *os.File) (*models.Curso, error) {
	// Establecer valores por defecto si no se proporcionan
	if estado == "" {
		estado = "Borrador"
	}

	// Manejar la imagen si se proporcionó un archivo
	finalImagenURL := imagenURL
	if imagenFile != nil {
		// Crear directorio si no existe
		if err := utils.CreateDirIfNotExists("./static/images"); err != nil {
			return nil, err
		}

		// Generar nombre único para el archivo
		filename := utils.GenerateUniqueFileName(imagenFile.Name())
		destPath := filepath.Join("./static/images", filename)

		// Leer contenido del archivo original
		fileData, err := os.ReadFile(imagenFile.Name())
		if err != nil {
			return nil, fmt.Errorf("error al leer archivo imagen: %v", err)
		}

		// Escribir a nuevo archivo
		if err := os.WriteFile(destPath, fileData, 0644); err != nil {
			return nil, fmt.Errorf("error al guardar imagen: %v", err)
		}

		// Actualizar URL
		finalImagenURL = "/static/images/" + filename
	}

	curso := models.Curso{
		Titulo:      titulo,
		Descripcion: descripcion,
		Contenido:   contenido,
		Precio:      precio,
		Estado:      estado,
		ImagenURL:   finalImagenURL,
	}

	if err := config.DB.Create(&curso).Error; err != nil {
		return nil, fmt.Errorf("error al crear curso: %v", err)
	}

	// Buscar curso completo con capítulos
	var cursoCompleto models.Curso
	if err := config.DB.Preload("Capitulos").First(&cursoCompleto, curso.ID).Error; err != nil {
		// Si hay error, asegurar que capítulos no sea nil
		curso.Capitulos = []models.Capitulo{}
		return &curso, nil
	}

	// Asegurar que capítulos no sea nil para evitar errores en el frontend
	if cursoCompleto.Capitulos == nil {
		cursoCompleto.Capitulos = []models.Capitulo{}
	}

	return &cursoCompleto, nil
}

// Update actualiza un curso existente
func (s *CursoService) Update(id uint, titulo, descripcion, contenido, estado string, precio float64, imagenURL string, imagenFile *os.File) (*models.Curso, error) {
	var curso models.Curso
	if err := config.DB.First(&curso, id).Error; err != nil {
		return nil, fmt.Errorf("curso no encontrado: %v", err)
	}

	// Actualizar campos si se proporcionan
	if titulo != "" {
		curso.Titulo = titulo
	}
	if descripcion != "" {
		curso.Descripcion = descripcion
	}
	if contenido != "" {
		curso.Contenido = contenido
	}
	if estado != "" {
		curso.Estado = estado
	}
	if precio > 0 {
		curso.Precio = precio
	}

	// Manejar la imagen si se proporcionó un archivo
	if imagenFile != nil {
		// Crear directorio si no existe
		if err := utils.CreateDirIfNotExists("./static/images"); err != nil {
			return nil, err
		}

		// Si hay una imagen anterior, intentar eliminarla
		if curso.ImagenURL != "" && strings.HasPrefix(curso.ImagenURL, "/static/images/") {
			oldImagePath := "." + curso.ImagenURL
			if _, err := os.Stat(oldImagePath); err == nil {
				os.Remove(oldImagePath)
			}
		}

		// Generar nombre único para el archivo
		filename := utils.GenerateUniqueFileName(imagenFile.Name())
		destPath := filepath.Join("./static/images", filename)

		// Leer contenido del archivo original
		fileData, err := os.ReadFile(imagenFile.Name())
		if err != nil {
			return nil, fmt.Errorf("error al leer archivo imagen: %v", err)
		}

		// Escribir a nuevo archivo
		if err := os.WriteFile(destPath, fileData, 0644); err != nil {
			return nil, fmt.Errorf("error al guardar imagen: %v", err)
		}

		// Actualizar URL
		curso.ImagenURL = "/static/images/" + filename
	} else if imagenURL != "" {
		// Si se proporcionó solo URL
		curso.ImagenURL = imagenURL
	}

	if err := config.DB.Save(&curso).Error; err != nil {
		return nil, fmt.Errorf("error al actualizar curso: %v", err)
	}

	// Buscar curso completo con capítulos
	var cursoCompleto models.Curso
	if err := config.DB.Preload("Capitulos").First(&cursoCompleto, curso.ID).Error; err != nil {
		// Si hay error, asegurar que capítulos no sea nil
		curso.Capitulos = []models.Capitulo{}
		return &curso, nil
	}

	// Asegurar que capítulos no sea nil para evitar errores en el frontend
	if cursoCompleto.Capitulos == nil {
		cursoCompleto.Capitulos = []models.Capitulo{}
	}

	return &cursoCompleto, nil
}

// Delete elimina un curso y sus recursos asociados
func (s *CursoService) Delete(id uint) error {
	var curso models.Curso
	if err := config.DB.First(&curso, id).Error; err != nil {
		return fmt.Errorf("curso no encontrado: %v", err)
	}

	// Cargar los capítulos para eliminar archivos de video
	var capitulos []models.Capitulo
	if err := config.DB.Where("curso_id = ?", id).Find(&capitulos).Error; err == nil {
		for _, capitulo := range capitulos {
			if capitulo.VideoNombre != "" {
				// Construir la ruta del archivo
				cursoIDStr := fmt.Sprintf("%d", curso.ID)
				filePath := filepath.Join("./static/videos", cursoIDStr, capitulo.VideoNombre)

				// Intentar eliminar el archivo
				if _, err := os.Stat(filePath); !os.IsNotExist(err) {
					if err := os.Remove(filePath); err != nil {
						log.Printf("Error al eliminar archivo de video del capítulo %d: %v", capitulo.ID, err)
					} else {
						log.Printf("Archivo de video del capítulo %d eliminado: %s", capitulo.ID, filePath)
					}
				}
			}
		}
	}

	// Eliminar progreso asociado al curso
	if err := config.DB.Where("curso_id = ?", id).Delete(&models.ProgresoUsuario{}).Error; err != nil {
		log.Printf("Error al eliminar progreso de usuarios del curso: %v", err)
	}

	if err := config.DB.Where("curso_id = ?", id).Delete(&models.ProgresoCapitulo{}).Error; err != nil {
		log.Printf("Error al eliminar progreso de capítulos del curso: %v", err)
	}

	// Eliminar todos los capítulos relacionados
	if err := config.DB.Where("curso_id = ?", id).Delete(&models.Capitulo{}).Error; err != nil {
		return fmt.Errorf("error al eliminar capítulos del curso: %v", err)
	}

	// Eliminar imagen del curso si existe
	if curso.ImagenURL != "" && strings.Contains(curso.ImagenURL, "/static/images/") {
		// Extraer el nombre del archivo de la URL
		imagenNombre := filepath.Base(curso.ImagenURL)
		imagenPath := filepath.Join("./static/images", imagenNombre)

		if _, err := os.Stat(imagenPath); !os.IsNotExist(err) {
			if err := os.Remove(imagenPath); err != nil {
				log.Printf("Error al eliminar imagen del curso: %v", err)
			} else {
				log.Printf("Imagen del curso eliminada: %s", imagenPath)
			}
		}
	}

	// Eliminar el curso
	if err := config.DB.Delete(&curso).Error; err != nil {
		return fmt.Errorf("error al eliminar curso: %v", err)
	}

	// Intentar eliminar el directorio de videos del curso si existe
	cursoIDStr := fmt.Sprintf("%d", curso.ID)
	dirPath := filepath.Join("./static/videos", cursoIDStr)
	if _, err := os.Stat(dirPath); !os.IsNotExist(err) {
		if err := os.RemoveAll(dirPath); err != nil {
			log.Printf("Error al eliminar directorio de videos del curso: %v", err)
		} else {
			log.Printf("Directorio de videos del curso eliminado: %s", dirPath)
		}
	}

	return nil
}

// UploadVideo sube un archivo de video para un capítulo
func (s *CursoService) UploadVideo(cursoID, capituloID uint, file *os.File, fileSize int64, fileName string) (*string, error) {
	// Validar tamaño del archivo
	const MAX_UPLOAD_SIZE = 100 * 1024 * 1024 // 100 MB
	if fileSize > MAX_UPLOAD_SIZE {
		return nil, fmt.Errorf("el archivo no debe superar los 100 MB")
	}

	// Validar tipo de archivo
	fileExt := strings.ToLower(filepath.Ext(fileName))
	if fileExt != ".mp4" && fileExt != ".webm" && fileExt != ".ogg" {
		return nil, fmt.Errorf("solo se permiten archivos MP4, WebM y OGG")
	}

	// Verificar si el curso existe
	var curso models.Curso
	if err := config.DB.First(&curso, cursoID).Error; err != nil {
		return nil, fmt.Errorf("curso no encontrado")
	}

	// Generar nombre de archivo único
	uniqueID := uuid.New().String()
	filename := fmt.Sprintf("%d-%d-%s%s",
		cursoID,
		capituloID,
		uniqueID,
		fileExt)

	// Crear directorio específico para el curso si no existe
	cursoDir := filepath.Join("./static/videos", fmt.Sprintf("%d", cursoID))
	if err := utils.CreateDirIfNotExists(cursoDir); err != nil {
		return nil, err
	}

	// Ruta completa del archivo
	filePath := filepath.Join(cursoDir, filename)

	// Leer contenido del archivo original
	fileData, err := os.ReadFile(file.Name())
	if err != nil {
		return nil, fmt.Errorf("error al leer archivo: %v", err)
	}

	// Escribir a nuevo archivo
	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		return nil, fmt.Errorf("error al guardar video: %v", err)
	}

	// URL donde se puede acceder al video
	videoURL := fmt.Sprintf("/static/videos/%d/%s", cursoID, filename)

	return &videoURL, nil
}

// DeleteVideo elimina un archivo de video
func (s *CursoService) DeleteVideo(cursoID uint, filename string) error {
	if filename == "" {
		return errors.New("nombre de archivo requerido")
	}

	filePath := filepath.Join("./static/videos", fmt.Sprintf("%d", cursoID), filename)

	// Verificar si el archivo existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("video no encontrado: %s", filePath)
	}

	// Eliminar el archivo
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("error al eliminar el video: %v", err)
	}

	return nil
}