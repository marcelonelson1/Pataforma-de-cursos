// controllers/progress_controller.go
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

type ProgressController struct {
	db              *gorm.DB
	config          *config.Config
	progressService *services.ProgressService
	courseService   *services.CourseService
}

func NewProgressController(db *gorm.DB, cfg *config.Config) *ProgressController {
	return &ProgressController{
		db:              db,
		config:          cfg,
		progressService: services.NewProgressService(db, cfg),
		courseService:   services.NewCourseService(db, cfg),
	}
}

// GetUserProgress obtiene el progreso de un usuario en un curso (migrado de getProgresoUsuario)
func (pc *ProgressController) GetUserProgress(c *gin.Context) {
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
	cursoIDStr := c.Param("id")
	if cursoIDStr == "" {
		utils.SendErrorResponse(c, "ID de curso requerido", http.StatusBadRequest)
		return
	}

	cursoID, err := strconv.ParseUint(cursoIDStr, 10, 32)
	if err != nil {
		utils.SendErrorResponse(c, "ID de curso inválido", http.StatusBadRequest)
		return
	}

	log.Printf("Obteniendo progreso para curso ID: %d, usuario ID: %d", cursoID, userIDUint)

	// Verificar que el curso existe
	_, err = pc.courseService.GetCourseByID(uint(cursoID), false)
	if err != nil {
		utils.SendErrorResponse(c, "Curso no encontrado", http.StatusNotFound)
		return
	}

	// Verificar acceso al curso
	hasAccess := pc.courseService.CheckUserAccessToCourse(uint(cursoID), userIDUint)
	if !hasAccess {
		utils.SendErrorResponse(c, "No tienes acceso a este curso", http.StatusForbidden)
		return
	}

	// Obtener progreso del usuario
	progress, err := pc.progressService.GetUserProgress(userIDUint, uint(cursoID))
	if err != nil {
		log.Printf("Error al obtener progreso: %v", err)
			utils.SendErrorResponse(c, "Error al obtener progreso", http.StatusInternalServerError)
			return
		}
	
		utils.SendSuccessResponse(c, gin.H{
			"success":  true,
			"message":  "Progreso obtenido correctamente",
			"progreso": progress,
		})
	}

// MarkChapterCompleted marca un capítulo como completado (migrado de marcarCapituloCompletado)
func (pc *ProgressController) MarkChapterCompleted(c *gin.Context) {
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

	// Obtener datos de la solicitud
	var req models.ChapterProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, "Datos inválidos: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Marcando capítulo completado: UserID=%d, CursoID=%d, CapituloID=%d, Completado=%v", 
		userIDUint, req.CursoID, req.CapituloID, req.Completado)

	// Verificar que el curso existe
	_, err := pc.courseService.GetCourseByID(req.CursoID, false)
	if err != nil {
		utils.SendErrorResponse(c, "Curso no encontrado", http.StatusNotFound)
		return
	}

	// Verificar acceso al curso
	hasAccess := pc.courseService.CheckUserAccessToCourse(req.CursoID, userIDUint)
	if !hasAccess {
		utils.SendErrorResponse(c, "No tienes acceso a este curso", http.StatusForbidden)
		return
	}

	// Actualizar progreso del capítulo
	progress, err := pc.progressService.UpdateChapterProgress(userIDUint, &req)
	if err != nil {
		log.Printf("Error al actualizar progreso de capítulo: %v", err)
		utils.SendErrorResponse(c, "Error al actualizar progreso: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, gin.H{
		"success":  true,
		"message":  "Progreso actualizado correctamente",
		"progreso": progress,
	})
}

// SaveLastChapter guarda el último capítulo visto (migrado de guardarUltimoCapitulo)
func (pc *ProgressController) SaveLastChapter(c *gin.Context) {
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

	// Obtener datos de la solicitud
	var req models.LastChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, "Datos inválidos: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Guardando último capítulo: UserID=%d, CursoID=%d, CapituloID=%d", 
		userIDUint, req.CursoID, req.CapituloID)

	// Verificar que el curso existe
	_, err := pc.courseService.GetCourseByID(req.CursoID, false)
	if err != nil {
		utils.SendErrorResponse(c, "Curso no encontrado", http.StatusNotFound)
		return
	}

	// Verificar acceso al curso
	hasAccess := pc.courseService.CheckUserAccessToCourse(req.CursoID, userIDUint)
	if !hasAccess {
		utils.SendErrorResponse(c, "No tienes acceso a este curso", http.StatusForbidden)
		return
	}

	// Actualizar último capítulo
	err = pc.progressService.UpdateLastChapter(userIDUint, req.CursoID, req.CapituloID)
	if err != nil {
		log.Printf("Error al actualizar último capítulo: %v", err)
		utils.SendErrorResponse(c, "Error al actualizar progreso: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, gin.H{
		"success": true,
		"message": "Último capítulo actualizado correctamente",
	})
}

// GetUserProgressSummary obtiene resumen del progreso del usuario en todos los cursos
func (pc *ProgressController) GetUserProgressSummary(c *gin.Context) {
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

	// Obtener resumen de progreso
	summary, err := pc.progressService.GetUserProgressSummary(userIDUint)
	if err != nil {
		log.Printf("Error al obtener resumen de progreso: %v", err)
		utils.SendErrorResponse(c, "Error al obtener resumen de progreso", http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, summary)
}

// UpdateChapterWatchTime actualiza el tiempo visto de un capítulo
func (pc *ProgressController) UpdateChapterWatchTime(c *gin.Context) {
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

	// Obtener datos de la solicitud
	var req struct {
		CursoID    uint `json:"curso_id" binding:"required"`
		CapituloID uint `json:"capitulo_id" binding:"required"`
		TiempoVisto int `json:"tiempo_visto" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, "Datos inválidos: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Verificar acceso al curso
	hasAccess := pc.courseService.CheckUserAccessToCourse(req.CursoID, userIDUint)
	if !hasAccess {
		utils.SendErrorResponse(c, "No tienes acceso a este curso", http.StatusForbidden)
		return
	}

	// Actualizar tiempo visto
	err := pc.progressService.UpdateChapterWatchTime(userIDUint, req.CursoID, req.CapituloID, req.TiempoVisto)
	if err != nil {
		log.Printf("Error al actualizar tiempo visto: %v", err)
		utils.SendErrorResponse(c, "Error al actualizar tiempo visto", http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, gin.H{
		"success": true,
		"message": "Tiempo visto actualizado correctamente",
	})
}

// GetCourseStatistics obtiene estadísticas de un curso
func (pc *ProgressController) GetCourseStatistics(c *gin.Context) {
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
	cursoIDStr := c.Param("id")
	cursoID, err := strconv.ParseUint(cursoIDStr, 10, 32)
	if err != nil {
		utils.SendErrorResponse(c, "ID de curso inválido", http.StatusBadRequest)
		return
	}

	// Verificar acceso al curso
	hasAccess := pc.courseService.CheckUserAccessToCourse(uint(cursoID), userIDUint)
	if !hasAccess {
		utils.SendErrorResponse(c, "No tienes acceso a este curso", http.StatusForbidden)
		return
	}

	// Obtener estadísticas
	stats, err := pc.progressService.GetCourseStatistics(userIDUint, uint(cursoID))
	if err != nil {
		log.Printf("Error al obtener estadísticas del curso: %v", err)
		utils.SendErrorResponse(c, "Error al obtener estadísticas", http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(c, stats)
}