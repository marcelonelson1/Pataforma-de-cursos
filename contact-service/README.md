# Contact Service

Microservicio para gestion de mensajes de contacto del sistema de cursos.

## Funcionalidades

- ✅ Recepcion de mensajes de contacto (publico)
- ✅ Gestion de mensajes para administradores
- ✅ Sistema de notificaciones por email
- ✅ Marcado de mensajes como leidos/estrella
- ✅ Respuesta a mensajes de contacto
- ✅ Estadisticas de mensajes

## Endpoints

### Publicos
- `POST /api/contact` - Enviar mensaje de contacto

### Administracion (requiere autenticacion + rol admin)
- `GET /api/admin/messages` - Listar todos los mensajes
- `GET /api/admin/messages/:id` - Obtener mensaje especifico
- `PATCH /api/admin/messages/:id/read` - Marcar como leido/no leido
- `PATCH /api/admin/messages/:id/star` - Marcar/desmarcar estrella
- `DELETE /api/admin/messages/:id` - Eliminar mensaje
- `POST /api/admin/messages/:id/reply` - Responder mensaje
- `GET /api/admin/messages/stats` - Estadisticas de mensajes
- `GET /api/admin/messages/starred` - Mensajes con estrella

### Health Check
- `GET /api/health` - Verificacion de salud del servicio

## Configuracion

Variables de entorno requeridas:

```env
PORT=8004
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=yourpassword
DB_NAME=cursos_db
JWT_SECRET=your-secret-key
APP_ENV=development

# Configuracion de email
EMAIL_FROM=noreply@example.com
EMAIL_PASSWORD=your-email-password
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
CONTACT_EMAIL=admin@example.com
```

## Instalacion y Ejecucion

### Con Docker
```bash
docker-compose up --build
```

### Manual
```bash
go mod tidy
go run main.go
```

## Estructura

```
contact-service/
├── config/          # Configuracion
├── controllers/     # Controladores HTTP
├── models/          # Modelos de datos
├── services/        # Logica de negocio
├── routes/          # Definicion de rutas
├── utils/           # Utilidades
├── middleware/      # Middlewares
├── templates/       # Templates de email
├── static/          # Archivos estaticos
├── main.go         # Punto de entrada
├── Dockerfile      # Imagen Docker
└── docker-compose.yml
```

## Modelo de Datos

### ContactMessage
```go
type ContactMessage struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Phone     string    `json:"phone"`
    Message   string    `json:"message"`
    Read      bool      `json:"read"`
    Starred   bool      `json:"starred"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

## Integracion con otros servicios

- **Auth-User-Service**: Validacion de tokens JWT y roles de admin
- **Frontend**: Recepcion de mensajes de contacto desde formularios web

## Desarrollo

En modo desarrollo (`APP_ENV=development`), los emails se guardan como archivos HTML locales en lugar de enviarse por SMTP.