package models

import (
	"database/sql"
	"errors"
	"time"
)

// Curso representa un curso en la plataforma
type Curso struct {
	ID          int       `json:"id"`
	Titulo      string    `json:"titulo"`
	Descripcion string    `json:"descripcion"`
	Imagen      string    `json:"imagen"`
	Precio      float64   `json:"precio"`
	Duracion    string    `json:"duracion"`
	Nivel       string    `json:"nivel"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GetCursos obtiene todos los cursos disponibles
func GetCursos(db *sql.DB) ([]*Curso, error) {
	rows, err := db.Query("SELECT id, titulo, descripcion, imagen, precio, duracion, nivel, created_at, updated_at FROM cursos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cursos []*Curso
	for rows.Next() {
		curso := &Curso{}
		err := rows.Scan(
			&curso.ID,
			&curso.Titulo,
			&curso.Descripcion,
			&curso.Imagen,
			&curso.Precio,
			&curso.Duracion,
			&curso.Nivel,
			&curso.CreatedAt,
			&curso.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		cursos = append(cursos, curso)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cursos, nil
}

// GetCursoByID obtiene un curso por su ID
func GetCursoByID(db *sql.DB, id int) (*Curso, error) {
	curso := &Curso{}
	err := db.QueryRow(
		"SELECT id, titulo, descripcion, imagen, precio, duracion, nivel, created_at, updated_at FROM cursos WHERE id = ?",
		id,
	).Scan(
		&curso.ID,
		&curso.Titulo,
		&curso.Descripcion,
		&curso.Imagen,
		&curso.Precio,
		&curso.Duracion,
		&curso.Nivel,
		&curso.CreatedAt,
		&curso.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("curso no encontrado")
		}
		return nil, err
	}
	return curso, nil
}

// CreateCurso crea un nuevo curso
func CreateCurso(db *sql.DB, curso *Curso) (*Curso, error) {
	result, err := db.Exec(
		"INSERT INTO cursos (titulo, descripcion, imagen, precio, duracion, nivel) VALUES (?, ?, ?, ?, ?, ?)",
		curso.Titulo, curso.Descripcion, curso.Imagen, curso.Precio, curso.Duracion, curso.Nivel,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	curso.ID = int(id)
	curso.CreatedAt = time.Now()
	curso.UpdatedAt = time.Now()
	return curso, nil
}

// InscribirUsuarioEnCurso inscribe a un usuario en un curso
func InscribirUsuarioEnCurso(db *sql.DB, usuarioID, cursoID int) error {
	_, err := db.Exec(
		"INSERT INTO inscripciones (usuario_id, curso_id) VALUES (?, ?)",
		usuarioID, cursoID,
	)
	return err
}

// GetCursosDeUsuario obtiene todos los cursos en los que está inscrito un usuario
func GetCursosDeUsuario(db *sql.DB, usuarioID int) ([]*Curso, error) {
	rows, err := db.Query(`
		SELECT c.id, c.titulo, c.descripcion, c.imagen, c.precio, c.duracion, c.nivel, c.created_at, c.updated_at 
		FROM cursos c
		JOIN inscripciones i ON c.id = i.curso_id
		WHERE i.usuario_id = ?
	`, usuarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cursos []*Curso
	for rows.Next() {
		curso := &Curso{}
		err := rows.Scan(
			&curso.ID,
			&curso.Titulo,
			&curso.Descripcion,
			&curso.Imagen,
			&curso.Precio,
			&curso.Duracion,
			&curso.Nivel,
			&curso.CreatedAt,
			&curso.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		cursos = append(cursos, curso)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cursos, nil
}