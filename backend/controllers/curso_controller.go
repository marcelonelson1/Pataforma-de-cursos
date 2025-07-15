package controllers

import (
	"curso-platform/middleware"
	"curso-platform/models"
	"curso-platform/services"
	"curso-platform/utils"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"fmt"

	"github.com/gin-gonic/gin"
)

// CursoController gestiona las operaciones relacionadas con los cursos
type CursoController struct {
	cursoService *services.CursoService
}

// NewCursoController crea una nueva instancia del controlador de cursos
func NewCursoController(cursoService *services.CursoService) *CursoController {
	return &CursoController{
		cursoService: cursoService,
	}
}

// GetCursos obtiene todos los cursos
func (c *CursoController) GetCursos(ctx *gin.Context) {
	cursos, err := c.cursoService.GetAll()
	if err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, cursos)
}

// GetCursoById obtiene un curso por su ID
func (c *CursoController) GetCursoById(ctx *gin.Context) {
	id := ctx.Param("id")
	cursoID, err := strconv.Atoi(id)
	if err != nil {
		utils.SendErrorResponse(ctx, utils.ErrBadRequest, http.StatusBadRequest)
		return
	}

	curso, err := c.cursoService.GetByID(uint(cursoID))
	if err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, curso)
}

// CreateCurso crea un nuevo curso
func (c *CursoController) CreateCurso(ctx *gin.Context) {
	// Verificar que el usuario está autenticado y obtenerlo
	user, exists := ctx.Get("user")
	if !exists {
		log.Print("Usuario no encontrado en el contexto - token puede ser inválido o expirado")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado. Por favor, inicie sesión nuevamente."})
		return
	}

	usuarioActual, ok := user.(models.Usuario)
	if !ok {
		log.Printf("Error al convertir usuario del contexto: %T", user)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	log.Printf("Usuario autenticado ID: %v, Nombre: %s", usuarioActual.ID, usuarioActual.Nombre)

	// Obtener los datos del formulario
	titulo := ctx.PostForm("titulo")
	descripcion := ctx.PostForm("descripcion")
	contenido := ctx.PostForm("contenido")
	precioStr := ctx.PostForm("precio")
	estado := ctx.PostForm("estado")
	imagenURL := ctx.PostForm("imagen_url")

	// Validar campos requeridos
	if titulo == "" || descripcion == "" || contenido == "" {
		log.Print("Campos requeridos faltantes")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Campos requeridos incompletos"})
		return
	}

	// Convertir precio a float64
	precio, err := strconv.ParseFloat(precioStr, 64)
	if err != nil {
		log.Printf("Error al convertir precio: %v", err)
		precio = 0.0
	}

	// Manejar la carga de la imagen
	var imageFile *os.File
	file, err := ctx.FormFile("imagen")
	if err == nil {
		// Si hay un archivo de imagen, crear un archivo temporal
		tempFile, err := os.CreateTemp("", "curso-*.tmp")
		if err != nil {
			utils.SendErrorResponse(ctx, utils.ErrServerError, http.StatusInternalServerError)
			return
		}
		defer os.Remove(tempFile.Name())
		defer tempFile.Close()

		// Guardar FormFile en archivo temporal
		if err := ctx.SaveUploadedFile(file, tempFile.Name()); err != nil {
			utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
			return
		}

		imageFile = tempFile
	}

	// Crear curso
	curso, err := c.cursoService.Create(titulo, descripcion, contenido, estado, precio, imagenURL, imageFile)
	if err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, curso)
}

// UpdateCurso actualiza un curso existente
func (c *CursoController) UpdateCurso(ctx *gin.Context) {
	id := ctx.Param("id")
	cursoID, err := strconv.Atoi(id)
	if err != nil {
		utils.SendErrorResponse(ctx, utils.ErrBadRequest, http.StatusBadRequest)
		return
	}

	// Obtener los datos del formulario
	titulo := ctx.PostForm("titulo")
	descripcion := ctx.PostForm("descripcion")
	contenido := ctx.PostForm("contenido")
	precioStr := ctx.PostForm("precio")
	estado := ctx.PostForm("estado")
	imagenURL := ctx.PostForm("imagen_url")

	// Validar campos requeridos
	if titulo == "" || descripcion == "" || contenido == "" {
		log.Print("Campos requeridos faltantes")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Campos requeridos incompletos"})
		return
	}

	// Convertir precio a float64
	precio, err := strconv.ParseFloat(precioStr, 64)
	if err != nil {
		log.Printf("Error al convertir precio: %v", err)
		precio = 0.0
	}

	// Manejar la carga de la imagen
	var imageFile *os.File
	file, err := ctx.FormFile("imagen")
	if err == nil {
		// Si hay un archivo de imagen, crear un archivo temporal
		tempFile, err := os.CreateTemp("", "curso-*.tmp")
		if err != nil {
			utils.SendErrorResponse(ctx, utils.ErrServerError, http.StatusInternalServerError)
			return
		}
		defer os.Remove(tempFile.Name())
		defer tempFile.Close()

		// Guardar FormFile en archivo temporal
		if err := ctx.SaveUploadedFile(file, tempFile.Name()); err != nil {
			utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
			return
		}

		imageFile = tempFile
	}

	// Actualizar curso
	curso, err := c.cursoService.Update(uint(cursoID), titulo, descripcion, contenido, estado, precio, imagenURL, imageFile)
	if err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, curso)
}

