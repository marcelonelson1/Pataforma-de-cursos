package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	// Configuración inicial
	setupLogging()
	loadEnv()

	// Base de datos
	if err := setupDatabase(); err != nil {
		log.Fatalf("Error al configurar la base de datos: %v", err)
	}

	// Router
	router := setupRouter()
	registerRoutes(router)

	// Iniciar servidor
	startServer(router)
}

func setupLogging() {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(logFile)
	}
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env, usando variables de entorno del sistema")
	}
}

func setupDatabase() error {
	var err error
	maxRetries := 5

	for i := 0; i < maxRetries; i++ {
		log.Printf("Conectando a la base de datos (intento %d/%d)...", i+1, maxRetries)

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			getEnv("DB_USER", "root"),
			getEnv("DB_PASSWORD", ""),
			getEnv("DB_HOST", "localhost"),
			getEnv("DB_PORT", "3306"),
			getEnv("DB_NAME", "cursos_db"))

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err == nil {
			sqlDB, err := db.DB()
			if err != nil {
				return err
			}

			sqlDB.SetMaxIdleConns(10)
			sqlDB.SetMaxOpenConns(100)
			sqlDB.SetConnMaxLifetime(time.Hour)

			return db.AutoMigrate(&Usuario{}, &Curso{}, &PasswordReset{})
		}

		time.Sleep(5 * time.Second)
	}

	return err
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{getEnv("FRONTEND_URL", "http://localhost:3000")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(errorHandler())

	return router
}

func errorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			log.Printf("Error en la solicitud: %v", c.Errors)
		}
	}
}

func registerRoutes(router *gin.Engine) {
	// Auth
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", register)
		auth.POST("/login", login)
		auth.POST("/forgot-password", forgotPassword)
		auth.GET("/reset-password/:token/validate", validateResetToken)
		auth.POST("/reset-password", resetPassword)
	}

	// Cursos
	cursos := router.Group("/api/cursos")
	{
		cursos.GET("", getCursos)
		cursos.GET("/:id", getCursoById)
		cursos.POST("", authMiddleware(), createCurso)
		cursos.PUT("/:id", authMiddleware(), updateCurso)
		cursos.DELETE("/:id", authMiddleware(), deleteCurso)
	}

	// Contacto
	router.POST("/api/contact", contactHandler)

	// Health Check
	router.GET("/api/health", healthCheck)
}

func healthCheck(c *gin.Context) {
	sqlDB, err := db.DB()
	if err != nil {
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}

	if err := sqlDB.Ping(); err != nil {
		SendErrorResponse(c, ErrDatabaseError, http.StatusServiceUnavailable)
		return
	}

	SendSuccessResponse(c, gin.H{"status": "ok"})
}

func startServer(router *gin.Engine) {
	port := getEnv("PORT", "5000")
	log.Printf("Iniciando servidor en puerto %s", port)
	log.Fatal(router.Run(":" + port))
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}