// services/payment_service_nuevo.go - PAYMENT SERVICE COMPLETO CON MERCADOPAGO ACTUALIZADO
package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"

	"payment-service/config"
	"payment-service/models"
)

type PaymentService struct {
	db     *gorm.DB
	config *config.Config
}

func NewPaymentService(db *gorm.DB, cfg *config.Config) *PaymentService {
	return &PaymentService{
		db:     db,
		config: cfg,
	}
}

// ProcessPayment procesa un pago segun el metodo seleccionado
func (ps *PaymentService) ProcessPayment(payment *models.Payment, req *models.PaymentRequest, course *CourseInfo) (*models.PaymentResponse, error) {
	log.Printf("[INFO] [PAYMENT_SERVICE] Iniciando ProcessPayment - Metodo: %s, Pago ID: %d", req.Metodo, payment.ID)

	// Validacion: Validar informacion antes de procesar
	log.Printf("[VALIDATION] Validando datos del pago")
	if err := ps.validatePaymentData(payment, course); err != nil {
		log.Printf("[ERROR] Validacion de pago fallida: %v", err)
		return nil, err
	}
	log.Printf("[SUCCESS] [VALIDATION] Datos del pago validados")

	response := &models.PaymentResponse{
		ID:        payment.ID,
		UsuarioID: payment.UsuarioID,
		CursoID:   payment.CursoID,
		Monto:     payment.Monto,
		Metodo:    payment.Metodo,
		Estado:    payment.Estado,
		Moneda:    payment.Moneda,
		CreatedAt: payment.CreatedAt,
	}

	log.Printf("[CONFIG] [PAYMENT_SERVICE] Procesando con metodo: %s", req.Metodo)

	switch req.Metodo {
	case models.PaymentMethodDev:
		log.Printf("[INFO] [DEV] Procesando pago en modo desarrollo")
		return ps.processDevPayment(payment, response)

	case models.PaymentMethodPayPal:
		log.Printf("[INFO] [PAYPAL] Procesando pago con PayPal")
		return ps.processPayPalPayment(payment, response)

	case models.PaymentMethodCoinbase:
		log.Printf("[INFO] [COINBASE] Procesando pago con Coinbase")
		return ps.processCoinbasePayment(payment, response, course)

	case models.PaymentMethodMercadoPago:
		log.Printf("[INFO] [MERCADOPAGO] Procesando pago con Mercado Pago")
		return ps.processMercadoPagoPayment(payment, response, course)

	case models.PaymentMethodStripe:
		log.Printf("[INFO] [STRIPE] Procesando pago con Stripe")
		return ps.processStripePayment(payment, response)

	case models.PaymentMethodCard, models.PaymentMethodTransfer:
		log.Printf("[INFO] [SIMULATED] Procesando pago simulado: %s", req.Metodo)
		return ps.processSimulatedPayment(payment, response, req.Metodo)

	default:
		log.Printf("[ERROR] Metodo de pago no soportado: %s", req.Metodo)
		return nil, fmt.Errorf("metodo de pago no soportado: %s", req.Metodo)
	}
}

// validatePaymentData valida datos del pago
func (ps *PaymentService) validatePaymentData(payment *models.Payment, course *CourseInfo) error {
	log.Printf("[VALIDATE] Validando monto: pago=%.2f, curso=%.2f", payment.Monto, course.Precio)

	// Validar que el monto coincide con el precio del curso
	if payment.Monto != course.Precio {
		return fmt.Errorf("monto del pago (%.2f) no coincide con precio del curso (%.2f)",
			payment.Monto, course.Precio)
	}

	// Validar que el pago no este duplicado
	log.Printf("[VALIDATE] Verificando pagos duplicados para usuario %d, curso %d", payment.UsuarioID, payment.CursoID)
	var existingPayment models.Payment
	result := ps.db.Where("usuario_id = ? AND curso_id = ? AND estado = ?",
		payment.UsuarioID, payment.CursoID, models.PaymentStatusApproved).First(&existingPayment)

	if result.Error == nil {
		log.Printf("[ERROR] Pago duplicado encontrado: ID %d", existingPayment.ID)
		return fmt.Errorf("el usuario ya tiene un pago aprobado para este curso")
	}

	// Validar monto minimo
	if payment.Monto <= 0 {
		return fmt.Errorf("el monto debe ser mayor a 0")
	}

	// Validar curso activo
	log.Printf("[VALIDATE] Verificando estado del curso: %s", course.Estado)
	if course.Estado != "published" && course.Estado != "active" && course.Estado != "disponible" && course.Estado != "Publicado" {
		return fmt.Errorf("el curso no esta disponible para compra (estado: %s)", course.Estado)
	}

	return nil
}

