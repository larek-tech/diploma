PROTOC=protoc
PROTO_SRC=./proto

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
	@python -m grpc_tools.protoc -I$(PROTO_SRC) --python_out=ml --pyi_out=ml --grpc_python_out=ml $(ML_PROTO_SRC);

.PHONY: proto
proto: proto-auth proto-ml
