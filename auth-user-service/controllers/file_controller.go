package controllers

import (
	"auth-user-service/models"
	"auth-user-service/services"
	"auth-user-service/utils"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type FileController struct {
	fileService *services.FileService
}

func NewFileController() *FileController {
	return &FileController{
		fileService: services.NewFileService(),
	}
}

// UploadProfileImage sube una imagen de perfil
func (f *FileController) UploadProfileImage(c *gin.Context) {
	// Obtener el usuario del contexto
	userValue, exists := c.Get("user")
	if !exists {
		utils.SendErrorResponse(c, utils.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	user, ok := userValue.(*models.Usuario)
	if !ok {
		utils.SendErrorMessage(c, "error al obtener información del usuario", http.StatusInternalServerError)
		return
	}

	// Obtener el archivo del formulario
	file, err := c.FormFile("image")
	if err != nil {
		utils.SendErrorMessage(c, "No se proporcionó ningún archivo", http.StatusBadRequest)
		return
	}

	// Subir el archivo
	imageURL, err := f.fileService.UploadProfileImage(user.ID, file)
	if err != nil {
		switch err {
		case utils.ErrInvalidFileType:
			utils.SendErrorResponse(c, err, http.StatusBadRequest)
		case utils.ErrFileTooLarge:
			utils.SendErrorResponse(c, err, http.StatusBadRequest)
		default:
			utils.SendErrorMessage(c, "Error al subir la imagen", http.StatusInternalServerError)
		}
		return
	}

	utils.SendSuccessResponse(c, gin.H{
		"message":   "Imagen subida correctamente",
		"image_url": imageURL,
	})
}

// ServeProfileImage sirve imágenes de perfil
func (f *FileController) ServeProfileImage(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		utils.SendErrorMessage(c, "Nombre de archivo requerido", http.StatusBadRequest)
		return
	}

	// Construir ruta del archivo
	filePath := filepath.Join("./static/profiles", filename)

	// Verificar si el archivo existe
	if _, err := filepath.Abs(filePath); err != nil {
		utils.SendErrorMessage(c, "Archivo no encontrado", http.StatusNotFound)
		return
	}

	// Configurar headers de cache
	c.Header("Cache-Control", "public, max-age=31536000") // 1 año
	c.Header("Content-Type", "image/*")

	// Servir el archivo
	c.File(filePath)
}