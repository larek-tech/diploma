LOCAL_BIN:=$(CURDIR)/bin
MIGRATION_DIR=$(CURDIR)/internal/migrations

.PHONY: migration-create
migration-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Usage: make migration-create NAME=your_migration_name"; \
		exit 1; \
	fi
	goose -dir ${MIGRATION_DIR} create $(NAME) sql
