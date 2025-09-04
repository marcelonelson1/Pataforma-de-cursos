#!/bin/bash

# Script mejorado para compilar y iniciar microservicios localmente
set -e  # Salir si hay errores

echo "🚀 Iniciando microservicios en modo desarrollo local..."
echo "📅 $(date)"
echo ""

# Colores para output
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() { echo -e "${GREEN}✅ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠️  $1${NC}"; }
print_error() { echo -e "${RED}❌ $1${NC}"; }
print_info() { echo -e "${BLUE}ℹ️  $1${NC}"; }

# Verificar dependencias del sistema
check_dependencies() {
    print_info "Verificando dependencias..."
    
    if ! command -v go &> /dev/null; then
        print_error "Go no está instalado. Instalar con: sudo apt install golang-go"
        exit 1
    fi
    
    if ! command -v node &> /dev/null; then
        print_error "Node.js no está instalado. Instalar con: sudo apt install nodejs npm"
        exit 1
    fi
    
    print_status "Dependencias verificadas"
}

# Limpiar procesos anteriores
cleanup_previous() {
    print_info "Limpiando procesos anteriores..."
    
    # Limpiar PIDs anteriores
    if [ -f "microservices.pids" ]; then
        while read -r pid; do
            if [ -n "$pid" ] && ps -p "$pid" > /dev/null 2>&1; then
                kill "$pid" 2>/dev/null || true
            fi
        done < microservices.pids
    fi
    
    rm -f *.pid microservices.pids
    pkill -f "service-binary" 2>/dev/null || true
    sleep 2
    print_status "Limpieza completada"
}

# Función mejorada para servicios Go
start_go_service() {
    local service_name=$1
    local port=$2
    local description=$3
    
    print_info "Compilando $description ($service_name)..."
    
    if [ ! -d "$service_name" ]; then
        print_error "Directorio $service_name no encontrado"
        return 1
    fi
    
    cd "$service_name"
    
    # Limpiar compilaciones anteriores
    rm -f "${service_name}-binary" "${service_name}" "main"
    
    # Actualizar dependencias y compilar
    print_info "  📦 Actualizando dependencias..."
    go mod tidy
    
    print_info "  🔨 Compilando..."
    if ! go build -ldflags="-s -w" -o "${service_name}-binary" .; then
        print_error "  Error compilando $service_name"
        cd ..
        return 1
    fi
    
    chmod +x "${service_name}-binary"
    
    # Ejecutar en background
    print_info "  🚀 Iniciando $service_name en puerto $port..."
    ./"${service_name}-binary" > app.log 2>&1 &
    local service_pid=$!
    
    # Verificar inicio
    sleep 1
    if ! ps -p $service_pid > /dev/null; then
        print_error "  Error iniciando $service_name"
        cd ..
        return 1
    fi
    
    echo $service_pid > "../${service_name}.pid"
    print_status "  $description iniciado (PID: $service_pid)"
    
    cd ..
    return 0
}

# Función para frontend
start_frontend() {
    print_info "Preparando Frontend React..."
    
    cd frontend
    
    # Instalar dependencias si es necesario
    if [ ! -d "node_modules" ] || [ "package.json" -nt "node_modules" ]; then
        print_info "  📦 Instalando dependencias..."
        npm install
    fi
    
    print_info "  🚀 Iniciando servidor React..."
    BROWSER=none npm start > ../frontend.log 2>&1 &
    local frontend_pid=$!
    
    sleep 2
    if ! ps -p $frontend_pid > /dev/null; then
        print_error "  Error iniciando frontend"
        cd ..
        return 1
    fi
    
    echo $frontend_pid > "../frontend.pid"
    print_status "  Frontend iniciado (PID: $frontend_pid)"
    cd ..
}

# MAIN EXECUTION
main() {
    echo "╔═══════════════════════════════════════╗"
    echo "║     MICROSERVICIOS LOCAL STARTER      ║"
    echo "╚═══════════════════════════════════════╝"
    echo ""
    
    check_dependencies
    cleanup_previous
    touch microservices.pids
    
    print_info "🏗️  INICIANDO SERVICIOS BACKEND..."
    
    # Lista de servicios
    local services=(
        "auth-user-service:8001:Auth Service"
        "payment-service:8002:Payment Service" 
        "course-service:8003:Course Service"
        "contact-service:8004:Contact Service"
        "portfolio-service:8005:Portfolio Service"
        "home-service:8006:Home Service"
        "analytics-service:8007:Analytics Service"
    )
    
    local failed_services=()
    
    for service_info in "${services[@]}"; do
        IFS=':' read -r service_name port description <<< "$service_info"
        if ! start_go_service "$service_name" "$port" "$description"; then
            failed_services+=("$description")
        fi
    done
    
    echo ""
    print_info "🌐 INICIANDO FRONTEND..."
    if ! start_frontend; then
        failed_services+=("Frontend React")
    fi
    
    # Recopilar PIDs
    cat *.pid > microservices.pids 2>/dev/null || true
    
    echo ""
    print_info "📊 REPORTE FINAL"
    
    if [ ${#failed_services[@]} -eq 0 ]; then
        print_status "¡Todos los servicios iniciados exitosamente! 🎉"
    else
        print_warning "Servicios que fallaron:"
        for service in "${failed_services[@]}"; do
            echo "  ❌ $service"
        done
    fi
    
    echo ""
    echo "🌐 URLs DISPONIBLES:"
    echo "   🖥️  Frontend:         http://localhost:3000"
    echo "   🔐 Auth Service:     http://localhost:8001"
    echo "   💳 Payment Service:  http://localhost:8002"
    echo "   📚 Course Service:   http://localhost:8003"
    echo "   📧 Contact Service:  http://localhost:8004"
    echo "   🎨 Portfolio Service: http://localhost:8005"
    echo "   🏠 Home Service:     http://localhost:8006"
    echo "   📊 Analytics Service: http://localhost:8007"
    echo ""
    echo "📋 COMANDOS ÚTILES:"
    echo "   🛑 Detener todo:     ./stop-local.sh"
    echo "   👀 Ver procesos:     ps aux | grep -E 'service-binary|node'"
    echo "   📄 Ver logs:         tail -f service-name/app.log"
    echo "   🔄 Reiniciar:        ./stop-local.sh && ./start-local.sh"
    echo ""
    print_status "¡Sistema listo! 🚀"
}

main "$@"