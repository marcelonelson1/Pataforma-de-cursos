// services/coinbase_service.go
package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"payment-service/config"
	"payment-service/models"
)

// CoinbaseCharge estructuras para Coinbase (migradas del código original)
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

type CoinbaseErrorResponse struct {
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
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

// CreateCoinbaseCharge crea un cargo en Coinbase (migrado de crearCargoCoinbase)
func CreateCoinbaseCharge(payment *models.Payment, course *CourseInfo, cfg *config.Config) (*CoinbaseCharge, error) {
	log.Printf("🔧 [COINBASE] Iniciando creación de cargo para pago ID: %d", payment.ID)
	
	// 🔥 VALIDACIÓN: Verificar API Key
	if cfg.Coinbase.APIKey == "" {
		log.Printf("❌ [COINBASE] CRÍTICO: Coinbase API key no configurada")
		log.Printf("💡 [COINBASE] Configura COINBASE_COMMERCE_API_KEY en tu .env")
		log.Printf("💡 [COINBASE] Obtén tu API Key en: https://commerce.coinbase.com/settings/api-keys")
		return nil, fmt.Errorf("Coinbase API key no configurada")
	}
	
	log.Printf("✅ [COINBASE] API Key configurada: %s... (length: %d)", 
		maskCoinbaseKey(cfg.Coinbase.APIKey), len(cfg.Coinbase.APIKey))

	// Validar parámetros de entrada
	if payment.Monto <= 0 {
		log.Printf("❌ [COINBASE] Monto inválido: %.2f", payment.Monto)
		return nil, fmt.Errorf("monto debe ser mayor a 0")
	}

	if payment.Moneda == "" {
		payment.Moneda = "USD"
		log.Printf("⚠️ [COINBASE] Moneda no especificada, usando USD por defecto")
	}

	// Crear estructura del cargo
	charge := CoinbaseCharge{
		Name:        fmt.Sprintf("Curso: %s", course.Titulo),
		Description: fmt.Sprintf("Acceso al curso %s", course.Titulo),
		PricingType: "fixed_price",
		LocalPrice: CoinbasePrice{
			Amount:   fmt.Sprintf("%.2f", payment.Monto),
			Currency: payment.Moneda,
		},
		Metadata: CoinbaseMetadata{
			PagoID:    payment.ID,
			CursoID:   payment.CursoID,
			UsuarioID: payment.UsuarioID,
		},
		// URLs usando el controlador genérico para consistencia
		RedirectURL: fmt.Sprintf("%s/api/pagos/return?pago_id=%d&provider=coinbase&status=success", cfg.BaseURL, payment.ID),
		CancelURL:   fmt.Sprintf("%s/api/pagos/return?pago_id=%d&provider=coinbase&status=returned", cfg.BaseURL, payment.ID),
	}

	log.Printf("✅ [COINBASE] Cargo configurado:")
	log.Printf("   Nombre: %s", charge.Name)
	log.Printf("   Descripción: %s", charge.Description)
	log.Printf("   Monto: %s %s", charge.LocalPrice.Amount, charge.LocalPrice.Currency)
	log.Printf("   Redirect URL: %s", charge.RedirectURL)
	log.Printf("   Cancel URL: %s", charge.CancelURL)
	log.Printf("   Metadata: PagoID=%d, CursoID=%d, UsuarioID=%d", 
		charge.Metadata.PagoID, charge.Metadata.CursoID, charge.Metadata.UsuarioID)

	// Serializar payload
	payload, err := json.Marshal(charge)
	if err != nil {
		log.Printf("❌ [COINBASE] Error al serializar datos: %v", err)
		return nil, fmt.Errorf("error al serializar datos: %v", err)
	}

	log.Printf("📦 [COINBASE] Payload JSON: %s", string(payload))

	// Crear cliente HTTP
	client := &http.Client{
		Timeout: 30 * time.Second, // Aumentar timeout
	}

	// Crear request
	log.Printf("📤 [COINBASE] Enviando solicitud a Coinbase Commerce API...")
	req, err := http.NewRequest("POST", "https://api.commerce.coinbase.com/charges", bytes.NewReader(payload))
	if err != nil {
		log.Printf("❌ [COINBASE] Error al crear request: %v", err)
		return nil, fmt.Errorf("error al crear request: %v", err)
	}

	// Configurar headers requeridos por Coinbase
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CC-Api-Key", cfg.Coinbase.APIKey)
	req.Header.Set("X-CC-Version", "2018-03-22")
	req.Header.Set("User-Agent", "Payment-Service/1.0")

	log.Printf("🔍 [COINBASE] Headers configurados:")
	log.Printf("   Content-Type: %s", req.Header.Get("Content-Type"))
	log.Printf("   X-CC-Api-Key: %s...", maskCoinbaseKey(cfg.Coinbase.APIKey))
	log.Printf("   X-CC-Version: %s", req.Header.Get("X-CC-Version"))
	log.Printf("   User-Agent: %s", req.Header.Get("User-Agent"))

	// Ejecutar request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ [COINBASE] Error al llamar API de Coinbase: %v", err)
		return nil, fmt.Errorf("error al llamar API de Coinbase: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("📥 [COINBASE] Respuesta recibida:")
	log.Printf("   Status Code: %d", resp.StatusCode)
	log.Printf("   Status: %s", resp.Status)

	// Leer el cuerpo de la respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("❌ [COINBASE] Error al leer respuesta: %v", err)
		return nil, fmt.Errorf("error al leer respuesta: %v", err)
	}

	log.Printf("📄 [COINBASE] Response Body: %s", string(body))

	// Verificar status code
	if resp.StatusCode != http.StatusCreated {
		log.Printf("❌ [COINBASE] Error detallado:")
		log.Printf("   Status Code: %d", resp.StatusCode)
		log.Printf("   Response Body: %s", string(body))
		
		// Intentar decodificar error específico de Coinbase
		var errorResp CoinbaseErrorResponse
		if err := json.Unmarshal(body, &errorResp); err == nil {
			log.Printf("   Error Type: %s", errorResp.Error.Type)
			log.Printf("   Error Message: %s", errorResp.Error.Message)
		}
		
		// Diagnósticos específicos por status code
		switch resp.StatusCode {
		case 401:
			log.Printf("💡 [COINBASE] Error 401: API Key inválida o expirada")
			log.Printf("💡 [COINBASE] Solución: Verifica tu API Key en Coinbase Commerce Dashboard")
		case 412:
			log.Printf("💡 [COINBASE] Error 412: Precondición falló")
			log.Printf("💡 [COINBASE] Posibles causas:")
			log.Printf("   - Moneda no soportada (%s)", payment.Moneda)
			log.Printf("   - Monto fuera del rango permitido (%.2f)", payment.Monto)
			log.Printf("   - Configuración de negocio incompleta en Coinbase")
			log.Printf("   - Restricciones geográficas (Argentina)")
		case 400:
			log.Printf("💡 [COINBASE] Error 400: Bad Request")
			log.Printf("💡 [COINBASE] Revisa el formato del payload JSON")
		case 403:
			log.Printf("💡 [COINBASE] Error 403: Forbidden")
			log.Printf("💡 [COINBASE] Permisos insuficientes en la API Key")
		case 422:
			log.Printf("💡 [COINBASE] Error 422: Unprocessable Entity")
			log.Printf("💡 [COINBASE] Datos inválidos en la solicitud")
		default:
			log.Printf("💡 [COINBASE] Error %d: Error no documentado", resp.StatusCode)
		}
		
		return nil, fmt.Errorf("coinbase API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Decodificar respuesta exitosa
	var response CoinbaseChargeResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("❌ [COINBASE] Error al decodificar respuesta JSON: %v", err)
		log.Printf("📄 [COINBASE] Raw response: %s", string(body))
		return nil, fmt.Errorf("error al decodificar respuesta: %v", err)
	}

	log.Printf("✅ [COINBASE] Cargo creado exitosamente:")
	log.Printf("   Code: %s", response.Data.Code)
	log.Printf("   Name: %s", response.Data.Name)
	log.Printf("   Hosted URL: %s", response.Data.HostedURL)
	log.Printf("   Pricing Type: %s", response.Data.PricingType)

	return &response.Data, nil
}

// 🔥 FUNCIÓN: Enmascarar API Key de Coinbase
func maskCoinbaseKey(apiKey string) string {
	if len(apiKey) < 8 {
		return "****"
	}
	if len(apiKey) < 16 {
		return apiKey[:4] + "****"
	}
	return apiKey[:8] + "****" + apiKey[len(apiKey)-8:]
}

// 🔥 FUNCIÓN: Validar configuración de Coinbase
func ValidateCoinbaseConfig(cfg *config.Config) error {
	if cfg.Coinbase.APIKey == "" {
		return fmt.Errorf("COINBASE_COMMERCE_API_KEY no configurada")
	}
	
	// Validar formato básico de API Key (Coinbase usa UUIDs)
	if len(cfg.Coinbase.APIKey) < 32 {
		return fmt.Errorf("API Key de Coinbase parece incorrecta (muy corta)")
	}
	
	log.Printf("✅ [COINBASE] Configuración válida")
	return nil
}

// 🔥 FUNCIÓN: Test de conectividad con Coinbase
func TestCoinbaseConnection(cfg *config.Config) error {
	if err := ValidateCoinbaseConfig(cfg); err != nil {
		return err
	}
	
	log.Printf("🔍 [COINBASE] Probando conectividad con Coinbase Commerce...")
	
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	// Endpoint de health check o información básica
	req, err := http.NewRequest("GET", "https://api.commerce.coinbase.com/charges", nil)
	if err != nil {
		return fmt.Errorf("error creando request de test: %v", err)
	}
	
	req.Header.Set("X-CC-Api-Key", cfg.Coinbase.APIKey)
	req.Header.Set("X-CC-Version", "2018-03-22")
	
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error de conectividad: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == 401 {
		return fmt.Errorf("API Key inválida")
	}
	
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("✅ [COINBASE] Conectividad exitosa")
		return nil
	}
	
	return fmt.Errorf("respuesta inesperada: %d", resp.StatusCode)
}