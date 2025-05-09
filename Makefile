include .env

BASE_IMAGE=diploma
PROTOC=protoc
PROTO_SRC=./proto
GOLANGCI_LINT=golangci-lint
GOOSE=bin/goose
SWAG=swag

.PHONY: goose-download
goose-download:
	mkdir -p bin
	curl -fsSL \
		https://raw.githubusercontent.com/pressly/goose/master/install.sh |\
		GOOSE_INSTALL=. sh

.PHONY: docker-build-up
docker-build-up:
	docker compose up --build -d

.PHONY: docker-up
docker-up:
	docker compose up -d

.PHONY: docker-down
docker-down:
	docker compose down

.PHONY: docker-remove
docker-remove:
	docker volume rm ${BASE_IMAGE}_pg_data
	docker volume rm ${BASE_IMAGE}_jaeger_data
	docker volume rm ${BASE_IMAGE}_zoo_data
	docker volume rm ${BASE_IMAGE}_kafka_data
	docker image rm ${BASE_IMAGE}-auth
	docker image rm ${BASE_IMAGE}-api
	docker image rm ${BASE_IMAGE}-chat
	docker image rm ${BASE_IMAGE}-domain

.PHONY: lint
lint:
	@echo "Starting linter"
	@for dir in $(shell find . -type f -name go.mod -exec dirname {} \;); do \
		echo "Running linter in $$dir"; \
		cd $$dir; \
		$(GOLANGCI_LINT) run --config ../.golangci.yml; \
		cd ..; \
	done


.PHONY: vendor
vendor:
	@for dir in $(shell find . -type f -name go.mod -exec dirname {} \;); do \
		echo "Vendoring $$dir"; \
		cd $$dir; \
		go mod tidy; \
		go mod vendor; \
		cd ..; \
	done
	@echo "Vendoring complete\n"

.PHONY: migrate-up
migrate-up:
	mkdir -p migrations
	cd migrations && ../$(GOOSE) postgres "user=${POSTGRES_USER} \
		password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} sslmode=disable \
		host=${POSTGRES_HOST} port=${POSTGRES_PORT}" up

.PHONY: migrate-down
migrate-down:
	mkdir -p migrations
	cd migrations && ../$(GOOSE) postgres "user=${POSTGRES_USER} \
		password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} sslmode=disable \
		host=localhost port=${POSTGRES_PORT}" down

.PHONY: migrate-new
migrate-new:
	mkdir -p migrations
	cd migrations && ../$(GOOSE) create $(name) sql

.PHONY: proto-auth
proto-auth: AUTH_PROTO_SRC=$(PROTO_SRC)/auth/v1/service.proto
proto-auth:
	@for dir in $(shell find . -type f -name go.mod -exec dirname {} \;); do \
		$(PROTOC) --proto_path=$(PROTO_SRC) --go_out=$$dir --go-grpc_out=$$dir $(AUTH_PROTO_SRC); \
	done

.PHONY: proto-ml
proto-ml: ML_PROTO_SRC=$(PROTO_SRC)/ml/v1/service.proto
proto-ml:
	@for dir in $(shell find . -type f -name go.mod -exec dirname {} \;); do \
		$(PROTOC) --proto_path=$(PROTO_SRC) --go_out=$$dir --go-grpc_out=$$dir $(ML_PROTO_SRC); \
	done
	mkdir -p ml/pb
	@python -m grpc_tools.protoc -I$(PROTO_SRC) --python_out=ml/pb --pyi_out=ml/pb --grpc_python_out=ml/pb $(ML_PROTO_SRC);

.PHONY: proto
proto: proto-auth proto-ml
