	// models/payment.go
	package models

	import (
		"time"
		"gorm.io/gorm"
	)

	// Payment representa un pago en el sistema (migrado del modelo Pago original)
	type Payment struct {
		ID            uint           `gorm:"primaryKey" json:"id"`
		UsuarioID     uint           `gorm:"not null;index" json:"usuario_id"`
		CursoID       uint           `gorm:"not null;index" json:"curso_id"`
		Monto         float64        `gorm:"type:decimal(10,2);not null" json:"monto"`
		Metodo        string         `gorm:"size:50;not null;index" json:"metodo"`
		Estado        string         `gorm:"size:20;not null;default:'pendiente';index:idx_estado_created_updated" json:"estado"`
		TransaccionID string         `gorm:"size:100;index" json:"transaccion_id"`
		Moneda        string         `gorm:"size:10;default:'USD'" json:"moneda"`
		ExpiresAt     *time.Time     `gorm:"index" json:"expires_at,omitempty"`
		Metadata      string         `gorm:"type:text" json:"metadata,omitempty"`
		CreatedAt     time.Time      `gorm:"index:idx_estado_created_updated" json:"created_at"`
		UpdatedAt     time.Time      `gorm:"index:idx_estado_created_updated" json:"updated_at"`
		DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	}

	// TableName especifica el nombre de la tabla para Payment
	func (Payment) TableName() string {
		return "pagos"
	}

	// Transaction representa los detalles de una transacción específica
	type Transaction struct {
		ID            uint           `gorm:"primaryKey" json:"id"`
		PaymentID     uint           `gorm:"not null;index" json:"payment_id"`
		Provider      string         `gorm:"size:50;not null" json:"provider"` // paypal, coinbase, stripe, etc.
		ExternalID    string         `gorm:"size:100;not null;index" json:"external_id"`
		Status        string         `gorm:"size:20;not null" json:"status"`
		Amount        float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
		Currency      string         `gorm:"size:10;not null" json:"currency"`
		RawResponse   string         `gorm:"type:text" json:"raw_response,omitempty"`
		CreatedAt     time.Time      `json:"created_at"`
		UpdatedAt     time.Time      `json:"updated_at"`
		DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
		
		// Relación
		Payment Payment `gorm:"foreignKey:PaymentID" json:"payment,omitempty"`
	}

	// PaymentRequest estructura para crear un pago (migrado de PagoRequest original)
	type PaymentRequest struct {
		CursoID         uint                  `json:"curso_id" binding:"required"`
		Monto           float64               `json:"monto" binding:"required"`
		Metodo          string                `json:"metodo" binding:"required"`
		DetallesTarjeta *DetallesTarjeta      `json:"detalles_tarjeta,omitempty"`
		Moneda          string                `json:"moneda,omitempty"`
	}

	// DetallesTarjeta para pagos con tarjeta (migrado del original)
	type DetallesTarjeta struct {
		Numero     string `json:"numero"`
		Expiracion string `json:"expiracion"`
		CVV        string `json:"cvv"`
	}

	// PaymentResponse estructura para respuestas de pago
	type PaymentResponse struct {
		ID          uint    `json:"id"`
		UsuarioID   uint    `json:"usuario_id"`
		CursoID     uint    `json:"curso_id"`
		Monto       float64 `json:"monto"`
		Metodo      string  `json:"metodo"`
		Estado      string  `json:"estado"`
		Moneda      string  `json:"moneda"`
		CheckoutURL string  `json:"checkout_url,omitempty"`
		Message     string  `json:"message,omitempty"`
		CreatedAt   time.Time `json:"created_at"`
	}

	// Constantes para estados de pago
	const (
		PaymentStatusPending   = "pendiente"
		PaymentStatusApproved  = "aprobado"
		PaymentStatusRejected  = "rechazado"
		PaymentStatusCancelled = "cancelado"
		PaymentStatusReturned  = "regresado"  // Usuario regresó sin completar el pago
		PaymentStatusRefunded  = "reembolsado"
	)

	// Constantes para métodos de pago
	const (
		PaymentMethodCard         = "tarjeta"
		PaymentMethodPayPal       = "paypal"
		PaymentMethodCoinbase     = "coinbase"
		PaymentMethodTransfer     = "transferencia"
		PaymentMethodStripe       = "stripe"
		PaymentMethodMercadoPago  = "mercadopago"
		PaymentMethodDev          = "dev"
	)

	// ValidPaymentMethods mapa de métodos de pago válidos
	var ValidPaymentMethods = map[string]bool{
		PaymentMethodCard:        true,
		PaymentMethodPayPal:      true,
		PaymentMethodCoinbase:    true,
		PaymentMethodTransfer:    true,
		PaymentMethodStripe:      true,
		PaymentMethodMercadoPago: true,
		PaymentMethodDev:         true,
	}

	// IsValidPaymentMethod valida si un método de pago es válido
	func IsValidPaymentMethod(method string) bool {
		return ValidPaymentMethods[method]
	}

	// IsCompleted verifica si el pago está completado
	func (p *Payment) IsCompleted() bool {
		return p.Estado == PaymentStatusApproved
	}

	// IsPending verifica si el pago está pendiente
	func (p *Payment) IsPending() bool {
		return p.Estado == PaymentStatusPending
	}

	// UpdateStatus actualiza el estado del pago
	func (p *Payment) UpdateStatus(status string) {
		p.Estado = status
		p.UpdatedAt = time.Now()
	}

	// IsExpired verifica si el pago ha expirado
	func (p *Payment) IsExpired() bool {
		if p.ExpiresAt == nil {
			return false
		}
		return time.Now().After(*p.ExpiresAt)
	}

	// SetExpiration establece la fecha de expiración del pago
	func (p *Payment) SetExpiration(duration time.Duration) {
		expirationTime := time.Now().Add(duration)
		p.ExpiresAt = &expirationTime
	}

	// CanBeExpired verifica si el pago puede ser expirado (está pendiente y tiene fecha de expiración)
	func (p *Payment) CanBeExpired() bool {
		return p.IsPending() && p.ExpiresAt != nil
	}