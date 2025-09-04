# Guía de Deployment para Microservicios

## Estructura de Microservicios

```
microservicios/
├── auth-user-service/     # Puerto 8001 - Autenticación y usuarios
├── payment-service/       # Puerto 8002 - Pagos y transacciones  
├── course-service/        # Puerto 8003 - Cursos y contenido
├── contact-service/       # Puerto 8004 - Mensajes de contacto
├── portfolio-service/     # Puerto 8005 - Portfolio y proyectos
├── home-service/          # Puerto 8006 - Imágenes del home
├── analytics-service/     # Puerto 8007 - Estadísticas y métricas
└── frontend/              # Frontend React
```

## Deployment Individual de Microservicios

### 1. Variables de Entorno por Servicio

Cada microservicio debe tener su propio `.env`:

```bash
# auth-user-service/.env
PORT=8001
DB_NAME=auth_db
JWT_SECRET=tu_jwt_secret
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=tu_email@gmail.com
SMTP_PASS=tu_password

# payment-service/.env  
PORT=8002
DB_NAME=payment_db
AUTH_SERVICE_URL=http://auth-service:8001

# course-service/.env
PORT=8003
DB_NAME=course_db
AUTH_SERVICE_URL=http://auth-service:8001

# contact-service/.env
PORT=8004
DB_NAME=contact_db
AUTH_SERVICE_URL=http://auth-service:8001
CONTACT_EMAIL=contacto@tudominio.com

# portfolio-service/.env
PORT=8005
DB_NAME=portfolio_db
AUTH_SERVICE_URL=http://auth-service:8001

# home-service/.env
PORT=8006
DB_NAME=home_db
AUTH_SERVICE_URL=http://auth-service:8001

# analytics-service/.env
PORT=8007
DB_NAME=analytics_db
AUTH_SERVICE_URL=http://auth-service:8001
```

### 2. Docker Deployment

#### Opción A: Docker individual para cada servicio

```bash
# Construcción individual
cd auth-user-service && docker build -t auth-service .
cd payment-service && docker build -t payment-service .
cd course-service && docker build -t course-service .
cd contact-service && docker build -t contact-service .
cd portfolio-service && docker build -t portfolio-service .
cd home-service && docker build -t home-service .
cd analytics-service && docker build -t analytics-service .

# Ejecución individual
docker run -d -p 8001:8001 --name auth-service auth-service
docker run -d -p 8002:8002 --name payment-service payment-service
docker run -d -p 8003:8003 --name course-service course-service
docker run -d -p 8004:8004 --name contact-service contact-service
docker run -d -p 8005:8005 --name portfolio-service portfolio-service
docker run -d -p 8006:8006 --name home-service home-service
docker run -d -p 8007:8007 --name analytics-service analytics-service
```

#### Opción B: Docker Compose (recomendado)

Crear `docker-compose.yml` en la raíz:

