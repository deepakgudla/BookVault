.PHONY: help build run dev lint migrate-up migrate down

help:
	@echo "Available commands"
	@echo " make build        -    build the app"
	@echo " make run          -    run the app"
	@echo " make dev          -    run the application in dev mode"
	@echo "make lint          -    run linter on the codebase"
	@echo "make migrate up    -    apply db migration"
	@echo "make migrate-down  -    rollback database migration"
	
build:
	go build -o bin/app ./cmd/api

run:
	go run ./cmd/api

dev:
	go run ./cmd/api

lint:
	golangci-lint run ./...

migrate-up:
	migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/bookvault?sslmode=disable" up

migrate-down:
	migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/bookvault?sslmode=disable" down

docker-up:
	docker compose -f docker/docker-compose.yml up -d

docker-down:
	docker compose -f docker/docker-compose.yml down

