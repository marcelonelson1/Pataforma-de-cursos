#!/bin/bash

# Script de deployment para VPS - bpnnneuq.co
# IP: 181.85.173.171

set -e

echo "🚀 Iniciando deployment en VPS para bpnnneuq.co..."

# Variables
DOMAIN="bpnnneuq.co"
VPS_IP="181.85.173.171"
VPS_USER="kali"
VPS_PASS="pilarLeo123"
PROJECT_DIR="/home/kali/microservicios"
BACKUP_DIR="/home/kali/backups"

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Función para verificar si el comando existe
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Verificar dependencias locales
check_dependencies() {
    print_status "Verificando dependencias locales..."
    
    if ! command_exists rsync; then
        print_error "rsync no está instalado. Instálalo con: apt install rsync"
        exit 1
    fi
    
    if ! command_exists ssh; then
        print_error "ssh no está instalado"
        exit 1
    fi
    
    if ! command_exists sshpass; then
        print_error "sshpass no está instalado. Instálalo con: apt install sshpass"
        exit 1
    fi
    
    print_success "Dependencias verificadas"
}

# Crear archivo de configuración para producción
create_production_env() {
    print_status "Creando archivo .env para producción..."
    
    cat > .env.production << EOF
# Base de datos
DB_HOST=localhost
DB_PORT=3306
DB_USER=microuser
DB_PASSWORD=micropass_prod_$(openssl rand -hex 8)
DB_ROOT_PASSWORD=root_$(openssl rand -hex 12)

# JWT
JWT_SECRET=$(openssl rand -hex 32)

# Email (configurar con tus credenciales)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=tu_email@gmail.com
SMTP_PASS=tu_app_password
CONTACT_EMAIL=contacto@${DOMAIN}

# PayPal (configurar con tus credenciales)
PAYPAL_CLIENT_ID=tu_paypal_client_id
PAYPAL_CLIENT_SECRET=tu_paypal_client_secret

# MercadoPago (configurar con tus credenciales)
MERCADOPAGO_ACCESS_TOKEN=tu_mercadopago_token

# URLs de producción
REACT_APP_AUTH_API_URL=https://${DOMAIN}/api/auth
REACT_APP_PAYMENT_API_URL=https://${DOMAIN}/api/payment
REACT_APP_COURSE_API_URL=https://${DOMAIN}/api/course
REACT_APP_CONTACT_API_URL=https://${DOMAIN}/api/contact
REACT_APP_PORTFOLIO_API_URL=https://${DOMAIN}/api/portfolio
REACT_APP_HOME_API_URL=https://${DOMAIN}/api/home
REACT_APP_ANALYTICS_API_URL=https://${DOMAIN}/api/analytics
REACT_APP_PAYPAL_CLIENT_ID=\${PAYPAL_CLIENT_ID}
EOF

    print_success "Archivo .env.production creado"
    print_warning "IMPORTANTE: Edita .env.production con tus credenciales reales antes de continuar"
    
    read -p "¿Has configurado las credenciales en .env.production? (y/n): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_error "Configura las credenciales y ejecuta el script nuevamente"
        exit 1
    fi
}

