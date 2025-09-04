#!/bin/bash

# Script completo para deployment de microservicios en VPS
# Uso: ./deploy-vps-complete.sh tudominio.com tu@email.com

set -e  # Salir si hay error

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Función para imprimir mensajes
print_message() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Verificar argumentos
if [ "$#" -ne 2 ]; then
    print_error "Uso: $0 <dominio> <email>"
    print_error "Ejemplo: $0 miapp.com admin@miapp.com"
    exit 1
fi

DOMAIN=$1
EMAIL=$2
PROJECT_DIR="/var/www/microservicios"

print_message "=== DEPLOYMENT DE MICROSERVICIOS EN VPS ==="
print_message "Dominio: $DOMAIN"
print_message "Email: $EMAIL"
print_message "Directorio: $PROJECT_DIR"

# 1. Actualizar sistema
print_message "Actualizando sistema..."
sudo apt update && sudo apt upgrade -y

# 2. Instalar dependencias necesarias
print_message "Instalando dependencias..."
sudo apt install -y \
    docker.io \
    docker-compose \
    nginx \
    certbot \
    python3-certbot-nginx \
    git \
    curl \
    wget \
    unzip \
    ufw

# 3. Configurar Docker
print_message "Configurando Docker..."
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker $USER

# 4. Crear estructura de directorios
print_message "Creando estructura de directorios..."
sudo mkdir -p $PROJECT_DIR
sudo chown -R $USER:$USER $PROJECT_DIR

# 5. Configurar variables de entorno para producción
print_message "Configurando variables de entorno..."

cat > $PROJECT_DIR/.env << EOF
# Configuración de producción
DOMAIN=$DOMAIN
EMAIL=$EMAIL

# Base de datos
DB_HOST=mysql
DB_PORT=3306
DB_USER=microuser
DB_PASSWORD=$(openssl rand -base64 32)
MYSQL_ROOT_PASSWORD=$(openssl rand -base64 32)

# JWT
JWT_SECRET=$(openssl rand -base64 64)

# Email (configurar después)
SMTP_USER=
SMTP_PASS=
CONTACT_EMAIL=$EMAIL

# PayPal (configurar después)
PAYPAL_CLIENT_ID=
PAYPAL_CLIENT_SECRET=

# MercadoPago (configurar después)
MERCADOPAGO_ACCESS_TOKEN=
EOF

# 6. Crear docker-compose para producción
print_message "Creando docker-compose para producción..."

cat > $PROJECT_DIR/docker-compose.prod.yml << 'EOF'
version: '3.8'

