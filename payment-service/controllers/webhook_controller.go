// controllers/webhook_controller.go
package controllers

import (
	"context"
	"encoding/json"
	"fmt"
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

type WebhookController struct {
	db             *gorm.DB
	config         *config.Config
	paymentService *services.PaymentService
}

func NewWebhookController(db *gorm.DB, cfg *config.Config) *WebhookController {
	return &WebhookController{
		db:             db,
		config:         cfg,
		paymentService: services.NewPaymentService(db, cfg),
	}
}

// PayPalCallback maneja el callback de PayPal (migrado de callbackPayPal)
func (wc *WebhookController) PayPalCallback(c *gin.Context) {
	// Extraer parÃ¡metros
	pagoID := c.Query("pago_id")
	cursoID := c.Query("curso_id")
	token := c.Query("token")

	log.Printf("Callback PayPal recibido: pagoID=%s, cursoID=%s, token=%s", pagoID, cursoID, token)

	if pagoID == "" || token == "" {
		log.Printf("Error en callback PayPal: parÃ¡metros invÃ¡lidos")
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "ParÃ¡metros invÃ¡lidos. Por favor intenta nuevamente.",
		})
		return
	}

	pagoIDUint, err := strconv.ParseUint(pagoID, 10, 32)
	if err != nil {
		log.Printf("Error en callback PayPal: ID de pago invÃ¡lido - %v", err)
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "ID de pago invÃ¡lido. Por favor intenta nuevamente.",
		})
		return
	}

	var cursoIDUint uint
	if cursoID != "" {
		if cursoIDParsed, err := strconv.ParseUint(cursoID, 10, 32); err == nil {
			cursoIDUint = uint(cursoIDParsed)
		}
	}

	// Buscar el pago
	var payment models.Payment
	if err := wc.db.First(&payment, uint(pagoIDUint)).Error; err != nil {
		log.Printf("Error en callback PayPal: Pago no encontrado (ID: %d) - %v", pagoIDUint, err)
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "Pago no encontrado. Por favor contacta a soporte.",
		})
		return
	}

	// Usar el ID del curso del pago si no se proporcionÃ³
	if cursoIDUint == 0 {
		cursoIDUint = payment.CursoID
	}

	// Crear contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Printf("Intentando capturar orden PayPal con token: %s para pago ID: %d", token, payment.ID)

	// Capturar la orden de PayPal
	captureResult, err := services.CapturePayPalOrder(ctx, token)
	if err != nil {
		log.Printf("Error al capturar orden PayPal: %v", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Error al procesar pago con PayPal. Por favor intenta nuevamente o contacta a soporte.",
		})
		return
	}

	log.Printf("Orden PayPal capturada. Estado: %s", captureResult.Status)

	// Actualizar estado del pago
	estadoActualizado := false
	if captureResult.Status == "COMPLETED" || captureResult.Status == "APPROVED" {
		payment.Estado = models.PaymentStatusApproved
		estadoActualizado = true
	} else if captureResult.Status == "DECLINED" || captureResult.Status == "FAILED" {
		payment.Estado = models.PaymentStatusRejected
		estadoActualizado = true
	}

	// Guardar cambios
	if estadoActualizado {
		if err := wc.db.Save(&payment).Error; err != nil {
			log.Printf("Error al actualizar estado de pago ID %d: %v", payment.ID, err)
		} else {
			log.Printf("Pago ID %d actualizado a estado '%s'", payment.ID, payment.Estado)
		}
	}

	// Redirigir al usuario
	var redirectURL string
	if payment.Estado == models.PaymentStatusApproved {
		redirectURL = fmt.Sprintf("%s/curso/%d", wc.config.FrontendURL, cursoIDUint)
		log.Printf("Pago aprobado, redirigiendo a la pÃ¡gina del curso: %s", redirectURL)
	} else {
		redirectURL = fmt.Sprintf("%s/pagos/fallido?pago_id=%d", wc.config.FrontendURL, pagoIDUint)
		log.Printf("Pago no aprobado, redirigiendo a: %s", redirectURL)
	}

	c.Redirect(http.StatusFound, redirectURL)
}

