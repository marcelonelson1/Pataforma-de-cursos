# 🚀 Guía de Deployment para bpnnneuq.co

## 📋 Resumen del Proyecto

**Dominio:** bpnnneuq.co  
**IP VPS:** 181.85.173.171  
**Arquitectura:** Microservicios con React Frontend  
**Base de Datos:** MariaDB  

### Servicios incluidos:
- **auth-user-service** (Puerto 8001) - Autenticación y usuarios
- **payment-service** (Puerto 8002) - Pagos con PayPal/MercadoPago  
- **course-service** (Puerto 8003) - Gestión de cursos
- **contact-service** (Puerto 8004) - Formulario de contacto
- **portfolio-service** (Puerto 8005) - Portfolio de proyectos
- **home-service** (Puerto 8006) - Imágenes del home
- **analytics-service** (Puerto 8007) - Estadísticas y métricas
- **frontend** (React) - Aplicación web

## 🛠️ Archivos de Deployment Creados

### 1. Script Principal de Deployment
```bash
./deploy-vps-production.sh
```
**Funciones:**
- ✅ Configuración automática del VPS
- ✅ Instalación de Docker y dependencias
- ✅ Configuración de MariaDB
- ✅ Setup de Nginx con SSL
- ✅ Deployment completo de todos los servicios

### 2. Configuración SSL y Nginx
```bash
./nginx-ssl-setup.sh
```
**Características:**
- ✅ SSL/TLS con Let's Encrypt
- ✅ HTTP/2 y compresión Gzip
- ✅ Rate limiting por API
- ✅ Headers de seguridad
- ✅ Optimizaciones de rendimiento

### 3. Docker Compose para Producción
```bash
docker-compose-production.yml
```
**Incluye:**
- ✅ MariaDB configurada
- ✅ Todos los microservicios
- ✅ Frontend optimizado
- ✅ Volúmenes persistentes
- ✅ Red interna segura

## 🚀 Pasos para el Deployment

### Paso 1: Preparación Local
```bash
# 1. Clonar/descargar el proyecto
cd /ruta/a/tu/proyecto

# 2. Dar permisos de ejecución
chmod +x deploy-vps-production.sh
chmod +x nginx-ssl-setup.sh

# 3. Verificar que tienes acceso SSH al VPS
ssh root@181.85.173.171
```

### Paso 2: Configurar Credenciales
Antes de ejecutar el deployment, necesitas configurar:

**Editar `.env.production` con tus credenciales reales:**
```bash
# Email (obligatorio para contacto y SSL)
SMTP_USER=tu_email@gmail.com
SMTP_PASS=tu_app_password_de_gmail
CONTACT_EMAIL=contacto@bpnnneuq.co

# PayPal (opcional)
PAYPAL_CLIENT_ID=tu_client_id
PAYPAL_CLIENT_SECRET=tu_client_secret

# MercadoPago (opcional)
MERCADOPAGO_ACCESS_TOKEN=tu_access_token
```

### Paso 3: Ejecutar Deployment
```bash
# Ejecutar el script principal
./deploy-vps-production.sh
```

El script automáticamente:
1. ✅ Verifica dependencias locales
2. ✅ Crea archivos de configuración
3. ✅ Copia archivos al VPS
4. ✅ Instala Docker, Nginx, Certbot
5. ✅ Configura MariaDB
6. ✅ Construye y ejecuta contenedores
7. ✅ Configura SSL con Let's Encrypt
8. ✅ Verifica que todo funcione

## 🔧 Comandos Útiles Post-Deployment

### Acceso al VPS
```bash
ssh root@181.85.173.171
cd /opt/microservicios
```

### Gestión de Contenedores
```bash
# Ver estado de servicios
docker-compose -f docker-compose-production.yml ps

# Ver logs de un servicio específico
docker-compose -f docker-compose-production.yml logs auth-service
docker-compose -f docker-compose-production.yml logs frontend

# Reiniciar un servicio
docker-compose -f docker-compose-production.yml restart payment-service

# Reiniciar todos los servicios
docker-compose -f docker-compose-production.yml restart

# Actualizar servicios (después de cambios de código)
docker-compose -f docker-compose-production.yml down
docker-compose -f docker-compose-production.yml build --no-cache
docker-compose -f docker-compose-production.yml up -d
```

### Gestión de Nginx
```bash
# Verificar configuración
nginx -t

# Recargar configuración
systemctl reload nginx

# Ver logs
tail -f /var/log/nginx/bpnnneuq.co_error.log
tail -f /var/log/nginx/bpnnneuq.co_access.log
```

