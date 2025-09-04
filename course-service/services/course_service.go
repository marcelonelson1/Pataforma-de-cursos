// services/course_service.go
package services

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"

	"course-service/config"
	"course-service/models"
)

type CourseService struct {
	db     *gorm.DB
	config *config.Config
}

func NewCourseService(db *gorm.DB, cfg *config.Config) *CourseService {
	return &CourseService{
		db:     db,
		config: cfg,
	}
}

// GetCourses obtiene cursos con filtros y paginación
func (cs *CourseService) GetCourses(page, limit int, filters map[string]string) ([]models.CourseResponse, int64, error) {
	query := cs.db.Model(&models.Course{})

	// Aplicar filtros
	if categoria := filters["categoria"]; categoria != "" {
		query = query.Where("categoria_id = ?", categoria)
	}

	if nivel := filters["nivel"]; nivel != "" {
		query = query.Where("nivel = ?", nivel)
	}

	if gratuito := filters["gratuito"]; gratuito == "true" {
		query = query.Where("gratuito = true OR precio = 0")
	}

	if destacado := filters["destacado"]; destacado == "true" {
		query = query.Where("destacado = true")
	}

	if search := filters["search"]; search != "" {
		query = query.Where("titulo LIKE ? OR descripcion LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Solo cursos publicados para consultas públicas
	query = query.Where("estado = ?", models.CourseStatusPublished)

	// Contar total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Obtener cursos con paginación
	var courses []models.Course
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&courses).Error; err != nil {
		return nil, 0, err
	}

	// Convertir a respuesta
	var responses []models.CourseResponse
	for _, course := range courses {
		responses = append(responses, cs.prepareCourseListResponse(&course))
	}

	return responses, total, nil
}

// GetCourseByID obtiene un curso por ID
func (cs *CourseService) GetCourseByID(courseID uint, includeChapters bool) (*models.Course, error) {
	query := cs.db

	if includeChapters {
		query = query.Preload("Capitulos", func(db *gorm.DB) *gorm.DB {
			return db.Order("orden ASC")
		})
	}

	var course models.Course
	if err := query.First(&course, courseID).Error; err != nil {
		return nil, err
	}

	return &course, nil
}

// CreateCourse crea un nuevo curso
func (cs *CourseService) CreateCourse(req *models.CourseRequest, instructorID uint) (*models.Course, error) {
	// Validaciones
	if err := cs.validateCourseRequest(req); err != nil {
		return nil, err
	}

	course := &models.Course{
		Titulo:       req.Titulo,
		Descripcion:  req.Descripcion,
		Contenido:    req.Contenido,
		Precio:       req.Precio,
		Estado:       req.Estado,
		ImagenURL:    req.ImagenURL,
		InstructorID: instructorID,
		Nivel:        req.Nivel,
		CategoriaID:  req.CategoriaID,
		Destacado:    req.Destacado,
		Gratuito:     req.Gratuito,
		Capitulos:    []models.Chapter{}, // Inicializar slice vacío
	}

	if err := cs.db.Create(course).Error; err != nil {
		return nil, fmt.Errorf("error al crear curso: %v", err)
	}

	// Recargar con capítulos
	if err := cs.db.Preload("Capitulos").First(course, course.ID).Error; err != nil {
		log.Printf("Advertencia: No se pudieron cargar capítulos del curso creado: %v", err)
	}

	log.Printf("Curso creado exitosamente: ID %v", course.ID)
	return course, nil
}

// UpdateCourse actualiza un curso existente
func (cs *CourseService) UpdateCourse(courseID uint, req *models.CourseRequest) (*models.Course, error) {
	// Validaciones
	if err := cs.validateCourseRequest(req); err != nil {
		return nil, err
	}

	// Obtener curso existente
	var course models.Course
	if err := cs.db.First(&course, courseID).Error; err != nil {
		return nil, fmt.Errorf("curso no encontrado")
	}

	// Actualizar campos
	course.Titulo = req.Titulo
	course.Descripcion = req.Descripcion
	course.Contenido = req.Contenido
	course.Precio = req.Precio
	course.Estado = req.Estado
	course.ImagenURL = req.ImagenURL
	course.Nivel = req.Nivel
	course.CategoriaID = req.CategoriaID
	course.Destacado = req.Destacado
	course.Gratuito = req.Gratuito

	if err := cs.db.Save(&course).Error; err != nil {
		return nil, fmt.Errorf("error al actualizar curso: %v", err)
	}

	// Recargar con capítulos
	if err := cs.db.Preload("Capitulos", func(db *gorm.DB) *gorm.DB {
		return db.Order("orden ASC")
	}).First(&course, course.ID).Error; err != nil {
		log.Printf("Advertencia: No se pudieron cargar capítulos del curso actualizado: %v", err)
	}

	log.Printf("Curso actualizado exitosamente: ID %v", course.ID)
	return &course, nil
}

// DeleteCourse elimina un curso y todos sus recursos asociados
func (cs *CourseService) DeleteCourse(courseID uint) error {
	// Obtener curso con capítulos
	var course models.Course
	if err := cs.db.Preload("Capitulos").First(&course, courseID).Error; err != nil {
		return fmt.Errorf("curso no encontrado")
	}

	// Eliminar archivos de capítulos
	fileService := NewFileService(cs.db, cs.config)
	for _, chapter := range course.Capitulos {
		if chapter.VideoNombre != "" {
			courseIDStr := fmt.Sprintf("%d", course.ID)
			if err := fileService.DeleteFile(courseIDStr, chapter.VideoNombre); err != nil {
				log.Printf("Error al eliminar archivo del capítulo %d: %v", chapter.ID, err)
			}
		}
	}

	// Eliminar progreso asociado
	if err := cs.db.Where("curso_id = ?", courseID).Delete(&models.UserProgress{}).Error; err != nil {
		log.Printf("Error al eliminar progreso de usuarios: %v", err)
	}

	if err := cs.db.Where("curso_id = ?", courseID).Delete(&models.ChapterProgress{}).Error; err != nil {
		log.Printf("Error al eliminar progreso de capítulos: %v", err)
	}

	// Eliminar capítulos
	if err := cs.db.Where("curso_id = ?", courseID).Delete(&models.Chapter{}).Error; err != nil {
		return fmt.Errorf("error al eliminar capítulos: %v", err)
	}

	// Eliminar imagen del curso
	if course.ImagenURL != "" && strings.Contains(course.ImagenURL, "/static/images/") {
		if err := fileService.DeleteImageFile(course.ImagenURL); err != nil {
			log.Printf("Error al eliminar imagen del curso: %v", err)
		}
	}

	// Eliminar curso
	if err := cs.db.Delete(&course).Error; err != nil {
		return fmt.Errorf("error al eliminar curso: %v", err)
	}

	// Eliminar directorio de archivos del curso
	courseIDStr := fmt.Sprintf("%d", course.ID)
	if err := fileService.DeleteCourseDirectory(courseIDStr); err != nil {
		log.Printf("Error al eliminar directorio del curso: %v", err)
	}

	log.Printf("Curso ID %d eliminado exitosamente", courseID)
	return nil
}

// CheckUserAccessToCourse verifica si un usuario tiene acceso a un curso
func (cs *CourseService) CheckUserAccessToCourse(courseID, userID uint) bool {
	// Obtener información del curso
	var course models.Course
	if err := cs.db.First(&course, courseID).Error; err != nil {
		log.Printf("Error al obtener curso: %v", err)
		return false
	}

	// Si el curso es gratuito, el usuario tiene acceso
	if course.IsFree() {
		return true
	}

	// Si el curso es de pago, verificar con Payment Service
	return CheckUserPaymentForCourse(userID, courseID, cs.config)
}

// CheckUserAccessToCourses verifica acceso para múltiples cursos
func (cs *CourseService) CheckUserAccessToCourses(courses []models.CourseResponse, userID uint) []models.CourseResponse {
	for i := range courses {
		courses[i].TieneAcceso = cs.CheckUserAccessToCourse(uint(courses[i].ID), userID)
	}
	return courses
}

// PrepareCourseResponse prepara la respuesta de un curso con información de acceso
func (cs *CourseService) PrepareCourseResponse(course *models.Course, hasAccess bool) models.CourseResponse {
	// Contar capítulos publicados
	publishedChapters := 0
	for _, chapter := range course.Capitulos {
		if chapter.Publicado {
			publishedChapters++
		}
	}

	return models.CourseResponse{
		ID:            course.ID,
		Titulo:        course.Titulo,
		Descripcion:   course.Descripcion,
		Precio:        course.Precio,
		Estado:        course.Estado,
		ImagenURL:     course.ImagenURL,
		Nivel:         course.Nivel,
		Gratuito:      course.IsFree(),
		TieneAcceso:   hasAccess,
		TotalChapters: publishedChapters,
		CreatedAt:     course.CreatedAt,
	}
}

// GetCategories obtiene todas las categorías activas
func (cs *CourseService) GetCategories() ([]models.Category, error) {
	var categories []models.Category
	if err := cs.db.Where("activa = true").Order("orden ASC, nombre ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// validateCourseRequest valida los datos de un curso
func (cs *CourseService) validateCourseRequest(req *models.CourseRequest) error {
	if req.Titulo == "" {
		return fmt.Errorf("título es requerido")
	}

	if req.Descripcion == "" {
		return fmt.Errorf("descripción es requerida")
	}

	if req.Contenido == "" {
		return fmt.Errorf("contenido es requerido")
	}

	if req.Estado != "" && !models.IsValidCourseState(req.Estado) {
		return fmt.Errorf("estado de curso inválido")
	}

	if req.Nivel != "" && !models.IsValidCourseLevel(req.Nivel) {
		return fmt.Errorf("nivel de curso inválido")
	}

	if req.Precio < 0 {
		return fmt.Errorf("precio no puede ser negativo")
	}

	return nil
}

// prepareCourseListResponse prepara respuesta para lista de cursos
func (cs *CourseService) prepareCourseListResponse(course *models.Course) models.CourseResponse {
	return models.CourseResponse{
		ID:            course.ID,
		Titulo:        course.Titulo,
		Descripcion:   course.Descripcion,
		Precio:        course.Precio,
		Estado:        course.Estado,
		ImagenURL:     course.ImagenURL,
		Nivel:         course.Nivel,
		Gratuito:      course.IsFree(),
		TieneAcceso:   false, // Se actualiza después si hay usuario
		TotalChapters: len(course.Capitulos),
		CreatedAt:     course.CreatedAt,
	}
}

// GetAllCoursesAdmin obtiene todos los cursos para administradores (incluye borrador e inactivos)
func (cs *CourseService) GetAllCoursesAdmin(page, limit int) ([]models.Course, int64, error) {
	var courses []models.Course
	var total int64

	// Contar total sin filtros de estado (admin ve todo)
	if err := cs.db.Model(&models.Course{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("error al contar cursos: %v", err)
	}

	// Calcular offset
	offset := (page - 1) * limit

	// Obtener cursos con paginación, precargar capítulos y categoría
	query := cs.db.Preload("Capitulos").Preload("Categoria").
		Offset(offset).Limit(limit).
		Order("created_at DESC")

	if err := query.Find(&courses).Error; err != nil {
		return nil, 0, fmt.Errorf("error al obtener cursos: %v", err)
	}

	// Devolver cursos completos con capítulos para admin
	return courses, total, nil
}

// prepareCourseAdminResponse prepara respuesta detallada para administradores
func (cs *CourseService) prepareCourseAdminResponse(course *models.Course) models.CourseResponse {
	return models.CourseResponse{
		ID:            course.ID,
		Titulo:        course.Titulo,
		Descripcion:   course.Descripcion,
		Precio:        course.Precio,
		Estado:        course.Estado,
		ImagenURL:     course.ImagenURL,
		Nivel:         course.Nivel,
		Gratuito:      course.IsFree(),
		TieneAcceso:   true, // Admin siempre tiene acceso
		TotalChapters: len(course.Capitulos),
		CreatedAt:     course.CreatedAt,
	}
}

// ================ MÉTODOS DE CATEGORÍAS ================

// CreateCategory crea una nueva categoría
func (cs *CourseService) CreateCategory(categoryReq interface{}) (*models.Category, error) {
	// Convertir interface{} a struct específico
	req := categoryReq.(*struct {
		Nombre      string `json:"nombre" binding:"required"`
		Descripcion string `json:"descripcion"`
		Slug        string `json:"slug"`
		PadreID     *uint  `json:"padre_id"`
		Activa      *bool  `json:"activa"`
		Orden       *int   `json:"orden"`
	})

	// Generar slug si no se proporciona
	slug := req.Slug
	if slug == "" {
		slug = generateSlug(req.Nombre)
	}

	// Valores por defecto
	activa := true
	if req.Activa != nil {
		activa = *req.Activa
	}

	orden := 0
	if req.Orden != nil {
		orden = *req.Orden
	}

	category := &models.Category{
		Nombre:      req.Nombre,
		Descripcion: req.Descripcion,
		Slug:        slug,
		PadreID:     req.PadreID,
		Activa:      activa,
		Orden:       orden,
	}

	// Verificar que el slug sea único
	var existingCategory models.Category
	if err := cs.db.Where("slug = ?", slug).First(&existingCategory).Error; err == nil {
		return nil, fmt.Errorf("ya existe una categoría con el slug '%s'", slug)
	}

	if err := cs.db.Create(category).Error; err != nil {
		return nil, fmt.Errorf("error al crear categoría: %v", err)
	}

	return category, nil
}

// UpdateCategory actualiza una categoría existente
func (cs *CourseService) UpdateCategory(categoryID uint, categoryReq interface{}) (*models.Category, error) {
	// Convertir interface{} a struct específico
	req := categoryReq.(*struct {
		Nombre      string `json:"nombre"`
		Descripcion string `json:"descripcion"`
		Slug        string `json:"slug"`
		PadreID     *uint  `json:"padre_id"`
		Activa      *bool  `json:"activa"`
		Orden       *int   `json:"orden"`
	})

	var category models.Category
	if err := cs.db.First(&category, categoryID).Error; err != nil {
		return nil, fmt.Errorf("categoría no encontrada")
	}

	// Actualizar campos si se proporcionan
	if req.Nombre != "" {
		category.Nombre = req.Nombre
	}
	if req.Descripcion != "" {
		category.Descripcion = req.Descripcion
	}
	if req.Slug != "" {
		// Verificar que el slug sea único (excepto para esta categoría)
		var existingCategory models.Category
		if err := cs.db.Where("slug = ? AND id != ?", req.Slug, categoryID).First(&existingCategory).Error; err == nil {
			return nil, fmt.Errorf("ya existe una categoría con el slug '%s'", req.Slug)
		}
		category.Slug = req.Slug
	}
	if req.PadreID != nil {
		category.PadreID = req.PadreID
	}
	if req.Activa != nil {
		category.Activa = *req.Activa
	}
	if req.Orden != nil {
		category.Orden = *req.Orden
	}

	if err := cs.db.Save(&category).Error; err != nil {
		return nil, fmt.Errorf("error al actualizar categoría: %v", err)
	}

	return &category, nil
}

// DeleteCategory elimina una categoría
func (cs *CourseService) DeleteCategory(categoryID uint) error {
	// Verificar que no tenga cursos asociados
	var courseCount int64
	if err := cs.db.Model(&models.Course{}).Where("categoria_id = ?", categoryID).Count(&courseCount).Error; err != nil {
		return fmt.Errorf("error al verificar cursos asociados: %v", err)
	}

	if courseCount > 0 {
		return fmt.Errorf("no se puede eliminar la categoría porque tiene %d cursos asociados", courseCount)
	}

	// Verificar que no tenga subcategorías
	var subcategoryCount int64
	if err := cs.db.Model(&models.Category{}).Where("padre_id = ?", categoryID).Count(&subcategoryCount).Error; err != nil {
		return fmt.Errorf("error al verificar subcategorías: %v", err)
	}

	if subcategoryCount > 0 {
		return fmt.Errorf("no se puede eliminar la categoría porque tiene %d subcategorías", subcategoryCount)
	}

	if err := cs.db.Delete(&models.Category{}, categoryID).Error; err != nil {
		return fmt.Errorf("error al eliminar categoría: %v", err)
	}

	return nil
}

// generateSlug genera un slug a partir del nombre
func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "á", "a")
	slug = strings.ReplaceAll(slug, "é", "e")
	slug = strings.ReplaceAll(slug, "í", "i")
	slug = strings.ReplaceAll(slug, "ó", "o")
	slug = strings.ReplaceAll(slug, "ú", "u")
	slug = strings.ReplaceAll(slug, "ñ", "n")
	// Remover caracteres especiales
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	return result.String()
}