// PayPalWebhook maneja webhooks de PayPal (migrado de webhookPayPal)
func (wc *WebhookController) PayPalWebhook(c *gin.Context) {
	// Verificar headers de webhook
	paypalEvent := c.GetHeader("Paypal-Transmission-Id")
	if paypalEvent == "" {
		log.Println("Advertencia: Posible llamada no autorizada a webhook de PayPal")
	}

	// Leer y validar el evento
	var event struct {
		EventType string `json:"event_type"`
		Resource  struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		} `json:"resource"`
	}

	body, err := c.GetRawData()
	if err != nil {
		log.Printf("Error al leer cuerpo de webhook PayPal: %v", err)
		utils.SendErrorResponse(c, "Error al leer cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("Error al parsear evento de PayPal: %v", err)
		utils.SendErrorResponse(c, "Error al parsear evento", http.StatusBadRequest)
		return
	}

	log.Printf("Webhook PayPal recibido: Tipo %s, ID %s, Estado %s",
		event.EventType, event.Resource.ID, event.Resource.Status)

	// Procesar eventos relevantes incluyendo cancelaciones
	switch event.EventType {
	case "PAYMENT.CAPTURE.COMPLETED", "CHECKOUT.ORDER.APPROVED":
		// Pago completado/aprobado
	case "PAYMENT.CAPTURE.DENIED", "CHECKOUT.ORDER.CANCELLED", "PAYMENT.CAPTURE.PENDING":
		// Pago denegado/cancelado/pendiente
	default:
		log.Printf("Evento PayPal no manejado: %s", event.EventType)
		c.JSON(http.StatusOK, gin.H{"message": "Evento no manejado"})
		return
	}

	// Buscar pago asociado
	var payment models.Payment
	if err := wc.db.Where("transaccion_id = ?", event.Resource.ID).First(&payment).Error; err != nil {
		log.Printf("Error en webhook PayPal: No se encontrÃ³ pago con transacciÃ³n ID %s", event.Resource.ID)
		utils.SendErrorResponse(c, "Pago no encontrado", http.StatusNotFound)
		return
	}

	// Actualizar estado según el tipo de evento
	estadoAnterior := payment.Estado
	
	switch event.EventType {
	case "PAYMENT.CAPTURE.COMPLETED", "CHECKOUT.ORDER.APPROVED":
		payment.Estado = models.PaymentStatusApproved
	case "PAYMENT.CAPTURE.DENIED", "CHECKOUT.ORDER.CANCELLED":
		payment.Estado = models.PaymentStatusCancelled
	case "PAYMENT.CAPTURE.PENDING":
		payment.Estado = models.PaymentStatusPending
	default:
		// No cambiar estado para eventos no reconocidos
		payment.Estado = estadoAnterior
	}

	if err := wc.db.Save(&payment).Error; err != nil {
		log.Printf("Error al actualizar estado de pago ID %d: %v", payment.ID, err)
		utils.SendErrorResponse(c, "Error de base de datos", http.StatusInternalServerError)
		return
	}

	log.Printf("Pago ID %d actualizado de '%s' a 'aprobado' mediante webhook PayPal",
		payment.ID, estadoAnterior)

	c.JSON(http.StatusOK, gin.H{
		"message": "Webhook de PayPal procesado correctamente",
		"pago_id": payment.ID,
		"estado":  payment.Estado,
	})
}