# Crear docker-compose para producción
create_production_compose() {
    print_status "Creando docker-compose-production.yml..."
    
    cat > docker-compose-production.yml << 'EOF'
version: '3.8'

services:
  # Base de datos MariaDB
  mariadb:
    image: mariadb:10.11
    container_name: microservices_mariadb
    restart: unless-stopped
    environment:
      MARIADB_ROOT_PASSWORD: ${DB_ROOT_PASSWORD}
      MARIADB_DATABASE: microservices_db
      MARIADB_USER: ${DB_USER}
      MARIADB_PASSWORD: ${DB_PASSWORD}
    ports:
      - "3306:3306"
    volumes:
      - mariadb_data:/var/lib/mysql
      - ./init-db.sql:/docker-entrypoint-initdb.d/init-db.sql
    networks:
      - microservices_network

  # Auth User Service
  auth-service:
    build: ./auth-user-service
    container_name: auth_service
    restart: unless-stopped
    environment:
      - PORT=8001
      - DB_HOST=mariadb
      - DB_PORT=3306
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=auth_db
      - JWT_SECRET=${JWT_SECRET}
      - SMTP_HOST=${SMTP_HOST}
      - SMTP_PORT=${SMTP_PORT}
      - SMTP_USER=${SMTP_USER}
      - SMTP_PASS=${SMTP_PASS}
    depends_on:
      - mariadb
    networks:
      - microservices_network

  # Payment Service
  payment-service:
    build: ./payment-service
    container_name: payment_service
    restart: unless-stopped
    environment:
      - PORT=8002
      - DB_HOST=mariadb
      - DB_PORT=3306
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=payment_db
      - AUTH_SERVICE_URL=http://auth-service:8001
      - PAYPAL_CLIENT_ID=${PAYPAL_CLIENT_ID}
      - PAYPAL_CLIENT_SECRET=${PAYPAL_CLIENT_SECRET}
      - MERCADOPAGO_ACCESS_TOKEN=${MERCADOPAGO_ACCESS_TOKEN}
    depends_on:
      - mariadb
      - auth-service
    networks:
      - microservices_network

  # Course Service
  course-service:
    build: ./course-service
    container_name: course_service
    restart: unless-stopped
    environment:
      - PORT=8003
      - DB_HOST=mariadb
      - DB_PORT=3306
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=course_db
      - AUTH_SERVICE_URL=http://auth-service:8001
    depends_on:
      - mariadb
      - auth-service
    networks:
      - microservices_network
    volumes:
      - course_uploads:/app/static

  # Contact Service
  contact-service:
    build: ./contact-service
    container_name: contact_service
    restart: unless-stopped
    environment:
      - PORT=8004
      - DB_HOST=mariadb
      - DB_PORT=3306
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=contact_db
      - AUTH_SERVICE_URL=http://auth-service:8001
      - CONTACT_EMAIL=${CONTACT_EMAIL}
      - SMTP_HOST=${SMTP_HOST}
      - SMTP_PORT=${SMTP_PORT}
      - SMTP_USER=${SMTP_USER}
      - SMTP_PASS=${SMTP_PASS}
    depends_on:
      - mariadb
      - auth-service
    networks:
      - microservices_network

  # Portfolio Service
  portfolio-service:
    build: ./portfolio-service
    container_name: portfolio_service
    restart: unless-stopped
    environment:
      - PORT=8005
      - DB_HOST=mariadb
      - DB_PORT=3306
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=portfolio_db
      - AUTH_SERVICE_URL=http://auth-service:8001
    depends_on:
      - mariadb
      - auth-service
    networks:
      - microservices_network
    volumes:
      - portfolio_uploads:/app/static

  # Home Service
  home-service:
    build: ./home-service
    container_name: home_service
    restart: unless-stopped
    environment:
      - PORT=8006
      - DB_HOST=mariadb
      - DB_PORT=3306
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=home_db
      - AUTH_SERVICE_URL=http://auth-service:8001
    depends_on:
      - mariadb
      - auth-service
    networks:
      - microservices_network
    volumes:
      - home_uploads:/app/static

  # Analytics Service
  analytics-service:
    build: ./analytics-service
    container_name: analytics_service
    restart: unless-stopped
    environment:
      - PORT=8007
      - DB_HOST=mariadb
      - DB_PORT=3306
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=analytics_db
      - AUTH_SERVICE_URL=http://auth-service:8001
    depends_on:
      - mariadb
      - auth-service
    networks:
      - microservices_network

  # Frontend React
  frontend:
    build: 
      context: ./frontend
      args:
        - REACT_APP_AUTH_API_URL=${REACT_APP_AUTH_API_URL}
        - REACT_APP_PAYMENT_API_URL=${REACT_APP_PAYMENT_API_URL}
        - REACT_APP_COURSE_API_URL=${REACT_APP_COURSE_API_URL}
        - REACT_APP_CONTACT_API_URL=${REACT_APP_CONTACT_API_URL}
        - REACT_APP_PORTFOLIO_API_URL=${REACT_APP_PORTFOLIO_API_URL}
        - REACT_APP_HOME_API_URL=${REACT_APP_HOME_API_URL}
        - REACT_APP_ANALYTICS_API_URL=${REACT_APP_ANALYTICS_API_URL}
        - REACT_APP_PAYPAL_CLIENT_ID=${REACT_APP_PAYPAL_CLIENT_ID}
    container_name: frontend_app
    restart: unless-stopped
    depends_on:
      - auth-service
      - payment-service
      - course-service
      - contact-service
      - portfolio-service
      - home-service
      - analytics-service
    networks:
      - microservices_network

networks:
  microservices_network:
    driver: bridge

volumes:
  mariadb_data:
  course_uploads:
  portfolio_uploads:
  home_uploads:
EOF

    print_success "docker-compose-production.yml creado"
}

