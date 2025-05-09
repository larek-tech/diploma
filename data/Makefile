LOCAL_BIN:=$(CURDIR)/bin
MIGRATION_DIR=$(CURDIR)/internal/migrations
PROTOC=protoc
PROTO_SRC=./proto

.PHONY: migration-create
migration-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Usage: make migration-create NAME=your_migration_name"; \
		exit 1; \
	fi
	goose -dir ${MIGRATION_DIR} create $(NAME) sql

.PHONY: proto
proto: DATA_PROTO_SRC=$(PROTO_SRC)/data/v1
proto:
	@for dir in $(shell find . -type f -name go.mod -exec dirname {} \;); do \
		echo "Generating stubs in $$dir";\
		$(PROTOC) --proto_path=$(PROTO_SRC) --go_out=$$dir --go-grpc_out=$$dir \
			$(DATA_PROTO_SRC)/service.proto $(DATA_PROTO_SRC)/model.proto; \
	done