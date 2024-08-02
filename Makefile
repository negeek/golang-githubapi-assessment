POSTGRESQL_URL = postgres://postgresuser:postgrespass@db:5432/postgres?sslmode=disable

build:
	docker compose build

run:
	docker compose up

create_table:
	migrate create -ext sql -dir db/migrations -seq $(FILE_NAME)

migrate_up:
	docker compose run --rm app migrate -database ${POSTGRESQL_URL} -path db/migrations up

migrate_down:
	docker compose run --rm app migrate -database ${POSTGRESQL_URL} -path db/migrations down

test:
	docker compose run --rm app go test ./tests/v1 
