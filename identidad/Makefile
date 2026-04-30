.PHONY: docker-up docker-down docker-logs docker-shell docker-clean build-test test

include .env

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f postgres

docker-shell:
	docker-compose exec postgres psql -U ${POSTGRES_USER} -d ${POSTGRES_DB}

docker-clean:
	docker-compose down -v

build-test:
	go build -o bin/test ./cmd/test/main.go

test: build-test
	./bin/test
