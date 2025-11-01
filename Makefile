.PHONY: help install test test-unit test-integration test-all clean docker-up docker-down \
        docker-test docker-clean lint fmt vet coverage deps build

# Variables
GOTEST := go test
GOVET := go vet
GOFMT := gofmt
GOLINT := golangci-lint

# Colors for output
GREEN  := \033[0;32m
YELLOW := \033[1;33m
RED    := \033[0;31m
NC     := \033[0m # No Color

# Test database configuration
export TEST_DB_HOST ?= localhost
export TEST_DB_PORT ?= 5432
export TEST_DB_USER ?= postgres
export TEST_DB_PASSWORD ?= postgres
export TEST_DB_NAME ?= outbox_test

# Default target
.DEFAULT_GOAL := help

help: ## Show this help message
	@echo '$(GREEN)Outbox Pattern - Makefile Commands$(NC)'
	@echo ''
	@echo 'Usage:'
	@echo '  $(YELLOW)make$(NC) $(GREEN)<target>$(NC)'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install dependencies
	@echo "$(GREEN)Installing dependencies...$(NC)"
	cd outbox_uow_abstraction && go mod download
	cd outbox_uow_sql && GOINSECURE="proxy.golang.org" go mod download
	@echo "$(GREEN)Dependencies installed successfully!$(NC)"

deps: install ## Alias for install

tidy: ## Tidy go modules
	@echo "$(GREEN)Tidying go modules...$(NC)"
	cd outbox_uow_abstraction && go mod tidy
	cd outbox_uow_sql && go mod tidy
	@echo "$(GREEN)Go modules tidied!$(NC)"

test-unit: ## Run unit tests only
	@echo "$(GREEN)Running unit tests...$(NC)"
	@echo ""
	@echo "$(YELLOW)Testing Abstraction Module...$(NC)"
	cd outbox_uow_abstraction && $(GOTEST) -v ./abstraction
	@echo ""
	@echo "$(YELLOW)Testing PostgreSQL Module (Unit Tests)...$(NC)"
	cd outbox_uow_sql && $(GOTEST) -v ./sql_channel
	@echo ""
	@echo "$(GREEN)Unit tests completed!$(NC)"

test-integration: ## Run integration tests (requires PostgreSQL)
	@echo "$(GREEN)Running integration tests...$(NC)"
	@echo "$(YELLOW)Note: Requires PostgreSQL to be running$(NC)"
	@if ! pg_isready -h $(TEST_DB_HOST) -p $(TEST_DB_PORT) -U $(TEST_DB_USER) > /dev/null 2>&1; then \
		echo "$(RED)Error: PostgreSQL is not ready at $(TEST_DB_HOST):$(TEST_DB_PORT)$(NC)"; \
		echo "$(YELLOW)Start PostgreSQL with: make docker-up$(NC)"; \
		exit 1; \
	fi
	@echo ""
	cd outbox_uow_sql && $(GOTEST) -tags=integration -v ./sql_channel
	@echo ""
	@echo "$(GREEN)Integration tests completed!$(NC)"

test-integration-docker: ## Run integration tests inside Docker (used by docker-compose)
	@echo "$(GREEN)Running integration tests in Docker...$(NC)"
	@echo "Waiting for PostgreSQL to be ready..."
	@until pg_isready -h $(TEST_DB_HOST) -p $(TEST_DB_PORT) -U $(TEST_DB_USER) > /dev/null 2>&1; do \
		echo "Waiting for PostgreSQL..."; \
		sleep 1; \
	done
	@echo "PostgreSQL is ready!"
	@echo ""
	cd outbox_uow_sql && $(GOTEST) -tags=integration -v ./sql_channel
	@echo ""
	@echo "$(GREEN)Integration tests completed!$(NC)"

test-all: test-unit test-integration ## Run all tests (unit + integration)

test: test-unit ## Alias for test-unit (default test target)

coverage: ## Generate test coverage report
	@echo "$(GREEN)Generating coverage reports...$(NC)"
	@echo ""
	@echo "$(YELLOW)Abstraction Module Coverage...$(NC)"
	cd outbox_uow_abstraction && $(GOTEST) -coverprofile=coverage.out ./abstraction
	cd outbox_uow_abstraction && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: outbox_uow_abstraction/coverage.html"
	@echo ""
	@echo "$(YELLOW)PostgreSQL Module Coverage...$(NC)"
	cd outbox_uow_sql && $(GOTEST) -coverprofile=coverage.out ./sql_channel
	cd outbox_uow_sql && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: outbox_uow_sql/coverage.html"
	@echo ""
	@echo "$(GREEN)Coverage reports generated!$(NC)"

coverage-integration: ## Generate integration test coverage
	@echo "$(GREEN)Generating integration test coverage...$(NC)"
	cd outbox_uow_sql && $(GOTEST) -tags=integration -coverprofile=coverage-integration.out ./sql_channel
	cd outbox_uow_sql && go tool cover -html=coverage-integration.out -o coverage-integration.html
	@echo "Coverage report: outbox_uow_sql/coverage-integration.html"
	@echo "$(GREEN)Integration coverage report generated!$(NC)"

fmt: ## Format Go code
	@echo "$(GREEN)Formatting Go code...$(NC)"
	$(GOFMT) -w -s outbox_uow_abstraction/
	$(GOFMT) -w -s outbox_uow_sql/
	@echo "$(GREEN)Code formatted!$(NC)"

vet: ## Run go vet
	@echo "$(GREEN)Running go vet...$(NC)"
	cd outbox_uow_abstraction && $(GOVET) ./...
	cd outbox_uow_sql && $(GOVET) ./...
	@echo "$(GREEN)go vet completed!$(NC)"

lint: ## Run linter (requires golangci-lint)
	@echo "$(GREEN)Running linter...$(NC)"
	@if command -v $(GOLINT) > /dev/null; then \
		cd outbox_uow_abstraction && $(GOLINT) run ./...; \
		cd ../outbox_uow_sql && $(GOLINT) run ./...; \
		echo "$(GREEN)Linting completed!$(NC)"; \
	else \
		echo "$(YELLOW)golangci-lint not installed. Skipping...$(NC)"; \
		echo "Install with: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$$(go env GOPATH)/bin"; \
	fi

build: ## Build the modules
	@echo "$(GREEN)Building modules...$(NC)"
	cd outbox_uow_abstraction && go build ./...
	cd outbox_uow_sql && go build ./...
	@echo "$(GREEN)Build completed!$(NC)"

clean: ## Clean build artifacts and test cache
	@echo "$(GREEN)Cleaning...$(NC)"
	go clean -testcache
	rm -f outbox_uow_abstraction/coverage.out outbox_uow_abstraction/coverage.html
	rm -f outbox_uow_sql/coverage.out outbox_uow_sql/coverage.html
	rm -f outbox_uow_sql/coverage-integration.out outbox_uow_sql/coverage-integration.html
	@echo "$(GREEN)Cleaned!$(NC)"

docker-up: ## Start PostgreSQL using Docker Compose
	@echo "$(GREEN)Starting PostgreSQL with Docker Compose...$(NC)"
	docker-compose up -d postgres
	@echo "$(YELLOW)Waiting for PostgreSQL to be ready...$(NC)"
	@until docker-compose exec -T postgres pg_isready -U postgres > /dev/null 2>&1; do \
		echo "Waiting for PostgreSQL..."; \
		sleep 1; \
	done
	@echo "$(GREEN)PostgreSQL is ready!$(NC)"
	@echo ""
	@echo "Connection details:"
	@echo "  Host: localhost"
	@echo "  Port: 5432"
	@echo "  User: postgres"
	@echo "  Password: postgres"
	@echo "  Database: outbox_test"

docker-down: ## Stop and remove Docker containers
	@echo "$(GREEN)Stopping Docker containers...$(NC)"
	docker-compose down
	@echo "$(GREEN)Containers stopped!$(NC)"

docker-clean: ## Stop containers and remove volumes
	@echo "$(GREEN)Cleaning Docker resources...$(NC)"
	docker-compose down -v
	@echo "$(GREEN)Docker resources cleaned!$(NC)"

docker-test: docker-up ## Start PostgreSQL and run all tests
	@echo "$(GREEN)Running tests with Docker PostgreSQL...$(NC)"
	@$(MAKE) test-all
	@echo "$(GREEN)Tests completed!$(NC)"

docker-integration: docker-up ## Start PostgreSQL and run integration tests
	@echo "$(GREEN)Running integration tests with Docker PostgreSQL...$(NC)"
	@$(MAKE) test-integration
	@echo "$(GREEN)Integration tests completed!$(NC)"

docker-test-full: ## Run tests in Docker container (isolated environment)
	@echo "$(GREEN)Running tests in Docker container...$(NC)"
	docker-compose up --build --abort-on-container-exit test-runner
	@echo "$(GREEN)Docker tests completed!$(NC)"

docker-logs: ## Show Docker container logs
	docker-compose logs -f

docker-ps: ## Show running Docker containers
	docker-compose ps

db-shell: ## Open PostgreSQL shell
	@echo "$(GREEN)Opening PostgreSQL shell...$(NC)"
	docker-compose exec postgres psql -U postgres -d outbox_test

db-reset: ## Reset the test database
	@echo "$(GREEN)Resetting test database...$(NC)"
	docker-compose exec -T postgres psql -U postgres -c "DROP DATABASE IF EXISTS outbox_test;"
	docker-compose exec -T postgres psql -U postgres -c "CREATE DATABASE outbox_test;"
	@echo "$(GREEN)Database reset completed!$(NC)"

verify: fmt vet build test-unit ## Run all verification steps (fmt, vet, build, test)
	@echo "$(GREEN)All verification steps completed successfully!$(NC)"

ci: install verify test-integration ## Run full CI pipeline
	@echo "$(GREEN)CI pipeline completed successfully!$(NC)"

quick: ## Quick test (unit tests only)
	@./run_tests.sh

all: clean install build verify test-all ## Run everything from scratch
	@echo "$(GREEN)All tasks completed successfully!$(NC)"

# Info targets
info: ## Show configuration information
	@echo "$(GREEN)Configuration:$(NC)"
	@echo "  TEST_DB_HOST:     $(TEST_DB_HOST)"
	@echo "  TEST_DB_PORT:     $(TEST_DB_PORT)"
	@echo "  TEST_DB_USER:     $(TEST_DB_USER)"
	@echo "  TEST_DB_NAME:     $(TEST_DB_NAME)"
	@echo ""
	@echo "$(GREEN)Go Environment:$(NC)"
	@go version
	@echo "  GOPATH: $$(go env GOPATH)"
	@echo "  GOROOT: $$(go env GOROOT)"

version: ## Show version information
	@echo "Outbox Pattern Implementation"
	@echo "Go version: $$(go version)"
	@echo "Docker version: $$(docker --version 2>/dev/null || echo 'Not installed')"
	@echo "Docker Compose version: $$(docker-compose --version 2>/dev/null || echo 'Not installed')"

