package database

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Init inicializa la conexión a la base de datos
func Init(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Configurar pool de conexiones
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verificar conexión
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Ejecutar migraciones
	if err := runMigrations(db); err != nil {
		return nil, err
	}

	return db, nil
}

// runMigrations ejecuta las migraciones de la base de datos
func runMigrations(db *sql.DB) error {
	// Crear tabla de usuarios si no existe
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS usuarios (
			id INT AUTO_INCREMENT PRIMARY KEY,
			nombre VARCHAR(100) NOT NULL,
			email VARCHAR(100) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// Crear tabla de cursos si no existe
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS cursos (
			id INT AUTO_INCREMENT PRIMARY KEY,
			titulo VARCHAR(255) NOT NULL,
			descripcion TEXT,
			imagen VARCHAR(255),
			precio DECIMAL(10,2) NOT NULL,
			duracion VARCHAR(50),
			nivel VARCHAR(50),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// Crear tabla de inscripciones si no existe
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS inscripciones (
			id INT AUTO_INCREMENT PRIMARY KEY,
			usuario_id INT NOT NULL,
			curso_id INT NOT NULL,
			fecha_inscripcion TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (usuario_id) REFERENCES usuarios(id),
			FOREIGN KEY (curso_id) REFERENCES cursos(id),
			UNIQUE KEY usuario_curso (usuario_id, curso_id)
		)
	`)
	return err
}