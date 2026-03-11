# Task Manager API — justfile

set dotenv-load

default:
    @just --list

# Build the application binary
build:
    go build -ldflags="-s -w" -o bin/api ./cmd/api

# Run the application locally
run:
    go run ./cmd/api

# Run unit tests
test:
    go test -v -race -count=1 ./internal/...

# Run unit tests with coverage
test-cover:
    go test -v -race -coverprofile=coverage.out ./internal/...
    go tool cover -html=coverage.out -o coverage.html
    @echo "Coverage report: coverage.html"

# Run integration tests (requires Docker)
test-integ:
    go test -v -race -count=1 -tags=integration ./tests/...

# Run all tests
test-all: test test-integ

# Run linter
lint:
    golangci-lint run ./...

# Generate Swagger documentation
swagger:
    swag init -g cmd/api/main.go -o docs --parseDependency

# Apply database migrations
migrate-up:
    goose -dir migrations postgres "host=${DB_HOST} port=${DB_PORT} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=${DB_SSL_MODE}" up

# Rollback last migration
migrate-down:
    goose -dir migrations postgres "host=${DB_HOST} port=${DB_PORT} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=${DB_SSL_MODE}" down

# Create a new migration (usage: just migrate-create name)
migrate-create name:
    goose -dir migrations create {{name}} sql

# Start infrastructure (PostgreSQL)
infra-up:
    docker-compose up -d postgres

# Stop infrastructure
infra-down:
    docker-compose down

# Start everything with Docker Compose
up:
    docker-compose up --build

# Format code
fmt:
    gofmt -s -w .
    goimports -w .
