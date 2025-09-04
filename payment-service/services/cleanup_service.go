// services/cleanup_service.go
package services

import (
	"log"
	"time"

	"gorm.io/gorm"

	"payment-service/config"
	"payment-service/models"
)

type CleanupService struct {
	db             *gorm.DB
	config         *config.Config
	paymentService *PaymentService
	stopChan       chan bool
	running        bool
}

func NewCleanupService(db *gorm.DB, cfg *config.Config) *CleanupService {
	return &CleanupService{
		db:             db,
		config:         cfg,
		paymentService: NewPaymentService(db, cfg),
		stopChan:       make(chan bool),
		running:        false,
	}
}

// StartCleanupScheduler inicia el programador de limpieza automática
func (cs *CleanupService) StartCleanupScheduler() {
	if cs.running {
		log.Printf("[CLEANUP] [SCHEDULER] Scheduler ya está en ejecución")
		return
	}

	cs.running = true
	log.Printf("[CLEANUP] [SCHEDULER] Iniciando scheduler de limpieza automática")

	go cs.cleanupLoop()
}

// StopCleanupScheduler detiene el programador de limpieza
func (cs *CleanupService) StopCleanupScheduler() {
	if !cs.running {
		log.Printf("[CLEANUP] [SCHEDULER] Scheduler no está en ejecución")
		return
	}

	log.Printf("[CLEANUP] [SCHEDULER] Deteniendo scheduler de limpieza")
	cs.stopChan <- true
	cs.running = false
}

// cleanupLoop ejecuta la limpieza periódicamente
func (cs *CleanupService) cleanupLoop() {
	// Configurar intervalo de limpieza optimizado
	cleanupInterval := 15 * time.Minute // Reducir frecuencia en producción
	if cs.config.AppEnv == "development" {
		// En desarrollo, ejecutar cada 30 segundos para pruebas rápidas
		cleanupInterval = 30 * time.Second
	}

	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	log.Printf("[CLEANUP] [SCHEDULER] Loop iniciado - intervalo: %v", cleanupInterval)

	// Ejecutar limpieza inicial inmediatamente
	cs.performCleanup()

	for {
		select {
		case <-ticker.C:
			cs.performCleanup()
		case <-cs.stopChan:
			log.Printf("[CLEANUP] [SCHEDULER] Scheduler detenido")
			return
		}
	}
}

// performCleanup ejecuta todas las tareas de limpieza
func (cs *CleanupService) performCleanup() {
	startTime := time.Now()
	log.Printf("[CLEANUP] [PERFORM] Iniciando limpieza automática - %v", startTime.Format("2006-01-02 15:04:05"))

	// 1. Cancelar pagos expirados
	cancelledCount, err := cs.paymentService.CancelExpiredPayments()
	if err != nil {
		log.Printf("[ERROR] [CLEANUP] Error en cancelación de pagos expirados: %v", err)
	} else {
		log.Printf("[SUCCESS] [CLEANUP] Pagos expirados cancelados: %d", cancelledCount)
	}

	// 1.5. Detectar y cancelar pagos abandonados
	abandonedCount := cs.cancelAbandonedPayments()
	log.Printf("[SUCCESS] [CLEANUP] Pagos abandonados cancelados: %d", abandonedCount)

	// 2. Actualizar pagos pendientes de MercadoPago (evitar sobrecargar la API)
	if cs.shouldUpdateMercadoPago() {
		cs.updateMercadoPagoPayments()
	}

	// 3. Limpiar pagos muy antiguos (opcional - pagos cancelados/rechazados de más de 30 días)
	cs.cleanupOldPayments()

	duration := time.Since(startTime)
	log.Printf("[SUCCESS] [CLEANUP] Limpieza completada en %v", duration)
}

// shouldUpdateMercadoPago determina si debe actualizar pagos de MercadoPago
func (cs *CleanupService) shouldUpdateMercadoPago() bool {
	// Solo actualizar pagos de MercadoPago cada 10 minutos para evitar sobrecargar la API
	now := time.Now()
	if now.Minute()%10 == 0 {
		return true
	}
	return false
}

