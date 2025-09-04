# Auth User Service

Microservicio de autenticacion y gestion de usuarios para la plataforma de cursos.

## Funcionalidades

- **Autenticacion**: Login/logout de usuarios
- **Registro**: Registro de nuevos usuarios
- **JWT**: Gestion de tokens JWT
- **Email**: Recuperacion de contraseñas por email
- **Admin**: Panel de administracion
- **Logs**: Registro de actividad de usuarios
- **Archivos**: Gestion de archivos de perfil
- **Seguridad**: Middleware de autorizacion

## Endpoints Principales

### Autenticacion
- `POST /api/auth/login` - Iniciar sesion
- `POST /api/auth/register` - Registrar usuario
- `POST /api/auth/logout` - Cerrar sesion
- `POST /api/auth/refresh` - Renovar token

### Usuarios
- `GET /api/users/profile` - Obtener perfil
- `PUT /api/users/profile` - Actualizar perfil
- `POST /api/users/upload-avatar` - Subir foto de perfil

### Recuperacion de Contraseña
- `POST /api/password/forgot` - Solicitar recuperacion
- `POST /api/password/reset` - Restablecer contraseña

### Admin
- `GET /api/admin/users` - Listar usuarios (Admin)
- `PUT /api/admin/users/:id` - Gestionar usuarios (Admin)
- `GET /api/admin/activity` - Logs de actividad (Admin)

## Configuracion

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

## Instalacion

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
- `usuarios` - Informacion de usuarios
- `password_resets` - Tokens de recuperacion
- `activity_logs` - Logs de actividad
- `notification_settings` - Configuracion de notificaciones

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
- **Gomail** - Envio de emails
- **MySQL Driver** - Conexion a MySQL

## Puerto

Servidor ejecutandose en puerto **8081**