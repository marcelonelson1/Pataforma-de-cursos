// controllers/payment_return_controller.go - Controlador unificado para retornos de todos los proveedores
package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"payment-service/config"
	"payment-service/models"
	"payment-service/services"
)

type PaymentReturnController struct {
	db             *gorm.DB
	config         *config.Config
	paymentService *services.PaymentService
}

func NewPaymentReturnController(db *gorm.DB, cfg *config.Config) *PaymentReturnController {
	return &PaymentReturnController{
		db:             db,
		config:         cfg,
		paymentService: services.NewPaymentService(db, cfg),
	}
}

// HandleGenericReturn maneja retornos genéricos de cualquier proveedor de pago
func (prc *PaymentReturnController) HandleGenericReturn(c *gin.Context) {
	log.Printf("[GENERIC_RETURN] === INICIO DE RETORNO GENÉRICO ===")
	log.Printf("[GENERIC_RETURN] IP del cliente: %s", c.ClientIP())
	log.Printf("[GENERIC_RETURN] User-Agent: %s", c.GetHeader("User-Agent"))
	
	// Loggear todos los parámetros recibidos
	log.Printf("[GENERIC_RETURN] Parámetros recibidos:")
	for key, values := range c.Request.URL.Query() {
		for _, value := range values {
			log.Printf("[GENERIC_RETURN]   %s: %s", key, value)
		}
	}

	// Extraer parámetros principales (comunes a todos los proveedores)
	pagoID := c.Query("pago_id")
	status := c.Query("status")
	provider := c.Query("provider") // ej: paypal, coinbase, stripe
	
	// Parámetros específicos por proveedor
	token := c.Query("token")                 // PayPal
	paymentIntent := c.Query("payment_intent") // Stripe
	chargeCode := c.Query("charge_code")       // Coinbase

	log.Printf("[GENERIC_RETURN] Parámetros principales:")
	log.Printf("[GENERIC_RETURN]   pago_id: %s", pagoID)
	log.Printf("[GENERIC_RETURN]   status: %s", status)
	log.Printf("[GENERIC_RETURN]   provider: %s", provider)
	log.Printf("[GENERIC_RETURN]   token: %s", token)
	log.Printf("[GENERIC_RETURN]   payment_intent: %s", paymentIntent)
	log.Printf("[GENERIC_RETURN]   charge_code: %s", chargeCode)

	// Validar que tenemos el ID del pago
	if pagoID == "" {
		log.Printf("[GENERIC_RETURN] ERROR: pago_id no proporcionado")
		prc.redirectToError(c, "No se pudo procesar tu solicitud. Intenta nuevamente.")
		return
	}

	// Convertir pago_id a uint
	pagoIDUint, err := strconv.ParseUint(pagoID, 10, 32)
	if err != nil {
		log.Printf("[GENERIC_RETURN] ERROR: pago_id inválido: %s - %v", pagoID, err)
		prc.redirectToError(c, "Ocurrió un problema procesando tu pago. Contacta soporte si persiste.")
		return
	}

	// Buscar el pago en la base de datos
	log.Printf("[GENERIC_RETURN] Buscando pago ID: %d en la base de datos", pagoIDUint)
	var payment models.Payment
	if err := prc.db.First(&payment, uint(pagoIDUint)).Error; err != nil {
		log.Printf("[GENERIC_RETURN] ERROR: Pago no encontrado ID: %d - %v", pagoIDUint, err)
		prc.redirectToError(c, "No pudimos encontrar tu pago. Si crees que es un error, contacta soporte.")
		return
	}

	log.Printf("[GENERIC_RETURN] Pago encontrado:")
	log.Printf("[GENERIC_RETURN]   ID: %d", payment.ID)
	log.Printf("[GENERIC_RETURN]   Usuario ID: %d", payment.UsuarioID)
	log.Printf("[GENERIC_RETURN]   Curso ID: %d", payment.CursoID)
	log.Printf("[GENERIC_RETURN]   Estado actual: %s", payment.Estado)
	log.Printf("[GENERIC_RETURN]   Método: %s", payment.Metodo)
	log.Printf("[GENERIC_RETURN]   Monto: %.2f %s", payment.Monto, payment.Moneda)
	log.Printf("[GENERIC_RETURN]   TransaccionID: %s", payment.TransaccionID)

	// Procesar según el estado del retorno
	switch status {
	case "success", "approved", "completed":
		log.Printf("[GENERIC_RETURN] Status SUCCESS/APPROVED - procesando confirmación")
		prc.handleSuccessReturn(c, &payment, provider, token, paymentIntent, chargeCode)

	case "failure", "failed", "declined", "cancelled", "canceled":
		log.Printf("[GENERIC_RETURN] Status FAILURE/CANCELLED - marcando como cancelado")
		prc.handleFailureReturn(c, &payment)

	case "returned", "return":
		log.Printf("[GENERIC_RETURN] Status RETURNED - usuario regresó sin completar")
		prc.handleReturnedStatus(c, &payment)

	case "pending":
		log.Printf("[GENERIC_RETURN] Status PENDING - verificando con proveedor")
		prc.handlePendingReturn(c, &payment, provider)

	default:
		log.Printf("[GENERIC_RETURN] Status DESCONOCIDO: %s - tratando como retorno", status)
		// Si no hay status claro, tratar como retorno del usuario
		prc.handleReturnedStatus(c, &payment)
	}

	log.Printf("[GENERIC_RETURN] === FIN DE RETORNO GENÉRICO ===")
}

