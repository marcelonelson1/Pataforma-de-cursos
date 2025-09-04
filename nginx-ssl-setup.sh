#!/bin/bash

# Script para configurar SSL y optimizar Nginx para bpnnneuq.co
set -e

DOMAIN="bpnnneuq.co"
EMAIL="admin@${DOMAIN}"

echo "🔒 Configurando SSL y optimizando Nginx para ${DOMAIN}..."

# Crear configuración optimizada de Nginx con SSL
create_nginx_ssl_config() {
    cat > /etc/nginx/sites-available/${DOMAIN} << EOF
# Configuración para ${DOMAIN}
server {
    listen 80;
    server_name ${DOMAIN} www.${DOMAIN};
    
    # Redirección a HTTPS
    return 301 https://\$server_name\$request_uri;
}

server {
    listen 443 ssl http2;
    server_name ${DOMAIN} www.${DOMAIN};

    # Certificados SSL (se configurarán con Certbot)
    ssl_certificate /etc/letsencrypt/live/${DOMAIN}/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/${DOMAIN}/privkey.pem;
    ssl_trusted_certificate /etc/letsencrypt/live/${DOMAIN}/chain.pem;

    # Configuración SSL moderna
    ssl_session_timeout 1d;
    ssl_session_cache shared:SSL:50m;
    ssl_session_tickets off;
    
    # Protocolos y cifrados modernos
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;

    # HSTS (HTTP Strict Transport Security)
    add_header Strict-Transport-Security "max-age=63072000; includeSubDomains; preload" always;
    
    # Headers de seguridad
    add_header X-Frame-Options DENY always;
    add_header X-Content-Type-Options nosniff always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
    add_header Content-Security-Policy "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self' https:; media-src 'self'; object-src 'none'; frame-src 'none';" always;

    # OCSP Stapling
    ssl_stapling on;
    ssl_stapling_verify on;
    resolver 8.8.8.8 8.8.4.4 valid=300s;
    resolver_timeout 5s;

    # Configuración de archivos
    client_max_body_size 100M;
    client_body_timeout 60s;
    client_header_timeout 60s;
    keepalive_timeout 65s;
    send_timeout 60s;

    # Compresión Gzip
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_comp_level 6;
    gzip_types
        text/plain
        text/css
        text/xml
        text/javascript
        application/json
        application/javascript
        application/xml+rss
        application/atom+xml
        image/svg+xml;

    # Rate limiting
    limit_req_zone \$binary_remote_addr zone=auth:10m rate=5r/s;
    limit_req_zone \$binary_remote_addr zone=api:10m rate=10r/s;
    limit_req_zone \$binary_remote_addr zone=contact:10m rate=2r/s;

    # Logs
    access_log /var/log/nginx/${DOMAIN}_access.log;
    error_log /var/log/nginx/${DOMAIN}_error.log;

    # API Endpoints con rate limiting
    location /api/auth/ {
        limit_req zone=auth burst=10 nodelay;
        proxy_pass http://localhost:8001/;
        include /etc/nginx/proxy_params;
    }

    location /api/payment/ {
        limit_req zone=api burst=5 nodelay;
        proxy_pass http://localhost:8002/;
        include /etc/nginx/proxy_params;
    }

    location /api/course/ {
        limit_req zone=api burst=20 nodelay;
        proxy_pass http://localhost:8003/;
        include /etc/nginx/proxy_params;
    }

    location /api/contact/ {
        limit_req zone=contact burst=3 nodelay;
        proxy_pass http://localhost:8004/;
        include /etc/nginx/proxy_params;
    }

    location /api/portfolio/ {
        limit_req zone=api burst=15 nodelay;
        proxy_pass http://localhost:8005/;
        include /etc/nginx/proxy_params;
    }

    location /api/home/ {
        limit_req zone=api burst=15 nodelay;
        proxy_pass http://localhost:8006/;
        include /etc/nginx/proxy_params;
    }

    location /api/analytics/ {
        limit_req zone=api burst=10 nodelay;
        proxy_pass http://localhost:8007/;
        include /etc/nginx/proxy_params;
    }

    # Archivos estáticos con cache
    location ~* \.(jpg|jpeg|png|gif|ico|css|js|woff|woff2|ttf|svg)\$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
        add_header Vary Accept-Encoding;
        access_log off;
    }

    # Frontend React
    location / {
        proxy_pass http://localhost:3000;
        include /etc/nginx/proxy_params;
        
        # Soporte para client-side routing
        try_files \$uri \$uri/ @fallback;
    }

    location @fallback {
        proxy_pass http://localhost:3000;
        include /etc/nginx/proxy_params;
    }

    # Health check
    location /health {
        access_log off;
        return 200 "healthy\\n";
        add_header Content-Type text/plain;
    }

    # Bloquear acceso a archivos sensibles
    location ~ /\. {
        deny all;
        access_log off;
        log_not_found off;
    }

    location ~ /\.(env|git|svn) {
        deny all;
        access_log off;
        log_not_found off;
    }
}
EOF

    echo "✅ Configuración de Nginx creada"
}