```yaml
version: '3.8'
services:
  auth-service:
    build: ./auth-user-service
    ports:
      - "8001:8001"
    environment:
      - PORT=8001
      - DB_HOST=mysql
    depends_on:
      - mysql

  payment-service:
    build: ./payment-service
    ports:
      - "8002:8002"
    environment:
      - PORT=8002
      - AUTH_SERVICE_URL=http://auth-service:8001
    depends_on:
      - auth-service

  course-service:
    build: ./course-service
    ports:
      - "8003:8003"
    environment:
      - PORT=8003
      - AUTH_SERVICE_URL=http://auth-service:8001
    depends_on:
      - auth-service

  contact-service:
    build: ./contact-service
    ports:
      - "8004:8004"
    environment:
      - PORT=8004
      - AUTH_SERVICE_URL=http://auth-service:8001
    depends_on:
      - auth-service

  portfolio-service:
    build: ./portfolio-service
    ports:
      - "8005:8005"
    environment:
      - PORT=8005
      - AUTH_SERVICE_URL=http://auth-service:8001
    depends_on:
      - auth-service

  home-service:
    build: ./home-service
    ports:
      - "8006:8006"
    environment:
      - PORT=8006
      - AUTH_SERVICE_URL=http://auth-service:8001
    depends_on:
      - auth-service

  analytics-service:
    build: ./analytics-service
    ports:
      - "8007:8007"
    environment:
      - PORT=8007
      - AUTH_SERVICE_URL=http://auth-service:8001
    depends_on:
      - auth-service

  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
    environment:
      - REACT_APP_AUTH_API_URL=http://localhost:8001
      - REACT_APP_PAYMENT_API_URL=http://localhost:8002
      - REACT_APP_COURSE_API_URL=http://localhost:8003
      - REACT_APP_CONTACT_API_URL=http://localhost:8004
      - REACT_APP_PORTFOLIO_API_URL=http://localhost:8005
      - REACT_APP_HOME_API_URL=http://localhost:8006
      - REACT_APP_ANALYTICS_API_URL=http://localhost:8007

  mysql:
    image: mysql:8.0
    ports:
      - "3306:3306"
    environment:
      - MYSQL_ROOT_PASSWORD=rootpassword
      - MYSQL_DATABASE=microservices_db
    volumes:
      - mysql_data:/var/lib/mysql

volumes:
  mysql_data:
```

### 3. Deployment en Cloud

#### AWS ECS/Fargate
- Cada microservicio como un task separado
- Application Load Balancer para routing
- RDS para MySQL
- ElastiCache para Redis (si se usa)

#### Kubernetes
- Cada microservicio como un deployment separado
- Services para comunicación interna
- Ingress para routing externo

#### Digital Ocean/Linode
- Droplets separados para cada servicio
- Load balancer para distribuir tráfico
- Managed databases

### 4. Frontend Deployment

#### Variables de entorno por ambiente:

```bash
# Desarrollo
cp .env.example .env

# Staging  
cp .env.staging .env

# Producción
cp .env.production .env
```

#### Build y deployment:

```bash
# Build para producción
npm run build

# Deployment a Netlify/Vercel
# Configurar variables de entorno en la plataforma

# Deployment manual a servidor
scp -r build/* usuario@servidor:/var/www/html/
```

### 5. Base de Datos

#### Opción A: Base de datos separada por servicio
- auth_db
- payment_db  
- course_db
- contact_db
- portfolio_db
- home_db
- analytics_db

#### Opción B: Base de datos compartida con schemas separados
- microservices_db con prefijos por servicio

### 6. Monitoreo y Logs

- Configurar logs centralizados (ELK Stack, Fluentd)
- Métricas con Prometheus + Grafana
- Health checks en `/health` para cada servicio
- Alertas para servicios caídos

### 7. Seguridad

- HTTPS en todos los servicios
- Rate limiting
- API Gateway para autenticación centralizada
- Secrets management (AWS Secrets, HashiCorp Vault)

### 8. CI/CD Pipeline

```yaml
# .github/workflows/deploy.yml
name: Deploy Microservices
on:
  push:
    branches: [main]

jobs:
  deploy-auth:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Deploy Auth Service
        run: |
          cd auth-user-service
          docker build -t auth-service .
          # Deploy commands

  deploy-payment:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Deploy Payment Service
        run: |
          cd payment-service
          docker build -t payment-service .
          # Deploy commands
```

## Comandos Útiles

```bash
# Iniciar todos los servicios localmente
./start-all-services.sh

# Verificar health de servicios
curl http://localhost:8001/health
curl http://localhost:8002/health
# ... etc

# Logs de servicios
docker logs auth-service
docker logs payment-service

# Scaling horizontal
docker-compose up --scale payment-service=3
```

## Troubleshooting

1. **Servicios no se comunican**: Verificar variables AUTH_SERVICE_URL
2. **Frontend no conecta**: Verificar variables REACT_APP_*_API_URL
3. **Base de datos**: Verificar conexiones y credenciales
4. **CORS**: Configurar origins permitidos en cada servicio