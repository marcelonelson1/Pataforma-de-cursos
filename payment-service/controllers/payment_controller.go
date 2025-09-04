// controllers/payment_controller.go
package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"payment-service/config"
	"payment-service/models"
	"payment-service/services"
	"payment-service/utils"
)

type PaymentController struct {
	db             *gorm.DB
	config         *config.Config
	paymentService *services.PaymentService
	cleanupService *services.CleanupService
}

func NewPaymentController(db *gorm.DB, cfg *config.Config) *PaymentController {
	return &PaymentController{
		db:             db,
		config:         cfg,
		paymentService: services.NewPaymentService(db, cfg),
		cleanupService: services.NewCleanupService(db, cfg),
	}
}

// VerifyCourseAccess verifica si un usuario tiene acceso a un curso
func (pc *PaymentController) VerifyCourseAccess(c *gin.Context) {
	userIDStr := c.Query("user_id")
	courseIDStr := c.Query("course_id")

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.SendErrorResponse(c, "user_id inválido", http.StatusBadRequest)
		return
	}

	courseID, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		utils.SendErrorResponse(c, "course_id inválido", http.StatusBadRequest)
		return
	}

	log.Printf("[COURSE_ACCESS] Verificando acceso: Usuario %d, Curso %d", userID, courseID)

	// Buscar pago aprobado para este usuario y curso
	var payment models.Payment
	err = pc.db.Where("usuario_id = ? AND curso_id = ? AND estado = ?", 
		uint(userID), uint(courseID), models.PaymentStatusApproved).First(&payment).Error

	hasAccess := err == nil
	log.Printf("[COURSE_ACCESS] Usuario %d tiene acceso al curso %d: %t", userID, courseID, hasAccess)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"has_access":  hasAccess,
			"user_id":     userID,
			"course_id":   courseID,
		},
	})
}

// GetCourseAccessInfo obtiene información detallada de acceso al curso
func (pc *PaymentController) GetCourseAccessInfo(c *gin.Context) {
	userIDStr := c.Query("user_id")
	courseIDStr := c.Query("course_id")

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.SendErrorResponse(c, "user_id inválido", http.StatusBadRequest)
		return
	}

	courseID, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		utils.SendErrorResponse(c, "course_id inválido", http.StatusBadRequest)
		return
	}

	log.Printf("[COURSE_ACCESS_INFO] Obteniendo info de acceso: Usuario %d, Curso %d", userID, courseID)

	// Buscar el pago para este usuario y curso
	var payment models.Payment
	err = pc.db.Where("usuario_id = ? AND curso_id = ?", 
		uint(userID), uint(courseID)).Order("created_at DESC").First(&payment).Error

	if err != nil {
		// No hay pago registrado
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"has_access":      false,
				"payment_status":  "no_payment",
				"access_type":     "none",
			},
		})
		return
	}

	hasAccess := payment.Estado == models.PaymentStatusApproved
	var paymentDate *time.Time
	if hasAccess {
		paymentDate = &payment.UpdatedAt
	}

	log.Printf("[COURSE_ACCESS_INFO] Usuario %d - Estado: %s, Acceso: %t", userID, payment.Estado, hasAccess)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"has_access":      hasAccess,
			"payment_status":  payment.Estado,
			"payment_date":    paymentDate,
			"access_type":     "paid",
			"payment_id":      payment.ID,
			"amount":          payment.Monto,
			"currency":        payment.Moneda,
			"payment_method":  payment.Metodo,
		},
	})
}

