GOOSE_DBSTRING ?= "root:12345@tcp(127.0.0.1:33306)/shopGO"
GOOSE_MIGRATION_DIR ?= sql/schema
GOOSE_DRIVER ?= mysql

# Tên của ứng dụng của bạn
APP_NAME := server

# Chạy ứng dụng

docker_build:
	docker-compose up -d --build
	docker-compose ps

docker_stop:
	docker-compose down

dev:
	go run ./cmd/$(APP_NAME)

docker_up:
	docker-compose -f docker-compose.yaml compose up

up_by_one:
	goose -dir=$(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) $(GOOSE_DBSTRING) up-by-one

# Create a new migration
create_migration:
	goose -dir=$(GOOSE_MIGRATION_DIR) create $(name) sql

upse:
	goose -dir=$(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) $(GOOSE_DBSTRING) up

downse:
	goose -dir=$(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) $(GOOSE_DBSTRING) down

resetse:
	goose -dir=$(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) $(GOOSE_DBSTRING) reset

sqlgen:
	sqlc generate

swag:
	swag init -g ./cmd/server/main.go -o ./cmd/swag/docs

.PHONY: dev downse upse resetse docker_build docker_stop docker_up swag

.PHONY: air
