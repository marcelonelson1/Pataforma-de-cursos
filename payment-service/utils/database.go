// utils/database.go
package utils

import (
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"payment-service/config"
	"payment-service/models"
)

// ConnectDatabase conecta a la base de datos MariaDB/MySQL
func ConnectDatabase(cfg *config.Config) *gorm.DB {
	var db *gorm.DB
	var err error
	maxRetries := 5
	retryDelay := 5 * time.Second

	dsn := cfg.Database.DSN()

	for i := 0; i < maxRetries; i++ {
		log.Printf("Conectando a la base de datos (intento %d/%d)...", i+1, maxRetries)

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err == nil {
			sqlDB, err := db.DB()
			if err != nil {
				log.Fatalf("Error al obtener instancia DB: %v", err)
			}

			// Configurar pool de conexiones
			sqlDB.SetMaxIdleConns(10)
			sqlDB.SetMaxOpenConns(100)
			sqlDB.SetConnMaxLifetime(time.Hour)

			// Verificar conexión
			if err := sqlDB.Ping(); err != nil {
				log.Printf("Error al hacer ping a la base de datos: %v", err)
				time.Sleep(retryDelay)
				continue
			}

			log.Println("✅ Conexión a la base de datos establecida correctamente")
			return db
		}

		log.Printf("❌ Error de conexión: %v. Reintentando en %v...", err, retryDelay)
		time.Sleep(retryDelay)
	}

	log.Fatalf("❌ No se pudo conectar a la base de datos después de %d intentos: %v", maxRetries, err)
	return nil
}

// MigrateDatabase ejecuta las migraciones de la base de datos
func MigrateDatabase(db *gorm.DB) {
	log.Println("🔄 Ejecutando migraciones de base de datos...")

	// Deshabilitar foreign key checks temporalmente
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		log.Printf("Advertencia: No se pudo deshabilitar FOREIGN_KEY_CHECKS: %v", err)
	}

	// GORM AutoMigrate - debe detectar automáticamente la nueva columna expires_at
	log.Println("🔧 Ejecutando GORM AutoMigrate...")
	err := db.AutoMigrate(
		&models.Payment{},
		&models.Transaction{},
		&models.PaymentMethod{},
	)

	if err != nil {
		log.Fatalf("❌ Error en migraciones: %v", err)
	}
	log.Println("✅ GORM AutoMigrate completado")

	// Crear métodos de pago por defecto si no existen
	createDefaultPaymentMethods(db)

	// Habilitar foreign key checks nuevamente
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error; err != nil {
		log.Printf("Advertencia: No se pudo habilitar FOREIGN_KEY_CHECKS: %v", err)
	}

	log.Println("✅ Migraciones completadas exitosamente")
}

// createDefaultPaymentMethods crea los métodos de pago por defecto
func createDefaultPaymentMethods(db *gorm.DB) {
	// Verificar si ya existen métodos de pago
	var count int64
	db.Model(&models.PaymentMethod{}).Count(&count)

	if count > 0 {
		log.Println("ℹ️  Métodos de pago ya existen, omitiendo creación")
		return
	}

	log.Println("🔄 Creando métodos de pago por defecto...")

	defaultMethods := models.DefaultPaymentMethods()
	for _, method := range defaultMethods {
		if err := db.Create(&method).Error; err != nil {
			log.Printf("Advertencia: No se pudo crear método de pago %s: %v", method.Name, err)
		} else {
			log.Printf("✅ Método de pago creado: %s", method.Name)
		}
	}
}

// CleanupTestData limpia datos de prueba (útil para testing)
func CleanupTestData(db *gorm.DB) {
	if err := db.Exec("DELETE FROM transactions WHERE 1=1").Error; err != nil {
		log.Printf("Error al limpiar transactions: %v", err)
	}

	if err := db.Exec("DELETE FROM pagos WHERE 1=1").Error; err != nil {
		log.Printf("Error al limpiar pagos: %v", err)
	}

	log.Println("🧹 Datos de prueba limpiados")
}

// GetDatabaseStats obtiene estadísticas de la base de datos
func GetDatabaseStats(db *gorm.DB) map[string]interface{} {
	stats := make(map[string]interface{})

	// Contar pagos por estado
	var paymentStats []struct {
		Estado string
		Count  int64
	}

	db.Model(&models.Payment{}).
		Select("estado, COUNT(*) as count").
		Group("estado").
		Find(&paymentStats)

	stats["payments_by_status"] = paymentStats

	// Contar total de pagos
	var totalPayments int64
	db.Model(&models.Payment{}).Count(&totalPayments)
	stats["total_payments"] = totalPayments

	// Contar transacciones
	var totalTransactions int64
	db.Model(&models.Transaction{}).Count(&totalTransactions)
	stats["total_transactions"] = totalTransactions

	// Métodos de pago activos
	var activePaymentMethods int64
	db.Model(&models.PaymentMethod{}).Where("is_active = ?", true).Count(&activePaymentMethods)
	stats["active_payment_methods"] = activePaymentMethods

	return stats
}

// runManualMigrations ejecuta migraciones manuales que GORM no detecta automáticamente
func runManualMigrations(db *gorm.DB) {
	log.Println("🔧 Ejecutando migraciones manuales...")

	// Migración 1: Agregar columna expires_at si no existe
	if !db.Migrator().HasColumn(&models.Payment{}, "expires_at") {
		log.Println("📅 Agregando columna expires_at a tabla pagos...")
		
		if err := db.Migrator().AddColumn(&models.Payment{}, "expires_at"); err != nil {
			log.Printf("❌ Error al agregar columna expires_at: %v", err)
		} else {
			log.Println("✅ Columna expires_at agregada exitosamente")
			
			// Crear índice para expires_at
			if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_expires_at ON pagos(expires_at)").Error; err != nil {
				log.Printf("⚠️ Advertencia: No se pudo crear índice para expires_at: %v", err)
			} else {
				log.Println("✅ Índice para expires_at creado exitosamente")
			}
		}
	} else {
		log.Println("ℹ️ Columna expires_at ya existe")
	}

	// Aquí puedes agregar más migraciones futuras
	// Migración 2: Ejemplo para futuras migraciones
	// if !db.Migrator().HasColumn(&models.Payment{}, "nueva_columna") {
	//     log.Println("Agregando nueva_columna...")
	//     db.Migrator().AddColumn(&models.Payment{}, "nueva_columna")
	// }

	log.Println("✅ Migraciones manuales completadas")
}