package services

import (
	"curso-platform/config"
	"curso-platform/models"
	"curso-platform/utils"
	"fmt"
	"log"
)

// ProgresoService maneja la lógica relacionada con el progreso del usuario
type ProgresoService struct{}

// NewProgresoService crea una nueva instancia de ProgresoService
func NewProgresoService() *ProgresoService {
	return &ProgresoService{}
}

// GetProgresoUsuario obtiene el progreso de un usuario en un curso específico
func (s *ProgresoService) GetProgresoUsuario(usuarioID, cursoID uint) (*models.ProgresoResponse, error) {
	// Verificar si el curso existe
	var curso models.Curso
	if err := config.DB.First(&curso, cursoID).Error; err != nil {
		return nil, utils.ErrResourceNotFound
	}

	// Obtener progreso general del usuario en el curso
	var progresoUsuario models.ProgresoUsuario
	if err := config.DB.Where("usuario_id = ? AND curso_id = ?", usuarioID, cursoID).
		First(&progresoUsuario).Error; err != nil {
		// Si no existe, creamos un registro de progreso vacío
		progresoUsuario = models.ProgresoUsuario{
			UsuarioID:       usuarioID,
			CursoID:         cursoID,
			PorcentajeTotal: 0,
			UltimoCapitulo:  0,
		}
	}

	// Obtener progreso de cada capítulo
	var progresosCapitulos []models.ProgresoCapitulo
	if err := config.DB.Where("usuario_id = ? AND curso_id = ?", usuarioID, cursoID).
		Find(&progresosCapitulos).Error; err != nil {
		log.Printf("Error al obtener progreso de capítulos: %v", err)
	}

	// Convertir a mapa para facilitar el acceso en el frontend
	capitulosProgreso := make(map[uint]models.ProgresoCapitulo)
	for _, p := range progresosCapitulos {
		capitulosProgreso[p.CapituloID] = p
	}

	// Armar respuesta
	response := &models.ProgresoResponse{
		ProgresoTotal:     progresoUsuario.PorcentajeTotal,
		UltimoCapitulo:    progresoUsuario.UltimoCapitulo,
		CapitulosProgreso: capitulosProgreso,
	}

	return response, nil
}

// MarcarCapituloCompletado marca un capítulo como completado/incompleto
func (s *ProgresoService) MarcarCapituloCompletado(usuarioID uint, req struct {
	CursoID    uint    `json:"curso_id" binding:"required"`
	CapituloID uint    `json:"capitulo_id" binding:"required"`
	Completado bool    `json:"completado"`
	Progreso   float64 `json:"progreso"`
}) (*models.ProgresoCapitulo, error) {
	// Verificar que el curso y capítulo existan
	var curso models.Curso
	if err := config.DB.First(&curso, req.CursoID).Error; err != nil {
		return nil, utils.ErrResourceNotFound
	}

	var capitulo models.Capitulo
	if err := config.DB.First(&capitulo, req.CapituloID).Error; err != nil {
		return nil, utils.ErrResourceNotFound
	}

	// Verificar que el capítulo pertenezca al curso
	if capitulo.CursoID != req.CursoID {
		return nil, fmt.Errorf("el capítulo no pertenece al curso indicado")
	}

	// Buscar si existe un registro de progreso para este capítulo
	var progresoCapitulo models.ProgresoCapitulo
	result := config.DB.Where("usuario_id = ? AND curso_id = ? AND capitulo_id = ?",
		usuarioID, req.CursoID, req.CapituloID).First(&progresoCapitulo)

	if result.Error != nil {
		// Si no existe, crear uno nuevo
		progresoCapitulo = models.ProgresoCapitulo{
			UsuarioID:  usuarioID,
			CursoID:    req.CursoID,
			CapituloID: req.CapituloID,
			Completado: req.Completado,
			Progreso:   req.Progreso,
		}

		if err := config.DB.Create(&progresoCapitulo).Error; err != nil {
			log.Printf("Error al crear progreso de capítulo: %v", err)
			return nil, utils.ErrDatabaseError
		}
	} else {
		// Si existe, actualizar
		progresoCapitulo.Completado = req.Completado
		progresoCapitulo.Progreso = req.Progreso

		if err := config.DB.Save(&progresoCapitulo).Error; err != nil {
			log.Printf("Error al actualizar progreso de capítulo: %v", err)
			return nil, utils.ErrDatabaseError
		}
	}

	// Actualizar progreso total del curso
	s.ActualizarProgresoTotal(usuarioID, req.CursoID)

	return &progresoCapitulo, nil
}

