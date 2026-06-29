.PHONY: docker-up docker-down docker-logs docker-shell docker-clean build-test test seed init-superuser run build-prod push-prod deploy-prod remove-prod ps-prod logs-prod

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

init-superuser:
	go run ./cmd/init_superuser

run:
	go run ./cmd

# ─── Targets para desarrollo local (BD dockerizada en puerto 5433) ───
run-dev:
	env $$(cat .env.develop | grep -v '^#' | xargs) go run ./cmd

seed-dev:
	env $$(cat .env.develop | grep -v '^#' | xargs) go run ./cmd/seed

# ─── Producción (Swarm) ───
build-prod:
	docker build -t identidad:release .

push-prod:
	docker tag identidad:release 172.31.36.189:5000/identidad:release
	docker push 172.31.36.189:5000/identidad:release

deploy-prod:
	docker stack deploy -c compose.prod.yml bunna-prod

remove-prod:
	docker stack rm bunna-prod

ps-prod:
	docker stack ps bunna-prod

logs-prod:
	docker service logs --tail 50 -f bunna-prod_identidad
