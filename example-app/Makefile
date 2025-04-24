MPAAS=../minipaas-cli/build/minipaas
ANSIBLE_PLAYBOOK=install.yml
INVENTORY=inventory.ini

.PHONY: prod

prod:
	$(MPAAS) code init --env prod -c compose.yaml -c compose.build.yaml -c prod/compose.infra.yaml --host localhost
	$(MPAAS) certs client --env prod --verbose --ca-dir infra/.certs
	echo postgres | $(MPAAS) secret create --verbose --env prod --name postgres_password --for postgres --for example --for example-migration --for example-consumer --for example-worker --for example-cron
	$(MPAAS) code expose --env prod example:8080 example.local --verbose
	$(MPAAS) code job --env prod example-migration
	$(MPAAS) code worker --env prod postgres example-worker example-consumer
	$(MPAAS) code cron --env prod example-cron

build:
	$(MPAAS) deploy build --verbose --env prod

rollout:
	$(MPAAS) deploy rollout --verbose --env prod

canary:
	$(MPAAS) deploy canary --verbose --env prod $(service)

routing:
	$(MPAAS) deploy routing --verbose --env prod

shell:
	$(MPAAS) shell --env prod --verbose
