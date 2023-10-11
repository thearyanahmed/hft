include .env

CONTAINER = app
DOCKER_COMPOSE_EXISTS := $(shell command -v docker-compose 2> /dev/null)

.PHONY: start stop ps ssh run build deps test

start:
	@echo "starting application environment"
	@docker compose up

stop:
	@docker compose stop

ps:
	@docker compose ps

ssh:
	@docker compose exec $(CONTAINER) bash

run:
	@CompileDaemon -build='make build' -graceful-kill -command='./build/app'

build:
	@echo "running build for development."
	@CGO_ENABLED=0 go build -o build/app -v cmd/pkg/main.go

deps:
	${call app_container, mod vendor}

test:
	${call app_container, test -v ./pkg/handler/}

#---- docker enviroment ----
ifdef DOCKER_COMPOSE_EXISTS
define app_container
	@docker compose exec ${CONTAINER} go ${1}
endef
endif
