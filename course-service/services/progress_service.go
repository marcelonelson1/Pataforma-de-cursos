// services/progress_service.go
package services

import (
	"fmt"
	"log"
	

	"gorm.io/gorm"

	"course-service/config"
	"course-service/models"
)

type ProgressService struct {
	db     *gorm.DB
	config *config.Config
}

func NewProgressService(db *gorm.DB, cfg *config.Config) *ProgressService {
	return &ProgressService{
		db:     db,
		config: cfg,
	}
}

// GetUserProgress obtiene el progreso completo de un usuario en un curso
func (ps *ProgressService) GetUserProgress(userID, courseID uint) (*models.ProgressResponse, error) {
	// Obtener progreso general del usuario en el curso
	var userProgress models.UserProgress
	result := ps.db.Where("usuario_id = ? AND curso_id = ?", userID, courseID).First(&userProgress)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Si no existe, crear uno nuevo
			userProgress = models.UserProgress{
				UsuarioID:       userID,
				CursoID:         courseID,
				PorcentajeTotal: 0,
				UltimoCapitulo:  0,
			}
		} else {
			return nil, fmt.Errorf("error al obtener progreso del usuario: %v", result.Error)
		}
	}

	// Obtener progreso de cada capítulo
	var chapterProgresses []models.ChapterProgress
	if err := ps.db.Where("usuario_id = ? AND curso_id = ?", userID, courseID).
		Find(&chapterProgresses).Error; err != nil {
		log.Printf("Error al obtener progreso de capítulos: %v", err)
	}

	// Convertir a mapa para facilitar acceso
	chapterProgressMap := make(map[uint]models.ChapterProgress)
	for _, cp := range chapterProgresses {
		chapterProgressMap[cp.CapituloID] = cp
	}

	// Obtener estadísticas del curso
	stats, err := ps.getCourseStatistics(userID, courseID)
	if err != nil {
		log.Printf("Error al obtener estadísticas del curso: %v", err)
		stats = models.CourseStatistics{}
	}

	// Construir respuesta
	response := &models.ProgressResponse{
		ProgresoTotal:     userProgress.PorcentajeTotal,
		UltimoCapitulo:    userProgress.UltimoCapitulo,
		TiempoTotal:       userProgress.TiempoTotal,
		FechaInicio:       userProgress.FechaInicio,
		FechaCompletado:   userProgress.FechaCompletado,
		CapitulosProgreso: chapterProgressMap,
		EstadisticasCurso: stats,
	}

	return response, nil
}

// UpdateChapterProgress actualiza el progreso de un capítulo
func (ps *ProgressService) UpdateChapterProgress(userID uint, req *models.ChapterProgressRequest) (*models.ChapterProgress, error) {
	// Verificar que el capítulo existe y pertenece al curso
	var chapter models.Chapter
	if err := ps.db.Where("id = ? AND curso_id = ?", req.CapituloID, req.CursoID).First(&chapter).Error; err != nil {
		return nil, fmt.Errorf("capítulo no encontrado o no pertenece al curso")
	}

	// Buscar progreso existente del capítulo
	var chapterProgress models.ChapterProgress
	result := ps.db.Where("usuario_id = ? AND curso_id = ? AND capitulo_id = ?",
		userID, req.CursoID, req.CapituloID).First(&chapterProgress)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Crear nuevo progreso
			chapterProgress = models.ChapterProgress{
				UsuarioID:  userID,
				CursoID:    req.CursoID,
				CapituloID: req.CapituloID,
			}
		} else {
			return nil, fmt.Errorf("error al buscar progreso del capítulo: %v", result.Error)
		}
	}

	// Actualizar progreso
	chapterProgress.UpdateProgress(req.Progreso, req.TiempoVisto)
	if req.Completado {
		chapterProgress.MarkAsCompleted()
	}

	// Guardar cambios
	if err := ps.db.Save(&chapterProgress).Error; err != nil {
		return nil, fmt.Errorf("error al guardar progreso del capítulo: %v", err)
	}

	// Actualizar progreso total del curso
	if err := ps.updateCourseProgress(userID, req.CursoID); err != nil {
		log.Printf("Error al actualizar progreso total del curso: %v", err)
	}

	log.Printf("Progreso del capítulo actualizado: UserID=%d, ChapterID=%d, Progress=%.2f", 
		userID, req.CapituloID, req.Progreso)

	return &chapterProgress, nil
}

