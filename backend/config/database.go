package config

import (
	"curso-platform/models"
	"curso-platform/utils"
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB es la instancia global de la base de datos
var DB *gorm.DB

// SetupDatabase configura la conexión a la base de datos y la devuelve
func SetupDatabase() *gorm.DB {
	var err error
	maxRetries := 5
	retryDelay := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		log.Printf("Conectando a MariaDB (intento %d/%d)...", i+1, maxRetries)

		// Obtener nombre de la base de datos del entorno o usar el predeterminado "gym_db"
		dbName := utils.GetEnv("DB_NAME", "gym_db")

		// DSN para MariaDB (similar a MySQL)
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			utils.GetEnv("DB_USER", "root"),
			utils.GetEnv("DB_PASSWORD", ""),
			utils.GetEnv("DB_HOST", "localhost"),
			utils.GetEnv("DB_PORT", "3306"),
			dbName)

		// Utilizamos el driver de MySQL, que también es compatible con MariaDB
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err == nil {
			sqlDB, err := DB.DB()
			if err != nil {
				log.Fatalf("Error al obtener instancia DB: %v", err)
			}

			// Configuración de conexiones
			sqlDB.SetMaxIdleConns(10)
			sqlDB.SetMaxOpenConns(100)
			sqlDB.SetConnMaxLifetime(time.Hour)

			log.Println("Conexión a MariaDB establecida correctamente")
			return DB
		}

		log.Printf("Error de conexión: %v. Reintentando en %v...", err, retryDelay)
		time.Sleep(retryDelay)
	}

	log.Fatalf("No se pudo conectar a MariaDB después de %d intentos: %v", maxRetries, err)
	return nil // Nunca se ejecuta debido al Fatalf anterior, pero necesario para compilar
}

