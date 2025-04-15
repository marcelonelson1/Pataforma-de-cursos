package main

import (
	"net/http"
	"fmt"
	"time"
	"log"
	"github.com/gin-gonic/gin"
	"strconv"
)

// Estructuras para respuestas de estadísticas

// StatValue representa un par de valores (actual y anterior) para comparar
type StatValue struct {
	Current  interface{} `json:"current"`
	Previous interface{} `json:"previous"`
}

// AdminStatsResponse estructura para la respuesta de estadísticas generales
type AdminStatsResponse struct {
	Stats struct {
		ActiveStudents   StatValue `json:"activeStudents"`
		PublishedCourses StatValue `json:"publishedCourses"`
		MonthlyRevenue   StatValue `json:"monthlyRevenue"`
		AverageRating    StatValue `json:"averageRating"`
	} `json:"stats"`
	Period string `json:"period"`
}

// CourseSale representa las ventas de un curso individual
type CourseSale struct {
	Name       string  `json:"name"`
	Sales      int     `json:"sales"`
	Revenue    float64 `json:"revenue"`
	Percentage int     `json:"percentage"`
}

// MonthlyData representa datos de ventas y usuarios por mes
type MonthlyData struct {
	Month string `json:"month"`
	Sales int    `json:"sales"`
	Users int    `json:"users"`
}

// UserStats representa estadísticas de usuarios
type UserStats struct {
	New       int `json:"new"`
	Returning int `json:"returning"`
	Premium   int `json:"premium"`
}

// PaymentMethods representa la distribución de métodos de pago
type PaymentMethods struct {
	Paypal    int `json:"paypal"`
	Card      int `json:"card"`
	Transfer  int `json:"transfer"`
}

// SalesStatsResponse estructura para la respuesta de estadísticas de ventas
type SalesStatsResponse struct {
	CoursesSales   []CourseSale   `json:"coursesSales"`
	MonthlyData    []MonthlyData  `json:"monthlyData"`
	UserStats      UserStats      `json:"userStats"`
	PaymentMethods PaymentMethods `json:"paymentMethods"`
	Period         string         `json:"period"`
}

// getAdminStats devuelve estadísticas generales para el panel de administración
func getAdminStats(c *gin.Context) {
	// Obtener el período de tiempo desde los parámetros de la consulta
	period := c.DefaultQuery("period", "month")
	
	// Calcular los rangos de fechas según el período
	startDate, endDate, prevStartDate, prevEndDate := calculateDateRanges(period)
	
	// Consultar estadísticas actuales
	activeStudents, err := getActiveStudentsCount(startDate, endDate)
	if err != nil {
		log.Printf("Error al obtener estudiantes activos: %v", err)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}
	
	// Consultar estadísticas anteriores para comparación
	prevActiveStudents, err := getActiveStudentsCount(prevStartDate, prevEndDate)
	if err != nil {
		log.Printf("Error al obtener estudiantes activos previos: %v", err)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}
	
	// Consultar cursos publicados
	publishedCourses, err := getPublishedCoursesCount(startDate, endDate)
	if err != nil {
		log.Printf("Error al obtener cursos publicados: %v", err)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}
	
	prevPublishedCourses, err := getPublishedCoursesCount(prevStartDate, prevEndDate)
	if err != nil {
		log.Printf("Error al obtener cursos publicados previos: %v", err)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}
	
	// Consultar ingresos mensuales
	monthlyRevenue, err := getMonthlyRevenue(startDate, endDate)
	if err != nil {
		log.Printf("Error al obtener ingresos mensuales: %v", err)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}
	
	prevMonthlyRevenue, err := getMonthlyRevenue(prevStartDate, prevEndDate)
	if err != nil {
		log.Printf("Error al obtener ingresos mensuales previos: %v", err)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}
	
	// Consultar valoración media
	averageRating, err := getAverageRating(startDate, endDate)
	if err != nil {
		log.Printf("Error al obtener valoración media: %v", err)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}
	
	prevAverageRating, err := getAverageRating(prevStartDate, prevEndDate)
	if err != nil {
		log.Printf("Error al obtener valoración media previa: %v", err)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}
	
	// Crear la respuesta
	response := AdminStatsResponse{}
	response.Stats.ActiveStudents = StatValue{Current: activeStudents, Previous: prevActiveStudents}
	response.Stats.PublishedCourses = StatValue{Current: publishedCourses, Previous: prevPublishedCourses}
	response.Stats.MonthlyRevenue = StatValue{Current: monthlyRevenue, Previous: prevMonthlyRevenue}
	response.Stats.AverageRating = StatValue{Current: averageRating, Previous: prevAverageRating}
	response.Period = period
	
	SendSuccessResponse(c, response)
}

