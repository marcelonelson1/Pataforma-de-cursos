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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Constantes para manejo de videos y imágenes
const (
	VIDEOS_DIR      = "./static/videos"
	IMAGES_DIR      = "./static/images"
	MAX_UPLOAD_SIZE = 100 * 1024 * 1024 // 100 MB
	MAX_IMAGE_SIZE  = 5 * 1024 * 1024   // 5 MB
)

// Estructura para el progreso del usuario en un curso
type ProgresoUsuario struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UsuarioID       uint      `gorm:"not null;index:idx_usuario_curso" json:"usuario_id"`
	CursoID         uint      `gorm:"not null;index:idx_usuario_curso" json:"curso_id"`
	PorcentajeTotal float64   `gorm:"type:decimal(5,2);default:0" json:"porcentaje_total"`
	UltimoCapitulo  uint      `gorm:"default:0" json:"ultimo_capitulo"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Estructura para el progreso del usuario en un capítulo específico
type ProgresoCapitulo struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UsuarioID  uint      `gorm:"not null;index:idx_usuario_capitulo" json:"usuario_id"`
	CursoID    uint      `gorm:"not null;index:idx_usuario_capitulo" json:"curso_id"`
	CapituloID uint      `gorm:"not null;index:idx_usuario_capitulo" json:"capitulo_id"`
	Completado bool      `gorm:"default:false" json:"completado"`
	Progreso   float64   `gorm:"type:decimal(5,2);default:0" json:"progreso"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Estructura para la respuesta de progreso
type ProgresoResponse struct {
	ProgresoTotal     float64                   `json:"progreso_total"`
	UltimoCapitulo    uint                      `json:"ultimo_capitulo"`
	CapitulosProgreso map[uint]ProgresoCapitulo `json:"capitulos_progreso"`
}

// Estructura para marcar un capítulo como completado o actualizar su progreso
type ProgresoCapituloRequest struct {
	CursoID    uint    `json:"curso_id" binding:"required"`
	CapituloID uint    `json:"capitulo_id" binding:"required"`
	Completado bool    `json:"completado"`
	Progreso   float64 `json:"progreso"`
}

// Estructura para actualizar el último capítulo visto
type UltimoCapituloRequest struct {
	CursoID    uint `json:"curso_id" binding:"required"`
	CapituloID uint `json:"capitulo_id" binding:"required"`
}

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

// Inicializa los directorios de videos e imágenes si no existen
func initVideosDir() {
	if _, err := os.Stat(VIDEOS_DIR); os.IsNotExist(err) {
		log.Printf("Creando directorio para videos: %s", VIDEOS_DIR)
		if err := os.MkdirAll(VIDEOS_DIR, 0755); err != nil {
			log.Printf("Error al crear directorio para videos: %v", err)
		}
	}

	if _, err := os.Stat(IMAGES_DIR); os.IsNotExist(err) {
		log.Printf("Creando directorio para imágenes: %s", IMAGES_DIR)
		if err := os.MkdirAll(IMAGES_DIR, 0755); err != nil {
			log.Printf("Error al crear directorio para imágenes: %v", err)
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
		Size:     header.Size,
	})
}

