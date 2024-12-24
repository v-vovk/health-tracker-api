include .env

# Variables
DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
MIGRATE_CMD := migrate -database "$(DB_URL)" -path migrations

# Commands
.PHONY: migrate-up migrate-down migrate-force migrate-create migrate-version migrate-status run test build lint clean dev

migrate-up:
	$(MIGRATE_CMD) up

migrate-down:
	$(MIGRATE_CMD) down

migrate-force:
	$(MIGRATE_CMD) force 1

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

migrate-version:
	$(MIGRATE_CMD) version

run:
	go run cmd/main.go

test:
	go test ./... -v

build:
	go build -o bin/health-tracker ./cmd/main.go

lint:
	@golangci-lint run --config .golangci.yml

clean:
	@rm -rf bin/ logs/app.log

dev:
	air
