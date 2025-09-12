# 🎓 Plataforma de Cursos Online - Arquitectura de Microservicios

Una plataforma completa de cursos online construida con arquitectura de microservicios, diseñada para escalabilidad, rendimiento y mantenibilidad.

## 🏗️ Arquitectura

### Microservicios Implementados

| Servicio | Puerto | Descripción | Tecnologías |
|----------|---------|-------------|-------------|
| **Auth & User Service** | 8081 | Autenticación, usuarios y perfiles | Go, MySQL, JWT |
| **Course Service** | 8083 | Gestión de cursos, capítulos y progreso | Go, MySQL, File Storage |
| **Payment Service** | 8084 | Procesamiento de pagos y facturación | Go, MySQL, MercadoPago |
| **Analytics Service** | 8085 | Métricas y analíticas de usuario | Go, MongoDB |
| **Contact Service** | 8086 | Sistema de contacto y mensajería | Go, MySQL, Email |
| **Home Service** | 8087 | Gestión de contenido del inicio | Go, MySQL |
| **Portfolio Service** | 8088 | Portafolio de proyectos | Go, MySQL |

### Stack Tecnológico

#### Backend
- **Lenguaje**: Go (Golang)
- **Framework**: Gorilla Mux
- **Base de Datos Principal**: MySQL
- **Base de Datos NoSQL**: MongoDB (Analytics)
- **Cache**: Redis
- **Autenticación**: JWT
- **API Gateway**: Nginx

#### Frontend
- **Framework**: React.js
- **Gestión de Estado**: Context API
- **HTTP Client**: Axios
- **UI**: CSS personalizado + componentes reutilizables

#### Infraestructura
- **Contenedores**: Docker
- **Orquestación**: Docker Compose
- **Proxy Reverso**: Apache/Nginx
- **Storage**: Sistema de archivos local
- **Email**: SMTP integrado

## 🚀 Características Principales

### Para Estudiantes
- ✅ Registro y autenticación segura
- 📚 Catálogo de cursos con búsqueda y filtros
- 🎥 Reproductor de video integrado
- 📊 Seguimiento de progreso por curso
- 💳 Pagos seguros con MercadoPago
- 📧 Notificaciones por email
- 👤 Perfil de usuario personalizable

### Para Instructores
- 📝 Panel de administración completo
- 🎬 Subida y gestión de contenido
- 📈 Analytics de cursos y estudiantes
- 💰 Dashboard de ingresos
- 📋 Gestión de estudiantes inscritos

### Para Administradores
- 🔧 Panel de administración avanzado
- 📊 Métricas y estadísticas globales
- 👥 Gestión de usuarios y roles
- 💼 Gestión del portafolio
- 📬 Sistema de mensajes y contacto

## 🛠️ Instalación y Despliegue

### Prerrequisitos
- Docker y Docker Compose
- Go 1.19+
- Node.js 16+
- MySQL
- Git

### Instalación Local

1. **Clonar el repositorio**
```bash
git clone <repository-url>
cd microservicios
```

2. **Configurar variables de entorno**
```bash
# Copiar archivos de configuración de ejemplo
cp auth-user-service/.env.example auth-user-service/.env
cp course-service/.env.example course-service/.env
# Repetir para cada servicio...
```

3. **Inicializar la base de datos**
```bash
mysql -u root -p < init-db.sql
```

4. **Construir y ejecutar con Docker**
```bash
docker-compose up --build
```

5. **Iniciar servicios individualmente (desarrollo)**
```bash
./start-microservices.sh
```

### Despliegue en Producción

Ver los scripts de despliegue incluidos:
- `deploy-vps-production.sh` - Despliegue completo en VPS
- `build-production.sh` - Build optimizado para producción
- `nginx-ssl-setup.sh` - Configuración SSL con Let's Encrypt

## 📁 Estructura del Proyecto

```
microservicios/
├── auth-user-service/          # Servicio de autenticación y usuarios
├── course-service/             # Servicio de cursos y contenido
├── payment-service/            # Servicio de pagos
├── analytics-service/          # Servicio de analíticas
├── contact-service/           # Servicio de contacto
├── home-service/              # Servicio de página principal
├── portfolio-service/         # Servicio de portafolio
├── frontend/                  # Aplicación React
├── backend/                   # Backend monolítico (legacy)
├── docker-compose.yml         # Configuración de contenedores
├── nginx-proxy.conf          # Configuración del proxy
├── init-db.sql              # Script de inicialización de BD
└── scripts de despliegue/
```

## 🔧 Configuración

### Variables de Entorno Principales

Cada microservicio requiere su propio archivo `.env` con:
- Configuración de base de datos
- Claves de JWT
- Configuración de email
- APIs externas (MercadoPago, etc.)
- Configuración de CORS

### Base de Datos

La aplicación utiliza MySQL como base de datos principal con las siguientes características:
- Esquema normalizado para datos transaccionales
- MongoDB para analytics y logs
- Redis para cache y sesiones

## 🚦 API Endpoints

### Servicios Principales

```
Auth Service (8081):
POST /api/auth/login
POST /api/auth/register
GET  /api/auth/profile
PUT  /api/auth/profile

Course Service (8083):
GET  /api/courses
POST /api/courses
GET  /api/courses/:id
POST /api/courses/:id/enroll
GET  /api/progress/:userId/:courseId

Payment Service (8084):
POST /api/payments/create
POST /api/payments/confirm
GET  /api/payments/history

Analytics Service (8085):
GET  /api/analytics/course/:id
GET  /api/analytics/user/:id
POST /api/analytics/track
```

## 📊 Monitoreo y Logs

- **Logs de aplicación**: Cada servicio genera logs en formato JSON
- **Métricas**: Sistema de analytics integrado
- **Health Checks**: Endpoints de salud en cada servicio
- **Error Tracking**: Logs centralizados con rotación automática

## 🔒 Seguridad

- **Autenticación**: JWT con expiración automática
- **Autorización**: Middleware de roles y permisos
- **CORS**: Configuración restrictiva por servicio
- **Validación**: Sanitización de inputs en todos los endpoints
- **HTTPS**: SSL/TLS en producción
- **Rate Limiting**: Protección contra ataques de fuerza bruta

## 📈 Escalabilidad

### Horizontal
- Cada microservicio puede escalarse independientemente
- Load balancer incluido para distribución de carga
- Base de datos optimizada con índices apropiados

### Vertical
- Configuración de recursos por servicio
- Cache multi-nivel (Redis + aplicación)
- Optimización de queries y conexiones de BD


#
## 👥 Equipo

Desarrollado con ❤️ por Marcelo Nelson

---


---

**Estado del Proyecto**: 🟢 Activo y en desarrollo



**Última Actualización**: Agosto 2025
