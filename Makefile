.PHONY: build run test lint clean docker-build docker-up docker-down help

# Variables
APP_NAME=quotebox
DOCKER_IMAGE=quotebox:latest
DOCKER_COMPOSE=docker compose

help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

build: ## Build the Go application
	@echo "Building $(APP_NAME)..."
	go build -o bin/$(APP_NAME) ./cmd/server

run: ## Run the application locally
	@echo "Running $(APP_NAME)..."
	go run ./cmd/server/main.go

test: ## Run all tests
	@echo "Running tests..."
	go test ./... -v -cover

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	go test ./tests/integration/... -v

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run ./...

security: ## Run security scan
	@echo "Running security scan..."
	gosec ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	go clean

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

docker-up: ## Start all services with Docker Compose
	@echo "Starting services..."
	$(DOCKER_COMPOSE) up -d

docker-up-build: ## Build and start all services
	@echo "Building and starting services..."
	$(DOCKER_COMPOSE) up --build -d

docker-down: ## Stop all services
	@echo "Stopping services..."
	$(DOCKER_COMPOSE) down

docker-logs: ## View logs from all services
	$(DOCKER_COMPOSE) logs -f

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

check-env: ## Check environment setup
	@echo "Checking environment..."
	@if [ -f scripts/check_env.sh ]; then bash scripts/check_env.sh; else powershell -ExecutionPolicy Bypass -File scripts/check_env.ps1; fi

.DEFAULT_GOAL := help
