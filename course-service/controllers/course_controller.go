// controllers/course_controller.go
package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"course-service/config"
	"course-service/models"
	"course-service/services"
	"course-service/utils"
)

type CourseController struct {
	db            *gorm.DB
	config        *config.Config
	courseService *services.CourseService
}

func NewCourseController(db *gorm.DB, cfg *config.Config) *CourseController {
	return &CourseController{
		db:            db,
		config:        cfg,
		courseService: services.NewCourseService(db, cfg),
	}
}

// GetCourses obtiene lista de cursos (migrado de getCursos)
func (cc *CourseController) GetCourses(c *gin.Context) {
	// Parámetros de filtrado
	categoria := c.Query("categoria")
	nivel := c.Query("nivel")
	gratuito := c.Query("gratuito")
	destacado := c.Query("destacado")
	search := c.Query("search")

	// Parámetros de paginación
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Obtener cursos con filtros
	courses, total, err := cc.courseService.GetCourses(page, limit, map[string]string{
		"categoria": categoria,
		"nivel":     nivel,
		"gratuito":  gratuito,
		"destacado": destacado,
		"search":    search,
	})

	if err != nil {
		log.Printf("Error al obtener cursos: %v", err)
		utils.SendErrorResponse(c, "Error al obtener cursos", http.StatusInternalServerError)
		return
	}

	// Si hay usuario autenticado, verificar acceso a cada curso
	userID, exists := c.Get("user_id")
	if exists {
		userIDUint, ok := userID.(uint)
		if ok {
			courses = cc.courseService.CheckUserAccessToCourses(courses, userIDUint)
		}
	}

	utils.SendSuccessResponse(c, gin.H{
		"courses": courses,
		"pagination": gin.H{
			"total":  total,
			"page":   page,
			"limit":  limit,
			"pages":  (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// GetCourseByID obtiene un curso específico (migrado de getCursoById)
func (cc *CourseController) GetCourseByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.SendErrorResponse(c, "ID de curso requerido", http.StatusBadRequest)
		return
	}

	courseID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.SendErrorResponse(c, "ID de curso inválido", http.StatusBadRequest)
		return
	}

	// Obtener curso con capítulos
	course, err := cc.courseService.GetCourseByID(uint(courseID), true)
	if err != nil {
		log.Printf("Curso no encontrado ID: %s, Error: %v", id, err)
		utils.SendErrorResponse(c, "Curso no encontrado", http.StatusNotFound)
		return
	}

	// Devolver curso completo con capítulos para el frontend
	utils.SendSuccessResponse(c, course)
}

// CreateCourse crea un nuevo curso (migrado de createCurso)
func (cc *CourseController) CreateCourse(c *gin.Context) {
	// Verificar autenticación
	userID, exists := c.Get("user_id")
	if !exists {
		utils.SendErrorResponse(c, "Usuario no autenticado", http.StatusUnauthorized)
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		utils.SendErrorResponse(c, "Error al obtener información del usuario", http.StatusInternalServerError)
		return
	}

	log.Printf("Usuario autenticado ID: %v creando curso", userIDUint)

	// Parsear datos del formulario multipart
	titulo := c.PostForm("titulo")
	descripcion := c.PostForm("descripcion")
	contenido := c.PostForm("contenido")
	precioStr := c.PostForm("precio")
	estado := c.PostForm("estado")
	nivel := c.PostForm("nivel")
	categoriaIDStr := c.PostForm("categoria_id")
	destacadoStr := c.PostForm("destacado")
	gratuitoStr := c.PostForm("gratuito")

	// Validar campos requeridos
	if titulo == "" || descripcion == "" || contenido == "" {
		utils.SendErrorResponse(c, "Campos requeridos incompletos", http.StatusBadRequest)
		return
	}

	// Convertir campos
	precio, _ := strconv.ParseFloat(precioStr, 64)
	categoriaID, _ := strconv.ParseUint(categoriaIDStr, 10, 32)
	destacado, _ := strconv.ParseBool(destacadoStr)
	gratuito, _ := strconv.ParseBool(gratuitoStr)

	// Valores por defecto
	if estado == "" {
		estado = models.CourseStatusDraft
	}
	if nivel == "" {
		nivel = models.CourseLevelBeginner
	}

	// Manejar imagen
	imagenURL, err := utils.SaveImageFile(c, "imagen", cc.config)
	if err != nil {
		log.Printf("Error al guardar imagen: %v", err)
		utils.SendErrorResponse(c, "Error con la imagen: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Si no hay nueva imagen, usar URL proporcionada
	if imagenURL == "" {
		imagenURL = c.PostForm("imagen_url")
	}

	// Crear curso
	courseReq := &models.CourseRequest{
		Titulo:      titulo,
		Descripcion: descripcion,
		Contenido:   contenido,
		Precio:      precio,
		Estado:      estado,
		ImagenURL:   imagenURL,
		Nivel:       nivel,
		CategoriaID: uint(categoriaID),
		Destacado:   destacado,
		Gratuito:    gratuito,
	}

	course, err := cc.courseService.CreateCourse(courseReq, userIDUint)
	if err != nil {
		log.Printf("Error al crear curso: %v", err)
		utils.SendErrorResponse(c, "Error al crear curso: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Curso creado exitosamente: ID %v", course.ID)
	utils.SendCreatedResponse(c, course)
}

// UpdateCourse actualiza un curso existente (migrado de updateCurso)
func (cc *CourseController) UpdateCourse(c *gin.Context) {
	id := c.Param("id")
	courseID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.SendErrorResponse(c, "ID de curso inválido", http.StatusBadRequest)
		return
	}

	// Verificar que el curso existe
	course, err := cc.courseService.GetCourseByID(uint(courseID), false)
	if err != nil {
		utils.SendErrorResponse(c, "Curso no encontrado", http.StatusNotFound)
		return
	}

	// Parsear datos del formulario
	titulo := c.PostForm("titulo")
	descripcion := c.PostForm("descripcion")
	contenido := c.PostForm("contenido")
	precioStr := c.PostForm("precio")
	estado := c.PostForm("estado")
	nivel := c.PostForm("nivel")
	categoriaIDStr := c.PostForm("categoria_id")
	destacadoStr := c.PostForm("destacado")
	gratuitoStr := c.PostForm("gratuito")

	// Validar campos requeridos
	if titulo == "" || descripcion == "" || contenido == "" {
		utils.SendErrorResponse(c, "Campos requeridos incompletos", http.StatusBadRequest)
		return
	}

	// Convertir campos
	precio, _ := strconv.ParseFloat(precioStr, 64)
	destacado, _ := strconv.ParseBool(destacadoStr)
	gratuito, _ := strconv.ParseBool(gratuitoStr)

	// Manejar imagen
	imagenURL, err := utils.SaveImageFile(c, "imagen", cc.config)
	if err != nil {
		log.Printf("Error al guardar imagen: %v", err)
		utils.SendErrorResponse(c, "Error con la imagen: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Si no hay nueva imagen, mantener la actual o usar URL proporcionada
	if imagenURL == "" {
		imagenURLPost := c.PostForm("imagen_url")
		if imagenURLPost != "" {
			imagenURL = imagenURLPost
		} else {
			imagenURL = course.ImagenURL
		}
	}

	// Manejar categoria_id: si no se proporciona o es vacío, mantener la categoría actual
	var finalCategoriaID uint
	if categoriaIDStr != "" {
		if parsedID, err := strconv.ParseUint(categoriaIDStr, 10, 32); err == nil {
			finalCategoriaID = uint(parsedID)
		} else {
			finalCategoriaID = course.CategoriaID // Mantener categoría actual si hay error de parseo
		}
	} else {
		finalCategoriaID = course.CategoriaID // Mantener categoría actual si no se envía
	}

	// Actualizar curso
	courseReq := &models.CourseRequest{
		Titulo:      titulo,
		Descripcion: descripcion,
		Contenido:   contenido,
		Precio:      precio,
		Estado:      estado,
		ImagenURL:   imagenURL,
		Nivel:       nivel,
		CategoriaID: finalCategoriaID,
		Destacado:   destacado,
		Gratuito:    gratuito,
	}

	updatedCourse, err := cc.courseService.UpdateCourse(uint(courseID), courseReq)
	if err != nil {
		log.Printf("Error al actualizar curso: %v", err)
		utils.SendErrorResponse(c, "Error al actualizar curso: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Curso actualizado exitosamente: ID %v", updatedCourse.ID)
	utils.SendSuccessResponse(c, updatedCourse)
}

// DeleteCourse elimina un curso (migrado de deleteCurso)
func (cc *CourseController) DeleteCourse(c *gin.Context) {
	id := c.Param("id")
	courseID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.SendErrorResponse(c, "ID de curso inválido", http.StatusBadRequest)
		return
	}

	err = cc.courseService.DeleteCourse(uint(courseID))
	if err != nil {
		log.Printf("Error al eliminar curso ID %s: %v", id, err)
		utils.SendErrorResponse(c, "Error al eliminar curso: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Curso ID %s eliminado exitosamente", id)
	utils.SendSuccessResponse(c, gin.H{"message": "Curso eliminado correctamente"})
}

// CheckCourseAccess verifica si el usuario tiene acceso al curso
func (cc *CourseController) CheckCourseAccess(c *gin.Context) {
	// Obtener usuario autenticado
	userID, exists := c.Get("user_id")
	if !exists {
		utils.SendErrorResponse(c, "Usuario no autenticado", http.StatusUnauthorized)
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		utils.SendErrorResponse(c, "Error al obtener información del usuario", http.StatusInternalServerError)
		return
	}

	// Obtener ID del curso
	id := c.Param("id")
	courseID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.SendErrorResponse(c, "ID de curso inválido", http.StatusBadRequest)
		return
	}

	// Verificar acceso
	hasAccess := cc.courseService.CheckUserAccessToCourse(uint(courseID), userIDUint)

	utils.SendSuccessResponse(c, gin.H{
		"course_id":   courseID,
		"user_id":     userIDUint,
		"has_access":  hasAccess,
	})
}

// GetCourseCategories obtiene todas las categorías
func (cc *CourseController) GetCourseCategories(c *gin.Context) {
	categories, err := cc.courseService.GetCategories()
	if err != nil {
		log.Printf("Error al obtener categorías: %v", err)
		utils.SendErrorResponse(c, "Error al obtener categorías", http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, gin.H{
		"categories": categories,
	})
}

// CreateCategory crea una nueva categoría (admin)
func (cc *CourseController) CreateCategory(c *gin.Context) {
	var categoryReq struct {
		Nombre      string `json:"nombre" binding:"required"`
		Descripcion string `json:"descripcion"`
		Slug        string `json:"slug"`
		PadreID     *uint  `json:"padre_id"`
		Activa      *bool  `json:"activa"`
		Orden       *int   `json:"orden"`
	}

	if err := c.ShouldBindJSON(&categoryReq); err != nil {
		utils.SendErrorResponse(c, "Datos inválidos: "+err.Error(), http.StatusBadRequest)
		return
	}

	category, err := cc.courseService.CreateCategory(&categoryReq)
	if err != nil {
		log.Printf("Error al crear categoría: %v", err)
		utils.SendErrorResponse(c, "Error al crear categoría: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendCreatedResponse(c, category)
}

// UpdateCategory actualiza una categoría existente (admin)
func (cc *CourseController) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	categoryID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.SendErrorResponse(c, "ID de categoría inválido", http.StatusBadRequest)
		return
	}

	var categoryReq struct {
		Nombre      string `json:"nombre"`
		Descripcion string `json:"descripcion"`
		Slug        string `json:"slug"`
		PadreID     *uint  `json:"padre_id"`
		Activa      *bool  `json:"activa"`
		Orden       *int   `json:"orden"`
	}

	if err := c.ShouldBindJSON(&categoryReq); err != nil {
		utils.SendErrorResponse(c, "Datos inválidos: "+err.Error(), http.StatusBadRequest)
		return
	}

	category, err := cc.courseService.UpdateCategory(uint(categoryID), &categoryReq)
	if err != nil {
		log.Printf("Error al actualizar categoría: %v", err)
		utils.SendErrorResponse(c, "Error al actualizar categoría: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, category)
}

// DeleteCategory elimina una categoría (admin)
func (cc *CourseController) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	categoryID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.SendErrorResponse(c, "ID de categoría inválido", http.StatusBadRequest)
		return
	}

	err = cc.courseService.DeleteCategory(uint(categoryID))
	if err != nil {
		log.Printf("Error al eliminar categoría: %v", err)
		utils.SendErrorResponse(c, "Error al eliminar categoría: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, gin.H{
		"message": "Categoría eliminada exitosamente",
	})
}

// HealthCheck endpoint para verificar la salud del servicio
func (cc *CourseController) HealthCheck(c *gin.Context) {
	// Verificar conexión a base de datos
	sqlDB, err := cc.db.DB()
	if err != nil {
		utils.SendErrorResponse(c, "Error de base de datos", http.StatusInternalServerError)
		return
	}

	if err := sqlDB.Ping(); err != nil {
		utils.SendErrorResponse(c, "Base de datos no disponible", http.StatusServiceUnavailable)
		return
	}

	// Verificar conectividad con Payment Service
	paymentStatus := "unknown"
	if err := services.CheckPaymentServiceHealth(cc.config); err == nil {
		paymentStatus = "connected"
	} else {
		paymentStatus = "disconnected"
	}

	utils.SendSuccessResponse(c, gin.H{
		"status":          "ok",
		"service":         "course-service",
		"version":         "1.0.0",
		"database":        "connected",
		"payment_service": paymentStatus,
	})
}

// GetAdminCourses obtiene todos los cursos para administradores
func (cc *CourseController) GetAdminCourses(c *gin.Context) {
	// Parámetros de paginación
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	// Obtener todos los cursos (incluidos inactivos/borrador)
	courses, total, err := cc.courseService.GetAllCoursesAdmin(page, limit)
	if err != nil {
		log.Printf("Error al obtener cursos para admin: %v", err)
		utils.SendErrorResponse(c, "Error al obtener cursos", http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, gin.H{
		"courses": courses,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}