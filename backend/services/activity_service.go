package services

import (
	"curso-platform/config"
	"curso-platform/models"
	"curso-platform/utils"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/gorm"
)

// ActivityService gestiona la lógica de negocio para actividades
type ActivityService struct {
	db *gorm.DB
}

// NewActivityService crea una nueva instancia del servicio de actividades
func NewActivityService(db *gorm.DB) *ActivityService {
	return &ActivityService{
		db: db,
	}
}

// GetActivities obtiene todas las actividades con filtros opcionales
func (s *ActivityService) GetActivities(search, category string) ([]models.Activity, error) {
	var activities []models.Activity
	
	// Verificar si la tabla existe, si no, intentar crearla
	if !s.db.Migrator().HasTable(&models.Activity{}) {
		log.Println("La tabla 'activities' no existe, intentando crear...")
		if err := s.db.Migrator().CreateTable(&models.Activity{}); err != nil {
			return nil, fmt.Errorf("error al crear tabla 'activities': %w", err)
		}
		log.Println("Tabla 'activities' creada correctamente")
		return []models.Activity{}, nil // Devolver lista vacía ya que la tabla se acaba de crear
	}
	
	query := s.db.Model(&models.Activity{})

	if search != "" {
		searchQuery := "%" + search + "%"
		query = query.Where("title LIKE ? OR description LIKE ? OR instructor LIKE ?", searchQuery, searchQuery, searchQuery)
	}

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Find(&activities).Error; err != nil {
		if strings.Contains(err.Error(), "Table") && strings.Contains(err.Error(), "doesn't exist") {
			// Intentar ejecutar migraciones si la tabla no existe
			if err := config.RunMigrations(s.db); err != nil {
				return nil, fmt.Errorf("error al migrar base de datos: %w", err)
			}
			// Intentar nuevamente después de migrar
			if err := query.Find(&activities).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return activities, nil
}

// GetActivityByID obtiene una actividad por su ID
func (s *ActivityService) GetActivityByID(id uint) (*models.Activity, error) {
	var activity models.Activity
	
	// Verificar si la tabla existe
	if !s.db.Migrator().HasTable(&models.Activity{}) {
		return nil, fmt.Errorf("la tabla 'activities' no existe")
	}
	
	if err := s.db.First(&activity, id).Error; err != nil {
		return nil, fmt.Errorf("actividad no encontrada: %w", err)
	}
	return &activity, nil
}

// CreateActivity crea una nueva actividad
func (s *ActivityService) CreateActivity(req models.CreateActivityRequest, file *multipart.FileHeader) (*models.Activity, error) {
	// Verificar si la tabla existe
	if !s.db.Migrator().HasTable(&models.Activity{}) {
		log.Println("La tabla 'activities' no existe, intentando crear...")
		if err := s.db.Migrator().CreateTable(&models.Activity{}); err != nil {
			return nil, fmt.Errorf("error al crear tabla 'activities': %w", err)
		}
		log.Println("Tabla 'activities' creada correctamente")
	}
	
	// Crear la actividad con los datos del request
	activity := models.Activity{
		Title:       req.Title,
		Description: req.Description,
		Day:         req.Day,
		Time:        req.Time,
		Duration:    req.Duration,
		Instructor:  req.Instructor,
		Category:    req.Category,
		Capacity:    req.Capacity,
		Enrolled:    0,
	}

	// Guardar la imagen si se proporciona
	if file != nil {
		imageURL, err := s.saveActivityImage(file)
		if err != nil {
			return nil, fmt.Errorf("error al guardar la imagen: %w", err)
		}
		activity.ImageUrl = imageURL
	}

	// Guardar la actividad en la base de datos
	if err := s.db.Create(&activity).Error; err != nil {
		return nil, fmt.Errorf("error al crear la actividad: %w", err)
	}

	return &activity, nil
}

// UpdateActivity actualiza una actividad existente
func (s *ActivityService) UpdateActivity(id uint, req models.UpdateActivityRequest, file *multipart.FileHeader) (*models.Activity, error) {
	// Verificar si la tabla existe
	if !s.db.Migrator().HasTable(&models.Activity{}) {
		return nil, fmt.Errorf("la tabla 'activities' no existe")
	}
	
	// Obtener la actividad existente
	var activity models.Activity
	if err := s.db.First(&activity, id).Error; err != nil {
		return nil, fmt.Errorf("actividad no encontrada: %w", err)
	}

	// Actualizar campos si se proporcionan en el request
	if req.Title != "" {
		activity.Title = req.Title
	}
	if req.Description != "" {
		activity.Description = req.Description
	}
	if req.Day != "" {
		activity.Day = req.Day
	}
	if req.Time != "" {
		activity.Time = req.Time
	}
	if req.Duration != 0 {
		activity.Duration = req.Duration
	}
	if req.Instructor != "" {
		activity.Instructor = req.Instructor
	}
	if req.Category != "" {
		activity.Category = req.Category
	}
	if req.Capacity != 0 {
		activity.Capacity = req.Capacity
	}

	// Guardar la nueva imagen si se proporciona
	if file != nil {
		// Eliminar la imagen anterior si existe
		if activity.ImageUrl != "" {
			oldImagePath := filepath.Join("static", "activities", filepath.Base(activity.ImageUrl))
			if utils.FileExists(oldImagePath) {
				os.Remove(oldImagePath)
			}
		}

		// Guardar la nueva imagen
		imageURL, err := s.saveActivityImage(file)
		if err != nil {
			return nil, fmt.Errorf("error al guardar la imagen: %w", err)
		}
		activity.ImageUrl = imageURL
	}

	// Actualizar la actividad en la base de datos
	if err := s.db.Save(&activity).Error; err != nil {
		return nil, fmt.Errorf("error al actualizar la actividad: %w", err)
	}

	return &activity, nil
}

// DeleteActivity elimina una actividad

	
	// Obtener la actividad para verificar si existe y obtener la ruta de la imagen
	// DeleteActivity elimina una actividad
func (s *ActivityService) DeleteActivity(id uint) error {
	// Verificar si la tabla existe
	if !s.db.Migrator().HasTable(&models.Activity{}) {
		return fmt.Errorf("la tabla 'activities' no existe")
	}
	
	// Obtener la actividad para verificar si existe y obtener la ruta de la imagen
	var activity models.Activity
	if err := s.db.First(&activity, id).Error; err != nil {
		return fmt.Errorf("actividad no encontrada: %w", err)
	}

	// Eliminar la imagen si existe
	if activity.ImageUrl != "" {
		imagePath := filepath.Join("static", "activities", filepath.Base(activity.ImageUrl))
		if utils.FileExists(imagePath) {
			os.Remove(imagePath)
		}
	}

	// Eliminar la actividad de la base de datos
	if err := s.db.Delete(&activity).Error; err != nil {
		return fmt.Errorf("error al eliminar la actividad: %w", err)
	}

	return nil
}

// saveActivityImage guarda la imagen de la actividad y devuelve la URL
func (s *ActivityService) saveActivityImage(file *multipart.FileHeader) (string, error) {
	// Validar extensión de archivo
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !utils.IsValidImageExtension(ext) {
		return "", errors.New("tipo de archivo no válido. Solo se permiten imágenes (jpg, jpeg, png, gif, webp)")
	}

	// Generar nombre único para el archivo
	fileName := utils.GenerateUniqueFileName(file.Filename)
	
	// Crear directorio si no existe
	uploadDir := filepath.Join("static", "activities")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("error al crear directorio: %w", err)
	}
	
	// Ruta completa del archivo
	filePath := filepath.Join(uploadDir, fileName)
	
	// Guardar archivo
	if err := utils.SaveUploadedFile(file, filePath); err != nil {
		return "", fmt.Errorf("error al guardar archivo: %w", err)
	}
	
	// Devolver URL relativa para acceso desde el frontend
	return fmt.Sprintf("/static/activities/%s", fileName), nil
}

// GetEnrollmentsByActivity obtiene todas las inscripciones para una actividad específica
func (s *ActivityService) GetEnrollmentsByActivity(activityID uint) ([]models.Enrollment, error) {
	// Verificar si la tabla existe
	if !s.db.Migrator().HasTable(&models.Enrollment{}) {
		log.Println("La tabla 'enrollments' no existe, intentando crear...")
		if err := s.db.Migrator().CreateTable(&models.Enrollment{}); err != nil {
			return nil, fmt.Errorf("error al crear tabla 'enrollments': %w", err)
		}
		log.Println("Tabla 'enrollments' creada correctamente")
		return []models.Enrollment{}, nil
	}
	
	var enrollments []models.Enrollment
	if err := s.db.Where("activity_id = ?", activityID).Find(&enrollments).Error; err != nil {
		return nil, fmt.Errorf("error al obtener inscripciones: %w", err)
	}
	
	return enrollments, nil
}

// GetUserEnrollments obtiene todas las inscripciones de un usuario
func (s *ActivityService) GetUserEnrollments(userID uint) ([]models.Enrollment, error) {
	// Verificar si la tabla existe
	if !s.db.Migrator().HasTable(&models.Enrollment{}) {
		log.Println("La tabla 'enrollments' no existe, intentando crear...")
		if err := s.db.Migrator().CreateTable(&models.Enrollment{}); err != nil {
			return nil, fmt.Errorf("error al crear tabla 'enrollments': %w", err)
		}
		log.Println("Tabla 'enrollments' creada correctamente")
		return []models.Enrollment{}, nil
	}
	
	var enrollments []models.Enrollment
	if err := s.db.Where("user_id = ?", userID).Find(&enrollments).Error; err != nil {
		return nil, fmt.Errorf("error al obtener inscripciones: %w", err)
	}
	
	return enrollments, nil
}

// EnrollUser inscribe a un usuario en una actividad
func (s *ActivityService) EnrollUser(userID, activityID uint) error {
	// Verificar si la tabla de actividades existe
	if !s.db.Migrator().HasTable(&models.Activity{}) {
		return fmt.Errorf("la tabla 'activities' no existe")
	}
	
	// Verificar si la tabla de inscripciones existe
	if !s.db.Migrator().HasTable(&models.Enrollment{}) {
		log.Println("La tabla 'enrollments' no existe, intentando crear...")
		if err := s.db.Migrator().CreateTable(&models.Enrollment{}); err != nil {
			return fmt.Errorf("error al crear tabla 'enrollments': %w", err)
		}
		log.Println("Tabla 'enrollments' creada correctamente")
	}
	
	// Verificar si la actividad existe y tiene capacidad
	var activity models.Activity
	if err := s.db.First(&activity, activityID).Error; err != nil {
		return fmt.Errorf("actividad no encontrada: %w", err)
	}
	
	// Verificar si hay cupo disponible
	if activity.Enrolled >= activity.Capacity {
		return errors.New("la actividad ya está llena")
	}
	
	// Verificar si el usuario ya está inscrito
	var existingEnrollment models.Enrollment
	result := s.db.Where("user_id = ? AND activity_id = ?", userID, activityID).First(&existingEnrollment)
	if result.Error == nil {
		return errors.New("el usuario ya está inscrito en esta actividad")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("error al verificar inscripción: %w", result.Error)
	}
	
	// Crear la inscripción
	enrollment := models.Enrollment{
		UserID:     userID,
		ActivityID: activityID,
		Status:     "active",
	}
	
	// Iniciar transacción
	tx := s.db.Begin()
	
	// Guardar la inscripción
	if err := tx.Create(&enrollment).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error al crear inscripción: %w", err)
	}
	
	// Actualizar contador de inscritos
	activity.Enrolled++
	if err := tx.Save(&activity).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error al actualizar contador de inscritos: %w", err)
	}
	
	// Confirmar transacción
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("error al confirmar transacción: %w", err)
	}
	
	return nil
}

