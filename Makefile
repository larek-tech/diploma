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
	docker volume rm ${BASE_IMAGE}_zoo_data_log
	docker volume rm ${BASE_IMAGE}_zoo_data_secrets
	docker volume rm ${BASE_IMAGE}_zoo_data
	docker volume rm ${BASE_IMAGE}_kafka_data_cfg
	docker volume rm ${BASE_IMAGE}_kafka_data
	docker volume rm ${BASE_IMAGE}_kafka_data_secrets
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
	@for dir in $(shell find . -type f -name go.mod -not -path "./data/*" -exec dirname {} \;); do \
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
proto-auth: AUTH_PROTO_SRC=$(PROTO_SRC)/auth/v1
proto-auth:
	@for dir in $(shell find . -type f -name go.mod -not -path "./data/*" -exec dirname {} \;); do \
		echo "Generating stubs in $$dir";\
		$(PROTOC) --proto_path=$(PROTO_SRC) --go_out=$$dir --go-grpc_out=$$dir \
			$(AUTH_PROTO_SRC)/*.proto; \
	done
	@echo "Protobuf stubs for auth service generated\n"

.PHONY: proto-ml
proto-ml: ML_PROTO_SRC=$(PROTO_SRC)/ml/v1
proto-ml:
	@for dir in ./chat ./domain ./api; do \
  		echo "Generating stubs in $$dir";\
		$(PROTOC) --proto_path=$(PROTO_SRC) --go_out=$$dir --go-grpc_out=$$dir \
			$(ML_PROTO_SRC)/*.proto; \
	done
	@mkdir -p ml/src
	@echo "Generating stubs in ./ml"
	@python3 -m grpc_tools.protoc -I$(PROTO_SRC) --python_out=ml/src --pyi_out=ml/src --grpc_python_out=ml/src \
		$(ML_PROTO_SRC)/*.proto;
	@echo "Protobuf stubs for ml service generated\n"

.PHONY: proto-data
proto-data: DATA_PROTO_SRC=$(PROTO_SRC)/data/v1
proto-data:
	@mkdir -p ml/src
	@echo "Generating stubs in ./ml"
	@python3 -m grpc_tools.protoc -I$(PROTO_SRC) --python_out=ml/src --pyi_out=ml/src --grpc_python_out=ml/src \
		$(DATA_PROTO_SRC)/*.proto;
	@echo "Protobuf stubs for data service generated\n"

.PHONY: proto-domain
proto-domain: DOMAIN_PROTO_SRC=$(PROTO_SRC)/domain/v1
proto-domain:
	@for dir in ./domain ./api; do \
		echo "Generating stubs in $$dir";\
		$(PROTOC) --proto_path=$(PROTO_SRC) --go_out=$$dir --go-grpc_out=$$dir \
			$(DOMAIN_PROTO_SRC)/*.proto \
			$(PROTO_SRC)/google/protobuf/*.proto; \
	done
	@echo "Protobuf stubs for domain service generated\n"

.PHONY: proto-chat
proto-chat: CHAT_PROTO_SRC=$(PROTO_SRC)/chat/v1
proto-chat:
	@for dir in ./chat ./api; do \
		echo "Generating stubs in $$dir";\
		$(PROTOC) --proto_path=$(PROTO_SRC) --go_out=$$dir --go-grpc_out=$$dir \
			$(CHAT_PROTO_SRC)/*.proto \
			$(PROTO_SRC)/google/protobuf/*.proto; \
	done
	@echo "Protobuf stubs for chat service generated\n"

.PHONY: proto
proto: proto-auth proto-ml proto-data proto-domain proto-chat
	@echo "All protobuf stubs generated"

.PHONY: swag
swag:
	mkdir -p api/docs
	cd api && SWAG init -g cmd/server/main.go -o ./docs && SWAG fmt
