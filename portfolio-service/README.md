# Portfolio Service

Microservicio para gestion del portafolio de proyectos del sistema.

## Funcionalidades

- ✅ Visualizacion publica de proyectos activos
- ✅ Filtrado por categorias
- ✅ CRUD completo de proyectos (admin)
- ✅ Gestion de imagenes de proyectos
- ✅ Sistema de ordenamiento personalizado
- ✅ Estados activo/inactivo de proyectos
- ✅ Estadisticas del portafolio

## Endpoints

### Publicos
- `GET /api/portfolio` - Listar proyectos activos
- `GET /api/portfolio/:id` - Obtener proyecto especifico
- `GET /api/portfolio/category/:category` - Proyectos por categoria
- `GET /api/portfolio/categories` - Listar categorias disponibles

### Administracion (requiere autenticacion + rol admin)
- `GET /api/admin/portfolio` - Listar todos los proyectos
- `POST /api/admin/portfolio` - Crear nuevo proyecto
- `PUT /api/admin/portfolio/:id` - Actualizar proyecto
- `DELETE /api/admin/portfolio/:id` - Eliminar proyecto
- `POST /api/admin/portfolio/upload-image` - Subir imagen de proyecto
- `DELETE /api/admin/portfolio/:id/image` - Eliminar imagen de proyecto
- `POST /api/admin/portfolio/reorder` - Reordenar proyectos
- `PATCH /api/admin/portfolio/:id/toggle` - Activar/desactivar proyecto
- `GET /api/admin/portfolio/stats` - Estadisticas del portafolio

### Archivos estaticos
- `GET /static/portfolio/:filename` - Acceso a imagenes

### Health Check
- `GET /api/health` - Verificacion de salud del servicio

## Configuracion

Variables de entorno requeridas:

```env
PORT=8005
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=yourpassword
DB_NAME=cursos_db
JWT_SECRET=your-secret-key
APP_ENV=development
```

## Categorias validas

- `web` - Proyectos web
- `mobile` - Aplicaciones moviles
- `desktop` - Aplicaciones de escritorio
- `design` - Proyectos de diseno
- `marketing` - Proyectos de marketing
- `other` - Otros proyectos

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
portfolio-service/
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

### ProjectPortfolio
```go
type ProjectPortfolio struct {
    ID          uint      `json:"id"`
    Title       string    `json:"title"`
    Category    string    `json:"category"`
    Description string    `json:"description"`
    ImageURL    string    `json:"image_url"`
    OrderIndex  int       `json:"order_index"`
    IsActive    bool      `json:"is_active"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

## Subida de imagenes

- **Formatos soportados**: JPEG, PNG, GIF, WebP
- **Tamano maximo**: 5MB
- **Validacion**: Tipo MIME y extension
- **Almacenamiento**: Local en `/static/portfolio/`

## Integracion con otros servicios

- **Auth-User-Service**: Validacion de tokens JWT y roles de admin
- **Frontend**: API para mostrar portafolio publico y administracion

## Ejemplo de uso

### Crear proyecto
```bash
POST /api/admin/portfolio
{
  "title": "Mi Proyecto Web",
  "category": "web",
  "description": "Descripcion del proyecto",
  "order_index": 1,
  "is_active": true
}
```

### Subir imagen
```bash
POST /api/admin/portfolio/upload-image
Content-Type: multipart/form-data

project_id: 1
image: [archivo de imagen]
```