// CancelEnrollment cancela la inscripción de un usuario a una actividad
func (s *ActivityService) CancelEnrollment(userID, activityID uint) error {
	// Verificar si las tablas existen
	if !s.db.Migrator().HasTable(&models.Activity{}) || !s.db.Migrator().HasTable(&models.Enrollment{}) {
		return fmt.Errorf("las tablas necesarias no existen")
	}
	
	// Verificar si la inscripción existe
	var enrollment models.Enrollment
	if err := s.db.Where("user_id = ? AND activity_id = ?", userID, activityID).First(&enrollment).Error; err != nil {
		return fmt.Errorf("inscripción no encontrada: %w", err)
	}
	
	// Obtener la actividad
	var activity models.Activity
	if err := s.db.First(&activity, activityID).Error; err != nil {
		return fmt.Errorf("actividad no encontrada: %w", err)
	}
	
	// Iniciar transacción
	tx := s.db.Begin()
	
	// Eliminar la inscripción
	if err := tx.Delete(&enrollment).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error al eliminar inscripción: %w", err)
	}
	
	// Actualizar contador de inscritos
	if activity.Enrolled > 0 {
		activity.Enrolled--
		if err := tx.Save(&activity).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("error al actualizar contador de inscritos: %w", err)
		}
	}
	
	// Confirmar transacción
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("error al confirmar transacción: %w", err)
	}
	
	return nil
}