# Crear configuración de Nginx para producción
create_nginx_config() {
    print_status "Creando configuración de Nginx..."
    
    cat > nginx-production.conf << EOF
events {
    worker_connections 1024;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    # Logs
    access_log /var/log/nginx/access.log;
    error_log /var/log/nginx/error.log;

    # Gzip compression
    gzip on;
    gzip_comp_level 6;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;

    # Rate limiting
    limit_req_zone \$binary_remote_addr zone=api:10m rate=10r/s;

    # Upstream definitions
    upstream frontend {
        server frontend_app:80;
    }

    upstream auth_api {
        server auth_service:8001;
    }

    upstream payment_api {
        server payment_service:8002;
    }

    upstream course_api {
        server course_service:8003;
    }

    upstream contact_api {
        server contact_service:8004;
    }

    upstream portfolio_api {
        server portfolio_service:8005;
    }

    upstream home_api {
        server home_service:8006;
    }

    upstream analytics_api {
        server analytics_service:8007;
    }

    # HTTP to HTTPS redirect
    server {
        listen 80;
        server_name ${DOMAIN} www.${DOMAIN};
        return 301 https://\$server_name\$request_uri;
    }

    # HTTPS server
    server {
        listen 443 ssl http2;
        server_name ${DOMAIN} www.${DOMAIN};

        # SSL certificates (will be configured with Certbot)
        ssl_certificate /etc/letsencrypt/live/${DOMAIN}/fullchain.pem;
        ssl_certificate_key /etc/letsencrypt/live/${DOMAIN}/privkey.pem;

        # SSL configuration
        ssl_session_timeout 1d;
        ssl_session_cache shared:MozTLS:10m;
        ssl_session_tickets off;
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384;
        ssl_prefer_server_ciphers off;

        # Security headers
        add_header Strict-Transport-Security "max-age=63072000" always;
        add_header X-Frame-Options DENY;
        add_header X-Content-Type-Options nosniff;
        add_header X-XSS-Protection "1; mode=block";

        # Client max body size for file uploads
        client_max_body_size 100M;

        # API routes with rate limiting
        location /api/auth/ {
            limit_req zone=api burst=20 nodelay;
            proxy_pass http://auth_api/;
            proxy_set_header Host \$host;
            proxy_set_header X-Real-IP \$remote_addr;
            proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto \$scheme;
        }

        location /api/payment/ {
            limit_req zone=api burst=10 nodelay;
            proxy_pass http://payment_api/;
            proxy_set_header Host \$host;
            proxy_set_header X-Real-IP \$remote_addr;
            proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto \$scheme;
        }

        location /api/course/ {
            limit_req zone=api burst=20 nodelay;
            proxy_pass http://course_api/;
            proxy_set_header Host \$host;
            proxy_set_header X-Real-IP \$remote_addr;
            proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto \$scheme;
        }

        location /api/contact/ {
            limit_req zone=api burst=5 nodelay;
            proxy_pass http://contact_api/;
            proxy_set_header Host \$host;
            proxy_set_header X-Real-IP \$remote_addr;
            proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto \$scheme;
        }

        location /api/portfolio/ {
            limit_req zone=api burst=15 nodelay;
            proxy_pass http://portfolio_api/;
            proxy_set_header Host \$host;
            proxy_set_header X-Real-IP \$remote_addr;
            proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto \$scheme;
        }

        location /api/home/ {
            limit_req zone=api burst=15 nodelay;
            proxy_pass http://home_api/;
            proxy_set_header Host \$host;
            proxy_set_header X-Real-IP \$remote_addr;
            proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto \$scheme;
        }

        location /api/analytics/ {
            limit_req zone=api burst=10 nodelay;
            proxy_pass http://analytics_api/;
            proxy_set_header Host \$host;
            proxy_set_header X-Real-IP \$remote_addr;
            proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto \$scheme;
        }

        # Frontend
        location / {
            proxy_pass http://frontend;
            proxy_set_header Host \$host;
            proxy_set_header X-Real-IP \$remote_addr;
            proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto \$scheme;
            
            # Handle client-side routing
            try_files \$uri \$uri/ /index.html;
        }

        # Health check endpoint
        location /health {
            access_log off;
            return 200 "healthy\n";
            add_header Content-Type text/plain;
        }
    }
}
EOF

    print_success "Configuración de Nginx creada"
}

