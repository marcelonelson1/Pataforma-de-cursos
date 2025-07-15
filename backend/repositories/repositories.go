package repositories

import "gorm.io/gorm"

// BaseRepository contiene la conexión a la base de datos
type BaseRepository struct {
	DB *gorm.DB
}

// NewBaseRepository crea una nueva instancia de BaseRepository
func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{DB: db}
}