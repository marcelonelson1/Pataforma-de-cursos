// controllers/chapter_controller.go
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

type ChapterController struct {
	db             *gorm.DB
	config         *config.Config
	chapterService *services.ChapterService
	courseService  *services.CourseService
}

func NewChapterController(db *gorm.DB, cfg *config.Config) *ChapterController {
	return &ChapterController{
		db:             db,
		config:         cfg,
		chapterService: services.NewChapterService(db, cfg),
		courseService:  services.NewCourseService(db, cfg),
	}
}

// GetChaptersByCourse obtiene capítulos de un curso (migrado de getCapitulosByCurso)
func (chc *ChapterController) GetChaptersByCourse(c *gin.Context) {
	cursoID := c.Param("id")
	if cursoID == "" {
		utils.SendErrorResponse(c, "ID de curso requerido", http.StatusBadRequest)
		return
	}

	courseID, err := strconv.ParseUint(cursoID, 10, 32)
	if err != nil {
		utils.SendErrorResponse(c, "ID de curso inválido", http.StatusBadRequest)
		return
	}

	// Verificar que el curso existe
	_, err = chc.courseService.GetCourseByID(uint(courseID), false)
	if err != nil {
		utils.SendErrorResponse(c, "Curso no encontrado", http.StatusNotFound)
		return
	}

	// Obtener capítulos
	chapters, err := chc.chapterService.GetChaptersByCourse(uint(courseID))
	if err != nil {
		log.Printf("Error al obtener capítulos del curso ID %s: %v", cursoID, err)
		utils.SendErrorResponse(c, "Error al obtener capítulos", http.StatusInternalServerError)
		return
	}

	// Si hay usuario autenticado, verificar acceso a cada capítulo
	userID, exists := c.Get("user_id")
	var hasAccess bool
	if exists {
		userIDUint, ok := userID.(uint)
		if ok {
			hasAccess = chc.courseService.CheckUserAccessToCourse(uint(courseID), userIDUint)
		}
	}

	// Preparar respuesta con acceso verificado
	chaptersResponse := chc.chapterService.PrepareChaptersResponse(chapters, hasAccess)

	utils.SendSuccessResponse(c, gin.H{
		"chapters":   chaptersResponse,
		"course_id":  courseID,
		"has_access": hasAccess,
	})
}

