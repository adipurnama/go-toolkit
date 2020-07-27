SHELL := /bin/bash

BASE := $(shell pwd)
PKGS := $(shell go list ./... | grep -v /vendor/)
ALL_PACKAGES := $(shell go list ./... | grep -v /vendor/ | grep -v /internal/mock)
GIT_COMMIT := $(shell git rev-list -1 HEAD)
BUILD_TIME := $(shell date +%FT%T%z)

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
	@env golangci-lint run && golint -set_exit_status $(ALL_PACKAGES)
	@echo -e "$(OK_COLOR)==> all is good$(NO_COLOR)..."

generate: fmt lint
	go generate ./...

test: generate
	go test -coverprofile=test_coverage.out $(PKGS)
	go tool cover -html=test_coverage.out -o test_coverage.html
	rm test_coverage.out
	@echo Open test_coverage.html file on your web browser for detailed coverage
