// services/external_client.go
package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"payment-service/config"
)

// UserInfo representa información básica de un usuario (debe coincidir con auth-user-service)
type UserInfo struct {
	ID     uint   `json:"id"`
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

// CourseInfo representa información básica de un curso
type CourseInfo struct {
	ID          uint    `json:"id"`
	Titulo      string  `json:"titulo"`
	Descripcion string  `json:"descripcion"`
	Precio      float64 `json:"precio"`
	Estado      string  `json:"estado"`
	InstructorID uint   `json:"instructor_id"`
}

// GetUserByID obtiene información de un usuario del User Service usando la ruta interna
func GetUserByID(userID uint, cfg *config.Config) (*UserInfo, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Usar la ruta interna que no requiere autenticación
	url := fmt.Sprintf("%s/api/internal/users/%d", cfg.UserServiceURL, userID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error al crear request: %v", err)
	}

	// Agregar header de identificación del servicio (opcional pero recomendado)
	req.Header.Set("X-Service-Name", "payment-service")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al conectar con User Service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("usuario no encontrado")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("User Service retornó status %d", resp.StatusCode)
	}

	var response struct {
		Success bool     `json:"success"`
		Data    UserInfo `json:"data"`
		Message string   `json:"message,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error al decodificar respuesta del User Service: %v", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("User Service retornó error: %s", response.Message)
	}

	return &response.Data, nil
}

// GetCourseByID obtiene información de un curso del Course Service
func GetCourseByID(courseID uint, cfg *config.Config) (*CourseInfo, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	url := fmt.Sprintf("%s/api/courses/%d", cfg.CourseServiceURL, courseID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error al crear request: %v", err)
	}

	// Agregar header de identificación del servicio
	req.Header.Set("X-Service-Name", "payment-service")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al conectar con Course Service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("curso no encontrado")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Course Service retornó status %d", resp.StatusCode)
	}

	// Intentar primero con estructura envuelta
	var response struct {
		Success bool       `json:"success"`
		Data    CourseInfo `json:"data"`
		Message string     `json:"message,omitempty"`
	}

	// Crear una copia del body para poder leerlo dos veces si es necesario
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err == nil && response.Success {
		return &response.Data, nil
	}

	// Si falla la primera decodificación, hacer nueva petición para intentar decodificación directa
	resp2, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al reconectar con Course Service: %v", err)
	}
	defer resp2.Body.Close()

	var course CourseInfo
	if err := json.NewDecoder(resp2.Body).Decode(&course); err != nil {
		return nil, fmt.Errorf("error al decodificar respuesta del Course Service: %v", err)
	}

	return &course, nil
}

// ValidateUserAccess valida si un usuario tiene acceso a realizar una acción
func ValidateUserAccess(userID uint, cfg *config.Config) error {
	user, err := GetUserByID(userID, cfg)
	if err != nil {
		return fmt.Errorf("error al validar acceso del usuario: %v", err)
	}

	if user == nil {
		return fmt.Errorf("usuario no encontrado")
	}

	// Aquí se pueden agregar más validaciones según las reglas de negocio
	// Por ejemplo: verificar si el usuario está activo, roles específicos, etc.
	
	return nil
}

// ValidateCourseAccess valida si un curso está disponible para compra
// ValidateCourseAccess valida si un curso está disponible para compra
func ValidateCourseAccess(courseID uint, cfg *config.Config) (*CourseInfo, error) {
	course, err := GetCourseByID(courseID, cfg)
	if err != nil {
		return nil, fmt.Errorf("error al validar acceso al curso: %v", err)
	}

	if course == nil {
		return nil, fmt.Errorf("curso no encontrado")
	}

	// VALIDACIÓN MEJORADA - Acepta "Publicado" y otros estados
	estadosValidos := []string{"published", "active", "disponible", "Publicado", "activo"}
	estadoValido := false
	
	for _, estadoAceptado := range estadosValidos {
		if course.Estado == estadoAceptado {
			estadoValido = true
			break
		}
	}

	if !estadoValido {
		return nil, fmt.Errorf("el curso no está disponible para compra (estado: %s)", course.Estado)
	}

	// Validar que tenga precio válido
	if course.Precio <= 0 {
		return nil, fmt.Errorf("el curso no tiene un precio válido")
	}

	return course, nil
}
// GetUserAndCourseInfo obtiene información tanto del usuario como del curso en una sola función
// Útil para validaciones de compra
func GetUserAndCourseInfo(userID, courseID uint, cfg *config.Config) (*UserInfo, *CourseInfo, error) {
	// Obtener información del usuario
	user, err := GetUserByID(userID, cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("error al obtener información del usuario: %v", err)
	}

	// Obtener información del curso
	course, err := ValidateCourseAccess(courseID, cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("error al obtener información del curso: %v", err)
	}

	return user, course, nil
}

// HealthCheckExternalServices verifica el estado de los servicios externos
func HealthCheckExternalServices(cfg *config.Config) map[string]bool {
	services := make(map[string]bool)
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Check User Service
	userServiceURL := fmt.Sprintf("%s/api/health", cfg.UserServiceURL)
	if resp, err := client.Get(userServiceURL); err == nil {
		services["user-service"] = resp.StatusCode == http.StatusOK
		resp.Body.Close()
	} else {
		services["user-service"] = false
	}

	// Check Course Service
	courseServiceURL := fmt.Sprintf("%s/api/health", cfg.CourseServiceURL)
	if resp, err := client.Get(courseServiceURL); err == nil {
		services["course-service"] = resp.StatusCode == http.StatusOK
		resp.Body.Close()
	} else {
		services["course-service"] = false
	}

	return services
}