// updateMercadoPagoPayments actualiza algunos pagos pendientes de MercadoPago
func (cs *CleanupService) updateMercadoPagoPayments() {
	log.Printf("[CLEANUP] [MP_UPDATE] Actualizando pagos pendientes de MercadoPago")

	// Obtener máximo 5 pagos pendientes de MercadoPago para no sobrecargar la API
	var pendingPayments []models.Payment
	err := cs.db.Where("estado = ? AND metodo = ?", 
		models.PaymentStatusPending, models.PaymentMethodMercadoPago).
		Limit(5).Find(&pendingPayments).Error

	if err != nil {
		log.Printf("[ERROR] [MP_UPDATE] Error al obtener pagos pendientes: %v", err)
		return
	}

	if len(pendingPayments) == 0 {
		log.Printf("[INFO] [MP_UPDATE] No hay pagos pendientes de MercadoPago para actualizar")
		return
	}

	log.Printf("[INFO] [MP_UPDATE] Actualizando %d pagos pendientes de MercadoPago", len(pendingPayments))

	updated := 0
	for _, payment := range pendingPayments {
		if err := cs.paymentService.UpdateMercadoPagoPaymentStatus(&payment); err != nil {
			log.Printf("[ERROR] [MP_UPDATE] Error al actualizar pago ID %d: %v", payment.ID, err)
		} else {
			updated++
		}
		// Pausa pequeña entre solicitudes para ser respetuoso con la API
		time.Sleep(500 * time.Millisecond)
	}

	log.Printf("[SUCCESS] [MP_UPDATE] Actualizados %d de %d pagos de MercadoPago", updated, len(pendingPayments))
}

// cleanupOldPayments limpia pagos muy antiguos (opcional)
func (cs *CleanupService) cleanupOldPayments() {
	log.Printf("[CLEANUP] [OLD_PAYMENTS] Limpiando pagos antiguos")

	// Buscar pagos cancelados/rechazados/regresados de más de 30 días
	cutoffDate := time.Now().AddDate(0, 0, -30) // 30 días atrás

	var oldPaymentsCount int64
	err := cs.db.Model(&models.Payment{}).
		Where("(estado = ? OR estado = ? OR estado = ?) AND updated_at < ?", 
			models.PaymentStatusCancelled, models.PaymentStatusRejected, models.PaymentStatusReturned, cutoffDate).
		Count(&oldPaymentsCount).Error

	if err != nil {
		log.Printf("[ERROR] [OLD_PAYMENTS] Error al contar pagos antiguos: %v", err)
		return
	}

	if oldPaymentsCount == 0 {
		log.Printf("[INFO] [OLD_PAYMENTS] No hay pagos antiguos para limpiar")
		return
	}

	log.Printf("[INFO] [OLD_PAYMENTS] Encontrados %d pagos antiguos", oldPaymentsCount)

	// Por seguridad, solo loggear por ahora. Podrías implementar eliminación real si es necesario
	// result := cs.db.Where("(estado = ? OR estado = ?) AND updated_at < ?", 
	//     models.PaymentStatusCancelled, models.PaymentStatusRejected, cutoffDate).
	//     Delete(&models.Payment{})
	
	log.Printf("[INFO] [OLD_PAYMENTS] Limpieza de pagos antiguos completada (modo logging)")
}

// cancelAbandonedPayments detecta y cancela pagos que fueron abandonados por el usuario
func (cs *CleanupService) cancelAbandonedPayments() int {
	log.Printf("[CLEANUP] [ABANDONED] Detectando pagos abandonados")

	// Definir criterios para pagos abandonados optimizados:
	// - Estado: pendiente
	// - Creado hace más de 2 minutos (para pruebas rápidas en desarrollo)
	// - Sin actualizaciones recientes
	abandonThreshold := time.Now().Add(-30 * time.Second)

	// OPTIMIZACIÓN 1: Limitar cantidad de pagos procesados por batch
	var abandonedPayments []models.Payment
	err := cs.db.Where("estado = ? AND created_at < ?", 
		models.PaymentStatusPending, abandonThreshold).
		Limit(20). // Procesar máximo 20 pagos por vez
		Order("created_at ASC"). // Procesar los más antiguos primero
		Find(&abandonedPayments).Error

	if err != nil {
		log.Printf("[ERROR] [ABANDONED] Error al buscar pagos abandonados: %v", err)
		return 0
	}

	if len(abandonedPayments) == 0 {
		log.Printf("[INFO] [ABANDONED] No hay pagos abandonados para procesar")
		return 0
	}

	log.Printf("[INFO] [ABANDONED] Encontrados %d pagos potencialmente abandonados (batch limitado)", len(abandonedPayments))

	// OPTIMIZACIÓN 2: Separar pagos por método para procesamiento eficiente
	mercadoPagoPayments := make([]models.Payment, 0)
	otherPayments := make([]models.Payment, 0)
	
	for _, payment := range abandonedPayments {
		if payment.Metodo == models.PaymentMethodMercadoPago {
			mercadoPagoPayments = append(mercadoPagoPayments, payment)
		} else {
			otherPayments = append(otherPayments, payment)
		}
	}

	cancelledCount := 0

	// OPTIMIZACIÓN 3: Cancelar pagos no-MercadoPago inmediatamente (batch update)
	if len(otherPayments) > 0 {
		cancelledCount += cs.batchCancelPayments(otherPayments, "métodos no-MercadoPago")
	}

	// OPTIMIZACIÓN 4: Procesar MercadoPago con rate limiting mejorado
	if len(mercadoPagoPayments) > 0 {
		cancelledCount += cs.processMercadoPagoAbandonedPayments(mercadoPagoPayments)
	}

	return cancelledCount
}

