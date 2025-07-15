package services

import (
	"context"
	"curso-platform/config"
	"curso-platform/models"
	"curso-platform/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/plutov/paypal/v4"
	"gorm.io/gorm"
)

// PagoService maneja la lógica relacionada con los pagos
type PagoService struct {
	PaypalClient *paypal.Client
	CoinbaseAPIKey string
}

// NewPagoService crea una nueva instancia de PagoService
func NewPagoService() *PagoService {
	// Inicializar cliente PayPal
	paypalClientID := utils.GetEnv("PAYPAL_CLIENT_ID", "")
	paypalSecret := utils.GetEnv("PAYPAL_SECRET", "")
	paypalEnv := utils.GetEnv("PAYPAL_ENV", "sandbox")

	var apiBase string
	if paypalEnv == "live" {
		apiBase = paypal.APIBaseLive
	} else {
		apiBase = paypal.APIBaseSandBox
	}

	client, err := paypal.NewClient(paypalClientID, paypalSecret, apiBase)
	if err != nil {
		log.Printf("Advertencia: Error al inicializar cliente PayPal: %v", err)
	}

	// Obtener API key de Coinbase
	coinbaseAPIKey := utils.GetEnv("COINBASE_COMMERCE_API_KEY", "")

	return &PagoService{
		PaypalClient: client,
		CoinbaseAPIKey: coinbaseAPIKey,
	}
}

// CrearPago crea un nuevo registro de pago y procesa según el método de pago
func (s *PagoService) CrearPago(usuarioID uint, req models.PagoRequest) (map[string]interface{}, error) {
	// Validar método de pago
	validMethods := map[string]bool{
		"tarjeta":       true,
		"paypal":        true,
		"coinbase":      true,
		"transferencia": true,
		"stripe":        true,
		"mercadopago":   true,
		"dev":           true,
	}
	if !validMethods[req.Metodo] {
		return nil, errors.New("método de pago no válido")
	}

	// Validar detalles específicos según el método
	if req.Metodo == "tarjeta" && req.DetallesTarjeta == nil {
		return nil, errors.New("se requieren detalles de tarjeta para este método de pago")
	}

	// Establecer moneda predeterminada si no se proporciona
	if req.Moneda == "" {
		req.Moneda = "USD"
	}

	// Buscar el curso en la base de datos
	var curso models.Curso
	if result := config.DB.First(&curso, req.CursoID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, utils.ErrResourceNotFound
		}
		return nil, utils.ErrDatabaseError
	}

	// Verificar si ya existe un pago aprobado para este curso y usuario
	var pagoExistente models.Pago
	if result := config.DB.Where("usuario_id = ? AND curso_id = ? AND estado = ?",
		usuarioID, req.CursoID, "aprobado").First(&pagoExistente); result.Error == nil {

		return map[string]interface{}{
			"message": "Ya tienes acceso a este curso",
			"estado":  "aprobado",
			"pago_id": pagoExistente.ID,
		}, nil
	}

	// Crear el nuevo registro de pago
	pago := models.Pago{
		UsuarioID:     usuarioID,
		CursoID:       req.CursoID,
		Monto:         req.Monto,
		Metodo:        req.Metodo,
		Estado:        "pendiente",
		TransaccionID: "",
		Moneda:        req.Moneda,
	}

	// Guardar el pago en la base de datos
	if result := config.DB.Create(&pago); result.Error != nil {
		return nil, utils.ErrDatabaseError
	}

	// Respuesta base
	response := map[string]interface{}{
		"pago_id": pago.ID,
		"estado":  pago.Estado,
	}

	// Procesar según el método de pago seleccionado
	switch req.Metodo {
	case "dev":
		// Modo desarrollo: simular proceso de pago
		go s.simularPasarelaPago(pago.ID, req.Metodo)
		response["message"] = "Pago en proceso (modo desarrollo)"

	case "paypal":
		// Procesar pago con PayPal
		paypalOrder, err := s.crearOrdenPayPalSimple(pago)
		if err != nil {
			log.Printf("Error al crear orden PayPal: %v", err)
			return nil, errors.New("error al procesar pago con PayPal")
		}

		// Obtener URL de aprobación de PayPal
		paypalApprovalURL := s.getPayPalApprovalURL(paypalOrder)
		if paypalApprovalURL == "" {
			log.Printf("Error: No se pudo obtener URL de aprobación de PayPal")
			return nil, errors.New("error al obtener URL de PayPal")
		}

		// Actualizar ID de transacción
		pago.TransaccionID = paypalOrder.ID
		config.DB.Save(&pago)

		// Añadir URL de redirección a la respuesta
		response["message"] = "Redirigir a PayPal para completar el pago"
		response["checkout_url"] = paypalApprovalURL

	case "coinbase":
		// Procesar pago con Coinbase
		charge, err := s.crearCargoCoinbase(pago, curso)
		if err != nil {
			log.Printf("Error al crear cargo Coinbase: %v", err)
			return nil, errors.New("error al procesar pago con Coinbase")
		}

		pago.TransaccionID = charge.Code
		config.DB.Save(&pago)

		response["message"] = "Redirigir a Coinbase para completar el pago"
		response["checkout_url"] = charge.HostedURL

	case "mercadopago":
		// Simulación de integración con Mercado Pago
		mockCheckoutURL := fmt.Sprintf("https://www.mercadopago.com.%s/checkout?pago_id=%d",
			utils.GetEnv("MERCADOPAGO_COUNTRY", "mx"), pago.ID)

		pago.TransaccionID = fmt.Sprintf("mp_%d", time.Now().Unix())
		config.DB.Save(&pago)

		response["message"] = "Redirigir a Mercado Pago para completar el pago"
		response["checkout_url"] = mockCheckoutURL

	case "stripe":
		// Simulación de integración con Stripe
		mockCheckoutURL := fmt.Sprintf("https://checkout.stripe.com/pay/%s",
			fmt.Sprintf("cs_%d", time.Now().Unix()))

		pago.TransaccionID = fmt.Sprintf("ch_%d", time.Now().Unix())
		config.DB.Save(&pago)

		response["message"] = "Redirigir a Stripe para completar el pago"
		response["checkout_url"] = mockCheckoutURL

	}

	return response, nil
}

