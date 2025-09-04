package models

// CoinbaseCharge representa una solicitud de pago a Coinbase Commerce
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

// CoinbasePrice representa el precio en una moneda específica
type CoinbasePrice struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

// CoinbaseMetadata representa metadatos para una transacción de Coinbase
type CoinbaseMetadata struct {
	PagoID    uint `json:"pago_id"`
	CursoID   uint `json:"curso_id"`
	UsuarioID uint `json:"usuario_id"`
}

// CoinbaseChargeResponse representa la respuesta de Coinbase Commerce al crear un cargo
type CoinbaseChargeResponse struct {
	Data CoinbaseCharge `json:"data"`
}

// CoinbaseWebhookEvent representa un evento de webhook recibido de Coinbase Commerce
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