# Actualizar Dockerfile del frontend para producción
update_frontend_dockerfile() {
    print_status "Actualizando Dockerfile del frontend..."
    
    cat > frontend/Dockerfile.production << 'EOF'
# Multi-stage build
FROM node:18-alpine as build

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci --only=production

# Copy source code
COPY . .

# Build arguments for environment variables
ARG REACT_APP_AUTH_API_URL
ARG REACT_APP_PAYMENT_API_URL
ARG REACT_APP_COURSE_API_URL
ARG REACT_APP_CONTACT_API_URL
ARG REACT_APP_PORTFOLIO_API_URL
ARG REACT_APP_HOME_API_URL
ARG REACT_APP_ANALYTICS_API_URL
ARG REACT_APP_PAYPAL_CLIENT_ID

# Set environment variables
ENV REACT_APP_AUTH_API_URL=$REACT_APP_AUTH_API_URL
ENV REACT_APP_PAYMENT_API_URL=$REACT_APP_PAYMENT_API_URL
ENV REACT_APP_COURSE_API_URL=$REACT_APP_COURSE_API_URL
ENV REACT_APP_CONTACT_API_URL=$REACT_APP_CONTACT_API_URL
ENV REACT_APP_PORTFOLIO_API_URL=$REACT_APP_PORTFOLIO_API_URL
ENV REACT_APP_HOME_API_URL=$REACT_APP_HOME_API_URL
ENV REACT_APP_ANALYTICS_API_URL=$REACT_APP_ANALYTICS_API_URL
ENV REACT_APP_PAYPAL_CLIENT_ID=$REACT_APP_PAYPAL_CLIENT_ID

# Build the app
RUN npm run build

# Production stage
FROM nginx:alpine

# Copy custom nginx config
COPY nginx.conf /etc/nginx/nginx.conf

# Copy built app
COPY --from=build /app/build /usr/share/nginx/html

# Expose port
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
EOF

    print_success "Dockerfile de producción del frontend creado"
}

# Crear script de instalación para el VPS
create_vps_install_script() {
    print_status "Creando script de instalación para VPS..."
    
    cat > vps-install.sh << EOF
#!/bin/bash

# Script de instalación en VPS
set -e

echo "🔧 Instalando dependencias en VPS..."

# Actualizar sistema
apt update && apt upgrade -y

# Instalar Docker
if ! command -v docker &> /dev/null; then
    echo "Instalando Docker..."
    curl -fsSL https://get.docker.com -o get-docker.sh
    sh get-docker.sh
    usermod -aG docker \$USER
fi

# Instalar Docker Compose
if ! command -v docker-compose &> /dev/null; then
    echo "Instalando Docker Compose..."
    curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-\$(uname -s)-\$(uname -m)" -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
fi

# Instalar Nginx
if ! command -v nginx &> /dev/null; then
    echo "Instalando Nginx..."
    apt install -y nginx
fi

# Instalar Certbot para SSL
if ! command -v certbot &> /dev/null; then
    echo "Instalando Certbot..."
    apt install -y certbot python3-certbot-nginx
fi

# Crear directorios necesarios
mkdir -p ${PROJECT_DIR}
mkdir -p ${BACKUP_DIR}
mkdir -p /var/log/microservices

# Configurar firewall
ufw allow 22
ufw allow 80
ufw allow 443
ufw --force enable

echo "✅ Instalación completada en VPS"
EOF

    chmod +x vps-install.sh
    print_success "Script de instalación para VPS creado"
}

