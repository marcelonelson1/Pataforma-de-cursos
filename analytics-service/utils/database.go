package utils

import (
	"analytics-service/config"
	"analytics-service/models"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabase(cfg *config.Config) error {
	var err error
	
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("✅ Database connection established")
	return nil
}

func RunMigrations() error {
	log.Println("🔄 Running database migrations...")
	
	err := DB.AutoMigrate(
		&models.User{},
		&models.SalesStats{},
		&models.DashboardMetrics{},
		&models.ActivityLog{},
	)
	
	if err != nil {
		return fmt.Errorf("migration failed: %v", err)
	}
	
	log.Println("✅ Database migrations completed")
	return nil
}