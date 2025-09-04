// controllers/mercadopago_return.go - Manejo de retornos de MercadoPago
package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"payment-service/models"
	"payment-service/services"
)

// HandleMercadoPagoReturn maneja el retorno desde MercadoPago después del pago
func (pc *PaymentController) HandleMercadoPagoReturn(c *gin.Context) {
	log.Printf("[MP_RETURN] === INICIO DE RETORNO MERCADOPAGO ===")
	log.Printf("[MP_RETURN] IP del cliente: %s", c.ClientIP())
	log.Printf("[MP_RETURN] User-Agent: %s", c.GetHeader("User-Agent"))
	
	// Loggear todos los parámetros recibidos
	log.Printf("[MP_RETURN] Parámetros recibidos:")
	for key, values := range c.Request.URL.Query() {
		for _, value := range values {
			log.Printf("[MP_RETURN]   %s: %s", key, value)
		}
	}

	// Extraer parámetros principales
	pagoIDStr := c.Query("pago_id")
	status := c.Query("status")
	collectionID := c.Query("collection_id")
	collectionStatus := c.Query("collection_status")
	preferenceID := c.Query("preference_id")
	externalReference := c.Query("external_reference")
	paymentID := c.Query("payment_id")
	merchantOrderID := c.Query("merchant_order_id")

	log.Printf("[MP_RETURN] Parámetros principales:")
	log.Printf("[MP_RETURN]   pago_id: %s", pagoIDStr)
	log.Printf("[MP_RETURN]   status: %s", status)
	log.Printf("[MP_RETURN]   collection_id: %s", collectionID)
	log.Printf("[MP_RETURN]   collection_status: %s", collectionStatus)
	log.Printf("[MP_RETURN]   preference_id: %s", preferenceID)
	log.Printf("[MP_RETURN]   external_reference: %s", externalReference)
	log.Printf("[MP_RETURN]   payment_id: %s", paymentID)
	log.Printf("[MP_RETURN]   merchant_order_id: %s", merchantOrderID)

	// Validar que tenemos el ID del pago
	if pagoIDStr == "" {
		log.Printf("[MP_RETURN] ERROR: pago_id no proporcionado")
		pc.redirectToError(c, "Parámetros de retorno inválidos")
		return
	}

	// Convertir pago_id a uint
	pagoID, err := strconv.ParseUint(pagoIDStr, 10, 32)
	if err != nil {
		log.Printf("[MP_RETURN] ERROR: pago_id inválido: %s - %v", pagoIDStr, err)
		pc.redirectToError(c, "ID de pago inválido")
		return
	}
	pagoIDUint := uint(pagoID)

	// Buscar el pago en la base de datos
	log.Printf("[MP_RETURN] Buscando pago ID: %d en la base de datos", pagoIDUint)
	var payment models.Payment
	if err := pc.db.First(&payment, pagoIDUint).Error; err != nil {
		log.Printf("[MP_RETURN] ERROR: Pago no encontrado ID: %d - %v", pagoIDUint, err)
		pc.redirectToError(c, "Pago no encontrado")
		return
	}

	log.Printf("[MP_RETURN] Pago encontrado:")
	log.Printf("[MP_RETURN]   ID: %d", payment.ID)
	log.Printf("[MP_RETURN]   Usuario ID: %d", payment.UsuarioID)
	log.Printf("[MP_RETURN]   Curso ID: %d", payment.CursoID)
	log.Printf("[MP_RETURN]   Estado actual: %s", payment.Estado)
	log.Printf("[MP_RETURN]   Método: %s", payment.Metodo)
	log.Printf("[MP_RETURN]   Monto: %.2f %s", payment.Monto, payment.Moneda)
	log.Printf("[MP_RETURN]   TransaccionID: %s", payment.TransaccionID)

	// Verificar que es un pago de MercadoPago
	if payment.Metodo != models.PaymentMethodMercadoPago {
		log.Printf("[MP_RETURN] ERROR: El pago no es de MercadoPago: %s", payment.Metodo)
		pc.redirectToError(c, "Tipo de pago incorrecto")
		return
	}

	// Si tenemos collection_id o payment_id, actualizar el estado desde MercadoPago
	var paymentToCheck string
	if collectionID != "" {
		paymentToCheck = collectionID
		log.Printf("[MP_RETURN] Usando collection_id para verificar: %s", collectionID)
	} else if paymentID != "" {
		paymentToCheck = paymentID
		log.Printf("[MP_RETURN] Usando payment_id para verificar: %s", paymentID)
	}

	// Manejar diferentes tipos de retorno
	switch status {
	case "failure":
		log.Printf("[MP_RETURN] Usuario retornó con status FAILURE - marcando como cancelado")
		payment.Estado = models.PaymentStatusCancelled
		payment.UpdatedAt = time.Now()
		if err := pc.db.Save(&payment).Error; err != nil {
			log.Printf("[MP_RETURN] ERROR: Error guardando pago cancelado: %v", err)
		} else {
			log.Printf("[MP_RETURN] SUCCESS: Pago marcado como cancelado por retorno FAILURE")
		}
		
	case "success", "approved":
		// Actualizar estado del pago consultando directamente a MercadoPago
		if paymentToCheck != "" {
			log.Printf("[MP_RETURN] Consultando estado actualizado desde MercadoPago...")
			if err := pc.updatePaymentFromMercadoPago(paymentToCheck, &payment); err != nil {
				log.Printf("[MP_RETURN] WARNING: Error al actualizar desde MercadoPago: %v", err)
				
				// Como fallback, usar el status de la URL si está disponible y el pago está approved
				if status == "approved" || collectionStatus == "approved" || status == "success" {
					log.Printf("[MP_RETURN] Usando fallback - status de URL para aprobar pago")
					payment.Estado = models.PaymentStatusApproved
					if collectionID != "" {
						payment.TransaccionID = fmt.Sprintf("mp_%s", collectionID)
					} else if paymentID != "" {
						payment.TransaccionID = fmt.Sprintf("mp_%s", paymentID)
					}
					payment.UpdatedAt = time.Now()
					if err := pc.db.Save(&payment).Error; err != nil {
						log.Printf("[MP_RETURN] ERROR: Error guardando pago con fallback: %v", err)
					} else {
						log.Printf("[MP_RETURN] SUCCESS: Pago actualizado por fallback tras fallo de API - Estado: %s, TransaccionID: %s", 
							payment.Estado, payment.TransaccionID)
					}
				}
			}
		} else {
			log.Printf("[MP_RETURN] WARNING: No hay collection_id ni payment_id para consultar MercadoPago")
			
			// Como fallback, usar el status de la URL si está disponible
			if status == "approved" || collectionStatus == "approved" || status == "success" {
				log.Printf("[MP_RETURN] Usando status de URL para aprobar pago")
				payment.Estado = models.PaymentStatusApproved
				// Usar collection_id o payment_id para transaction_id
				if collectionID != "" {
					payment.TransaccionID = fmt.Sprintf("mp_%s", collectionID)
				} else if paymentID != "" {
					payment.TransaccionID = fmt.Sprintf("mp_%s", paymentID)
				}
				payment.UpdatedAt = time.Now()
				if err := pc.db.Save(&payment).Error; err != nil {
					log.Printf("[MP_RETURN] ERROR: Error guardando pago: %v", err)
				} else {
					log.Printf("[MP_RETURN] SUCCESS: Pago actualizado por fallback - Estado: %s, TransaccionID: %s", 
						payment.Estado, payment.TransaccionID)
				}
			}
		}
		
	case "pending":
		log.Printf("[MP_RETURN] Usuario retornó con status PENDING - consultando MercadoPago para estado real")
		if paymentToCheck != "" {
			if err := pc.updatePaymentFromMercadoPago(paymentToCheck, &payment); err != nil {
				log.Printf("[MP_RETURN] WARNING: Error al actualizar desde MercadoPago para pending: %v", err)
			}
		}
		
	default:
		log.Printf("[MP_RETURN] Status desconocido o sin status: %s", status)
		// Para pagos sin información clara, verificar con MercadoPago si es posible
		if paymentToCheck != "" {
			log.Printf("[MP_RETURN] Consultando MercadoPago para status desconocido...")
			if err := pc.updatePaymentFromMercadoPago(paymentToCheck, &payment); err != nil {
				log.Printf("[MP_RETURN] WARNING: Error al actualizar desde MercadoPago: %v", err)
			}
		} else {
			// Si no tenemos información y el usuario retornó sin completar, puede ser abandono
			log.Printf("[MP_RETURN] Posible abandono de pago - manteniendo estado actual")
		}
	}

	// Recargar el pago para obtener el estado más reciente
	pc.db.First(&payment, pagoIDUint)
	log.Printf("[MP_RETURN] Estado final del pago: %s", payment.Estado)

	// Decidir redirección basada en el estado final
	switch payment.Estado {
	case models.PaymentStatusApproved:
		log.Printf("[MP_RETURN] SUCCESS: Pago aprobado, redirigiendo al curso")
		pc.redirectToSuccess(c, &payment)
	case models.PaymentStatusPending:
		log.Printf("[MP_RETURN] INFO: Pago pendiente, redirigiendo a página de espera")
		pc.redirectToPending(c, &payment)
	case models.PaymentStatusRejected, models.PaymentStatusCancelled:
		log.Printf("[MP_RETURN] INFO: Pago rechazado/cancelado, redirigiendo a error")
		pc.redirectToError(c, "El pago fue rechazado o cancelado")
	default:
		log.Printf("[MP_RETURN] WARNING: Estado desconocido: %s", payment.Estado)
		pc.redirectToPending(c, &payment)
	}

	log.Printf("[MP_RETURN] === FIN DE RETORNO MERCADOPAGO ===")
}

