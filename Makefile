.PHONY: build-agent build-ingestor build-all push-agent push-ingestor push-all \
        deploy-test deploy-prod remove-test remove-prod ps-test ps-prod logs-test logs-prod

REGISTRY = 172.31.36.189:5000

# ─── Build ───
build-agent:
	cd hardware-monitor-agent && docker build -t hardware-monitor-agent:release .

build-ingestor:
	cd servicio-monitoreo && docker build -t servicio-monitoreo:release .

build-all: build-agent build-ingestor

# ─── Push ───
push-agent:
	docker tag hardware-monitor-agent:release $(REGISTRY)/hardware-monitor-agent:release
	docker push $(REGISTRY)/hardware-monitor-agent:release

push-ingestor:
	docker tag servicio-monitoreo:release $(REGISTRY)/servicio-monitoreo:release
	docker push $(REGISTRY)/servicio-monitoreo:release

push-all: push-agent push-ingestor

# ─── Deploy test ───
deploy-test:
	bash -c 'set -a; . soporte/test/.env; set +a; docker stack deploy -c soporte/test/compose.yml bunna-test'

remove-test:
	docker stack rm bunna-test

ps-test:
	docker stack ps bunna-test

logs-test:
	docker service logs --tail 50 -f bunna-test_kafka bunna-test_clickhouse bunna-test_grafana

# ─── Deploy prod ───
deploy-prod:
	bash -c 'set -a; . soporte/prod/.env; set +a; docker stack deploy -c soporte/prod/compose.yml bunna-prod'

remove-prod:
	docker stack rm bunna-prod

ps-prod:
	docker stack ps bunna-prod

logs-prod:
	docker service logs --tail 50 -f bunna-prod_kafka bunna-prod_clickhouse bunna-prod_grafana
