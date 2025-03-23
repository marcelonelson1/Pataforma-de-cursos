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
	// Configurar logging
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Error al abrir archivo de log:", err)
	} else {
		log.SetOutput(logFile)
	}

	// Cargar variables de entorno
	err = godotenv.Load()
	if err != nil {
		log.Println("Archivo .env no encontrado, usando variables de entorno del sistema")
	}

	// Configurar conexión a base de datos con reintentos
	maxRetries := 5
	var dbError error

	for i := 0; i < maxRetries; i++ {
		log.Printf("Intentando conectar a la base de datos (intento %d/%d)...", i+1, maxRetries)

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			getEnv("DB_USER", "root"),
			getEnv("DB_PASSWORD", ""),
			getEnv("DB_HOST", "localhost"),
			getEnv("DB_PORT", "3306"),
			getEnv("DB_NAME", "cursos_db"))

		// Configurar logger de GORM
		dbLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Tiempo para considerar una consulta lenta
				LogLevel:                  logger.Warn, // Nivel de log
				IgnoreRecordNotFoundError: false,       // Registrar errores ErrRecordNotFound
				Colorful:                  true,        // Habilitar color
			},
		)

		db, dbError = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: dbLogger,
		})

		if dbError == nil {
			log.Println("Conexión a la base de datos establecida exitosamente")
			break
		}

		log.Printf("Error al conectar a la base de datos: %v. Reintentando en 5 segundos...", dbError)
		time.Sleep(5 * time.Second)
	}

	if dbError != nil {
		log.Fatalf("No se pudo conectar a la base de datos después de %d intentos: %v", maxRetries, dbError)
	}

	// Configuración de pool de conexiones
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error al obtener el objeto *sql.DB: %v", err)
	}

	// Establecer límites de conexión
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Migraciones automáticas
	err = db.AutoMigrate(&Usuario{}, &Curso{}, &PasswordReset{})
	if err != nil {
		log.Fatalf("Error en la migración: %v", err)
	}

	// Inicializar router
	gin.SetMode(getEnv("GIN_MODE", "debug"))
	router := gin.Default()

	// Configuración CORS más permisiva
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Permite solicitudes desde el frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Middleware para manejo global de errores
	router.Use(func(c *gin.Context) {
		c.Next()
		// Manejar errores después de la solicitud
		if len(c.Errors) > 0 {
			log.Printf("Errores en la solicitud: %v", c.Errors)
		}
	})

	// Rutas de autenticación
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", register)
		auth.POST("/login", login)

		// Rutas para recuperación de contraseña
		auth.POST("/forgot-password", forgotPassword)
		auth.GET("/reset-password/:token/validate", validateResetToken)
		auth.POST("/reset-password", resetPassword)
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

	// Ruta de verificación de estado
	router.GET("/api/health", func(c *gin.Context) {
		// Verificar conexión a DB
		sqlDB, err := db.DB()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Error de conexión a DB"})
			return
		}

		if err := sqlDB.Ping(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "No se puede conectar a la base de datos"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Servidor funcionando correctamente",
		})
	})

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
