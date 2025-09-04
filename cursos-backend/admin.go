package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Obtener estadísticas básicas para el panel de administración


// checkAdmin verifica si el usuario autenticado es administrador


// getAdminDashboard obtiene datos detallados para el dashboard del administrador


// Historial de actividades administrativas


// getSalesStats obtiene estadísticas detalladas de ventas para el panel de administración

// Listar todos los usuarios (solo admin)
func listUsers(c *gin.Context) {
	var users []Usuario
	query := db.Order("created_at desc")

	// Paginación
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	// Filtrado por rol
	role := c.Query("role")
	if role != "" {
		query = query.Where("role = ?", role)
	}

	// Búsqueda por nombre o email
	search := c.Query("search")
	if search != "" {
		query = query.Where("nombre LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Contar total de registros para paginación
	var total int64
	if err := query.Model(&Usuario{}).Count(&total).Error; err != nil {
		SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Obtener usuarios con paginación
	if err := query.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// No enviar contraseñas
	for i := range users {
		users[i].Password = ""
	}

	// Registrar actividad
	userValue, _ := c.Get("user")
	user := userValue.(Usuario)
	logActivity(c, user.ID, "list_users", 
		fmt.Sprintf("Admin visualizó lista de usuarios (total: %d, página: %d)", len(users), page))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"users": users,
			"pagination": gin.H{
				"total":  total,
				"page":   page,
				"limit":  limit,
				"pages":  (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// Obtener usuario por ID (solo admin)
func getUserById(c *gin.Context) {
	id := c.Param("id")

	var user Usuario
	if err := db.First(&user, id).Error; err != nil {
		SendErrorResponse(c, ErrResourceNotFound, http.StatusNotFound)
		return
	}

	// No enviar contraseña
	user.Password = ""

	// Registrar actividad
	currentUser, _ := c.Get("user")
	adminUser := currentUser.(Usuario)
	logActivity(c, adminUser.ID, "view_user_details", 
		fmt.Sprintf("Admin consultó detalles del usuario ID: %s, Email: %s", id, user.Email))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

// Actualizar usuario por ID (solo admin)
func updateUser(c *gin.Context) {
	id := c.Param("id")

	var user Usuario
	if err := db.First(&user, id).Error; err != nil {
		SendErrorResponse(c, ErrResourceNotFound, http.StatusNotFound)
		return
	}

	var req struct {
		Nombre string `json:"nombre"`
		Email  string `json:"email"`
		Phone  string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		SendErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	// Verificar si el nuevo email ya existe (si se cambia)
	if req.Email != "" && req.Email != user.Email {
		var existingUser Usuario
		if result := db.Where("email = ? AND id != ?", req.Email, id).First(&existingUser); result.Error == nil {
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

	// Guardar datos anteriores para el log
	emailAnterior := user.Email
	nombreAnterior := user.Nombre

	// Actualizar usuario
	if err := db.Model(&user).Updates(updates).Error; err != nil {
		SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Obtener usuario actualizado
	if err := db.First(&user, id).Error; err != nil {
		SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// No enviar contraseña
	user.Password = ""

	// Registrar actividad
	currentUser, _ := c.Get("user")
	adminUser := currentUser.(Usuario)
	logActivity(c, adminUser.ID, "update_user", 
		fmt.Sprintf("Admin actualizó usuario ID: %s, Email: '%s' -> '%s', Nombre: '%s' -> '%s'", 
			id, emailAnterior, user.Email, nombreAnterior, user.Nombre))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Usuario actualizado correctamente",
		"data":    user,
	})
}

// Eliminar usuario (solo admin)
func deleteUser(c *gin.Context) {
	id := c.Param("id")

	// Verificar si el usuario existe
	var user Usuario
	if err := db.First(&user, id).Error; err != nil {
		SendErrorResponse(c, ErrResourceNotFound, http.StatusNotFound)
		return
	}

	// Evitar que un admin elimine su propia cuenta
	currentUser, _ := c.Get("user")
	adminUser := currentUser.(Usuario)
	if adminUser.ID == user.ID {
		SendErrorResponse(c, errors.New("no puedes eliminar tu propia cuenta"), http.StatusBadRequest)
		return
	}

	// Guardar datos para el log
	userEmail := user.Email
	userName := user.Nombre

	// Eliminar usuario (borrado en cascada configurado en la base de datos)
	if err := db.Delete(&user).Error; err != nil {
		SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Registrar actividad
	logActivity(c, adminUser.ID, "delete_user", 
		fmt.Sprintf("Admin eliminó usuario ID: %s, Email: %s, Nombre: %s", id, userEmail, userName))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Usuario eliminado correctamente",
	})
}

// Estructura para cambio de rol
type ChangeRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

// Cambiar rol de usuario (solo admin)
func changeUserRole(c *gin.Context) {
	id := c.Param("id")

	// Verificar si el usuario existe
	var user Usuario
	if err := db.First(&user, id).Error; err != nil {
		SendErrorResponse(c, ErrResourceNotFound, http.StatusNotFound)
		return
	}

	var req ChangeRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	// Validar que el rol sea válido
	if req.Role != "admin" && req.Role != "user" {
		SendErrorResponse(c, errors.New("rol no válido"), http.StatusBadRequest)
		return
	}

	// Un admin no debería poder quitarse los privilegios a sí mismo
	currentUser, _ := c.Get("user")
	adminUser := currentUser.(Usuario)
	if adminUser.ID == user.ID && req.Role != "admin" {
		SendErrorResponse(c, errors.New("no puedes quitarte los privilegios de administrador a ti mismo"), http.StatusBadRequest)
		return
	}

	// Guardar rol anterior para el log
	rolAnterior := user.Role

	// Actualizar rol del usuario
	if err := db.Model(&user).Update("role", req.Role).Error; err != nil {
		SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Obtener usuario actualizado
	if err := db.First(&user, id).Error; err != nil {
		SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// No enviar contraseña
	user.Password = ""

	// Registrar actividad
	logActivity(c, adminUser.ID, "change_user_role", 
		fmt.Sprintf("Admin cambió rol de usuario ID: %s, Email: %s, Rol: '%s' -> '%s'", 
			id, user.Email, rolAnterior, user.Role))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Rol de usuario actualizado correctamente",
		"data":    user,
	})
}

// Función para registrar actividad en el sistema
func logActivity(c *gin.Context, userID uint, action, details string) {
	// Obtener IP y User-Agent
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	
	activityLog := ActivityLog{
		UserID:    userID,
		Action:    action,
		Details:   details,
		IP:        ip,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
	}
	
	if err := db.Create(&activityLog).Error; err != nil {
		log.Printf("Error al registrar actividad: %v", err)
	}
}