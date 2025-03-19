package api

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/tuusuario/proyecto-cursos/internal/config"
)

// SetupRouter configura todas las rutas de la API
func SetupRouter(db *sql.DB) http.Handler {
	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// Crear router
	r := mux.NewRouter()

	// Middleware
	r.Use(loggingMiddleware)

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// Rutas de autenticación (públicas)
	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", registerHandler(db)).Methods("POST")
	auth.HandleFunc("/login", loginHandler(db, cfg)).Methods("POST")

	// Rutas de cursos
	cursos := api.PathPrefix("/cursos").Subrouter()
	cursos.HandleFunc("", getCursosHandler(db)).Methods("GET")
	cursos.HandleFunc("/{id:[0-9]+}", getCursoHandler(db)).Methods("GET")

	// Rutas protegidas (requieren autenticación)
	protected := api.PathPrefix("/user").Subrouter()
	protected.Use(authMiddleware(cfg))
	protected.HandleFunc("/profile", getUserProfileHandler(db)).Methods("GET")
	protected.HandleFunc("/cursos", getUserCursosHandler(db)).Methods("GET")
	protected.HandleFunc("/inscribir/{id:[0-9]+}", inscribirCursoHandler(db)).Methods("POST")

	// Configurar CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	return c.Handler(r)
}