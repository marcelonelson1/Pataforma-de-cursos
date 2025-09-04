package services

import (
	"auth-user-service/config"
	"auth-user-service/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTService struct{}

func NewJWTService() *JWTService {
	return &JWTService{}
}

// GenerateToken genera un nuevo token JWT
func (j *JWTService) GenerateToken(userID uint, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}

// ValidateToken valida un token JWT y retorna los claims
func (j *JWTService) ValidateToken(tokenString string) (*models.Claims, error) {
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}

// RefreshToken genera un nuevo token basado en uno existente
func (j *JWTService) RefreshToken(userID uint, role string) (string, error) {
	return j.GenerateToken(userID, role)
}