// processDevPayment procesa pagos en modo desarrollo
func (ps *PaymentService) processDevPayment(payment *models.Payment, response *models.PaymentResponse) (*models.PaymentResponse, error) {
	log.Printf("[CONFIG] [DEV] Iniciando pago en modo desarrollo para ID: %d", payment.ID)

	go ps.simulatePaymentGateway(payment.ID, payment.Metodo)

	response.Message = "Pago en proceso (modo desarrollo)"
	log.Printf("[SUCCESS] [DEV] Pago en desarrollo iniciado correctamente")
	return response, nil
}

// processPayPalPayment procesa pagos con PayPal
func (ps *PaymentService) processPayPalPayment(payment *models.Payment, response *models.PaymentResponse) (*models.PaymentResponse, error) {
	log.Printf("[CONFIG] [PAYPAL] Iniciando creacion de orden PayPal para pago ID: %d", payment.ID)

	// Verificar configuracion de PayPal
	if ps.config.PayPal.ClientID == "" || ps.config.PayPal.Secret == "" {
		log.Printf("[ERROR] [PAYPAL] Credenciales PayPal no configuradas")
		return nil, fmt.Errorf("credenciales PayPal no configuradas")
	}

	log.Printf("[SUCCESS] [PAYPAL] Credenciales PayPal disponibles - ClientID: %s...", ps.config.PayPal.ClientID[:10])

	// Crear orden PayPal
	order, err := CreatePayPalOrder(payment, ps.config)
	if err != nil {
		log.Printf("[ERROR] [PAYPAL] Error al crear orden PayPal: %v", err)
		return nil, fmt.Errorf("error al procesar pago con PayPal: %v", err)
	}
	log.Printf("[SUCCESS] [PAYPAL] Orden PayPal creada: ID=%s", order.ID)

	// Obtener URL de aprobacion
	approvalURL := GetPayPalApprovalURL(order)
	if approvalURL == "" {
		log.Printf("[ERROR] [PAYPAL] No se pudo obtener URL de aprobacion")
		return nil, fmt.Errorf("error al obtener URL de PayPal - no se encontro enlace de aprobacion")
	}
	log.Printf("[SUCCESS] [PAYPAL] URL de aprobacion obtenida: %s", approvalURL)

	// Actualizar ID de transaccion
	payment.TransaccionID = order.ID
	if err := ps.db.Save(payment).Error; err != nil {
		log.Printf("[ERROR] [PAYPAL] Error al actualizar transaccion ID: %v", err)
	} else {
		log.Printf("[SUCCESS] [PAYPAL] Transaccion ID actualizada: %s", order.ID)
	}

	response.CheckoutURL = approvalURL
	response.Message = "Redirigir a PayPal para completar el pago"
	log.Printf("[SUCCESS] [PAYPAL] Proceso PayPal completado exitosamente")
	return response, nil
}

// processCoinbasePayment procesa pagos con Coinbase
func (ps *PaymentService) processCoinbasePayment(payment *models.Payment, response *models.PaymentResponse, course *CourseInfo) (*models.PaymentResponse, error) {
	log.Printf("[NEW] [COINBASE] Iniciando pago con Coinbase para ID: %d", payment.ID)

	// NUEVA VALIDACION: Verificar configuracion de Coinbase
	if err := ValidateCoinbaseConfig(ps.config); err != nil {
		log.Printf("[ERROR] [COINBASE] Configuracion invalida: %v", err)
		return nil, fmt.Errorf("configuracion de Coinbase invalida: %v", err)
	}

	// NUEVA FUNCION: Test de conectividad opcional
	if err := TestCoinbaseConnection(ps.config); err != nil {
		log.Printf("[WARNING] [COINBASE] Problema de conectividad: %v", err)
		// No retornar error aqui, intentar crear el cargo de todas formas
	}

	charge, err := CreateCoinbaseCharge(payment, course, ps.config)
	if err != nil {
		log.Printf("[ERROR] [COINBASE] Error al crear cargo: %v", err)

		// MEJORADO: Mensajes de error mas especificos
		if strings.Contains(err.Error(), "412") {
			return nil, fmt.Errorf("error de configuracion Coinbase: tu cuenta puede necesitar verificacion adicional o el servicio no esta disponible en tu region")
		}
		if strings.Contains(err.Error(), "401") {
			return nil, fmt.Errorf("credenciales Coinbase invalidas: verifica tu API Key")
		}

		return nil, fmt.Errorf("error al procesar pago con Coinbase: %v", err)
	}
	log.Printf("[SUCCESS] [COINBASE] Cargo creado: %s", charge.Code)

	payment.TransaccionID = charge.Code
	if err := ps.db.Save(payment).Error; err != nil {
		log.Printf("[ERROR] [COINBASE] Error al actualizar transaccion ID: %v", err)
	} else {
		log.Printf("[SUCCESS] [COINBASE] Transaccion ID actualizada: %s", charge.Code)
	}

	response.CheckoutURL = charge.HostedURL
	response.Message = "Redirigir a Coinbase para completar el pago"
	log.Printf("[SUCCESS] [COINBASE] Proceso completado exitosamente")
	return response, nil
}

