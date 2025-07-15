package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWTSecret clave secreta para firmar los tokens
var JWTSecret []byte

// Claims estructura de reclamaciones personalizada para JWT
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// InitJWT inicializa la clave secreta para JWT
func InitJWT() {
	JWTSecret = []byte(GetEnv("JWT_SECRET", "mi_clave_secreta_muy_segura"))
}

// GenerateToken genera un nuevo token JWT
func GenerateToken(userID uint, role string) (string, error) {
	// Aumentamos el tiempo de expiración para evitar desconexiones frecuentes
	expirationTime := time.Now().Add(24 * time.Hour)
	
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

// ValidateToken valida un token JWT y retorna sus claims
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}