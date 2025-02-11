include .env
export

CONTAINER = app
DOCKER_COMPOSE_EXISTS := $(shell command -v docker-compose 2> /dev/null)
BINARY_NAME=app
BUILD_DIR=build
SRC_DIR=cmd/pkg
DOCKER_IMAGE_NAME=hf_app
DOCKER_IMAGE_BUILD_VERSION ?= v$(shell date +%s)

.PHONY: start stop ps ssh run build deps run-prod-container deploy test

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
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(BINARY_NAME) $(SRC_DIR)/main.go

	@echo "build version $(DOCKER_IMAGE_BUILD_VERSION)"

	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_BUILD_VERSION) \
	  --build-arg SERVE_PORT=$(SERVE_PORT) \
	  --build-arg LOG_LEVEL=$(LOG_LEVEL) \
	  -f Dockerfile.prod .

run-prod-container:
	@echo "Running standalone container"
	docker run -p $(SERVE_PORT):$(SERVE_PORT) $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_BUILD_VERSION)

deploy:
	@if [ "$(ENVIRONMENT)" = "local" ]; then \
	    echo "Staring minikube and loading image"; \
		minikube start ; \
		minikube image load $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_BUILD_VERSION) ; \
	fi

	kubectl apply -f deployments/deployment.yaml -f deployments/service.yaml

	@if [ "$(ENVIRONMENT)" = "local" ]; then \
	    echo "Tunneling LoadBalancer through minikube"; \
		minikube tunnel --cleanup; \
	fi

shutdown:
	@echo "Stopping kubernetes environment"

	kubectl delete -f deployments/deployment.yaml -f deployments/service.yaml ; \

	@if [ "$(ENVIRONMENT)" = "local" ]; then \
		minikube stop; \
	fi

deps:
	${call app_container, mod vendor}

test:
	go test ./...

#---- docker enviroment ----
ifdef DOCKER_COMPOSE_EXISTS
define app_container
	@docker compose exec ${CONTAINER} go ${1}
endef
endif