// calculateDateRanges calcula los rangos de fechas para estadísticas según el período seleccionado
func calculateDateRanges(period string) (time.Time, time.Time, time.Time, time.Time) {
	now := time.Now()
	var startDate, endDate, prevStartDate, prevEndDate time.Time

	switch period {
	case "week":
		// Período actual: últimos 7 días
		endDate = now
		startDate = now.AddDate(0, 0, -7)
		// Período anterior: 7 días antes del período actual
		prevEndDate = startDate
		prevStartDate = startDate.AddDate(0, 0, -7)
	case "month":
		// Período actual: últimos 30 días
		endDate = now
		startDate = now.AddDate(0, 0, -30)
		// Período anterior: 30 días antes del período actual
		prevEndDate = startDate
		prevStartDate = startDate.AddDate(0, 0, -30)
	case "year":
		// Período actual: último año
		endDate = now
		startDate = now.AddDate(-1, 0, 0)
		// Período anterior: año anterior al período actual
		prevEndDate = startDate
		prevStartDate = startDate.AddDate(-1, 0, 0)
	default:
		// Por defecto, usar vista mensual
		endDate = now
		startDate = now.AddDate(0, 0, -30)
		prevEndDate = startDate
		prevStartDate = startDate.AddDate(0, 0, -30)
	}

	return startDate, endDate, prevStartDate, prevEndDate
}

// getSalesStats obtiene estadísticas detalladas de ventas para el panel de administración
func getSalesStats(c *gin.Context) {
	// Obtener el período de tiempo desde los parámetros de la consulta
	period := c.DefaultQuery("period", "month")
	
	// Calcular los rangos de fechas según el período
	startDate, endDate, _, _ := calculateDateRanges(period)
	
	// Consultar ventas por curso
	coursesSales, err := getCoursesSales(startDate, endDate)
	if err != nil {
		log.Printf("Error al obtener ventas por curso: %v", err)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}
	
	// Consultar datos mensuales
	monthlyData, err := getMonthlyData(startDate, endDate)
	if err != nil {
		log.Printf("Error al obtener datos mensuales: %v", err)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}
	
	// Consultar estadísticas de usuarios
	userStats, err := getUserStats(startDate, endDate)
	if err != nil {
		log.Printf("Error al obtener estadísticas de usuarios: %v", err)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}
	
	// Consultar métodos de pago
	paymentMethods, err := getPaymentMethods(startDate, endDate)
	if err != nil {
		log.Printf("Error al obtener métodos de pago: %v", err)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}
	
	// Crear la respuesta
	response := SalesStatsResponse{
		CoursesSales:   coursesSales,
		MonthlyData:    monthlyData,
		UserStats:      userStats,
		PaymentMethods: paymentMethods,
		Period:         period,
	}
	
	SendSuccessResponse(c, response)
}