// CreateChapter crea un nuevo capítulo (migrado de createCapitulo)
func (chc *ChapterController) CreateChapter(c *gin.Context) {
	var req models.ChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error en el binding JSON para capítulo: %v", err)
		utils.SendErrorResponse(c, "Datos inválidos: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validar que el curso existe
	_, err := chc.courseService.GetCourseByID(req.CursoID, false)
	if err != nil {
		log.Printf("Curso no encontrado ID: %d para crear capítulo", req.CursoID)
		utils.SendErrorResponse(c, "Curso no encontrado", http.StatusNotFound)
		return
	}

	// Crear capítulo
	chapter, err := chc.chapterService.CreateChapter(&req)
	if err != nil {
		log.Printf("Error al crear capítulo: %v", err)
		utils.SendErrorResponse(c, "Error al crear capítulo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Capítulo creado exitosamente: ID %v", chapter.ID)
	utils.SendCreatedResponse(c, chapter)
}

// UpdateChapter actualiza un capítulo existente (migrado de updateCapitulo)
func (chc *ChapterController) UpdateChapter(c *gin.Context) {
	id := c.Param("id")
	chapterID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.SendErrorResponse(c, "ID de capítulo inválido", http.StatusBadRequest)
		return
	}

	// Verificar que el capítulo existe
	_, err = chc.chapterService.GetChapterByID(uint(chapterID))
	if err != nil {
		log.Printf("Capítulo no encontrado ID: %s para actualizar", id)
		utils.SendErrorResponse(c, "Capítulo no encontrado", http.StatusNotFound)
		return
	}

	var req models.ChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error en el binding JSON para actualizar capítulo: %v", err)
		utils.SendErrorResponse(c, "Datos inválidos: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Verificar que el curso existe
	_, err = chc.courseService.GetCourseByID(req.CursoID, false)
	if err != nil {
		log.Printf("Curso no encontrado ID: %d para actualizar capítulo", req.CursoID)
		utils.SendErrorResponse(c, "Curso no encontrado", http.StatusNotFound)
		return
	}

	// Actualizar capítulo
	updatedChapter, err := chc.chapterService.UpdateChapter(uint(chapterID), &req)
	if err != nil {
		log.Printf("Error al actualizar capítulo: %v", err)
		utils.SendErrorResponse(c, "Error al actualizar capítulo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Capítulo actualizado exitosamente: ID %v", updatedChapter.ID)
	utils.SendSuccessResponse(c, updatedChapter)
}

// DeleteChapter elimina un capítulo (migrado de deleteCapitulo)
func (chc *ChapterController) DeleteChapter(c *gin.Context) {
	id := c.Param("id")
	chapterID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.SendErrorResponse(c, "ID de capítulo inválido", http.StatusBadRequest)
		return
	}

	// Eliminar capítulo (incluye archivos asociados)
	err = chc.chapterService.DeleteChapter(uint(chapterID))
	if err != nil {
		log.Printf("Error al eliminar capítulo ID %s: %v", id, err)
		utils.SendErrorResponse(c, "Error al eliminar capítulo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Capítulo ID %s eliminado exitosamente", id)
	utils.SendSuccessResponse(c, gin.H{"message": "Capítulo eliminado correctamente"})
}

// GetChapterContent obtiene el contenido de un capítulo verificando acceso
func (chc *ChapterController) GetChapterContent(c *gin.Context) {
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

	// Obtener ID del capítulo
	id := c.Param("id")
	chapterID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.SendErrorResponse(c, "ID de capítulo inválido", http.StatusBadRequest)
		return
	}

	// Obtener capítulo
	chapter, err := chc.chapterService.GetChapterByID(uint(chapterID))
	if err != nil {
		utils.SendErrorResponse(c, "Capítulo no encontrado", http.StatusNotFound)
		return
	}

	// Verificar acceso al curso
	hasAccess := chc.courseService.CheckUserAccessToCourse(chapter.CursoID, userIDUint)

	// Verificar si puede acceder a este capítulo específico
	canAccess := chapter.CanUserAccess(hasAccess)
	if !canAccess {
		utils.SendErrorResponse(c, "No tienes acceso a este capítulo", http.StatusForbidden)
		return
	}

	// Preparar respuesta con contenido completo
	chapterResponse := chc.chapterService.PrepareChapterResponse(chapter, true)

	utils.SendSuccessResponse(c, chapterResponse)
}

// ReorderChapters reordena los capítulos de un curso
func (chc *ChapterController) ReorderChapters(c *gin.Context) {
	var req struct {
		CourseID   uint `json:"course_id" binding:"required"`
		ChapterIDs []uint `json:"chapter_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, "Datos inválidos: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Verificar que el curso existe
	_, err := chc.courseService.GetCourseByID(req.CourseID, false)
	if err != nil {
		utils.SendErrorResponse(c, "Curso no encontrado", http.StatusNotFound)
		return
	}

	// Reordenar capítulos
	err = chc.chapterService.ReorderChapters(req.CourseID, req.ChapterIDs)
	if err != nil {
		log.Printf("Error al reordenar capítulos: %v", err)
		utils.SendErrorResponse(c, "Error al reordenar capítulos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, gin.H{
		"message": "Capítulos reordenados correctamente",
		"course_id": req.CourseID,
	})
}

// ToggleChapterPublished alterna el estado de publicación de un capítulo
func (chc *ChapterController) ToggleChapterPublished(c *gin.Context) {
	id := c.Param("id")
	chapterID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.SendErrorResponse(c, "ID de capítulo inválido", http.StatusBadRequest)
		return
	}

	// Obtener capítulo actual
	chapter, err := chc.chapterService.GetChapterByID(uint(chapterID))
	if err != nil {
		utils.SendErrorResponse(c, "Capítulo no encontrado", http.StatusNotFound)
		return
	}

	// Alternar estado
	newStatus := !chapter.Publicado
	err = chc.chapterService.UpdateChapterPublishedStatus(uint(chapterID), newStatus)
	if err != nil {
		log.Printf("Error al actualizar estado de publicación: %v", err)
		utils.SendErrorResponse(c, "Error al actualizar estado", http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, gin.H{
		"message": "Estado de publicación actualizado",
		"chapter_id": chapterID,
		"published": newStatus,
	})
}