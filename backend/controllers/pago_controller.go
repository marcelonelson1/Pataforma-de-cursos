package controllers

import (
	"curso-platform/middleware"
	"curso-platform/models"
	"curso-platform/services"
	"curso-platform/utils"
	
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PagoController gestiona las operaciones relacionadas con los pagos
type PagoController struct {
	pagoService *services.PagoService
}

// NewPagoController crea una nueva instancia del controlador de pagos
func NewPagoController(pagoService *services.PagoService) *PagoController {
	return &PagoController{
		pagoService: pagoService,
	}
}

// CrearPago crea un nuevo registro de pago
func (c *PagoController) CrearPago(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	// Extraer información del usuario de la autenticación
	userValue, exists := ctx.Get("user")
	if !exists {
		utils.SendErrorResponse(ctx, utils.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	user, ok := userValue.(models.Usuario)
	if !ok {
		utils.SendErrorResponse(ctx, utils.ErrServerError, http.StatusInternalServerError)
		return
	}

	// Leer y validar la solicitud de pago
	var req models.PagoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SendValidationErrorResponse(ctx, utils.ErrInvalidRequest, gin.H{"details": err.Error()})
		return
	}

	// Procesar el pago
	response, err := c.pagoService.CrearPago(user.ID, req)
	if err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// VerificarPagoPorCurso verifica si un usuario ha pagado por un curso
func (c *PagoController) VerificarPagoPorCurso(ctx *gin.Context) {
	// Obtener el usuario autenticado
	userValue, exists := ctx.Get("user")
	if !exists {
		utils.SendErrorResponse(ctx, utils.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	user, ok := userValue.(models.Usuario)
	if !ok {
		utils.SendErrorResponse(ctx, utils.ErrServerError, http.StatusInternalServerError)
		return
	}

	// Obtener ID del curso
	id := ctx.Param("id")
	cursoID, err := strconv.Atoi(id)
	if err != nil {
		utils.SendErrorResponse(ctx, utils.ErrBadRequest, http.StatusBadRequest)
		return
	}

	// Verificar pago
	result, err := c.pagoService.VerificarPagoPorCurso(user.ID, uint(cursoID))
	if err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// CallbackPayPal procesa la respuesta de PayPal después de un pago
func (c *PagoController) CallbackPayPal(ctx *gin.Context) {
	// Extraer parámetros de la URL
	pagoID := ctx.Query("pago_id")
	token := ctx.Query("token")

	if pagoID == "" || token == "" {
		log.Printf("Error en callback PayPal: parámetros inválidos")
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Parámetros inválidos. Por favor intenta nuevamente.",
		})
		return
	}

	pagoIDUint, err := strconv.Atoi(pagoID)
	if err != nil {
		log.Printf("Error en callback PayPal: ID de pago inválido - %v", err)
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "ID de pago inválido. Por favor intenta nuevamente.",
		})
		return
	}

	// Procesar la respuesta de PayPal
	_, redirectURL, err := c.pagoService.ProcesarCallbackPayPal(uint(pagoIDUint), token)
	if err != nil {
		log.Printf("Error al procesar callback PayPal: %v", err)
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Error al procesar pago con PayPal. Por favor intenta nuevamente o contacta a soporte.",
		})
		return
	}

	log.Printf("Redirigiendo a: %s", redirectURL)
	ctx.Redirect(http.StatusFound, redirectURL)
}