// CoinbaseWebhook maneja webhooks de Coinbase (migrado de webhookCoinbase)
func (wc *WebhookController) CoinbaseWebhook(c *gin.Context) {
	// Verificar firma de webhook
	signature := c.GetHeader("X-CC-Webhook-Signature")
	if signature == "" && wc.config.AppEnv != "development" {
		log.Println("Error en webhook Coinbase: Firma no proporcionada")
		utils.SendErrorResponse(c, "Firma no proporcionada", http.StatusUnauthorized)
		return
	}

	// Leer y validar evento
	body, err := c.GetRawData()
	if err != nil {
		log.Printf("Error al leer cuerpo de webhook Coinbase: %v", err)
		utils.SendErrorResponse(c, "Error al leer cuerpo del webhook", http.StatusBadRequest)
		return
	}

	var event services.CoinbaseWebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("Error al parsear evento Coinbase: %v", err)
		utils.SendErrorResponse(c, "Error al parsear evento", http.StatusBadRequest)
		return
	}

	log.Printf("Webhook Coinbase recibido: Tipo %s, CÃ³digo %s",
		event.Event.Type, event.Event.Data.Code)

	// Procesar eventos relevantes incluyendo cancelaciones
	switch event.Event.Type {
	case "charge:confirmed":
		// Pago confirmado
	case "charge:failed", "charge:cancelled":
		// Pago fallido o cancelado
	default:
		log.Printf("Evento Coinbase no manejado: %s", event.Event.Type)
		c.JSON(http.StatusOK, gin.H{"message": "Evento no manejado"})
		return
	}

	// Extraer metadata del pago
	metadata := event.Event.Data.Metadata
	var payment models.Payment

	// Buscar por ID de pago en metadata
	if metadata.PagoID > 0 {
		if err := wc.db.First(&payment, metadata.PagoID).Error; err != nil {
			log.Printf("Error en webhook Coinbase: Pago ID %d no encontrado", metadata.PagoID)

			// Buscar por cÃ³digo de transacciÃ³n
			if err := wc.db.Where("transaccion_id = ?", event.Event.Data.Code).First(&payment).Error; err != nil {
				log.Printf("Error en webhook Coinbase: No se encontrÃ³ pago con transacciÃ³n ID %s", event.Event.Data.Code)
				utils.SendErrorResponse(c, "Pago no encontrado", http.StatusNotFound)
				return
			}
		}
	} else {
		// Buscar por cÃ³digo de transacciÃ³n
		if err := wc.db.Where("transaccion_id = ?", event.Event.Data.Code).First(&payment).Error; err != nil {
			log.Printf("Error en webhook Coinbase: No se encontrÃ³ pago con transacciÃ³n ID %s", event.Event.Data.Code)
			utils.SendErrorResponse(c, "Pago no encontrado", http.StatusNotFound)
			return
		}
	}

	// Actualizar estado según el tipo de evento
	estadoAnterior := payment.Estado
	
	switch event.Event.Type {
	case "charge:confirmed":
		payment.Estado = models.PaymentStatusApproved
	case "charge:failed":
		payment.Estado = models.PaymentStatusRejected
	case "charge:cancelled":
		payment.Estado = models.PaymentStatusCancelled
	default:
		// No cambiar estado para eventos no reconocidos
		payment.Estado = estadoAnterior
	}

	// Asegurar ID de transacciÃ³n
	if payment.TransaccionID == "" {
		payment.TransaccionID = event.Event.Data.Code
	}

	if err := wc.db.Save(&payment).Error; err != nil {
		log.Printf("Error al actualizar estado de pago ID %d: %v", payment.ID, err)
		utils.SendErrorResponse(c, "Error de base de datos", http.StatusInternalServerError)
		return
	}

	log.Printf("Pago ID %d actualizado de '%s' a '%s' mediante webhook Coinbase",
		payment.ID, estadoAnterior, payment.Estado)

	c.JSON(http.StatusOK, gin.H{
		"message": "Webhook de Coinbase procesado correctamente",
		"pago_id": payment.ID,
		"estado":  payment.Estado,
	})
}