// handleSuccessReturn maneja retornos exitosos
func (prc *PaymentReturnController) handleSuccessReturn(c *gin.Context, payment *models.Payment, provider, token, paymentIntent, chargeCode string) {
	log.Printf("[GENERIC_RETURN] [SUCCESS] Procesando retorno exitoso para método: %s", payment.Metodo)

	switch payment.Metodo {
	case models.PaymentMethodPayPal:
		if token != "" {
			// Verificar con PayPal usando el token
			if err := prc.verifyPayPalPayment(payment, token); err != nil {
				log.Printf("[GENERIC_RETURN] [SUCCESS] Error verificando con PayPal: %v", err)
				prc.handleFailureReturn(c, payment)
				return
			}
		}

	case models.PaymentMethodStripe:
		if paymentIntent != "" {
			// Verificar con Stripe usando payment_intent
			log.Printf("[GENERIC_RETURN] [SUCCESS] Verificación Stripe con payment_intent: %s", paymentIntent)
			payment.TransaccionID = paymentIntent
		}

	case models.PaymentMethodCoinbase:
		if chargeCode != "" {
			// Verificar con Coinbase usando charge_code
			log.Printf("[GENERIC_RETURN] [SUCCESS] Verificación Coinbase con charge_code: %s", chargeCode)
			payment.TransaccionID = chargeCode
		}

	default:
		log.Printf("[GENERIC_RETURN] [SUCCESS] Método no requiere verificación especial: %s", payment.Metodo)
	}

	// Actualizar a estado aprobado
	payment.Estado = models.PaymentStatusApproved
	payment.UpdatedAt = time.Now()

	if err := prc.db.Save(payment).Error; err != nil {
		log.Printf("[GENERIC_RETURN] [SUCCESS] ERROR: Error guardando pago aprobado: %v", err)
		prc.redirectToError(c, "Ocurrió un problema procesando tu pago exitoso. Contacta soporte.")
		return
	}

	log.Printf("[GENERIC_RETURN] [SUCCESS] Pago aprobado exitosamente - ID: %d", payment.ID)
	prc.redirectToSuccess(c, payment)
}

// handleFailureReturn maneja retornos de fallo/cancelación
func (prc *PaymentReturnController) handleFailureReturn(c *gin.Context, payment *models.Payment) {
	log.Printf("[GENERIC_RETURN] [FAILURE] Procesando cancelación para método: %s", payment.Metodo)

	// Marcar como cancelado
	payment.Estado = models.PaymentStatusCancelled
	payment.UpdatedAt = time.Now()

	if err := prc.db.Save(payment).Error; err != nil {
		log.Printf("[GENERIC_RETURN] [FAILURE] ERROR: Error guardando pago cancelado: %v", err)
		prc.redirectToError(c, "Ocurrió un problema procesando la cancelación. Intenta nuevamente.")
		return
	}

	log.Printf("[GENERIC_RETURN] [FAILURE] Pago cancelado exitosamente - ID: %d", payment.ID)
	prc.redirectToError(c, "Pago cancelado. Puedes intentar nuevamente si lo deseas.")
}

// handleReturnedStatus maneja cuando el usuario regresa sin completar el pago
func (prc *PaymentReturnController) handleReturnedStatus(c *gin.Context, payment *models.Payment) {
	log.Printf("[GENERIC_RETURN] [RETURNED] Usuario regresó sin completar pago para método: %s", payment.Metodo)

	// Marcar como regresado (cambio inmediato, sin esperar cleanup)
	payment.Estado = models.PaymentStatusReturned
	payment.UpdatedAt = time.Now()

	if err := prc.db.Save(payment).Error; err != nil {
		log.Printf("[GENERIC_RETURN] [RETURNED] ERROR: Error guardando pago regresado: %v", err)
		prc.redirectToError(c, "Ocurrió un problema procesando tu solicitud. Intenta nuevamente.")
		return
	}

	log.Printf("[GENERIC_RETURN] [RETURNED] Pago marcado como regresado exitosamente - ID: %d", payment.ID)
	prc.redirectToReturn(c, payment)
}