// updatePaymentFromMercadoPago actualiza el estado del pago consultando MercadoPago
func (pc *PaymentController) updatePaymentFromMercadoPago(paymentID string, payment *models.Payment) error {
	log.Printf("[MP_RETURN] Actualizando pago desde MercadoPago API...")
	
	// Crear servicio MercadoPago
	mpService, err := services.NewMercadoPagoService(pc.config)
	if err != nil {
		return fmt.Errorf("error creando servicio MercadoPago: %v", err)
	}

	// Obtener detalles del pago
	paymentDetails, err := mpService.GetPaymentDetails(paymentID)
	if err != nil {
		return fmt.Errorf("error obteniendo detalles del pago: %v", err)
	}

	log.Printf("[MP_RETURN] Detalles obtenidos de MercadoPago:")
	log.Printf("[MP_RETURN]   PaymentID: %d", paymentDetails.ID)
	log.Printf("[MP_RETURN]   Status: %s", paymentDetails.Status)
	log.Printf("[MP_RETURN]   Status Detail: %s", paymentDetails.StatusDetail)
	log.Printf("[MP_RETURN]   Amount: %.2f %s", paymentDetails.TransactionAmount, paymentDetails.CurrencyID)
	log.Printf("[MP_RETURN]   External Reference: %s", paymentDetails.ExternalReference)
	log.Printf("[MP_RETURN]   Payment Method: %s", paymentDetails.PaymentMethodID)
	log.Printf("[MP_RETURN]   Date Created: %s", paymentDetails.DateCreated)
	log.Printf("[MP_RETURN]   Date Approved: %s", paymentDetails.DateApproved)

	// Mapear estado de MercadoPago a nuestro estado
	oldStatus := payment.Estado
	newStatus := mapMercadoPagoStatusToLocal(paymentDetails.Status)
	
	log.Printf("[MP_RETURN] Mapeo de estados:")
	log.Printf("[MP_RETURN]   Estado MP: %s", paymentDetails.Status)
	log.Printf("[MP_RETURN]   Estado anterior: %s", oldStatus)
	log.Printf("[MP_RETURN]   Estado nuevo: %s", newStatus)

	// Actualizar solo si hay cambio
	if oldStatus != newStatus {
		payment.Estado = newStatus
		payment.TransaccionID = fmt.Sprintf("mp_%d", paymentDetails.ID)
		payment.UpdatedAt = time.Now()

		if err := pc.db.Save(payment).Error; err != nil {
			log.Printf("[MP_RETURN] ERROR: Error guardando cambios: %v", err)
			return fmt.Errorf("error guardando cambios: %v", err)
		}

		log.Printf("[MP_RETURN] SUCCESS: Pago actualizado:")
		log.Printf("[MP_RETURN]   Estado: %s → %s", oldStatus, newStatus)
		log.Printf("[MP_RETURN]   TransaccionID: %s", payment.TransaccionID)

		if newStatus == models.PaymentStatusApproved {
			log.Printf("[MP_RETURN] 🎉 PAGO APROBADO - Usuario %d tiene acceso al curso %d",
				payment.UsuarioID, payment.CursoID)
		}
	} else {
		log.Printf("[MP_RETURN] INFO: Sin cambios de estado")
	}

	return nil
}

