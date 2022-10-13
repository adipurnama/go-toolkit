SHELL := /bin/bash

BASE := $(shell pwd)
PKGS := $(shell go list ./... | grep -v /vendor/)
ALL_PACKAGES := $(shell go list ./... | grep -v /vendor/ | grep -v /internal/mock)
GIT_COMMIT := $(shell git rev-list -1 HEAD)
BUILD_TIME := $(shell date +%FT%T%z)
PROTOC_CMD := $(shell command -v protoc 2> /dev/null)

NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

PHONY: fmt lint test run docker-image-push docker-image check-env generate
fmt:
	@echo -e "$(OK_COLOR)==> formatting code$(NO_COLOR)..."
	@go fmt $(ALL_PACKAGES)

lint: fmt
	@echo -e "$(OK_COLOR)==> linting source files$(NO_COLOR)..."
	@env golangci-lint run --fix && golint -set_exit_status $(ALL_PACKAGES)
	@echo -e "$(OK_COLOR)==> all is good$(NO_COLOR)..."

generate: fmt lint
	@go generate ./...

test: generate
	@go test -coverprofile=test_coverage.out $(PKGS)
	@go tool cover -html=test_coverage.out -o test_coverage.html
	@rm test_coverage.out
	@echo Open test_coverage.html file on your web browser for detailed coverage

gen-proto:
ifndef PROTOC_CMD
	$(error "protoc-gen-go is not installed. Run command 'go get -u google.golang.org/protobuf/proto && go install github.com/golang/protobuf/protoc-gen-go'")
endif
	@echo -e "$(OK_COLOR)==> Generate proto objects to grpckit/$(NO_COLOR)..."
	@protoc --proto_path=./grpckit --go_out=plugins=grpc:./ grpckit/health.proto
	@protoc --proto_path=./examples/grpc-server --go_out=plugins=grpc:./ examples/grpc-server/example_service.proto
	@echo -e "$(OK_COLOR)==> Done$(NO_COLOR)..."

run-pubsub: pubsub-local
	PUBSUB_EMULATOR_HOST=localhost:8085 go run examples/gcp-pubsub/main.go

run-springconfig-docker:
	SPRING_CLOUD_CONFIG_URL="http://localhost:8888/" \
	SPRING_CLOUD_CONFIG_PATHS="/go-config-app/dev/,/go-config-app/other/" \
							go run examples/springcloud-config/main.go

run-elasticapm-echo: apm-local
	ELASTIC_APM_SERVER_URL="http://localhost:8200" \
         ELASTIC_APM_ENVIRONMENT="local-test" \
         ELASTIC_APM_SERVICE_NAME="echo-test-apm-middleware" \
         go run examples/echo-restapi/main.go

apm-local:
	docker compose up -d elasticsearch kibana es-apm-server

pubsub-local:
	docker compose up -d googlecloud-pubsub

springcloud-config-docker:
	docker compose up -d spring-config-server

springcloud-config-localfile:
	SPRING_CLOUD_CONFIG_URL=file://$(PWD)/examples/springcloud-config/data/go-config-app-dev.yml \
	SPRING_CLOUD_CONFIG_PATHS="/" \
							go run examples/springcloud-config/main.go

run-echo-grpc:
	SPRING_CLOUD_CONFIG_URL=file://$(PWD)/examples/springcloud-config/data/go-config-app-dev.yml \
	SPRING_CLOUD_CONFIG_PATHS="/" \
							go run examples/echo-grpc-sample/main.go

run-echo-restapi:
	SPRING_CLOUD_CONFIG_URL=file://$(PWD)/examples/springcloud-config/data/go-config-app-dev.yml \
	SPRING_CLOUD_CONFIG_PATHS="/" \
							go run examples/echo-restapi/main.go
