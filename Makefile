BIN := "./.bin/app"
DOCKER_COMPOSE_FILE := "./deployments/docker-compose.yml"
DOCKER_COMPOSE_TEST_FILE := "./deployments/docker-compose.tests.yml"
APP_NAME := "image-previewer"
APP_TEST_NAME := "image-previewer_test"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -a -o $(BIN) -ldflags "$(LDFLAGS)" cmd/main.go

run: build 
	 $(BIN)

test: 
	go test --short -race ./internal/...

.PHONY: build test

mocks:
	mockgen -source=./internal/services/services.go -destination ./internal/services/mocks/mock.go

up:
	docker-compose -f ${DOCKER_COMPOSE_FILE} -p ${APP_NAME} up --build

down:
	docker-compose -f ${DOCKER_COMPOSE_FILE} down --volumes

integration-tests:
	docker-compose -f ${DOCKER_COMPOSE_TEST_FILE} -p ${APP_TEST_NAME} up --build --abort-on-container-exit --exit-code-from integration-app
	docker-compose -f ${DOCKER_COMPOSE_TEST_FILE} -p ${APP_TEST_NAME} down --volumes

reset-integration-tests:
	docker-compose -f ${DOCKER_COMPOSE_TEST_FILE} -p ${APP_TEST_NAME} down --volumes
