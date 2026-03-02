.PHONY: help build run test clean docker-build docker-up migrate-up migrate-down

# Variables
APP_NAME := alike
GO_VERSION := 1.23
DOCKER_IMAGE := alike-api
DOCKER_TAG := latest

# Colors
COLOR_RESET := \033[0m
COLOR_BOLD := \033[1m
COLOR_GREEN := \033[32m
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
	@go build -o bin/$(APP_NAME)-api cmd/api/main.go
	@go build -o bin/$(APP_NAME)-migrate cmd/migrate/main.go
	@echo '$(COLOR_GREEN)Build complete!$(COLOR_RESET)'

## run: Run the application
run:
	@echo '$(COLOR_BLUE)Running $(APP_NAME)...$(COLOR_RESET)'
	@go run cmd/api/main.go

## test: Run tests
test:
	@echo '$(COLOR_BLUE)Running tests...$(COLOR_RESET)'
	@go test -v ./...

## clean: Clean build artifacts
clean:
	@echo '$(COLOR_BLUE)Cleaning...$(COLOR_RESET)'
	@rm -rf bin/
	@echo '$(COLOR_GREEN)Clean complete!$(COLOR_RESET)'

## deps: Download dependencies
deps:
	@echo '$(COLOR_BLUE)Downloading dependencies...$(COLOR_RESET)'
	@go mod download
	@go mod tidy
	@echo '$(COLOR_GREEN)Dependencies updated!$(COLOR_RESET)'

## migrate-up: Run database migrations
migrate-up:
	@echo '$(COLOR_BLUE)Running migrations...$(COLOR_RESET)'
	@go run cmd/migrate/main.go up

## migrate-down: Rollback last migration
migrate-down:
	@echo '$(COLOR_BLUE)Rolling back migration...$(COLOR_RESET)'
	@go run cmd/migrate/main.go down

## seed: Seed database with sample data
seed:
	@echo '$(COLOR_BLUE)Seeding database...$(COLOR_RESET)'
	@psql -h localhost -U alike_user -d alike_db -f db/seeds/seed.sql

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

## docker-logs: Show Docker logs
docker-logs:
	@docker-compose -f deployments/docker/docker-compose.yml logs -f api

## deploy: Deploy to production (uses scripts/deploy.sh)
deploy:
	@./scripts/deploy.sh

## dev: Run in development mode
dev:
	@./scripts/dev.sh

.DEFAULT_GOAL := help