// ðŸ”¥ MEJORADO: MercadoPagoWebhook con debug completo
func (wc *WebhookController) MercadoPagoWebhook(c *gin.Context) {
	log.Printf("ðŸš€ [MERCADOPAGO_WEBHOOK] Webhook recibido desde IP: %s", c.ClientIP())
	log.Printf("ðŸ“‹ [MERCADOPAGO_WEBHOOK] Headers recibidos:")
	for name, values := range c.Request.Header {
		log.Printf("   %s: %s", name, strings.Join(values, ", "))
	}
	
	// Leer el cuerpo del webhook
	body, err := c.GetRawData()
	if err != nil {
		log.Printf("âŒ [MERCADOPAGO_WEBHOOK] Error al leer cuerpo: %v", err)
		utils.SendErrorResponse(c, "Error al leer webhook", http.StatusBadRequest)
		return
	}
	
	log.Printf("ðŸ“„ [MERCADOPAGO_WEBHOOK] Payload recibido: %s", string(body))
	
	// Parsear el evento de Mercado Pago
	var webhookEvent services.MercadoPagoWebhookPayment
	if err := json.Unmarshal(body, &webhookEvent); err != nil {
		log.Printf("âŒ [MERCADOPAGO_WEBHOOK] Error al parsear JSON: %v", err)
		
		// Intentar parsear como estructura simple para debug
		var simpleEvent map[string]interface{}
		if json.Unmarshal(body, &simpleEvent) == nil {
			log.Printf("ðŸ” [MERCADOPAGO_WEBHOOK] Estructura recibida: %+v", simpleEvent)
		}
		
		utils.SendErrorResponse(c, "Error al parsear webhook", http.StatusBadRequest)
		return
	}
	
	log.Printf("ðŸ“Š [MERCADOPAGO_WEBHOOK] Evento parseado:")
	log.Printf("   Tipo: %s", webhookEvent.Type)
	log.Printf("   AcciÃ³n: %s", webhookEvent.Action)
	log.Printf("   PaymentID: %s", webhookEvent.Data.ID)
	log.Printf("   Live Mode: %t", webhookEvent.LiveMode)
	
	// Solo procesar eventos de pago
	if webhookEvent.Type != "payment" {
		log.Printf("â„¹ï¸ [MERCADOPAGO_WEBHOOK] Evento no relevante: %s - Respondiendo OK", webhookEvent.Type)
		c.JSON(http.StatusOK, gin.H{"message": "Evento no procesado"})
		return
	}
	
	// Solo procesar acciones de actualizaciÃ³n de pago
	if webhookEvent.Action != "payment.created" && webhookEvent.Action != "payment.updated" {
		log.Printf("â„¹ï¸ [MERCADOPAGO_WEBHOOK] AcciÃ³n no relevante: %s - Respondiendo OK", webhookEvent.Action)
		c.JSON(http.StatusOK, gin.H{"message": "AcciÃ³n no procesada"})
		return
	}
	
	log.Printf("ðŸ” [MERCADOPAGO_WEBHOOK] Obteniendo detalles del pago ID: %s", webhookEvent.Data.ID)
	
	// Obtener detalles del pago desde la API de Mercado Pago
	paymentDetails, err := services.GetMercadoPagoPaymentDetails(webhookEvent.Data.ID, wc.config)
	if err != nil {
		log.Printf("âŒ [MERCADOPAGO_WEBHOOK] Error al obtener detalles del pago: %v", err)
		// AÃºn asÃ­ respondemos OK para evitar reenvÃ­os
		c.JSON(http.StatusOK, gin.H{"message": "Error al obtener detalles, pero webhook recibido"})
		return
	}
	
	log.Printf("ðŸ’° [MERCADOPAGO_WEBHOOK] Detalles del pago obtenidos:")
	log.Printf("   Status: %s", paymentDetails.Status)
	log.Printf("   ExternalRef: %s", paymentDetails.ExternalReference)
	log.Printf("   Amount: %.2f", paymentDetails.TransactionAmount)
	log.Printf("   Currency: %s", paymentDetails.CurrencyID)
	log.Printf("   Payment Method: %s", paymentDetails.PaymentMethodID)
	log.Printf("   Date Created: %s", paymentDetails.DateCreated)
	log.Printf("   Date Approved: %s", paymentDetails.DateApproved)
	
	// Buscar el pago en nuestra base de datos usando mÃºltiples estrategias
	var payment models.Payment
	var err_db error
	
	// Estrategia 1: Por referencia externa
	if paymentDetails.ExternalReference != "" {
		paymentID := extractPaymentIDFromReference(paymentDetails.ExternalReference)
		if paymentID > 0 {
			log.Printf("ðŸ” [MERCADOPAGO_WEBHOOK] Buscando pago por ID extraÃ­do de referencia: %d", paymentID)
			err_db = wc.db.First(&payment, paymentID).Error
		}
	}
	
	// Estrategia 2: Por transacciÃ³n ID si la primera fallÃ³
	if err_db != nil {
		log.Printf("ðŸ” [MERCADOPAGO_WEBHOOK] Buscando pago por transacciÃ³n ID que contenga: %s", paymentDetails.ExternalReference)
		err_db = wc.db.Where("transaccion_id LIKE ?", "%"+paymentDetails.ExternalReference+"%").First(&payment).Error
	}
	
	// Estrategia 3: Por metadata si las anteriores fallaron
	if err_db != nil && len(paymentDetails.Metadata) > 0 {
		if pagoID, exists := paymentDetails.Metadata["pago_id"]; exists {
			log.Printf("ðŸ” [MERCADOPAGO_WEBHOOK] Buscando pago por metadata pago_id: %v", pagoID)
			switch v := pagoID.(type) {
			case float64:
				err_db = wc.db.First(&payment, uint(v)).Error
			case string:
				if id, parseErr := strconv.ParseUint(v, 10, 32); parseErr == nil {
					err_db = wc.db.First(&payment, uint(id)).Error
				}
			}
		}
	}
	
	if err_db != nil {
		log.Printf("âŒ [MERCADOPAGO_WEBHOOK] Pago no encontrado en BD:")
		log.Printf("   ExternalRef: %s", paymentDetails.ExternalReference)
		log.Printf("   PaymentID MP: %d", paymentDetails.ID)
		log.Printf("   Metadata: %+v", paymentDetails.Metadata)
		log.Printf("   Error DB: %v", err_db)
		
		// Listar pagos pendientes para debug
		var pendingPayments []models.Payment
		wc.db.Where("estado = ?", models.PaymentStatusPending).Find(&pendingPayments)
		log.Printf("ðŸ“‹ [MERCADOPAGO_WEBHOOK] Pagos pendientes en BD:")
		for _, p := range pendingPayments {
			log.Printf("   ID: %d, TransaccionID: %s, MÃ©todo: %s", p.ID, p.TransaccionID, p.Metodo)
		}
		
		c.JSON(http.StatusOK, gin.H{"message": "Pago no encontrado en BD, pero webhook recibido"})
		return
	}
	
	log.Printf("ðŸŽ¯ [MERCADOPAGO_WEBHOOK] Pago encontrado en BD:")
	log.Printf("   ID: %d", payment.ID)
	log.Printf("   Estado actual: %s", payment.Estado)
	log.Printf("   TransaccionID: %s", payment.TransaccionID)
	log.Printf("   UsuarioID: %d", payment.UsuarioID)
	log.Printf("   CursoID: %d", payment.CursoID)
	log.Printf("   Monto: %.2f", payment.Monto)
	
	// Mapear estado de Mercado Pago a nuestro estado
	newStatus := mapMercadoPagoStatus(paymentDetails.Status)
	oldStatus := payment.Estado
	
	log.Printf("ðŸ”„ [MERCADOPAGO_WEBHOOK] Mapeo de estados:")
	log.Printf("   Estado MP: %s", paymentDetails.Status)
	log.Printf("   Estado anterior: %s", oldStatus)
	log.Printf("   Estado nuevo: %s", newStatus)
	
	// Solo actualizar si hay cambio de estado
	if oldStatus != newStatus {
		payment.Estado = newStatus
		payment.TransaccionID = fmt.Sprintf("mp_%d", paymentDetails.ID)
		
		if err := wc.db.Save(&payment).Error; err != nil {
			log.Printf("âŒ [MERCADOPAGO_WEBHOOK] Error al actualizar pago en BD: %v", err)
			c.JSON(http.StatusOK, gin.H{"message": "Error al actualizar pago, pero webhook recibido"})
			return
		}
		
		log.Printf("âœ… [MERCADOPAGO_WEBHOOK] Pago ID %d actualizado exitosamente:", payment.ID)
		log.Printf("   Estado: %s â†’ %s", oldStatus, newStatus)
		log.Printf("   TransacciÃ³nID: %s", payment.TransaccionID)
		
		// Log adicional para pagos aprobados
		if newStatus == models.PaymentStatusApproved {
			log.Printf("ðŸŽ‰ [MERCADOPAGO_WEBHOOK] PAGO APROBADO - Usuario %d ahora tiene acceso al curso %d", 
				payment.UsuarioID, payment.CursoID)
		}
		
	} else {
		log.Printf("â„¹ï¸ [MERCADOPAGO_WEBHOOK] Sin cambios de estado para pago ID %d", payment.ID)
	}
	
	// Respuesta exitosa (requerida por Mercado Pago)
	c.JSON(http.StatusOK, gin.H{
		"message": "Webhook procesado correctamente",
		"pago_id": payment.ID,
		"estado_anterior": oldStatus,
		"estado_nuevo": newStatus,
		"updated": oldStatus != newStatus,
	})
	
	log.Printf("âœ… [MERCADOPAGO_WEBHOOK] Webhook procesado completamente para pago ID: %d", payment.ID)
}

