package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tuusuario/proyecto-cursos/internal/models"
)

// getCursosHandler obtiene todos los cursos
func getCursosHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Obtener cursos
		cursos, err := models.GetCursos(db)
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

// getCursoHandler obtiene un curso específico por ID
func getCursoHandler(db *sql.DB) http.