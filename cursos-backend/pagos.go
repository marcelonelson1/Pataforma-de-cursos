// pagos.go
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/plutov/paypal/v4"
	"gorm.io/gorm"
)

type PagoRequest struct {
	CursoID         uint             `json:"curso_id" binding:"required"`
	Monto           float64          `json:"monto" binding:"required"`
	Metodo          string           `json:"metodo" binding:"required"`
	DetallesTarjeta *DetallesTarjeta `json:"detalles_tarjeta,omitempty"`
	Moneda          string           `json:"moneda,omitempty"`
}

type DetallesTarjeta struct {
	Numero     string `json:"numero"`
	Expiracion string `json:"expiracion"`
	CVV        string `json:"cvv"`
}

type CoinbaseCharge struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	PricingType string           `json:"pricing_type"`
	LocalPrice  CoinbasePrice    `json:"local_price"`
	Metadata    CoinbaseMetadata `json:"metadata"`
	HostedURL   string           `json:"hosted_url"`
	RedirectURL string           `json:"redirect_url"`
	CancelURL   string           `json:"cancel_url"`
	Code        string           `json:"code"`
}

type CoinbasePrice struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type CoinbaseMetadata struct {
	PagoID    uint `json:"pago_id"`
	CursoID   uint `json:"curso_id"`
	UsuarioID uint `json:"usuario_id"`
}

type CoinbaseChargeResponse struct {
	Data CoinbaseCharge `json:"data"`
}

type CoinbaseWebhookEvent struct {
	Event struct {
		Type string `json:"type"`
		Data struct {
			Code     string           `json:"code"`
			Metadata CoinbaseMetadata `json:"metadata"`
			Timeline []struct {
				Status string `json:"status"`
				Time   string `json:"time"`
			} `json:"timeline"`
		} `json:"data"`
	} `json:"event"`
}

var (
	paypalClient   *paypal.Client
	coinbaseAPIKey string
)

func initPaymentProviders() {
	// Inicializar el generador de números aleatorios
	rand.Seed(time.Now().UnixNano())

	paypalClientID := getEnv("PAYPAL_CLIENT_ID", "ASYN839bjb4gjMr6nRCc-7YYR8HutdM48kFMWhq-Sxp-PgB5c5R38yGiLBEPwDBIptFj8IJ71OPVXVUt")
	paypalSecret := getEnv("PAYPAL_SECRET", "EHGs6eflLFMvOWTrhWjCWmwBmbXIL8he0dM6bIbcVDFhwGStuz3PGFp_nreODGiJueoNyjxfZG1Hqi0-")
	paypalEnv := getEnv("PAYPAL_ENV", "sandbox")

	var err error
	if paypalEnv == "live" {
		paypalClient, err = paypal.NewClient(paypalClientID, paypalSecret, paypal.APIBaseLive)
	} else {
		paypalClient, err = paypal.NewClient(paypalClientID, paypalSecret, paypal.APIBaseSandBox)
	}
	if err != nil {
		log.Printf("Advertencia: Error al inicializar cliente PayPal: %v", err)
	} else {
		log.Printf("Cliente PayPal inicializado correctamente. Modo: %s", paypalEnv)
	}

	coinbaseAPIKey = getEnv("COINBASE_COMMERCE_API_KEY", "")
	if coinbaseAPIKey == "" {
		log.Println("Advertencia: COINBASE_COMMERCE_API_KEY no está configurada")
	}
}

