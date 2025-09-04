package models

import (
	"time"
)

// Pago representa una transacción de pago en el sistema
type Pago struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UsuarioID     uint      `gorm:"not null" json:"usuario_id"`
	CursoID       uint      `gorm:"not null" json:"curso_id"`
	Monto         float64   `gorm:"type:decimal(10,2);not null" json:"monto"`
	Metodo        string    `gorm:"size:50;not null" json:"metodo"`
	Estado        string    `gorm:"size:20;not null;default:'pendiente'" json:"estado"`
	TransaccionID string    `gorm:"size:100" json:"transaccion_id"`
	Moneda        string    `gorm:"size:10" json:"moneda"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TableName especifica el nombre de la tabla para el modelo Pago
func (Pago) TableName() string {
	return "pagos"
}

// PagoRequest representa una solicitud de pago enviada por el cliente
type PagoRequest struct {
	CursoID         uint             `json:"curso_id" binding:"required"`
	Monto           float64          `json:"monto" binding:"required"`
	Metodo          string           `json:"metodo" binding:"required"`
	DetallesTarjeta *DetallesTarjeta `json:"detalles_tarjeta,omitempty"`
	Moneda          string           `json:"moneda,omitempty"`
}

// DetallesTarjeta representa los datos de una tarjeta para procesar un pago
type DetallesTarjeta struct {
	Numero     string `json:"numero"`
	Expiracion string `json:"expiracion"`
	CVV        string `json:"cvv"`
}