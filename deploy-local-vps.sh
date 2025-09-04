#!/bin/bash

# Script para deployment cuando los archivos YA ESTÁN en el VPS
# Ejecutar este script DESDE DENTRO del VPS

set -e

echo "🚀 Iniciando deployment local en VPS para bpnnneuq.co..."
echo "📍 Ejecutándose desde: $(pwd)"

# Variables
DOMAIN="bpnnneuq.co"
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

# Verificar que estamos en el directorio correcto
check_location() {
    if [[ ! -f "docker-compose.yml" ]] || [[ ! -f ".env.production" ]]; then
        print_error "No se encuentran los archivos necesarios."
        print_error "Asegúrate de estar en el directorio del proyecto y que existan:"
        print_error "  - docker-compose.yml"
        print_error "  - .env.production"
        exit 1
    fi
    print_success "Archivos del proyecto encontrados"
}

# Verificar que Docker está instalado
check_docker() {
    print_status "Verificando Docker..."
    
    if ! command -v docker >/dev/null 2>&1; then
        print_error "Docker no está instalado. Instalando..."
        curl -fsSL https://get.docker.com -o get-docker.sh
        sudo sh get-docker.sh
        sudo usermod -aG docker $USER
        print_warning "Docker instalado. Es recomendable reiniciar la sesión."
    fi
    
    if ! command -v docker-compose >/dev/null 2>&1; then
        print_error "Docker Compose no está instalado. Instalando..."
        sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
        sudo chmod +x /usr/local/bin/docker-compose
    fi
    
    print_success "Docker verificado"
}

# Instalar dependencias del sistema
install_dependencies() {
    print_status "Instalando dependencias del sistema..."
    
    sudo apt update
    
    # Nginx
    if ! command -v nginx >/dev/null 2>&1; then
        print_status "Instalando Nginx..."
        sudo apt install -y nginx
    fi
    
    # Certbot
    if ! command -v certbot >/dev/null 2>&1; then
        print_status "Instalando Certbot..."
        sudo apt install -y certbot python3-certbot-nginx
    fi
    
    # Crear directorios necesarios
    sudo mkdir -p ${BACKUP_DIR}
    sudo mkdir -p /var/log/microservices
    sudo chown -R $USER:$USER ${BACKUP_DIR}
    
    print_success "Dependencias instaladas"
}

# Configurar firewall
setup_firewall() {
    print_status "Configurando firewall..."
    
    sudo ufw allow 22
    sudo ufw allow 80
    sudo ufw allow 443
    sudo ufw --force enable
    
    print_success "Firewall configurado"
}

# Configurar variables de entorno
setup_environment() {
    print_status "Configurando variables de entorno..."
    
    if [[ -f ".env.production" ]]; then
        cp .env.production .env
        print_success "Variables de entorno configuradas desde .env.production"
    else
        print_error "Archivo .env.production no encontrado"
        exit 1
    fi
}

# Detener servicios existentes
stop_existing_services() {
    print_status "Deteniendo servicios existentes..."
    
    # Detener contenedores Docker si existen
    docker-compose down 2>/dev/null || true
    docker-compose -f docker-compose-production.yml down 2>/dev/null || true
    
    # Detener procesos locales si existen
    pkill -f "auth-user-service" 2>/dev/null || true
    pkill -f "payment-service" 2>/dev/null || true
    pkill -f "course-service" 2>/dev/null || true
    pkill -f "contact-service" 2>/dev/null || true
    pkill -f "portfolio-service" 2>/dev/null || true
    pkill -f "home-service" 2>/dev/null || true
    pkill -f "analytics-service" 2>/dev/null || true
    
    print_success "Servicios existentes detenidos"
}

