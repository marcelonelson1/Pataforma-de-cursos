// services/auto_updater.go - Actualizador automático de pagos
package services

import (
	"log"
	"time"
	
	"payment-service/config"
	"payment-service/models"
	"gorm.io/gorm"
)

// AutoUpdater maneja la actualización automática de pagos pendientes
type AutoUpdater struct {
	db      *gorm.DB
	config  *config.Config
	service *PaymentService
	ticker  *time.Ticker
	done    chan bool
}

// NewAutoUpdater crea una nueva instancia del actualizador automático
func NewAutoUpdater(db *gorm.DB, cfg *config.Config) *AutoUpdater {
	return &AutoUpdater{
		db:      db,
		config:  cfg,
		service: NewPaymentService(db, cfg),
		done:    make(chan bool),
	}
}

// Start inicia el actualizador automático
func (au *AutoUpdater) Start() {
	log.Printf("[AUTO_UPDATER] Iniciando actualizador automático de pagos")
	
	// Actualizar cada 15 segundos para desarrollo local (más agresivo para MercadoPago)
	au.ticker = time.NewTicker(15 * time.Second)
	
	// Primera actualización inmediata
	go au.updatePendingPayments()
	
	// Loop principal
	go func() {
		for {
			select {
			case <-au.ticker.C:
				au.updatePendingPayments()
			case <-au.done:
				log.Printf("[AUTO_UPDATER] Deteniendo actualizador automático")
				return
			}
		}
	}()
	
	log.Printf("✅ [AUTO_UPDATER] Actualizador automático iniciado (intervalo: 30 segundos)")
}

// Stop detiene el actualizador automático
func (au *AutoUpdater) Stop() {
	if au.ticker != nil {
		au.ticker.Stop()
	}
	au.done <- true
}

// updatePendingPayments actualiza todos los pagos pendientes
func (au *AutoUpdater) updatePendingPayments() {
	log.Printf("[AUTO_UPDATER] Iniciando actualización automática de pagos pendientes")
	
	// Obtener pagos pendientes de MercadoPago de los últimos 30 minutos (más frecuente)
	var pendingPayments []models.Payment
	thirtyMinutesAgo := time.Now().Add(-30 * time.Minute)
	
	err := au.db.Where("estado = ? AND metodo = ? AND created_at > ?", 
		models.PaymentStatusPending, 
		models.PaymentMethodMercadoPago,
		thirtyMinutesAgo).
		Order("created_at DESC"). // Más recientes primero
		Find(&pendingPayments).Error
	
	if err != nil {
		log.Printf("[ERROR] [AUTO_UPDATER] Error al obtener pagos pendientes: %v", err)
		return
	}
	
	if len(pendingPayments) == 0 {
		log.Printf("[AUTO_UPDATER] No hay pagos pendientes recientes para actualizar")
		return
	}
	
	log.Printf("[AUTO_UPDATER] Encontrados %d pagos pendientes para actualizar", len(pendingPayments))
	
	updated := 0
	errors := 0
	
	for _, payment := range pendingPayments {
		log.Printf("[AUTO_UPDATER] Procesando pago ID: %d", payment.ID)
		
		// Usar la función existente del PaymentService
		if err := au.service.UpdateMercadoPagoPaymentStatus(&payment); err != nil {
			log.Printf("[WARNING] [AUTO_UPDATER] No se pudo actualizar pago %d: %v", payment.ID, err)
			errors++
		} else {
			// Verificar si cambió el estado
			var updatedPayment models.Payment
			if err := au.db.First(&updatedPayment, payment.ID).Error; err == nil {
				if updatedPayment.Estado != models.PaymentStatusPending {
					log.Printf("✅ [AUTO_UPDATER] Pago %d actualizado: %s → %s", 
						payment.ID, models.PaymentStatusPending, updatedPayment.Estado)
					updated++
				}
			}
		}
		
		// Pequeña pausa para no sobrecargar la API
		time.Sleep(500 * time.Millisecond)
	}
	
	log.Printf("🎯 [AUTO_UPDATER] Actualización completada: %d actualizados, %d errores de %d pagos", 
		updated, errors, len(pendingPayments))
}

// UpdateSinglePayment actualiza un pago específico (para uso inmediato)
func (au *AutoUpdater) UpdateSinglePayment(paymentID uint) error {
	log.Printf("⚡ [AUTO_UPDATER] Actualización inmediata de pago ID: %d", paymentID)
	
	var payment models.Payment
	if err := au.db.First(&payment, paymentID).Error; err != nil {
		return err
	}
	
	if payment.Metodo != models.PaymentMethodMercadoPago {
		log.Printf("ℹ️ [AUTO_UPDATER] Pago %d no es de MercadoPago, saltando", paymentID)
		return nil
	}
	
	return au.service.UpdateMercadoPagoPaymentStatus(&payment)
}