// VerificarPagoPorCurso verifica el estado de pago de un usuario para un curso específico
func (s *PagoService) VerificarPagoPorCurso(usuarioID, cursoID uint) (map[string]interface{}, error) {
	// Buscar el pago más reciente para este curso y usuario
	var pago models.Pago
	result := config.DB.Where("usuario_id = ? AND curso_id = ?", usuarioID, cursoID).
		Order("created_at desc").
		First(&pago)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// No hay pago registrado
			return map[string]interface{}{
				"estado":  "no_pagado",
				"message": "No se encontró pago para este curso",
			}, nil
		}
		return nil, utils.ErrDatabaseError
	}

	// Para métodos de pago externos, verificar estado actual si está pendiente
	if pago.Estado == "pendiente" {
		switch pago.Metodo {
		case "paypal":
			if pago.TransaccionID != "" && s.PaypalClient != nil {
				// Verificar estado actual con PayPal
				ctx := context.Background()
				orderDetail, err := s.PaypalClient.GetOrder(ctx, pago.TransaccionID)
				if err == nil && orderDetail != nil {
					log.Printf("Estado actual de PayPal para orden %s: %s", pago.TransaccionID, orderDetail.Status)
					if orderDetail.Status == "COMPLETED" ||
						orderDetail.Status == "APPROVED" ||
						orderDetail.Status == "PAYER_ACTION_REQUIRED" {
						pago.Estado = "aprobado"
						config.DB.Save(&pago)
					}
				} else if err != nil {
					log.Printf("Error al verificar estado con PayPal: %v", err)
				}
			}
		}
	}

	// Responder con el estado del pago
	return map[string]interface{}{
		"estado": pago.Estado,
		"pago":   pago,
	}, nil
}

