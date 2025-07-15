package repositories

import (
	"curso-platform/models"
	"errors"

	"gorm.io/gorm"
)

// ActivityRepository define las operaciones para actividades deportivas
type ActivityRepository interface {
	FindAll(filters map[string]interface{}) ([]models.Activity, error)
	FindByID(id uint) (*models.Activity, error)
	FindByIDs(ids []uint) ([]models.Activity, error)
	Create(activity *models.Activity) (*models.Activity, error)
	Update(activity *models.Activity) (*models.Activity, error)
	Delete(id uint) error
}

type activityRepository struct {
	*BaseRepository
}

// NewActivityRepository crea una nueva instancia de ActivityRepository
func NewActivityRepository(db *gorm.DB) ActivityRepository {
	return &activityRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *activityRepository) FindAll(filters map[string]interface{}) ([]models.Activity, error) {
	var activities []models.Activity
	query := r.DB.Model(&models.Activity{})

	if search, ok := filters["search"]; ok {
		query = query.Where("title LIKE ?", "%"+search.(string)+"%")
	}

	if category, ok := filters["category"]; ok {
		query = query.Where("category = ?", category.(string))
	}

	if err := query.Find(&activities).Error; err != nil {
		return nil, err
	}

	return activities, nil
}

func (r *activityRepository) FindByID(id uint) (*models.Activity, error) {
	var activity models.Activity
	if err := r.DB.First(&activity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &activity, nil
}

func (r *activityRepository) FindByIDs(ids []uint) ([]models.Activity, error) {
	var activities []models.Activity
	if err := r.DB.Where("id IN ?", ids).Find(&activities).Error; err != nil {
		return nil, err
	}
	return activities, nil
}

func (r *activityRepository) Create(activity *models.Activity) (*models.Activity, error) {
	if err := r.DB.Create(activity).Error; err != nil {
		return nil, err
	}
	return activity, nil
}

func (r *activityRepository) Update(activity *models.Activity) (*models.Activity, error) {
	if err := r.DB.Save(activity).Error; err != nil {
		return nil, err
	}
	return activity, nil
}

func (r *activityRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Activity{}, id).Error
}