// GuardarUltimoCapitulo guarda el último capítulo visto por el usuario
func (s *ProgresoService) GuardarUltimoCapitulo(usuarioID uint, req struct {
	CursoID    uint `json:"curso_id" binding:"required"`
	CapituloID uint `json:"capitulo_id" binding:"required"`
}) error {
	// Verificar que el curso y capítulo existan
	var curso models.Curso
	if err := config.DB.First(&curso, req.CursoID).Error; err != nil {
		return utils.ErrResourceNotFound
	}

	var capitulo models.Capitulo
	if err := config.DB.First(&capitulo, req.CapituloID).Error; err != nil {
		return utils.ErrResourceNotFound
	}

	// Verificar que el capítulo pertenezca al curso
	if capitulo.CursoID != req.CursoID {
		return fmt.Errorf("el capítulo no pertenece al curso indicado")
	}

	// Buscar o crear el registro de progreso del usuario
	var progresoUsuario models.ProgresoUsuario
	result := config.DB.Where("usuario_id = ? AND curso_id = ?", usuarioID, req.CursoID).First(&progresoUsuario)

	if result.Error != nil {
		// Si no existe, crear uno nuevo
		progresoUsuario = models.ProgresoUsuario{
			UsuarioID:       usuarioID,
			CursoID:         req.CursoID,
			UltimoCapitulo:  req.CapituloID,
			PorcentajeTotal: 0, // Se calculará en la siguiente función
		}

		if err := config.DB.Create(&progresoUsuario).Error; err != nil {
			log.Printf("Error al crear progreso de usuario: %v", err)
			return utils.ErrDatabaseError
		}
	} else {
		// Si existe, actualizar
		progresoUsuario.UltimoCapitulo = req.CapituloID

		if err := config.DB.Save(&progresoUsuario).Error; err != nil {
			log.Printf("Error al actualizar último capítulo: %v", err)
			return utils.ErrDatabaseError
		}
	}

	return nil
}

// ActualizarProgresoTotal actualiza el progreso total de un usuario en un curso
func (s *ProgresoService) ActualizarProgresoTotal(usuarioID, cursoID uint) {
	// Obtener el total de capítulos del curso
	var totalCapitulos int64
	if err := config.DB.Model(&models.Capitulo{}).Where("curso_id = ?", cursoID).Count(&totalCapitulos).Error; err != nil {
		log.Printf("Error al contar capítulos del curso: %v", err)
		return
	}

	if totalCapitulos == 0 {
		return // No hay capítulos para calcular progreso
	}

	// Obtener cantidad de capítulos completados
	var completados int64
	if err := config.DB.Model(&models.ProgresoCapitulo{}).
		Where("usuario_id = ? AND curso_id = ? AND completado = ?",
			usuarioID, cursoID, true).
		Count(&completados).Error; err != nil {
		log.Printf("Error al contar capítulos completados: %v", err)
		return
	}

	// Calcular porcentaje
	porcentaje := float64(completados) / float64(totalCapitulos) * 100

	// Actualizar o crear registro de progreso del usuario
	var progresoUsuario models.ProgresoUsuario
	result := config.DB.Where("usuario_id = ? AND curso_id = ?", usuarioID, cursoID).First(&progresoUsuario)

	if result.Error != nil {
		// Si no existe, crear uno nuevo
		progresoUsuario = models.ProgresoUsuario{
			UsuarioID:       usuarioID,
			CursoID:         cursoID,
			PorcentajeTotal: porcentaje,
		}

		if err := config.DB.Create(&progresoUsuario).Error; err != nil {
			log.Printf("Error al crear progreso de usuario: %v", err)
		}
	} else {
		// Si existe, actualizar
		progresoUsuario.PorcentajeTotal = porcentaje

		if err := config.DB.Save(&progresoUsuario).Error; err != nil {
			log.Printf("Error al actualizar progreso total: %v", err)
		}
	}
}