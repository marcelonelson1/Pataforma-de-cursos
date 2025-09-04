// services/paypal_service.go
package services

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/plutov/paypal/v4"

	"payment-service/config"
	"payment-service/models"
)

var paypalClient *paypal.Client

// InitPaymentProviders inicializa los proveedores de pago (migrado de initPaymentProviders)
func InitPaymentProviders(cfg *config.Config) {
	// 🔥 NUEVO: Verificar credenciales antes de inicializar
	log.Printf("🔧 [PAYPAL_INIT] Verificando credenciales PayPal:")
	log.Printf("   Client ID: %s... (length: %d)", maskCredential(cfg.PayPal.ClientID), len(cfg.PayPal.ClientID))
	log.Printf("   Secret: %s... (length: %d)", maskCredential(cfg.PayPal.Secret), len(cfg.PayPal.Secret))
	log.Printf("   Environment: %s", cfg.PayPal.Env)

	if cfg.PayPal.ClientID == "" {
		log.Printf("❌ [PAYPAL_INIT] CRÍTICO: Client ID está vacío!")
		return
	}
	if cfg.PayPal.Secret == "" {
		log.Printf("❌ [PAYPAL_INIT] CRÍTICO: Secret está vacío!")
		return
	}

	var err error
	if cfg.PayPal.Env == "live" {
		paypalClient, err = paypal.NewClient(cfg.PayPal.ClientID, cfg.PayPal.Secret, paypal.APIBaseLive)
	} else {
		paypalClient, err = paypal.NewClient(cfg.PayPal.ClientID, cfg.PayPal.Secret, paypal.APIBaseSandBox)
	}
	
	if err != nil {
		log.Printf("❌ [PAYPAL_INIT] Error al inicializar cliente PayPal: %v", err)
	} else {
		log.Printf("✅ [PAYPAL_INIT] Cliente PayPal inicializado correctamente. Modo: %s", cfg.PayPal.Env)
		
		// 🔥 NUEVO: Test de autenticación obligatorio
		if err := testPayPalAuth(); err != nil {
			log.Printf("❌ [PAYPAL_INIT] ⚠️  FALLO DE AUTENTICACIÓN: %v", err)
			log.Printf("💡 [PAYPAL_INIT] Verifica tus credenciales en PayPal Developer Dashboard")
		} else {
			log.Printf("🎉 [PAYPAL_INIT] ✅ AUTENTICACIÓN EXITOSA - PayPal listo para usar")
		}
	}
}

// 🔥 NUEVA FUNCIÓN: Verificar autenticación con PayPal
func testPayPalAuth() error {
	if paypalClient == nil {
		return fmt.Errorf("cliente PayPal no inicializado")
	}
	
	log.Printf("🔍 [PAYPAL_AUTH] Probando autenticación...")
	ctx := context.Background()
	token, err := paypalClient.GetAccessToken(ctx)
	if err != nil {
		return fmt.Errorf("error de autenticación: %v", err)
	}
	
	// Verificar que el token no esté vacío
	if token == nil {
		return fmt.Errorf("token es nil")
	}
	
	log.Printf("✅ [PAYPAL_AUTH] Token obtenido exitosamente")
	return nil
}

// 🔥 NUEVA FUNCIÓN: Enmascarar credenciales
func maskCredential(credential string) string {
	if len(credential) < 8 {
		return "****"
	}
	return credential[:4] + "****" + credential[len(credential)-4:]
}

