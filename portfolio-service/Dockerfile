# Dockerfile para Portfolio Service
FROM golang:1.21-alpine AS builder

# Instalar dependencias necesarias
RUN apk add --no-cache git

# Establecer directorio de trabajo
WORKDIR /app

# Copiar archivos de modulos y descargar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar codigo fuente
COPY . .

# Compilar la aplicacion
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Etapa final
FROM alpine:latest

# Instalar ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates

# Crear directorio de trabajo
WORKDIR /root/

# Copiar binario desde el builder
COPY --from=builder /app/main .

# Crear directorios necesarios
RUN mkdir -p ./static/portfolio

# Exponer puerto
EXPOSE 8005

# Comando para ejecutar la aplicacion
CMD ["./main"]