func crearPago(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	// Extraer información del usuario de la autenticación
	userValue, exists := c.Get("user")
	if !exists {
		log.Println("Error: Usuario no encontrado en el contexto")
		SendErrorResponse(c, ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	user, ok := userValue.(Usuario)
	if !ok {
		log.Println("Error: No se pudo convertir el valor del usuario al tipo Usuario")
		SendErrorResponse(c, errors.New("error al obtener información del usuario"), http.StatusInternalServerError)
		return
	}

	// Leer y validar la solicitud de pago
	var req PagoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error al parsear solicitud de pago: %v", err)
		SendValidationErrorResponse(c, ErrInvalidRequest, gin.H{"details": err.Error()})
		return
	}

	log.Printf("Solicitud de pago recibida: %+v", req)

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
		SendErrorResponse(c, errors.New("método de pago no válido"), http.StatusBadRequest)
		return
	}

	// Validar detalles específicos según el método
	if req.Metodo == "tarjeta" && req.DetallesTarjeta == nil {
		SendErrorResponse(c, errors.New("se requieren detalles de tarjeta para este método de pago"), http.StatusBadRequest)
		return
	}

	// Establecer moneda predeterminada si no se proporciona
	if req.Moneda == "" {
		req.Moneda = "USD"
	}

	// Buscar el curso en la base de datos
	var curso Curso
	if result := db.First(&curso, req.CursoID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("Curso no encontrado para ID: %d", req.CursoID)
			SendErrorResponse(c, ErrResourceNotFound, http.StatusNotFound)
		} else {
			log.Printf("Error de base de datos al buscar curso: %v", result.Error)
			SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		}
		return
	}

	// Verificar si ya existe un pago aprobado para este curso y usuario
	var pagoExistente Pago
	if result := db.Where("usuario_id = ? AND curso_id = ? AND estado = ?",
		user.ID, req.CursoID, "aprobado").First(&pagoExistente); result.Error == nil {

		log.Printf("Usuario %d ya tiene acceso aprobado al curso %d", user.ID, req.CursoID)
		SendSuccessResponse(c, gin.H{
			"message": "Ya tienes acceso a este curso",
			"estado":  "aprobado",
			"pago_id": pagoExistente.ID,
		})
		return
	}

	// Crear el nuevo registro de pago
	pago := Pago{
		UsuarioID:     user.ID,
		CursoID:       req.CursoID,
		Monto:         req.Monto,
		Metodo:        req.Metodo,
		Estado:        "pendiente",
		TransaccionID: "",
		Moneda:        req.Moneda,
	}

	// Guardar el pago en la base de datos
	if result := db.Create(&pago); result.Error != nil {
		log.Printf("Error al guardar pago en base de datos: %v", result.Error)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}

	// Procesar según el método de pago seleccionado
	switch req.Metodo {
	case "dev":
		// Modo desarrollo: simular proceso de pago
		go simularPasarelaPago(pago.ID, req.Metodo)
		SendSuccessResponse(c, gin.H{
			"message": "Pago en proceso (modo desarrollo)",
			"pago_id": pago.ID,
			"estado":  pago.Estado,
		})

	case "paypal":
		// Procesar pago con PayPal
		paypalOrder, err := crearOrdenPayPalSimple(pago)
		if err != nil {
			log.Printf("Error al crear orden PayPal: %v", err)
			SendErrorResponse(c, errors.New("error al procesar pago con PayPal"), http.StatusInternalServerError)
			return
		}

		// Obtener URL de aprobación de PayPal
		paypalApprovalURL := getPayPalApprovalURL(paypalOrder)
		if paypalApprovalURL == "" {
			log.Printf("Error: No se pudo obtener URL de aprobación de PayPal")
			SendErrorResponse(c, errors.New("error al obtener URL de PayPal"), http.StatusInternalServerError)
			return
		}

		log.Printf("URL de redirección PayPal: %s", paypalApprovalURL)

		// Actualizar ID de transacción
		pago.TransaccionID = paypalOrder.ID
		db.Save(&pago)

		// Responder con URL de redirección
		SendSuccessResponse(c, gin.H{
			"message":      "Redirigir a PayPal para completar el pago",
			"pago_id":      pago.ID,
			"estado":       pago.Estado,
			"checkout_url": paypalApprovalURL,
		})

	case "coinbase":
		// Procesar pago con Coinbase
		charge, err := crearCargoCoinbase(pago, curso)
		if err != nil {
			log.Printf("Error al crear cargo Coinbase: %v", err)
			SendErrorResponse(c, errors.New("error al procesar pago con Coinbase"), http.StatusInternalServerError)
			return
		}

		pago.TransaccionID = charge.Code
		db.Save(&pago)

		SendSuccessResponse(c, gin.H{
			"message":      "Redirigir a Coinbase para completar el pago",
			"pago_id":      pago.ID,
			"estado":       pago.Estado,
			"checkout_url": charge.HostedURL,
		})

	case "mercadopago":
		// Simulación de integración con Mercado Pago
		mockCheckoutURL := fmt.Sprintf("https://www.mercadopago.com.%s/checkout?pago_id=%d",
			getEnv("MERCADOPAGO_COUNTRY", "mx"), pago.ID)

		pago.TransaccionID = fmt.Sprintf("mp_%d_%d", time.Now().Unix(), rand.Intn(10000))
		db.Save(&pago)

		SendSuccessResponse(c, gin.H{
			"message":      "Redirigir a Mercado Pago para completar el pago",
			"pago_id":      pago.ID,
			"estado":       pago.Estado,
			"checkout_url": mockCheckoutURL,
		})

	case "stripe":
		// Simulación de integración con Stripe
		mockCheckoutURL := fmt.Sprintf("https://checkout.stripe.com/pay/%s",
			fmt.Sprintf("cs_%d_%d", time.Now().Unix(), rand.Intn(10000)))

		pago.TransaccionID = fmt.Sprintf("ch_%d_%d", time.Now().Unix(), rand.Intn(10000))
		db.Save(&pago)

		SendSuccessResponse(c, gin.H{
			"message":      "Redirigir a Stripe para completar el pago",
			"pago_id":      pago.ID,
			"estado":       pago.Estado,
			"checkout_url": mockCheckoutURL,
		})

	}
}

