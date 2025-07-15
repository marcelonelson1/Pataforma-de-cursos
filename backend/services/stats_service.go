package services

import (
	"curso-platform/config"
	"curso-platform/models"
	"curso-platform/utils"
	"fmt"
	"time"
	
)

// StatsService maneja la lógica relacionada con estadísticas
type StatsService struct{}

// NewStatsService crea una nueva instancia de StatsService
func NewStatsService() *StatsService {
	return &StatsService{}
}

// GetAdminStats obtiene estadísticas generales para el panel de administración
func (s *StatsService) GetAdminStats(period string) (*models.AdminStatsResponse, error) {
	// Calcular los rangos de fechas según el período
	startDate, endDate, prevStartDate, prevEndDate := utils.CalculateDateRanges(period)
	
	// Consultar estadísticas actuales
	activeStudents, err := s.getActiveStudentsCount(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error al obtener estudiantes activos: %v", err)
	}
	
	// Consultar estadísticas anteriores para comparación
	prevActiveStudents, err := s.getActiveStudentsCount(prevStartDate, prevEndDate)
	if err != nil {
		return nil, fmt.Errorf("error al obtener estudiantes activos previos: %v", err)
	}
	
	// Consultar cursos publicados
	publishedCourses, err := s.getPublishedCoursesCount(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error al obtener cursos publicados: %v", err)
	}
	
	prevPublishedCourses, err := s.getPublishedCoursesCount(prevStartDate, prevEndDate)
	if err != nil {
		return nil, fmt.Errorf("error al obtener cursos publicados previos: %v", err)
	}
	
	// Consultar ingresos mensuales
	monthlyRevenue, err := s.getMonthlyRevenue(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error al obtener ingresos mensuales: %v", err)
	}
	
	prevMonthlyRevenue, err := s.getMonthlyRevenue(prevStartDate, prevEndDate)
	if err != nil {
		return nil, fmt.Errorf("error al obtener ingresos mensuales previos: %v", err)
	}
	
	// Consultar valoración media
	averageRating, err := s.getAverageRating(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error al obtener valoración media: %v", err)
	}
	
	prevAverageRating, err := s.getAverageRating(prevStartDate, prevEndDate)
	if err != nil {
		return nil, fmt.Errorf("error al obtener valoración media previa: %v", err)
	}
	
	// Crear la respuesta
	response := &models.AdminStatsResponse{}
	response.Stats.ActiveStudents = models.StatValue{Current: activeStudents, Previous: prevActiveStudents}
	response.Stats.PublishedCourses = models.StatValue{Current: publishedCourses, Previous: prevPublishedCourses}
	response.Stats.MonthlyRevenue = models.StatValue{Current: monthlyRevenue, Previous: prevMonthlyRevenue}
	response.Stats.AverageRating = models.StatValue{Current: averageRating, Previous: prevAverageRating}
	response.Period = period
	
	return response, nil
}

// GetSalesStats obtiene estadísticas detalladas de ventas
func (s *StatsService) GetSalesStats(period string) (*models.SalesStatsResponse, error) {
	// Calcular los rangos de fechas según el período
	startDate, endDate, _, _ := utils.CalculateDateRanges(period)
	
	// Consultar ventas por curso
	coursesSales, err := s.getCoursesSales(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error al obtener ventas por curso: %v", err)
	}
	
	// Consultar datos mensuales
	monthlyData, err := s.getMonthlyData(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error al obtener datos mensuales: %v", err)
	}
	
	// Consultar estadísticas de usuarios
	userStats, err := s.getUserStats(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error al obtener estadísticas de usuarios: %v", err)
	}
	
	// Consultar métodos de pago
	paymentMethods, err := s.getPaymentMethods(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error al obtener métodos de pago: %v", err)
	}
	
	// Crear la respuesta
	response := &models.SalesStatsResponse{
		CoursesSales:   coursesSales,
		MonthlyData:    monthlyData,
		UserStats:      userStats,
		PaymentMethods: paymentMethods,
		Period:         period,
	}
	
	return response, nil
}

