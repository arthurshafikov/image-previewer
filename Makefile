BIN := "./.bin/app"
DOCKER_COMPOSE_FILE := "./deployments/docker-compose.yml"

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
	docker-compose -f ${DOCKER_COMPOSE_FILE} up --build

down:
	docker-compose -f ${DOCKER_COMPOSE_FILE} down --volumes