// Función mejorada para crear una orden de PayPal
func crearOrdenPayPalSimple(pago Pago) (*paypal.Order, error) {
	// Crear un contexto para la petición
	ctx := context.Background()

	// Verificar que el cliente PayPal esté inicializado
	if paypalClient == nil {
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
	baseURL := getEnv("BASE_URL", "")
	if baseURL == "" {
		// Si BASE_URL no está configurada, intentar usar una URL basada en la IP del servidor
		baseURL = fmt.Sprintf("http://%s:5000", getEnv("SERVER_IP", "localhost"))
		log.Printf("BASE_URL no configurada, usando: %s", baseURL)
	}

	frontendURL := getEnv("FRONTEND_URL", "")
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
	order, err := paypalClient.CreateOrder(ctx, "CAPTURE", []paypal.PurchaseUnitRequest{purchaseUnit}, nil, appContext)

	if err != nil {
		return nil, fmt.Errorf("error al crear orden de PayPal: %v", err)
	}

	log.Printf("Orden PayPal creada. ID: %s, Estado: %s", order.ID, order.Status)
	return order, nil
}

// Función mejorada para extraer URL de aprobación de PayPal
func getPayPalApprovalURL(order *paypal.Order) string {
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

func crearCargoCoinbase(pago Pago, curso Curso) (*CoinbaseCharge, error) {
	charge := CoinbaseCharge{
		Name:        fmt.Sprintf("Curso: %s", curso.Titulo),
		Description: fmt.Sprintf("Acceso al curso %s", curso.Titulo),
		PricingType: "fixed_price",
		LocalPrice: CoinbasePrice{
			Amount:   fmt.Sprintf("%.2f", pago.Monto),
			Currency: pago.Moneda,
		},
		Metadata: CoinbaseMetadata{
			PagoID:    pago.ID,
			CursoID:   pago.CursoID,
			UsuarioID: pago.UsuarioID,
		},
		RedirectURL: fmt.Sprintf("%s/pagos/completado", getEnv("FRONTEND_URL", "")),
		CancelURL:   fmt.Sprintf("%s/pagos/cancelado", getEnv("FRONTEND_URL", "")),
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
	req.Header.Set("X-CC-Api-Key", coinbaseAPIKey)
	req.Header.Set("X-CC-Version", "2018-03-22")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("coinbase API returned status %d", resp.StatusCode)
	}

	var response CoinbaseChargeResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

// FUNCIÓN MEJORADA: Verificar estado de pago sin requerir token en la URL
func verificarPagoPorCurso(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	// Extraer información del usuario autenticado
	userValue, exists := c.Get("user")
	if !exists {
		log.Println("Error en verificarPagoPorCurso: Usuario no encontrado en el contexto")
		SendErrorResponse(c, ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	user, ok := userValue.(Usuario)
	if !ok {
		log.Println("Error en verificarPagoPorCurso: No se pudo convertir el valor del usuario")
		SendErrorResponse(c, errors.New("error al obtener información del usuario"), http.StatusInternalServerError)
		return
	}

	// Obtener ID del curso de los parámetros
	cursoID := c.Param("id")
	if cursoID == "" {
		SendErrorResponse(c, errors.New("se requiere el ID del curso"), http.StatusBadRequest)
		return
	}

	var cursoUint uint
	if _, err := fmt.Sscanf(cursoID, "%d", &cursoUint); err != nil {
		SendErrorResponse(c, errors.New("ID de curso inválido"), http.StatusBadRequest)
		return
	}

	log.Printf("Verificando pago para curso ID: %s, usuario ID: %d", cursoID, user.ID)

	// Buscar el pago más reciente para este curso y usuario
	var pago Pago
	result := db.Where("usuario_id = ? AND curso_id = ?", user.ID, cursoUint).
		Order("created_at desc").
		First(&pago)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// No hay pago registrado
			c.JSON(http.StatusOK, gin.H{
				"estado":  "no_pagado",
				"message": "No se encontró pago para este curso",
			})
			return
		}
		log.Printf("Error de base de datos al verificar pago: %v", result.Error)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}

	// Para métodos de pago externos, verificar estado actual si está pendiente
	if pago.Estado == "pendiente" {
		switch pago.Metodo {
		case "paypal":
			if pago.TransaccionID != "" {
				// Verificar estado actual con PayPal
				ctx := context.Background()
				orderDetail, err := paypalClient.GetOrder(ctx, pago.TransaccionID)
				if err == nil && orderDetail != nil {
					log.Printf("Estado actual de PayPal para orden %s: %s", pago.TransaccionID, orderDetail.Status)
					if orderDetail.Status == "COMPLETED" ||
						orderDetail.Status == "APPROVED" ||
						orderDetail.Status == "PAYER_ACTION_REQUIRED" {
						pago.Estado = "aprobado"
						db.Save(&pago)
						log.Printf("Actualizado estado de pago ID %d a 'aprobado' según PayPal", pago.ID)
					}
				} else if err != nil {
					log.Printf("Error al verificar estado con PayPal: %v", err)
				}
			}

		}
	}

	// Responder con el estado del pago
	c.JSON(http.StatusOK, gin.H{
		"estado": pago.Estado,
		"pago":   pago,
	})
}

// FUNCIÓN MODIFICADA: Callback de PayPal redirige a la página del curso
func callbackPayPal(c *gin.Context) {
	// Extraer parámetros de la URL
	pagoID := c.Query("pago_id")
	cursoID := c.Query("curso_id")
	token := c.Query("token")

	log.Printf("Recibida callback PayPal con pagoID: %s, cursoID: %s, token: %s", pagoID, cursoID, token)

	if pagoID == "" || token == "" {
		log.Printf("Error en callback PayPal: parámetros inválidos")
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Parámetros inválidos. Por favor intenta nuevamente.",
		})
		return
	}

	var pagoIDUint uint
	if _, err := fmt.Sscanf(pagoID, "%d", &pagoIDUint); err != nil {
		log.Printf("Error en callback PayPal: ID de pago inválido - %v", err)
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "ID de pago inválido. Por favor intenta nuevamente.",
		})
		return
	}

	var cursoIDUint uint
	if cursoID != "" {
		if _, err := fmt.Sscanf(cursoID, "%d", &cursoIDUint); err != nil {
			log.Printf("Error en callback PayPal: ID de curso inválido - %v", err)
			// Continuamos porque el ID del curso no es crítico, podemos obtenerlo del pago
		}
	}

	// Buscar el pago en la base de datos
	var pago Pago
	if result := db.First(&pago, pagoIDUint); result.Error != nil {
		log.Printf("Error en callback PayPal: Pago no encontrado (ID: %d) - %v", pagoIDUint, result.Error)
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "Pago no encontrado. Por favor contacta a soporte.",
		})
		return
	}

	// Usar el ID del curso del pago si no se proporcionó en la URL
	if cursoIDUint == 0 {
		cursoIDUint = pago.CursoID
	}

	// Crear un contexto para la petición a PayPal
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Printf("Intentando capturar orden PayPal con token: %s para pago ID: %d", token, pago.ID)

	// Capturar la orden de PayPal con el token
	captureResult, err := paypalClient.CaptureOrder(ctx, token, paypal.CaptureOrderRequest{})
	if err != nil {
		log.Printf("Error al capturar orden PayPal: %v", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Error al procesar pago con PayPal. Por favor intenta nuevamente o contacta a soporte.",
		})
		return
	}

	log.Printf("Orden PayPal capturada. Estado: %s", captureResult.Status)

	// Actualizar el estado del pago según la respuesta de PayPal
	estadoActualizado := false
	if captureResult.Status == "COMPLETED" || captureResult.Status == "APPROVED" {
		pago.Estado = "aprobado"
		estadoActualizado = true
	} else if captureResult.Status == "DECLINED" || captureResult.Status == "FAILED" {
		pago.Estado = "rechazado"
		estadoActualizado = true
	}

	// Guardar cambios en la base de datos si se actualizó el estado
	if estadoActualizado {
		if result := db.Save(&pago); result.Error != nil {
			log.Printf("Error al actualizar estado de pago ID %d: %v", pago.ID, result.Error)
			// Continuar a pesar del error, para no bloquear al usuario
		} else {
			log.Printf("Pago ID %d actualizado a estado '%s'", pago.ID, pago.Estado)
		}
	}

	// Redirigir al usuario a la página del curso si el pago fue aprobado
	frontendURL := getEnv("FRONTEND_URL", "http://localhost:3000")
	var redirectURL string

	if pago.Estado == "aprobado" {
		// Redirigir a la página del curso en lugar de la página de pago completado
		redirectURL = fmt.Sprintf("%s/curso/%d", frontendURL, cursoIDUint)
		log.Printf("Pago aprobado, redirigiendo a la página del curso: %s", redirectURL)
	} else {
		// En caso de fallo, redirigir a la página de pago fallido
		redirectURL = fmt.Sprintf("%s/pagos/fallido?pago_id=%d", frontendURL, pagoIDUint)
		log.Printf("Pago no aprobado, redirigiendo a: %s", redirectURL)
	}

	log.Printf("Redirigiendo a: %s", redirectURL)
	c.Redirect(http.StatusFound, redirectURL)
}

