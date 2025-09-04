# Auth User Service

Microservicio de autenticación y gestión de usuarios para la plataforma de cursos.

## Funcionalidades

- = Autenticación de usuarios (login/logout)
- =d Registro de nuevos usuarios
- = Gestión de tokens JWT
- =ç Recuperación de contraseñas por email
- =h=¼ Panel de administración
- =Ê Logs de actividad de usuarios
- =Á Gestión de archivos de perfil
- =á Middleware de autorización

## Endpoints Principales

### Autenticación
- `POST /api/auth/login` - Iniciar sesión
- `POST /api/auth/register` - Registrar usuario
- `POST /api/auth/logout` - Cerrar sesión
- `POST /api/auth/refresh` - Renovar token

### Usuarios
- `GET /api/users/profile` - Obtener perfil
- `PUT /api/users/profile` - Actualizar perfil
- `POST /api/users/upload-avatar` - Subir foto de perfil

### Recuperación de Contraseña
- `POST /api/password/forgot` - Solicitar recuperación
- `POST /api/password/reset` - Restablecer contraseña

### Admin
- `GET /api/admin/users` - Listar usuarios (Admin)
- `PUT /api/admin/users/:id` - Gestionar usuarios (Admin)
- `GET /api/admin/activity` - Logs de actividad (Admin)

## Configuración

### Variables de Entorno
```bash
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=cursos_db
JWT_SECRET=your-secret-key
EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=587
EMAIL_USER=your-email@gmail.com
EMAIL_PASS=your-password
```

## Instalación

```bash
# Instalar dependencias
go mod tidy

# Compilar
go build -o auth-user-service-binary

# Ejecutar
./auth-user-service-binary
```

## Base de Datos

El servicio utiliza MySQL y gestiona las siguientes tablas:
- `usuarios` - Información de usuarios
- `password_resets` - Tokens de recuperación
- `activity_logs` - Logs de actividad
- `notification_settings` - Configuración de notificaciones

## Docker

```bash
# Construir imagen
docker build -t auth-user-service .

# Ejecutar container
docker run -p 8081:8081 auth-user-service
```

## Dependencias

- **Gin** - Framework web
- **GORM** - ORM para base de datos
- **JWT-Go** - Manejo de tokens JWT
- **Gomail** - Envío de emails
- **MySQL Driver** - Conexión a MySQL

## Puerto

Servidor ejecutándose en puerto **8081**