// Función para manejar la subida de imágenes
func saveImageFile(c *gin.Context, fieldName string) (string, error) {
	file, header, err := c.Request.FormFile(fieldName)
	if err != nil {
		// Si no hay archivo, no es un error
		if err == http.ErrMissingFile {
			return "", nil
		}
		return "", err
	}
	defer file.Close()

	// Validar tamaño del archivo
	if header.Size > MAX_IMAGE_SIZE {
		return "", fmt.Errorf("La imagen no debe superar los 5 MB")
	}

	// Validar tipo de archivo
	fileExt := strings.ToLower(filepath.Ext(header.Filename))
	if fileExt != ".jpg" && fileExt != ".jpeg" && fileExt != ".png" && fileExt != ".gif" && fileExt != ".webp" {
		return "", fmt.Errorf("Solo se permiten archivos JPG, PNG, GIF y WEBP")
	}

	// Generar nombre de archivo único
	uniqueID := uuid.New().String()
	filename := fmt.Sprintf("%s%s", uniqueID, fileExt)

	// Crear directorio específico para imágenes si no existe
	if _, err := os.Stat(IMAGES_DIR); os.IsNotExist(err) {
		if err := os.MkdirAll(IMAGES_DIR, 0755); err != nil {
			return "", fmt.Errorf("Error al crear directorio para imágenes: %v", err)
		}
	}

	// Ruta completa del archivo
	filePath := filepath.Join(IMAGES_DIR, filename)

	// Crear el archivo en el servidor
	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("Error al crear el archivo: %v", err)
	}
	defer out.Close()

	// Copiar el contenido del archivo
	if _, err = io.Copy(out, file); err != nil {
		return "", fmt.Errorf("Error al copiar el archivo: %v", err)
	}

	// URL donde se puede acceder a la imagen
	imageURL := fmt.Sprintf("/static/images/%s", filename)

	return imageURL, nil
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

	// Obtener los datos del formulario
	titulo := c.PostForm("titulo")
	descripcion := c.PostForm("descripcion")
	contenido := c.PostForm("contenido")
	precioStr := c.PostForm("precio")
	estado := c.PostForm("estado")

	// Validar campos requeridos
	if titulo == "" || descripcion == "" || contenido == "" {
		log.Print("Campos requeridos faltantes")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campos requeridos incompletos"})
		return
	}

	// Convertir precio a float64
	precio, err := strconv.ParseFloat(precioStr, 64)
	if err != nil {
		log.Printf("Error al convertir precio: %v", err)
		precio = 0.0
	}

	// Establecer valores por defecto si no se proporcionan
	if estado == "" {
		estado = "Borrador"
	}

	// Manejar la subida de la imagen
	imagenURL, err := saveImageFile(c, "imagen")
	if err != nil {
		log.Printf("Error al guardar imagen: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error con la imagen: %v", err)})
		return
	}

	// Si no hay nueva imagen, usar la URL proporcionada si existe
	if imagenURL == "" {
		imagenURL = c.PostForm("imagen_url")
	}

	curso := Curso{
		Titulo:      titulo,
		Descripcion: descripcion,
		Contenido:   contenido,
		Precio:      precio,
		Estado:      estado,
		ImagenURL:   imagenURL,
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

	// Obtener los datos del formulario
	titulo := c.PostForm("titulo")
	descripcion := c.PostForm("descripcion")
	contenido := c.PostForm("contenido")
	precioStr := c.PostForm("precio")
	estado := c.PostForm("estado")

	// Validar campos requeridos
	if titulo == "" || descripcion == "" || contenido == "" {
		log.Print("Campos requeridos faltantes")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campos requeridos incompletos"})
		return
	}

	// Convertir precio a float64
	precio, err := strconv.ParseFloat(precioStr, 64)
	if err != nil {
		log.Printf("Error al convertir precio: %v", err)
		// Mantener el precio anterior en caso de error
		precio = curso.Precio
	}

	// Actualizar campos del curso
	curso.Titulo = titulo
	curso.Descripcion = descripcion
	curso.Contenido = contenido
	curso.Precio = precio
	curso.Estado = estado

	// Manejar la subida de la imagen si hay una nueva
	imagenURL, err := saveImageFile(c, "imagen")
	if err != nil {
		log.Printf("Error al guardar imagen: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error con la imagen: %v", err)})
		return
	}

	// Si hay una nueva imagen, actualizar la URL
	if imagenURL != "" {
		curso.ImagenURL = imagenURL
	} else {
		// Si no hay nueva imagen, usar la URL proporcionada si existe
		imagenURLPost := c.PostForm("imagen_url")
		if imagenURLPost != "" {
			curso.ImagenURL = imagenURLPost
		}
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

	// Eliminar progreso asociado al curso
	if err := db.Where("curso_id = ?", id).Delete(&ProgresoUsuario{}).Error; err != nil {
		log.Printf("Error al eliminar progreso de usuarios del curso: %v", err)
	}

	if err := db.Where("curso_id = ?", id).Delete(&ProgresoCapitulo{}).Error; err != nil {
		log.Printf("Error al eliminar progreso de capítulos del curso: %v", err)
	}

	// Primero eliminamos todos los capítulos relacionados para evitar problemas de integridad
	if result := db.Where("curso_id = ?", id).Delete(&Capitulo{}); result.Error != nil {
		log.Printf("Error al eliminar capítulos: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar capítulos del curso: " + result.Error.Error()})
		return
	}

	// Eliminar imagen del curso si existe
	if curso.ImagenURL != "" && strings.Contains(curso.ImagenURL, "/static/images/") {
		// Extraer el nombre del archivo de la URL
		imagenNombre := filepath.Base(curso.ImagenURL)
		imagenPath := filepath.Join(IMAGES_DIR, imagenNombre)

		if _, err := os.Stat(imagenPath); !os.IsNotExist(err) {
			if err := os.Remove(imagenPath); err != nil {
				log.Printf("Error al eliminar imagen del curso: %v", err)
			} else {
				log.Printf("Imagen del curso eliminada: %s", imagenPath)
			}
		}
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

	// Eliminar progreso asociado al capítulo
	if err := db.Where("capitulo_id = ?", id).Delete(&ProgresoCapitulo{}).Error; err != nil {
		log.Printf("Error al eliminar progreso del capítulo: %v", err)
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

// ==============================================
// Funciones para el manejo de progreso de cursos
// ==============================================

// Obtener el progreso de un usuario en un curso específico
func getProgresoUsuario(c *gin.Context) {
	// Obtener el usuario autenticado
	userValue, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	usuario, ok := userValue.(Usuario)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener información del usuario"})
		return
	}

	// Obtener ID del curso
	cursoID, err := strconv.ParseUint(c.Param("cursoId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de curso inválido"})
		return
	}

	// Verificar si el curso existe
	var curso Curso
	if result := db.First(&curso, cursoID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Curso no encontrado"})
		return
	}

	// Obtener progreso general del usuario en el curso
	var progresoUsuario ProgresoUsuario
	if result := db.Where("usuario_id = ? AND curso_id = ?", usuario.ID, cursoID).
		First(&progresoUsuario); result.Error != nil {
		// Si no existe, creamos un registro de progreso vacío
		progresoUsuario = ProgresoUsuario{
			UsuarioID:       usuario.ID,
			CursoID:         uint(cursoID),
			PorcentajeTotal: 0,
			UltimoCapitulo:  0,
		}
	}

	// Obtener progreso de cada capítulo
	var progresosCapitulos []ProgresoCapitulo
	if err := db.Where("usuario_id = ? AND curso_id = ?", usuario.ID, cursoID).
		Find(&progresosCapitulos).Error; err != nil {
		log.Printf("Error al obtener progreso de capítulos: %v", err)
	}

	// Convertir a mapa para facilitar el acceso en el frontend
	capitulosProgreso := make(map[uint]ProgresoCapitulo)
	for _, p := range progresosCapitulos {
		capitulosProgreso[p.CapituloID] = p
	}

	// Armar respuesta
	response := ProgresoResponse{
		ProgresoTotal:     progresoUsuario.PorcentajeTotal,
		UltimoCapitulo:    progresoUsuario.UltimoCapitulo,
		CapitulosProgreso: capitulosProgreso,
	}

	c.JSON(http.StatusOK, response)
}

