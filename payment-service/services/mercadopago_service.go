// services/mercadopago_service.go - IMPLEMENTACION MERCADOPAGO DESDE CERO
package services

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

	"payment-service/config"
	"payment-service/models"
)

// MercadoPagoService - Servicio MercadoPago implementado desde cero
type MercadoPagoService struct {
	config *config.Config
	client *http.Client
}

// NewMercadoPagoService crea una nueva instancia del servicio
func NewMercadoPagoService(cfg *config.Config) (*MercadoPagoService, error) {
	log.Printf("[MP_SERVICE] Iniciando nuevo servicio MercadoPago")
	
	if cfg.MercadoPago.AccessToken == "" {
		log.Printf("[MP_SERVICE] ERROR: Access Token no configurado")
		return nil, fmt.Errorf("MERCADOPAGO_ACCESS_TOKEN no configurado")
	}

	if len(cfg.MercadoPago.AccessToken) < 50 {
		log.Printf("[MP_SERVICE] WARNING: Access Token parece muy corto: %d caracteres", len(cfg.MercadoPago.AccessToken))
	}

	// Crear cliente HTTP con timeout apropiado
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	log.Printf("[MP_SERVICE] Servicio MercadoPago creado exitosamente")
	log.Printf("[MP_SERVICE] Ambiente: %s", cfg.MercadoPago.Environment)
	log.Printf("[MP_SERVICE] Token: %s...", maskToken(cfg.MercadoPago.AccessToken))

	return &MercadoPagoService{
		config: cfg,
		client: httpClient,
	}, nil
}

// maskToken enmascara el token para logging seguro
func maskToken(token string) string {
	if len(token) < 10 {
		return "***INVALID***"
	}
	return token[:10] + "..."
}

// MercadoPagoPreferenceResponse estructura para respuesta de preferencia
type MercadoPagoPreferenceResponse struct {
	ID                string `json:"id"`
	InitPoint         string `json:"init_point"`
	SandboxInitPoint  string `json:"sandbox_init_point"`
	CollectorID       int    `json:"collector_id"`
	ExternalReference string `json:"external_reference"`
}

