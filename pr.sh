#!/bin/ash

# Actualiza la lista de paquetes
echo "Actualizando repositorios..."
apk update

# Habilita el repositorio edge para obtener paquetes más recientes (opcional)
echo "Habilitando repositorios edge..."
echo "http://dl-cdn.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories
echo "http://dl-cdn.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories

# Actualiza nuevamente para incluir los nuevos repositorios
apk update

# Instala Chromium y las dependencias necesarias
echo "Instalando Chromium y dependencias..."
apk --no-cache add ca-certificates chromium nss freetype freetype-dev harfbuzz ttf-freefont

# Verifica la instalación de Chromium
if which chromium-browser > /dev/null 2>&1; then
    echo "Chromium se ha instalado correctamente."
else
    echo "Error: Chromium no se ha instalado correctamente."
fi
