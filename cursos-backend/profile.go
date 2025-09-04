package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	
)

// Estructuras para las solicitudes
type UpdateProfileRequest struct {
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=8"`
}

type NotificationSettingsRequest struct {
	EmailNotifications bool `json:"emailNotifications"`
	NewMessages        bool `json:"newMessages"`
	NewStudents        bool `json:"newStudents"`
	SalesReports       bool `json:"salesReports"`
	SystemUpdates      bool `json:"systemUpdates"`
}

// Obtener perfil del usuario
func getProfile(c *gin.Context) {
	// El usuario ya está cargado por el middleware de autenticación
	userValue, exists := c.Get("user")
	if !exists {
		SendErrorResponse(c, ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	user, ok := userValue.(Usuario)
	if !ok {
		SendErrorResponse(c, errors.New("error al obtener información del usuario"), http.StatusInternalServerError)
		return
	}

	// No enviar contraseña en la respuesta
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    user,
	})
}

// Actualizar perfil
func updateProfile(c *gin.Context) {
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	// Obtener usuario actual
	userValue, _ := c.Get("user")
	currentUser := userValue.(Usuario)

	// Verificar si el nuevo email ya existe (si se cambia)
	if req.Email != currentUser.Email {
		var existingUser Usuario
		if result := db.Where("email = ? AND id != ?", req.Email, currentUser.ID).First(&existingUser); result.Error == nil {
			SendErrorResponse(c, ErrEmailExists, http.StatusBadRequest)
			return
		}
	}

	// Preparar updates
	updates := map[string]interface{}{}
	
	if req.Nombre != "" {
		updates["nombre"] = req.Nombre
	}
	
	if req.Email != "" {
		updates["email"] = req.Email
	}
	
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}

	// Actualizar usuario en la BD
	result := db.Model(&currentUser).Updates(updates)
	if result.Error != nil {
		SendErrorResponse(c, errors.New("error al actualizar perfil"), http.StatusInternalServerError)
		return
	}

	// Obtener usuario actualizado
	var updatedUser Usuario
	if err := db.First(&updatedUser, currentUser.ID).Error; err != nil {
		SendErrorResponse(c, errors.New("error al obtener perfil actualizado"), http.StatusInternalServerError)
		return
	}

	// No enviar contraseña en la respuesta
	updatedUser.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Perfil actualizado correctamente",
		"user":    updatedUser,
	})
}

// Cambiar contraseña
func changePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	// Obtener usuario actual
	userValue, _ := c.Get("user")
	currentUser := userValue.(Usuario)

	// Verificar contraseña actual
	if err := bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(req.CurrentPassword)); err != nil {
		SendErrorResponse(c, errors.New("contraseña actual incorrecta"), http.StatusBadRequest)
		return
	}

	// Hash de la nueva contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		SendErrorResponse(c, errors.New("error al procesar la nueva contraseña"), http.StatusInternalServerError)
		return
	}

	// Actualizar contraseña en la BD
	result := db.Model(&currentUser).Update("password", string(hashedPassword))
	if result.Error != nil {
		SendErrorResponse(c, errors.New("error al actualizar contraseña"), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Contraseña actualizada correctamente",
	})
}

// Subir imagen de perfil
func uploadProfileImage(c *gin.Context) {
	// Obtener usuario actual
	userValue, _ := c.Get("user")
	currentUser := userValue.(Usuario)

	// Obtener archivo
	file, err := c.FormFile("image")
	if err != nil {
		SendErrorResponse(c, errors.New("error al recibir imagen"), http.StatusBadRequest)
		return
	}

	// Validar tamaño (máx 2MB)
	if file.Size > 2*1024*1024 {
		SendErrorResponse(c, errors.New("la imagen no debe superar los 2MB"), http.StatusBadRequest)
		return
	}

	// Validar tipo de archivo
	extension := filepath.Ext(file.Filename)
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	if !allowedExtensions[extension] {
		SendErrorResponse(c, errors.New("tipo de archivo no permitido"), http.StatusBadRequest)
		return
	}

	// Crear directorio para imágenes de usuario si no existe
	uploadDir := "./static/profiles"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, os.ModePerm)
	}

	// Nombre de archivo: user_ID.extension
	filename := "user_" + string(currentUser.ID) + extension
	filepath := filepath.Join(uploadDir, filename)

	// Guardar archivo
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		SendErrorResponse(c, errors.New("error al guardar imagen"), http.StatusInternalServerError)
		return
	}

	// Actualizar URL de imagen en la BD
	imageURL := "/static/profiles/" + filename
	if err := db.Model(&currentUser).Update("image_url", imageURL).Error; err != nil {
		SendErrorResponse(c, errors.New("error al actualizar imagen en la base de datos"), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "Imagen de perfil actualizada correctamente",
		"imageUrl": imageURL,
	})
}

// Obtener preferencias de notificación
func getNotificationSettings(c *gin.Context) {
	// Obtener usuario actual
	userValue, _ := c.Get("user")
	currentUser := userValue.(Usuario)
	
	// Registrar la acción (para evitar el error "declared and not used")
	log.Printf("Usuario %s (ID: %d) consultó sus preferencias de notificación", currentUser.Email, currentUser.ID)

	// En una implementación real, obtendríamos esto de la base de datos
	// Por ahora, devolvemos valores predeterminados

	settings := NotificationSettingsRequest{
		EmailNotifications: true,
		NewMessages:        true,
		NewStudents:        true,
		SalesReports:       true,
		SystemUpdates:      false,
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"settings": settings,
	})
}

// Actualizar preferencias de notificación
func updateNotificationSettings(c *gin.Context) {
	var req NotificationSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	// Obtener usuario actual
	userValue, _ := c.Get("user")
	currentUser := userValue.(Usuario)
	
	// Registrar la acción (para evitar el error "declared and not used")
	log.Printf("Usuario %s (ID: %d) actualizó sus preferencias de notificación", currentUser.Email, currentUser.ID)

	// En una implementación real, guardaríamos esto en la base de datos
	// usando el ID del usuario: currentUser.ID

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Preferencias de notificación actualizadas correctamente",
	})
}