// mapMercadoPagoStatusToLocal mapea estados de MercadoPago a nuestros estados
func mapMercadoPagoStatusToLocal(mpStatus string) string {
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
		log.Printf("[MP_RETURN] WARNING: Estado MP desconocido: %s", mpStatus)
		return models.PaymentStatusPending
	}
}

// redirectToSuccess redirige al usuario a la página del curso cuando el pago es exitoso
func (pc *PaymentController) redirectToSuccess(c *gin.Context, payment *models.Payment) {
	redirectURL := fmt.Sprintf("%s/curso/%d?pago_aprobado=true&pago_id=%d",
		pc.config.FrontendURL, payment.CursoID, payment.ID)
	log.Printf("[MP_RETURN] Redirigiendo a éxito: %s", redirectURL)
	c.Redirect(http.StatusFound, redirectURL)
}

// redirectToPending redirige al usuario a una página de espera cuando el pago está pendiente
func (pc *PaymentController) redirectToPending(c *gin.Context, payment *models.Payment) {
	redirectURL := fmt.Sprintf("%s/pago-pendiente?pago_id=%d",
		pc.config.FrontendURL, payment.ID)
	log.Printf("[MP_RETURN] Redirigiendo a pendiente: %s", redirectURL)
	c.Redirect(http.StatusFound, redirectURL)
}

// redirectToError redirige al usuario a una página de error
func (pc *PaymentController) redirectToError(c *gin.Context, message string) {
	redirectURL := fmt.Sprintf("%s/pago-error?mensaje=%s",
		pc.config.FrontendURL, message)
	log.Printf("[MP_RETURN] Redirigiendo a error: %s", redirectURL)
	c.Redirect(http.StatusFound, redirectURL)
}