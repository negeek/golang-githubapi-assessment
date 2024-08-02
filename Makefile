.PHONY: install run create_table migrate_up migrate_down help

help:
	@echo "Available commands:"
	@echo "  install        - Install dependencies"
	@echo "  run            - Run the application"
	@echo "  create_table   - Create a new migration file"
	@echo "  migrate_up     - Apply migrations"
	@echo "  migrate_down   - Revert migrations"

install:
	go mod tidy
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

run:
	go run main.go

create_table:
	@if [ -z "${MIGRATION_NAME}" ]; then \
		echo "FILE_NAME is not set"; \
		exit 1; \
	fi
	migrate create -ext sql -dir db/migrations -seq ${MIGRATION_NAME}

migrate_up:
	migrate -database ${POSTGRESQL_URL} -path db/migrations up

migrate_down:
	migrate -database ${POSTGRESQL_URL} -path db/migrations down

test:
	go test ./...