// GenericWebhook maneja webhooks genÃ©ricos (migrado de webhookPago)
func (wc *WebhookController) GenericWebhook(c *gin.Context) {
	var payload struct {
		PagoID        uint   `json:"pago_id"`
		Estado        string `json:"estado"`
		TransaccionID string `json:"transaccion_id"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Printf("Error en webhook genÃ©rico: JSON invÃ¡lido - %v", err)
		utils.SendErrorResponse(c, "JSON invÃ¡lido", http.StatusBadRequest)
		return
	}

	log.Printf("Webhook genÃ©rico recibido: Pago ID %d, Estado %s", payload.PagoID, payload.Estado)

	var payment models.Payment
	if err := wc.db.First(&payment, payload.PagoID).Error; err != nil {
		log.Printf("Error en webhook genÃ©rico: Pago ID %d no encontrado - %v", payload.PagoID, err)
		utils.SendErrorResponse(c, "Pago no encontrado", http.StatusNotFound)
		return
	}

	// Actualizar estado y posiblemente el ID de transacciÃ³n
	payment.Estado = payload.Estado
	if payload.TransaccionID != "" {
		payment.TransaccionID = payload.TransaccionID
	}

	if err := wc.db.Save(&payment).Error; err != nil {
		log.Printf("Error en webhook genÃ©rico: No se pudo actualizar Pago ID %d - %v", payload.PagoID, err)
		utils.SendErrorResponse(c, "Error de base de datos", http.StatusInternalServerError)
		return
	}

	log.Printf("Pago ID %d actualizado a estado %s mediante webhook genÃ©rico", payment.ID, payment.Estado)

	utils.SendSuccessResponse(c, gin.H{
		"message": "Estado de pago actualizado correctamente",
		"pago_id": payment.ID,
		"estado":  payment.Estado,
	})
}

// ðŸ”¥ FUNCIONES HELPER PARA MERCADO PAGO

// Helper function para extraer ID de pago de la referencia externa
func extractPaymentIDFromReference(externalRef string) uint {
	// Formato esperado: "pago_123"
	parts := strings.Split(externalRef, "_")
	if len(parts) >= 2 {
		if id, err := strconv.ParseUint(parts[len(parts)-1], 10, 32); err == nil {
			return uint(id)
		}
	}
	return 0
}

// Helper function para mapear estados de Mercado Pago a nuestros estados
func mapMercadoPagoStatus(mpStatus string) string {
	switch mpStatus {
	case "approved":
		return models.PaymentStatusApproved
	case "pending":
		return models.PaymentStatusPending
	case "cancelled", "rejected":
		return models.PaymentStatusRejected
	case "refunded":
		return models.PaymentStatusRefunded
	default:
		log.Printf("âš ï¸ [MERCADOPAGO_WEBHOOK] Estado de MP no reconocido: %s", mpStatus)
		return models.PaymentStatusPending
	}
}