// Marcar un capítulo como completado/incompleto
func marcarCapituloCompletado(c *gin.Context) {
	// Obtener el usuario autenticado
	userValue, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	usuario, ok := userValue.(Usuario)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener información del usuario"})
		return
	}

	// Obtener datos de la solicitud
	var req ProgresoCapituloRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	// Verificar que el curso y capítulo existan
	var curso Curso
	if result := db.First(&curso, req.CursoID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Curso no encontrado"})
		return
	}

	var capitulo Capitulo
	if result := db.First(&capitulo, req.CapituloID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Capítulo no encontrado"})
		return
	}

	// Verificar que el capítulo pertenezca al curso
	if capitulo.CursoID != req.CursoID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El capítulo no pertenece al curso indicado"})
		return
	}

	// Buscar si existe un registro de progreso para este capítulo
	var progresoCapitulo ProgresoCapitulo
	result := db.Where("usuario_id = ? AND curso_id = ? AND capitulo_id = ?",
		usuario.ID, req.CursoID, req.CapituloID).First(&progresoCapitulo)

	if result.Error != nil {
		// Si no existe, crear uno nuevo
		progresoCapitulo = ProgresoCapitulo{
			UsuarioID:  usuario.ID,
			CursoID:    req.CursoID,
			CapituloID: req.CapituloID,
			Completado: req.Completado,
			Progreso:   req.Progreso,
		}

		if err := db.Create(&progresoCapitulo).Error; err != nil {
			log.Printf("Error al crear progreso de capítulo: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar progreso"})
			return
		}
	} else {
		// Si existe, actualizar
		progresoCapitulo.Completado = req.Completado
		progresoCapitulo.Progreso = req.Progreso

		if err := db.Save(&progresoCapitulo).Error; err != nil {
			log.Printf("Error al actualizar progreso de capítulo: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar progreso"})
			return
		}
	}

	// Actualizar progreso total del curso
	actualizarProgresoTotal(usuario.ID, req.CursoID)

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "Progreso actualizado correctamente",
		"progreso": progresoCapitulo,
	})
}

