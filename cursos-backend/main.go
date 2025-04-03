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

// Estructuras de modelos
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
	Precio      float64   `gorm:"type:decimal(10,2);default:29.99" json:"precio"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Pago struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UsuarioID     uint      `gorm:"not null" json:"usuario_id"`
	CursoID       uint      `gorm:"not null" json:"curso_id"`
	Monto         float64   `gorm:"type:decimal(10,2);not null" json:"monto"`
	Metodo        string    `gorm:"size:50;not null" json:"metodo"`
	Estado        string    `gorm:"size:20;not null;default:'pendiente'" json:"estado"`
	TransaccionID string    `gorm:"size:100" json:"transaccion_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

var db *gorm.DB

func main() {
	// Configuración inicial
	setupLogging()
	loadEnv()

	// Configuración de la base de datos
	if err := setupDatabase(); err != nil {
		log.Fatalf("Error al configurar la base de datos: %v", err)
	}

	// Configuración del router
	router := setupRouter()
	registerRoutes(router)

	// Iniciar servidor
	startServer(router)
}

func setupLogging() {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(logFile)
	} else {
		log.Println("No se pudo abrir el archivo de log, usando salida estándar")
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
	retryDelay := 5 * time.Second

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
				return fmt.Errorf("error al obtener instancia DB: %v", err)
			}

			sqlDB.SetMaxIdleConns(10)
			sqlDB.SetMaxOpenConns(100)
			sqlDB.SetConnMaxLifetime(time.Hour)

			if err := migrateWithConstraints(); err != nil {
				return fmt.Errorf("error en migración: %v", err)
			}

			log.Println("Conexión a la base de datos establecida correctamente")
			return nil
		}

		log.Printf("Error de conexión: %v. Reintentando en %v...", err, retryDelay)
		time.Sleep(retryDelay)
	}

	return fmt.Errorf("no se pudo conectar a la base de datos después de %d intentos: %v", maxRetries, err)
}

func migrateWithConstraints() error {
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		return fmt.Errorf("error al deshabilitar FOREIGN_KEY_CHECKS: %v", err)
	}

	constraintNames := []string{"fk_pagos_usuario", "fk_pago_usuario", "fk_pagos_curso", "fk_pago_curso"}
	for _, constraint := range constraintNames {
		if db.Migrator().HasConstraint(&Pago{}, constraint) {
			if err := db.Migrator().DropConstraint(&Pago{}, constraint); err != nil {
				log.Printf("Advertencia: No se pudo eliminar constraint %s: %v", constraint, err)
			}
		}
	}

	if err := db.Exec("ALTER TABLE pagos DROP FOREIGN KEY IF EXISTS fk_pagos_usuario").Error; err != nil {
		log.Printf("Advertencia: No se pudo eliminar constraint fk_pagos_usuario: %v", err)
	}
	if err := db.Exec("ALTER TABLE pagos DROP FOREIGN KEY IF EXISTS fk_pago_usuario").Error; err != nil {
		log.Printf("Advertencia: No se pudo eliminar constraint fk_pago_usuario: %v", err)
	}
	if err := db.Exec("ALTER TABLE pagos DROP FOREIGN KEY IF EXISTS fk_pagos_curso").Error; err != nil {
		log.Printf("Advertencia: No se pudo eliminar constraint fk_pagos_curso: %v", err)
	}
	if err := db.Exec("ALTER TABLE pagos DROP FOREIGN KEY IF EXISTS fk_pago_curso").Error; err != nil {
		log.Printf("Advertencia: No se pudo eliminar constraint fk_pago_curso: %v", err)
	}

	if err := db.AutoMigrate(&Usuario{}, &Curso{}); err != nil {
		return fmt.Errorf("error al migrar tablas base: %v", err)
	}

	if !db.Migrator().HasTable(&Pago{}) {
		if err := db.Migrator().CreateTable(&Pago{}); err != nil {
			return fmt.Errorf("error al crear tabla Pago: %v", err)
		}
	} else {
		if err := db.Exec("ALTER TABLE pagos MODIFY COLUMN estado varchar(20) NOT NULL DEFAULT 'pendiente'").Error; err != nil {
			log.Printf("Advertencia: No se pudo modificar columna estado: %v", err)
		}
	}

	if err := db.Exec(
		"ALTER TABLE pagos ADD CONSTRAINT fk_pagos_usuario FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE",
	).Error; err != nil {
		log.Printf("Advertencia: No se pudo crear constraint fk_pagos_usuario: %v", err)
	}

	if err := db.Exec(
		"ALTER TABLE pagos ADD CONSTRAINT fk_pagos_curso FOREIGN KEY (curso_id) REFERENCES cursos(id) ON DELETE CASCADE",
	).Error; err != nil {
		log.Printf("Advertencia: No se pudo crear constraint fk_pagos_curso: %v", err)
	}

	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error; err != nil {
		return fmt.Errorf("error al habilitar FOREIGN_KEY_CHECKS: %v", err)
	}

	return nil
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	// Obtener la URL del frontend desde variables de entorno
	frontendURL := getEnv("FRONTEND_URL", "http://localhost:3000")
	log.Printf("Configurando CORS para permitir origen: %s", frontendURL)

	// Configuración CORS corregida para permitir credenciales
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(JSONMiddleware())
	router.Use(RecoveryMiddleware())
	router.Use(ErrorHandler())

	return router
}

func registerRoutes(router *gin.Engine) {
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", register)
		auth.POST("/login", login)
		auth.POST("/forgot-password", forgotPassword)
		auth.GET("/reset-password/:token/validate", validateResetToken)
		auth.POST("/reset-password", resetPassword)
	}

	cursos := router.Group("/api/cursos")
	{
		cursos.GET("", getCursos)
		cursos.GET("/:id", getCursoById)
		cursos.POST("", authMiddleware(), createCurso)
		cursos.PUT("/:id", authMiddleware(), updateCurso)
		cursos.DELETE("/:id", authMiddleware(), deleteCurso)
	}

	pagos := router.Group("/api/pagos")
	{
		pagos.Use(authMiddleware())
		pagos.POST("", crearPago)
		pagos.GET("/:id", verificarPagoPorCurso) // Cambiado para coincidir con frontend
		pagos.POST("/webhook", webhookPago)
	}

	router.POST("/api/contact", contactHandler)
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

	SendSuccessResponse(c, gin.H{
		"status":  "ok",
		"version": "1.0.0",
		"time":    time.Now().Format(time.RFC3339),
	})
}

func startServer(router *gin.Engine) {
	port := getEnv("PORT", "5000")
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Iniciando servidor en puerto %s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error al iniciar servidor: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Funciones externas (implementadas en otros archivos)
