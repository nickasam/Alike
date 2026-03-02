.PHONY: help build run test clean docker-build docker-up migrate-up migrate-down

# Variables
APP_NAME := alike
GO_VERSION := 1.21
DOCKER_IMAGE := alike-api
DOCKER_TAG := latest

# Colors
COLOR_RESET := \033[0m
COLOR_BOLD := \033[1m
COLOR_GREEN := \033[32m
COLOR_YELLOW := \033[33m
COLOR_BLUE := \033[34m

## help: Show this help message
help:
	@echo '$(COLOR_BOLD)Alike Development Commands$(COLOR_RESET)'
	@echo ''
	@echo '$(COLOR_GREEN)Usage:$(COLOR_RESET)'
	@echo '  make [target]'
	@echo ''
	@echo '$(COLOR_GREEN)Targets:$(COLOR_RESET)'
	@grep -E '^## ' Makefile | sed 's/## /  /' | column -t -s ':'

## build: Build the application
build:
	@echo '$(COLOR_BLUE)Building $(APP_NAME)...$(COLOR_RESET)'
	@go build -o bin/$(APP_NAME) cmd/api/main.go
	@echo '$(COLOR_GREEN)Build complete!$(COLOR_RESET)'

## run: Run the application
run:
	@echo '$(COLOR_BLUE)Running $(APP_NAME)...$(COLOR_RESET)'
	@go run cmd/api/main.go

## test: Run tests
test:
	@echo '$(COLOR_BLUE)Running tests...$(COLOR_RESET)'
	@go test -v ./...

## test-coverage: Run tests with coverage
test-coverage:
	@echo '$(COLOR_BLUE)Running tests with coverage...$(COLOR_RESET)'
	@go test -v -coverprofile=coverage.txt ./...
	@go tool cover -html=coverage.txt -o coverage.html
	@echo '$(COLOR_GREEN)Coverage report generated: coverage.html$(COLOR_RESET)'

## clean: Clean build artifacts
clean:
	@echo '$(COLOR_BLUE)Cleaning...$(COLOR_RESET)'
	@rm -rf bin/
	@rm -rf coverage.txt coverage.html
	@rm -rf tmp/
	@echo '$(COLOR_GREEN)Clean complete!$(COLOR_RESET)'

## deps: Download dependencies
deps:
	@echo '$(COLOR_BLUE)Downloading dependencies...$(COLOR_RESET)'
	@go mod download
	@go mod tidy
	@echo '$(COLOR_GREEN)Dependencies updated!$(COLOR_RESET)'

## fmt: Format code
fmt:
	@echo '$(COLOR_BLUE)Formatting code...$(COLOR_RESET)'
	@go fmt ./...

## migrate-up: Run database migrations
migrate-up:
	@echo '$(COLOR_BLUE)Running migrations...$(COLOR_RESET)'
	@go run cmd/migrate/main.go up

## migrate-down: Rollback last migration
migrate-down:
	@echo '$(COLOR_BLUE)Rolling back migration...$(COLOR_RESET)'
	@go run cmd/migrate/main.go down

## docker-build: Build Docker image
docker-build:
	@echo '$(COLOR_BLUE)Building Docker image...$(COLOR_RESET)'
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) -f deployments/docker/Dockerfile .
	@echo '$(COLOR_GREEN)Docker image built!$(COLOR_RESET)'

## docker-up: Start Docker containers
docker-up:
	@echo '$(COLOR_BLUE)Starting Docker containers...$(COLOR_RESET)'
	@docker-compose -f deployments/docker/docker-compose.yml up -d

## docker-down: Stop Docker containers
docker-down:
	@echo '$(COLOR_BLUE)Stopping Docker containers...$(COLOR_RESET)'
	@docker-compose -f deployments/docker/docker-compose.yml down

## setup: Set up development environment
setup: deps
	@echo '$(COLOR_GREEN)Development environment ready!$(COLOR_RESET)'

.DEFAULT_GOAL := help
