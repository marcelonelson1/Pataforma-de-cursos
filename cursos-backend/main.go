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
	Role      string    `gorm:"size:20;default:'user'" json:"role"`
	Phone     string    `gorm:"size:20" json:"phone"`
	ImageURL  string    `gorm:"size:255" json:"image_url"`
	LastLogin time.Time `json:"last_login"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Curso struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Titulo      string     `gorm:"size:200;not null" json:"titulo"`
	Descripcion string     `gorm:"size:500" json:"descripcion"`
	Contenido   string     `gorm:"type:text" json:"contenido"`
	Precio      float64    `gorm:"type:decimal(10,2);default:29.99" json:"precio"`
	Estado      string     `gorm:"size:20;default:'Borrador'" json:"estado"`
	ImagenURL   string     `gorm:"size:255" json:"imagen_url"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Capitulos   []Capitulo `gorm:"foreignKey:CursoID" json:"capitulos,omitempty"`
}

type Capitulo struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CursoID     uint      `gorm:"not null;index" json:"curso_id"`
	Titulo      string    `gorm:"size:200;not null" json:"titulo"`
	Descripcion string    `gorm:"size:500" json:"descripcion"`
	Duracion    string    `gorm:"size:10" json:"duracion"`
	VideoURL    string    `gorm:"size:255" json:"video_url"`
	VideoNombre string    `gorm:"size:255" json:"video_nombre"`
	Orden       int       `gorm:"default:0" json:"orden"`
	Publicado   bool      `gorm:"default:false" json:"publicado"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Curso       *Curso    `gorm:"foreignKey:CursoID" json:"-"`
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

type ActivityLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Action    string    `gorm:"size:50;not null" json:"action"`
	Details   string    `gorm:"size:500" json:"details"`
	IP        string    `gorm:"size:50" json:"ip"`
	UserAgent string    `gorm:"size:255" json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

// ContactMessage representa un mensaje enviado a través del formulario de contacto

var db *gorm.DB

func main() {
	setupLogging()
	loadEnv()

	if err := setupDatabase(); err != nil {
		log.Fatalf("Error al configurar la base de datos: %v", err)
	}
	
	// Inicializar directorios necesarios para almacenamiento
	initStaticDirs()

	router := setupRouter()
	registerRoutes(router)

	startServer(router)
}

// initStaticDirs inicializa todos los directorios necesarios para almacenamiento
func initStaticDirs() {
	// Crear directorio principal si no existe
	createDirIfNotExists("./static")
	
	// Inicializar directorios específicos
	initVideosDir()
	initProfilesDir()
	initPortfolioDir()
}

// initVideosDir inicializa el directorio para videos de cursos


// initProfilesDir inicializa el directorio para imágenes de perfil
func initProfilesDir() {
	createDirIfNotExists("./static/profiles")
}

