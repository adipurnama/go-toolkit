SHELL := /bin/bash

BASE := $(shell pwd)
PKGS := $(shell go list ./... | grep -v /vendor/)
ALL_PACKAGES := $(shell go list ./... | grep -v /vendor/ | grep -v /internal/mock)
GIT_COMMIT := $(shell git rev-list -1 HEAD)
BUILD_TIME := $(shell date +%FT%T%z)

PHONY: fmt lint test run docker-image-push docker-image check-env generate
fmt:
	go fmt $(ALL_PACKAGES)

lint:
	env GO111MODULE=on golint -set_exit_status $(ALL_PACKAGES)

generate: fmt lint
	go generate ./...

test: generate
	go test -coverprofile=test_coverage.out $(PKGS)
	go tool cover -html=test_coverage.out -o test_coverage.html
	rm test_coverage.out
	@echo Open test_coverage.html file on your web browser for detailed coverage
