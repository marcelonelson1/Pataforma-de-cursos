package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Modelos
type Usuario struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Nombre    string    `gorm:"size:100;not null" json:"nombre"`
	Email     string    `gorm:"size:100;not null;uniqueIndex" json:"email"`
	Password  string    `gorm:"size:100;not null" json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Curso struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Titulo      string    `gorm:"size:200;not null" json:"titulo"`
	Descripcion string    `gorm:"size:500" json:"descripcion"`
	Contenido   string    `gorm:"type:text" json:"contenido"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

var db *gorm.DB

func main() {
	// Cargar variables de entorno
	err := godotenv.Load()
	if err != nil {
		log.Println("Archivo .env no encontrado, usando variables de entorno del sistema")
	}

	// Configurar conexión a base de datos
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		getEnv("DB_USER", "root"),
		getEnv("DB_PASSWORD", ""),
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "3306"),
		getEnv("DB_NAME", "cursos_db"))

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	log.Println("Conexión a la base de datos establecida")

	// Migraciones automáticas
	err = db.AutoMigrate(&Usuario{}, &Curso{})
	if err != nil {
		log.Fatalf("Error en la migración: %v", err)
	}

	// Inicializar router
	router := gin.Default()

	// Opción 1: Configuración más permisiva para desarrollo
	router.Use(cors.Default()) // Permite todos los orígenes

	// Opción 2: Configuración más específica
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Permite todos los orígenes
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Rutas de autenticación
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", register)
		auth.POST("/login", login)
	}

	// Rutas de cursos
	cursos := router.Group("/api/cursos")
	{
		cursos.GET("", getCursos)
		cursos.GET("/:id", getCursoById)
		cursos.POST("", authMiddleware(), createCurso)
		cursos.PUT("/:id", authMiddleware(), updateCurso)
		cursos.DELETE("/:id", authMiddleware(), deleteCurso)
	}

	// Iniciar servidor
	port := getEnv("PORT", "5000")
	log.Printf("Servidor iniciado en el puerto %s", port)
	log.Fatal(router.Run(":" + port))
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