services:
  # Base de datos MySQL
  mysql:
    image: mysql:8.0
    container_name: microservices_mysql
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}
      MYSQL_DATABASE: microservices_db
      MYSQL_USER: ${DB_USER}
      MYSQL_PASSWORD: ${DB_PASSWORD}
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./init-db.sql:/docker-entrypoint-initdb.d/init-db.sql
    networks:
      - microservices_network

  # Auth User Service
  auth-service:
    build: ./auth-user-service
    container_name: auth_service
    restart: unless-stopped
    ports:
      - "8001:8001"
    environment:
      - PORT=8001
      - DB_HOST=${DB_HOST}
      - DB_PORT=${DB_PORT}
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=auth_db
      - JWT_SECRET=${JWT_SECRET}
      - SMTP_HOST=smtp.gmail.com
      - SMTP_PORT=587
      - SMTP_USER=${SMTP_USER}
      - SMTP_PASS=${SMTP_PASS}
    depends_on:
      - mysql
    networks:
      - microservices_network

  # Payment Service
  payment-service:
    build: ./payment-service
    container_name: payment_service
    restart: unless-stopped
    ports:
      - "8002:8002"
    environment:
      - PORT=8002
      - DB_HOST=${DB_HOST}
      - DB_PORT=${DB_PORT}
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=payment_db
      - AUTH_SERVICE_URL=http://auth-service:8001
      - PAYPAL_CLIENT_ID=${PAYPAL_CLIENT_ID}
      - PAYPAL_CLIENT_SECRET=${PAYPAL_CLIENT_SECRET}
      - MERCADOPAGO_ACCESS_TOKEN=${MERCADOPAGO_ACCESS_TOKEN}
    depends_on:
      - mysql
      - auth-service
    networks:
      - microservices_network

  # Course Service
  course-service:
    build: ./course-service
    container_name: course_service
    restart: unless-stopped
    ports:
      - "8003:8003"
    environment:
      - PORT=8003
      - DB_HOST=${DB_HOST}
      - DB_PORT=${DB_PORT}
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=course_db
      - AUTH_SERVICE_URL=http://auth-service:8001
    depends_on:
      - mysql
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
    ports:
      - "8004:8004"
    environment:
      - PORT=8004
      - DB_HOST=${DB_HOST}
      - DB_PORT=${DB_PORT}
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=contact_db
      - AUTH_SERVICE_URL=http://auth-service:8001
      - CONTACT_EMAIL=${CONTACT_EMAIL}
      - SMTP_HOST=smtp.gmail.com
      - SMTP_PORT=587
      - SMTP_USER=${SMTP_USER}
      - SMTP_PASS=${SMTP_PASS}
    depends_on:
      - mysql
      - auth-service
    networks:
      - microservices_network

  # Portfolio Service
  portfolio-service:
    build: ./portfolio-service
    container_name: portfolio_service
    restart: unless-stopped
    ports:
      - "8005:8005"
    environment:
      - PORT=8005
      - DB_HOST=${DB_HOST}
      - DB_PORT=${DB_PORT}
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=portfolio_db
      - AUTH_SERVICE_URL=http://auth-service:8001
    depends_on:
      - mysql
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
    ports:
      - "8006:8006"
    environment:
      - PORT=8006
      - DB_HOST=${DB_HOST}
      - DB_PORT=${DB_PORT}
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=home_db
      - AUTH_SERVICE_URL=http://auth-service:8001
    depends_on:
      - mysql
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
    ports:
      - "8007:8007"
    environment:
      - PORT=8007
      - DB_HOST=${DB_HOST}
      - DB_PORT=${DB_PORT}
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=analytics_db
      - AUTH_SERVICE_URL=http://auth-service:8001
    depends_on:
      - mysql
      - auth-service
    networks:
      - microservices_network

  # Frontend React
  frontend:
    build: 
      context: ./frontend
      dockerfile: Dockerfile.prod
    container_name: frontend_app
    restart: unless-stopped
    ports:
      - "3000:80"
    environment:
      - REACT_APP_AUTH_API_URL=https://${DOMAIN}/api/auth
      - REACT_APP_PAYMENT_API_URL=https://${DOMAIN}/api/payment
      - REACT_APP_COURSE_API_URL=https://${DOMAIN}/api/course
      - REACT_APP_CONTACT_API_URL=https://${DOMAIN}/api/contact
      - REACT_APP_PORTFOLIO_API_URL=https://${DOMAIN}/api/portfolio
      - REACT_APP_HOME_API_URL=https://${DOMAIN}/api/home
      - REACT_APP_ANALYTICS_API_URL=https://${DOMAIN}/api/analytics
      - REACT_APP_PAYPAL_CLIENT_ID=${PAYPAL_CLIENT_ID}
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
  mysql_data:
  course_uploads:
  portfolio_uploads:
  home_uploads:
EOF

# 7. Configurar Nginx como reverse proxy
print_message "Configurando Nginx..."

cat > /tmp/nginx-microservices << EOF
server {
    listen 80;
    server_name $DOMAIN www.$DOMAIN;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header Referrer-Policy "no-referrer-when-downgrade" always;
    add_header Content-Security-Policy "default-src 'self' http: https: data: blob: 'unsafe-inline'" always;

    # Frontend React
    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_cache_bypass \$http_upgrade;
    }

    # API Routes
    location /api/auth {
        proxy_pass http://localhost:8001;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }

    location /api/payment {
        proxy_pass http://localhost:8002;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }

    location /api/course {
        proxy_pass http://localhost:8003;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }

    location /api/contact {
        proxy_pass http://localhost:8004;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }

    location /api/portfolio {
        proxy_pass http://localhost:8005;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }

    location /api/home {
        proxy_pass http://localhost:8006;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }

    location /api/analytics {
        proxy_pass http://localhost:8007;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }

    # Static files for uploads
    location /uploads/ {
        alias /var/www/microservicios/uploads/;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
EOF

sudo mv /tmp/nginx-microservices /etc/nginx/sites-available/$DOMAIN
sudo ln -sf /etc/nginx/sites-available/$DOMAIN /etc/nginx/sites-enabled/
sudo rm -f /etc/nginx/sites-enabled/default

# 8. Configurar firewall
print_message "Configurando firewall..."
sudo ufw --force enable
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw reload

# 9. Reiniciar Nginx
print_message "Reiniciando Nginx..."
sudo nginx -t && sudo systemctl restart nginx

print_message "=== DEPLOYMENT INICIAL COMPLETADO ==="
print_warning "PASOS MANUALES PENDIENTES:"
print_warning "1. Copiar código fuente a $PROJECT_DIR"
print_warning "2. Configurar variables de entorno en $PROJECT_DIR/.env"
print_warning "3. Ejecutar: cd $PROJECT_DIR && docker-compose -f docker-compose.prod.yml up -d"
print_warning "4. Configurar SSL: sudo certbot --nginx -d $DOMAIN -d www.$DOMAIN"

print_message "Variables de entorno creadas en: $PROJECT_DIR/.env"
print_message "Configurar las credenciales de email, PayPal y MercadoPago antes del deployment final"