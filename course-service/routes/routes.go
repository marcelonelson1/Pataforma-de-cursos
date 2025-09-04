// routes/routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"course-service/config"
	"course-service/controllers"
	"course-service/utils"
)

// SetupCourseRoutes configura todas las rutas del Course Service
func SetupCourseRoutes(router *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// Inicializar controladores
	courseController := controllers.NewCourseController(db, cfg)
	chapterController := controllers.NewChapterController(db, cfg)
	progressController := controllers.NewProgressController(db, cfg)
	fileController := controllers.NewFileController(db, cfg)

	// API Group
	api := router.Group("/api")

	// Rutas públicas de cursos
	courses := api.Group("/courses")
	{
		courses.GET("", courseController.GetCourses)           // Listar cursos públicos
		courses.GET("/:id", courseController.GetCourseByID)   // Obtener curso específico
	}

	// Rutas autenticadas de cursos
	coursesAuth := api.Group("/courses")
	coursesAuth.Use(utils.AuthMiddleware(cfg))
	{
		coursesAuth.POST("", courseController.CreateCourse)                    // Crear curso
		coursesAuth.PUT("/:id", courseController.UpdateCourse)                // Actualizar curso
		coursesAuth.DELETE("/:id", courseController.DeleteCourse)             // Eliminar curso
		coursesAuth.GET("/:id/access", courseController.CheckCourseAccess)    // Verificar acceso
	}

	// Rutas de capítulos (requieren autenticación)
	chapters := api.Group("/chapters")
	chapters.Use(utils.AuthMiddleware(cfg))
	{
		chapters.GET("/course/:id", chapterController.GetChaptersByCourse)    // Capítulos de un curso
		chapters.POST("", chapterController.CreateChapter)                   // Crear capítulo
		chapters.PUT("/:id", chapterController.UpdateChapter)                // Actualizar capítulo
		chapters.DELETE("/:id", chapterController.DeleteChapter)             // Eliminar capítulo
		chapters.GET("/:id/content", chapterController.GetChapterContent)    // Contenido del capítulo
		chapters.POST("/reorder", chapterController.ReorderChapters)         // Reordenar capítulos
		chapters.PATCH("/:id/publish", chapterController.ToggleChapterPublished) // Toggle publicación
	}

	// Rutas de progreso (requieren autenticación)
	progress := api.Group("/progress")
	progress.Use(utils.AuthMiddleware(cfg))
	{
		progress.GET("/course/:id", progressController.GetUserProgress)           // Progreso en curso
		progress.POST("/chapter/complete", progressController.MarkChapterCompleted) // Marcar completado
		progress.POST("/last-chapter", progressController.SaveLastChapter)        // Guardar último capítulo
		progress.GET("/user/summary", progressController.GetUserProgressSummary)  // Resumen del usuario
		progress.POST("/chapter/watch-time", progressController.UpdateChapterWatchTime) // Tiempo visto
		progress.GET("/course/:id/stats", progressController.GetCourseStatistics) // Estadísticas del curso
	}

	// Rutas de archivos
	files := api.Group("/files")
	{
		// Upload (requiere autenticación)
		filesAuth := files.Group("")
		filesAuth.Use(utils.AuthMiddleware(cfg))
		{
			filesAuth.POST("/upload-video", fileController.UploadVideo)  // Subir video
			filesAuth.POST("/upload-image", fileController.UploadImage)  // Subir imagen
			filesAuth.DELETE("/:courseId/:filename", fileController.DeleteFile) // Eliminar archivo
		}
		
		// Acceso a archivos (algunos requieren verificación de acceso)
		files.GET("/video/:courseId/:filename", fileController.GetVideo)     // Streaming de video
		files.GET("/image/:filename", fileController.GetImage)               // Imágenes públicas
		files.GET("/thumbnail/:courseId/:filename", fileController.GetThumbnail) // Miniaturas
	}

	// Rutas de categorías
	categories := api.Group("/categories")
	{
		categories.GET("", courseController.GetCourseCategories) // Listar categorías (público)
		
		// Admin routes para categorías (requieren autenticación y permisos de admin)
		categoriesAdmin := categories.Group("")
		categoriesAdmin.Use(utils.AuthMiddleware(cfg), utils.AdminMiddleware())
		{
			categoriesAdmin.POST("", courseController.CreateCategory)           // Crear categoría
			categoriesAdmin.PUT("/:id", courseController.UpdateCategory)        // Actualizar categoría
			categoriesAdmin.DELETE("/:id", courseController.DeleteCategory)     // Eliminar categoría
		}
	}

	// Health check y estadísticas
	api.GET("/health", courseController.HealthCheck)

	// Rutas de administración (requieren autenticación y permisos de admin)
	admin := api.Group("/admin")
	admin.Use(utils.AuthMiddleware(cfg), utils.AdminMiddleware())
	{
		admin.GET("/courses", courseController.GetAdminCourses) // Obtener todos los cursos para admin
	}
}