// processMercadoPagoPayment procesa pagos con Mercado Pago (NUEVA IMPLEMENTACION LIMPIA)
func (ps *PaymentService) processMercadoPagoPayment(payment *models.Payment, response *models.PaymentResponse, course *CourseInfo) (*models.PaymentResponse, error) {
	log.Printf("[MP] Iniciando proceso MercadoPago para pago ID: %d", payment.ID)

	// Crear servicio MercadoPago
	mpService, err := NewMercadoPagoService(ps.config)
	if err != nil {
		log.Printf("[ERROR] Error creando servicio MercadoPago: %v", err)
		return nil, fmt.Errorf("error de configuracion MercadoPago: %v", err)
	}

	// Validar configuracion
	if err := mpService.ValidateConfig(); err != nil {
		log.Printf("[ERROR] Configuracion invalida: %v", err)
		return nil, fmt.Errorf("configuracion de MercadoPago invalida: %v", err)
	}

	// Crear preferencia
	preference, err := mpService.CreatePreference(payment, course)
	if err != nil {
		log.Printf("[ERROR] Error creando preferencia: %v", err)
		return nil, fmt.Errorf("error al crear preferencia MercadoPago: %v", err)
	}

	// Actualizar pago con ID de preferencia
	payment.TransaccionID = preference.ID
	if err := ps.db.Save(payment).Error; err != nil {
		log.Printf("[ERROR] Error actualizando pago: %v", err)
	}

	// Configurar respuesta
	checkoutURL := mpService.GetCheckoutURL(preference)
	response.CheckoutURL = checkoutURL
	response.Message = "Redirigir a MercadoPago para completar el pago"

	log.Printf("[SUCCESS] Proceso MercadoPago completado - Checkout URL: %s", checkoutURL)
	return response, nil
}

// processStripePayment simula pago con Stripe
func (ps *PaymentService) processStripePayment(payment *models.Payment, response *models.PaymentResponse) (*models.PaymentResponse, error) {
	log.Printf("[CONFIG] [STRIPE] Iniciando pago con Stripe para ID: %d", payment.ID)

	mockCheckoutURL := fmt.Sprintf("https://checkout.stripe.com/pay/cs_%d_%d",
		time.Now().Unix(), rand.Intn(10000))

	payment.TransaccionID = fmt.Sprintf("ch_%d_%d", time.Now().Unix(), rand.Intn(10000))
	if err := ps.db.Save(payment).Error; err != nil {
		log.Printf("[ERROR] [STRIPE] Error al actualizar transaccion ID: %v", err)
	} else {
		log.Printf("[SUCCESS] [STRIPE] Transaccion ID generada: %s", payment.TransaccionID)
	}

	response.CheckoutURL = mockCheckoutURL
	response.Message = "Redirigir a Stripe para completar el pago"
	log.Printf("[SUCCESS] [STRIPE] Proceso completado exitosamente")
	return response, nil
}

// processSimulatedPayment procesa pagos simulados (tarjeta, transferencia)
func (ps *PaymentService) processSimulatedPayment(payment *models.Payment, response *models.PaymentResponse, method string) (*models.PaymentResponse, error) {
	log.Printf("[CONFIG] [SIMULATED] Iniciando pago simulado (%s) para ID: %d", method, payment.ID)

	go ps.simulatePaymentGateway(payment.ID, method)

	response.Message = fmt.Sprintf("Pago con %s en proceso", method)
	log.Printf("[SUCCESS] [SIMULATED] Pago simulado iniciado correctamente")
	return response, nil
}

// GetApprovedPayment busca un pago aprobado para un usuario y curso
func (ps *PaymentService) GetApprovedPayment(userID, courseID uint) (*models.Payment, error) {
	log.Printf("[SEARCH] [GET_APPROVED] Buscando pago aprobado - Usuario: %d, Curso: %d", userID, courseID)

	var payment models.Payment
	err := ps.db.Where("usuario_id = ? AND curso_id = ? AND estado = ?",
		userID, courseID, models.PaymentStatusApproved).First(&payment).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("[INFO] [GET_APPROVED] No se encontro pago aprobado")
		} else {
			log.Printf("[ERROR] [GET_APPROVED] Error en consulta: %v", err)
		}
		return nil, err
	}

	log.Printf("[SUCCESS] [GET_APPROVED] Pago aprobado encontrado: ID %d", payment.ID)
	return &payment, nil
}