// DashboardData representa los datos generales para el panel de control de administración
type DashboardData struct {
	TotalUsers       int     `json:"totalUsers"`
	TotalCourses     int     `json:"totalCourses"`
	TotalRevenue     float64 `json:"totalRevenue"`
	PendingPayments  int     `json:"pendingPayments"`
	RecentUsers      []struct {
		ID        uint      `json:"id"`
		Nombre    string    `json:"nombre"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"recentUsers"`
	RecentPayments []struct {
		ID            uint      `json:"id"`
		UsuarioNombre string    `json:"usuario_nombre"`
		CursoTitulo   string    `json:"curso_titulo"`
		Monto         float64   `json:"monto"`
		Estado        string    `json:"estado"`
		CreatedAt     time.Time `json:"created_at"`
	} `json:"recentPayments"`
}

// getAdminDashboard devuelve un resumen general para el panel de administración
func getAdminDashboard(c *gin.Context) {
	// Obtener datos del dashboard
	dashboardData, err := fetchDashboardData()
	if err != nil {
		log.Printf("Error al obtener datos para el dashboard: %v", err)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}
	
	SendSuccessResponse(c, dashboardData)
}

// Estructura para los datos del registro de actividad
type ActivityLogData struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	UserName  string    `json:"user_name"`
	Action    string    `json:"action"`
	Details   string    `json:"details"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

// getActivityLog devuelve el registro de actividad para el panel de administración
func getActivityLog(c *gin.Context) {
	// Obtener parámetros de paginación
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "50")
	
	// Obtener parámetros de filtrado
	userID := c.Query("user_id")
	action := c.Query("action")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	
	// Consultar registro de actividad
	activityLog, totalEntries, err := fetchActivityLog(page, limit, userID, action, startDate, endDate)
	if err != nil {
		log.Printf("Error al obtener registro de actividad: %v", err)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}
	
	// Crear la respuesta
	response := gin.H{
		"activityLog":  activityLog,
		"totalEntries": totalEntries,
		"currentPage":  page,
		"limit":        limit,
	}
	
	SendSuccessResponse(c, response)
}

// Funciones auxiliares para consultar datos estadísticos reales de la base de datos

// getActiveStudentsCount obtiene el número de estudiantes activos en un rango de fechas
func getActiveStudentsCount(startDate, endDate time.Time) (int, error) {
	// Contamos los usuarios con progreso registrado en el rango de fechas
	var count int64
	
	// Contar usuarios distintos con actividad de progreso en el período
	err := db.Model(&ProgresoUsuario{}).
		Where("updated_at BETWEEN ? AND ?", startDate, endDate).
		Distinct("usuario_id").
		Count(&count).Error
		
	if err != nil {
		return 0, fmt.Errorf("error al contar estudiantes activos: %v", err)
	}
	
	return int(count), nil
}

// getPublishedCoursesCount obtiene el número de cursos publicados en un rango de fechas
func getPublishedCoursesCount(startDate, endDate time.Time) (int, error) {
	var count int64
	
	// Contar cursos publicados en el período
	err := db.Model(&Curso{}).
		Where("estado = 'Publicado' AND updated_at BETWEEN ? AND ?", startDate, endDate).
		Count(&count).Error
		
	if err != nil {
		return 0, fmt.Errorf("error al contar cursos publicados: %v", err)
	}
	
	return int(count), nil
}

// getMonthlyRevenue obtiene los ingresos en un rango de fechas
func getMonthlyRevenue(startDate, endDate time.Time) (float64, error) {
	var totalRevenue float64
	
	// Consultar la suma de los montos de los pagos aprobados en el período
	err := db.Model(&Pago{}).
		Select("COALESCE(SUM(monto), 0) as total").
		Where("estado = 'aprobado' AND created_at BETWEEN ? AND ?", startDate, endDate).
		Scan(&totalRevenue).Error
		
	if err != nil {
		return 0, fmt.Errorf("error al calcular ingresos mensuales: %v", err)
	}
	
	return totalRevenue, nil
}

// getAverageRating obtiene la valoración media de los cursos en un rango de fechas
// Nota: Como no se tiene una tabla de valoraciones definida en los archivos proporcionados,
// crearemos una consulta que podría funcionar si existiera dicha tabla.
// En este caso, usaremos un valor predeterminado o se podría implementar cuando se tenga la tabla.
func getAverageRating(startDate, endDate time.Time) (float64, error) {
	// Como no hay tabla de valoraciones en los archivos proporcionados, 
	// podemos implementar una lógica alternativa o usar un valor predeterminado
	
	// Versión simulada - se podría implementar cuando se tenga la tabla de valoraciones
	var averageRating float64 = 4.5 // Valor predeterminado
	
	// La consulta real sería algo como:
	/*
	err := db.Table("valoraciones").
		Select("COALESCE(AVG(puntuacion), 0) as promedio").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Scan(&averageRating).Error
	
	if err != nil {
		return 0, fmt.Errorf("error al calcular valoración media: %v", err)
	}
	*/
	
	return averageRating, nil
}

// getCoursesSales obtiene las ventas por curso en un rango de fechas
func getCoursesSales(startDate, endDate time.Time) ([]CourseSale, error) {
	type ResultRow struct {
		CursoID   uint
		Nombre    string
		Ventas    int
		Ingresos  float64
	}
	
	var results []ResultRow
	
	// Consultar ventas agrupadas por curso
	err := db.Table("pagos").
		Select("pagos.curso_id, cursos.titulo as nombre, COUNT(pagos.id) as ventas, SUM(pagos.monto) as ingresos").
		Joins("JOIN cursos ON pagos.curso_id = cursos.id").
		Where("pagos.estado = 'aprobado' AND pagos.created_at BETWEEN ? AND ?", startDate, endDate).
		Group("pagos.curso_id, cursos.titulo").
		Order("ingresos DESC").
		Limit(10).
		Scan(&results).Error
		
	if err != nil {
		return nil, fmt.Errorf("error al obtener ventas por curso: %v", err)
	}
	
	// Calcular el total de ingresos para obtener porcentajes
	var totalRevenue float64 = 0
	for _, r := range results {
		totalRevenue += r.Ingresos
	}
	
	// Convertir resultados al formato de respuesta
	var coursesSales []CourseSale
	for _, r := range results {
		percentage := 0
		if totalRevenue > 0 {
			percentage = int((r.Ingresos / totalRevenue) * 100)
		}
		
		coursesSales = append(coursesSales, CourseSale{
			Name:       r.Nombre,
			Sales:      r.Ventas,
			Revenue:    r.Ingresos,
			Percentage: percentage,
		})
	}
	
	return coursesSales, nil
}

// getMonthlyData obtiene los datos mensuales de ventas y usuarios
func getMonthlyData(startDate, endDate time.Time) ([]MonthlyData, error) {
	// Determinar cuántos meses mostrar según el período
	monthsDiff := int(endDate.Sub(startDate).Hours() / 24 / 30) + 1
	if monthsDiff > 12 {
		monthsDiff = 12 // Limitar a máximo 12 meses
	}
	
	var monthlyData []MonthlyData
	months := []string{"Ene", "Feb", "Mar", "Abr", "May", "Jun", "Jul", "Ago", "Sep", "Oct", "Nov", "Dic"}
	
	currentMonth := int(time.Now().Month())
	
	// Generar datos para los últimos X meses
	for i := 0; i < monthsDiff; i++ {
		monthIndex := (currentMonth - i - 1 + 12) % 12 // Ajustar para array base 0 y meses anteriores
		
		// Calcular fechas para este mes
		year := time.Now().Year()
		if currentMonth-i <= 0 {
			year-- // Ajustar año si necesitamos meses del año anterior
		}
		
		monthStart := time.Date(year, time.Month((monthIndex+1)), 1, 0, 0, 0, 0, time.UTC)
		var monthEnd time.Time
		if i == 0 {
			monthEnd = time.Now() // Para el mes actual, usar la fecha actual como fin
		} else {
			// Para meses anteriores, usar el último día del mes
			if monthIndex == 11 { // Diciembre
				monthEnd = time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.UTC).Add(-time.Second)
			} else {
				monthEnd = time.Date(year, time.Month(monthIndex+2), 1, 0, 0, 0, 0, time.UTC).Add(-time.Second)
			}
		}
		
		// Obtener ventas del mes
		var sales int64
		if err := db.Model(&Pago{}).
			Where("estado = 'aprobado' AND created_at BETWEEN ? AND ?", monthStart, monthEnd).
			Count(&sales).Error; err != nil {
			return nil, fmt.Errorf("error al obtener ventas mensuales: %v", err)
		}
		
		// Obtener usuarios activos del mes
		var users int64
		if err := db.Model(&ProgresoUsuario{}).
			Where("updated_at BETWEEN ? AND ?", monthStart, monthEnd).
			Distinct("usuario_id").
			Count(&users).Error; err != nil {
			return nil, fmt.Errorf("error al obtener usuarios mensuales: %v", err)
		}
		
		monthlyData = append([]MonthlyData{
			{
				Month: months[monthIndex],
				Sales: int(sales),
				Users: int(users),
			},
		}, monthlyData...) // Insertar al inicio para mantener orden cronológico
	}
	
	return monthlyData, nil
}

// getUserStats obtiene estadísticas de usuarios
func getUserStats(startDate, endDate time.Time) (UserStats, error) {
	// Contadores para diferentes tipos de usuarios
	var newUsers, returningUsers, premiumUsers int64
	
	// Usuarios nuevos registrados en el período
	if err := db.Model(&Usuario{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&newUsers).Error; err != nil {
		return UserStats{}, fmt.Errorf("error al contar usuarios nuevos: %v", err)
	}
	
	// Usuarios que retornan (con login en el período pero creados antes)
	if err := db.Model(&Usuario{}).
		Where("last_login BETWEEN ? AND ? AND created_at < ?", startDate, endDate, startDate).
		Count(&returningUsers).Error; err != nil {
		return UserStats{}, fmt.Errorf("error al contar usuarios que retornan: %v", err)
	}
	
	// Usuarios premium (con pagos aprobados en el período)
	if err := db.Model(&Pago{}).
		Select("COUNT(DISTINCT usuario_id)").
		Where("estado = 'aprobado' AND created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&premiumUsers).Error; err != nil {
		return UserStats{}, fmt.Errorf("error al contar usuarios premium: %v", err)
	}
	
	userStats := UserStats{
		New:       int(newUsers),
		Returning: int(returningUsers),
		Premium:   int(premiumUsers),
	}
	
	return userStats, nil
}

// getPaymentMethods obtiene la distribución de métodos de pago
func getPaymentMethods(startDate, endDate time.Time) (PaymentMethods, error) {
	// Consultar cantidades por método de pago
	var paypalCount, cardCount, transferCount int64
	
	// Contar pagos con PayPal
	if err := db.Model(&Pago{}).
		Where("metodo = 'paypal' AND estado = 'aprobado' AND created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&paypalCount).Error; err != nil {
		return PaymentMethods{}, fmt.Errorf("error al contar pagos con PayPal: %v", err)
	}
	
	// Contar pagos con tarjeta
	if err := db.Model(&Pago{}).
		Where("metodo = 'tarjeta' AND estado = 'aprobado' AND created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&cardCount).Error; err != nil {
		return PaymentMethods{}, fmt.Errorf("error al contar pagos con tarjeta: %v", err)
	}
	
	// Contar pagos con transferencia
	if err := db.Model(&Pago{}).
		Where("metodo = 'transferencia' AND estado = 'aprobado' AND created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&transferCount).Error; err != nil {
		return PaymentMethods{}, fmt.Errorf("error al contar pagos con transferencia: %v", err)
	}
	
	// Calcular total
	total := paypalCount + cardCount + transferCount
	
	// Calcular porcentajes
	paymentMethods := PaymentMethods{}
	
	if total > 0 {
		paymentMethods.Paypal = int((float64(paypalCount) / float64(total)) * 100)
		paymentMethods.Card = int((float64(cardCount) / float64(total)) * 100)
		paymentMethods.Transfer = int((float64(transferCount) / float64(total)) * 100)
		
		// Ajustar para asegurar que sumen 100%
		remaining := 100 - paymentMethods.Paypal - paymentMethods.Card - paymentMethods.Transfer
		
		// Asignar el residuo al método con mayor cantidad
		if remaining != 0 {
			if paypalCount >= cardCount && paypalCount >= transferCount {
				paymentMethods.Paypal += remaining
			} else if cardCount >= paypalCount && cardCount >= transferCount {
				paymentMethods.Card += remaining
			} else {
				paymentMethods.Transfer += remaining
			}
		}
	}
	
	return paymentMethods, nil
}

// fetchDashboardData obtiene datos generales para el dashboard
func fetchDashboardData() (DashboardData, error) {
	dashboardData := DashboardData{}
	
	// Total de usuarios
	var totalUsers int64
	if err := db.Model(&Usuario{}).Count(&totalUsers).Error; err != nil {
		return dashboardData, fmt.Errorf("error al contar usuarios: %v", err)
	}
	dashboardData.TotalUsers = int(totalUsers)
	
	// Total de cursos
	var totalCourses int64
	if err := db.Model(&Curso{}).Count(&totalCourses).Error; err != nil {
		return dashboardData, fmt.Errorf("error al contar cursos: %v", err)
	}
	dashboardData.TotalCourses = int(totalCourses)
	
	// Total de ingresos
	if err := db.Model(&Pago{}).
		Select("COALESCE(SUM(monto), 0) as total").
		Where("estado = 'aprobado'").
		Scan(&dashboardData.TotalRevenue).Error; err != nil {
		return dashboardData, fmt.Errorf("error al calcular ingresos totales: %v", err)
	}
	
	// Pagos pendientes
	var pendingPayments int64
	if err := db.Model(&Pago{}).
		Where("estado = 'pendiente'").
		Count(&pendingPayments).Error; err != nil {
		return dashboardData, fmt.Errorf("error al contar pagos pendientes: %v", err)
	}
	dashboardData.PendingPayments = int(pendingPayments)
	
	// Usuarios recientes
	type RecentUser struct {
		ID        uint      `json:"id"`
		Nombre    string    `json:"nombre"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
	}
	
	var recentUsers []RecentUser
	if err := db.Model(&Usuario{}).
		Select("id, nombre, email, created_at").
		Order("created_at DESC").
		Limit(5).
		Find(&recentUsers).Error; err != nil {
		return dashboardData, fmt.Errorf("error al obtener usuarios recientes: %v", err)
	}
	
	// Convertir a estructura esperada
	for _, user := range recentUsers {
		dashboardData.RecentUsers = append(dashboardData.RecentUsers, struct {
			ID        uint      `json:"id"`
			Nombre    string    `json:"nombre"`
			Email     string    `json:"email"`
			CreatedAt time.Time `json:"created_at"`
		}{
			ID:        user.ID,
			Nombre:    user.Nombre,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		})
	}
	
	// Pagos recientes con nombres de usuario y curso
	type RecentPayment struct {
		ID            uint      `json:"id"`
		UsuarioNombre string    `json:"usuario_nombre"`
		CursoTitulo   string    `json:"curso_titulo"`
		Monto         float64   `json:"monto"`
		Estado        string    `json:"estado"`
		CreatedAt     time.Time `json:"created_at"`
	}
	
	var recentPayments []RecentPayment
	if err := db.Table("pagos").
		Select("pagos.id, usuarios.nombre as usuario_nombre, cursos.titulo as curso_titulo, pagos.monto, pagos.estado, pagos.created_at").
		Joins("JOIN usuarios ON pagos.usuario_id = usuarios.id").
		Joins("JOIN cursos ON pagos.curso_id = cursos.id").
		Order("pagos.created_at DESC").
		Limit(5).
		Scan(&recentPayments).Error; err != nil {
		return dashboardData, fmt.Errorf("error al obtener pagos recientes: %v", err)
	}
	
	// Convertir a estructura esperada
	for _, payment := range recentPayments {
		dashboardData.RecentPayments = append(dashboardData.RecentPayments, struct {
			ID            uint      `json:"id"`
			UsuarioNombre string    `json:"usuario_nombre"`
			CursoTitulo   string    `json:"curso_titulo"`
			Monto         float64   `json:"monto"`
			Estado        string    `json:"estado"`
			CreatedAt     time.Time `json:"created_at"`
		}{
			ID:            payment.ID,
			UsuarioNombre: payment.UsuarioNombre,
			CursoTitulo:   payment.CursoTitulo,
			Monto:         payment.Monto,
			Estado:        payment.Estado,
			CreatedAt:     payment.CreatedAt,
		})
	}
	
	return dashboardData, nil
}