// Función utilizada en modo desarrollo o para tarjeta y transferencia
func simularPasarelaPago(pagoID uint, metodo string) {
	if getEnv("APP_ENV", "") != "development" && metodo != "tarjeta" && metodo != "transferencia" {
		return
	}

	// Simulamos un retardo para el procesamiento
	time.Sleep(3 * time.Second)

	var pago Pago
	if result := db.First(&pago, pagoID); result.Error != nil {
		log.Printf("Error al recuperar pago ID %d para simulación: %v", pagoID, result.Error)
		return
	}

	// En modo desarrollo, siempre aprobamos el pago
	if getEnv("APP_ENV", "") == "development" {
		pago.Estado = "aprobado"
		pago.TransaccionID = generarIDTransaccion(metodo)
	} else {
		// En producción, simulamos una tasa de aprobación del 80%
		if rand.Intn(100) < 80 {
			pago.Estado = "aprobado"
			pago.TransaccionID = generarIDTransaccion(metodo)
		} else {
			pago.Estado = "rechazado"
		}
	}

	if result := db.Save(&pago); result.Error != nil {
		log.Printf("Error al actualizar estado de pago ID %d: %v", pagoID, result.Error)
	} else {
		log.Printf("Pago ID %d actualizado a estado: %s", pagoID, pago.Estado)
	}
}