// ProcesarCallbackPayPal procesa la respuesta de PayPal y actualiza el estado del pago
func (s *PagoService) ProcesarCallbackPayPal(pagoID uint, token string) (*models.Pago, string, error) {
	// Buscar el pago en la base de datos
	var pago models.Pago
	if result := config.DB.First(&pago, pagoID); result.Error != nil {
		return nil, "", utils.ErrResourceNotFound
	}

	// Crear un contexto para la petición a PayPal
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Capturar la orden de PayPal con el token
	captureResult, err := s.PaypalClient.CaptureOrder(ctx, token, paypal.CaptureOrderRequest{})
	if err != nil {
		return nil, "", fmt.Errorf("error al capturar orden PayPal: %v", err)
	}

	// Actualizar el estado del pago según la respuesta de PayPal
	if captureResult.Status == "COMPLETED" || captureResult.Status == "APPROVED" {
		pago.Estado = "aprobado"
	} else if captureResult.Status == "DECLINED" || captureResult.Status == "FAILED" {
		pago.Estado = "rechazado"
	}

	// Guardar cambios en la base de datos
	if err := config.DB.Save(&pago).Error; err != nil {
		return nil, "", utils.ErrDatabaseError
	}

	// Determinar URL de redirección
	frontendURL := utils.GetEnv("FRONTEND_URL", "http://localhost:3000")
	var redirectURL string

	if pago.Estado == "aprobado" {
		// Redirigir a la página del curso
		redirectURL = fmt.Sprintf("%s/curso/%d", frontendURL, pago.CursoID)
	} else {
		// En caso de fallo, redirigir a la página de pago fallido
		redirectURL = fmt.Sprintf("%s/pagos/fallido?pago_id=%d", frontendURL, pago.ID)
	}

	return &pago, redirectURL, nil
}

// ProcesarWebhookPago procesa un webhook genérico de pasarela de pago
func (s *PagoService) ProcesarWebhookPago(pagoID uint, estado, transaccionID string) (*models.Pago, error) {
	var pago models.Pago
	if result := config.DB.First(&pago, pagoID); result.Error != nil {
		return nil, utils.ErrResourceNotFound
	}

	// Actualizar estado y posiblemente el ID de transacción
	pago.Estado = estado
	if transaccionID != "" {
		pago.TransaccionID = transaccionID
	}

	if err := config.DB.Save(&pago).Error; err != nil {
		return nil, utils.ErrDatabaseError
	}

	return &pago, nil
}

// ProcesarWebhookPayPal procesa un webhook específico de PayPal
func (s *PagoService) ProcesarWebhookPayPal(eventType, resourceID, resourceStatus string) (*models.Pago, error) {
	// Procesar solo eventos relevantes
	if eventType != "PAYMENT.CAPTURE.COMPLETED" &&
		eventType != "CHECKOUT.ORDER.APPROVED" {
		return nil, fmt.Errorf("evento PayPal no manejado: %s", eventType)
	}

	// Buscar el pago asociado a esta transacción
	var pago models.Pago
	if result := config.DB.Where("transaccion_id = ?", resourceID).First(&pago); result.Error != nil {
		return nil, utils.ErrResourceNotFound
	}

	// Actualizar estado según el evento
	pago.Estado = "aprobado"

	if err := config.DB.Save(&pago).Error; err != nil {
		return nil, utils.ErrDatabaseError
	}

	return &pago, nil
}

