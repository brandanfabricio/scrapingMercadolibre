# Nombre del ejecutable
EXEC = webScraping.exe

# Directorios
BUILD_DIR = build


# Comando para compilar el proyecto en Go
BUILD_CMD = go build -o $(BUILD_DIR)/$(EXEC)


# Meta principal
all:  
	@echo "Compilando el proyecto..."
	$(BUILD_CMD)
	@echo "Build completado y archivos copiados."