// CreatePayment maneja la creación de pagos con validaciones mejoradas
func (pc *PaymentController) CreatePayment(c *gin.Context) {
	startTime := time.Now()
	log.Printf("🟡 [INICIO] CreatePayment - Timestamp: %v", startTime)
	
	// Obtener usuario autenticado
	userID, exists := c.Get("user_id")
	if !exists {
		log.Printf("❌ [ERROR] Usuario no autenticado")
		utils.SendErrorResponse(c, "Usuario no autenticado", http.StatusUnauthorized)
		return
	}
	log.Printf("✅ [AUTH] Usuario autenticado: %v", userID)

	userIDUint, ok := userID.(uint)
	if !ok {
		log.Printf("❌ [ERROR] Error al convertir user_id: %v (tipo: %T)", userID, userID)
		utils.SendErrorResponse(c, "Error al obtener información del usuario", http.StatusInternalServerError)
		return
	}
	log.Printf("✅ [AUTH] UserID convertido: %d", userIDUint)

	// Parsear request
	var req models.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ [ERROR] Error al parsear solicitud: %v", err)
		utils.SendErrorResponse(c, "Datos de solicitud inválidos: "+err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("✅ [REQUEST] Solicitud parseada: UserID=%d, CursoID=%d, Método=%s, Monto=%.2f",
		userIDUint, req.CursoID, req.Metodo, req.Monto)

	// Validar método de pago
	log.Printf("🔍 [VALIDATION] Validando método de pago: %s", req.Metodo)
	if !models.IsValidPaymentMethod(req.Metodo) {
		log.Printf("❌ [ERROR] Método de pago no válido: %s", req.Metodo)
		utils.SendErrorResponse(c, "Método de pago no válido", http.StatusBadRequest)
		return
	}
	log.Printf("✅ [VALIDATION] Método de pago válido: %s", req.Metodo)

	// Validar detalles de tarjeta si es necesario
	if req.Metodo == models.PaymentMethodCard && req.DetallesTarjeta == nil {
		log.Printf("❌ [ERROR] Se requieren detalles de tarjeta para método: %s", req.Metodo)
		utils.SendErrorResponse(c, "Se requieren detalles de tarjeta para este método de pago", http.StatusBadRequest)
		return
	}

	// Establecer moneda por defecto
	if req.Moneda == "" {
		req.Moneda = "USD"
		log.Printf("🔧 [CONFIG] Moneda establecida por defecto: %s", req.Moneda)
	}

	// 🔥 VALIDACIÓN: Verificar que el usuario existe y tiene acceso
	log.Printf("🔍 [USER_VALIDATION] Validando usuario ID: %d", userIDUint)
	userStartTime := time.Now()
	user, err := services.GetUserByID(userIDUint, pc.config)
	userDuration := time.Since(userStartTime)
	log.Printf("⏱️ [USER_VALIDATION] Tiempo de validación de usuario: %v", userDuration)
	
	if err != nil {
		log.Printf("❌ [ERROR] Error al validar usuario %d: %v", userIDUint, err)
		utils.SendErrorResponse(c, "Usuario no encontrado o no autorizado", http.StatusForbidden)
		return
	}
	log.Printf("✅ [USER_VALIDATION] Usuario validado: %s (%s) - Role: %s", user.Nombre, user.Email, user.Role)

	// 🔥 VALIDACIÓN: Verificar que el curso existe y está disponible
	log.Printf("🔍 [COURSE_VALIDATION] Validando curso ID: %d", req.CursoID)
	courseStartTime := time.Now()
	course, err := services.ValidateCourseAccess(req.CursoID, pc.config)
	courseDuration := time.Since(courseStartTime)
	log.Printf("⏱️ [COURSE_VALIDATION] Tiempo de validación de curso: %v", courseDuration)
	
	if err != nil {
		log.Printf("❌ [ERROR] Error al verificar curso %d: %v", req.CursoID, err)
		utils.SendErrorResponse(c, "Curso no disponible para compra", http.StatusNotFound)
		return
	}
	log.Printf("✅ [COURSE_VALIDATION] Curso validado: %s - Precio: %.2f %s", course.Titulo, course.Precio, req.Moneda)

	// 🔥 VALIDACIÓN: Verificar que el monto coincide con el precio del curso
	if req.Monto != course.Precio {
		log.Printf("❌ [ERROR] Monto inválido para curso %d: recibido %.2f, esperado %.2f",
			req.CursoID, req.Monto, course.Precio)
		utils.SendErrorResponse(c, "El monto no coincide con el precio del curso", http.StatusBadRequest)
		return
	}
	log.Printf("✅ [VALIDATION] Monto válido: %.2f", req.Monto)

	// Verificar si ya existe un pago aprobado
	log.Printf("🔍 [PAYMENT_CHECK] Verificando pagos existentes para usuario %d, curso %d", userIDUint, req.CursoID)
	existingPayment, err := pc.paymentService.GetApprovedPayment(userIDUint, req.CursoID)
	if err == nil && existingPayment != nil {
		log.Printf("⚠️ [WARNING] Usuario %d ya tiene acceso aprobado al curso %d", userIDUint, req.CursoID)
		utils.SendSuccessResponse(c, gin.H{
			"message": "Ya tienes acceso a este curso",
			"estado":  models.PaymentStatusApproved,
			"pago_id": existingPayment.ID,
		})
		return
	}
	log.Printf("✅ [PAYMENT_CHECK] No hay pagos aprobados previos")

	// Crear el nuevo pago
	log.Printf("🔧 [PAYMENT_CREATE] Creando nuevo pago")
	payment := &models.Payment{
		UsuarioID:     userIDUint,
		CursoID:       req.CursoID,
		Monto:         req.Monto,
		Metodo:        req.Metodo,
		Estado:        models.PaymentStatusPending,
		TransaccionID: "",
		Moneda:        req.Moneda,
	}

	// Establecer fecha de expiración (30 minutos por defecto)
	payment.SetExpiration(30 * time.Minute)
	log.Printf("⏰ [PAYMENT_CREATE] Fecha de expiración establecida: %v", payment.ExpiresAt)

	// Guardar en base de datos
	dbStartTime := time.Now()
	if err := pc.db.Create(payment).Error; err != nil {
		log.Printf("❌ [ERROR] Error al guardar pago en base de datos: %v", err)
		utils.SendErrorResponse(c, "Error al crear pago", http.StatusInternalServerError)
		return
	}
	dbDuration := time.Since(dbStartTime)
	log.Printf("✅ [PAYMENT_CREATE] Pago guardado en BD - ID: %d, Tiempo: %v", payment.ID, dbDuration)

	log.Printf("✅ [SUCCESS] Pago creado exitosamente: ID=%d, Usuario=%s, Curso=%s, Monto=%.2f",
		payment.ID, user.Email, course.Titulo, payment.Monto)

	// Procesar según el método de pago
	log.Printf("🔧 [PAYMENT_PROCESS] Procesando pago con método: %s", req.Metodo)
	processStartTime := time.Now()
	
	response, err := pc.paymentService.ProcessPayment(payment, &req, course)
	processDuration := time.Since(processStartTime)
	log.Printf("⏱️ [PAYMENT_PROCESS] Tiempo de procesamiento: %v", processDuration)
	
	if err != nil {
		log.Printf("❌ [ERROR] Error al procesar pago: %v", err)
		utils.SendErrorResponse(c, "Error al procesar pago: "+err.Error(), http.StatusInternalServerError)
		return
	}

	totalDuration := time.Since(startTime)
	log.Printf("🟢 [SUCCESS] CreatePayment completado - Tiempo total: %v", totalDuration)
	
	utils.SendSuccessResponse(c, response)
}

// GetPaymentStatus verifica el estado de un pago por curso (migrado de verificarPagoPorCurso)
func (pc *PaymentController) GetPaymentStatus(c *gin.Context) {
	log.Printf("🟡 [INICIO] GetPaymentStatus")
	
	// Obtener usuario autenticado
	userID, exists := c.Get("user_id")
	if !exists {
		log.Printf("❌ [ERROR] Usuario no autenticado en GetPaymentStatus")
		utils.SendErrorResponse(c, "Usuario no autenticado", http.StatusUnauthorized)
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		log.Printf("❌ [ERROR] Error al obtener información del usuario en GetPaymentStatus")
		utils.SendErrorResponse(c, "Error al obtener información del usuario", http.StatusInternalServerError)
		return
	}

	// Obtener ID del curso
	cursoIDStr := c.Param("id")
	if cursoIDStr == "" {
		log.Printf("❌ [ERROR] ID de curso no proporcionado")
		utils.SendErrorResponse(c, "Se requiere el ID del curso", http.StatusBadRequest)
		return
	}

	cursoID, err := strconv.ParseUint(cursoIDStr, 10, 32)
	if err != nil {
		log.Printf("❌ [ERROR] ID de curso inválido: %s", cursoIDStr)
		utils.SendErrorResponse(c, "ID de curso inválido", http.StatusBadRequest)
		return
	}

	log.Printf("🔍 [PAYMENT_STATUS] Verificando pago para curso ID: %d, usuario ID: %d", cursoID, userIDUint)

	// 🔥 VALIDACIÓN: Verificar que el usuario sigue existiendo
	err = services.ValidateUserAccess(userIDUint, pc.config)
	if err != nil {
		log.Printf("❌ [ERROR] Usuario %d no válido al verificar pago: %v", userIDUint, err)
		utils.SendErrorResponse(c, "Usuario no autorizado", http.StatusForbidden)
		return
	}

	// Buscar el pago más reciente
	var payment models.Payment
	result := pc.db.Where("usuario_id = ? AND curso_id = ?", userIDUint, uint(cursoID)).
		Order("created_at desc").
		First(&payment)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("ℹ️ [INFO] No se encontró pago para usuario %d, curso %d", userIDUint, cursoID)
			c.JSON(http.StatusOK, gin.H{
				"estado":  "no_pagado",
				"message": "No se encontró pago para este curso",
			})
			return
		}
		log.Printf("❌ [ERROR] Error de base de datos al verificar pago: %v", result.Error)
		utils.SendErrorResponse(c, "Error de base de datos", http.StatusInternalServerError)
		return
	}

	// Para métodos externos, verificar estado actual si está pendiente
	if payment.IsPending() {
		log.Printf("🔄 [UPDATE] Actualizando estado de pago pendiente: %d", payment.ID)
		pc.paymentService.UpdatePendingPayment(&payment)
		
		// 🔥 RECARGAR EL PAGO DESPUÉS DE LA ACTUALIZACIÓN
		if err := pc.db.First(&payment, payment.ID).Error; err != nil {
			log.Printf("❌ [UPDATE] Error al recargar pago actualizado: %v", err)
		} else {
			log.Printf("✅ [UPDATE] Estado actualizado automáticamente: %s", payment.Estado)
		}
	}

	log.Printf("✅ [SUCCESS] Estado de pago obtenido: %s", payment.Estado)
	c.JSON(http.StatusOK, gin.H{
		"estado": payment.Estado,
		"pago":   payment,
	})
}