// GetAdminDashboard obtiene un resumen general para el panel de administración
func (s *StatsService) GetAdminDashboard() (*models.DashboardData, error) {
	// Obtener datos del dashboard
	return s.fetchDashboardData()
}

// GetActivityLog obtiene el registro de actividad con paginación y filtros
func (s *StatsService) GetActivityLog(page, limit string, userID, action, startDate, endDate string) ([]models.ActivityLogData, int, error) {
	return s.fetchActivityLog(page, limit, userID, action, startDate, endDate)
}

// Métodos auxiliares para consultar datos estadísticos

// getActiveStudentsCount obtiene el número de estudiantes activos en un rango de fechas
func (s *StatsService) getActiveStudentsCount(startDate, endDate time.Time) (int, error) {
	// Contamos los usuarios con progreso registrado en el rango de fechas
	var count int64
	
	// Contar usuarios distintos con actividad de progreso en el período
	err := config.DB.Model(&models.ProgresoUsuario{}).
		Where("updated_at BETWEEN ? AND ?", startDate, endDate).
		Distinct("usuario_id").
		Count(&count).Error
		
	if err != nil {
		return 0, fmt.Errorf("error al contar estudiantes activos: %v", err)
	}
	
	return int(count), nil
}

// getPublishedCoursesCount obtiene el número de cursos publicados en un rango de fechas
func (s *StatsService) getPublishedCoursesCount(startDate, endDate time.Time) (int, error) {
	var count int64
	
	// Contar cursos publicados en el período
	err := config.DB.Model(&models.Curso{}).
		Where("estado = 'Publicado' AND updated_at BETWEEN ? AND ?", startDate, endDate).
		Count(&count).Error
		
	if err != nil {
		return 0, fmt.Errorf("error al contar cursos publicados: %v", err)
	}
	
	return int(count), nil
}

// getMonthlyRevenue obtiene los ingresos en un rango de fechas
func (s *StatsService) getMonthlyRevenue(startDate, endDate time.Time) (float64, error) {
	var totalRevenue float64
	
	// Consultar la suma de los montos de los pagos aprobados en el período
	err := config.DB.Model(&models.Pago{}).
		Select("COALESCE(SUM(monto), 0) as total").
		Where("estado = 'aprobado' AND created_at BETWEEN ? AND ?", startDate, endDate).
		Scan(&totalRevenue).Error
		
	if err != nil {
		return 0, fmt.Errorf("error al calcular ingresos mensuales: %v", err)
	}
	
	return totalRevenue, nil
}

// getAverageRating obtiene la valoración media de los cursos en un rango de fechas
func (s *StatsService) getAverageRating(startDate, endDate time.Time) (float64, error) {
	// Como no hay tabla de valoraciones definida, usamos un valor predeterminado
	// En una implementación real, esto consultaría de una tabla de valoraciones
	return 4.5, nil
}

// getCoursesSales obtiene las ventas por curso en un rango de fechas
func (s *StatsService) getCoursesSales(startDate, endDate time.Time) ([]models.CourseSale, error) {
	type ResultRow struct {
		CursoID   uint
		Nombre    string
		Ventas    int
		Ingresos  float64
	}
	
	var results []ResultRow
	
	// Consultar ventas agrupadas por curso
	err := config.DB.Table("pagos").
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
	var coursesSales []models.CourseSale
	for _, r := range results {
		percentage := 0
		if totalRevenue > 0 {
			percentage = int((r.Ingresos / totalRevenue) * 100)
		}
		
		coursesSales = append(coursesSales, models.CourseSale{
			Name:       r.Nombre,
			Sales:      r.Ventas,
			Revenue:    r.Ingresos,
			Percentage: percentage,
		})
	}
	
	return coursesSales, nil
}