### Gestión de SSL
```bash
# Ver certificados
certbot certificates

# Renovar certificados (automático, pero se puede hacer manual)
certbot renew

# Probar renovación
certbot renew --dry-run
```

### Base de Datos
```bash
# Acceder a MariaDB
docker exec -it microservices_mariadb mysql -u root -p

# Backup de base de datos
docker exec microservices_mariadb mysqldump -u root -p microservices_db > backup.sql

# Restaurar backup
docker exec -i microservices_mariadb mysql -u root -p microservices_db < backup.sql
```

## 🔍 Verificación del Deployment

### URLs a verificar:
- **Frontend:** https://bpnnneuq.co
- **APIs:**
  - https://bpnnneuq.co/api/auth/health
  - https://bpnnneuq.co/api/payment/health
  - https://bpnnneuq.co/api/course/health
  - https://bpnnneuq.co/api/contact/health
  - https://bpnnneuq.co/api/portfolio/health
  - https://bpnnneuq.co/api/home/health
  - https://bpnnneuq.co/api/analytics/health

### Verificación SSL:
```bash
# Probar SSL
curl -I https://bpnnneuq.co

# Verificar rating SSL
# Visitar: https://www.ssllabs.com/ssltest/analyze.html?d=bpnnneuq.co
```

## 🔒 Seguridad Implementada

### SSL/TLS
- ✅ Certificados Let's Encrypt
- ✅ TLS 1.2 y 1.3
- ✅ HSTS habilitado
- ✅ Renovación automática

### Headers de Seguridad
- ✅ X-Frame-Options: DENY
- ✅ X-Content-Type-Options: nosniff
- ✅ X-XSS-Protection: enabled
- ✅ Content-Security-Policy configurado
- ✅ Referrer-Policy: strict-origin

### Rate Limiting
- ✅ Auth API: 5 req/s
- ✅ APIs generales: 10 req/s
- ✅ Contact API: 2 req/s (anti-spam)

### Firewall
- ✅ Solo puertos 22, 80, 443 abiertos
- ✅ UFW habilitado

## 📊 Monitoreo y Logs

### Ubicación de Logs
```bash
# Logs de aplicación
/var/log/microservices/

# Logs de Nginx
/var/log/nginx/bpnnneuq.co_access.log
/var/log/nginx/bpnnneuq.co_error.log

# Logs de Docker
docker logs <container_name>
```

### Monitoreo de Recursos
```bash
# Ver uso de recursos
htop
docker stats

# Ver espacio en disco
df -h

# Ver uso de memoria
free -m
```

## 🔄 Actualizaciones y Mantenimiento

### Actualizar Código
```bash
# En tu máquina local, después de hacer cambios
./deploy-vps-production.sh
```

### Backups Recomendados
```bash
# Backup diario de base de datos (cron job)
0 2 * * * cd /opt/microservicios && docker exec microservices_mariadb mysqldump -u root -p'password' microservices_db > /opt/backups/db_$(date +\%Y\%m\%d).sql

# Backup de archivos subidos
0 3 * * * rsync -av /opt/microservicios/volumes/ /opt/backups/files/
```

### Actualizaciones del Sistema
```bash
# Actualizar VPS (hacer con cuidado)
apt update && apt upgrade -y

# Actualizar Docker images
docker-compose -f docker-compose-production.yml pull
docker-compose -f docker-compose-production.yml up -d
```

## 🆘 Troubleshooting

### Problemas Comunes

**1. Servicios no inician:**
```bash
# Verificar logs
docker-compose -f docker-compose-production.yml logs

# Verificar recursos
docker stats
df -h
```

**2. SSL no funciona:**
```bash
# Verificar certificados
certbot certificates

# Renovar manualmente
certbot renew --force-renewal
```

**3. Base de datos no conecta:**
```bash
# Verificar contenedor de MariaDB
docker exec -it microservices_mariadb mysql -u root -p

# Verificar variables de entorno
docker exec microservices_mariadb env | grep DB
```

**4. APIs no responden:**
```bash
# Verificar cada servicio
curl https://bpnnneuq.co/api/auth/health
curl https://bpnnneuq.co/api/payment/health
# etc...

# Verificar logs específicos
docker logs auth_service
```

## 📞 Soporte

Si encuentras problemas:

1. **Revisa los logs** mencionados arriba
2. **Verifica la configuración** de .env.production
3. **Comprueba el estado** de los contenedores
4. **Confirma la conectividad** de red

## 🎉 ¡Deployment Completado!

Tu aplicación debería estar funcionando en:
**https://bpnnneuq.co**

Con todos los microservicios, SSL, optimizaciones de rendimiento y medidas de seguridad implementadas.