// CreatePreference crea una preferencia de pago en MercadoPago
func (mp *MercadoPagoService) CreatePreference(payment *models.Payment, course *CourseInfo) (*MercadoPagoPreferenceResponse, error) {
	log.Printf("[MP_SERVICE] Creando preferencia para pago ID: %d", payment.ID)
	log.Printf("[MP_SERVICE] Curso: %s, Monto: %.2f %s", course.Titulo, payment.Monto, payment.Moneda)

	// Determinar moneda y precio
	price := payment.Monto
	currency := payment.Moneda

	// Para Argentina, convertir USD a ARS
	if payment.Moneda == "USD" {
		// Usar tipo de cambio fijo para demo (en producción usar API de cambio)
		exchangeRate := 350.0
		price = payment.Monto * exchangeRate
		currency = "ARS"
		log.Printf("[MP_SERVICE] Conversión USD->ARS: %.2f USD -> %.2f ARS (tasa: %.0f)", 
			payment.Monto, price, exchangeRate)
	}

	// Configurar URLs de retorno usando el controlador genérico
	baseReturnURL := fmt.Sprintf("%s/api/pagos/return", mp.config.BaseURL)
	successURL := fmt.Sprintf("%s?pago_id=%d&provider=mercadopago&status=success", baseReturnURL, payment.ID)
	failureURL := fmt.Sprintf("%s?pago_id=%d&provider=mercadopago&status=failure", baseReturnURL, payment.ID)
	pendingURL := fmt.Sprintf("%s?pago_id=%d&provider=mercadopago&status=pending", baseReturnURL, payment.ID)

	// URL de webhook para notificaciones de MercadoPago
	notificationURL := fmt.Sprintf("%s/api/pagos/mercadopago/webhook", mp.config.BaseURL)
	
	log.Printf("[MP_SERVICE] URLs configuradas:")
	log.Printf("[MP_SERVICE]   Success: %s", successURL)
	log.Printf("[MP_SERVICE]   Failure: %s", failureURL)
	log.Printf("[MP_SERVICE]   Pending: %s", pendingURL)
	log.Printf("[MP_SERVICE]   Webhook: %s", notificationURL)

	// Crear estructura básica de preferencia según documentación oficial de MercadoPago
	preferenceRequest := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"title":       course.Titulo,
				"quantity":    1,
				"unit_price":  price,
				"currency_id": currency,
			},
		},
		"external_reference": fmt.Sprintf("pago_%d", payment.ID),
		"metadata": map[string]interface{}{
			"payment_id": payment.ID,
			"user_id":    payment.UsuarioID,
			"course_id":  payment.CursoID,
		},
	}
	
	// Configurar back_urls independientemente del protocolo para manejar cancelaciones
	preferenceRequest["back_urls"] = map[string]interface{}{
		"success": successURL,
		"failure": failureURL,
		"pending": pendingURL,
	}
	preferenceRequest["auto_return"] = "approved"
	
	// Solo agregar notification_url si tenemos una URL pública válida
	if strings.HasPrefix(mp.config.BaseURL, "https://") && !strings.Contains(mp.config.BaseURL, "localhost") {
		preferenceRequest["notification_url"] = notificationURL
		log.Printf("[MP_SERVICE] Modo PRODUCCIÓN: back_urls, auto_return y notification_url habilitados")
	} else {
		log.Printf("[MP_SERVICE] Modo DESARROLLO: back_urls y auto_return habilitados, notification_url deshabilitado")
		log.Printf("[MP_SERVICE] Para habilitar webhooks, usa ngrok o una URL pública")
	}
	
	log.Printf("[MP_SERVICE] Estructura de preferencia:")
	log.Printf("[MP_SERVICE]   Items: 1 item, precio: %.2f %s", price, currency)
	log.Printf("[MP_SERVICE]   External ref: pago_%d", payment.ID)
	log.Printf("[MP_SERVICE]   Auto return: approved")

	// Serializar la preferencia a JSON
	payload, err := json.Marshal(preferenceRequest)
	if err != nil {
		log.Printf("[MP_SERVICE] ERROR: Error al serializar preferencia: %v", err)
		return nil, fmt.Errorf("error al serializar preferencia: %v", err)
	}

	log.Printf("[MP_SERVICE] Enviando request a API de MercadoPago...")
	log.Printf("[MP_SERVICE] Payload size: %d bytes", len(payload))
	// No loggear el payload completo por seguridad

	// Crear request HTTP a la API de MercadoPago
	apiURL := "https://api.mercadopago.com/checkout/preferences"
	req, err := http.NewRequest("POST", apiURL, bytes.NewReader(payload))
	if err != nil {
		log.Printf("[MP_SERVICE] ERROR: Error creando HTTP request: %v", err)
		return nil, fmt.Errorf("error creando request: %v", err)
	}

	// Configurar headers requeridos
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+mp.config.MercadoPago.AccessToken)
	req.Header.Set("X-Idempotency-Key", fmt.Sprintf("pago_%d_%d", payment.ID, time.Now().Unix()))
	req.Header.Set("User-Agent", "PaymentService/1.0")
	
	log.Printf("[MP_SERVICE] Request configurado para: %s", apiURL)

	// Ejecutar request HTTP
	log.Printf("[MP_SERVICE] Ejecutando request HTTP...")
	resp, err := mp.client.Do(req)
	if err != nil {
		log.Printf("[MP_SERVICE] ERROR: Error en comunicación HTTP: %v", err)
		return nil, fmt.Errorf("error al comunicarse con MercadoPago: %v", err)
	}
	defer resp.Body.Close()

	// Leer respuesta completa
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[MP_SERVICE] ERROR: Error leyendo respuesta: %v", err)
		return nil, fmt.Errorf("error leyendo respuesta: %v", err)
	}

	log.Printf("[MP_SERVICE] Response Status: %d", resp.StatusCode)
	log.Printf("[MP_SERVICE] Response Size: %d bytes", len(body))
	
	// Loggear respuesta solo en modo debug
	if mp.config.AppEnv == "development" {
		log.Printf("[MP_SERVICE] Response Body: %s", string(body))
	}

	// Verificar si la respuesta fue exitosa
	if resp.StatusCode != http.StatusCreated {
		log.Printf("[MP_SERVICE] ERROR: Status code no exitoso: %d", resp.StatusCode)
		log.Printf("[MP_SERVICE] ERROR: Response body: %s", string(body))
		
		// Intentar parsear error de MercadoPago
		var errorResp map[string]interface{}
		if json.Unmarshal(body, &errorResp) == nil {
			if message, exists := errorResp["message"]; exists {
				return nil, fmt.Errorf("error de MercadoPago (%d): %v", resp.StatusCode, message)
			}
		}
		
		return nil, fmt.Errorf("error de MercadoPago (%d): %s", resp.StatusCode, string(body))
	}

	// Decodificar respuesta exitosa de MercadoPago
	var response MercadoPagoPreferenceResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("[MP_SERVICE] ERROR: Error decodificando respuesta JSON: %v", err)
		return nil, fmt.Errorf("error decodificando respuesta: %v", err)
	}

	// Validar que tenemos los datos esenciales
	if response.ID == "" {
		log.Printf("[MP_SERVICE] ERROR: Respuesta no contiene ID de preferencia")
		return nil, fmt.Errorf("respuesta inválida: falta ID de preferencia")
	}

	log.Printf("[MP_SERVICE] SUCCESS: Preferencia creada exitosamente")
	log.Printf("[MP_SERVICE]   Preference ID: %s", response.ID)
	log.Printf("[MP_SERVICE]   Collector ID: %d", response.CollectorID)
	log.Printf("[MP_SERVICE]   External Ref: %s", response.ExternalReference)
	
	if mp.config.MercadoPago.Environment == "sandbox" {
		log.Printf("[MP_SERVICE]   Sandbox URL: %s", response.SandboxInitPoint)
	} else {
		log.Printf("[MP_SERVICE]   Production URL: %s", response.InitPoint)
	}

	return &response, nil
}

