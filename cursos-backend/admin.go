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
func getAdminStats(c *gin.Context) {
	// Obtener el usuario actual para registro de actividad
	userValue, _ := c.Get("user")
	user := userValue.(Usuario)
	
	// Registrar actividad
	logActivity(c, user.ID, "view_stats", "Visualización de estadísticas del panel admin")
	
	// Estadísticas de usuarios
	var totalUsers int64
	db.Model(&Usuario{}).Count(&totalUsers)
	
	// Usuarios nuevos en el último mes
	var newUsers int64
	oneMonthAgo := time.Now().AddDate(0, -1, 0)
	db.Model(&Usuario{}).Where("created_at >= ?", oneMonthAgo).Count(&newUsers)
	
	// Estadísticas de cursos
	var totalCursos int64
	db.Model(&Curso{}).Count(&totalCursos)
	
	var cursosPublicados int64
	db.Model(&Curso{}).Where("estado = ?", "Publicado").Count(&cursosPublicados)
	
	// Estadísticas de ventas/pagos
	var ingresosMensuales float64
	primerDiaMes := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)
	db.Model(&Pago{}).Where("created_at >= ? AND estado = ?", primerDiaMes, "completado").
		Select("COALESCE(SUM(monto), 0)").Scan(&ingresosMensuales)
	
	var totalPagos int64
	db.Model(&Pago{}).Where("estado = ?", "completado").Count(&totalPagos)
	
	// Respuesta
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"usuarios": gin.H{
				"total": totalUsers,
				"nuevos": newUsers,
			},
			"cursos": gin.H{
				"total": totalCursos,
				"publicados": cursosPublicados,
			},
			"ventas": gin.H{
				"ingresosMensuales": ingresosMensuales,
				"totalTransacciones": totalPagos,
			},
		},
	})
}

// checkAdmin verifica si el usuario autenticado es administrador