// RunMigrations ejecuta las migraciones automáticas para los modelos
func RunMigrations(db *gorm.DB) error {
	log.Println("Iniciando migración de la base de datos...")

	// Verificar si la base de datos existe, si no, crear
	dbName := utils.GetEnv("DB_NAME", "gym_db")
	if err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName)).Error; err != nil {
		return fmt.Errorf("error al crear la base de datos: %v", err)
	}

	// Desactivar temporalmente las restricciones de clave foránea
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		return fmt.Errorf("error al deshabilitar FOREIGN_KEY_CHECKS: %v", err)
	}

	// Lista de nombres de restricciones que podrían existir
	constraintNames := []string{
		"fk_pagos_usuario", "fk_pago_usuario",
		"fk_pagos_curso", "fk_pago_curso",
		"fk_capitulos_curso", "fk_capitulo_curso",
		"fk_progreso_usuario", "fk_progreso_curso",
		"fk_progreso_capitulo_usuario", "fk_progreso_capitulo_curso", "fk_progreso_capitulo_capitulo",
	}

	// Intentar eliminar las restricciones existentes para evitar errores
	for _, constraint := range constraintNames {
		if db.Migrator().HasConstraint(&models.Pago{}, constraint) {
			if err := db.Migrator().DropConstraint(&models.Pago{}, constraint); err != nil {
				log.Printf("Advertencia: No se pudo eliminar constraint %s: %v", constraint, err)
			}
		}

		if db.Migrator().HasConstraint(&models.Capitulo{}, constraint) {
			if err := db.Migrator().DropConstraint(&models.Capitulo{}, constraint); err != nil {
				log.Printf("Advertencia: No se pudo eliminar constraint %s: %v", constraint, err)
			}
		}

		if db.Migrator().HasConstraint(&models.ProgresoUsuario{}, constraint) {
			if err := db.Migrator().DropConstraint(&models.ProgresoUsuario{}, constraint); err != nil {
				log.Printf("Advertencia: No se pudo eliminar constraint %s: %v", constraint, err)
			}
		}

		if db.Migrator().HasConstraint(&models.ProgresoCapitulo{}, constraint) {
			if err := db.Migrator().DropConstraint(&models.ProgresoCapitulo{}, constraint); err != nil {
				log.Printf("Advertencia: No se pudo eliminar constraint %s: %v", constraint, err)
			}
		}
	}

	// Intentar eliminar las restricciones con SQL directo
	tablesWithConstraints := []string{"pagos", "capitulos", "progreso_usuarios", "progreso_capitulos"}
	for _, table := range tablesWithConstraints {
		for _, constraint := range constraintNames {
			if err := db.Exec(fmt.Sprintf("ALTER TABLE %s DROP FOREIGN KEY IF EXISTS %s", table, constraint)).Error; err != nil {
				log.Printf("Advertencia: No se pudo eliminar constraint %s de %s: %v", constraint, table, err)
			}
		}
	}

	// Verificar y migrar la tabla de usuarios
	log.Println("Verificando y migrando tabla de usuarios...")

	// Verificar si la tabla usuarios existe
	if !db.Migrator().HasTable("usuarios") {
		log.Println("Tabla 'usuarios' no encontrada, creándola...")
		if err := db.Migrator().CreateTable(&models.Usuario{}); err != nil {
			return fmt.Errorf("error al crear tabla 'usuarios': %v", err)
		}
		log.Println("Tabla 'usuarios' creada correctamente")
	} else {
		// La tabla existe, verificar si necesitamos migrar de nombre completo a nombre/apellido
		if err := migrateNombreToNombreApellido(db); err != nil {
			log.Printf("Error durante migración de nombres: %v", err)
			return err
		}
	}

	// Verificar si la tabla 'password_resets' existe, y crearla si no
	if !db.Migrator().HasTable("password_resets") {
		log.Println("Tabla 'password_resets' no encontrada, creándola...")
		if err := db.Migrator().CreateTable(&models.PasswordReset{}); err != nil {
			return fmt.Errorf("error al crear tabla 'password_resets': %v", err)
		}
		log.Println("Tabla 'password_resets' creada correctamente")
	}

	// Verificar si la tabla 'activities' existe, y crearla si no
	if !db.Migrator().HasTable("activities") {
		log.Println("Tabla 'activities' no encontrada, creándola...")
		if err := db.Migrator().CreateTable(&models.Activity{}); err != nil {
			return fmt.Errorf("error al crear tabla 'activities': %v", err)
		}
		log.Println("Tabla 'activities' creada correctamente")
	}

	// Verificar si la tabla 'enrollments' existe, y crearla si no
	if !db.Migrator().HasTable("enrollments") {
		log.Println("Tabla 'enrollments' no encontrada, creándola...")
		if err := db.Migrator().CreateTable(&models.Enrollment{}); err != nil {
			return fmt.Errorf("error al crear tabla 'enrollments': %v", err)
		}
		log.Println("Tabla 'enrollments' creada correctamente")
	}

	// Realizar migración automática de todos los modelos
	log.Println("Ejecutando auto-migración de modelos...")
	if err := db.AutoMigrate(
		&models.Usuario{},
		//&models.Curso{}, //
		//&models.Capitulo{}, //
		//&models.Pago{}, //
		//&models.ProgresoUsuario{}, //
		//&models.ProgresoCapitulo{}, //
		&models.ActivityLog{},
		&models.ContactMessage{},
		//&models.ProjectPortfolio{},
		//&models.HomeImage{},
		&models.PasswordReset{},
		&models.Activity{},
		&models.Enrollment{},
	); err != nil {
		return fmt.Errorf("error al migrar tablas base: %v", err)
	}

	// Crear tablas si no existen
	/*tables := []interface{}{&models.Pago{}, &models.ProgresoUsuario{}, &models.ProgresoCapitulo{}}
	for _, table := range tables {
		if !db.Migrator().HasTable(table) {
			if err := db.Migrator().CreateTable(table); err != nil {
				return fmt.Errorf("error al crear tabla %T: %v", table, err)
			}
		}
	}*/

	// Modificar columnas específicas si es necesario
	if err := db.Exec("ALTER TABLE pagos MODIFY COLUMN estado varchar(20) NOT NULL DEFAULT 'pendiente'").Error; err != nil {
		log.Printf("Advertencia: No se pudo modificar columna estado: %v", err)
	}

	// Crear restricciones de clave externa
	log.Println("Creando restricciones de clave externa...")
	constraints := []struct {
		table      string
		constraint string
		definition string
	}{
		{"pagos", "fk_pagos_usuario", "FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE"},
		{"pagos", "fk_pagos_curso", "FOREIGN KEY (curso_id) REFERENCES cursos(id) ON DELETE CASCADE"},
		{"capitulos", "fk_capitulos_curso", "FOREIGN KEY (curso_id) REFERENCES cursos(id) ON DELETE CASCADE"},
		{"progreso_usuarios", "fk_progreso_usuario", "FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE"},
		{"progreso_usuarios", "fk_progreso_curso", "FOREIGN KEY (curso_id) REFERENCES cursos(id) ON DELETE CASCADE"},
		{"progreso_capitulos", "fk_progreso_capitulo_usuario", "FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE"},
		{"progreso_capitulos", "fk_progreso_capitulo_curso", "FOREIGN KEY (curso_id) REFERENCES cursos(id) ON DELETE CASCADE"},
		{"progreso_capitulos", "fk_progreso_capitulo_capitulo", "FOREIGN KEY (capitulo_id) REFERENCES capitulos(id) ON DELETE CASCADE"},
	}

	for _, c := range constraints {
		if err := db.Exec(fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s %s", c.table, c.constraint, c.definition)).Error; err != nil {
			log.Printf("Advertencia: No se pudo crear constraint %s: %v", c.constraint, err)
		}
	}

	// Crear índices
	log.Println("Creando índices...")
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_usuario_curso ON progreso_usuarios(usuario_id, curso_id)").Error; err != nil {
		log.Printf("Advertencia: No se pudo crear índice idx_usuario_curso: %v", err)
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_usuario_capitulo ON progreso_capitulos(usuario_id, curso_id, capitulo_id)").Error; err != nil {
		log.Printf("Advertencia: No se pudo crear índice idx_usuario_capitulo: %v", err)
	}

	// Crear índice para los tokens de restablecimiento de contraseña
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_password_resets_token ON password_resets(token)").Error; err != nil {
		log.Printf("Advertencia: No se pudo crear índice idx_password_resets_token: %v", err)
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_password_resets_email ON password_resets(email)").Error; err != nil {
		log.Printf("Advertencia: No se pudo crear índice idx_password_resets_email: %v", err)
	}

	// Volver a habilitar las comprobaciones de clave externa
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error; err != nil {
		return fmt.Errorf("error al habilitar FOREIGN_KEY_CHECKS: %v", err)
	}

	log.Println("Migración de la base de datos completada con éxito")
	return nil
}