// ProcesarWebhookCoinbase procesa un webhook de Coinbase Commerce
func (s *PagoService) ProcesarWebhookCoinbase(event models.CoinbaseWebhookEvent) (*models.Pago, error) {
	// Procesar solo eventos relevantes
	if event.Event.Type != "charge:confirmed" && event.Event.Type != "charge:failed" {
		return nil, fmt.Errorf("evento Coinbase no manejado: %s", event.Event.Type)
	}

	// Extraer metadata del pago
	metadata := event.Event.Data.Metadata
	var pago models.Pago

	// Buscar primero por ID de pago en metadata
	if metadata.PagoID > 0 {
		if result := config.DB.First(&pago, metadata.PagoID); result.Error != nil {
			// Intentar buscar por código de transacción como alternativa
			if result := config.DB.Where("transaccion_id = ?", event.Event.Data.Code).First(&pago); result.Error != nil {
				return nil, utils.ErrResourceNotFound
			}
		}
	} else {
		// Buscar por código de transacción si no hay ID de pago
		if result := config.DB.Where("transaccion_id = ?", event.Event.Data.Code).First(&pago); result.Error != nil {
			return nil, utils.ErrResourceNotFound
		}
	}

	// Actualizar estado según el evento
	if event.Event.Type == "charge:confirmed" {
		pago.Estado = "aprobado"
	} else {
		pago.Estado = "rechazado"
	}

	// Asegurar que tenemos el código de transacción
	if pago.TransaccionID == "" {
		pago.TransaccionID = event.Event.Data.Code
	}

	if err := config.DB.Save(&pago).Error; err != nil {
		return nil, utils.ErrDatabaseError
	}

	return &pago, nil
}

// Métodos privados auxiliares

// simularPasarelaPago simula un proceso de pago (modo de desarrollo)
func (s *PagoService) simularPasarelaPago(pagoID uint, metodo string) {
	if utils.GetEnv("APP_ENV", "") != "development" && metodo != "tarjeta" && metodo != "transferencia" {
		return
	}

	// Simulamos un retardo para el procesamiento
	time.Sleep(3 * time.Second)

	var pago models.Pago
	if result := config.DB.First(&pago, pagoID); result.Error != nil {
		log.Printf("Error al recuperar pago ID %d para simulación: %v", pagoID, result.Error)
		return
	}

	// En modo desarrollo, siempre aprobamos el pago
	if utils.GetEnv("APP_ENV", "") == "development" {
		pago.Estado = "aprobado"
		pago.TransaccionID = utils.GenerateIDTransaccion(metodo)
	} else {
		// En producción, simulamos una tasa de aprobación del 80%
		if utils.RandomInt(100) < 80 {
			pago.Estado = "aprobado"
			pago.TransaccionID = utils.GenerateIDTransaccion(metodo)
		} else {
			pago.Estado = "rechazado"
		}
	}

	if result := config.DB.Save(&pago); result.Error != nil {
		log.Printf("Error al actualizar estado de pago ID %d: %v", pagoID, result.Error)
	} else {
		log.Printf("Pago ID %d actualizado a estado: %s", pagoID, pago.Estado)
	}
}

// crearOrdenPayPalSimple crea una orden de PayPal
func (s *PagoService) crearOrdenPayPalSimple(pago models.Pago) (*paypal.Order, error) {
	// Crear un contexto para la petición
	ctx := context.Background()

	// Verificar que el cliente PayPal esté inicializado
	if s.PaypalClient == nil {
		return nil, errors.New("cliente PayPal no inicializado")
	}

	// Definir la unidad de compra con los detalles del pago
	purchaseUnit := paypal.PurchaseUnitRequest{
		ReferenceID: fmt.Sprintf("pago_%d", pago.ID),
		Amount: &paypal.PurchaseUnitAmount{
			Currency: pago.Moneda,
			Value:    fmt.Sprintf("%.2f", pago.Monto),
		},
		Description: fmt.Sprintf("Pago para curso ID: %d", pago.CursoID),
	}

	// Obtener las URLs base
	baseURL := utils.GetEnv("BASE_URL", "")
	if baseURL == "" {
		// Si BASE_URL no está configurada, intentar usar una URL basada en la IP del servidor
		baseURL = fmt.Sprintf("http://%s:5000", utils.GetEnv("SERVER_IP", "localhost"))
		log.Printf("BASE_URL no configurada, usando: %s", baseURL)
	}

	frontendURL := utils.GetEnv("FRONTEND_URL", "")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
		log.Printf("FRONTEND_URL no configurada, usando: %s", frontendURL)
	}
	
	// Definir el contexto de la aplicación con URLs de retorno
	appContext := &paypal.ApplicationContext{
		ReturnURL:          fmt.Sprintf("%s/api/pagos/paypal/callback?pago_id=%d&curso_id=%d", baseURL, pago.ID, pago.CursoID),
		CancelURL:          fmt.Sprintf("%s/pagos/cancelado", frontendURL),
		UserAction:         "PAY_NOW",     // Forzar acción de pago inmediato
		ShippingPreference: "NO_SHIPPING", // No requerir dirección de envío
	}

	log.Printf("Creando orden PayPal con ReturnURL: %s, CancelURL: %s",
		appContext.ReturnURL, appContext.CancelURL)

	// Crear la orden usando la API de PayPal
	order, err := s.PaypalClient.CreateOrder(ctx, "CAPTURE", []paypal.PurchaseUnitRequest{purchaseUnit}, nil, appContext)

	if err != nil {
		return nil, fmt.Errorf("error al crear orden de PayPal: %v", err)
	}

	log.Printf("Orden PayPal creada. ID: %s, Estado: %s", order.ID, order.Status)
	return order, nil
}

