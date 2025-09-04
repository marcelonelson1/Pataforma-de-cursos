# Analytics Service

Microservicio de análisis y métricas que centraliza todas las estadísticas del sistema.

## Funcionalidades

### 📊 Dashboard y Métricas
- `/api/admin/dashboard` - Métricas completas del dashboard
- `/api/admin/stats` - Estadísticas generales del sistema
- `/api/admin/sales-stats` - Estadísticas de ventas y pagos
- `/api/admin/activity-log` - Registro de actividades del sistema

### 👤 Perfil de Usuario
- `/api/auth/profile` - Obtener y actualizar perfil de usuario
- `/api/auth/notification-settings` - Configuraciones de notificaciones
- `/api/auth/check-admin` - Verificar permisos de administrador

### 🔍 Características
- **Agregación de datos** de todos los microservicios
- **Mock data inteligente** cuando los servicios no están disponibles
- **Logging de actividades** centralizado
- **Métricas en tiempo real** con fallback a datos simulados
- **Gestión de perfiles** con sincronización multi-servicio

## Endpoints Principales

### Admin (Requiere autenticación + rol admin)
```
GET  /api/admin/dashboard       - Dashboard completo
GET  /api/admin/stats          - Estadísticas generales  
GET  /api/admin/sales-stats    - Estadísticas de ventas
GET  /api/admin/activity-log   - Logs de actividad
```

### Usuario (Requiere autenticación)
```
GET  /api/auth/profile              - Obtener perfil
PUT  /api/auth/profile              - Actualizar perfil
GET  /api/auth/notification-settings - Obtener configuraciones
PUT  /api/auth/notification-settings - Actualizar configuraciones
GET  /api/auth/check-admin          - Verificar si es admin
```

### Público
```
GET  /api/health                    - Health check
GET  /                              - Info del servicio
```

## Configuración

### Variables de Entorno
```bash
PORT=8007
DB_NAME=analytics_db
JWT_SECRET=yoyo

# URLs de otros servicios
AUTH_SERVICE_URL=http://localhost:8001
PAYMENT_SERVICE_URL=http://localhost:8002
COURSE_SERVICE_URL=http://localhost:8003
CONTACT_SERVICE_URL=http://localhost:8004
PORTFOLIO_SERVICE_URL=http://localhost:8005
HOME_SERVICE_URL=http://localhost:8006
```

## Instalación y Uso

### Desarrollo Local
```bash
# Instalar dependencias
go mod tidy

# Ejecutar
go run main.go
```

### Docker
```bash
# Construir imagen
docker build -t analytics-service .

# Ejecutar
docker run -p 8007:8007 analytics-service
```

## Integración con Frontend

Este servicio resuelve los errores 404 del frontend proporcionando:

1. **Estadísticas de ventas** que el StatsDashboard necesita
2. **Perfiles de usuario** que AdminPage y ProfileAdmin requieren  
3. **Verificación de admin** para el sistema de permisos
4. **Dashboard metrics** agregadas de todos los servicios

## Datos Mock Inteligentes

Cuando los servicios externos no están disponibles, proporciona:
- Estadísticas realistas de ventas por período
- Métricas de usuario consistentes
- Datos de portfolio y cursos simulados
- Logs de actividad con formato correcto

Esto garantiza que el frontend siempre funcione, incluso si algunos microservicios están inactivos.