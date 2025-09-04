// models/payment_method.go
package models

import (
	"time"
	"gorm.io/gorm"
)

// PaymentMethod representa los métodos de pago disponibles
type PaymentMethod struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:50;not null" json:"name"`
	Provider    string         `gorm:"size:50;not null" json:"provider"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	Config      string         `gorm:"type:text" json:"config,omitempty"` // JSON config específico
	Description string         `gorm:"size:255" json:"description,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// DefaultPaymentMethods retorna los métodos de pago por defecto
func DefaultPaymentMethods() []PaymentMethod {
	return []PaymentMethod{
		{
			Name:        "PayPal",
			Provider:    "paypal",
			IsActive:    true,
			Description: "Pago seguro con PayPal",
		},
		{
			Name:        "Coinbase Commerce",
			Provider:    "coinbase",
			IsActive:    true,
			Description: "Pago con criptomonedas",
		},
		{
			Name:        "Tarjeta de Crédito",
			Provider:    "stripe",
			IsActive:    true,
			Description: "Visa, Mastercard, American Express",
		},
		{
			Name:        "Transferencia Bancaria",
			Provider:    "transfer",
			IsActive:    true,
			Description: "Transferencia bancaria directa",
		},
		{
			Name:        "Mercado Pago",
			Provider:    "mercadopago",
			IsActive:    true,
			Description: "Pago con Mercado Pago",
		},
		{
			Name:        "Modo Desarrollo",
			Provider:    "dev",
			IsActive:    true,
			Description: "Solo para testing",
		},
	}
}