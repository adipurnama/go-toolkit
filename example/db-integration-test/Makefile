SHELL := /bin/bash

.PHONY: postgres

PG_CONTAINER:=$(shell docker ps -aq --format '{{.Names}}' | grep -w pg_go_db_integration_test)

postgres:
ifeq ("$(PG_CONTAINER)", "pg_go_db_integration_test")
	docker start pg_go_db_integration_test
else
	docker run --name pg_go_db_integration_test \
	-e POSTGRES_USER=postgres \
	-e POSTGRES_PASSWORD=postgres \
	-e POSTGRES_DB=go_db_integration_test_dev \
	-p 5432:5432 \
	-d postgres:alpine
endif
