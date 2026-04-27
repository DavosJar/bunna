.PHONY: docker-up docker-down docker-logs docker-clean build-test test

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f postgres

docker-clean:
	docker-compose down -v

build-test:
	go build -o bin/test ./cmd/test/main.go

test: build-test
	./bin/test
