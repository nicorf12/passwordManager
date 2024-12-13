# Variables
APP_NAME := password_manager.exe
SRC_FILE := main.go

# Comandos
BUILD_CMD := go build -o $(APP_NAME) $(SRC_FILE)
RUN_CMD := $(APP_NAME)
TEST_CMD := go test .\test

# Reglas
.PHONY: all build run test clean

all: build run

build:
	@echo "Construyendo el programa..."
	$(BUILD_CMD)
	@echo "Programa construido con éxito."

run: build
	@echo "Ejecutando el programa..."
	$(RUN_CMD)

test:
	@echo "Ejecutando tests..."
	$(TEST_CMD)

clean:
	@echo "Limpiando archivos generados..."
	rm -rf $(OUTPUT_DIR)