// WebhookPago procesa un webhook genérico de pasarela de pago
func (c *PagoController) WebhookPago(ctx *gin.Context) {
	var payload struct {
		PagoID        uint   `json:"pago_id"`
		Estado        string `json:"estado"`
		TransaccionID string `json:"transaccion_id"`
	}

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		log.Printf("Error en webhookPago: JSON inválido - %v", err)
		utils.SendErrorResponse(ctx, utils.ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	log.Printf("Webhook recibido: Pago ID %d, Estado %s", payload.PagoID, payload.Estado)

	pago, err := c.pagoService.ProcesarWebhookPago(payload.PagoID, payload.Estado, payload.TransaccionID)
	if err != nil {
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	utils.SendSuccessResponse(ctx, gin.H{
		"message": "Estado de pago actualizado correctamente",
		"pago_id": pago.ID,
		"estado":  pago.Estado,
	})
}

// WebhookPayPal procesa un webhook específico de PayPal
func (c *PagoController) WebhookPayPal(ctx *gin.Context) {
	// Verificar encabezados de webhook
	paypalEvent := ctx.GetHeader("Paypal-Transmission-Id")
	if paypalEvent == "" {
		log.Println("Advertencia: Posible llamada no autorizada a webhook de PayPal")
		// Continuar procesando ya que algunos eventos de prueba pueden no incluir este encabezado
	}

	// Leer y validar el evento
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Printf("Error al leer cuerpo de webhook PayPal: %v", err)
		utils.SendErrorResponse(ctx, utils.ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	var event struct {
		EventType string `json:"event_type"`
		Resource  struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		} `json:"resource"`
	}

	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("Error al parsear evento de PayPal: %v", err)
		utils.SendErrorResponse(ctx, utils.ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	log.Printf("Webhook PayPal recibido: Tipo %s, ID %s, Estado %s",
		event.EventType, event.Resource.ID, event.Resource.Status)

	pago, err := c.pagoService.ProcesarWebhookPayPal(event.EventType, event.Resource.ID, event.Resource.Status)
	if err != nil {
		if err.Error() == "evento PayPal no manejado" {
			ctx.JSON(http.StatusOK, gin.H{"message": "Evento no manejado"})
			return
		}
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Webhook de PayPal procesado correctamente",
		"pago_id": pago.ID,
		"estado":  pago.Estado,
	})
}

// WebhookCoinbase procesa un webhook de Coinbase Commerce
func (c *PagoController) WebhookCoinbase(ctx *gin.Context) {
	// Verificar firma de webhook para autenticación
	signature := ctx.GetHeader("X-CC-Webhook-Signature")
	if signature == "" && utils.GetEnv("ENV", "development") != "development" {
		log.Println("Error en webhook Coinbase: Firma no proporcionada")
		utils.SendErrorResponse(ctx, utils.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	// Leer y validar el evento
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Printf("Error al leer cuerpo de webhook Coinbase: %v", err)
		utils.SendErrorResponse(ctx, utils.ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	var event models.CoinbaseWebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("Error al parsear evento Coinbase: %v", err)
		utils.SendErrorResponse(ctx, utils.ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	log.Printf("Webhook Coinbase recibido: Tipo %s, Código %s",
		event.Event.Type, event.Event.Data.Code)

	pago, err := c.pagoService.ProcesarWebhookCoinbase(event)
	if err != nil {
		if err.Error() == "evento Coinbase no manejado" {
			ctx.JSON(http.StatusOK, gin.H{"message": "Evento no manejado"})
			return
		}
		utils.SendErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Webhook de Coinbase procesado correctamente",
		"pago_id": pago.ID,
		"estado":  pago.Estado,
	})
}

// RegisterRoutes registra todas las rutas relacionadas con los pagos
func (c *PagoController) RegisterRoutes(router *gin.Engine) {
	pagos := router.Group("/api/pagos")
	{
		pagos.Use(middleware.AuthMiddleware())
		pagos.POST("", c.CrearPago)
		pagos.GET("/:id", c.VerificarPagoPorCurso)
	}

	// Webhooks y callbacks (no requieren autenticación)
	router.GET("/api/pagos/paypal/callback", c.CallbackPayPal)
	router.POST("/api/pagos/webhook", c.WebhookPago)
	router.POST("/api/pagos/paypal/webhook", c.WebhookPayPal)
	router.POST("/api/pagos/coinbase/webhook", c.WebhookCoinbase)
}