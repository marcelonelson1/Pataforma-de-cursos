package utils

import (
	"contact-service/config"
	"contact-service/models"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// GetDB retorna la instancia de la base de datos
func GetDB() *gorm.DB {
	return db
}

// SetupDatabase configura la conexion a la base de datos
func SetupDatabase() error {
	var err error
	maxRetries := 5
	retryDelay := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		log.Printf("Conectando a la base de datos (intento %d/%d)...", i+1, maxRetries)

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.AppConfig.DBUser,
			config.AppConfig.DBPassword,
			config.AppConfig.DBHost,
			config.AppConfig.DBPort,
			config.AppConfig.DBName)

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err == nil {
			sqlDB, err := db.DB()
			if err != nil {
				return fmt.Errorf("error al obtener instancia DB: %v", err)
			}

			sqlDB.SetMaxIdleConns(10)
			sqlDB.SetMaxOpenConns(100)
			sqlDB.SetConnMaxLifetime(time.Hour)

			if err := runMigrations(); err != nil {
				return fmt.Errorf("error en migraciones: %v", err)
			}

			log.Println("Conexion a la base de datos establecida correctamente")
			return nil
		}

		log.Printf("Error de conexion: %v. Reintentando en %v...", err, retryDelay)
		time.Sleep(retryDelay)
	}

	return fmt.Errorf("no se pudo conectar a la base de datos despues de %d intentos: %v", maxRetries, err)
}

// runMigrations ejecuta las migraciones de la base de datos
func runMigrations() error {
	// Migrar modelos
	if err := db.AutoMigrate(&models.ContactMessage{}); err != nil {
		return fmt.Errorf("error al migrar ContactMessage: %v", err)
	}

	log.Println("Migraciones ejecutadas correctamente")
	return nil
}