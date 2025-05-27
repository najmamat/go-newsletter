.PHONY: help build test clean generate run dev docker-build docker-run

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Build the application
build: ## Build the Go application
	@echo "Building the application..."
	go build -o bin/server ./cmd/server

# Run tests
test: ## Run all tests
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf pkg/generated/

# Generate code from OpenAPI specification
generate: ## Generate Go code from OpenAPI specification
	@echo "Generating code from OpenAPI specification..."
	mkdir -p pkg/generated
	$(shell go env GOPATH)/bin/oapi-codegen -config api/oapi-config.yaml api/openapi.yaml

# Run the application locally
run: build ## Build and run the application
	@echo "Starting the server..."
	./bin/server

# Run in development mode with auto-reload (requires air)
dev: ## Run in development mode with auto-reload
	@echo "Starting development server..."
	@if command -v air > /dev/null 2>&1; then \
		air; \
	else \
		echo "Air not installed. Install with: go install github.com/air-verse/air@latest"; \
		echo "Running with go run instead..."; \
		go run ./cmd/server; \
	fi

# Install development dependencies
dev-deps: ## Install development dependencies
	@echo "Installing development dependencies..."
	go install github.com/air-verse/air@latest
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	go install github.com/go-delve/delve/cmd/dlv@master

# Format code
fmt: ## Format Go code
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint: ## Lint Go code (requires golangci-lint)
	@echo "Linting code..."
	@if command -v golangci-lint > /dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Tidy dependencies
tidy: ## Clean up dependencies
	@echo "Tidying dependencies..."
	go mod tidy

# Docker build
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t go-newsletter .

# Docker run
docker-run: docker-build ## Build and run Docker container
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file .env go-newsletter

# Database migration (placeholder for future implementation)
migrate-up: ## Run database migrations up
	@echo "Database migrations not yet implemented"

migrate-down: ## Run database migrations down
	@echo "Database migrations not yet implemented"

# Full development setup
setup: dev-deps tidy generate fmt ## Set up development environment
	@echo "Development environment set up complete!" 