// UpdateLastChapter actualiza el último capítulo visto
func (ps *ProgressService) UpdateLastChapter(userID, courseID, chapterID uint) error {
	// Verificar que el capítulo existe y pertenece al curso
	var chapter models.Chapter
	if err := ps.db.Where("id = ? AND curso_id = ?", chapterID, courseID).First(&chapter).Error; err != nil {
		return fmt.Errorf("capítulo no encontrado o no pertenece al curso")
	}

	// Buscar o crear progreso del usuario
	var userProgress models.UserProgress
	result := ps.db.Where("usuario_id = ? AND curso_id = ?", userID, courseID).First(&userProgress)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Crear nuevo progreso
			userProgress = models.UserProgress{
				UsuarioID:      userID,
				CursoID:        courseID,
				UltimoCapitulo: chapterID,
			}
			userProgress.MarkAsStarted()
		} else {
			return fmt.Errorf("error al buscar progreso del usuario: %v", result.Error)
		}
	} else {
		// Actualizar último capítulo
		userProgress.UltimoCapitulo = chapterID
	}

	// Guardar cambios
	if err := ps.db.Save(&userProgress).Error; err != nil {
		return fmt.Errorf("error al actualizar último capítulo: %v", err)
	}

	log.Printf("Último capítulo actualizado: UserID=%d, CourseID=%d, ChapterID=%d", 
		userID, courseID, chapterID)

	return nil
}

// UpdateChapterWatchTime actualiza solo el tiempo visto de un capítulo
func (ps *ProgressService) UpdateChapterWatchTime(userID, courseID, chapterID uint, timeWatched int) error {
	// Buscar progreso existente del capítulo
	var chapterProgress models.ChapterProgress
	result := ps.db.Where("usuario_id = ? AND curso_id = ? AND capitulo_id = ?",
		userID, courseID, chapterID).First(&chapterProgress)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Crear nuevo progreso
			chapterProgress = models.ChapterProgress{
				UsuarioID:   userID,
				CursoID:     courseID,
				CapituloID:  chapterID,
				TiempoVisto: timeWatched,
			}
		} else {
			return fmt.Errorf("error al buscar progreso del capítulo: %v", result.Error)
		}
	} else {
		// Actualizar tiempo visto (solo si es mayor al actual)
		if timeWatched > chapterProgress.TiempoVisto {
			chapterProgress.TiempoVisto = timeWatched
		}
	}

	// Guardar cambios
	if err := ps.db.Save(&chapterProgress).Error; err != nil {
		return fmt.Errorf("error al actualizar tiempo visto: %v", err)
	}

	return nil
}

// GetUserProgressSummary obtiene resumen del progreso del usuario en todos los cursos
func (ps *ProgressService) GetUserProgressSummary(userID uint) (*models.UserProgressSummary, error) {
	// Obtener todos los progresos del usuario
	var userProgresses []models.UserProgress
	if err := ps.db.Where("usuario_id = ?", userID).Find(&userProgresses).Error; err != nil {
		return nil, fmt.Errorf("error al obtener progresos del usuario: %v", err)
	}

	// Calcular estadísticas
	totalCursos := len(userProgresses)
	cursosCompletados := 0
	cursosEnProgreso := 0
	tiempoTotalEstudio := 0
	var sumaCompletitud float64

	for _, progress := range userProgresses {
		tiempoTotalEstudio += progress.TiempoTotal
		sumaCompletitud += progress.PorcentajeTotal

		if progress.IsCompleted() {
			cursosCompletados++
		} else if progress.PorcentajeTotal > 0 {
			cursosEnProgreso++
		}
	}

	// Promedio de completitud
	var promedioCompletitud float64
	if totalCursos > 0 {
		promedioCompletitud = sumaCompletitud / float64(totalCursos)
	}

	// Obtener cursos recientes (últimos 5)
	cursosRecientes, err := ps.getRecentCourses(userID, 5)
	if err != nil {
		log.Printf("Error al obtener cursos recientes: %v", err)
		cursosRecientes = []models.CourseProgress{}
	}

	return &models.UserProgressSummary{
		TotalCursos:         totalCursos,
		CursosCompletados:   cursosCompletados,
		CursosEnProgreso:    cursosEnProgreso,
		TiempoTotalEstudio:  tiempoTotalEstudio,
		PromedioCompletitud: promedioCompletitud,
		CursosRecientes:     cursosRecientes,
	}, nil
}