// GetCheckoutURL retorna la URL de checkout según el ambiente configurado
func (mp *MercadoPagoService) GetCheckoutURL(response *MercadoPagoPreferenceResponse) string {
	var checkoutURL string
	
	if mp.config.MercadoPago.Environment == "sandbox" {
		checkoutURL = response.SandboxInitPoint
		log.Printf("[MP_SERVICE] Usando URL de sandbox: %s", checkoutURL)
	} else {
		checkoutURL = response.InitPoint
		log.Printf("[MP_SERVICE] Usando URL de producción: %s", checkoutURL)
	}
	
	if checkoutURL == "" {
		log.Printf("[MP_SERVICE] WARNING: URL de checkout vacía")
	}
	
	return checkoutURL
}

// ValidateConfig valida la configuración completa de MercadoPago
func (mp *MercadoPagoService) ValidateConfig() error {
	log.Printf("[MP_SERVICE] Validando configuración...")
	
	if mp.config.MercadoPago.AccessToken == "" {
		log.Printf("[MP_SERVICE] ERROR: Access Token no configurado")
		return fmt.Errorf("MERCADOPAGO_ACCESS_TOKEN no configurado")
	}

	if len(mp.config.MercadoPago.AccessToken) < 50 {
		log.Printf("[MP_SERVICE] ERROR: Access Token muy corto (%d caracteres)", len(mp.config.MercadoPago.AccessToken))
		return fmt.Errorf("Access Token de MercadoPago parece inválido (muy corto)")
	}

	if mp.config.BaseURL == "" {
		log.Printf("[MP_SERVICE] ERROR: BASE_URL no configurado")
		return fmt.Errorf("BASE_URL no configurado")
	}

	env := mp.config.MercadoPago.Environment
	if env != "sandbox" && env != "production" {
		log.Printf("[MP_SERVICE] WARNING: Environment no válido: %s, usando sandbox", env)
		mp.config.MercadoPago.Environment = "sandbox"
	}

	log.Printf("[MP_SERVICE] Configuración válida")
	log.Printf("[MP_SERVICE]   Environment: %s", mp.config.MercadoPago.Environment)
	log.Printf("[MP_SERVICE]   Base URL: %s", mp.config.BaseURL)
	log.Printf("[MP_SERVICE]   Token length: %d", len(mp.config.MercadoPago.AccessToken))
	
	return nil
}

// ProcessPayment integra con el PaymentService principal
func (mp *MercadoPagoService) ProcessPayment(payment *models.Payment, course *CourseInfo) (map[string]interface{}, error) {
	log.Printf("[MP] Procesando pago MercadoPago ID: %d", payment.ID)

	// Validar configuracion
	if err := mp.ValidateConfig(); err != nil {
		return nil, err
	}

	// Crear preferencia
	preference, err := mp.CreatePreference(payment, course)
	if err != nil {
		return nil, err
	}

	// Obtener URL de checkout
	checkoutURL := mp.GetCheckoutURL(preference)

	// Actualizar el pago con el ID de preferencia
	payment.TransaccionID = preference.ID

	log.Printf("[SUCCESS] Pago procesado exitosamente:")
	log.Printf("   Preferencia ID: %s", preference.ID)
	log.Printf("   Checkout URL: %s", checkoutURL)

	return map[string]interface{}{
		"id":                payment.ID,
		"estado":            payment.Estado,
		"checkout_url":      checkoutURL,
		"preference_id":     preference.ID,
		"message":           "Pago creado exitosamente. Redirige al usuario al checkout.",
		"currency":          "ARS",
		"amount":            payment.Monto * 350, // Mostrar el monto en ARS
		"original_amount":   payment.Monto,
		"original_currency": payment.Moneda,
	}, nil
}

