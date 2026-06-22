.PHONY: deploy-dev deploy-prod remove-dev remove-prod ps-dev ps-prod logs-dev logs-prod

# ─── Dev ───
deploy-dev:
	docker stack deploy -c soporte/dev/compose.yml bunna-dev

remove-dev:
	docker stack rm bunna-dev

ps-dev:
	docker stack ps bunna-dev

logs-dev:
	docker service logs --tail 50 -f bunna-dev_kafka bunna-dev_clickhouse bunna-dev_grafana

# ─── Prod ───
deploy-prod:
	docker stack deploy -c soporte/prod/compose.yml bunna-prod

remove-prod:
	docker stack rm bunna-prod

ps-prod:
	docker stack ps bunna-prod

logs-prod:
	docker service logs --tail 50 -f bunna-prod_kafka bunna-prod_clickhouse bunna-prod_grafana