// GetPendingPayments valida si un usuario tiene pagos pendientes
func (ps *PaymentService) GetPendingPayments(userID uint) ([]models.Payment, error) {
	log.Printf("[SEARCH] [GET_PENDING] Obteniendo pagos pendientes para usuario: %d", userID)

	var payments []models.Payment
	err := ps.db.Where("usuario_id = ? AND estado = ?",
		userID, models.PaymentStatusPending).Find(&payments).Error

	if err != nil {
		log.Printf("[ERROR] [GET_PENDING] Error al obtener pagos pendientes: %v", err)
	} else {
		log.Printf("[SUCCESS] [GET_PENDING] Encontrados %d pagos pendientes", len(payments))
	}

	return payments, err
}

// GetUserPaymentStats obtiene estadisticas de pagos de un usuario
func (ps *PaymentService) GetUserPaymentStats(userID uint) (map[string]interface{}, error) {
	log.Printf("[SEARCH] [GET_STATS] Obteniendo estadisticas para usuario: %d", userID)

	stats := make(map[string]interface{})

	// Total de pagos por estado
	var statusCounts []struct {
		Estado string
		Count  int64
	}

	err := ps.db.Model(&models.Payment{}).
		Where("usuario_id = ?", userID).
		Select("estado, COUNT(*) as count").
		Group("estado").
		Find(&statusCounts).Error

	if err != nil {
		log.Printf("[ERROR] [GET_STATS] Error al obtener estadisticas por estado: %v", err)
		return nil, err
	}

	stats["payments_by_status"] = statusCounts

	// Total gastado en pagos aprobados
	var totalSpent float64
	ps.db.Model(&models.Payment{}).
		Where("usuario_id = ? AND estado = ?", userID, models.PaymentStatusApproved).
		Select("COALESCE(SUM(monto), 0)").
		Scan(&totalSpent)

	stats["total_spent"] = totalSpent

	// Ultimo pago
	var lastPayment models.Payment
	err = ps.db.Where("usuario_id = ?", userID).
		Order("created_at desc").
		First(&lastPayment).Error

	if err == nil {
		stats["last_payment"] = lastPayment
	}

	log.Printf("[SUCCESS] [GET_STATS] Estadisticas obtenidas - Total gastado: %.2f", totalSpent)
	return stats, nil
}

