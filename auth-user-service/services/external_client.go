package services

import (
	"auth-user-service/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type ExternalClient struct {
	httpClient *http.Client
}

func NewExternalClient() *ExternalClient {
	return &ExternalClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// PaymentServiceClient maneja llamadas al servicio de pagos
type PaymentServiceClient struct {
	client  *ExternalClient
	baseURL string
}

func NewPaymentServiceClient() *PaymentServiceClient {
	return &PaymentServiceClient{
		client:  NewExternalClient(),
		baseURL: config.AppConfig.PaymentServiceURL,
	}
}

// CourseServiceClient maneja llamadas al servicio de cursos
type CourseServiceClient struct {
	client  *ExternalClient
	baseURL string
}

func NewCourseServiceClient() *CourseServiceClient {
	return &CourseServiceClient{
		client:  NewExternalClient(),
		baseURL: config.AppConfig.CourseServiceURL,
	}
}

// makeRequest realiza una petición HTTP genérica
func (e *ExternalClient) makeRequest(method, url string, body interface{}, headers map[string]string) (*http.Response, error) {
	var reqBody io.Reader
	
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error al serializar body: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error al crear request: %v", err)
	}

	// Headers por defecto
	req.Header.Set("Content-Type", "application/json")
	
	// Headers adicionales
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al realizar request: %v", err)
	}

	return resp, nil
}

// NotifyUserCreated notifica a otros servicios que se creó un usuario
func (p *PaymentServiceClient) NotifyUserCreated(userID uint, email string) error {
	url := fmt.Sprintf("%s/api/internal/user-created", p.baseURL)
	
	payload := map[string]interface{}{
		"user_id": userID,
		"email":   email,
	}

	resp, err := p.client.makeRequest("POST", url, payload, nil)
	if err != nil {
		log.Printf("Error al notificar creación de usuario al Payment Service: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Payment Service respondió con status %d para notificación de usuario", resp.StatusCode)
	}

	return nil
}

// NotifyUserDeleted notifica que se eliminó un usuario
func (p *PaymentServiceClient) NotifyUserDeleted(userID uint) error {
	url := fmt.Sprintf("%s/api/internal/user-deleted/%d", p.baseURL, userID)
	
	resp, err := p.client.makeRequest("DELETE", url, nil, nil)
	if err != nil {
		log.Printf("Error al notificar eliminación de usuario al Payment Service: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Payment Service respondió con status %d para eliminación de usuario", resp.StatusCode)
	}

	return nil
}

// GetUserPayments obtiene los pagos de un usuario
func (p *PaymentServiceClient) GetUserPayments(userID uint, authToken string) (interface{}, error) {
	url := fmt.Sprintf("%s/api/pagos/user/history", p.baseURL)
	
	headers := map[string]string{
		"Authorization": "Bearer " + authToken,
	}

	resp, err := p.client.makeRequest("GET", url, nil, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("payment service respondió con status %d", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// NotifyUserCreated notifica al servicio de cursos sobre un nuevo usuario
func (c *CourseServiceClient) NotifyUserCreated(userID uint, email string) error {
	url := fmt.Sprintf("%s/api/internal/user-created", c.baseURL)
	
	payload := map[string]interface{}{
		"user_id": userID,
		"email":   email,
	}

	resp, err := c.client.makeRequest("POST", url, payload, nil)
	if err != nil {
		log.Printf("Error al notificar creación de usuario al Course Service: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Course Service respondió con status %d para notificación de usuario", resp.StatusCode)
	}

	return nil
}

// NotifyUserDeleted notifica al servicio de cursos sobre eliminación de usuario
func (c *CourseServiceClient) NotifyUserDeleted(userID uint) error {
	url := fmt.Sprintf("%s/api/internal/user-deleted/%d", c.baseURL, userID)
	
	resp, err := c.client.makeRequest("DELETE", url, nil, nil)
	if err != nil {
		log.Printf("Error al notificar eliminación de usuario al Course Service: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Course Service respondió con status %d para eliminación de usuario", resp.StatusCode)
	}

	return nil
}

// GetUserCourses obtiene los cursos de un usuario
func (c *CourseServiceClient) GetUserCourses(userID uint, authToken string) (interface{}, error) {
	url := fmt.Sprintf("%s/api/courses/user/%d", c.baseURL, userID)
	
	headers := map[string]string{
		"Authorization": "Bearer " + authToken,
	}

	resp, err := c.client.makeRequest("GET", url, nil, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("course service respondió con status %d", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}