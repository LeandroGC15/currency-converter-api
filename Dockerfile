# Etapa 1: Construir la aplicación Go
FROM golang:1.18-alpine AS builder

WORKDIR /app

# Copiar los archivos de módulos de Go y descargar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código fuente
COPY . .

# Compilar la aplicación Go
# -buildvcs=false es para evitar problemas con git en el build
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/conversor -buildvcs=false main.go

# Etapa 2: Crear la imagen final con Python y la aplicación Go
FROM python:3.9-slim

WORKDIR /app

# Instalar las dependencias de Python
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copiar el script de Python
COPY get_dolar.py .

# Copiar el ejecutable de Go desde la etapa de construcción
COPY --from=builder /go/bin/conversor .

# Exponer el puerto que la aplicación usará (Render lo gestiona, pero es buena práctica)
EXPOSE 10000

# Comando para iniciar la aplicación
CMD ["./conversor"]