# Función principal de deployment
deploy_to_vps() {
    print_status "Iniciando deployment a VPS..."
    
    # Copiar archivos al VPS
    print_status "Copiando archivos al VPS..."
    sshpass -p "${VPS_PASS}" rsync -avz --delete \
        --exclude 'node_modules' \
        --exclude '.git' \
        --exclude '*.log' \
        --exclude '*.pid' \
        ./ ${VPS_USER}@${VPS_IP}:${PROJECT_DIR}/
    
    # Ejecutar comandos en el VPS
    sshpass -p "${VPS_PASS}" ssh ${VPS_USER}@${VPS_IP} << ENDSSH
        cd ${PROJECT_DIR}
        
        # Ejecutar script de instalación
        chmod +x vps-install.sh
        sudo ./vps-install.sh
        
        # Configurar variables de entorno
        cp .env.production .env
        
        # Construir y ejecutar contenedores
        docker-compose -f docker-compose-production.yml down || true
        docker-compose -f docker-compose-production.yml build --no-cache
        docker-compose -f docker-compose-production.yml up -d
        
        # Configurar Nginx
        sudo cp nginx-production.conf /etc/nginx/nginx.conf
        sudo nginx -t
        sudo systemctl restart nginx
        
        # Configurar SSL
        sudo certbot --nginx -d ${DOMAIN} -d www.${DOMAIN} --non-interactive --agree-tos --email admin@${DOMAIN}
        
        # Verificar servicios
        docker-compose -f docker-compose-production.yml ps
        
        echo "✅ Deployment completado"
ENDSSH

    print_success "Deployment a VPS completado"
}

# Función para verificar el deployment
verify_deployment() {
    print_status "Verificando deployment..."
    
    # Verificar que los servicios respondan
    sleep 10  # Esperar un poco más para que inicien
    for service in auth payment course contact portfolio home analytics; do
        if curl -f -s "https://${DOMAIN}/api/${service}/health" > /dev/null 2>&1; then
            print_success "✅ ${service} service - OK"
        else
            print_warning "⚠️ ${service} service - FAIL (puede estar iniciando)"
        fi
    done
    
    # Verificar frontend
    if curl -f -s "https://${DOMAIN}" > /dev/null; then
        print_success "✅ Frontend - OK"
    else
        print_error "❌ Frontend - FAIL"
    fi
}

# Función principal
main() {
    print_status "=== DEPLOYMENT A PRODUCCIÓN ==="
    print_status "Dominio: ${DOMAIN}"
    print_status "IP VPS: ${VPS_IP}"
    echo
    
    # Verificar dependencias
    check_dependencies
    
    # Crear archivos de configuración
    create_production_env
    create_production_compose
    create_nginx_config
    update_frontend_dockerfile
    create_vps_install_script
    
    # Confirmar deployment
    read -p "¿Continuar con el deployment? (y/n): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_warning "Deployment cancelado"
        exit 0
    fi
    
    # Realizar deployment
    deploy_to_vps
    
    # Esperar un momento para que los servicios inicien
    print_status "Esperando que los servicios inicien..."
    sleep 30
    
    # Verificar deployment
    verify_deployment
    
    print_success "🎉 ¡Deployment completado!"
    print_success "Tu aplicación está disponible en: https://${DOMAIN}"
    
    echo
    print_status "Comandos útiles:"
    echo "  - Ver logs: ssh ${VPS_USER}@${VPS_IP} 'cd ${PROJECT_DIR} && docker-compose -f docker-compose-production.yml logs'"
    echo "  - Reiniciar servicios: ssh ${VPS_USER}@${VPS_IP} 'cd ${PROJECT_DIR} && docker-compose -f docker-compose-production.yml restart'"
    echo "  - Ver estado: ssh ${VPS_USER}@${VPS_IP} 'cd ${PROJECT_DIR} && docker-compose -f docker-compose-production.yml ps'"
}

# Ejecutar función principal
main "$@"