// UpdateMercadoPagoPaymentStatus - Consulta directa a MercadoPago
func (ps *PaymentService) UpdateMercadoPagoPaymentStatus(payment *models.Payment) error {
	log.Printf("[UPDATE] [UPDATE_MP] Actualizando estado de pago MercadoPago: ID %d", payment.ID)

	// Extraer el payment ID de MercadoPago del transaction ID
	var mpPaymentID string

	if strings.HasPrefix(payment.TransaccionID, "mp_") {
		mpPaymentID = strings.TrimPrefix(payment.TransaccionID, "mp_")
	} else {
		// Si es una preference ID, necesitamos buscar el payment ID
		log.Printf("[SEARCH] [UPDATE_MP] TransaccionID es preference ID: %s", payment.TransaccionID)

		// Buscar pagos por external reference usando la API de MercadoPago
		externalRef := fmt.Sprintf("pago_%d", payment.ID)
		mpPayments, err := ps.searchMercadoPagoPaymentsByReference(externalRef)
		if err != nil || len(mpPayments) == 0 {
			log.Printf("[INFO] [UPDATE_MP] Aún no hay pagos en MercadoPago para referencia: %s (normal si el usuario no ha pagado)", externalRef)
			return nil // No es un error, simplemente el usuario no ha pagado aún
		}

		// Tomar el primer pago encontrado
		mpPaymentID = fmt.Sprintf("%d", mpPayments[0].ID)
		log.Printf("[SUCCESS] [UPDATE_MP] Encontrado payment ID: %s", mpPaymentID)
	}

	// Crear servicio para obtener detalles
	mpService, err := NewMercadoPagoService(ps.config)
	if err != nil {
		return err
	}

	// Obtener detalles del pago desde MercadoPago
	paymentDetails, err := mpService.GetPaymentDetails(mpPaymentID)
	if err != nil {
		log.Printf("[ERROR] [UPDATE_MP] Error al obtener detalles del pago %s: %v", mpPaymentID, err)
		return err
	}

	log.Printf("[INFO] [UPDATE_MP] Detalles obtenidos de MercadoPago:")
	log.Printf("   PaymentID: %d", paymentDetails.ID)
	log.Printf("   Status: %s", paymentDetails.Status)
	log.Printf("   Amount: %.2f %s", paymentDetails.TransactionAmount, paymentDetails.CurrencyID)
	log.Printf("   External Reference: %s", paymentDetails.ExternalReference)
	log.Printf("   Date Created: %s", paymentDetails.DateCreated)
	log.Printf("   Date Approved: %s", paymentDetails.DateApproved)

	// Mapear estado de MercadoPago a nuestro estado
	newStatus := ps.mapMercadoPagoStatusToLocal(paymentDetails.Status)
	oldStatus := payment.Estado

	log.Printf("[UPDATE] [UPDATE_MP] Mapeo de estados:")
	log.Printf("   Estado MP: %s", paymentDetails.Status)
	log.Printf("   Estado actual BD: %s", oldStatus)
	log.Printf("   Estado nuevo: %s", newStatus)

	// Solo actualizar si hay cambio de estado
	if oldStatus != newStatus {
		payment.Estado = newStatus
		payment.TransaccionID = fmt.Sprintf("mp_%d", paymentDetails.ID)
		payment.UpdatedAt = time.Now()

		if err := ps.db.Save(payment).Error; err != nil {
			log.Printf("[ERROR] [UPDATE_MP] Error al guardar cambios en BD: %v", err)
			return err
		}

		log.Printf("[SUCCESS] [UPDATE_MP] Pago ID %d actualizado exitosamente:", payment.ID)
		log.Printf("   Estado: %s -> %s", oldStatus, newStatus)
		log.Printf("   TransaccionID: %s", payment.TransaccionID)

		// Log especial para pagos aprobados
		if newStatus == models.PaymentStatusApproved {
			log.Printf("[SUCCESS] [UPDATE_MP] PAGO APROBADO - Usuario %d ahora tiene acceso al curso %d",
				payment.UsuarioID, payment.CursoID)
		}

	} else {
		log.Printf("[INFO] [UPDATE_MP] Sin cambios de estado para pago ID %d", payment.ID)
	}

	return nil
}

