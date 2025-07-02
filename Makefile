.PHONY: build run clean deps tidy docker-build docker-run docker-stop migrate-up migrate-down swagger

export DATABASE_URL ?= postgres://postgres:password@localhost:5432/voting_db?sslmode=disable

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application locally
run:
	go run cmd/server/main.go

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Install dependencies
deps:
	go mod download

# Update go.mod
tidy:
	go mod tidy

# Generate Swagger documentation
swagger:
	swag init -g cmd/server/main.go

# Build Docker image
docker-build:
	docker build -t real-time-voting .

# Run with Docker Compose
docker-run:
	docker-compose up --build

# Stop Docker Compose
docker-stop:
	docker-compose down

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run



# Development setup (install dependencies, generate docs, run migrations)
dev-setup: deps swagger migrate-up

# Full development workflow
dev: dev-setup run 