.PHONY: build run test clean help

help:
	@echo "Image Service Makefile"
	@echo ""
	@echo "Comandos disponibles:"
	@echo "  make up       - Levanta MQTT + compila + ejecuta servicio"
	@echo "  make down     - Apaga MQTT broker"
	@echo "  make build    - Compilar el servicio"
	@echo "  make run      - Ejecutar con directorio ./images"
	@echo "  make clean    - Limpiar binarios"
	@echo "  make mqtt-up  - Levanta MQTT broker"
	@echo "  make mqtt-down - Apaga MQTT broker"
	@echo ""

build:
	go build -o image-service

run: build
	./image-service -dir ./images -broker tcp://localhost:1883 -workers 4

test:
	go test -v ./...

clean:
	rm -f image-service
	go clean

.PHONY: deps
deps:
	go get github.com/eclipse/paho.mqtt.golang
	go get github.com/nfnt/resize
	go mod tidy

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: mqtt-up
mqtt-up:
	sudo docker-compose up -d mosquitto
	@echo "MQTT broker levantado en localhost:1883"

.PHONY: mqtt-down
mqtt-down:
	sudo docker-compose down

.PHONY: mqtt-logs
mqtt-logs:
	docker-compose logs -f mosquitto

.PHONY: up
up: mqtt-up build
	@echo "Esperando que MQTT esté listo..."
	@sleep 2
	@echo "Iniciando servicio (mantiene MQTT activo)..."
	./image-service -dir ./images -broker tcp://localhost:1883 -workers 4

.PHONY: down
down: mqtt-down
	@echo "Servicios detenidos"