// GetUserPayments obtiene todos los pagos de un usuario
func (pc *PaymentController) GetUserPayments(c *gin.Context) {
	log.Printf("🟡 [INICIO] GetUserPayments")
	
	userID, exists := c.Get("user_id")
	if !exists {
		log.Printf("❌ [ERROR] Usuario no autenticado en GetUserPayments")
		utils.SendErrorResponse(c, "Usuario no autenticado", http.StatusUnauthorized)
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		log.Printf("❌ [ERROR] Error al obtener información del usuario en GetUserPayments")
		utils.SendErrorResponse(c, "Error al obtener información del usuario", http.StatusInternalServerError)
		return
	}

	// 🔥 VALIDACIÓN: Verificar que el usuario existe
	err := services.ValidateUserAccess(userIDUint, pc.config)
	if err != nil {
		log.Printf("❌ [ERROR] Usuario %d no válido al obtener pagos: %v", userIDUint, err)
		utils.SendErrorResponse(c, "Usuario no autorizado", http.StatusForbidden)
		return
	}

	var payments []models.Payment
	if err := pc.db.Where("usuario_id = ?", userIDUint).
		Order("created_at desc").
		Find(&payments).Error; err != nil {
		log.Printf("❌ [ERROR] Error al obtener pagos del usuario: %v", err)
		utils.SendErrorResponse(c, "Error al obtener pagos", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ [SUCCESS] Pagos obtenidos para usuario %d: %d pagos", userIDUint, len(payments))
	utils.SendSuccessResponse(c, gin.H{
		"payments": payments,
		"total":    len(payments),
	})
}

// UpdatePaymentStatus actualiza el estado de un pago
func (pc *PaymentController) UpdatePaymentStatus(c *gin.Context) {
	log.Printf("🟡 [INICIO] UpdatePaymentStatus")
	
	paymentIDStr := c.Param("id")
	paymentID, err := strconv.ParseUint(paymentIDStr, 10, 32)
	if err != nil {
		log.Printf("❌ [ERROR] ID de pago inválido: %s", paymentIDStr)
		utils.SendErrorResponse(c, "ID de pago inválido", http.StatusBadRequest)
		return
	}

	var req struct {
		Estado string `json:"estado" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ [ERROR] Datos inválidos en UpdatePaymentStatus: %v", err)
		utils.SendErrorResponse(c, "Datos inválidos", http.StatusBadRequest)
		return
	}

	var payment models.Payment
	if err := pc.db.First(&payment, uint(paymentID)).Error; err != nil {
		log.Printf("❌ [ERROR] Pago no encontrado: %d", paymentID)
		utils.SendErrorResponse(c, "Pago no encontrado", http.StatusNotFound)
		return
	}

	// 🔥 VALIDACIÓN: Verificar que el usuario del pago sigue existiendo
	err = services.ValidateUserAccess(payment.UsuarioID, pc.config)
	if err != nil {
		log.Printf("❌ [ERROR] Usuario %d del pago %d no válido: %v", payment.UsuarioID, paymentID, err)
		utils.SendErrorResponse(c, "Usuario del pago no autorizado", http.StatusForbidden)
		return
	}

	oldStatus := payment.Estado
	payment.UpdateStatus(req.Estado)
	if err := pc.db.Save(&payment).Error; err != nil {
		log.Printf("❌ [ERROR] Error al actualizar estado de pago: %v", err)
		utils.SendErrorResponse(c, "Error al actualizar pago", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ [SUCCESS] Estado de pago %d actualizado de '%s' a '%s'", paymentID, oldStatus, req.Estado)
	utils.SendSuccessResponse(c, gin.H{
		"message": "Estado de pago actualizado correctamente",
		"payment": payment,
	})
}

// HealthCheck endpoint para verificar la salud del servicio
func (pc *PaymentController) HealthCheck(c *gin.Context) {
	log.Printf("🟡 [INICIO] HealthCheck")
	
	// Verificar conexión a base de datos
	sqlDB, err := pc.db.DB()
	if err != nil {
		log.Printf("❌ [ERROR] Error de base de datos en HealthCheck: %v", err)
		utils.SendErrorResponse(c, "Error de base de datos", http.StatusInternalServerError)
		return
	}

	if err := sqlDB.Ping(); err != nil {
		log.Printf("❌ [ERROR] Base de datos no disponible: %v", err)
		utils.SendErrorResponse(c, "Base de datos no disponible", http.StatusServiceUnavailable)
		return
	}

	// 🔥 NUEVO: Verificar conexión con servicios externos
	log.Printf("🔍 [HEALTH] Verificando servicios externos")
	externalServicesStatus := services.HealthCheckExternalServices(pc.config)

	log.Printf("✅ [SUCCESS] HealthCheck completado")
	utils.SendSuccessResponse(c, gin.H{
		"status":            "ok",
		"service":           "payment-service",
		"version":           "1.0.0",
		"database":          "connected",
		"external_services": externalServicesStatus,
	})
}

// PayPalHealthCheck endpoint específico para verificar PayPal
func (pc *PaymentController) PayPalHealthCheck(c *gin.Context) {
	log.Printf("🟡 [INICIO] PayPalHealthCheck")
	
	// Verificar configuración de PayPal
	if pc.config.PayPal.ClientID == "" || pc.config.PayPal.Secret == "" {
		log.Printf("❌ [ERROR] Credenciales PayPal no configuradas")
		utils.SendErrorResponse(c, "Credenciales PayPal no configuradas", http.StatusServiceUnavailable)
		return
	}

	log.Printf("✅ [SUCCESS] PayPal configurado correctamente")
	utils.SendSuccessResponse(c, gin.H{
		"status":    "ok",
		"message":   "PayPal configurado correctamente",
		"client_id": pc.config.PayPal.ClientID[:10] + "...", // Solo mostrar primeros 10 caracteres
		"env":       pc.config.PayPal.Env,
	})
}

// 🔥 NUEVO: MercadoPagoHealthCheck endpoint específico para verificar Mercado Pago
func (pc *PaymentController) MercadoPagoHealthCheck(c *gin.Context) {
	log.Printf("🟡 [INICIO] MercadoPagoHealthCheck")
	
	// Verificar configuración de Mercado Pago
	if err := services.ValidateMercadoPagoConfig(pc.config); err != nil {
		log.Printf("❌ [MERCADOPAGO] Configuración inválida: %v", err)
		utils.SendErrorResponse(c, "Configuración de Mercado Pago inválida: "+err.Error(), http.StatusServiceUnavailable)
		return
	}

	log.Printf("✅ [SUCCESS] Mercado Pago configurado correctamente")
	utils.SendSuccessResponse(c, gin.H{
		"status":       "ok",
		"message":      "Mercado Pago configurado correctamente",
		"access_token": services.MaskMercadoPagoToken(pc.config.MercadoPago.AccessToken),
		"environment":  pc.config.MercadoPago.Environment,
		"accept_usd":   pc.config.MercadoPago.AcceptUSD,
	})
}

// 🔥 NUEVO: TestMercadoPago endpoint para testing de Mercado Pago (solo admin)
func (pc *PaymentController) TestMercadoPago(c *gin.Context) {
	log.Printf("🟡 [INICIO] TestMercadoPago")
	
	// Validar configuración
	if err := services.ValidateMercadoPagoConfig(pc.config); err != nil {
		log.Printf("❌ [MERCADOPAGO_TEST] Configuración inválida: %v", err)
		utils.SendErrorResponse(c, "Configuración de Mercado Pago inválida", http.StatusServiceUnavailable)
		return
	}

	// Crear un pago de prueba
	testPayment := &models.Payment{
		UsuarioID:     1, // ID de prueba
		CursoID:       1, // ID de prueba
		Monto:         10.00,
		Metodo:        models.PaymentMethodMercadoPago,
		Estado:        models.PaymentStatusPending,
		Moneda:        "USD",
		TransaccionID: "test_mp_" + strconv.FormatInt(time.Now().Unix(), 10),
	}

	// Información del curso de prueba
	testCourse := &services.CourseInfo{
		ID:          1,
		Titulo:      "Curso de Prueba MercadoPago",
		Descripcion: "Curso para testing de integración con MercadoPago",
		Precio:      10.00,
		Estado:      "published",
	}

	// Intentar crear una preferencia de Mercado Pago
	preference, err := services.CreateMercadoPagoPreference(testPayment, testCourse, pc.config)
	if err != nil {
		log.Printf("❌ [MERCADOPAGO_TEST] Error al crear preferencia de prueba: %v", err)
		utils.SendErrorResponse(c, "Error al probar Mercado Pago: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("✅ [MERCADOPAGO_TEST] Test exitoso - Preferencia creada: %s", preference.ID)
	utils.SendSuccessResponse(c, gin.H{
		"status":        "ok",
		"message":       "Test de Mercado Pago exitoso",
		"preference_id": preference.ID,
		"checkout_url":  preference.InitPoint,
		"test_payment":  testPayment,
	})
}

// 🔥 NUEVO: GetPaymentStats obtiene estadísticas de pagos (admin)
func (pc *PaymentController) GetPaymentStats(c *gin.Context) {
	log.Printf("🟡 [INICIO] GetPaymentStats")

	// Estadísticas generales
	var stats struct {
		TotalPayments      int64                    `json:"total_payments"`
		PaymentsByStatus   []map[string]interface{} `json:"payments_by_status"`
		PaymentsByMethod   []map[string]interface{} `json:"payments_by_method"`
		TotalRevenue       float64                  `json:"total_revenue"`
		RevenueByMonth     []map[string]interface{} `json:"revenue_by_month"`
		AveragePayment     float64                  `json:"average_payment"`
	}

	// Total de pagos
	pc.db.Model(&models.Payment{}).Count(&stats.TotalPayments)

	// Pagos por estado
	var statusStats []struct {
		Estado string  `json:"estado"`
		Count  int64   `json:"count"`
		Total  float64 `json:"total"`
	}
	pc.db.Model(&models.Payment{}).
		Select("estado, COUNT(*) as count, COALESCE(SUM(monto), 0) as total").
		Group("estado").
		Find(&statusStats)

	for _, stat := range statusStats {
		stats.PaymentsByStatus = append(stats.PaymentsByStatus, map[string]interface{}{
			"estado": stat.Estado,
			"count":  stat.Count,
			"total":  stat.Total,
		})
	}

	// Pagos por método
	var methodStats []struct {
		Metodo string  `json:"metodo"`
		Count  int64   `json:"count"`
		Total  float64 `json:"total"`
	}
	pc.db.Model(&models.Payment{}).
		Select("metodo, COUNT(*) as count, COALESCE(SUM(monto), 0) as total").
		Group("metodo").
		Find(&methodStats)

	for _, stat := range methodStats {
		stats.PaymentsByMethod = append(stats.PaymentsByMethod, map[string]interface{}{
			"metodo": stat.Metodo,
			"count":  stat.Count,
			"total":  stat.Total,
		})
	}

	// Revenue total (solo pagos aprobados)
	pc.db.Model(&models.Payment{}).
		Where("estado = ?", models.PaymentStatusApproved).
		Select("COALESCE(SUM(monto), 0)").
		Scan(&stats.TotalRevenue)

	// Promedio de pago
	if stats.TotalPayments > 0 {
		stats.AveragePayment = stats.TotalRevenue / float64(stats.TotalPayments)
	}

	log.Printf("✅ [SUCCESS] Estadísticas obtenidas: %d pagos totales, $%.2f revenue total", 
		stats.TotalPayments, stats.TotalRevenue)

	utils.SendSuccessResponse(c, stats)
}

// 🔥 NUEVO: GetAllPayments obtiene todos los pagos (admin)
func (pc *PaymentController) GetAllPayments(c *gin.Context) {
	log.Printf("🟡 [INICIO] GetAllPayments")

	// Parámetros de paginación
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset := (page - 1) * limit

	// Filtros opcionales
	status := c.Query("status")
	method := c.Query("method")

	query := pc.db.Model(&models.Payment{})

	if status != "" {
		query = query.Where("estado = ?", status)
	}
	if method != "" {
		query = query.Where("metodo = ?", method)
	}

	// Obtener total
	var total int64
	query.Count(&total)

	// Obtener pagos
	var payments []models.Payment
	err := query.Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&payments).Error

	if err != nil {
		log.Printf("❌ [ERROR] Error al obtener pagos: %v", err)
		utils.SendErrorResponse(c, "Error al obtener pagos", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ [SUCCESS] Pagos obtenidos: %d de %d total", len(payments), total)

	utils.SendSuccessResponse(c, gin.H{
		"payments": payments,
		"total":    total,
		"page":     page,
		"limit":    limit,
		"pages":    (total + int64(limit) - 1) / int64(limit),
	})
}

// 🔥 NUEVO: DiagnosticMercadoPago endpoint para debug completo
func (pc *PaymentController) DiagnosticMercadoPago(c *gin.Context) {
	log.Printf("🔍 [DIAGNOSTIC] Iniciando diagnóstico completo de MercadoPago")
	
	// 1. Verificar configuración básica
	diagnosticResult := make(map[string]interface{})
	
	// Configuración
	diagnosticResult["config"] = map[string]interface{}{
		"access_token_masked": services.MaskMercadoPagoToken(pc.config.MercadoPago.AccessToken),
		"access_token_length": len(pc.config.MercadoPago.AccessToken),
		"environment":         pc.config.MercadoPago.Environment,
		"accept_usd":          pc.config.MercadoPago.AcceptUSD,
		"token_format_valid":  len(pc.config.MercadoPago.AccessToken) > 50,
		"is_test_token":       strings.HasPrefix(pc.config.MercadoPago.AccessToken, "TEST-"),
	}
	
	// 2. Test de conectividad
	connectivityError := services.TestMercadoPagoConnection(pc.config)
	diagnosticResult["connectivity"] = map[string]interface{}{
		"success": connectivityError == nil,
		"error":   func() interface{} {
			if connectivityError != nil {
				return connectivityError.Error()
			}
			return nil
		}(),
	}
	
	// 3. Información de la cuenta
	accountInfo, accountError := services.GetMercadoPagoAccountInfo(pc.config)
	if accountError == nil {
		diagnosticResult["account"] = accountInfo
	} else {
		diagnosticResult["account_error"] = accountError.Error()
	}
	
	// 4. Test de creación de preferencia (simulado)
	testPayment := &models.Payment{
		ID:        999999,
		UsuarioID: 1,
		CursoID:   1,
		Monto:     1.00,
		Metodo:    models.PaymentMethodMercadoPago,
		Estado:    models.PaymentStatusPending,
		Moneda:    "USD",
	}
	
	testCourse := &services.CourseInfo{
		ID:          1,
		Titulo:      "Curso de Prueba Diagnóstico",
		Descripcion: "Curso para testing de MercadoPago",
		Precio:      1.00,
		Estado:      "published",
	}
	
	preference, prefError := services.CreateMercadoPagoPreference(testPayment, testCourse, pc.config)
	if prefError == nil {
		diagnosticResult["test_preference"] = map[string]interface{}{
			"success":        true,
			"preference_id":  preference.ID,
			"checkout_url":   preference.InitPoint,
			"collector_id":   preference.CollectorID,
		}
	} else {
		diagnosticResult["test_preference"] = map[string]interface{}{
			"success": false,
			"error":   prefError.Error(),
		}
	}
	
	// 5. Recomendaciones
	recommendations := []string{}
	
	if !strings.HasPrefix(pc.config.MercadoPago.AccessToken, "TEST-") {
		recommendations = append(recommendations, "⚠️ Estás usando un token de PRODUCCIÓN. Para pruebas usa un token TEST-")
	}
	
	if pc.config.MercadoPago.Environment != "sandbox" {
		recommendations = append(recommendations, "⚠️ Estás en modo PRODUCCIÓN. Para pruebas usa environment=sandbox")
	}
	
	if connectivityError != nil {
		recommendations = append(recommendations, "❌ Problema de conectividad. Verifica tu Access Token")
	}
	
	if len(recommendations) == 0 {
		recommendations = append(recommendations, "✅ Todo parece estar configurado correctamente")
		
		if prefError != nil {
			recommendations = append(recommendations, "🔍 Para resolver el error de preferencia:")
			recommendations = append(recommendations, "1. Usa cuentas de PRUEBA para pagar (no tu cuenta real)")
			recommendations = append(recommendations, "2. Crea cuentas de prueba en MercadoPago Developer Panel")
			recommendations = append(recommendations, "3. Usa tarjetas de prueba: 4509 9535 6623 3704")
		}
	}
	
	diagnosticResult["recommendations"] = recommendations
	
	log.Printf("✅ [DIAGNOSTIC] Diagnóstico completado")
	utils.SendSuccessResponse(c, diagnosticResult)
}

// 🔥 NUEVO: RefreshPaymentStatus - Actualizar manualmente el estado de un pago
func (pc *PaymentController) RefreshPaymentStatus(c *gin.Context) {
	log.Printf("🔄 [REFRESH] Iniciando actualización manual de estado de pago")
	
	paymentIDStr := c.Param("id")
	paymentID, err := strconv.ParseUint(paymentIDStr, 10, 32)
	if err != nil {
		log.Printf("❌ [REFRESH] ID de pago inválido: %s", paymentIDStr)
		utils.SendErrorResponse(c, "ID de pago inválido", http.StatusBadRequest)
		return
	}

	// Obtener usuario autenticado
	userID, exists := c.Get("user_id")
	if !exists {
		log.Printf("❌ [REFRESH] Usuario no autenticado")
		utils.SendErrorResponse(c, "Usuario no autenticado", http.StatusUnauthorized)
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		log.Printf("❌ [REFRESH] Error al obtener información del usuario")
		utils.SendErrorResponse(c, "Error al obtener información del usuario", http.StatusInternalServerError)
		return
	}

	// Verificar que el pago pertenece al usuario
	var payment models.Payment
	if err := pc.db.Where("id = ? AND usuario_id = ?", uint(paymentID), userIDUint).First(&payment).Error; err != nil {
		log.Printf("❌ [REFRESH] Pago no encontrado o no pertenece al usuario: %d", paymentID)
		utils.SendErrorResponse(c, "Pago no encontrado", http.StatusNotFound)
		return
	}

	log.Printf("🔍 [REFRESH] Actualizando estado de pago ID: %d, Método: %s, Estado actual: %s", 
		payment.ID, payment.Metodo, payment.Estado)

	// Solo actualizar si está pendiente
	if payment.Estado != models.PaymentStatusPending {
		log.Printf("ℹ️ [REFRESH] Pago ya no está pendiente, estado actual: %s", payment.Estado)
		utils.SendSuccessResponse(c, gin.H{
			"message": "El pago ya fue procesado",
			"estado":  payment.Estado,
			"pago":    payment,
		})
		return
	}

	// Actualizar según el método
	switch payment.Metodo {
	case models.PaymentMethodMercadoPago:
		log.Printf("🔄 [REFRESH] Actualizando estado de MercadoPago...")
		if err := pc.paymentService.UpdateMercadoPagoPaymentStatus(&payment); err != nil {
			log.Printf("❌ [REFRESH] Error al actualizar estado: %v", err)
			utils.SendErrorResponse(c, "Error al actualizar estado del pago: "+err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		log.Printf("⚠️ [REFRESH] Método %s no soporta actualización manual", payment.Metodo)
		utils.SendErrorResponse(c, "Este método de pago no soporta actualización manual", http.StatusBadRequest)
		return
	}

	// Recargar el pago para obtener el estado actualizado
	if err := pc.db.First(&payment, payment.ID).Error; err != nil {
		log.Printf("❌ [REFRESH] Error al recargar pago: %v", err)
		utils.SendErrorResponse(c, "Error al recargar pago", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ [REFRESH] Estado actualizado exitosamente: %s", payment.Estado)
	utils.SendSuccessResponse(c, gin.H{
		"message": "Estado de pago actualizado exitosamente",
		"estado":  payment.Estado,
		"pago":    payment,
	})
}

// 🔥 NUEVO: DeepDiagnosticMercadoPago - Diagnóstico completo para debugging
func (pc *PaymentController) DeepDiagnosticMercadoPago(c *gin.Context) {
	log.Printf("🔍 [DEEP_DIAGNOSTIC] Iniciando diagnóstico completo de MercadoPago")
	
	diagnostic := make(map[string]interface{})
	
	// 1. Verificar configuración básica
	config := map[string]interface{}{
		"access_token_masked": services.MaskMercadoPagoToken(pc.config.MercadoPago.AccessToken),
		"access_token_length": len(pc.config.MercadoPago.AccessToken),
		"access_token_prefix": func() string {
			if len(pc.config.MercadoPago.AccessToken) >= 10 {
				return pc.config.MercadoPago.AccessToken[:10]
			}
			return pc.config.MercadoPago.AccessToken
		}(),
		"environment":        pc.config.MercadoPago.Environment,
		"accept_usd":         pc.config.MercadoPago.AcceptUSD,
		"is_test_token":      strings.HasPrefix(pc.config.MercadoPago.AccessToken, "TEST-"),
		"is_sandbox_env":     pc.config.MercadoPago.Environment == "sandbox",
		"token_format_valid": len(pc.config.MercadoPago.AccessToken) > 50,
	}
	diagnostic["config"] = config
	
	// 2. Test de conectividad básica
	client := &http.Client{Timeout: 10 * time.Second}
	
	// Test 1: Obtener información de la cuenta
	req, _ := http.NewRequest("GET", "https://api.mercadopago.com/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+pc.config.MercadoPago.AccessToken)
	
	resp, err := client.Do(req)
	accountTest := map[string]interface{}{
		"success": err == nil && resp != nil,
		"status_code": func() int {
			if resp != nil { return resp.StatusCode }
			return 0
		}(),
	}
	
	if resp != nil {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		
		var accountInfo map[string]interface{}
		if json.Unmarshal(body, &accountInfo) == nil {
			accountTest["account_id"] = accountInfo["id"]
			accountTest["site_id"] = accountInfo["site_id"]
			accountTest["country_id"] = accountInfo["country_id"]
			accountTest["account_type"] = accountInfo["account_type"]
			accountTest["nickname"] = accountInfo["nickname"]
			accountTest["first_name"] = accountInfo["first_name"]
			accountTest["last_name"] = accountInfo["last_name"]
		} else {
			accountTest["raw_response"] = string(body)
		}
	}
	
	if err != nil {
		accountTest["error"] = err.Error()
	}
	
	diagnostic["account_test"] = accountTest
	
	// 3. Test de creación de preferencia MÍNIMA
	minimalPreference := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"title":       "Test Item",
				"description": "Item de prueba para diagnóstico",
				"quantity":    1,
				"unit_price":  1.0,
				"currency_id": "ARS",
			},
		},
		"external_reference": "test_diagnostic_" + strconv.FormatInt(time.Now().Unix(), 10),
	}
	
	payload, _ := json.Marshal(minimalPreference)
	
	req2, _ := http.NewRequest("POST", "https://api.mercadopago.com/checkout/preferences", bytes.NewReader(payload))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", "Bearer "+pc.config.MercadoPago.AccessToken)
	req2.Header.Set("X-Idempotency-Key", "diagnostic_"+strconv.FormatInt(time.Now().Unix(), 10))
	
	resp2, err2 := client.Do(req2)
	preferenceTest := map[string]interface{}{
		"success": err2 == nil && resp2 != nil && resp2.StatusCode == 201,
		"payload_sent": string(payload),
	}
	
	if resp2 != nil {
		defer resp2.Body.Close()
		body2, _ := io.ReadAll(resp2.Body)
		preferenceTest["status_code"] = resp2.StatusCode
		preferenceTest["response_body"] = string(body2)
		
		if resp2.StatusCode == 201 {
			var prefResp map[string]interface{}
			if json.Unmarshal(body2, &prefResp) == nil {
				preferenceTest["preference_id"] = prefResp["id"]
				preferenceTest["init_point"] = prefResp["init_point"]
				preferenceTest["sandbox_init_point"] = prefResp["sandbox_init_point"]
			}
		}
	}
	
	if err2 != nil {
		preferenceTest["error"] = err2.Error()
	}
	
	diagnostic["preference_test"] = preferenceTest
	
	// 4. Test con diferentes monedas
	currencyTests := make(map[string]interface{})
	
	currencies := []string{"ARS", "USD"}
	for _, currency := range currencies {
		currencyPref := map[string]interface{}{
			"items": []map[string]interface{}{
				{
					"title":       "Test " + currency,
					"description": "Prueba de moneda " + currency,
					"quantity":    1,
					"unit_price":  func() float64 {
						if currency == "USD" { return 1.0 }
						return 100.0
					}(),
					"currency_id": currency,
				},
			},
			"external_reference": "test_" + currency + "_" + strconv.FormatInt(time.Now().Unix(), 10),
		}
		
		currencyPayload, _ := json.Marshal(currencyPref)
		
		currencyReq, _ := http.NewRequest("POST", "https://api.mercadopago.com/checkout/preferences", bytes.NewReader(currencyPayload))
		currencyReq.Header.Set("Content-Type", "application/json")
		currencyReq.Header.Set("Authorization", "Bearer "+pc.config.MercadoPago.AccessToken)
		currencyReq.Header.Set("X-Idempotency-Key", "curr_"+currency+"_"+strconv.FormatInt(time.Now().Unix(), 10))
		
		currencyResp, currencyErr := client.Do(currencyReq)
		
		currencyResult := map[string]interface{}{
			"currency": currency,
			"success": currencyErr == nil && currencyResp != nil && currencyResp.StatusCode == 201,
		}
		
		if currencyResp != nil {
			defer currencyResp.Body.Close()
			currencyBody, _ := io.ReadAll(currencyResp.Body)
			currencyResult["status_code"] = currencyResp.StatusCode
			currencyResult["response"] = string(currencyBody)
		}
		
		if currencyErr != nil {
			currencyResult["error"] = currencyErr.Error()
		}
		
		currencyTests[currency] = currencyResult
	}
	
	diagnostic["currency_tests"] = currencyTests
	
	// 5. Test de payment methods disponibles
	pmReq, _ := http.NewRequest("GET", "https://api.mercadopago.com/v1/payment_methods", nil)
	pmReq.Header.Set("Authorization", "Bearer "+pc.config.MercadoPago.AccessToken)
	
	pmResp, pmErr := client.Do(pmReq)
	paymentMethodsTest := map[string]interface{}{
		"success": pmErr == nil && pmResp != nil && pmResp.StatusCode == 200,
	}
	
	if pmResp != nil {
		defer pmResp.Body.Close()
		pmBody, _ := io.ReadAll(pmResp.Body)
		paymentMethodsTest["status_code"] = pmResp.StatusCode
		
		var paymentMethods []interface{}
		if json.Unmarshal(pmBody, &paymentMethods) == nil {
			paymentMethodsTest["available_methods_count"] = len(paymentMethods)
			paymentMethodsTest["methods"] = paymentMethods
		} else {
			paymentMethodsTest["raw_response"] = string(pmBody)
		}
	}
	
	if pmErr != nil {
		paymentMethodsTest["error"] = pmErr.Error()
	}
	
	diagnostic["payment_methods_test"] = paymentMethodsTest
	
	// 6. Análisis de errores comunes
	analysis := []string{}
	
	if !strings.HasPrefix(pc.config.MercadoPago.AccessToken, "TEST-") {
		analysis = append(analysis, "⚠️ Estás usando un token de PRODUCCIÓN, no de prueba")
	}
	
	if pc.config.MercadoPago.Environment != "sandbox" {
		analysis = append(analysis, "⚠️ Environment no está en 'sandbox'")
	}
	
	if accountTest["success"] == false {
		analysis = append(analysis, "❌ No se puede conectar con la API de MercadoPago")
	}
	
	if preferenceTest["success"] == false {
		analysis = append(analysis, "❌ No se puede crear preferencias de pago")
		analysis = append(analysis, "💡 Revisa que tu cuenta de prueba VENDEDOR tenga permisos")
	}
	
	if len(analysis) == 0 {
		analysis = append(analysis, "✅ Configuración básica parece correcta")
		analysis = append(analysis, "💡 Si aún tienes problemas:")
		analysis = append(analysis, "   1. Asegúrate de usar cuentas de PRUEBA para pagar")
		analysis = append(analysis, "   2. Usa tarjetas de prueba: 4509 9535 6623 3704")
		analysis = append(analysis, "   3. No uses tu cuenta personal de MercadoPago")
	}
	
	diagnostic["analysis"] = analysis
	
	// 7. Recomendaciones específicas
	recommendations := map[string]interface{}{
		"token_type": func() string {
			if strings.HasPrefix(pc.config.MercadoPago.AccessToken, "TEST-") {
				return "✅ Correcto - Token de prueba"
			}
			return "❌ Incorrecto - Token de producción"
		}(),
		"environment": func() string {
			if pc.config.MercadoPago.Environment == "sandbox" {
				return "✅ Correcto - Ambiente sandbox"
			}
			return "❌ Incorrecto - Ambiente production"
		}(),
		"next_steps": []string{
			"1. Ve a https://www.mercadopago.com.ar/developers/panel",
			"2. En 'Cuentas de prueba', verifica que tengas:",
			"   - Cuenta VENDEDOR (para generar el token)",
			"   - Cuenta COMPRADOR (para hacer pagos)",
			"3. Al pagar, usa la cuenta COMPRADOR de prueba",
			"4. Usa tarjetas de prueba, no reales",
		},
	}
	
	diagnostic["recommendations"] = recommendations
	
	log.Printf("✅ [DEEP_DIAGNOSTIC] Diagnóstico completado")
	utils.SendSuccessResponse(c, diagnostic)
}

// ================ ENDPOINTS DE GESTIÓN DE EXPIRACIÓN ================

// CleanupExpiredPayments ejecuta una limpieza manual de pagos expirados
func (pc *PaymentController) CleanupExpiredPayments(c *gin.Context) {
	log.Printf("[CLEANUP] [MANUAL] Iniciando limpieza manual de pagos expirados")

	stats := pc.cleanupService.ManualCleanup()
	
	log.Printf("[SUCCESS] [CLEANUP] Limpieza manual completada")
	utils.SendSuccessResponse(c, gin.H{
		"message": "Limpieza de pagos expirados completada",
		"stats":   stats,
	})
}

// GetCleanupStats obtiene estadísticas del servicio de limpieza
func (pc *PaymentController) GetCleanupStats(c *gin.Context) {
	log.Printf("[CLEANUP] [STATS] Obteniendo estadísticas de limpieza")

	stats := pc.cleanupService.GetCleanupStats()
	
	utils.SendSuccessResponse(c, gin.H{
		"message": "Estadísticas de limpieza obtenidas",
		"stats":   stats,
	})
}

// StartCleanupScheduler inicia el programador automático de limpieza
func (pc *PaymentController) StartCleanupScheduler(c *gin.Context) {
	log.Printf("[CLEANUP] [SCHEDULER] Iniciando programador de limpieza")

	pc.cleanupService.StartCleanupScheduler()
	
	utils.SendSuccessResponse(c, gin.H{
		"message": "Programador de limpieza iniciado",
		"status":  "running",
	})
}

// StopCleanupScheduler detiene el programador automático de limpieza
func (pc *PaymentController) StopCleanupScheduler(c *gin.Context) {
	log.Printf("[CLEANUP] [SCHEDULER] Deteniendo programador de limpieza")

	pc.cleanupService.StopCleanupScheduler()
	
	utils.SendSuccessResponse(c, gin.H{
		"message": "Programador de limpieza detenido",
		"status":  "stopped",
	})
}

// ExpirePayment marca un pago específico como expirado
func (pc *PaymentController) ExpirePayment(c *gin.Context) {
	paymentIDStr := c.Param("id")
	paymentID, err := strconv.ParseUint(paymentIDStr, 10, 32)
	if err != nil {
		utils.SendErrorResponse(c, "ID de pago inválido", http.StatusBadRequest)
		return
	}

	log.Printf("[EXPIRE] [MANUAL] Marcando pago ID %d como expirado manualmente", paymentID)

	err = pc.paymentService.MarkPaymentAsExpired(uint(paymentID))
	if err != nil {
		log.Printf("[ERROR] [EXPIRE] Error al marcar pago como expirado: %v", err)
		utils.SendErrorResponse(c, "Error al marcar pago como expirado: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[SUCCESS] [EXPIRE] Pago ID %d marcado como expirado exitosamente", paymentID)
	utils.SendSuccessResponse(c, gin.H{
		"message":    "Pago marcado como expirado exitosamente",
		"payment_id": paymentID,
	})
}

// ExtendPaymentExpiration extiende la fecha de expiración de un pago
func (pc *PaymentController) ExtendPaymentExpiration(c *gin.Context) {
	paymentIDStr := c.Param("id")
	paymentID, err := strconv.ParseUint(paymentIDStr, 10, 32)
	if err != nil {
		utils.SendErrorResponse(c, "ID de pago inválido", http.StatusBadRequest)
		return
	}

	var req struct {
		Minutes int `json:"minutes" binding:"required,min=1,max=1440"` // máximo 24 horas
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, "Datos inválidos: "+err.Error(), http.StatusBadRequest)
		return
	}

	duration := time.Duration(req.Minutes) * time.Minute
	log.Printf("[EXTEND] [MANUAL] Extendiendo pago ID %d por %v", paymentID, duration)

	err = pc.paymentService.ExtendPaymentExpiration(uint(paymentID), duration)
	if err != nil {
		log.Printf("[ERROR] [EXTEND] Error al extender expiración: %v", err)
		utils.SendErrorResponse(c, "Error al extender expiración: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[SUCCESS] [EXTEND] Expiración del pago ID %d extendida exitosamente", paymentID)
	utils.SendSuccessResponse(c, gin.H{
		"message":    "Expiración del pago extendida exitosamente",
		"payment_id": paymentID,
		"extension":  fmt.Sprintf("%d minutos", req.Minutes),
	})
}