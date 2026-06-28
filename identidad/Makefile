.PHONY: docker-up docker-down docker-logs docker-shell docker-clean seed run

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


seed:
	go run ./cmd/seed

run:
	go run ./cmd

# ─── Targets para desarrollo local (BD dockerizada en puerto 5433) ───
run-dev:
	env $$(cat .env.develop | grep -v '^#' | xargs) go run ./cmd

seed-dev:
	env $$(cat .env.develop | grep -v '^#' | xargs) go run ./cmd/seed