// GetCourseStatistics obtiene estadísticas específicas de un curso para un usuario
func (ps *ProgressService) GetCourseStatistics(userID, courseID uint) (*models.CourseStatistics, error) {
	stats, err := ps.getCourseStatistics(userID, courseID)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

// updateCourseProgress actualiza el progreso total de un curso basado en capítulos completados
func (ps *ProgressService) updateCourseProgress(userID, courseID uint) error {
	// Obtener total de capítulos publicados del curso
	var totalChapters int64
	if err := ps.db.Model(&models.Chapter{}).
		Where("curso_id = ? AND publicado = true", courseID).
		Count(&totalChapters).Error; err != nil {
		return fmt.Errorf("error al contar capítulos: %v", err)
	}

	if totalChapters == 0 {
		return nil // No hay capítulos para calcular progreso
	}

	// Obtener capítulos completados por el usuario
	var completedChapters int64
	if err := ps.db.Model(&models.ChapterProgress{}).
		Where("usuario_id = ? AND curso_id = ? AND completado = true", userID, courseID).
		Count(&completedChapters).Error; err != nil {
		return fmt.Errorf("error al contar capítulos completados: %v", err)
	}

	// Calcular porcentaje
	porcentaje := models.GetProgressPercentage(int(totalChapters), int(completedChapters))

	// Obtener tiempo total visto
	var tiempoTotal int
	if err := ps.db.Model(&models.ChapterProgress{}).
		Where("usuario_id = ? AND curso_id = ?", userID, courseID).
		Select("COALESCE(SUM(tiempo_visto), 0)").
		Scan(&tiempoTotal).Error; err != nil {
		log.Printf("Error al calcular tiempo total: %v", err)
	}

	// Actualizar o crear progreso del usuario
	var userProgress models.UserProgress
	result := ps.db.Where("usuario_id = ? AND curso_id = ?", userID, courseID).First(&userProgress)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Crear nuevo progreso
			userProgress = models.UserProgress{
				UsuarioID:       userID,
				CursoID:         courseID,
				PorcentajeTotal: porcentaje,
				TiempoTotal:     tiempoTotal,
			}
			userProgress.MarkAsStarted()
		} else {
			return fmt.Errorf("error al buscar progreso del usuario: %v", result.Error)
		}
	} else {
		// Actualizar progreso existente
		userProgress.UpdateTotalProgress(int(totalChapters), int(completedChapters))
		userProgress.TiempoTotal = tiempoTotal
	}

	// Guardar cambios
	if err := ps.db.Save(&userProgress).Error; err != nil {
		return fmt.Errorf("error al actualizar progreso del curso: %v", err)
	}

	return nil
}

// getCourseStatistics obtiene estadísticas de un curso para un usuario
func (ps *ProgressService) getCourseStatistics(userID, courseID uint) (models.CourseStatistics, error) {
	// Total de capítulos publicados
	var totalChapters int64
	if err := ps.db.Model(&models.Chapter{}).
		Where("curso_id = ? AND publicado = true", courseID).
		Count(&totalChapters).Error; err != nil {
		return models.CourseStatistics{}, err
	}

	// Capítulos completados por el usuario
	var completedChapters int64
	if err := ps.db.Model(&models.ChapterProgress{}).
		Where("usuario_id = ? AND curso_id = ? AND completado = true", userID, courseID).
		Count(&completedChapters).Error; err != nil {
		return models.CourseStatistics{}, err
	}

	return models.CourseStatistics{
		TotalCapitulos:       int(totalChapters),
		CapitulosCompletados: int(completedChapters),
		CapitulosPendientes:  int(totalChapters) - int(completedChapters),
	}, nil
}

// getRecentCourses obtiene los cursos recientes del usuario
func (ps *ProgressService) getRecentCourses(userID uint, limit int) ([]models.CourseProgress, error) {
	var progresses []models.UserProgress
	if err := ps.db.Where("usuario_id = ?", userID).
		Order("updated_at DESC").
		Limit(limit).
		Find(&progresses).Error; err != nil {
		return nil, err
	}

	var courseProgresses []models.CourseProgress
	for _, progress := range progresses {
		// Obtener información del curso
		var course models.Course
		if err := ps.db.Select("titulo, imagen_url").First(&course, progress.CursoID).Error; err != nil {
			log.Printf("Error al obtener curso %d: %v", progress.CursoID, err)
			continue
		}

		courseProgress := models.CourseProgress{
			CursoID:        progress.CursoID,
			TituloCurso:    course.Titulo,
			ImagenURL:      course.ImagenURL,
			Progreso:       progress.PorcentajeTotal,
			UltimoCapitulo: progress.UltimoCapitulo,
			FechaUltimaVez: progress.UpdatedAt,
		}

		courseProgresses = append(courseProgresses, courseProgress)
	}

	return courseProgresses, nil
}