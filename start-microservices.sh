#!/bin/bash

# Script para iniciar todos los microservicios
echo "🚀 Iniciando todos los microservicios..."

# Detectar si usar docker-compose o docker compose
DOCKER_COMPOSE="docker-compose"
if ! command -v docker-compose &> /dev/null; then
    if command -v docker &> /dev/null && docker compose version &> /dev/null; then
        DOCKER_COMPOSE="docker compose"
        echo "📦 Usando Docker Compose v2 (docker compose)"
    else
        echo "❌ Error: Docker Compose no está instalado"
        echo "🔧 Instalar con: sudo apt install docker-compose"
        echo "   O usar Docker Desktop que incluye Compose v2"
        exit 1
    fi
else
    echo "📦 Usando Docker Compose v1 (docker-compose)"
fi

# Verificar que Docker está corriendo
if ! docker info &> /dev/null; then
    echo "❌ Error: Docker no está corriendo"
    echo "🔧 Iniciar con: sudo systemctl start docker"
    exit 1
fi

# Crear archivo .env si no existe
if [ ! -f .env ]; then
    echo "📝 Creando archivo .env desde .env.docker..."
    cp .env.docker .env
    echo "⚠️  IMPORTANTE: Edita el archivo .env con tus credenciales reales antes de continuar"
    echo "⚠️  Presiona Enter para continuar o Ctrl+C para cancelar..."
    read
fi

# Detener contenedores existentes
echo "🛑 Deteniendo contenedores existentes..."
$DOCKER_COMPOSE down

# Construir e iniciar servicios
echo "🏗️  Construyendo e iniciando servicios..."
$DOCKER_COMPOSE up --build -d

# Mostrar estado de los servicios
echo "📊 Estado de los servicios:"
$DOCKER_COMPOSE ps

echo ""
echo "✅ Microservicios iniciados!"
echo ""
echo "🌐 URLs disponibles:"
echo "   Frontend:        http://localhost:3000"
echo "   Auth Service:    http://localhost:8001"
echo "   Payment Service: http://localhost:8002"
echo "   Course Service:  http://localhost:8003"
echo "   Contact Service: http://localhost:8004"
echo "   Portfolio Service: http://localhost:8005"
echo "   Home Service:    http://localhost:8006"
echo ""
echo "📋 Comandos útiles:"
echo "   Ver logs:        docker-compose logs -f [servicio]"
echo "   Detener todo:    docker-compose down"
echo "   Reiniciar:       docker-compose restart [servicio]"
echo ""