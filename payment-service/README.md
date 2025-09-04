# Payment Service

Microservicio de gestion de pagos que integra multiples proveedores de pago para la plataforma de cursos.

## Funcionalidades

- **MercadoPago**: Integracion con MercadoPago
- **PayPal**: Soporte para PayPal  
- **Coinbase**: Integracion con Coinbase (criptomonedas)
- **Webhooks**: Webhooks para actualizaciones automaticas
- **Auto-updater**: Auto-updater para seguimiento de pagos
- **JWT**: Autenticacion JWT
- **Email**: Notificaciones de estado de pago
- **Limpieza**: Limpieza automatica de pagos antiguos

## Proveedores de Pago

### MercadoPago
- Integracion completa con API v1
- Retorno automatico despues del pago
- Webhook para actualizaciones en tiempo real
- Auto-updater cada 15 segundos
- Soporte para USD y ARS

### PayPal
- Pagos internacionales
- Sandbox y produccion
- Webhooks integrados

### Coinbase
- Pagos con criptomonedas
- Bitcoin, Ethereum, etc.
- Conversion automatica a USD

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

## Configuracion

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

# Aplicacion
APP_PORT=8084
APP_ENV=development
AUTH_SERVICE_URL=http://localhost:8081
COURSE_SERVICE_URL=http://localhost:8083
```

## Instalacion

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

### 2. Redireccion
- Usuario es redirigido al proveedor de pago
- Completa el pago en la plataforma externa

### 3. Retorno/Webhook
- Webhook actualiza estado automaticamente
- Auto-updater verifica pagos pendientes cada 15s
- Usuario es redirigido al curso o pagina de estado

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
- **UUID** - IDs unicos para pagos
- **HTTP Client** - Comunicacion con APIs externas
- **Godotenv** - Variables de entorno

## Integracion

Se integra con:
- **Auth Service** (puerto 8081) - Validacion de usuarios
- **Course Service** (puerto 8083) - Informacion de cursos
- **Frontend** - Redirecciones y estados de pago

## Monitoreo

- **Logs**: Logs detallados de transacciones
- **Auto-limpieza**: Auto-limpieza de pagos antiguos (>30 dias)
- **Auto-updater**: Auto-updater para pagos recientes
- **Alertas**: Alertas de fallos de pago

## Puerto

Servidor ejecutandose en puerto **8084**