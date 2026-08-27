.PHONY: help build run clean dev lint format migrate-up migrate-down docs-generate docker-up docker-down docker-reset

BIN_DIR := bin
DB_URL := postgresql://postgres:password@localhost:5433/bookvault?sslmode=disable
MAIN := ./cmd/api
COMPOSE_FILE := docker/docker-compose.yml

help:
	@echo "Available commands"
	@echo "make build           -    Build the app"
	@echo "make run             -    Run the app"
	@echo "make dev             -    Run the application in dev mode"
	@echo "make lint            -    Run linter on the codebase"
	@echo "make format          -    Format the code"
	@echo "make docs-generate   -    Generate Swagger Documentation"
	@echo "make migrate-up      -    Apply db migration"
	@echo "make migrate-down    -    Rollback database migration"
	@echo "make docker-up       -    Start Docker services"
	@echo "make docker-down     -    Stop Docker services"	
	@echo "make docker-reset    -    Stop Docker services and remove volumes"
	@echo "make generate-graph" -    Starts generating graphql schema 

build:
	@echo "Building all binaries..."
	@mkdir -p $(BIN_DIR)
	@for cmd in cmd/*/; do \
		binary=$$(basename $$cmd); \
		echo "Building $$binary..."; \
		go build -o $(BIN_DIR)/$$binary ./$$cmd; \
	done
	

run:  
	go run $(MAIN)

clean:
	rm -rf $(BIN_DIR)

dev:
	go run $(MAIN)

lint: format
	golangci-lint run ./...

format:
	@gofmt -s -w .
	@goimports -w .

docs-generate:
	@mkdir -p docs
	swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal --exclude .git,docs,docker,db

migrate-up:
	migrate -path db/migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path db/migrations -database "$(DB_URL)" down

docker-up:
	docker compose -f $(COMPOSE_FILE) up -d

docker-down:
	docker compose -f $(COMPOSE_FILE) down

docker-reset:
	docker compose -f $(COMPOSE_FILE) down -v

generate-graph:
	@go get github.com/99designs/gqlgen@v0.17.94
	@go run github.com/99designs/gqlgen generate

