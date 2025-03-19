package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tuusuario/proyecto-cursos/internal/models"
)

// getUserCursosHandler obtiene los cursos en los que está inscrito un usuario
func getUserCursosHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Obtener ID de usuario del contexto
		userID := getUserIDFromContext(r)
		if userID == 0 {
			respondWithError(w, http.StatusUnauthorized, "no autorizado")
			return
		}

		// Obtener cursos del usuario
		cursos, err := models.GetCursosDeUsuario(db, userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error al obtener cursos")
			return
		}

		// Responder con cursos
		respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"status": "success",
			"cursos": cursos,
		})
	}
}

// inscribirCursoHandler inscribe a un usuario en un curso
func inscribirCursoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Obtener ID de usuario del contexto
		userID := getUserIDFromContext(r)
		if userID == 0 {
			respondWithError(w, http.StatusUnauthorized, "no autorizado")
			return
		}

		// Obtener ID del curso desde la URL
		vars := mux.Vars(r)
		cursoID, err := strconv.Atoi(vars["id"])
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "ID de curso inválido")
			return
		}

		// Inscribir usuario en el curso
		err = models.InscribirUsuarioEnCurso(db, userID, cursoID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error al inscribirse en el curso")
			return
		}

		// Responder con éxito
		respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"status":  "success",
			"message": "Inscripción exitosa",
		})
	}
}