// CreatePayPalOrder crea una orden de PayPal (migrado de crearOrdenPayPalSimple)
func CreatePayPalOrder(payment *models.Payment, cfg *config.Config) (*paypal.Order, error) {
	log.Printf("🔧 [PAYPAL_ORDER] Iniciando creación de orden para pago ID: %d", payment.ID)
	
	ctx := context.Background()

	if paypalClient == nil {
		log.Printf("❌ [PAYPAL_ORDER] Cliente PayPal no inicializado")
		return nil, fmt.Errorf("cliente PayPal no inicializado")
	}

	// 🔥 NUEVO: Verificar autenticación antes de crear orden
	log.Printf("🔍 [PAYPAL_ORDER] Verificando autenticación antes de crear orden...")
	_, err := paypalClient.GetAccessToken(ctx)
	if err != nil {
		log.Printf("❌ [PAYPAL_ORDER] Error de autenticación: %v", err)
		return nil, fmt.Errorf("error de autenticación PayPal: %v", err)
	}
	log.Printf("✅ [PAYPAL_ORDER] Autenticación verificada")

	// Definir la unidad de compra
	purchaseUnit := paypal.PurchaseUnitRequest{
		ReferenceID: fmt.Sprintf("pago_%d", payment.ID),
		Amount: &paypal.PurchaseUnitAmount{
			Currency: payment.Moneda,
			Value:    fmt.Sprintf("%.2f", payment.Monto),
		},
		Description: fmt.Sprintf("Pago para curso ID: %d", payment.CursoID),
	}

	// URLs de retorno
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = fmt.Sprintf("http://localhost:%s", cfg.Port)
		log.Printf("BASE_URL no configurada, usando: %s", baseURL)
	}

	frontendURL := cfg.FrontendURL
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
		log.Printf("FRONTEND_URL no configurada, usando: %s", frontendURL)
	}

	// URLs para PayPal usando el controlador genérico existente
	returnURL := fmt.Sprintf("%s/api/pagos/paypal/callback?pago_id=%d&curso_id=%d", baseURL, payment.ID, payment.CursoID)
	cancelURL := fmt.Sprintf("%s/api/pagos/return?pago_id=%d&provider=paypal&status=returned", baseURL, payment.ID)
	
	log.Printf("PayPal URLs preparadas - Return: %s | Cancel: %s", returnURL, cancelURL)

	// Contexto de la aplicación
	appContext := &paypal.ApplicationContext{
		ReturnURL:          returnURL,
		CancelURL:          cancelURL,
		UserAction:         "PAY_NOW",
		ShippingPreference: "NO_SHIPPING",
	}

	log.Printf("Creando orden PayPal con ReturnURL: %s, CancelURL: %s",
		appContext.ReturnURL, appContext.CancelURL)

	// Crear la orden
	order, err := paypalClient.CreateOrder(ctx, "CAPTURE", []paypal.PurchaseUnitRequest{purchaseUnit}, nil, appContext)
	if err != nil {
		log.Printf("❌ [PAYPAL_ORDER] Error detallado al crear orden: %v", err)
		return nil, fmt.Errorf("error al crear orden de PayPal: %v", err)
	}

	log.Printf("✅ [PAYPAL_ORDER] Orden PayPal creada. ID: %s, Estado: %s", order.ID, order.Status)
	return order, nil
}

// GetPayPalApprovalURL extrae la URL de aprobación de PayPal (migrado de getPayPalApprovalURL)
func GetPayPalApprovalURL(order *paypal.Order) string {
	if order == nil || len(order.Links) == 0 {
		log.Println("Error: orden de PayPal nula o sin enlaces")
		return ""
	}

	// Buscar enlace de aprobación específico
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

// CapturePayPalOrder captura una orden de PayPal
func CapturePayPalOrder(ctx context.Context, token string) (*paypal.CaptureOrderResponse, error) {
	if paypalClient == nil {
		return nil, fmt.Errorf("cliente PayPal no inicializado")
	}

	captureResult, err := paypalClient.CaptureOrder(ctx, token, paypal.CaptureOrderRequest{})
	if err != nil {
		return nil, fmt.Errorf("error al capturar orden PayPal: %v", err)
	}

	return captureResult, nil
}

// GetPayPalOrderDetails obtiene los detalles de una orden de PayPal
func GetPayPalOrderDetails(orderID string) (*paypal.Order, error) {
	if paypalClient == nil {
		return nil, fmt.Errorf("cliente PayPal no inicializado")
	}

	ctx := context.Background()
	orderDetail, err := paypalClient.GetOrder(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener detalles de orden PayPal: %v", err)
	}

	return orderDetail, nil
}