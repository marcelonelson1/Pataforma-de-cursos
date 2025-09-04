# Home Service

Microservicio para gestion de imagenes del home/carousel del sistema.

## Funcionalidades

- ✅ Visualizacion publica de imagenes del home activas
- ✅ CRUD completo de imagenes del home (admin)
- ✅ Gestion de archivos de imagenes
- ✅ Sistema de ordenamiento personalizado
- ✅ Estados activo/inactivo de imagenes
- ✅ Estadisticas de imagenes del home

## Endpoints

### Publicos
- `GET /api/home/images` - Listar imagenes activas del home
- `GET /api/home/images/:id` - Obtener imagen especifica

### Administracion (requiere autenticacion + rol admin)
- `GET /api/admin/home/images` - Listar todas las imagenes
- `POST /api/admin/home/images` - Crear nueva imagen
- `PUT /api/admin/home/images/:id` - Actualizar imagen
- `DELETE /api/admin/home/images/:id` - Eliminar imagen
- `POST /api/admin/home/upload-image` - Subir archivo de imagen
- `DELETE /api/admin/home/images/:id/file` - Eliminar archivo de imagen
- `POST /api/admin/home/images/reorder` - Reordenar imagenes
- `PATCH /api/admin/home/images/:id/toggle` - Activar/desactivar imagen
- `GET /api/admin/home/stats` - Estadisticas de imagenes

### Archivos estaticos
- `GET /static/home/:filename` - Acceso a imagenes

### Health Check
- `GET /api/health` - Verificacion de salud del servicio

## Configuracion

Variables de entorno requeridas:

```env
PORT=8006
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=yourpassword
DB_NAME=home_db
JWT_SECRET=your-secret-key
APP_ENV=development
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
home-service/
├── config/          # Configuracion
├── controllers/     # Controladores HTTP
├── models/          # Modelos de datos
├── services/        # Logica de negocio
├── routes/          # Definicion de rutas
├── utils/           # Utilidades
├── middleware/      # Middlewares
├── static/          # Archivos estaticos
├── main.go         # Punto de entrada
├── Dockerfile      # Imagen Docker
└── docker-compose.yml
```

## Modelo de Datos

### HomeImage
```go
type HomeImage struct {
    ID        uint      `json:"id"`
    ImageURL  string    `json:"image_url"`
    Title     string    `json:"title"`
    Subtitle  string    `json:"subtitle"`
    Order     int       `json:"order"`
    IsActive  bool      `json:"is_active"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

## Subida de imagenes

- **Formatos soportados**: JPEG, PNG, GIF, WebP
- **Tamano maximo**: 10MB
- **Validacion**: Tipo MIME y extension
- **Almacenamiento**: Local en `/static/home/`

## Integracion con otros servicios

- **Auth-User-Service**: Validacion de tokens JWT y roles de admin
- **Frontend**: API para carousel/slider del home

## Ejemplo de uso

### Crear imagen del home
```bash
POST /api/admin/home/images
{
  "title": "Bienvenido a nuestro sitio",
  "subtitle": "La mejor plataforma de cursos online",
  "order": 1,
  "is_active": true
}
```

### Subir imagen
```bash
POST /api/admin/home/upload-image
Content-Type: multipart/form-data

image_id: 1
image: [archivo de imagen]
```