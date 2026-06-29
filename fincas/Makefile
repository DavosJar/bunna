include .env
export

PG  := $(shell docker compose config --services 2>/dev/null | grep -i postgres || echo postgres)
KF  := $(shell docker compose config --services 2>/dev/null | grep -i kafka || echo kafka)
.PHONY: up down clean logs ps logs-kafka logs-postgres build run consume topics infra up-prod down-prod logs-prod

up: build
	docker compose up --build

down:
	docker compose down

clean:
	docker compose down -v

logs:
	docker compose logs -f

logs-kafka:
	docker compose logs -f $(KF)

logs-postgres:
	docker compose logs -f $(PG)

ps:
	docker compose ps

up-prod: build
	docker compose -f docker-compose.prod.yml up -d --build

down-prod:
	docker compose -f docker-compose.prod.yml down

logs-prod:
	docker compose -f docker-compose.prod.yml logs -f

build:
	go build -o bin/fincas ./cmd/main.go

infra:
	docker compose up -d --no-build $(PG) $(KF)

run: infra
	go run ./cmd --env $(or $(ENV),dev)

consume:
	docker compose exec $(KF) kafka-console-consumer --bootstrap-server localhost:9092 --topic $(or $(TOPIC),dev.permisos) --from-beginning

topics:
	docker compose exec $(KF) kafka-topics --bootstrap-server localhost:9092 --list