// WebhookPayment estructura para webhooks de MercadoPago
type WebhookPayment struct {
	Action string `json:"action"`
	Type   string `json:"type"`
	Data   struct {
		ID string `json:"id"`
	} `json:"data"`
	LiveMode    bool   `json:"live_mode"`
	DateCreated string `json:"date_created"`
	UserID      string `json:"user_id"`
}

// PaymentDetail estructura para detalles de pago de MercadoPago
type PaymentDetail struct {
	ID                int64                  `json:"id"`
	Status            string                 `json:"status"`
	StatusDetail      string                 `json:"status_detail"`
	ExternalReference string                 `json:"external_reference"`
	TransactionAmount float64                `json:"transaction_amount"`
	CurrencyID        string                 `json:"currency_id"`
	PaymentMethodID   string                 `json:"payment_method_id"`
	PaymentTypeID     string                 `json:"payment_type_id"`
	DateCreated       string                 `json:"date_created"`
	DateApproved      string                 `json:"date_approved"`
	Metadata          map[string]interface{} `json:"metadata"`
}

// GetPaymentDetails obtiene los detalles de un pago desde la API de MercadoPago
func (mp *MercadoPagoService) GetPaymentDetails(paymentID string) (*PaymentDetail, error) {
	log.Printf("[MP_SERVICE] Obteniendo detalles del pago ID: %s", paymentID)

	if paymentID == "" {
		return nil, fmt.Errorf("paymentID no puede estar vacío")
	}

	// Construir URL de la API
	apiURL := fmt.Sprintf("https://api.mercadopago.com/v1/payments/%s", paymentID)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Printf("[MP_SERVICE] ERROR: Error creando request: %v", err)
		return nil, fmt.Errorf("error creando request: %v", err)
	}

	// Configurar headers
	req.Header.Set("Authorization", "Bearer "+mp.config.MercadoPago.AccessToken)
	req.Header.Set("User-Agent", "PaymentService/1.0")
	req.Header.Set("Content-Type", "application/json")

	log.Printf("[MP_SERVICE] Ejecutando GET request a: %s", apiURL)
	
	// Ejecutar request
	resp, err := mp.client.Do(req)
	if err != nil {
		log.Printf("[MP_SERVICE] ERROR: Error en request HTTP: %v", err)
		return nil, fmt.Errorf("error en request: %v", err)
	}
	defer resp.Body.Close()

	// Leer respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[MP_SERVICE] ERROR: Error leyendo respuesta: %v", err)
		return nil, fmt.Errorf("error leyendo respuesta: %v", err)
	}

	log.Printf("[MP_SERVICE] Response status: %d", resp.StatusCode)
	
	if resp.StatusCode != http.StatusOK {
		log.Printf("[MP_SERVICE] ERROR: Status no exitoso: %d", resp.StatusCode)
		log.Printf("[MP_SERVICE] ERROR: Response body: %s", string(body))
		return nil, fmt.Errorf("error %d al obtener detalles del pago: %s", resp.StatusCode, string(body))
	}

	// Decodificar respuesta
	var payment PaymentDetail
	if err := json.Unmarshal(body, &payment); err != nil {
		log.Printf("[MP_SERVICE] ERROR: Error decodificando JSON: %v", err)
		log.Printf("[MP_SERVICE] ERROR: Response body: %s", string(body))
		return nil, fmt.Errorf("error decodificando respuesta: %v", err)
	}

	log.Printf("[MP_SERVICE] SUCCESS: Detalles del pago obtenidos")
	log.Printf("[MP_SERVICE]   Payment ID: %d", payment.ID)
	log.Printf("[MP_SERVICE]   Status: %s", payment.Status)
	log.Printf("[MP_SERVICE]   Amount: %.2f %s", payment.TransactionAmount, payment.CurrencyID)
	log.Printf("[MP_SERVICE]   External Ref: %s", payment.ExternalReference)
	
	return &payment, nil
}

