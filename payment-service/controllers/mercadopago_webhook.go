// controllers/mercadopago_webhook.go - Manejo de webhooks de MercadoPago
package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"payment-service/models"
	"payment-service/services"
	"payment-service/utils"
)

// MercadoPagoWebhookData estructura para webhooks de MercadoPago
type MercadoPagoWebhookData struct {
	Action   string `json:"action"`
	Type     string `json:"type"`
	LiveMode bool   `json:"live_mode"`
	Data     struct {
		ID string `json:"id"`
	} `json:"data"`
	DateCreated string `json:"date_created"`
	UserID      string `json:"user_id"`
}

// HandleMercadoPagoWebhookNew maneja webhooks de MercadoPago con implementación completa
func (wc *WebhookController) HandleMercadoPagoWebhookNew(c *gin.Context) {
	log.Printf("[MP_WEBHOOK] === INICIO WEBHOOK MERCADOPAGO ===")
	log.Printf("[MP_WEBHOOK] Timestamp: %v", time.Now())
	log.Printf("[MP_WEBHOOK] IP del cliente: %s", c.ClientIP())
	log.Printf("[MP_WEBHOOK] User-Agent: %s", c.GetHeader("User-Agent"))

	// Loggear todos los headers recibidos
	log.Printf("[MP_WEBHOOK] Headers recibidos:")
	for name, values := range c.Request.Header {
		log.Printf("[MP_WEBHOOK]   %s: %s", name, strings.Join(values, ", "))
	}

	// Leer el cuerpo del webhook
	body, err := c.GetRawData()
	if err != nil {
		log.Printf("[MP_WEBHOOK] ERROR: Error al leer cuerpo del webhook: %v", err)
		utils.SendErrorResponse(c, "Error al leer webhook", http.StatusBadRequest)
		return
	}

	log.Printf("[MP_WEBHOOK] Payload recibido (%d bytes): %s", len(body), string(body))

	// Parsear el evento de MercadoPago
	var webhookEvent MercadoPagoWebhookData
	if err := json.Unmarshal(body, &webhookEvent); err != nil {
		log.Printf("[MP_WEBHOOK] ERROR: Error al parsear JSON: %v", err)
		
		// Intentar parsear como estructura genérica para debug
		var genericEvent map[string]interface{}
		if json.Unmarshal(body, &genericEvent) == nil {
			log.Printf("[MP_WEBHOOK] Estructura genérica recibida: %+v", genericEvent)
		}
		
		utils.SendErrorResponse(c, "Error al parsear webhook", http.StatusBadRequest)
		return
	}

	log.Printf("[MP_WEBHOOK] Evento parseado exitosamente:")
	log.Printf("[MP_WEBHOOK]   Tipo: %s", webhookEvent.Type)
	log.Printf("[MP_WEBHOOK]   Acción: %s", webhookEvent.Action)
	log.Printf("[MP_WEBHOOK]   PaymentID: %s", webhookEvent.Data.ID)
	log.Printf("[MP_WEBHOOK]   Live Mode: %t", webhookEvent.LiveMode)
	log.Printf("[MP_WEBHOOK]   Date Created: %s", webhookEvent.DateCreated)
	log.Printf("[MP_WEBHOOK]   User ID: %s", webhookEvent.UserID)

	// Validar que es un evento de pago
	if webhookEvent.Type != "payment" {
		log.Printf("[MP_WEBHOOK] INFO: Evento no relevante (tipo: %s), respondiendo OK", webhookEvent.Type)
		c.JSON(http.StatusOK, gin.H{
			"message": "Evento no procesado",
			"type":    webhookEvent.Type,
		})
		return
	}

	// Validar que es una acción que nos interesa
	relevantActions := []string{"payment.created", "payment.updated"}
	isRelevantAction := false
	for _, action := range relevantActions {
		if webhookEvent.Action == action {
			isRelevantAction = true
			break
		}
	}

	if !isRelevantAction {
		log.Printf("[MP_WEBHOOK] INFO: Acción no relevante (%s), respondiendo OK", webhookEvent.Action)
		c.JSON(http.StatusOK, gin.H{
			"message": "Acción no procesada",
			"action":  webhookEvent.Action,
		})
		return
	}

	// Validar que tenemos un ID de pago
	if webhookEvent.Data.ID == "" {
		log.Printf("[MP_WEBHOOK] ERROR: No se proporcionó ID de pago")
		c.JSON(http.StatusOK, gin.H{
			"message": "ID de pago faltante",
		})
		return
	}

	log.Printf("[MP_WEBHOOK] Procesando pago ID: %s", webhookEvent.Data.ID)

	// Obtener detalles del pago desde MercadoPago
	paymentDetails, err := wc.getMercadoPagoPaymentDetails(webhookEvent.Data.ID)
	if err != nil {
		log.Printf("[MP_WEBHOOK] ERROR: Error obteniendo detalles del pago: %v", err)
		// Aún así respondemos OK para evitar reenvíos
		c.JSON(http.StatusOK, gin.H{
			"message": "Error obteniendo detalles, pero webhook recibido",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("[MP_WEBHOOK] Detalles del pago obtenidos:")
	log.Printf("[MP_WEBHOOK]   PaymentID: %d", paymentDetails.ID)
	log.Printf("[MP_WEBHOOK]   Status: %s", paymentDetails.Status)
	log.Printf("[MP_WEBHOOK]   Status Detail: %s", paymentDetails.StatusDetail)
	log.Printf("[MP_WEBHOOK]   External Ref: %s", paymentDetails.ExternalReference)
	log.Printf("[MP_WEBHOOK]   Amount: %.2f %s", paymentDetails.TransactionAmount, paymentDetails.CurrencyID)
	log.Printf("[MP_WEBHOOK]   Payment Method: %s", paymentDetails.PaymentMethodID)
	log.Printf("[MP_WEBHOOK]   Payment Type: %s", paymentDetails.PaymentTypeID)
	log.Printf("[MP_WEBHOOK]   Date Created: %s", paymentDetails.DateCreated)
	log.Printf("[MP_WEBHOOK]   Date Approved: %s", paymentDetails.DateApproved)

	// Buscar el pago en nuestra base de datos
	payment, err := wc.findPaymentInDatabase(paymentDetails)
	if err != nil {
		log.Printf("[MP_WEBHOOK] ERROR: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"message": "Pago no encontrado en BD, pero webhook recibido",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("[MP_WEBHOOK] Pago encontrado en BD:")
	log.Printf("[MP_WEBHOOK]   ID: %d", payment.ID)
	log.Printf("[MP_WEBHOOK]   Usuario ID: %d", payment.UsuarioID)
	log.Printf("[MP_WEBHOOK]   Curso ID: %d", payment.CursoID)
	log.Printf("[MP_WEBHOOK]   Estado actual: %s", payment.Estado)
	log.Printf("[MP_WEBHOOK]   TransaccionID: %s", payment.TransaccionID)
	log.Printf("[MP_WEBHOOK]   Monto: %.2f %s", payment.Monto, payment.Moneda)

	// Mapear y actualizar estado
	oldStatus := payment.Estado
	newStatus := mapMercadoPagoStatusToLocalWebhook(paymentDetails.Status)

	log.Printf("[MP_WEBHOOK] Mapeo de estados:")
	log.Printf("[MP_WEBHOOK]   Estado MP: %s", paymentDetails.Status)
	log.Printf("[MP_WEBHOOK]   Estado anterior: %s", oldStatus)
	log.Printf("[MP_WEBHOOK]   Estado nuevo: %s", newStatus)

	// Actualizar solo si hay cambio de estado
	if oldStatus != newStatus {
		payment.Estado = newStatus
		payment.TransaccionID = fmt.Sprintf("mp_%d", paymentDetails.ID)
		payment.UpdatedAt = time.Now()

		if err := wc.db.Save(payment).Error; err != nil {
			log.Printf("[MP_WEBHOOK] ERROR: Error guardando cambios en BD: %v", err)
			c.JSON(http.StatusOK, gin.H{
				"message": "Error guardando cambios, pero webhook recibido",
				"error":   err.Error(),
			})
			return
		}

		log.Printf("[MP_WEBHOOK] SUCCESS: Pago actualizado exitosamente")
		log.Printf("[MP_WEBHOOK]   Estado: %s → %s", oldStatus, newStatus)
		log.Printf("[MP_WEBHOOK]   TransaccionID: %s", payment.TransaccionID)

		// Log especial para pagos aprobados
		if newStatus == models.PaymentStatusApproved {
			log.Printf("[MP_WEBHOOK] 🎉 PAGO APROBADO - Usuario %d ahora tiene acceso al curso %d",
				payment.UsuarioID, payment.CursoID)
		}
	} else {
		log.Printf("[MP_WEBHOOK] INFO: Sin cambios de estado para pago ID %d", payment.ID)
	}

	// Respuesta exitosa (requerida por MercadoPago)
	response := gin.H{
		"message":          "Webhook procesado correctamente",
		"payment_id":       payment.ID,
		"mp_payment_id":    paymentDetails.ID,
		"estado_anterior":  oldStatus,
		"estado_nuevo":     newStatus,
		"actualizado":      oldStatus != newStatus,
		"timestamp":        time.Now(),
	}

	c.JSON(http.StatusOK, response)
	log.Printf("[MP_WEBHOOK] Respuesta enviada: %+v", response)
	log.Printf("[MP_WEBHOOK] === FIN WEBHOOK MERCADOPAGO ===")
}

// getMercadoPagoPaymentDetails obtiene detalles del pago desde MercadoPago
func (wc *WebhookController) getMercadoPagoPaymentDetails(paymentID string) (*services.PaymentDetail, error) {
	log.Printf("[MP_WEBHOOK] Obteniendo detalles del pago %s desde MercadoPago", paymentID)
	
	mpService, err := services.NewMercadoPagoService(wc.config)
	if err != nil {
		return nil, fmt.Errorf("error creando servicio MercadoPago: %v", err)
	}

	return mpService.GetPaymentDetails(paymentID)
}

// findPaymentInDatabase busca el pago en la base de datos usando múltiples estrategias
func (wc *WebhookController) findPaymentInDatabase(paymentDetails *services.PaymentDetail) (*models.Payment, error) {
	log.Printf("[MP_WEBHOOK] Buscando pago en base de datos...")
	
	var payment models.Payment
	var err error

	// Estrategia 1: Por referencia externa
	if paymentDetails.ExternalReference != "" {
		paymentID := extractPaymentIDFromReferenceWebhook(paymentDetails.ExternalReference)
		if paymentID > 0 {
			log.Printf("[MP_WEBHOOK] Buscando por ID extraído de referencia: %d", paymentID)
			err = wc.db.First(&payment, paymentID).Error
			if err == nil {
				log.Printf("[MP_WEBHOOK] Pago encontrado por referencia externa")
				return &payment, nil
			}
			log.Printf("[MP_WEBHOOK] No encontrado por referencia externa: %v", err)
		}
	}

	// Estrategia 2: Por TransaccionID que contenga la referencia externa
	if paymentDetails.ExternalReference != "" {
		log.Printf("[MP_WEBHOOK] Buscando por TransaccionID que contenga: %s", paymentDetails.ExternalReference)
		err = wc.db.Where("transaccion_id LIKE ?", "%"+paymentDetails.ExternalReference+"%").First(&payment).Error
		if err == nil {
			log.Printf("[MP_WEBHOOK] Pago encontrado por TransaccionID")
			return &payment, nil
		}
		log.Printf("[MP_WEBHOOK] No encontrado por TransaccionID: %v", err)
	}

	// Estrategia 3: Por metadata si está disponible
	if len(paymentDetails.Metadata) > 0 {
		if pagoID, exists := paymentDetails.Metadata["payment_id"]; exists {
			log.Printf("[MP_WEBHOOK] Buscando por metadata payment_id: %v", pagoID)
			switch v := pagoID.(type) {
			case float64:
				err = wc.db.First(&payment, uint(v)).Error
			case string:
				if id, parseErr := strconv.ParseUint(v, 10, 32); parseErr == nil {
					err = wc.db.First(&payment, uint(id)).Error
				}
			}
			if err == nil {
				log.Printf("[MP_WEBHOOK] Pago encontrado por metadata")
				return &payment, nil
			}
			log.Printf("[MP_WEBHOOK] No encontrado por metadata: %v", err)
		}
	}

	// Estrategia 4: Por preference_id en TransaccionID (nuevo)
	log.Printf("[MP_WEBHOOK] Buscando por preference_id en TransaccionID...")
	var allPayments []models.Payment
	wc.db.Where("metodo = ? AND estado = ?", models.PaymentMethodMercadoPago, models.PaymentStatusPending).Find(&allPayments)
	
	for _, p := range allPayments {
		if p.TransaccionID != "" {
			log.Printf("[MP_WEBHOOK] Verificando pago ID %d con preference: %s", p.ID, p.TransaccionID)
			// Verificar si este preference_id corresponde al pago que llegó en el webhook
			// Para esto, podemos usar el payment ID de MercadoPago para actualizar
			if paymentDetails.ID > 0 {
				payment = p
				log.Printf("[MP_WEBHOOK] ✅ PAGO ENCONTRADO POR PREFERENCE - ID: %d", p.ID)
				return &payment, nil
			}
		}
	}

	// Estrategia 5: Por TransaccionID con formato mp_
	mpTransactionID := fmt.Sprintf("mp_%d", paymentDetails.ID)
	log.Printf("[MP_WEBHOOK] Buscando por TransaccionID exacto: %s", mpTransactionID)
	err = wc.db.Where("transaccion_id = ?", mpTransactionID).First(&payment).Error
	if err == nil {
		log.Printf("[MP_WEBHOOK] Pago encontrado por TransaccionID exacto")
		return &payment, nil
	}

	// Log de debug - mostrar pagos pendientes
	var pendingPayments []models.Payment
	wc.db.Where("estado = ? AND metodo = ?", models.PaymentStatusPending, models.PaymentMethodMercadoPago).Find(&pendingPayments)
	log.Printf("[MP_WEBHOOK] Pagos pendientes de MercadoPago en BD:")
	for _, p := range pendingPayments {
		log.Printf("[MP_WEBHOOK]   ID: %d, TransaccionID: %s, Usuario: %d, Curso: %d",
			p.ID, p.TransaccionID, p.UsuarioID, p.CursoID)
	}

	return nil, fmt.Errorf("pago no encontrado en BD - External Ref: %s, MP PaymentID: %d",
		paymentDetails.ExternalReference, paymentDetails.ID)
}

// extractPaymentIDFromReferenceWebhook extrae el ID del pago de la referencia externa
func extractPaymentIDFromReferenceWebhook(externalRef string) uint {
	log.Printf("[MP_WEBHOOK] Extrayendo ID de referencia: %s", externalRef)
	
	// Formato esperado: "pago_123"
	parts := strings.Split(externalRef, "_")
	if len(parts) >= 2 {
		if id, err := strconv.ParseUint(parts[len(parts)-1], 10, 32); err == nil {
			log.Printf("[MP_WEBHOOK] ID extraído: %d", uint(id))
			return uint(id)
		}
	}
	
	log.Printf("[MP_WEBHOOK] No se pudo extraer ID de: %s", externalRef)
	return 0
}

// mapMercadoPagoStatusToLocalWebhook mapea estados de MercadoPago a nuestros estados
func mapMercadoPagoStatusToLocalWebhook(mpStatus string) string {
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
		log.Printf("[MP_WEBHOOK] WARNING: Estado MP desconocido: %s", mpStatus)
		return models.PaymentStatusPending
	}
}