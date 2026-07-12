.PHONY: compile compile-identidad compile-fincas \
        build build-agent build-ingestor build-identidad build-fincas build-all \
        push-agent push-ingestor push-identidad push-fincas push-all \
        deploy-infra deploy-test deploy-prod deploy-all \
        remove-infra remove-test remove-prod \
        ps-infra ps-test ps-prod \
        logs-infra logs-test logs-prod \
        dev-up dev-down dev-logs \
        up down

REGISTRY = 172.31.36.189:5000

# ─── Compile (binarios locales, compilador permanente) ───
compile: compile-identidad compile-fincas

compile-identidad:
	cd identidad && CGO_ENABLED=0 go build -o bin/server ./cmd/main.go

compile-fincas:
	cd fincas && CGO_ENABLED=0 go build -o bin/fincas ./cmd/main.go

# ─── Build imagenes Docker (desde binarios precompilados) ───
build: build-agent build-ingestor build-identidad build-fincas

build-agent:
	cd hardware-monitor-agent && docker build -t hardware-monitor-agent:release .

build-ingestor:
	cd servicio-monitoreo && docker build -t servicio-monitoreo:release .

build-identidad: compile-identidad
	cd identidad && docker build -t identidad:release .

build-fincas: compile-fincas
	cd fincas && docker build -t fincas:release .

build-all: compile build

# ─── Push ───
push-agent:
	docker tag hardware-monitor-agent:release $(REGISTRY)/hardware-monitor-agent:release
	docker push $(REGISTRY)/hardware-monitor-agent:release

push-ingestor:
	docker tag servicio-monitoreo:release $(REGISTRY)/servicio-monitoreo:release
	docker push $(REGISTRY)/servicio-monitoreo:release

push-identidad:
	docker tag identidad:release $(REGISTRY)/identidad:release
	docker push $(REGISTRY)/identidad:release

push-fincas:
	docker tag fincas:release $(REGISTRY)/fincas:release
	docker push $(REGISTRY)/fincas:release

push-all: push-agent push-ingestor push-identidad push-fincas

# ─── Deploy infra (monitoring stack) ───
deploy-infra:
	REGISTRY=$(REGISTRY) bash -c 'set -a; . soporte/infra/.env 2>/dev/null || true; set +a; docker stack deploy -c soporte/infra/compose.yml bunna-infra'

remove-infra:
	docker stack rm bunna-infra

ps-infra:
	docker stack ps bunna-infra

logs-infra:
	docker service logs --tail 50 -f bunna-infra_kafka bunna-infra_clickhouse bunna-infra_grafana

# ─── Deploy test ───
deploy-test:
	bash -c 'set -a; . soporte/test/.env; set +a; docker stack deploy -c soporte/test/compose.yml bunna-test'

remove-test:
	docker stack rm bunna-test

ps-test:
	docker stack ps bunna-test

logs-test:
	docker service logs --tail 50 -f bunna-test_kafka bunna-test_clickhouse bunna-test_grafana

# ─── Deploy prod (app services) ───
deploy-prod:
	docker stack deploy -c soporte/prod/compose.yml bunna-prod

remove-prod:
	docker stack rm bunna-prod

ps-prod:
	docker stack ps bunna-prod

# ─── Deploy all (infra + prod apps) ───
deploy-all: deploy-infra deploy-prod

# ─── Desarrollo Local (Unificado) ───
dev-up:
	docker compose -f docker-compose.dev.yml up -d --build

dev-down:
	docker compose -f docker-compose.dev.yml down

dev-logs:
	docker compose -f docker-compose.dev.yml logs -f

# ─── Up / Down: compile + todo en Docker ───
up: compile
	docker compose -f docker-compose.dev.yml up -d --build
	docker compose -f docker-compose.dev.yml logs -f

down:
	docker compose -f docker-compose.dev.yml down