// Guardar el último capítulo visto
func guardarUltimoCapitulo(c *gin.Context) {
	// Obtener el usuario autenticado
	userValue, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	usuario, ok := userValue.(Usuario)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener información del usuario"})
		return
	}

	// Obtener datos de la solicitud
	var req UltimoCapituloRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	// Verificar que el curso y capítulo existan
	var curso Curso
	if result := db.First(&curso, req.CursoID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Curso no encontrado"})
		return
	}

	var capitulo Capitulo
	if result := db.First(&capitulo, req.CapituloID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Capítulo no encontrado"})
		return
	}

	// Verificar que el capítulo pertenezca al curso
	if capitulo.CursoID != req.CursoID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El capítulo no pertenece al curso indicado"})
		return
	}

	// Buscar o crear el registro de progreso del usuario
	var progresoUsuario ProgresoUsuario
	result := db.Where("usuario_id = ? AND curso_id = ?", usuario.ID, req.CursoID).First(&progresoUsuario)

	if result.Error != nil {
		// Si no existe, crear uno nuevo
		progresoUsuario = ProgresoUsuario{
			UsuarioID:       usuario.ID,
			CursoID:         req.CursoID,
			UltimoCapitulo:  req.CapituloID,
			PorcentajeTotal: 0, // Se calculará en la siguiente función
		}

		if err := db.Create(&progresoUsuario).Error; err != nil {
			log.Printf("Error al crear progreso de usuario: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar progreso"})
			return
		}
	} else {
		// Si existe, actualizar
		progresoUsuario.UltimoCapitulo = req.CapituloID

		if err := db.Save(&progresoUsuario).Error; err != nil {
			log.Printf("Error al actualizar último capítulo: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar progreso"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Último capítulo actualizado correctamente",
	})
}

// Función para actualizar el progreso total de un curso
func actualizarProgresoTotal(usuarioID uint, cursoID uint) {
	// Obtener el total de capítulos del curso
	var totalCapitulos int64
	if err := db.Model(&Capitulo{}).Where("curso_id = ?", cursoID).Count(&totalCapitulos).Error; err != nil {
		log.Printf("Error al contar capítulos del curso: %v", err)
		return
	}

	if totalCapitulos == 0 {
		return // No hay capítulos para calcular progreso
	}

	// Obtener cantidad de capítulos completados
	var completados int64
	if err := db.Model(&ProgresoCapitulo{}).
		Where("usuario_id = ? AND curso_id = ? AND completado = ?",
			usuarioID, cursoID, true).
		Count(&completados).Error; err != nil {
		log.Printf("Error al contar capítulos completados: %v", err)
		return
	}

	// Calcular porcentaje
	porcentaje := float64(completados) / float64(totalCapitulos) * 100

	// Actualizar o crear registro de progreso del usuario
	var progresoUsuario ProgresoUsuario
	result := db.Where("usuario_id = ? AND curso_id = ?", usuarioID, cursoID).First(&progresoUsuario)

	if result.Error != nil {
		// Si no existe, crear uno nuevo
		progresoUsuario = ProgresoUsuario{
			UsuarioID:       usuarioID,
			CursoID:         cursoID,
			PorcentajeTotal: porcentaje,
		}

		if err := db.Create(&progresoUsuario).Error; err != nil {
			log.Printf("Error al crear progreso de usuario: %v", err)
		}
	} else {
		// Si existe, actualizar
		progresoUsuario.PorcentajeTotal = porcentaje

		if err := db.Save(&progresoUsuario).Error; err != nil {
			log.Printf("Error al actualizar progreso total: %v", err)
		}
	}
}
