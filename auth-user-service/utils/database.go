package utils

import (
	"auth-user-service/config"
	"auth-user-service/models"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// SetupDatabase configura la conexión a la base de datos
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

		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err == nil {
			sqlDB, err := DB.DB()
			if err != nil {
				return fmt.Errorf("error al obtener instancia DB: %v", err)
			}

			sqlDB.SetMaxIdleConns(10)
			sqlDB.SetMaxOpenConns(100)
			sqlDB.SetConnMaxLifetime(time.Hour)

			if err := runMigrations(); err != nil {
				return fmt.Errorf("error en migración: %v", err)
			}

			log.Println("Conexión a la base de datos establecida correctamente")
			return nil
		}

		log.Printf("Error de conexión: %v. Reintentando en %v...", err, retryDelay)
		time.Sleep(retryDelay)
	}

	return fmt.Errorf("no se pudo conectar a la base de datos después de %d intentos: %v", maxRetries, err)
}

// runMigrations ejecuta las migraciones de la base de datos
func runMigrations() error {
	log.Println("Ejecutando migraciones...")

	// Migrar modelos
	err := DB.AutoMigrate(
		&models.Usuario{},
		&models.PasswordReset{},
		&models.ActivityLog{},
		&models.NotificationSettings{},
	)

	if err != nil {
		return fmt.Errorf("error al migrar tablas: %v", err)
	}

	// Crear índices adicionales
	if err := createIndexes(); err != nil {
		return fmt.Errorf("error al crear índices: %v", err)
	}

	log.Println("Migraciones completadas exitosamente")
	return nil
}

// createIndexes crea índices adicionales para optimizar consultas
func createIndexes() error {
	// Índice para búsquedas de activity logs por usuario
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_activity_user_action ON activity_logs(user_id, action)").Error; err != nil {
		log.Printf("Advertencia: No se pudo crear índice idx_activity_user_action: %v", err)
	}

	// Índice para búsquedas de password resets por email y estado
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_password_reset_email_used ON password_resets(email, used)").Error; err != nil {
		log.Printf("Advertencia: No se pudo crear índice idx_password_reset_email_used: %v", err)
	}

	// Índice para búsquedas por rol
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_usuario_role ON usuarios(role)").Error; err != nil {
		log.Printf("Advertencia: No se pudo crear índice idx_usuario_role: %v", err)
	}

	return nil
}

// GetDB retorna la instancia de la base de datos
func GetDB() *gorm.DB {
	return DB
}