# Crear parámetros de proxy optimizados
create_proxy_params() {
    cat > /etc/nginx/proxy_params << EOF
proxy_set_header Host \$http_host;
proxy_set_header X-Real-IP \$remote_addr;
proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
proxy_set_header X-Forwarded-Proto \$scheme;
proxy_set_header X-Forwarded-Host \$server_name;
proxy_set_header X-Forwarded-Port \$server_port;

# Timeouts
proxy_connect_timeout 60s;
proxy_send_timeout 60s;
proxy_read_timeout 60s;

# Buffer settings
proxy_buffering on;
proxy_buffer_size 128k;
proxy_buffers 4 256k;
proxy_busy_buffers_size 256k;

# Headers
proxy_http_version 1.1;
proxy_set_header Upgrade \$http_upgrade;
proxy_set_header Connection 'upgrade';
proxy_cache_bypass \$http_upgrade;
EOF

    echo "✅ Parámetros de proxy creados"
}

# Optimizar configuración general de Nginx
optimize_nginx() {
    # Backup de configuración original
    cp /etc/nginx/nginx.conf /etc/nginx/nginx.conf.backup

    cat > /etc/nginx/nginx.conf << EOF
user www-data;
worker_processes auto;
pid /run/nginx.pid;

events {
    worker_connections 2048;
    use epoll;
    multi_accept on;
}

http {
    # Tipos MIME
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    # Logs
    log_format main '\$remote_addr - \$remote_user [\$time_local] "\$request" '
                    '\$status \$body_bytes_sent "\$http_referer" '
                    '"\$http_user_agent" "\$http_x_forwarded_for" '
                    'rt=\$request_time uct="\$upstream_connect_time" '
                    'uht="\$upstream_header_time" urt="\$upstream_response_time"';

    access_log /var/log/nginx/access.log main;
    error_log /var/log/nginx/error.log warn;

    # Optimizaciones de rendimiento
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    types_hash_max_size 2048;
    server_tokens off;

    # Buffer sizes
    client_body_buffer_size 128k;
    client_max_body_size 100m;
    client_header_buffer_size 1k;
    large_client_header_buffers 4 4k;
    output_buffers 1 32k;
    postpone_output 1460;

    # Timeouts
    client_header_timeout 3m;
    client_body_timeout 3m;
    send_timeout 3m;

    # Gzip
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_comp_level 6;
    gzip_types
        text/plain
        text/css
        text/xml
        text/javascript
        application/json
        application/javascript
        application/xml+rss
        application/atom+xml
        image/svg+xml
        application/x-font-ttf
        font/opentype;

    # Rate limiting zones
    limit_req_zone \$binary_remote_addr zone=global:10m rate=10r/s;

    # Virtual Host Configs
    include /etc/nginx/conf.d/*.conf;
    include /etc/nginx/sites-enabled/*;
}
EOF

    echo "✅ Configuración general de Nginx optimizada"
}

# Configurar SSL con Certbot
setup_ssl() {
    echo "🔒 Configurando SSL con Let's Encrypt..."
    
    # Instalar Certbot si no está instalado
    if ! command -v certbot &> /dev/null; then
        apt update
        apt install -y certbot python3-certbot-nginx
    fi

    # Habilitar el sitio
    ln -sf /etc/nginx/sites-available/${DOMAIN} /etc/nginx/sites-enabled/
    rm -f /etc/nginx/sites-enabled/default

    # Probar configuración
    nginx -t

    # Recargar Nginx
    systemctl reload nginx

    # Obtener certificado SSL
    certbot --nginx -d ${DOMAIN} -d www.${DOMAIN} \
        --non-interactive \
        --agree-tos \
        --email ${EMAIL} \
        --redirect

    # Configurar renovación automática
    systemctl enable certbot.timer
    systemctl start certbot.timer

    echo "✅ SSL configurado correctamente"
}

# Configurar logrotate para logs de Nginx
setup_logrotate() {
    cat > /etc/logrotate.d/nginx-${DOMAIN} << EOF
/var/log/nginx/${DOMAIN}_*.log {
    daily
    missingok
    rotate 52
    compress
    delaycompress
    notifempty
    create 644 www-data adm
    postrotate
        systemctl reload nginx
    endscript
}
EOF

    echo "✅ Logrotate configurado"
}

# Función principal
main() {
    echo "🚀 Configurando SSL y optimizando Nginx..."
    
    # Verificar que se está ejecutando como root
    if [[ $EUID -ne 0 ]]; then
        echo "❌ Este script debe ejecutarse como root"
        exit 1
    fi

    # Crear configuraciones
    create_nginx_ssl_config
    create_proxy_params
    optimize_nginx
    setup_logrotate

    # Configurar SSL
    setup_ssl

    # Verificar configuración
    nginx -t
    systemctl restart nginx

    echo "✅ ¡SSL y Nginx configurados correctamente!"
    echo "🌐 Tu sitio está disponible en: https://${DOMAIN}"
    
    # Mostrar estado de certificados
    certbot certificates

    echo ""
    echo "📋 Comandos útiles:"
    echo "  - Verificar SSL: curl -I https://${DOMAIN}"
    echo "  - Renovar certificados: certbot renew"
    echo "  - Ver logs: tail -f /var/log/nginx/${DOMAIN}_error.log"
    echo "  - Probar configuración: nginx -t"
}

# Ejecutar si se llama directamente
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi