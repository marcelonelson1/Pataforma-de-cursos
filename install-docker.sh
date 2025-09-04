#!/bin/bash

echo "🐳 Instalando Docker y Docker Compose en Kali Linux..."

# Actualizar repositorios
echo "📦 Actualizando repositorios..."
sudo apt update

# Instalar dependencias
echo "🔧 Instalando dependencias..."
sudo apt install -y apt-transport-https ca-certificates curl gnupg lsb-release

# Agregar clave GPG oficial de Docker
echo "🔑 Agregando clave GPG de Docker..."
curl -fsSL https://download.docker.com/linux/debian/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

# Agregar repositorio de Docker
echo "📋 Agregando repositorio de Docker..."
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/debian bullseye stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Actualizar repositorios nuevamente
echo "🔄 Actualizando repositorios..."
sudo apt update

# Instalar Docker Engine
echo "🐳 Instalando Docker Engine..."
sudo apt install -y docker-ce docker-ce-cli containerd.io

# Instalar Docker Compose (método oficial)
echo "🔗 Instalando Docker Compose..."
DOCKER_COMPOSE_VERSION=$(curl -s https://api.github.com/repos/docker/compose/releases/latest | grep -Po '"tag_name": "\K.*?(?=")')
sudo curl -L "https://github.com/docker/compose/releases/download/${DOCKER_COMPOSE_VERSION}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Agregar usuario actual al grupo docker
echo "👤 Agregando usuario al grupo docker..."
sudo usermod -aG docker $USER

# Iniciar y habilitar Docker
echo "🚀 Iniciando Docker..."
sudo systemctl start docker
sudo systemctl enable docker

# Verificar instalación
echo "✅ Verificando instalación..."
docker --version
docker-compose --version

echo ""
echo "✅ ¡Docker y Docker Compose instalados correctamente!"
echo ""
echo "⚠️  IMPORTANTE: Reinicia tu sesión o ejecuta:"
echo "   newgrp docker"
echo ""
echo "🔄 Después ejecuta:"
echo "   ./start-microservices.sh"
echo ""