package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/tuusuario/proyecto-cursos/internal/config"
	"github.com/tuusuario/proyecto-cursos/internal/models"
	"github.com/tuusuario/proyecto-cursos/internal/utils"
)

// registerRequest representa los datos de registro
type registerRequest struct {
	Nombre   string `json:"nombre"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// loginRequest representa los datos de inicio de sesión
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// authResponse representa la respuesta de autenticación
type authResponse struct {
	Token  string         `json:"token"`
	User   *models.Usuario `json:"user"`
	Status string         `json:"status"`
}

// errorResponse representa una respuesta de error
type errorResponse struct {
	Message string `json:"message"`
}

// registerHandler maneja el registro de nuevos usuarios
func registerHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decodificar solicitud
		var req registerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "formato de solicitud inválido")
			return
		}

		// Validar datos
		if err := utils.ValidateNombre(req.Nombre); err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err := utils.ValidateEmail(req.Email); err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err := utils.ValidatePassword(req.Password); err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Crear usuario
		usuario, err := models.CreateUsuario(db, req.Nombre, req.Email, req.Password)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Responder con éxito
		respondWithJSON(w, http.StatusCreated, map[string]interface{}{
			"status":  "success",
			"message": "Usuario registrado correctamente",
			"user": map[string]interface{}{
				"id":     usuario.ID,
				"nombre": usuario.Nombre,
				"email":  usuario.Email,
			},
		})
	}
}

// loginHandler maneja el inicio de sesión
func loginHandler(db *sql.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decodificar solicitud
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "formato de solicitud inválido")
			return
		}

		// Buscar usuario
		usuario, err := models.GetUsuarioByEmail(db, req.Email)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "email o contraseña incorrectos")
			return
		}

		// Verificar contraseña
		if !usuario.CheckPassword(req.Password) {
			respondWithError(w, http.StatusUnauthorized, "email o contraseña incorrectos")
			return
		}

		// Generar token
		token, err := utils.GenerateToken(usuario.ID, cfg)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error al generar token")
			return
		}

		// Datos de usuario para respuesta (sin contraseña)
		userResponse := &models.Usuario{
			ID:        usuario.ID,
			Nombre:    usuario.Nombre,
			Email:     usuario.Email,
			CreatedAt: usuario.CreatedAt,
			UpdatedAt: usuario.UpdatedAt,
		}

		// Responder con token
		respondWithJSON(w, http.StatusOK, &authResponse{
			Token:  token,
			User:   userResponse,
			Status: "success",
		})
	}
}

// getUserProfileHandler obtiene el perfil del usuario autenticado
func getUserProfileHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Obtener ID de usuario del contexto
		userID := getUserIDFromContext(r)
		if userID == 0 {
			respondWithError(w, http.StatusUnauthorized, "no autorizado")
			return
		}

		// Buscar usuario
		usuario, err := models.GetUsuarioByID(db, userID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "usuario no encontrado")
			return
		}

		// Responder con datos del usuario
		respondWithJSON(w, http.StatusOK, usuario)
	}
}

// Función auxiliar para responder con error
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, &errorResponse{Message: message})
}

// Función auxiliar para responder con JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}