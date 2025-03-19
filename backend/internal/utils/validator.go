package utils

import (
	"errors"
	"net/mail"
	"regexp"
	"unicode"
)

// ValidateEmail valida el formato de un email
func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("email inválido")
	}
	return nil
}

// ValidatePassword verifica que la contraseña cumpla con los requisitos mínimos
func ValidatePassword(password string) error {
	if len(password) < 6 {
		return errors.New("la contraseña debe tener al menos 6 caracteres")
	}

	// Verificar si contiene al menos una letra
	hasLetter := false
	for _, char := range password {
		if unicode.IsLetter(char) {
			hasLetter = true
			break
		}
	}
	if !hasLetter {
		return errors.New("la contraseña debe contener al menos una letra")
	}

	// Verificar si contiene al menos un número
	hasNumber := false
	for _, char := range password {
		if unicode.IsNumber(char) {
			hasNumber = true
			break
		}
	}
	if !hasNumber {
		return errors.New("la contraseña debe contener al menos un número")
	}

	return nil
}

// ValidateNombre verifica que el nombre sea válido
func ValidateNombre(nombre string) error {
	if len(nombre) < 2 {
		return errors.New("el nombre debe tener al menos 2 caracteres")
	}

	// Permitir solo letras, espacios y algunos caracteres especiales
	r := regexp.MustCompile(`^[a-zA-ZáéíóúÁÉÍÓÚñÑ\s'.-]+$`)
	if !r.MatchString(nombre) {
		return errors.New("el nombre contiene caracteres no válidos")
	}

	return nil
}