package models

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Usuario representa un usuario en el sistema
type Usuario struct {
	ID        int       `json:"id"`
	Nombre    string    `json:"nombre"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // No se envía en respuestas JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUsuario crea un nuevo usuario en la base de datos
func CreateUsuario(db *sql.DB, nombre, email, password string) (*Usuario, error) {
	// Verificar si el email ya existe
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM usuarios WHERE email = ?", email).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("el email ya está registrado")
	}

	// Encriptar contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Insertar usuario
	result, err := db.Exec(
		"INSERT INTO usuarios (nombre, email, password) VALUES (?, ?, ?)",
		nombre, email, string(hashedPassword),
	)
	if err != nil {
		return nil, err
	}

	// Obtener ID generado
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Devolver usuario creado
	return &Usuario{
		ID:        int(id),
		Nombre:    nombre,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// GetUsuarioByEmail busca un usuario por su email
func GetUsuarioByEmail(db *sql.DB, email string) (*Usuario, error) {
	usuario := &Usuario{}
	err := db.QueryRow(
		"SELECT id, nombre, email, password, created_at, updated_at FROM usuarios WHERE email = ?",
		email,
	).Scan(&usuario.ID, &usuario.Nombre, &usuario.Email, &usuario.Password, &usuario.CreatedAt, &usuario.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}
	return usuario, nil
}

// GetUsuarioByID busca un usuario por su ID
func GetUsuarioByID(db *sql.DB, id int) (*Usuario, error) {
	usuario := &Usuario{}
	err := db.QueryRow(
		"SELECT id, nombre, email, created_at, updated_at FROM usuarios WHERE id = ?",
		id,
	).Scan(&usuario.ID, &usuario.Nombre, &usuario.Email, &usuario.CreatedAt, &usuario.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}
	return usuario, nil
}

// CheckPassword verifica si la contraseña es correcta
func (u *Usuario) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}