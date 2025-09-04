#!/bin/bash

echo "🛑 Deteniendo todos los microservicios locales..."

# Función para matar proceso por PID
kill_process() {
    local pid_file=$1
    local service_name=$2
    
    if [ -f "$pid_file" ]; then
        local pid=$(cat $pid_file)
        if kill -0 $pid 2>/dev/null; then
            echo "🔴 Deteniendo $service_name (PID: $pid)..."
            kill $pid
            sleep 1
            # Force kill si es necesario
            if kill -0 $pid 2>/dev/null; then
                kill -9 $pid
            fi
        fi
        rm $pid_file
    fi
}

# Detener cada servicio
kill_process "auth-user-service.pid" "Auth Service"
kill_process "payment-service.pid" "Payment Service"
kill_process "course-service.pid" "Course Service"
kill_process "contact-service.pid" "Contact Service"
kill_process "portfolio-service.pid" "Portfolio Service"
kill_process "home-service.pid" "Home Service"
kill_process "frontend.pid" "Frontend"

# Limpiar archivos PID
rm -f *.pid microservices.pids

# Matar cualquier proceso Go o Node que pueda quedar
echo "🧹 Limpiando procesos restantes..."
pkill -f "auth-user-service-binary" 2>/dev/null
pkill -f "payment-service-binary" 2>/dev/null
pkill -f "course-service-binary" 2>/dev/null
pkill -f "contact-service-binary" 2>/dev/null
pkill -f "portfolio-service-binary" 2>/dev/null
pkill -f "home-service-binary" 2>/dev/null
pkill -f "npm start" 2>/dev/null

echo ""
echo "✅ Todos los microservicios han sido detenidos"
echo ""