// migrateNombreToNombreApellido migra los datos existentes de nombre completo a nombre y apellido separados
func migrateNombreToNombreApellido(db *gorm.DB) error {
	log.Println("Verificando si es necesario migrar nombres...")

	// Verificar si la columna apellido ya existe
	if db.Migrator().HasColumn(&models.Usuario{}, "apellido") {
		log.Println("La columna 'apellido' ya existe, no es necesario migrar")
		return nil
	}

	log.Println("Iniciando migración de nombres completos a nombre y apellido...")

	// Primero, agregar la nueva columna apellido
	if err := db.Migrator().AddColumn(&models.Usuario{}, "apellido"); err != nil {
		return fmt.Errorf("error al agregar columna apellido: %v", err)
	}

	// Obtener todos los usuarios existentes
	var usuarios []models.Usuario
	if err := db.Find(&usuarios).Error; err != nil {
		return fmt.Errorf("error al obtener usuarios existentes: %v", err)
	}

	// Migrar cada usuario
	for _, usuario := range usuarios {
		if usuario.Nombre == "" {
			continue // Saltar usuarios sin nombre
		}

		// Dividir el nombre completo
		partes := strings.Fields(strings.TrimSpace(usuario.Nombre))

		var nombre, apellido string
		if len(partes) == 0 {
			// Si no hay partes, usar valores por defecto
			nombre = "Usuario"
			apellido = "Sin Apellido"
		} else if len(partes) == 1 {
			// Solo un nombre
			nombre = partes[0]
			apellido = "Sin Apellido"
		} else {
			// Múltiples partes: primera parte es nombre, resto es apellido
			nombre = partes[0]
			apellido = strings.Join(partes[1:], " ")
		}

		// Actualizar el usuario
		if err := db.Model(&usuario).Updates(map[string]interface{}{
			"nombre":   nombre,
			"apellido": apellido,
		}).Error; err != nil {
			log.Printf("Error al migrar usuario ID %d: %v", usuario.ID, err)
			continue
		}

		log.Printf("Migrado usuario ID %d: '%s' -> Nombre: '%s', Apellido: '%s'",
			usuario.ID, usuario.Nombre, nombre, apellido)
	}

	// Actualizar el tamaño de la columna nombre
	if err := db.Exec("ALTER TABLE usuarios MODIFY COLUMN nombre varchar(50) NOT NULL").Error; err != nil {
		log.Printf("Advertencia: No se pudo modificar el tamaño de la columna nombre: %v", err)
	}

	// Hacer que la columna apellido sea NOT NULL
	if err := db.Exec("ALTER TABLE usuarios MODIFY COLUMN apellido varchar(50) NOT NULL").Error; err != nil {
		log.Printf("Advertencia: No se pudo hacer NOT NULL la columna apellido: %v", err)
	}

	log.Printf("Migración de nombres completada. Se migraron %d usuarios", len(usuarios))
	return nil
}
