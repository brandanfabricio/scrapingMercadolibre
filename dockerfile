# Primera fase: Build del binario
FROM golang:1.22.5-alpine AS build

# Instalar dependencias necesarias para la compilación
RUN apk --no-cache add \
    git upx ca-certificates

# Definir el directorio de trabajo
WORKDIR /app

# Copiar los archivos de dependencias
COPY ["go.mod", "go.sum", "./"]

# Descargar los módulos necesarios
RUN go mod download

# Copiar el resto del código
COPY . .

# Construir el binario
RUN go build -ldflags="-s -w" -o app .

# Comprimir el binario con UPX (opcional)
RUN upx app || true  # `|| true` para que no falle si `upx` no está disponible

# Segunda fase: Imagen más liviana para ejecución
FROM alpine:3.18

# Instalar solo las dependencias necesarias sin repositorios Edge
RUN apk --no-cache add \
    chromium \
    harfbuzz \
    nss \
    freetype \
    ttf-freefont \
    alsa-lib

# Copiar el binario desde la primera fase
COPY --from=build /app/app /app/app

# Definir el directorio de trabajo
WORKDIR /app

# Ejecutar el binario al iniciar el contenedor
CMD ["./app"]
