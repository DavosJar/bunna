.PHONY: docker-up docker-down docker-logs docker-shell docker-clean build test run run-test

include .env
export

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f postgres

docker-shell:
	docker compose exec postgres psql -U ${POSTGRES_USER} -d ${POSTGRES_DB}

docker-clean:
	docker compose down -v

build:
	go build -o bin/fincas ./cmd

test:
	go test ./...

run:
	go run ./cmd

run-test:
	go run ./cmd/test