// DeleteCurso elimina un curso
func (c *CursoController) DeleteCurso(ctx *gin.Context) {
	id := ctx.Param("id")
	cursoID, err := strconv.Atoi(id)
	if err != nil {
		utils.SendErrorResponse(ctx, utils.ErrBadRequest, http.StatusBadRequest)
		return
	}

	if err := c.cursoService.Delete(uint(cursoID)); err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Curso eliminado correctamente"})
}

// UploadVideo sube un video para un capítulo
func (c *CursoController) UploadVideo(ctx *gin.Context) {
	// Autenticación
	user, exists := ctx.Get("user")
	if !exists {
		log.Print("Usuario no autenticado intentando subir video")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}
	
	usuarioActual, ok := user.(models.Usuario)
	if !ok {
		log.Printf("Error al convertir usuario del contexto: %T", user)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	log.Printf("Usuario %s (ID: %d) iniciando subida de video", usuarioActual.Nombre, usuarioActual.ID)

	// Obtener curso_id y capitulo_id de los parámetros del formulario
	cursoID := ctx.PostForm("curso_id")
	capituloID := ctx.PostForm("capitulo_id")

	if cursoID == "" {
		log.Print("Falta el ID del curso")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID del curso requerido"})
		return
	}

	cursoIDInt, err := strconv.Atoi(cursoID)
	if err != nil {
		log.Printf("ID de curso inválido: %s", cursoID)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de curso inválido"})
		return
	}

	capituloIDInt := 0
	if capituloID != "" {
		capituloIDInt, err = strconv.Atoi(capituloID)
		if err != nil {
			log.Printf("ID de capítulo inválido: %s", capituloID)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de capítulo inválido"})
			return
		}
	}

	// Obtener el archivo
	file, header, err := ctx.Request.FormFile("video")
	if err != nil {
		log.Printf("Error al obtener el archivo: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Archivo de video no proporcionado"})
		return
	}
	defer file.Close()

	// Crear archivo temporal para procesar
	tempFile, err := os.CreateTemp("", "video-*.tmp")
	if err != nil {
		utils.SendErrorResponse(ctx, utils.ErrServerError, http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Copiar el contenido del archivo al archivo temporal
	_, err = io.Copy(tempFile, file)
	if err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	// Subir video
	videoURL, err := c.cursoService.UploadVideo(uint(cursoIDInt), uint(capituloIDInt), tempFile, header.Size, header.Filename)
	if err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	// Responder con la URL del video
	ctx.JSON(http.StatusOK, gin.H{
		"videoURL": *videoURL,
		"filename": filepath.Base(*videoURL),
		"size":     header.Size,
	})
}

// DeleteVideo elimina un archivo de video
func (c *CursoController) DeleteVideo(ctx *gin.Context) {
	cursoID := ctx.Param("cursoId")
	filename := ctx.Param("filename")

	if cursoID == "" || filename == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Parámetros incompletos"})
		return
	}

	cursoIDInt, err := strconv.Atoi(cursoID)
	if err != nil {
		utils.SendErrorResponse(ctx, utils.ErrBadRequest, http.StatusBadRequest)
		return
	}

	if err := c.cursoService.DeleteVideo(uint(cursoIDInt), filename); err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Video eliminado correctamente"})
}

// GetVideo entrega un archivo de video
func (c *CursoController) GetVideo(ctx *gin.Context) {
	cursoID := ctx.Param("cursoId")
	filename := ctx.Param("filename")

	if cursoID == "" || filename == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Parámetros incompletos"})
		return
	}

	filePath := filepath.Join("./static/videos", cursoID, filename)

	// Verificar si el archivo existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("Video no encontrado: %s", filePath)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Video no encontrado"})
		return
	}

	// Establecer headers para streaming
	ctx.Header("Content-Type", "video/mp4") // Ajustar según el tipo de archivo
	ctx.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", filename))
	ctx.Header("Cache-Control", "public, max-age=31536000")

	// Entregar el archivo
	ctx.File(filePath)
}

// RegisterRoutes registra todas las rutas relacionadas con los cursos
func (c *CursoController) RegisterRoutes(router *gin.Engine) {
	cursos := router.Group("/api/cursos")
	{
		cursos.GET("", c.GetCursos)
		cursos.GET("/:id", c.GetCursoById)
		cursos.POST("", middleware.AuthMiddleware(), c.CreateCurso)
		cursos.PUT("/:id", middleware.AuthMiddleware(), c.UpdateCurso)
		cursos.DELETE("/:id", middleware.AuthMiddleware(), c.DeleteCurso)
	}

	videos := router.Group("/api/videos")
	{
		videos.Use(middleware.AuthMiddleware())
		videos.POST("/upload", c.UploadVideo)
		videos.DELETE("/:cursoId/:filename", c.DeleteVideo)
	}

	// Ruta pública para obtener videos
	router.GET("/static/videos/:cursoId/:filename", c.GetVideo)
}