// Generar ID de transacción único para cada método de pago
func generarIDTransaccion(metodo string) string {
	timestamp := time.Now().Unix()
	randomPart := rand.Intn(100000)

	// Prefijo según el método de pago
	prefix := "txn"
	switch metodo {
	case "tarjeta":
		prefix = "card"
	case "paypal":
		prefix = "pp"
	case "coinbase":
		prefix = "cb"
	case "transferencia":
		prefix = "trf"
	case "dev":
		prefix = "dev"
	case "stripe":
		prefix = "ch"
	case "mercadopago":
		prefix = "mp"
	}

	return fmt.Sprintf("%s_%d_%d", prefix, timestamp, randomPart)
}

// Webhook genérico para pasarelas de pago
func webhookPago(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	// No requerimos autenticación para webhooks ya que vienen de la pasarela de pago

	var payload struct {
		PagoID        uint   `json:"pago_id"`
		Estado        string `json:"estado"`
		TransaccionID string `json:"transaccion_id"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Printf("Error en webhookPago: JSON inválido - %v", err)
		SendErrorResponse(c, ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	log.Printf("Webhook recibido: Pago ID %d, Estado %s", payload.PagoID, payload.Estado)

	var pago Pago
	if result := db.First(&pago, payload.PagoID); result.Error != nil {
		log.Printf("Error en webhookPago: Pago ID %d no encontrado - %v", payload.PagoID, result.Error)
		SendErrorResponse(c, ErrPaymentNotFound, http.StatusNotFound)
		return
	}

	// Actualizar estado y posiblemente el ID de transacción
	pago.Estado = payload.Estado
	if payload.TransaccionID != "" {
		pago.TransaccionID = payload.TransaccionID
	}

	if result := db.Save(&pago); result.Error != nil {
		log.Printf("Error en webhookPago: No se pudo actualizar Pago ID %d - %v", payload.PagoID, result.Error)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}

	// Registrar la actualización en logs
	log.Printf("Pago ID %d actualizado a estado %s mediante webhook", pago.ID, pago.Estado)

	SendSuccessResponse(c, gin.H{
		"message": "Estado de pago actualizado correctamente",
		"pago_id": pago.ID,
		"estado":  pago.Estado,
	})
}

// Webhook específico para PayPal
func webhookPayPal(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	// Verificar encabezados de webhook
	paypalEvent := c.GetHeader("Paypal-Transmission-Id")
	if paypalEvent == "" {
		log.Println("Advertencia: Posible llamada no autorizada a webhook de PayPal")
		// Continuar procesando ya que algunos eventos de prueba pueden no incluir este encabezado
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
		SendErrorResponse(c, errors.New("error al leer cuerpo de la solicitud"), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("Error al parsear evento de PayPal: %v", err)
		SendErrorResponse(c, ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	log.Printf("Webhook PayPal recibido: Tipo %s, ID %s, Estado %s",
		event.EventType, event.Resource.ID, event.Resource.Status)

	// Procesar solo eventos de captura completada
	if event.EventType != "PAYMENT.CAPTURE.COMPLETED" &&
		event.EventType != "CHECKOUT.ORDER.APPROVED" {
		log.Printf("Evento PayPal no manejado: %s", event.EventType)
		c.JSON(http.StatusOK, gin.H{"message": "Evento no manejado"})
		return
	}

	// Buscar el pago asociado a esta transacción
	var pago Pago
	if result := db.Where("transaccion_id = ?", event.Resource.ID).First(&pago); result.Error != nil {
		log.Printf("Error en webhook PayPal: No se encontró pago con transacción ID %s", event.Resource.ID)
		SendErrorResponse(c, ErrPaymentNotFound, http.StatusNotFound)
		return
	}

	// Actualizar estado según el evento
	estadoAnterior := pago.Estado
	pago.Estado = "aprobado"

	if result := db.Save(&pago); result.Error != nil {
		log.Printf("Error al actualizar estado de pago ID %d: %v", pago.ID, result.Error)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}

	log.Printf("Pago ID %d actualizado de '%s' a 'aprobado' mediante webhook PayPal",
		pago.ID, estadoAnterior)

	c.JSON(http.StatusOK, gin.H{
		"message": "Webhook de PayPal procesado correctamente",
		"pago_id": pago.ID,
		"estado":  pago.Estado,
	})
}

// Webhook específico para Coinbase
func webhookCoinbase(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	// Verificar firma de webhook para autenticación
	signature := c.GetHeader("X-CC-Webhook-Signature")
	if signature == "" && getEnv("ENV", "development") != "development" {
		log.Println("Error en webhook Coinbase: Firma no proporcionada")
		SendErrorResponse(c, errors.New("firma no proporcionada"), http.StatusUnauthorized)
		return
	}

	// Leer y validar el evento
	body, err := c.GetRawData()
	if err != nil {
		log.Printf("Error al leer cuerpo de webhook Coinbase: %v", err)
		SendErrorResponse(c, errors.New("error al leer cuerpo del webhook"), http.StatusBadRequest)
		return
	}

	var event CoinbaseWebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("Error al parsear evento Coinbase: %v", err)
		SendErrorResponse(c, errors.New("error al parsear evento"), http.StatusBadRequest)
		return
	}

	log.Printf("Webhook Coinbase recibido: Tipo %s, Código %s",
		event.Event.Type, event.Event.Data.Code)

	// Procesar solo eventos relevantes
	if event.Event.Type != "charge:confirmed" && event.Event.Type != "charge:failed" {
		log.Printf("Evento Coinbase no manejado: %s", event.Event.Type)
		c.JSON(http.StatusOK, gin.H{"message": "Evento no manejado"})
		return
	}

	// Extraer metadata del pago
	metadata := event.Event.Data.Metadata
	var pago Pago

	// Buscar primero por ID de pago en metadata
	if metadata.PagoID > 0 {
		if result := db.First(&pago, metadata.PagoID); result.Error != nil {
			log.Printf("Error en webhook Coinbase: Pago ID %d no encontrado", metadata.PagoID)

			// Intentar buscar por código de transacción como alternativa
			if result := db.Where("transaccion_id = ?", event.Event.Data.Code).First(&pago); result.Error != nil {
				log.Printf("Error en webhook Coinbase: No se encontró pago con transacción ID %s", event.Event.Data.Code)
				SendErrorResponse(c, ErrPaymentNotFound, http.StatusNotFound)
				return
			}
		}
	} else {
		// Buscar por código de transacción si no hay ID de pago
		if result := db.Where("transaccion_id = ?", event.Event.Data.Code).First(&pago); result.Error != nil {
			log.Printf("Error en webhook Coinbase: No se encontró pago con transacción ID %s", event.Event.Data.Code)
			SendErrorResponse(c, ErrPaymentNotFound, http.StatusNotFound)
			return
		}
	}

	// Actualizar estado según el evento
	estadoAnterior := pago.Estado
	if event.Event.Type == "charge:confirmed" {
		pago.Estado = "aprobado"
	} else {
		pago.Estado = "rechazado"
	}

	// Asegurar que tenemos el código de transacción
	if pago.TransaccionID == "" {
		pago.TransaccionID = event.Event.Data.Code
	}

	if result := db.Save(&pago); result.Error != nil {
		log.Printf("Error al actualizar estado de pago ID %d: %v", pago.ID, result.Error)
		SendErrorResponse(c, ErrDatabaseError, http.StatusInternalServerError)
		return
	}

	log.Printf("Pago ID %d actualizado de '%s' a '%s' mediante webhook Coinbase",
		pago.ID, estadoAnterior, pago.Estado)

	c.JSON(http.StatusOK, gin.H{
		"message": "Webhook de Coinbase procesado correctamente",
		"pago_id": pago.ID,
		"estado":  pago.Estado,
	})
}