// fetchActivityLog obtiene el registro de actividad con paginación y filtros
func fetchActivityLog(page, limit, userID, action, startDate, endDate string) ([]ActivityLogData, int, error) {
	// Convertir parámetros de paginación a enteros
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}
	
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		limitInt = 50
	}
	
	// Calcular offset para paginación
	offset := (pageInt - 1) * limitInt
	
	// Preparar consulta base
	query := db.Model(&ActivityLog{})
	countQuery := db.Model(&ActivityLog{})
	
	// Aplicar filtros
	if userID != "" {
		userIDInt, err := strconv.Atoi(userID)
		if err == nil {
			query = query.Where("user_id = ?", userIDInt)
			countQuery = countQuery.Where("user_id = ?", userIDInt)
		}
	}
	
	if action != "" {
		query = query.Where("action = ?", action)
		countQuery = countQuery.Where("action = ?", action)
	}
	
	// Filtrar por fechas si se proporcionan
	if startDate != "" {
		startTime, err := time.Parse("2006-01-02", startDate)
		if err == nil {
			query = query.Where("created_at >= ?", startTime)
			countQuery = countQuery.Where("created_at >= ?", startTime)
		}
	}
	
	if endDate != "" {
		endTime, err := time.Parse("2006-01-02", endDate)
		if err == nil {
			// Añadir un día para incluir todo el día final
			endTime = endTime.Add(24 * time.Hour)
			query = query.Where("created_at < ?", endTime)
			countQuery = countQuery.Where("created_at < ?", endTime)
		}
	}
	
	// Contar total de entradas para metadatos de paginación
	var totalEntries int64
	if err := countQuery.Count(&totalEntries).Error; err != nil {
		return nil, 0, fmt.Errorf("error al contar registros de actividad: %v", err)
	}
	
	// Ejecutar consulta principal con paginación
	var logs []ActivityLog
	if err := query.Order("created_at DESC").
		Limit(limitInt).
		Offset(offset).
		Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("error al obtener registros de actividad: %v", err)
	}
	
	// Preparar la respuesta con nombres de usuario
	var result []ActivityLogData
	
	// Mapa para cachear nombres de usuario y evitar consultas repetidas
	userNames := make(map[uint]string)
	
	for _, log := range logs {
		// Obtener nombre de usuario si no está en caché
		userName, exists := userNames[log.UserID]
		if !exists {
			var user Usuario
			if err := db.Select("nombre").Where("id = ?", log.UserID).First(&user).Error; err == nil {
				userName = user.Nombre
				userNames[log.UserID] = userName
			} else {
				userName = "Usuario Desconocido"
			}
		}
		
		// Añadir a resultados
		result = append(result, ActivityLogData{
			ID:        log.ID,
			UserID:    log.UserID,
			UserName:  userName,
			Action:    log.Action,
			Details:   log.Details,
			IP:        log.IP,
			UserAgent: log.UserAgent,
			CreatedAt: log.CreatedAt,
		})
	}
	
	return result, int(totalEntries), nil
}