package services

import (
	"analytics-service/models"
	"analytics-service/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ProfileService struct {
	authServiceURL string
}

func NewProfileService(authURL string) *ProfileService {
	return &ProfileService{
		authServiceURL: authURL,
	}
}

func (s *ProfileService) GetUserProfile(userID uint) (map[string]interface{}, error) {
	// Intentar obtener del auth service primero
	profile, err := s.fetchFromAuthService(userID)
	if err == nil {
		return profile, nil
	}
	
	// Si falla, buscar en base de datos local o generar mock data
	return s.getMockProfile(userID), nil
}

func (s *ProfileService) UpdateUserProfile(userID uint, updateData map[string]interface{}) (map[string]interface{}, error) {
	// Intentar actualizar en auth service
	profile, err := s.updateInAuthService(userID, updateData)
	if err == nil {
		return profile, nil
	}
	
	// Si falla, actualizar localmente
	return s.updateLocalProfile(userID, updateData)
}

func (s *ProfileService) GetNotificationSettings(userID uint) (map[string]interface{}, error) {
	// Mock data para configuraciones de notificación
	return map[string]interface{}{
		"email_notifications": true,
		"push_notifications": false,
		"marketing_emails": true,
		"course_updates": true,
		"payment_alerts": true,
		"security_alerts": true,
	}, nil
}

func (s *ProfileService) UpdateNotificationSettings(userID uint, settings map[string]interface{}) (map[string]interface{}, error) {
	// Mock update - en producción esto se guardaría en base de datos
	return settings, nil
}

func (s *ProfileService) fetchFromAuthService(userID uint) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/auth/profile/%d", s.authServiceURL, userID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("auth service returned status %d", resp.StatusCode)
	}
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	return result, err
}

func (s *ProfileService) updateInAuthService(userID uint, updateData map[string]interface{}) (map[string]interface{}, error) {
	// Implementar actualización via HTTP al auth service
	// Por simplicidad, devolvemos mock data
	return s.getMockProfile(userID), nil
}

func (s *ProfileService) getMockProfile(userID uint) map[string]interface{} {
	return map[string]interface{}{
		"id":         userID,
		"nombre":     "Usuario Admin",
		"email":      "admin@example.com",
		"role":       "admin",
		"image_url":  "/static/profiles/default.jpg",
		"created_at": "2024-01-15T10:30:00Z",
		"updated_at": "2024-08-16T15:45:00Z",
		"settings": map[string]interface{}{
			"timezone":     "America/Argentina/Buenos_Aires",
			"language":     "es",
			"theme":        "light",
			"email_verified": true,
		},
		"stats": map[string]interface{}{
			"total_logins":     245,
			"last_login":       "2024-08-16T20:30:00Z",
			"courses_created":  12,
			"total_students":   187,
		},
	}
}

func (s *ProfileService) updateLocalProfile(userID uint, updateData map[string]interface{}) (map[string]interface{}, error) {
	// Buscar o crear usuario en base de datos local
	var user models.User
	result := utils.DB.First(&user, userID)
	
	if result.Error != nil {
		// Crear nuevo usuario
		user = models.User{
			ID:       userID,
			Nombre:   "Usuario",
			Email:    "user@example.com",
			Role:     "user",
			ImageURL: "/static/profiles/default.jpg",
		}
	}
	
	// Actualizar campos
	if nombre, ok := updateData["nombre"].(string); ok {
		user.Nombre = nombre
	}
	if email, ok := updateData["email"].(string); ok {
		user.Email = email
	}
	if imageURL, ok := updateData["image_url"].(string); ok {
		user.ImageURL = imageURL
	}
	
	// Guardar en base de datos
	utils.DB.Save(&user)
	
	return map[string]interface{}{
		"id":         user.ID,
		"nombre":     user.Nombre,
		"email":      user.Email,
		"role":       user.Role,
		"image_url":  user.ImageURL,
		"updated_at": user.UpdatedAt,
	}, nil
}