# Payment Service

Microservicio de gestión de pagos que integra múltiples proveedores de pago para la plataforma de cursos.

## Funcionalidades

- =³ Integración con MercadoPago
- < Soporte para PayPal  
- ¿ Integración con Coinbase (criptomonedas)
- = Webhooks para actualizaciones automáticas
- =Ê Auto-updater para seguimiento de pagos
- = Autenticación JWT
- =ç Notificaciones de estado de pago
- =á Limpieza automática de pagos antiguos

## Proveedores de Pago

### MercadoPago <æ<÷
-  Integración completa con API v1
-  Retorno automático después del pago
-  Webhook para actualizaciones en tiempo real
-  Auto-updater cada 15 segundos
-  Soporte para USD y ARS

### PayPal <
-  Pagos internacionales
-  Sandbox y producción
-  Webhooks integrados

### Coinbase ¿
-  Pagos con criptomonedas
-  Bitcoin, Ethereum, etc.
-  Conversión automática a USD

## Endpoints Principales

### Pagos Generales
- `POST /api/pagos` - Crear nuevo pago
- `GET /api/pagos/:id` - Obtener estado del pago
- `GET /api/pagos/user/:userId` - Historial de pagos del usuario

### MercadoPago
- `POST /api/pagos/mercadopago/preference` - Crear preferencia
- `GET /api/pagos/mercadopago/return` - Handler de retorno
- `POST /api/pagos/mercadopago/webhook` - Webhook notifications

### PayPal
- `POST /api/pagos/paypal/create` - Crear pago PayPal
- `POST /api/pagos/paypal/execute` - Ejecutar pago
- `POST /api/pagos/paypal/webhook` - Webhook PayPal

### Coinbase
- `POST /api/pagos/coinbase/create` - Crear pago cripto
- `GET /api/pagos/coinbase/status/:id` - Estado pago cripto
- `POST /api/pagos/coinbase/webhook` - Webhook Coinbase

## Configuración

### Variables de Entorno
```bash
# Base de datos
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=cursos_db

# JWT
JWT_SECRET=your-secret-key

# MercadoPago
MERCADOPAGO_ACCESS_TOKEN=your-access-token
MERCADOPAGO_ENVIRONMENT=sandbox
MERCADOPAGO_ACCEPT_USD=true

# PayPal
PAYPAL_CLIENT_ID=your-client-id
PAYPAL_CLIENT_SECRET=your-secret
PAYPAL_ENV=sandbox

# Coinbase
COINBASE_API_KEY=your-api-key
COINBASE_WEBHOOK_SECRET=your-webhook-secret

# Aplicación
APP_PORT=8084
APP_ENV=development
AUTH_SERVICE_URL=http://localhost:8081
COURSE_SERVICE_URL=http://localhost:8083
```

## Instalación

```bash
# Instalar dependencias
go mod tidy

# Compilar
go build -o payment-service-binary

# Ejecutar
./payment-service-binary
```

## Flujo de Pago

### 1. Crear Pago
```json
POST /api/pagos
{
  "amount": 29.99,
  "currency": "USD",
  "course_id": 123,
  "user_id": 456,
  "payment_method": "mercadopago"
}
```

### 2. Redirección
- Usuario es redirigido al proveedor de pago
- Completa el pago en la plataforma externa

### 3. Retorno/Webhook
- Webhook actualiza estado automáticamente
- Auto-updater verifica pagos pendientes cada 15s
- Usuario es redirigido al curso o página de estado

## Base de Datos

### Tabla `payments`
```sql
- id (UUID)
- user_id (INT)
- course_id (INT) 
- amount (DECIMAL)
- currency (VARCHAR)
- status (ENUM: pending, approved, rejected, cancelled)
- payment_method (ENUM: mercadopago, paypal, coinbase)
- external_id (VARCHAR)
- created_at, updated_at (TIMESTAMP)
```

## Docker

```bash
# Construir imagen
docker build -t payment-service .

# Ejecutar container
docker run -p 8084:8084 payment-service
```

## Dependencias

- **Gin** - Framework web
- **GORM** - ORM para base de datos
- **UUID** - IDs únicos para pagos
- **HTTP Client** - Comunicación con APIs externas
- **Godotenv** - Variables de entorno

## Integración

Se integra con:
- **Auth Service** (puerto 8081) - Validación de usuarios
- **Course Service** (puerto 8083) - Información de cursos
- **Frontend** - Redirecciones y estados de pago

## Monitoreo

- =Ê Logs detallados de transacciones
- = Auto-limpieza de pagos antiguos (>30 días)
- ¡ Auto-updater para pagos recientes
- =¨ Alertas de fallos de pago

## Puerto

Servidor ejecutándose en puerto **8084**