// getMonthlyData obtiene los datos mensuales de ventas y usuarios
func (s *StatsService) getMonthlyData(startDate, endDate time.Time) ([]models.MonthlyData, error) {
	// Determinar cuántos meses mostrar según el período
	monthsDiff := int(endDate.Sub(startDate).Hours() / 24 / 30) + 1
	if monthsDiff > 12 {
		monthsDiff = 12 // Limitar a máximo 12 meses
	}
	
	var monthlyData []models.MonthlyData
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
		if err := config.DB.Model(&models.Pago{}).
			Where("estado = 'aprobado' AND created_at BETWEEN ? AND ?", monthStart, monthEnd).
			Count(&sales).Error; err != nil {
			return nil, fmt.Errorf("error al obtener ventas mensuales: %v", err)
		}
		
		// Obtener usuarios activos del mes
		var users int64
		if err := config.DB.Model(&models.ProgresoUsuario{}).
			Where("updated_at BETWEEN ? AND ?", monthStart, monthEnd).
			Distinct("usuario_id").
			Count(&users).Error; err != nil {
			return nil, fmt.Errorf("error al obtener usuarios mensuales: %v", err)
		}
		
		monthlyData = append([]models.MonthlyData{
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
func (s *StatsService) getUserStats(startDate, endDate time.Time) (models.UserStats, error) {
	// Contadores para diferentes tipos de usuarios
	var newUsers, returningUsers, premiumUsers int64
	
	// Usuarios nuevos registrados en el período
	if err := config.DB.Model(&models.Usuario{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&newUsers).Error; err != nil {
		return models.UserStats{}, fmt.Errorf("error al contar usuarios nuevos: %v", err)
	}
	
	// Usuarios que retornan (con login en el período pero creados antes)
	if err := config.DB.Model(&models.Usuario{}).
		Where("last_login BETWEEN ? AND ? AND created_at < ?", startDate, endDate, startDate).
		Count(&returningUsers).Error; err != nil {
		return models.UserStats{}, fmt.Errorf("error al contar usuarios que retornan: %v", err)
	}
	
	// Usuarios premium (con pagos aprobados en el período)
	if err := config.DB.Model(&models.Pago{}).
		Select("COUNT(DISTINCT usuario_id)").
		Where("estado = 'aprobado' AND created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&premiumUsers).Error; err != nil {
		return models.UserStats{}, fmt.Errorf("error al contar usuarios premium: %v", err)
	}
	
	userStats := models.UserStats{
		New:       int(newUsers),
		Returning: int(returningUsers),
		Premium:   int(premiumUsers),
	}
	
	return userStats, nil
}

// getPaymentMethods obtiene la distribución de métodos de pago
func (s *StatsService) getPaymentMethods(startDate, endDate time.Time) (models.PaymentMethods, error) {
	// Consultar cantidades por método de pago
	var paypalCount, cardCount, transferCount int64
	
	// Contar pagos con PayPal
	if err := config.DB.Model(&models.Pago{}).
		Where("metodo = 'paypal' AND estado = 'aprobado' AND created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&paypalCount).Error; err != nil {
		return models.PaymentMethods{}, fmt.Errorf("error al contar pagos con PayPal: %v", err)
	}
	
	// Contar pagos con tarjeta
	if err := config.DB.Model(&models.Pago{}).
		Where("metodo = 'tarjeta' AND estado = 'aprobado' AND created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&cardCount).Error; err != nil {
		return models.PaymentMethods{}, fmt.Errorf("error al contar pagos con tarjeta: %v", err)
	}
	
	// Contar pagos con transferencia
	if err := config.DB.Model(&models.Pago{}).
		Where("metodo = 'transferencia' AND estado = 'aprobado' AND created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&transferCount).Error; err != nil {
		return models.PaymentMethods{}, fmt.Errorf("error al contar pagos con transferencia: %v", err)
	}
	
	// Calcular total
	total := paypalCount + cardCount + transferCount
	
	// Calcular porcentajes
	paymentMethods := models.PaymentMethods{}
	
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
func (s *StatsService) fetchDashboardData() (*models.DashboardData, error) {
	dashboardData := &models.DashboardData{}
	
	// Total de usuarios
	var totalUsers int64
	if err := config.DB.Model(&models.Usuario{}).Count(&totalUsers).Error; err != nil {
		return nil, fmt.Errorf("error al contar usuarios: %v", err)
	}
	dashboardData.TotalUsers = int(totalUsers)
	
	// Total de cursos
	var totalCourses int64
	if err := config.DB.Model(&models.Curso{}).Count(&totalCourses).Error; err != nil {
		return nil, fmt.Errorf("error al contar cursos: %v", err)
	}
	dashboardData.TotalCourses = int(totalCourses)
	
	// Total de ingresos
	if err := config.DB.Model(&models.Pago{}).
		Select("COALESCE(SUM(monto), 0) as total").
		Where("estado = 'aprobado'").
		Scan(&dashboardData.TotalRevenue).Error; err != nil {
		return nil, fmt.Errorf("error al calcular ingresos totales: %v", err)
	}
	
	// Pagos pendientes
	var pendingPayments int64
	if err := config.DB.Model(&models.Pago{}).
		Where("estado = 'pendiente'").
		Count(&pendingPayments).Error; err != nil {
		return nil, fmt.Errorf("error al contar pagos pendientes: %v", err)
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
	if err := config.DB.Model(&models.Usuario{}).
		Select("id, nombre, email, created_at").
		Order("created_at DESC").
		Limit(5).
		Find(&recentUsers).Error; err != nil {
		return nil, fmt.Errorf("error al obtener usuarios recientes: %v", err)
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
	if err := config.DB.Table("pagos").
		Select("pagos.id, usuarios.nombre as usuario_nombre, cursos.titulo as curso_titulo, pagos.monto, pagos.estado, pagos.created_at").
		Joins("JOIN usuarios ON pagos.usuario_id = usuarios.id").
		Joins("JOIN cursos ON pagos.curso_id = cursos.id").
		Order("pagos.created_at DESC").
		Limit(5).
		Scan(&recentPayments).Error; err != nil {
		return nil, fmt.Errorf("error al obtener pagos recientes: %v", err)
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
func (s *StatsService) fetchActivityLog(page, limit, userID, action, startDate, endDate string) ([]models.ActivityLogData, int, error) {
	// Convertir parámetros de paginación a enteros
	pageInt, _ := utils.StringToInt(page, 1)
	limitInt, _ := utils.StringToInt(limit, 50)
	
	// Calcular offset para paginación
	offset := (pageInt - 1) * limitInt
	
	// Preparar consulta base
	query := config.DB.Model(&models.ActivityLog{})
	countQuery := config.DB.Model(&models.ActivityLog{})
	
	// Aplicar filtros
	if userID != "" {
		userIDInt, err := utils.StringToInt(userID, 0)
		if err == nil && userIDInt > 0 {
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
	var logs []models.ActivityLog
	if err := query.Order("created_at DESC").
		Limit(limitInt).
		Offset(offset).
		Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("error al obtener registros de actividad: %v", err)
	}
	
	// Preparar la respuesta con nombres de usuario
	var result []models.ActivityLogData
	
	// Mapa para cachear nombres de usuario y evitar consultas repetidas
	userNames := make(map[uint]string)
	
	for _, log := range logs {
		// Obtener nombre de usuario si no está en caché
		userName, exists := userNames[log.UserID]
		if !exists {
			var user models.Usuario
			if err := config.DB.Select("nombre").Where("id = ?", log.UserID).First(&user).Error; err == nil {
				userName = user.Nombre
				userNames[log.UserID] = userName
			} else {
				userName = "Usuario Desconocido"
			}
		}
		
		// Añadir a resultados
		result = append(result, models.ActivityLogData{
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