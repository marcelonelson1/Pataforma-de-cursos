package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/tuusuario/proyecto-cursos/internal/config"
	"github.com/tuusuario/proyecto-cursos/internal/utils"
)

// loggingMiddleware registra las solicitudes entrantes
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Inicio: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Fin: %s %s (%v)", r.Method, r.URL.Path, time.Since(start))
	})
}

// authMiddleware verifica la autenticación del usuario
func authMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Obtener token del header
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				respondWithError(w, http.StatusUnauthorized, "Token de autorización requerido")
				return
			}

			// Validar token
			claims, err := utils.ValidateToken(tokenString, cfg)
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "Token inválido")
				return
			}

			// Añadir ID de usuario al contexto
			ctx := context.WithValue(r.Context(), "userID", claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// getUserIDFromContext obtiene el ID de usuario del contexto
func getUserIDFromContext(r *http.Request) int {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		return 0
	}
	return userID
}