// searchMercadoPagoPaymentsByReference helper function
func (ps *PaymentService) searchMercadoPagoPaymentsByReference(externalRef string) ([]PaymentDetail, error) {
	log.Printf("[SEARCH] [SEARCH_MP] Buscando pagos por referencia: %s", externalRef)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Usar la API de busqueda de MercadoPago
	url := fmt.Sprintf("https://api.mercadopago.com/v1/payments/search?external_reference=%s", externalRef)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+ps.config.MercadoPago.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error en request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[ERROR] [SEARCH_MP] Error en busqueda: Status %d, Body: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("error %d en busqueda", resp.StatusCode)
	}

	var searchResult struct {
		Results []PaymentDetail `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta: %v", err)
	}

	log.Printf("[SUCCESS] [SEARCH_MP] Encontrados %d pagos para referencia %s", len(searchResult.Results), externalRef)

	return searchResult.Results, nil
}

// mapMercadoPagoStatusToLocal helper function
func (ps *PaymentService) mapMercadoPagoStatusToLocal(mpStatus string) string {
	switch strings.ToLower(mpStatus) {
	case "approved":
		return models.PaymentStatusApproved
	case "pending":
		return models.PaymentStatusPending
	case "cancelled", "rejected":
		return models.PaymentStatusRejected
	case "refunded":
		return models.PaymentStatusRefunded
	default:
		log.Printf("[WARNING] [MAP_STATUS] Estado de MP no reconocido: %s", mpStatus)
		return models.PaymentStatusPending
	}
}

// UpdatePendingPayment actualiza el estado de pagos pendientes consultando el proveedor
func (ps *PaymentService) UpdatePendingPayment(payment *models.Payment) {
	log.Printf("[UPDATE] [UPDATE_PENDING] Actualizando pago pendiente: ID %d, Metodo: %s", payment.ID, payment.Metodo)

	// Verificar que el usuario del pago sigue existiendo
	if err := ValidateUserAccess(payment.UsuarioID, ps.config); err != nil {
		log.Printf("[ERROR] [UPDATE_PENDING] Usuario %d del pago %d ya no es valido, marcando como rechazado: %v",
			payment.UsuarioID, payment.ID, err)
		payment.Estado = models.PaymentStatusRejected
		ps.db.Save(payment)
		return
	}

	switch payment.Metodo {
	case models.PaymentMethodPayPal:
		if payment.TransaccionID != "" {
			log.Printf("[SEARCH] [UPDATE_PENDING] Verificando estado con PayPal para transaccion: %s", payment.TransaccionID)
			// Verificar estado actual con PayPal
			orderDetail, err := GetPayPalOrderDetails(payment.TransaccionID)
			if err == nil && orderDetail != nil {
				log.Printf("[SUCCESS] [UPDATE_PENDING] Estado actual de PayPal: %s", orderDetail.Status)
				if orderDetail.Status == "COMPLETED" ||
					orderDetail.Status == "APPROVED" ||
					orderDetail.Status == "PAYER_ACTION_REQUIRED" {
					payment.Estado = models.PaymentStatusApproved
					ps.db.Save(payment)
					log.Printf("[SUCCESS] [UPDATE_PENDING] Pago ID %d actualizado a 'aprobado' segun PayPal", payment.ID)
				}
			} else if err != nil {
				log.Printf("[ERROR] [UPDATE_PENDING] Error al verificar estado con PayPal: %v", err)
			}
		}
	case models.PaymentMethodMercadoPago:
		log.Printf("[SEARCH] [UPDATE_PENDING] Verificando estado con Mercado Pago para pago ID: %d", payment.ID)
		// USAR LA NUEVA FUNCION DE ACTUALIZACION
		if err := ps.UpdateMercadoPagoPaymentStatus(payment); err != nil {
			log.Printf("[ERROR] [UPDATE_PENDING] Error al actualizar estado de MercadoPago: %v", err)
		}
	case models.PaymentMethodCoinbase:
		if payment.TransaccionID != "" {
			log.Printf("[SEARCH] [UPDATE_PENDING] Verificando estado con Coinbase para transaccion: %s", payment.TransaccionID)
			// Aqui podrias implementar verificacion de estado con Coinbase si tienen API para ello
			// Por ahora, Coinbase usa principalmente webhooks
		}
	default:
		log.Printf("[INFO] [UPDATE_PENDING] Metodo %s no requiere actualizacion automatica", payment.Metodo)
	}
}

// ForceUpdatePaymentStatus - Para actualizacion manual
func (ps *PaymentService) ForceUpdatePaymentStatus(paymentID uint) error {
	log.Printf("[UPDATE] [FORCE_UPDATE] Forzando actualizacion de pago ID: %d", paymentID)

	var payment models.Payment
	if err := ps.db.First(&payment, paymentID).Error; err != nil {
		log.Printf("[ERROR] [FORCE_UPDATE] Pago no encontrado: %d", paymentID)
		return fmt.Errorf("pago no encontrado: %d", paymentID)
	}

	log.Printf("[SEARCH] [FORCE_UPDATE] Pago encontrado: ID=%d, Metodo=%s, Estado=%s",
		payment.ID, payment.Metodo, payment.Estado)

	switch payment.Metodo {
	case models.PaymentMethodMercadoPago:
		return ps.UpdateMercadoPagoPaymentStatus(&payment)
	case models.PaymentMethodPayPal:
		// Implementar actualizacion de PayPal si es necesario
		log.Printf("[WARNING] [FORCE_UPDATE] Actualizacion manual de PayPal no implementada")
		return fmt.Errorf("actualizacion manual de PayPal no implementada")
	default:
		log.Printf("[WARNING] [FORCE_UPDATE] Metodo %s no soporta actualizacion automatica", payment.Metodo)
		return fmt.Errorf("metodo %s no soporta actualizacion automatica", payment.Metodo)
	}
}

// BatchUpdatePendingPayments - Actualizar todos los pagos pendientes
func (ps *PaymentService) BatchUpdatePendingPayments() {
	log.Printf("[UPDATE] [BATCH_UPDATE] Iniciando actualizacion masiva de pagos pendientes")

	var pendingPayments []models.Payment
	err := ps.db.Where("estado = ? AND metodo = ?",
		models.PaymentStatusPending, models.PaymentMethodMercadoPago).
		Find(&pendingPayments).Error

	if err != nil {
		log.Printf("[ERROR] [BATCH_UPDATE] Error al obtener pagos pendientes: %v", err)
		return
	}

	log.Printf("[INFO] [BATCH_UPDATE] Encontrados %d pagos pendientes de MercadoPago", len(pendingPayments))

	updated := 0
	errors := 0

	for _, payment := range pendingPayments {
		log.Printf("[UPDATE] [BATCH_UPDATE] Procesando pago ID: %d", payment.ID)

		if err := ps.UpdateMercadoPagoPaymentStatus(&payment); err != nil {
			log.Printf("[ERROR] [BATCH_UPDATE] Error al actualizar pago %d: %v", payment.ID, err)
			errors++
		} else {
			updated++
		}

		// Pequena pausa para no sobrecargar la API
		time.Sleep(500 * time.Millisecond)
	}

	log.Printf("[SUCCESS] [BATCH_UPDATE] Actualizacion completada: %d actualizados, %d errores", updated, errors)
}

// simulatePaymentGateway simula el procesamiento de una pasarela de pago
func (ps *PaymentService) simulatePaymentGateway(paymentID uint, method string) {
	log.Printf("[SIMULATE] Iniciando simulacion para pago ID: %d, metodo: %s", paymentID, method)

	if ps.config.AppEnv != "development" && method != models.PaymentMethodCard && method != models.PaymentMethodTransfer {
		log.Printf("[INFO] [SIMULATE] Simulacion no aplicable en entorno: %s", ps.config.AppEnv)
		return
	}

	// Simular retardo de procesamiento
	time.Sleep(3 * time.Second)

	var payment models.Payment
	if err := ps.db.First(&payment, paymentID).Error; err != nil {
		log.Printf("[ERROR] [SIMULATE] Error al recuperar pago ID %d: %v", paymentID, err)
		return
	}

	// Verificar que el usuario sigue siendo valido antes de aprobar
	if err := ValidateUserAccess(payment.UsuarioID, ps.config); err != nil {
		log.Printf("[ERROR] [SIMULATE] Usuario %d del pago %d ya no es valido: %v",
			payment.UsuarioID, paymentID, err)
		payment.Estado = models.PaymentStatusRejected
		ps.db.Save(&payment)
		return
	}

	// En desarrollo, siempre aprobar el pago
	if ps.config.AppEnv == "development" {
		payment.Estado = models.PaymentStatusApproved
		payment.TransaccionID = ps.generateTransactionID(method)
		log.Printf("[SUCCESS] [SIMULATE] Pago aprobado en modo desarrollo")
	} else {
		// En produccion, simular tasa de aprobacion del 80%
		if rand.Intn(100) < 80 {
			payment.Estado = models.PaymentStatusApproved
			payment.TransaccionID = ps.generateTransactionID(method)
			log.Printf("[SUCCESS] [SIMULATE] Pago aprobado en simulacion (80 por ciento)")
		} else {
			payment.Estado = models.PaymentStatusRejected
			log.Printf("[ERROR] [SIMULATE] Pago rechazado en simulacion (20 por ciento)")
		}
	}

	if err := ps.db.Save(&payment).Error; err != nil {
		log.Printf("[ERROR] [SIMULATE] Error al actualizar estado: %v", err)
	} else {
		log.Printf("[SUCCESS] [SIMULATE] Pago ID %d actualizado a estado: %s", paymentID, payment.Estado)
	}
}

// generateTransactionID genera un ID de transaccion unico
func (ps *PaymentService) generateTransactionID(method string) string {
	timestamp := time.Now().Unix()
	randomPart := rand.Intn(100000)

	// Prefijo segun el metodo de pago
	prefix := "txn"
	switch method {
	case models.PaymentMethodCard:
		prefix = "card"
	case models.PaymentMethodPayPal:
		prefix = "pp"
	case models.PaymentMethodCoinbase:
		prefix = "cb"
	case models.PaymentMethodTransfer:
		prefix = "trf"
	case models.PaymentMethodDev:
		prefix = "dev"
	case models.PaymentMethodStripe:
		prefix = "ch"
	case models.PaymentMethodMercadoPago:
		prefix = "mp"
	}

	transactionID := fmt.Sprintf("%s_%d_%d", prefix, timestamp, randomPart)
	log.Printf("[GENERATE] [GENERATE_TXN] ID de transaccion generado: %s", transactionID)
	return transactionID
}

// ================ MÉTODOS DE GESTIÓN DE EXPIRACIÓN ================

// CancelExpiredPayments cancela todos los pagos que han expirado
func (ps *PaymentService) CancelExpiredPayments() (int, error) {
	log.Printf("[CLEANUP] [EXPIRED] Iniciando cancelación de pagos expirados")

	var expiredPayments []models.Payment
	now := time.Now()

	// Buscar pagos pendientes que han expirado
	err := ps.db.Where("estado = ? AND expires_at IS NOT NULL AND expires_at < ?",
		models.PaymentStatusPending, now).Find(&expiredPayments).Error

	if err != nil {
		log.Printf("[ERROR] [EXPIRED] Error al buscar pagos expirados: %v", err)
		return 0, err
	}

	cancelledCount := 0
	for _, payment := range expiredPayments {
		log.Printf("[CLEANUP] [EXPIRED] Cancelando pago expirado ID: %d (expiró: %v)", 
			payment.ID, payment.ExpiresAt)

		// Actualizar estado a cancelado
		payment.UpdateStatus(models.PaymentStatusCancelled)
		
		if err := ps.db.Save(&payment).Error; err != nil {
			log.Printf("[ERROR] [EXPIRED] Error al cancelar pago ID %d: %v", payment.ID, err)
			continue
		}

		cancelledCount++
		log.Printf("[SUCCESS] [EXPIRED] Pago ID %d cancelado por expiración", payment.ID)
	}

	log.Printf("[SUCCESS] [EXPIRED] Cancelación completada: %d pagos cancelados de %d encontrados", 
		cancelledCount, len(expiredPayments))

	return cancelledCount, nil
}

// MarkPaymentAsExpired marca un pago específico como expirado si cumple las condiciones
func (ps *PaymentService) MarkPaymentAsExpired(paymentID uint) error {
	log.Printf("[EXPIRE] [MARK] Marcando pago ID %d como expirado", paymentID)

	var payment models.Payment
	if err := ps.db.First(&payment, paymentID).Error; err != nil {
		log.Printf("[ERROR] [EXPIRE] Pago no encontrado: %d", paymentID)
		return fmt.Errorf("pago no encontrado: %d", paymentID)
	}

	// Verificar si el pago puede ser expirado
	if !payment.CanBeExpired() {
		log.Printf("[WARNING] [EXPIRE] Pago ID %d no puede ser expirado (estado: %s, tiene expiración: %v)", 
			payment.ID, payment.Estado, payment.ExpiresAt != nil)
		return fmt.Errorf("el pago no puede ser expirado")
	}

	// Verificar si realmente ha expirado
	if !payment.IsExpired() {
		log.Printf("[WARNING] [EXPIRE] Pago ID %d aún no ha expirado (expira: %v)", 
			payment.ID, payment.ExpiresAt)
		return fmt.Errorf("el pago aún no ha expirado")
	}

	// Marcar como cancelado
	payment.UpdateStatus(models.PaymentStatusCancelled)
	
	if err := ps.db.Save(&payment).Error; err != nil {
		log.Printf("[ERROR] [EXPIRE] Error al marcar pago como expirado: %v", err)
		return err
	}

	log.Printf("[SUCCESS] [EXPIRE] Pago ID %d marcado como cancelado por expiración", payment.ID)
	return nil
}

// GetExpiredPayments obtiene todos los pagos que han expirado pero aún no han sido cancelados
func (ps *PaymentService) GetExpiredPayments() ([]models.Payment, error) {
	log.Printf("[SEARCH] [EXPIRED] Buscando pagos expirados")

	var expiredPayments []models.Payment
	now := time.Now()

	err := ps.db.Where("estado = ? AND expires_at IS NOT NULL AND expires_at < ?",
		models.PaymentStatusPending, now).Find(&expiredPayments).Error

	if err != nil {
		log.Printf("[ERROR] [EXPIRED] Error al buscar pagos expirados: %v", err)
		return nil, err
	}

	log.Printf("[SUCCESS] [EXPIRED] Encontrados %d pagos expirados", len(expiredPayments))
	return expiredPayments, nil
}

// ExtendPaymentExpiration extiende la fecha de expiración de un pago
func (ps *PaymentService) ExtendPaymentExpiration(paymentID uint, duration time.Duration) error {
	log.Printf("[EXTEND] [EXPIRY] Extendiendo expiración del pago ID %d por %v", paymentID, duration)

	var payment models.Payment
	if err := ps.db.First(&payment, paymentID).Error; err != nil {
		log.Printf("[ERROR] [EXTEND] Pago no encontrado: %d", paymentID)
		return fmt.Errorf("pago no encontrado: %d", paymentID)
	}

	if !payment.IsPending() {
		log.Printf("[WARNING] [EXTEND] Pago ID %d no está pendiente (estado: %s)", payment.ID, payment.Estado)
		return fmt.Errorf("solo se puede extender la expiración de pagos pendientes")
	}

	// Extender la expiración
	payment.SetExpiration(duration)
	
	if err := ps.db.Save(&payment).Error; err != nil {
		log.Printf("[ERROR] [EXTEND] Error al extender expiración: %v", err)
		return err
	}

	log.Printf("[SUCCESS] [EXTEND] Expiración del pago ID %d extendida hasta: %v", payment.ID, payment.ExpiresAt)
	return nil
}
