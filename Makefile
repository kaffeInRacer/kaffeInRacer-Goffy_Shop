include .env

docker-up:
	docker run --name meh_db_postgres \
		-e POSTGRES_USER=$(POSTGRES_USER) \
		-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
		-e POSTGRES_DB=$(POSTGRES_DB) \
		-p 5432:5432 \
		-d postgres:17-alpine

migrate-up:
	@migrate -path pkg/database/postgresql/migrations -database "postgres://$(strip ${POSTGRES_USER}):$(strip ${POSTGRES_PASSWORD})@$(strip ${POSTGRES_HOSTNAME}):$(strip ${POSTGRES_PORT})/$(strip ${POSTGRES_DB})?sslmode=disable" -verbose up

migrate-down:
	@migrate -path pkg/database/postgresql/migrations -database "postgres://$(strip ${POSTGRES_USER}):$(strip ${POSTGRES_PASSWORD})@$(strip ${POSTGRES_HOSTNAME}):$(strip ${POSTGRES_PORT})/$(strip ${POSTGRES_DB})?sslmode=disable" -verbose down

migrate-create:
	@migrate create -ext sql -dir ./pkg/database/postgresql/migrations/ -seq $(name)

migrate-force:
	@migrate -path pkg/database/postgresql/migrations -database "postgres://$(strip ${POSTGRES_USER}):$(strip ${POSTGRES_PASSWORD})@$(strip ${POSTGRES_HOSTNAME}):$(strip ${POSTGRES_PORT})/$(strip ${POSTGRES_DB})?sslmode=disable" force 5