// getAdminDashboard obtiene datos detallados para el dashboard del administrador
func getAdminDashboard(c *gin.Context) {
	// Obtener usuario para registro de actividad
	userValue, _ := c.Get("user")
	user := userValue.(Usuario)
	
	// Registrar actividad
	logActivity(c, user.ID, "view_dashboard", "Acceso al panel principal de administración")
	
	// Estadísticas de usuarios
	var totalUsers int64
	if err := db.Model(&Usuario{}).Count(&totalUsers).Error; err != nil {
		SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Usuarios nuevos en la última semana
	oneWeekAgo := time.Now().AddDate(0, 0, -7)
	var newLastWeek int64
	if err := db.Model(&Usuario{}).Where("created_at >= ?", oneWeekAgo).Count(&newLastWeek).Error; err != nil {
		SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Usuarios nuevos en el último mes
	oneMonthAgo := time.Now().AddDate(0, -1, 0)
	var newLastMonth int64
	if err := db.Model(&Usuario{}).Where("created_at >= ?", oneMonthAgo).Count(&newLastMonth).Error; err != nil {
		SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Usuarios activos (login en la última semana)
	var activeUsers int64
	if err := db.Model(&Usuario{}).Where("last_login >= ?", oneWeekAgo).Count(&activeUsers).Error; err != nil {
		SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Estadísticas de cursos
	var totalCourses int64
	if err := db.Model(&Curso{}).Count(&totalCourses).Error; err != nil {
		SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	var publishedCourses int64
	if err := db.Model(&Curso{}).Where("estado = ?", "Publicado").Count(&publishedCourses).Error; err != nil {
		SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	var draftCourses int64
	if err := db.Model(&Curso{}).Where("estado = ?", "Borrador").Count(&draftCourses).Error; err != nil {
		SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Precio promedio de cursos
	var avgPrice float64
	if err := db.Model(&Curso{}).Select("COALESCE(AVG(precio), 0) as avg_price").Scan(&avgPrice).Error; err != nil {
		SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Estadísticas de ventas
	var totalSales float64
	if err := db.Model(&Pago{}).Where("estado = ?", "completado").Select("COALESCE(SUM(monto), 0) as total_sales").Scan(&totalSales).Error; err != nil {
		SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Obtener ventas mensuales (últimos 6 meses)
	var monthlySales []float64
	var monthlyLabels []string

	for i := 5; i >= 0; i-- {
		startDate := time.Now().AddDate(0, -i, 0).Format("2006-01")
		endDate := time.Now().AddDate(0, -i+1, 0).Format("2006-01")
		
		var monthSales float64
		if err := db.Model(&Pago{}).Where("estado = ? AND created_at >= ? AND created_at < ?", 
			"completado", startDate, endDate).Select("COALESCE(SUM(monto), 0) as month_sales").Scan(&monthSales).Error; err != nil {
			SendErrorResponse(c, err, http.StatusInternalServerError)
			return
		}
		
		monthlySales = append(monthlySales, monthSales)
		monthlyLabels = append(monthlyLabels, time.Now().AddDate(0, -i, 0).Format("Jan"))
	}

	// Conteo de ventas recientes (última semana)
	var recentSales int64
	if err := db.Model(&Pago{}).Where("created_at >= ?", oneWeekAgo).Count(&recentSales).Error; err != nil {
		SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Conteo de ventas pendientes y completadas
	var pendingSales, completedSales int64
	if err := db.Model(&Pago{}).Where("estado = ?", "pendiente").Count(&pendingSales).Error; err != nil {
		SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}
	if err := db.Model(&Pago{}).Where("estado = ?", "completado").Count(&completedSales).Error; err != nil {
		SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Obtener usuarios recientes
	var recentUsers []Usuario
	if err := db.Order("created_at desc").Limit(5).Find(&recentUsers).Error; err != nil {
		SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// No enviar contraseñas
	for i := range recentUsers {
		recentUsers[i].Password = ""
	}

	// Construir respuesta completa
	dashboardResponse := gin.H{
		"userStats": gin.H{
			"total":        int(totalUsers),
			"newLastWeek":  int(newLastWeek),
			"newLastMonth": int(newLastMonth),
			"active":       int(activeUsers),
		},
		"courseStats": gin.H{
			"total":     int(totalCourses),
			"published": int(publishedCourses),
			"draft":     int(draftCourses),
			"avgPrice":  avgPrice,
		},
		"salesStats": gin.H{
			"totalSales":     totalSales,
			"monthlySales":   monthlySales,
			"monthlyLabels":  monthlyLabels,
			"recentSales":    int(recentSales),
			"pendingSales":   int(pendingSales),
			"completedSales": int(completedSales),
		},
		"recentUsers": recentUsers,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    dashboardResponse,
	})
}

// Historial de actividades administrativas
func getActivityLog(c *gin.Context) {
	// Obtener usuario para registro de actividad
	userValue, _ := c.Get("user")
	user := userValue.(Usuario)

	// Registrar esta actividad
	logActivity(c, user.ID, "view_activity_log", "Consulta al registro de actividades del sistema")

	// Paginación
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset := (page - 1) * limit

	// Construir query
	query := db.Model(&ActivityLog{}).Order("created_at desc")

	// Filtros
	userID := c.Query("user_id")
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	action := c.Query("action")
	if action != "" {
		query = query.Where("action = ?", action)
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if startDate != "" && endDate != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate+" 23:59:59")
	} else if startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	} else if endDate != "" {
		query = query.Where("created_at <= ?", endDate+" 23:59:59")
	}

	// Total para paginación
	var total int64
	if err := query.Count(&total).Error; err != nil {
		SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Obtener logs con paginación
	var logs []ActivityLog
	if err := query.Limit(limit).Offset(offset).Find(&logs).Error; err != nil {
		SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Obtener información adicional de los usuarios
	type LogWithUserInfo struct {
		Log         ActivityLog `json:"log"`
		UserName    string      `json:"userName"`
		UserEmail   string      `json:"userEmail"`
	}

	var enhancedLogs []LogWithUserInfo

	for _, log := range logs {
		var userName, userEmail string

		// Obtener información del usuario
		var logUser Usuario
		if err := db.Select("nombre, email").First(&logUser, log.UserID).Error; err == nil {
			userName = logUser.Nombre
			userEmail = logUser.Email
		}

		enhancedLogs = append(enhancedLogs, LogWithUserInfo{
			Log:         log,
			UserName:    userName,
			UserEmail:   userEmail,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"logs": enhancedLogs,
			"pagination": gin.H{
				"total": total,
				"page":  page,
				"limit": limit,
				"pages": (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// getSalesStats obtiene estadísticas detalladas de ventas para el panel de administración
func getSalesStats(c *gin.Context) {
	// Obtener usuario para registro de actividad
	userValue, _ := c.Get("user")
	user := userValue.(Usuario)

	// Registrar actividad
	logActivity(c, user.ID, "view_sales_stats", "Consulta de estadísticas de ventas")

	// Periodo (últimos 30 días por defecto, o personalizado)
	startDateStr := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDateStr := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		SendErrorResponse(c, errors.New("formato de fecha de inicio inválido"), http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		SendErrorResponse(c, errors.New("formato de fecha de fin inválido"), http.StatusBadRequest)
		return
	}
	
	// Ajustar endDate para incluir todo el día
	endDate = endDate.Add(24 * time.Hour - time.Second)

	// Ventas totales en el periodo
	var totalSales float64
	if err := db.Model(&Pago{}).Where("created_at BETWEEN ? AND ? AND estado = ?", 
		startDate, endDate, "completado").Select("COALESCE(SUM(monto), 0)").Scan(&totalSales).Error; err != nil {
		SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Número de transacciones en el periodo
	var totalTransactions int64
	if err := db.Model(&Pago{}).Where("created_at BETWEEN ? AND ?", 
		startDate, endDate).Count(&totalTransactions).Error; err != nil {
		SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Valor promedio de transacción
	var avgTransaction float64
	if totalTransactions > 0 {
		avgTransaction = totalSales / float64(totalTransactions)
	}

	// Datos por día
	type DailySales struct {
		Date  string  `json:"date"`
		Sales float64 `json:"sales"`
		Count int     `json:"count"`
	}

	// Calcular el número de días entre las fechas
	days := int(endDate.Sub(startDate).Hours() / 24) + 1
	
	// Preparar slice para datos diarios
	dailyData := make([]DailySales, 0, days)
	
	// Para cada día en el rango
	for d := 0; d < days; d++ {
		currentDate := startDate.AddDate(0, 0, d)
		nextDate := currentDate.AddDate(0, 0, 1)
		
		// Ventas del día
		var daySales float64
		if err := db.Model(&Pago{}).Where("created_at >= ? AND created_at < ? AND estado = ?", 
			currentDate, nextDate, "completado").Select("COALESCE(SUM(monto), 0)").Scan(&daySales).Error; err != nil {
			SendErrorResponse(c, err, http.StatusInternalServerError)
			return
		}
		
		// Número de transacciones del día
		var dayTransactions int64
		if err := db.Model(&Pago{}).Where("created_at >= ? AND created_at < ?", 
			currentDate, nextDate).Count(&dayTransactions).Error; err != nil {
			SendErrorResponse(c, err, http.StatusInternalServerError)
			return
		}
		
		dailyData = append(dailyData, DailySales{
			Date:  currentDate.Format("2006-01-02"),
			Sales: daySales,
			Count: int(dayTransactions),
		})
	}

	// Top 5 cursos más vendidos
	type TopCourse struct {
		ID      uint    `json:"id"`
		Titulo  string  `json:"titulo"`
		Ventas  int     `json:"ventas"`
		Ingreso float64 `json:"ingreso"`
	}

	var topCourses []TopCourse
	if err := db.Model(&Pago{}).
		Select("curso_id as id, COUNT(*) as ventas, SUM(monto) as ingreso").
		Where("created_at BETWEEN ? AND ? AND estado = ?", startDate, endDate, "completado").
		Group("curso_id").
		Order("ventas DESC").
		Limit(5).
		Scan(&topCourses).Error; err != nil {
		SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Obtener títulos de los cursos
	for i, course := range topCourses {
		var curso Curso
		if err := db.Select("titulo").First(&curso, course.ID).Error; err == nil {
			topCourses[i].Titulo = curso.Titulo
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"period": gin.H{
				"startDate": startDate.Format("2006-01-02"),
				"endDate":   endDate.Format("2006-01-02"),
				"days":      days,
			},
			"summary": gin.H{
				"totalSales":       totalSales,
				"totalTransactions": totalTransactions,
				"avgTransaction":   avgTransaction,
			},
			"dailyData":  dailyData,
			"topCourses": topCourses,
		},
	})
}

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