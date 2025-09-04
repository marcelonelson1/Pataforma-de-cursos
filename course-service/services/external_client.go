// services/external_client.go
package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"course-service/config"
)

// PaymentServiceClient cliente para comunicarse con Payment Service
type PaymentServiceClient struct {
	baseURL string
	client  *http.Client
}

// UserServiceClient cliente para comunicarse con User Service
type UserServiceClient struct {
	baseURL string
	client  *http.Client
}

// PaymentResponse estructura de respuesta del Payment Service
type PaymentResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error,omitempty"`
}

// UserPaymentInfo información de pago del usuario
type UserPaymentInfo struct {
	UserID   uint   `json:"user_id"`
	CourseID uint   `json:"course_id"`
	HasPaid  bool   `json:"has_paid"`
	Status   string `json:"status"`
}

// UserInfo información del usuario
type UserInfo struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// NewPaymentServiceClient crea un nuevo cliente para Payment Service
func NewPaymentServiceClient(cfg *config.Config) *PaymentServiceClient {
	return &PaymentServiceClient{
		baseURL: cfg.PaymentServiceURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// NewUserServiceClient crea un nuevo cliente para User Service
func NewUserServiceClient(cfg *config.Config) *UserServiceClient {
	return &UserServiceClient{
		baseURL: cfg.UserServiceURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// CheckUserPaymentForCourse verifica si un usuario ha pagado por un curso
func CheckUserPaymentForCourse(userID, courseID uint, cfg *config.Config) bool {
	client := NewPaymentServiceClient(cfg)
	return client.CheckUserPaymentForCourse(userID, courseID)
}

// CheckUserPaymentForCourse verifica el pago del usuario para un curso específico
func (p *PaymentServiceClient) CheckUserPaymentForCourse(userID, courseID uint) bool {
	url := fmt.Sprintf("%s/api/pagos/verify-course-access?user_id=%d&course_id=%d", 
		p.baseURL, userID, courseID)

	resp, err := p.client.Get(url)
	if err != nil {
		log.Printf("Error al consultar Payment Service: %v", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Payment Service respondió con status: %d", resp.StatusCode)
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error al leer respuesta de Payment Service: %v", err)
		return false
	}

	var paymentResp PaymentResponse
	if err := json.Unmarshal(body, &paymentResp); err != nil {
		log.Printf("Error al parsear respuesta de Payment Service: %v", err)
		return false
	}

	if !paymentResp.Success {
		log.Printf("Payment Service error: %s", paymentResp.Error)
		return false
	}

	// Extraer información de pago
	if dataMap, ok := paymentResp.Data.(map[string]interface{}); ok {
		if hasPaid, exists := dataMap["has_access"]; exists {
			if paid, ok := hasPaid.(bool); ok {
				return paid
			}
		}
	}

	return false
}

// GetUserPaymentHistory obtiene el historial de pagos del usuario
func (p *PaymentServiceClient) GetUserPaymentHistory(userID uint) ([]UserPaymentInfo, error) {
	url := fmt.Sprintf("%s/api/pagos/user/%d/history", p.baseURL, userID)

	resp, err := p.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error al consultar Payment Service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Payment Service respondió con status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer respuesta: %v", err)
	}

	var paymentResp PaymentResponse
	if err := json.Unmarshal(body, &paymentResp); err != nil {
		return nil, fmt.Errorf("error al parsear respuesta: %v", err)
	}

	if !paymentResp.Success {
		return nil, fmt.Errorf("error del Payment Service: %s", paymentResp.Error)
	}

	// Convertir data a slice de UserPaymentInfo
	var payments []UserPaymentInfo
	if dataBytes, err := json.Marshal(paymentResp.Data); err == nil {
		json.Unmarshal(dataBytes, &payments)
	}

	return payments, nil
}

// CheckPaymentServiceHealth verifica la salud del Payment Service
func CheckPaymentServiceHealth(cfg *config.Config) error {
	client := NewPaymentServiceClient(cfg)
	return client.CheckHealth()
}

// CheckHealth verifica la salud del Payment Service
func (p *PaymentServiceClient) CheckHealth() error {
	url := fmt.Sprintf("%s/api/health", p.baseURL)

	resp, err := p.client.Get(url)
	if err != nil {
		return fmt.Errorf("error al conectar con Payment Service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Payment Service no saludable, status: %d", resp.StatusCode)
	}

	return nil
}

// GetUserInfo obtiene información del usuario desde User Service
func (u *UserServiceClient) GetUserInfo(userID uint) (*UserInfo, error) {
	url := fmt.Sprintf("%s/api/users/%d", u.baseURL, userID)

	resp, err := u.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error al consultar User Service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("usuario no encontrado")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("User Service respondió con status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer respuesta: %v", err)
	}

	var response struct {
		Success bool     `json:"success"`
		Data    UserInfo `json:"data"`
		Error   string   `json:"error,omitempty"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error al parsear respuesta: %v", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("error del User Service: %s", response.Error)
	}

	return &response.Data, nil
}

// CheckUserServiceHealth verifica la salud del User Service
func (u *UserServiceClient) CheckHealth() error {
	url := fmt.Sprintf("%s/api/health", u.baseURL)

	resp, err := u.client.Get(url)
	if err != nil {
		return fmt.Errorf("error al conectar con User Service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("User Service no saludable, status: %d", resp.StatusCode)
	}

	return nil
}

// NotifyPaymentCompleted notifica al Payment Service que un curso fue completado
func (p *PaymentServiceClient) NotifyPaymentCompleted(userID, courseID uint) error {
	url := fmt.Sprintf("%s/api/notifications/course-completed", p.baseURL)

	payload := map[string]interface{}{
		"user_id":   userID,
		"course_id": courseID,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error al serializar payload: %v", err)
	}

	resp, err := p.client.Post(url, "application/json", 
		io.NopCloser(strings.NewReader(string(payloadBytes))))
	if err != nil {
		log.Printf("Error al notificar Payment Service (no crítico): %v", err)
		return nil // No es crítico si falla
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Payment Service notification failed with status: %d", resp.StatusCode)
	}

	return nil
}

// GetCourseAccessInfo obtiene información detallada de acceso para un curso
func (p *PaymentServiceClient) GetCourseAccessInfo(userID, courseID uint) (*CourseAccessInfo, error) {
	url := fmt.Sprintf("%s/api/pagos/course-access-info?user_id=%d&course_id=%d", 
		p.baseURL, userID, courseID)

	resp, err := p.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error al consultar Payment Service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Payment Service respondió con status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer respuesta: %v", err)
	}

	var response struct {
		Success bool              `json:"success"`
		Data    CourseAccessInfo  `json:"data"`
		Error   string            `json:"error,omitempty"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error al parsear respuesta: %v", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("error del Payment Service: %s", response.Error)
	}

	return &response.Data, nil
}

// CourseAccessInfo información detallada de acceso al curso
type CourseAccessInfo struct {
	HasAccess       bool      `json:"has_access"`
	PaymentStatus   string    `json:"payment_status"`
	PaymentDate     *time.Time `json:"payment_date,omitempty"`
	ExpirationDate  *time.Time `json:"expiration_date,omitempty"`
	AccessType      string    `json:"access_type"` // "free", "paid", "trial"
}