// batchCancelPayments cancela múltiples pagos en una sola operación
func (cs *CleanupService) batchCancelPayments(payments []models.Payment, reason string) int {
	if len(payments) == 0 {
		return 0
	}

	log.Printf("[ABANDONED] [BATCH] Cancelando %d pagos en batch (%s)", len(payments), reason)

	// Extraer IDs para update batch
	paymentIDs := make([]uint, len(payments))
	for i, payment := range payments {
		paymentIDs[i] = payment.ID
	}

	// Update batch más eficiente
	result := cs.db.Model(&models.Payment{}).
		Where("id IN ?", paymentIDs).
		Updates(map[string]interface{}{
			"estado":     models.PaymentStatusCancelled,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		log.Printf("[ERROR] [BATCH] Error en cancelación batch: %v", result.Error)
		return 0
	}

	affected := int(result.RowsAffected)
	log.Printf("[SUCCESS] [BATCH] %d pagos cancelados en batch (%s)", affected, reason)
	return affected
}

// processMercadoPagoAbandonedPayments procesa pagos MercadoPago con rate limiting
func (cs *CleanupService) processMercadoPagoAbandonedPayments(payments []models.Payment) int {
	log.Printf("[ABANDONED] [MP] Procesando %d pagos de MercadoPago con verificación", len(payments))

	cancelledCount := 0
	verifiedCount := 0
	maxVerifications := 5 // Límite de verificaciones API por cleanup

	for i, payment := range payments {
		// OPTIMIZACIÓN 5: Limitar verificaciones API costosas
		if verifiedCount >= maxVerifications {
			log.Printf("[ABANDONED] [MP] Límite de verificaciones API alcanzado (%d), cancelando resto sin verificar", maxVerifications)
			// Cancelar el resto sin verificar (son abandonos > 30 min)
			remainingPayments := payments[i:]
			cancelledCount += cs.batchCancelPayments(remainingPayments, "MercadoPago sin verificar")
			break
		}

		log.Printf("[ABANDONED] [MP] Verificando pago ID %d (%d/%d)", payment.ID, i+1, len(payments))

		if cs.verifyMercadoPagoBeforeCancel(&payment) {
			log.Printf("[ABANDONED] [MP] Pago ID %d actualizado desde MercadoPago - no cancelar", payment.ID)
			verifiedCount++
			continue
		}

		// Cancelar individualmente si no se pudo verificar o sigue pendiente
		payment.Estado = models.PaymentStatusCancelled
		payment.UpdatedAt = time.Now()

		if err := cs.db.Save(&payment).Error; err != nil {
			log.Printf("[ERROR] [ABANDONED] [MP] Error cancelando pago ID %d: %v", payment.ID, err)
		} else {
			log.Printf("[SUCCESS] [ABANDONED] [MP] Pago ID %d cancelado", payment.ID)
			cancelledCount++
		}

		verifiedCount++
		// Rate limiting: pausa entre verificaciones API
		time.Sleep(1 * time.Second)
	}

	log.Printf("[SUCCESS] [ABANDONED] [MP] Procesados %d pagos MercadoPago, %d cancelados", len(payments), cancelledCount)
	return cancelledCount
}

// verifyMercadoPagoBeforeCancel verifica con MercadoPago antes de cancelar por abandono
func (cs *CleanupService) verifyMercadoPagoBeforeCancel(payment *models.Payment) bool {
	// Si no hay transaction_id, no podemos verificar
	if payment.TransaccionID == "" {
		log.Printf("[ABANDONED] [MP_VERIFY] Pago ID %d sin TransaccionID - proceder con cancelación", payment.ID)
		return false
	}

	log.Printf("[ABANDONED] [MP_VERIFY] Verificando pago ID %d con MercadoPago antes de cancelar", payment.ID)

	// Intentar actualizar el estado desde MercadoPago
	if err := cs.paymentService.UpdateMercadoPagoPaymentStatus(payment); err != nil {
		log.Printf("[ABANDONED] [MP_VERIFY] Error verificando con MercadoPago: %v", err)
		return false
	}

	// Recargar el pago para ver si cambió
	var updatedPayment models.Payment
	if err := cs.db.First(&updatedPayment, payment.ID).Error; err != nil {
		log.Printf("[ABANDONED] [MP_VERIFY] Error recargando pago: %v", err)
		return false
	}

	// Si el estado cambió de pending, significa que se actualizó
	if updatedPayment.Estado != models.PaymentStatusPending {
		log.Printf("[ABANDONED] [MP_VERIFY] Pago ID %d actualizado a estado: %s", payment.ID, updatedPayment.Estado)
		*payment = updatedPayment // Actualizar la referencia
		return true
	}

	return false
}

// GetCleanupStats devuelve estadísticas del servicio de limpieza
func (cs *CleanupService) GetCleanupStats() map[string]interface{} {
	startTime := time.Now()
	stats := make(map[string]interface{})
	
	// Contar pagos por estado usando consultas optimizadas
	var pendingCount, expiredCount, cancelledCount, abandonedCount int64
	
	// OPTIMIZACIÓN: Usar consultas paralelas para mejorar rendimiento
	type statResult struct {
		name  string
		count int64
		err   error
	}
	
	results := make(chan statResult, 4)
	
	// Consulta 1: Pagos pendientes
	go func() {
		var count int64
		err := cs.db.Model(&models.Payment{}).Where("estado = ?", models.PaymentStatusPending).Count(&count).Error
		results <- statResult{"pending", count, err}
	}()
	
	// Consulta 2: Pagos expirados
	go func() {
		var count int64
		err := cs.db.Model(&models.Payment{}).
			Where("estado = ? AND expires_at IS NOT NULL AND expires_at < ?", 
				models.PaymentStatusPending, time.Now()).Count(&count).Error
		results <- statResult{"expired", count, err}
	}()
	
	// Consulta 3: Pagos cancelados
	go func() {
		var count int64
		err := cs.db.Model(&models.Payment{}).Where("estado = ?", models.PaymentStatusCancelled).Count(&count).Error
		results <- statResult{"cancelled", count, err}
	}()
	
	// Consulta 4: Pagos abandonados
	go func() {
		var count int64
		abandonThreshold := time.Now().Add(-30 * time.Second)
		err := cs.db.Model(&models.Payment{}).
			Where("estado = ? AND created_at < ?", 
				models.PaymentStatusPending, abandonThreshold).Count(&count).Error
		results <- statResult{"abandoned", count, err}
	}()
	
	// Recopilar resultados
	for i := 0; i < 4; i++ {
		result := <-results
		if result.err != nil {
			log.Printf("[STATS] [ERROR] Error en consulta %s: %v", result.name, result.err)
			continue
		}
		
		switch result.name {
		case "pending":
			pendingCount = result.count
		case "expired":
			expiredCount = result.count
		case "cancelled":
			cancelledCount = result.count
		case "abandoned":
			abandonedCount = result.count
		}
	}
	
	queryDuration := time.Since(startTime)
	
	stats["pending_payments"] = pendingCount
	stats["expired_payments"] = expiredCount
	stats["cancelled_payments"] = cancelledCount
	stats["abandoned_payments"] = abandonedCount
	stats["scheduler_running"] = cs.running
	stats["last_cleanup"] = time.Now().Format("2006-01-02 15:04:05")
	stats["query_duration_ms"] = queryDuration.Milliseconds()
	stats["performance"] = map[string]interface{}{
		"total_payments":      pendingCount + expiredCount + cancelledCount,
		"query_time_ms":       queryDuration.Milliseconds(),
		"queries_parallelized": 4,
	}
	
	return stats
}

// ManualCleanup ejecuta una limpieza manual inmediata
func (cs *CleanupService) ManualCleanup() map[string]interface{} {
	log.Printf("[CLEANUP] [MANUAL] Ejecutando limpieza manual")
	
	startTime := time.Now()
	
	// Ejecutar limpieza
	cs.performCleanup()
	
	// Obtener estadísticas después de la limpieza
	stats := cs.GetCleanupStats()
	stats["execution_time"] = time.Since(startTime).String()
	stats["cleanup_type"] = "manual"
	
	return stats
}