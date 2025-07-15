package models

import "time"

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