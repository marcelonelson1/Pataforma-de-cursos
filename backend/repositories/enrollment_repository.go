package repositories

import (
	"curso-platform/models"
"errors"
	"gorm.io/gorm"
)

// EnrollmentRepository define las operaciones para inscripciones
type EnrollmentRepository interface {
	FindByUserAndActivity(userID, activityID uint) (*models.Enrollment, error)
	FindByUser(userID uint) ([]models.Enrollment, error)
	FindByActivity(activityID uint) ([]models.Enrollment, error)
	Create(enrollment *models.Enrollment) (*models.Enrollment, error)
}

type enrollmentRepository struct {
	*BaseRepository
}

// NewEnrollmentRepository crea una nueva instancia de EnrollmentRepository
func NewEnrollmentRepository(db *gorm.DB) EnrollmentRepository {
	return &enrollmentRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *enrollmentRepository) FindByUserAndActivity(userID, activityID uint) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	err := r.DB.Where("user_id = ? AND activity_id = ?", userID, activityID).First(&enrollment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &enrollment, nil
}

func (r *enrollmentRepository) FindByUser(userID uint) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	if err := r.DB.Where("user_id = ?", userID).Find(&enrollments).Error; err != nil {
		return nil, err
	}
	return enrollments, nil
}

func (r *enrollmentRepository) FindByActivity(activityID uint) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	if err := r.DB.Where("activity_id = ?", activityID).Find(&enrollments).Error; err != nil {
		return nil, err
	}
	return enrollments, nil
}

func (r *enrollmentRepository) Create(enrollment *models.Enrollment) (*models.Enrollment, error) {
	if err := r.DB.Create(enrollment).Error; err != nil {
		return nil, err
	}
	return enrollment, nil
}