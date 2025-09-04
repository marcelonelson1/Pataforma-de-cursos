#!/bin/bash

echo "=== Building Frontend for Production ==="

# Ir al directorio del frontend
cd /home/kali/microservicios/frontend

# Verificar que existe package.json
if [ ! -f "package.json" ]; then
    echo "❌ Error: No se encontró package.json en el directorio frontend"
    exit 1
fi

# Copiar archivo .env.production como .env para el build
echo "📋 Configurando variables de entorno para producción..."
cp .env.production .env

# Instalar dependencias si no existen
if [ ! -d "node_modules" ]; then
    echo "📦 Instalando dependencias..."
    npm install
fi

# Limpiar build anterior
echo "🧹 Limpiando build anterior..."
rm -rf build/

# Crear build de producción
echo "🏗️ Creando build de producción..."
npm run build

# Verificar que el build se creó correctamente
if [ -d "build" ]; then
    echo "✅ Build de producción creado exitosamente!"
    echo "📁 Archivos generados en: /home/kali/microservicios/frontend/build/"
    
    # Mostrar tamaño del build
    echo "📊 Tamaño del build:"
    du -sh build/
    
    # Listar archivos principales
    echo "📄 Archivos principales generados:"
    ls -la build/static/js/*.js 2>/dev/null || echo "No se encontraron archivos JS"
    ls -la build/static/css/*.css 2>/dev/null || echo "No se encontraron archivos CSS"
    
else
    echo "❌ Error: No se pudo crear el build de producción"
    exit 1
fi

echo ""
echo "🎉 ¡Build de producción listo!"
echo ""
echo "Próximos pasos:"
echo "1. Configurar Apache: sudo cp apache-microservices.conf /etc/apache2/sites-available/bpnnneuq.co.conf"
echo "2. Habilitar módulos: sudo a2enmod proxy proxy_http headers rewrite expires"
echo "3. Habilitar sitio: sudo a2ensite bpnnneuq.co.conf"
echo "4. Reiniciar Apache: sudo systemctl reload apache2"
echo "5. Configurar SSL: sudo certbot --apache -d bpnnneuq.co -d www.bpnnneuq.co"