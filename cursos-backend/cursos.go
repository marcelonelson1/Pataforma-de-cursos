package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Constantes para manejo de videos
const (
	VIDEOS_DIR     = "./static/videos"
	MAX_UPLOAD_SIZE = 100 * 1024 * 1024 // 100 MB
)

type CursoRequest struct {
	Titulo      string  `json:"titulo" binding:"required"`
	Descripcion string  `json:"descripcion" binding:"required"`
	Contenido   string  `json:"contenido" binding:"required"`
	Precio      float64 `json:"precio"`
	Estado      string  `json:"estado"`
	ImagenURL   string  `json:"imagen_url"`
}

type CapituloRequest struct {
	CursoID     uint   `json:"curso_id" binding:"required"`
	Titulo      string `json:"titulo" binding:"required"`
	Descripcion string `json:"descripcion" binding:"required"`
	Duracion    string `json:"duracion"`
	VideoURL    string `json:"video_url"`
	VideoNombre string `json:"video_nombre"`
	Publicado   bool   `json:"publicado"`
	Orden       int    `json:"orden"`
}

type VideoResponse struct {
	VideoURL string `json:"video_url"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

// Inicializa el directorio de videos si no existe
func initVideosDir() {
	if _, err := os.Stat(VIDEOS_DIR); os.IsNotExist(err) {
		log.Printf("Creando directorio para videos: %s", VIDEOS_DIR)
		if err := os.MkdirAll(VIDEOS_DIR, 0755); err != nil {
			log.Printf("Error al crear directorio para videos: %v", err)
		}
	}
}

// Maneja la subida de videos
func uploadVideo(c *gin.Context) {
	// Autenticación
	user, exists := c.Get("user")
	if !exists {
		log.Print("Usuario no autenticado intentando subir video")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}
	usuarioActual, ok := user.(Usuario)
	if !ok {
		log.Printf("Error al convertir usuario del contexto: %T", user)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}
	
	log.Printf("Usuario %s (ID: %d) iniciando subida de video", usuarioActual.Nombre, usuarioActual.ID)

	// Obtener curso_id y capitulo_id de los parámetros del formulario
	cursoID := c.PostForm("curso_id")
	capituloID := c.PostForm("capitulo_id")

	if cursoID == "" {
		log.Print("Falta el ID del curso")
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID del curso requerido"})
		return
	}

	// Verificar si el curso existe
	var curso Curso
	cursoIDInt, err := strconv.Atoi(cursoID)
	if err != nil {
		log.Printf("ID de curso inválido: %s", cursoID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de curso inválido"})
		return
	}

	if result := db.First(&curso, cursoIDInt); result.Error != nil {
		log.Printf("Curso no encontrado ID: %s", cursoID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Curso no encontrado"})
		return
	}

	// Obtener el archivo
	file, header, err := c.Request.FormFile("video")
	if err != nil {
		log.Printf("Error al obtener el archivo: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Archivo de video no proporcionado"})
		return
	}
	defer file.Close()

	// Validar tamaño del archivo
	if header.Size > MAX_UPLOAD_SIZE {
		log.Printf("Archivo demasiado grande: %d bytes", header.Size)
		c.JSON(http.StatusBadRequest, gin.H{"error": "El archivo no debe superar los 100 MB"})
		return
	}

	// Validar tipo de archivo
	fileExt := strings.ToLower(filepath.Ext(header.Filename))
	if fileExt != ".mp4" && fileExt != ".webm" && fileExt != ".ogg" {
		log.Printf("Tipo de archivo no permitido: %s", fileExt)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Solo se permiten archivos MP4, WebM y OGG"})
		return
	}

	// Generar nombre de archivo único
	uniqueID := uuid.New().String()
	filename := fmt.Sprintf("%s-%s-%s%s", 
		cursoID, 
		capituloID, 
		uniqueID, 
		fileExt)
	
	// Crear directorio específico para el curso si no existe
	cursoDir := filepath.Join(VIDEOS_DIR, cursoID)
	if _, err := os.Stat(cursoDir); os.IsNotExist(err) {
		if err := os.MkdirAll(cursoDir, 0755); err != nil {
			log.Printf("Error al crear directorio para el curso: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar el video"})
			return
		}
	}
	
	// Ruta completa del archivo
	filePath := filepath.Join(cursoDir, filename)
	
	// Crear el archivo en el servidor
	out, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error al crear el archivo: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar el video"})
		return
	}
	defer out.Close()

	// Copiar el contenido del archivo
	if _, err = io.Copy(out, file); err != nil {
		log.Printf("Error al copiar el archivo: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar el video"})
		return
	}

	// URL donde se puede acceder al video
	videoURL := fmt.Sprintf("/static/videos/%s/%s", cursoID, filename)
	
	log.Printf("Video subido exitosamente: %s", videoURL)
	
	// Responder con la URL del video
	c.JSON(http.StatusOK, VideoResponse{
		VideoURL: videoURL,
		Filename: filename,
		Size: header.Size,
	})
}

// Endpoint para obtener un video
func getVideo(c *gin.Context) {
	cursoID := c.Param("cursoId")
	filename := c.Param("filename")
	
	if cursoID == "" || filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetros incompletos"})
		return
	}
	
	filePath := filepath.Join(VIDEOS_DIR, cursoID, filename)
	
	// Verificar si el archivo existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("Video no encontrado: %s", filePath)
		c.JSON(http.StatusNotFound, gin.H{"error": "Video no encontrado"})
		return
	}
	
	// Establecer headers para streaming
	c.Header("Content-Type", "video/mp4") // Ajustar según el tipo de archivo
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", filename))
	c.Header("Cache-Control", "public, max-age=31536000")
	
	// Entregar el archivo
	c.File(filePath)
}

// Eliminar un video
func deleteVideo(c *gin.Context) {
	cursoID := c.Param("cursoId")
	filename := c.Param("filename")
	
	if cursoID == "" || filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetros incompletos"})
		return
	}
	
	filePath := filepath.Join(VIDEOS_DIR, cursoID, filename)
	
	// Verificar si el archivo existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("Video no encontrado para eliminar: %s", filePath)
		c.JSON(http.StatusNotFound, gin.H{"error": "Video no encontrado"})
		return
	}
	
	// Eliminar el archivo
	if err := os.Remove(filePath); err != nil {
		log.Printf("Error al eliminar el video: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el video"})
		return
	}
	
	log.Printf("Video eliminado exitosamente: %s", filePath)
	c.JSON(http.StatusOK, gin.H{"message": "Video eliminado correctamente"})
}

func getCursos(c *gin.Context) {
	var cursos []Curso
	
	// Utilizamos Preload para cargar también los capítulos de cada curso
	result := db.Preload("Capitulos").Find(&cursos)
	if result.Error != nil {
		log.Printf("Error al obtener cursos: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener cursos: " + result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, cursos)
}

func getCursoById(c *gin.Context) {
	id := c.Param("id")
	
	var curso Curso
	if result := db.Preload("Capitulos").First(&curso, id); result.Error != nil {
		log.Printf("Curso no encontrado ID: %s, Error: %v", id, result.Error)
		c.JSON(http.StatusNotFound, gin.H{"error": "Curso no encontrado"})
		return
	}

	c.JSON(http.StatusOK, curso)
}

func createCurso(c *gin.Context) {
	var req CursoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error en el binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar que el usuario está autenticado y obtenerlo
	user, exists := c.Get("user")
	if !exists {
		log.Print("Usuario no encontrado en el contexto - token puede ser inválido o expirado")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado. Por favor, inicie sesión nuevamente."})
		return
	}
	
	usuarioActual, ok := user.(Usuario)
	if !ok {
		log.Printf("Error al convertir usuario del contexto: %T", user)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}
	
	log.Printf("Usuario autenticado ID: %v, Nombre: %s", usuarioActual.ID, usuarioActual.Nombre)

	// Establecer valores por defecto si no se proporcionan
	if req.Estado == "" {
		req.Estado = "Borrador"
	}

	curso := Curso{
		Titulo:      req.Titulo,
		Descripcion: req.Descripcion,
		Contenido:   req.Contenido,
		Precio:      req.Precio,
		Estado:      req.Estado,
		ImagenURL:   req.ImagenURL,
	}

	log.Printf("Creando curso: %+v", curso)
	
	if result := db.Create(&curso); result.Error != nil {
		log.Printf("Error al crear curso en DB: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear curso: " + result.Error.Error()})
		return
	}

	// Devolvemos el curso con sus capítulos (que serán ninguno inicialmente)
	var cursoCompleto Curso
	if result := db.Preload("Capitulos").First(&cursoCompleto, curso.ID); result.Error != nil {
		log.Printf("Curso creado pero no se pudo recuperar con capítulos: %v", result.Error)
		// Si no podemos recuperar el curso completo, inicializamos el array de capítulos
		curso.Capitulos = []Capitulo{}
		c.JSON(http.StatusCreated, curso)
		return
	}
	
	// Asegurar que capitulos no sea nil para evitar errores en el frontend
	if cursoCompleto.Capitulos == nil {
		cursoCompleto.Capitulos = []Capitulo{}
	}
	
	log.Printf("Curso creado exitosamente: ID %v", cursoCompleto.ID)
	c.JSON(http.StatusCreated, cursoCompleto)
}

func updateCurso(c *gin.Context) {
	id := c.Param("id")
	
	var curso Curso
	if result := db.First(&curso, id); result.Error != nil {
		log.Printf("Curso no encontrado para actualizar ID: %s, Error: %v", id, result.Error)
		c.JSON(http.StatusNotFound, gin.H{"error": "Curso no encontrado"})
		return
	}

	var req CursoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error en el binding JSON para actualización: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Actualizar solo los campos proporcionados
	curso.Titulo = req.Titulo
	curso.Descripcion = req.Descripcion
	curso.Contenido = req.Contenido
	curso.Precio = req.Precio
	curso.Estado = req.Estado
	if req.ImagenURL != "" {
		curso.ImagenURL = req.ImagenURL
	}

	log.Printf("Actualizando curso ID %s: %+v", id, curso)
	
	if result := db.Save(&curso); result.Error != nil {
		log.Printf("Error al actualizar curso en DB: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar curso: " + result.Error.Error()})
		return
	}

	// Devolvemos el curso completo con sus capítulos
	var cursoCompleto Curso
	if result := db.Preload("Capitulos").First(&cursoCompleto, curso.ID); result.Error != nil {
		log.Printf("Curso actualizado pero no se pudo recuperar con capítulos: %v", result.Error)
		// Si no podemos recuperar el curso completo, inicializamos el array de capítulos
		curso.Capitulos = []Capitulo{}
		c.JSON(http.StatusOK, curso)
		return
	}
	
	// Asegurar que capitulos no sea nil para evitar errores en el frontend
	if cursoCompleto.Capitulos == nil {
		cursoCompleto.Capitulos = []Capitulo{}
	}
	
	log.Printf("Curso actualizado exitosamente: ID %v", cursoCompleto.ID)
	c.JSON(http.StatusOK, cursoCompleto)
}

func deleteCurso(c *gin.Context) {
	id := c.Param("id")
	
	var curso Curso
	if result := db.First(&curso, id); result.Error != nil {
		log.Printf("Curso no encontrado para eliminar ID: %s, Error: %v", id, result.Error)
		c.JSON(http.StatusNotFound, gin.H{"error": "Curso no encontrado"})
		return
	}

	// Cargar los capítulos para eliminar archivos de video
	var capitulos []Capitulo
	if result := db.Where("curso_id = ?", id).Find(&capitulos); result.Error == nil {
		for _, capitulo := range capitulos {
			if capitulo.VideoNombre != "" {
				// Construir la ruta del archivo
				cursoIDStr := strconv.FormatUint(uint64(curso.ID), 10)
				filePath := filepath.Join(VIDEOS_DIR, cursoIDStr, capitulo.VideoNombre)
				
				// Intentar eliminar el archivo
				if _, err := os.Stat(filePath); !os.IsNotExist(err) {
					if err := os.Remove(filePath); err != nil {
						log.Printf("Error al eliminar archivo de video del capítulo %d: %v", capitulo.ID, err)
					} else {
						log.Printf("Archivo de video del capítulo %d eliminado: %s", capitulo.ID, filePath)
					}
				}
			}
		}
	}

	log.Printf("Eliminando capítulos del curso ID: %s", id)
	
	// Primero eliminamos todos los capítulos relacionados para evitar problemas de integridad
	if result := db.Where("curso_id = ?", id).Delete(&Capitulo{}); result.Error != nil {
		log.Printf("Error al eliminar capítulos: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar capítulos del curso: " + result.Error.Error()})
		return
	}

	log.Printf("Eliminando curso ID: %s", id)
	
	// Luego eliminamos el curso
	if result := db.Delete(&curso); result.Error != nil {
		log.Printf("Error al eliminar curso: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar curso: " + result.Error.Error()})
		return
	}

	// Intentar eliminar el directorio de videos del curso si existe
	cursoIDStr := strconv.FormatUint(uint64(curso.ID), 10)
	dirPath := filepath.Join(VIDEOS_DIR, cursoIDStr)
	if _, err := os.Stat(dirPath); !os.IsNotExist(err) {
		if err := os.RemoveAll(dirPath); err != nil {
			log.Printf("Error al eliminar directorio de videos del curso: %v", err)
		} else {
			log.Printf("Directorio de videos del curso eliminado: %s", dirPath)
		}
	}

	log.Printf("Curso ID %s eliminado exitosamente", id)
	c.JSON(http.StatusOK, gin.H{"message": "Curso eliminado correctamente"})
}

func getCapitulosByCurso(c *gin.Context) {
	cursoId := c.Param("cursoId")
	
	var capitulos []Capitulo
	if result := db.Where("curso_id = ?", cursoId).Order("orden ASC").Find(&capitulos); result.Error != nil {
		log.Printf("Error al obtener capítulos del curso ID %s: %v", cursoId, result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener capítulos: " + result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, capitulos)
}

func createCapitulo(c *gin.Context) {
	var req CapituloRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error en el binding JSON para capítulo: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validamos que el curso exista
	var curso Curso
	if result := db.First(&curso, req.CursoID); result.Error != nil {
		log.Printf("Curso no encontrado ID: %d para crear capítulo", req.CursoID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Curso no encontrado"})
		return
	}

	// Extraer el nombre del archivo desde la URL si existe
	videoNombre := ""
	if req.VideoURL != "" {
		parts := strings.Split(req.VideoURL, "/")
		if len(parts) > 0 {
			videoNombre = parts[len(parts)-1]
		}
	}

	capitulo := Capitulo{
		CursoID:     req.CursoID,
		Titulo:      req.Titulo,
		Descripcion: req.Descripcion,
		Duracion:    req.Duracion,
		VideoURL:    req.VideoURL,
		VideoNombre: videoNombre,
		Publicado:   req.Publicado,
		Orden:       req.Orden,
	}

	log.Printf("Creando capítulo: %+v", capitulo)
	
	if result := db.Create(&capitulo); result.Error != nil {
		log.Printf("Error al crear capítulo: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear capítulo: " + result.Error.Error()})
		return
	}

	log.Printf("Capítulo creado exitosamente: ID %v", capitulo.ID)
	c.JSON(http.StatusCreated, capitulo)
}

func updateCapitulo(c *gin.Context) {
	id := c.Param("id")
	
	var capitulo Capitulo
	if result := db.First(&capitulo, id); result.Error != nil {
		log.Printf("Capítulo no encontrado ID: %s para actualizar", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Capítulo no encontrado"})
		return
	}

	var req CapituloRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error en el binding JSON para actualizar capítulo: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificamos que el curso exista (seguridad adicional)
	var curso Curso
	if result := db.First(&curso, req.CursoID); result.Error != nil {
		log.Printf("Curso no encontrado ID: %d para actualizar capítulo", req.CursoID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Curso no encontrado"})
		return
	}

	// Extraer el nombre del archivo desde la URL si cambió
	videoNombre := capitulo.VideoNombre
	if req.VideoURL != capitulo.VideoURL {
		videoNombre = ""
		if req.VideoURL != "" {
			parts := strings.Split(req.VideoURL, "/")
			if len(parts) > 0 {
				videoNombre = parts[len(parts)-1]
			}
		}
	}

	capitulo.CursoID = req.CursoID
	capitulo.Titulo = req.Titulo
	capitulo.Descripcion = req.Descripcion
	capitulo.Duracion = req.Duracion
	capitulo.VideoURL = req.VideoURL
	capitulo.VideoNombre = videoNombre
	capitulo.Publicado = req.Publicado
	capitulo.Orden = req.Orden

	log.Printf("Actualizando capítulo ID %s: %+v", id, capitulo)
	
	if result := db.Save(&capitulo); result.Error != nil {
		log.Printf("Error al actualizar capítulo: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar capítulo: " + result.Error.Error()})
		return
	}

	log.Printf("Capítulo actualizado exitosamente: ID %v", capitulo.ID)
	c.JSON(http.StatusOK, capitulo)
}

func deleteCapitulo(c *gin.Context) {
	id := c.Param("id")
	
	var capitulo Capitulo
	if result := db.First(&capitulo, id); result.Error != nil {
		log.Printf("Capítulo no encontrado ID: %s para eliminar", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Capítulo no encontrado"})
		return
	}

	// Si hay un video asociado, eliminarlo
	if capitulo.VideoNombre != "" && capitulo.CursoID > 0 {
		cursoID := strconv.FormatUint(uint64(capitulo.CursoID), 10)
		filePath := filepath.Join(VIDEOS_DIR, cursoID, capitulo.VideoNombre)
		
		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			log.Printf("Eliminando archivo de video: %s", filePath)
			if err := os.Remove(filePath); err != nil {
				log.Printf("Advertencia: No se pudo eliminar el archivo de video: %v", err)
			}
		}
	}

	log.Printf("Eliminando capítulo ID: %s", id)
	
	if result := db.Delete(&capitulo); result.Error != nil {
		log.Printf("Error al eliminar capítulo: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar capítulo: " + result.Error.Error()})
		return
	}

	log.Printf("Capítulo ID %s eliminado exitosamente", id)
	c.JSON(http.StatusOK, gin.H{"message": "Capítulo eliminado correctamente"})
}