// getPayPalApprovalURL extrae la URL de aprobación de una orden de PayPal
func (s *PagoService) getPayPalApprovalURL(order *paypal.Order) string {
	if order == nil || len(order.Links) == 0 {
		log.Println("Error: orden de PayPal nula o sin enlaces")
		return ""
	}

	// Primero buscar enlace de aprobación específico
	for _, link := range order.Links {
		if link.Rel == "approve" || link.Rel == "approval_url" {
			log.Printf("Enlace de aprobación encontrado: %s", link.Href)
			return link.Href
		}
	}

	// Si no se encuentra, buscar enlace de payer action
	for _, link := range order.Links {
		if link.Rel == "payer-action" {
			log.Printf("Enlace payer-action encontrado: %s", link.Href)
			return link.Href
		}
	}

	// Buscar cualquier enlace con "checkout" en la URL
	for _, link := range order.Links {
		if strings.Contains(link.Href, "checkout") &&
			(link.Method == "GET" || link.Method == "REDIRECT") {
			log.Printf("Enlace de checkout encontrado: %s", link.Href)
			return link.Href
		}
	}

	// Depuración: imprimir todos los enlaces disponibles
	log.Println("No se encontró enlace de aprobación. Enlaces disponibles:")
	for _, link := range order.Links {
		log.Printf("- Rel: %s, Href: %s, Method: %s", link.Rel, link.Href, link.Method)
	}

	// Último recurso: si hay un solo enlace con método GET, usarlo
	for _, link := range order.Links {
		if link.Method == "GET" && strings.Contains(link.Href, "paypal.com") {
			log.Printf("Usando enlace alternativo: %s", link.Href)
			return link.Href
		}
	}

	return ""
}

// crearCargoCoinbase crea un cargo en Coinbase Commerce
func (s *PagoService) crearCargoCoinbase(pago models.Pago, curso models.Curso) (*models.CoinbaseCharge, error) {
	charge := models.CoinbaseCharge{
		Name:        fmt.Sprintf("Curso: %s", curso.Titulo),
		Description: fmt.Sprintf("Acceso al curso %s", curso.Titulo),
		PricingType: "fixed_price",
		LocalPrice: models.CoinbasePrice{
			Amount:   fmt.Sprintf("%.2f", pago.Monto),
			Currency: pago.Moneda,
		},
		Metadata: models.CoinbaseMetadata{
			PagoID:    pago.ID,
			CursoID:   pago.CursoID,
			UsuarioID: pago.UsuarioID,
		},
		RedirectURL: fmt.Sprintf("%s/pagos/completado", utils.GetEnv("FRONTEND_URL", "")),
		CancelURL:   fmt.Sprintf("%s/pagos/cancelado", utils.GetEnv("FRONTEND_URL", "")),
	}

	client := &http.Client{
		Timeout: 15 * time.Second, // Agregar timeout para evitar bloqueos
	}

	payload, err := json.Marshal(charge)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.commerce.coinbase.com/charges", strings.NewReader(string(payload)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CC-Api-Key", s.CoinbaseAPIKey)
	req.Header.Set("X-CC-Version", "2018-03-22")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("coinbase API returned status %d", resp.StatusCode)
	}

	var response models.CoinbaseChargeResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}