// MapStatus mapea los estados de MercadoPago a nuestros estados
func (mp *MercadoPagoService) MapStatus(mpStatus string) string {
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
		log.Printf("[WARNING] Estado desconocido: %s", mpStatus)
		return models.PaymentStatusPending
	}
}

// ExtractPaymentIDFromReference extrae el ID de pago de la referencia externa
func ExtractPaymentIDFromReference(externalRef string) uint {
	// Formato esperado: "pago_123"
	parts := strings.Split(externalRef, "_")
	if len(parts) >= 2 {
		if id, err := strconv.ParseUint(parts[len(parts)-1], 10, 32); err == nil {
			return uint(id)
		}
	}
	return 0
}

// TestConnection verifica que la configuracion de MercadoPago funcione
func (mp *MercadoPagoService) TestConnection() error {
	log.Printf("[MP] Probando conexion con MercadoPago...")

	// Crear una preferencia de prueba muy simple
	testRequest := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"title":       "Test Item",
				"quantity":    1,
				"unit_price":  1.0,
				"currency_id": "ARS",
			},
		},
		"external_reference": fmt.Sprintf("test_%d", time.Now().Unix()),
	}

	payload, _ := json.Marshal(testRequest)
	req, err := http.NewRequest("POST", "https://api.mercadopago.com/checkout/preferences", bytes.NewReader(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+mp.config.MercadoPago.AccessToken)

	resp, err := mp.client.Do(req)
	if err != nil {
		log.Printf("[ERROR] Test de conexion fallo: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		log.Printf("[SUCCESS] Test de conexion exitoso")
		return nil
	}

	body, _ := io.ReadAll(resp.Body)
	log.Printf("[ERROR] Test fallo (%d): %s", resp.StatusCode, string(body))
	return fmt.Errorf("test de conexion fallo: %d", resp.StatusCode)
}

// FUNCIONES HELPER PARA COMPATIBILIDAD CON CONTROLLER

// ValidateMercadoPagoConfig valida la configuracion de MercadoPago
func ValidateMercadoPagoConfig(cfg *config.Config) error {
	service, err := NewMercadoPagoService(cfg)
	if err != nil {
		return err
	}
	return service.ValidateConfig()
}

// MaskMercadoPagoToken enmascara el token para logging seguro
func MaskMercadoPagoToken(token string) string {
	if len(token) < 10 {
		return "***INVALID***"
	}
	return token[:10] + "..." + token[len(token)-4:]
}

// CreateMercadoPagoPreference crea una preferencia de MercadoPago
func CreateMercadoPagoPreference(payment *models.Payment, course *CourseInfo, cfg *config.Config) (*MercadoPagoPreferenceResponse, error) {
	service, err := NewMercadoPagoService(cfg)
	if err != nil {
		return nil, err
	}
	return service.CreatePreference(payment, course)
}

// TestMercadoPagoConnection prueba la conexion con MercadoPago
func TestMercadoPagoConnection(cfg *config.Config) error {
	service, err := NewMercadoPagoService(cfg)
	if err != nil {
		return err
	}
	return service.TestConnection()
}

// GetMercadoPagoAccountInfo obtiene informacion de la cuenta
func GetMercadoPagoAccountInfo(cfg *config.Config) (map[string]interface{}, error) {
	_, err := NewMercadoPagoService(cfg)
	if err != nil {
		return nil, err
	}

	// Hacer una peticion simple para verificar la cuenta
	req, err := http.NewRequest("GET", "https://api.mercadopago.com/v1/account/settings", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+cfg.MercadoPago.AccessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var accountInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&accountInfo); err != nil {
		return nil, err
	}

	return accountInfo, nil
}

// MercadoPagoWebhookPayment estructura para webhooks
type MercadoPagoWebhookPayment struct {
	Action   string `json:"action"`
	Type     string `json:"type"`
	LiveMode bool   `json:"live_mode"`
	Data     struct {
		ID string `json:"id"`
	} `json:"data"`
}

// GetMercadoPagoPaymentDetails obtiene detalles de un pago por ID
func GetMercadoPagoPaymentDetails(paymentID string, cfg *config.Config) (*PaymentDetail, error) {
	service, err := NewMercadoPagoService(cfg)
	if err != nil {
		return nil, err
	}
	return service.GetPaymentDetails(paymentID)
}