// initPortfolioDir inicializa el directorio para imágenes del portfolio
func initPortfolioDir() {
	createDirIfNotExists("./static/portfolio")
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

	// Eliminar restricciones existentes para evitar errores de duplicación
	constraintNames := []string{
		"fk_pagos_usuario", "fk_pago_usuario", 
		"fk_pagos_curso", "fk_pago_curso",
		"fk_capitulos_curso", "fk_capitulo_curso",
	}
	
	for _, constraint := range constraintNames {
		if db.Migrator().HasConstraint(&Pago{}, constraint) {
			if err := db.Migrator().DropConstraint(&Pago{}, constraint); err != nil {
				log.Printf("Advertencia: No se pudo eliminar constraint %s: %v", constraint, err)
			}
		}
		
		if db.Migrator().HasConstraint(&Capitulo{}, constraint) {
			if err := db.Migrator().DropConstraint(&Capitulo{}, constraint); err != nil {
				log.Printf("Advertencia: No se pudo eliminar constraint %s: %v", constraint, err)
			}
		}
	}

	// Eliminar constraints de tablas específicas para Pagos
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

	// Eliminar constraints de tablas específicas para Capítulos
	if err := db.Exec("ALTER TABLE capitulos DROP FOREIGN KEY IF EXISTS fk_capitulos_curso").Error; err != nil {
		log.Printf("Advertencia: No se pudo eliminar constraint fk_capitulos_curso: %v", err)
	}
	if err := db.Exec("ALTER TABLE capitulos DROP FOREIGN KEY IF EXISTS fk_capitulo_curso").Error; err != nil {
		log.Printf("Advertencia: No se pudo eliminar constraint fk_capitulo_curso: %v", err)
	}

	// Migrar las tablas base
	if err := db.AutoMigrate(&Usuario{}, &Curso{}, &Capitulo{}, &ActivityLog{}, &ContactMessage{}, &ProjectPortfolio{}); err != nil {
		return fmt.Errorf("error al migrar tablas base: %v", err)
	}

	// Verificar y crear tabla de pagos si no existe
	if !db.Migrator().HasTable(&Pago{}) {
		if err := db.Migrator().CreateTable(&Pago{}); err != nil {
			return fmt.Errorf("error al crear tabla Pago: %v", err)
		}
	} else {
		if err := db.Exec("ALTER TABLE pagos MODIFY COLUMN estado varchar(20) NOT NULL DEFAULT 'pendiente'").Error; err != nil {
			log.Printf("Advertencia: No se pudo modificar columna estado: %v", err)
		}
	}

	// Crear o recrear las relaciones entre tablas
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

	if err := db.Exec(
		"ALTER TABLE capitulos ADD CONSTRAINT fk_capitulos_curso FOREIGN KEY (curso_id) REFERENCES cursos(id) ON DELETE CASCADE",
	).Error; err != nil {
		log.Printf("Advertencia: No se pudo crear constraint fk_capitulos_curso: %v", err)
	}

	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error; err != nil {
		return fmt.Errorf("error al habilitar FOREIGN_KEY_CHECKS: %v", err)
	}

	return nil
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	frontendURL := getEnv("FRONTEND_URL", "http://localhost:3000")
	log.Printf("Configurando CORS para permitir origen: %s", frontendURL)

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
		auth.GET("/check-admin", authMiddleware(), checkAdmin)
		
		profile := auth.Group("")
		profile.Use(authMiddleware())
		{
			profile.GET("/profile", getProfile)
			profile.PUT("/profile", updateProfile)
			profile.POST("/change-password", changePassword)
			profile.POST("/profile/image", uploadProfileImage)
			profile.GET("/notification-settings", getNotificationSettings)
			profile.PUT("/notification-settings", updateNotificationSettings)
			profile.POST("/refresh-token", refreshToken)
		}
	}

	// Rutas de administrador con middleware de admin
	admin := router.Group("/api/admin")
	admin.Use(authMiddleware(), adminMiddleware())
	{
		admin.GET("/stats", getAdminStats)
		admin.GET("/dashboard", getAdminDashboard)
		admin.GET("/activity-log", getActivityLog)
		admin.GET("/sales-stats", getSalesStats)
		admin.GET("/users", listUsers)
		admin.GET("/users/:id", getUserById)
		admin.PUT("/users/:id", updateUser)
		admin.DELETE("/users/:id", deleteUser)
		admin.PUT("/users/:id/role", changeUserRole)
		
		// Rutas para mensajes de contacto
		admin.GET("/messages", getContactMessages)
		admin.GET("/messages/:id", getContactMessage)
		admin.PATCH("/messages/:id/:action", updateMessageStatus)
		admin.DELETE("/messages/:id", deleteContactMessage)
		admin.POST("/messages/:id/reply", replyToMessage)
	}

	cursos := router.Group("/api/cursos")
	{
		cursos.GET("", getCursos)
		cursos.GET("/:id", getCursoById)
		cursos.POST("", authMiddleware(), createCurso)
		cursos.PUT("/:id", authMiddleware(), updateCurso)
		cursos.DELETE("/:id", authMiddleware(), deleteCurso)
	}

	capitulos := router.Group("/api/capitulos")
	{
		capitulos.Use(authMiddleware())
		capitulos.GET("/curso/:cursoId", getCapitulosByCurso)
		capitulos.POST("", createCapitulo)
		capitulos.PUT("/:id", updateCapitulo)
		capitulos.DELETE("/:id", deleteCapitulo)
	}

	videos := router.Group("/api/videos")
	{
		videos.Use(authMiddleware())
		videos.POST("/upload", uploadVideo)
		videos.DELETE("/:cursoId/:filename", deleteVideo)
	}

	// Rutas para el portfolio
	portfolio := router.Group("/api/portfolio")
	{
		portfolio.GET("", getAllProjects)
		portfolio.GET("/:id", getProjectById)
		portfolio.GET("/category/:category", getProjectsByCategory)
		portfolio.POST("", authMiddleware(), adminMiddleware(), createProject)
		portfolio.PUT("/:id", authMiddleware(), adminMiddleware(), updateProject)
		portfolio.DELETE("/:id", authMiddleware(), adminMiddleware(), deleteProject)
		portfolio.POST("/reorder", authMiddleware(), adminMiddleware(), reorderProjects)
		portfolio.GET("/stats", authMiddleware(), adminMiddleware(), getPortfolioStats)
	}

	// Rutas para archivos estáticos
	router.GET("/static/videos/:cursoId/:filename", getVideo)
	router.GET("/static/profiles/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		c.File("./static/profiles/" + filename)
	})
	// Ruta para servir imágenes del portfolio
	router.GET("/static/portfolio/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		c.File("./static/portfolio/" + filename)
	})

	pagos := router.Group("/api/pagos")
	{
		pagos.Use(authMiddleware())
		pagos.POST("", crearPago)
		pagos.GET("/:id", verificarPagoPorCurso)
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

// Función auxiliar para crear directorios si no existen
func createDirIfNotExists(dirPath string) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		log.Printf("Creando directorio: %s", dirPath)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			log.Printf("Error al crear directorio %s: %v", dirPath, err)
		}
	}
}