# Construir y ejecutar contenedores
deploy_containers() {
    print_status "Construyendo y ejecutando contenedores..."
    
    # Usar docker-compose-production.yml si existe, sino usar docker-compose.yml
    if [[ -f "docker-compose-production.yml" ]]; then
        COMPOSE_FILE="docker-compose-production.yml"
    else
        COMPOSE_FILE="docker-compose.yml"
    fi
    
    print_status "Usando archivo: ${COMPOSE_FILE}"
    
    # Construir imágenes
    docker-compose -f ${COMPOSE_FILE} build --no-cache
    
    # Ejecutar contenedores
    docker-compose -f ${COMPOSE_FILE} up -d
    
    # Esperar a que los servicios inicien
    print_status "Esperando que los servicios inicien..."
    sleep 30
    
    # Verificar estado
    docker-compose -f ${COMPOSE_FILE} ps
    
    print_success "Contenedores desplegados"
}

# Configurar Nginx
setup_nginx() {
    print_status "Configurando Nginx..."
    
    # Crear configuración de Nginx (forzar recreación)
    if [[ -f "nginx-production.conf" ]]; then
        print_status "Recreando configuración de Nginx..."
        rm -f nginx-production.conf
    fi
    
    print_status "Creando configuración básica de Nginx..."
    cat > nginx-production.conf << 'EOF'
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
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml;

    # Rate limiting
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;

    # Upstream definitions
    upstream frontend {
        server localhost:3000;
    }

    upstream auth_api {
        server localhost:8001;
    }

    upstream payment_api {
        server localhost:8002;
    }

    upstream course_api {
        server localhost:8003;
    }

    upstream contact_api {
        server localhost:8004;
    }

    upstream portfolio_api {
        server localhost:8005;
    }

    upstream home_api {
        server localhost:8006;
    }

    upstream analytics_api {
        server localhost:8007;
    }

    # HTTP to HTTPS redirect
    server {
        listen 80;
        server_name bpnnneuq.co www.bpnnneuq.co;
        return 301 https://$server_name$request_uri;
    }

    # HTTPS server
    server {
        listen 443 ssl http2;
        server_name bpnnneuq.co www.bpnnneuq.co;

        # SSL certificates (configurados por Certbot)
        ssl_certificate /etc/letsencrypt/live/bpnnneuq.co/fullchain.pem;
        ssl_certificate_key /etc/letsencrypt/live/bpnnneuq.co/privkey.pem;

        # SSL configuration
        ssl_session_timeout 1d;
        ssl_session_cache shared:MozTLS:10m;
        ssl_protocols TLSv1.2 TLSv1.3;

        # Security headers
        add_header Strict-Transport-Security "max-age=63072000" always;
        add_header X-Frame-Options DENY;
        add_header X-Content-Type-Options nosniff;

        # Client max body size
        client_max_body_size 100M;

        # API routes
        location /api/auth/ {
            limit_req zone=api burst=20 nodelay;
            rewrite ^/api/auth/(.*) /api/auth/$1 break;
            proxy_pass http://auth_api;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        location /api/payment/ {
            limit_req zone=api burst=10 nodelay;
            proxy_pass http://payment_api/;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        location /api/course/ {
            limit_req zone=api burst=20 nodelay;
            proxy_pass http://course_api/;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        location /api/contact/ {
            limit_req zone=api burst=5 nodelay;
            proxy_pass http://contact_api/;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        location /api/portfolio/ {
            limit_req zone=api burst=15 nodelay;
            proxy_pass http://portfolio_api/;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        location /api/home/ {
            limit_req zone=api burst=15 nodelay;
            proxy_pass http://home_api/;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        location /api/analytics/ {
            limit_req zone=api burst=10 nodelay;
            proxy_pass http://analytics_api/;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        # Frontend
        location / {
            proxy_pass http://frontend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        # Health check
        location /health {
            access_log off;
            return 200 "healthy\n";
            add_header Content-Type text/plain;
        }
    }
}
EOF
    
    # Hacer backup de configuración actual
    sudo cp /etc/nginx/nginx.conf /etc/nginx/nginx.conf.backup.$(date +%Y%m%d_%H%M%S) 2>/dev/null || true
    
    # Aplicar nueva configuración
    sudo cp nginx-production.conf /etc/nginx/nginx.conf
    
    # Probar configuración
    sudo nginx -t
    
    # Reiniciar Nginx
    sudo systemctl restart nginx
    
    print_success "Nginx configurado"
}

# Configurar SSL con Certbot
setup_ssl() {
    print_status "Configurando SSL con Let's Encrypt..."
    
    # Verificar que el dominio apunte al servidor
    print_warning "Asegúrate de que ${DOMAIN} apunte a esta IP antes de continuar"
    read -p "¿El dominio ${DOMAIN} apunta a este servidor? (y/n): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_warning "Configura el DNS y ejecuta el script nuevamente"
        return 0
    fi
    
    # Obtener certificado SSL
    sudo certbot --nginx -d ${DOMAIN} -d www.${DOMAIN} \
        --non-interactive \
        --agree-tos \
        --email admin@${DOMAIN} \
        --redirect || {
        print_warning "SSL falló. Puedes configurarlo manualmente después:"
        print_warning "sudo certbot --nginx -d ${DOMAIN}"
        return 0
    }
    
    # Configurar renovación automática
    sudo systemctl enable certbot.timer
    sudo systemctl start certbot.timer
    
    print_success "SSL configurado correctamente"
}

# Verificar deployment
verify_deployment() {
    print_status "Verificando deployment..."
    
    # Verificar que los contenedores estén corriendo
    print_status "Estado de contenedores:"
    docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
    
    echo ""
    
    # Verificar puertos locales
    print_status "Verificando puertos locales..."
    for port in 8001 8002 8003 8004 8005 8006 8007 3000; do
        if netstat -tuln | grep ":${port} " >/dev/null 2>&1; then
            print_success "✅ Puerto ${port} - ACTIVO"
        else
            print_warning "⚠️ Puerto ${port} - INACTIVO"
        fi
    done
    
    echo ""
    
    # Verificar Nginx
    if sudo nginx -t >/dev/null 2>&1; then
        print_success "✅ Nginx - CONFIGURACIÓN OK"
    else
        print_error "❌ Nginx - CONFIGURACIÓN ERRÓNEA"
    fi
    
    # Verificar servicios del sistema
    for service in nginx docker; do
        if systemctl is-active --quiet $service; then
            print_success "✅ $service - ACTIVO"
        else
            print_warning "⚠️ $service - INACTIVO"
        fi
    done
}

# Mostrar información útil
show_info() {
    echo ""
    print_success "🎉 ¡Deployment completado!"
    echo ""
    print_status "📋 Información del deployment:"
    echo "  🌐 Dominio: https://${DOMAIN}"
    echo "  📂 Directorio: $(pwd)"
    echo "  🐳 Contenedores: $(docker ps --format '{{.Names}}' | wc -l) activos"
    echo ""
    print_status "🔧 Comandos útiles:"
    echo "  Ver logs: docker-compose logs -f"
    echo "  Reiniciar: docker-compose restart"
    echo "  Estado: docker-compose ps"
    echo "  Logs Nginx: sudo tail -f /var/log/nginx/error.log"
    echo "  SSL: sudo certbot certificates"
    echo ""
    print_status "🔗 URLs a verificar:"
    echo "  https://${DOMAIN}"
    echo "  https://${DOMAIN}/api/auth/health"
    echo "  https://${DOMAIN}/api/payment/health"
    echo "  https://${DOMAIN}/health"
}

# Función principal
main() {
    print_status "=== DEPLOYMENT LOCAL EN VPS ==="
    print_status "Dominio: ${DOMAIN}"
    print_status "Directorio: $(pwd)"
    echo ""
    
    # Verificar ubicación y archivos
    check_location
    
    # Verificar e instalar dependencias
    check_docker
    install_dependencies
    setup_firewall
    
    # Configurar entorno
    setup_environment
    
    # Detener servicios existentes
    stop_existing_services
    
    # Desplegar contenedores
    deploy_containers
    
    # Configurar Nginx
    setup_nginx
    
    # Configurar SSL
    setup_ssl
    
    # Verificar deployment
    verify_deployment
    
    # Mostrar información
    show_info
}

# Ejecutar función principal
main "$@"