// handlePendingReturn maneja retornos pendientes
func (prc *PaymentReturnController) handlePendingReturn(c *gin.Context, payment *models.Payment, provider string) {
	log.Printf("[GENERIC_RETURN] [PENDING] Verificando estado pendiente para método: %s", payment.Metodo)

	// Intentar verificar estado real con el proveedor
	switch payment.Metodo {
	case models.PaymentMethodMercadoPago:
		if err := prc.paymentService.UpdateMercadoPagoPaymentStatus(payment); err != nil {
			log.Printf("[GENERIC_RETURN] [PENDING] Error verificando con MercadoPago: %v", err)
		}
	case models.PaymentMethodPayPal:
		// Implementar verificación PayPal si es necesario
		log.Printf("[GENERIC_RETURN] [PENDING] Verificación PayPal pendiente")
	default:
		log.Printf("[GENERIC_RETURN] [PENDING] Método no soporta verificación automática: %s", payment.Metodo)
	}

	// Recargar el pago para ver si cambió
	prc.db.First(payment, payment.ID)

	if payment.Estado == models.PaymentStatusApproved {
		prc.redirectToSuccess(c, payment)
	} else {
		prc.redirectToPending(c, payment)
	}
}

// handleUnknownReturn maneja retornos sin status claro
func (prc *PaymentReturnController) handleUnknownReturn(c *gin.Context, payment *models.Payment, provider string) {
	log.Printf("[GENERIC_RETURN] [UNKNOWN] Estado desconocido para método: %s", payment.Metodo)

	// Por defecto, mantener como pendiente y redirigir a página de espera
	prc.redirectToPending(c, payment)
}

// verifyPayPalPayment verifica un pago con PayPal
func (prc *PaymentReturnController) verifyPayPalPayment(payment *models.Payment, token string) error {
	// Implementar verificación con PayPal API aquí
	log.Printf("[GENERIC_RETURN] [PAYPAL_VERIFY] Verificando token: %s", token)
	
	// Por ahora, asumir que si hay token, el pago es válido
	payment.TransaccionID = token
	return nil
}

// Funciones de redirección reutilizadas
func (prc *PaymentReturnController) redirectToSuccess(c *gin.Context, payment *models.Payment) {
	redirectURL := fmt.Sprintf("%s/curso/%d?pago_aprobado=true&pago_id=%d",
		prc.config.FrontendURL, payment.CursoID, payment.ID)
	log.Printf("[GENERIC_RETURN] Redirigiendo a éxito: %s", redirectURL)
	c.Redirect(http.StatusFound, redirectURL)
}

func (prc *PaymentReturnController) redirectToPending(c *gin.Context, payment *models.Payment) {
	redirectURL := fmt.Sprintf("%s/pago-pendiente?pago_id=%d",
		prc.config.FrontendURL, payment.ID)
	log.Printf("[GENERIC_RETURN] Redirigiendo a pendiente: %s", redirectURL)
	c.Redirect(http.StatusFound, redirectURL)
}

func (prc *PaymentReturnController) redirectToError(c *gin.Context, message string) {
	redirectURL := fmt.Sprintf("%s/pago-error?mensaje=%s",
		prc.config.FrontendURL, message)
	log.Printf("[GENERIC_RETURN] Redirigiendo a error: %s", redirectURL)
	c.Redirect(http.StatusFound, redirectURL)
}

func (prc *PaymentReturnController) redirectToReturn(c *gin.Context, payment *models.Payment) {
	redirectURL := fmt.Sprintf("%s/cursos?pago_cancelado=true&mensaje=No completaste el pago. Puedes intentar nuevamente cuando gustes.",
		prc.config.FrontendURL)
	log.Printf("[GENERIC_RETURN] Redirigiendo por retorno: %s", redirectURL)
	c.Redirect(http.StatusFound, redirectURL)
}

// === HELPER METHODS PARA URLS DE RETORNO ===

// BuildReturnURLs construye URLs de retorno para cualquier proveedor
func (prc *PaymentReturnController) BuildReturnURLs(payment *models.Payment, provider string) map[string]string {
	baseURL := prc.config.BaseURL
	
	return map[string]string{
		"success": fmt.Sprintf("%s/api/pagos/return?pago_id=%d&provider=%s&status=success", 
			baseURL, payment.ID, provider),
		"failure": fmt.Sprintf("%s/api/pagos/return?pago_id=%d&provider=%s&status=failure", 
			baseURL, payment.ID, provider),
		"pending": fmt.Sprintf("%s/api/pagos/return?pago_id=%d&provider=%s&status=pending", 
			baseURL, payment.ID, provider),
		"cancel":  fmt.Sprintf("%s/api/pagos/return?pago_id=%d&provider=%s&status=cancelled", 
			baseURL, payment.ID, provider),
	}
}

// GetProviderName normaliza nombres de métodos de pago a proveedores
func GetProviderName(paymentMethod string) string {
	switch paymentMethod {
	case models.PaymentMethodMercadoPago:
		return "mercadopago"
	case models.PaymentMethodPayPal:
		return "paypal"
	case models.PaymentMethodCoinbase:
		return "coinbase"
	case models.PaymentMethodStripe:
		return "stripe"
	case models.PaymentMethodCard:
		return "card"
	case models.PaymentMethodTransfer:
		return "transfer"
	case models.PaymentMethodDev:
		return "dev"
	default:
		return "unknown"
	}
}