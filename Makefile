.PHONY: deploy-test deploy-prod remove-test remove-prod ps-test ps-prod logs-test logs-prod

# ─── Test ───
deploy-test:
	docker stack deploy -c soporte/test/compose.yml bunna-test

remove-test:
	docker stack rm bunna-test

ps-test:
	docker stack ps bunna-test

logs-test:
	docker service logs --tail 50 -f bunna-test_kafka bunna-test_clickhouse bunna-test_grafana

# ─── Prod ───
deploy-prod:
	docker stack deploy -c soporte/prod/compose.yml bunna-prod

remove-prod:
	docker stack rm bunna-prod

ps-prod:
	docker stack ps bunna-prod

logs-prod:
	docker service logs --tail 50 -f bunna-prod_kafka bunna-prod_clickhouse bunna-prod_grafana
