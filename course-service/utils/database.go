// utils/database.go
package utils

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"course-service/config"
	"course-service/models"
)

// ConnectDatabase conecta a la base de datos con reintentos
func ConnectDatabase(cfg *config.Config) *gorm.DB {
	var db *gorm.DB
	var err error
	maxRetries := 5
	retryDelay := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		log.Printf("Conectando a la base de datos (intento %d/%d)...", i+1, maxRetries)

		db, err = gorm.Open(mysql.Open(cfg.Database.DSN()), &gorm.Config{
			Logger: logger.Default.LogMode(getLogLevel(cfg.AppEnv)),
		})

		if err == nil {
			sqlDB, err := db.DB()
			if err != nil {
				log.Fatalf("Error al obtener instancia DB: %v", err)
			}

			// Configurar pool de conexiones
			sqlDB.SetMaxIdleConns(10)
			sqlDB.SetMaxOpenConns(100)
			sqlDB.SetConnMaxLifetime(time.Hour)

			// Verificar conectividad
			if err := sqlDB.Ping(); err != nil {
				log.Printf("Error al hacer ping a la base de datos: %v", err)
				time.Sleep(retryDelay)
				continue
			}

			log.Println("Conexión a la base de datos establecida correctamente")
			return db
		}

		log.Printf("Error de conexión: %v. Reintentando en %v...", err, retryDelay)
		time.Sleep(retryDelay)
	}

	log.Fatalf("No se pudo conectar a la base de datos después de %d intentos: %v", maxRetries, err)
	return nil
}

// MigrateDatabase ejecuta las migraciones de la base de datos
func MigrateDatabase(db *gorm.DB) {
	log.Println("Iniciando migración de base de datos...")

	// Deshabilitar foreign key checks temporalmente
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		log.Printf("Error al deshabilitar FOREIGN_KEY_CHECKS: %v", err)
	}

	// Eliminar constraints existentes que pueden causar problemas
	dropExistingConstraints(db)

	// Migrar modelos
	err := db.AutoMigrate(
		&models.Course{},
		&models.Chapter{},
		&models.Category{},
		&models.UserProgress{},
		&models.ChapterProgress{},
	)

	if err != nil {
		log.Fatalf("Error al migrar modelos: %v", err)
	}

	// Crear constraints personalizados
	createCustomConstraints(db)

	// Crear índices adicionales
	createAdditionalIndexes(db)

	// Rehabilitar foreign key checks
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error; err != nil {
		log.Printf("Error al habilitar FOREIGN_KEY_CHECKS: %v", err)
	}

	log.Println("Migración de base de datos completada exitosamente")
}

// dropExistingConstraints elimina constraints existentes que pueden causar problemas
func dropExistingConstraints(db *gorm.DB) {
	constraints := []string{
		"fk_cursos_instructor", "fk_courses_instructor",
		"fk_capitulos_curso", "fk_chapters_course",
		"fk_progreso_usuario", "fk_user_progress_user",
		"fk_progreso_curso", "fk_user_progress_course",
		"fk_progreso_capitulo_usuario", "fk_chapter_progress_user",
		"fk_progreso_capitulo_curso", "fk_chapter_progress_course",
		"fk_progreso_capitulo_capitulo", "fk_chapter_progress_chapter",
	}

	tables := []string{"cursos", "courses", "capitulos", "chapters", 
		"progreso_usuarios", "user_progress", "progreso_capitulos", "chapter_progress"}

	for _, table := range tables {
		for _, constraint := range constraints {
			query := fmt.Sprintf("ALTER TABLE %s DROP FOREIGN KEY IF EXISTS %s", table, constraint)
			if err := db.Exec(query).Error; err != nil {
				// No es crítico si falla
				log.Printf("Advertencia: No se pudo eliminar constraint %s de tabla %s: %v", constraint, table, err)
			}
		}
	}
}

// createCustomConstraints crea constraints personalizados
func createCustomConstraints(db *gorm.DB) {
	// Foreign keys para chapters
	if err := db.Exec(`
		ALTER TABLE chapters 
		ADD CONSTRAINT fk_chapters_course 
		FOREIGN KEY (curso_id) REFERENCES cursos(id) ON DELETE CASCADE
	`).Error; err != nil {
		log.Printf("Advertencia: No se pudo crear constraint fk_chapters_course: %v", err)
	}

	// Foreign keys para user_progress
	if err := db.Exec(`
		ALTER TABLE progreso_usuarios 
		ADD CONSTRAINT fk_user_progress_course 
		FOREIGN KEY (curso_id) REFERENCES cursos(id) ON DELETE CASCADE
	`).Error; err != nil {
		log.Printf("Advertencia: No se pudo crear constraint fk_user_progress_course: %v", err)
	}

	// Foreign keys para chapter_progress
	if err := db.Exec(`
		ALTER TABLE progreso_capitulos 
		ADD CONSTRAINT fk_chapter_progress_course 
		FOREIGN KEY (curso_id) REFERENCES cursos(id) ON DELETE CASCADE
	`).Error; err != nil {
		log.Printf("Advertencia: No se pudo crear constraint fk_chapter_progress_course: %v", err)
	}

	if err := db.Exec(`
		ALTER TABLE progreso_capitulos 
		ADD CONSTRAINT fk_chapter_progress_chapter 
		FOREIGN KEY (capitulo_id) REFERENCES capitulos(id) ON DELETE CASCADE
	`).Error; err != nil {
		log.Printf("Advertencia: No se pudo crear constraint fk_chapter_progress_chapter: %v", err)
	}
}

// createAdditionalIndexes crea índices adicionales para optimizar rendimiento
func createAdditionalIndexes(db *gorm.DB) {
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_usuario_curso ON progreso_usuarios(usuario_id, curso_id)",
		"CREATE INDEX IF NOT EXISTS idx_usuario_capitulo ON progreso_capitulos(usuario_id, curso_id, capitulo_id)",
		"CREATE INDEX IF NOT EXISTS idx_curso_publicado ON capitulos(curso_id, publicado)",
		"CREATE INDEX IF NOT EXISTS idx_curso_estado ON cursos(estado)",
		"CREATE INDEX IF NOT EXISTS idx_capitulo_orden ON capitulos(curso_id, orden)",
		"CREATE INDEX IF NOT EXISTS idx_progreso_completado ON progreso_capitulos(usuario_id, completado)",
	}

	for _, indexQuery := range indexes {
		if err := db.Exec(indexQuery).Error; err != nil {
			log.Printf("Advertencia: No se pudo crear índice: %v", err)
		}
	}
}

// getLogLevel retorna el nivel de log apropiado según el ambiente
func getLogLevel(env string) logger.LogLevel {
	switch env {
	case "production":
		return logger.Error
	case "staging":
		return logger.Warn
	default:
		return logger.Info
	}
}

// InitStaticDirs inicializa los directorios estáticos necesarios
func InitStaticDirs(cfg *config.Config) {
	dirs := []string{
		cfg.UploadPath,
		cfg.UploadPath + "/videos",
		cfg.UploadPath + "/images", 
		cfg.UploadPath + "/thumbnails",
	}

	for _, dir := range dirs {
		CreateDirIfNotExists(dir)
	}

	log.Printf("Directorios estáticos inicializados en: %s", cfg.UploadPath)
}