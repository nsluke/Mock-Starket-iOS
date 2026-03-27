.PHONY: dev dev-backend dev-web test-all test-backend test-ios test-android test-web lint clean setup

# ============================
# Development
# ============================

dev: ## Start full stack (backend + postgres + web)
	cd deploy && docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build

dev-backend: ## Start backend + postgres only
	cd deploy && docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build backend postgres

dev-web: ## Start web frontend in dev mode
	cd web && npm run dev

# ============================
# Testing
# ============================

test-all: test-backend test-web ## Run all test suites

test-backend: ## Run Go backend tests
	cd backend && go test ./... -race -coverprofile=coverage.out
	@echo "Backend coverage:"
	@cd backend && go tool cover -func=coverage.out | tail -1

test-ios: ## Run iOS tests (requires Xcode)
	cd ios && xcodebuild test \
		-scheme MockStarket \
		-destination 'platform=iOS Simulator,name=iPhone 16' \
		-resultBundlePath TestResults \
		| xcpretty

test-android: ## Run Android tests
	cd android && ./gradlew test

test-web: ## Run web frontend tests
	cd web && npm run test

# ============================
# Linting
# ============================

lint: lint-backend lint-web ## Lint all code

lint-backend: ## Lint Go code
	cd backend && golangci-lint run ./...

lint-web: ## Lint web code
	cd web && npm run lint

# ============================
# Database
# ============================

db-migrate: ## Run database migrations
	cd backend && go run cmd/migrate/main.go up

db-rollback: ## Rollback last migration
	cd backend && go run cmd/migrate/main.go down 1

db-seed: ## Seed database with stock data
	cd backend && go run cmd/seed/main.go

db-reset: db-rollback db-migrate db-seed ## Reset database

# ============================
# Build
# ============================

build-backend: ## Build Go backend binary
	cd backend && go build -o bin/server cmd/server/main.go

build-web: ## Build web frontend for production
	cd web && npm run build

# ============================
# Docker
# ============================

docker-build: ## Build all Docker images
	cd deploy && docker compose build

docker-up: ## Start production stack
	cd deploy && docker compose up -d

docker-down: ## Stop all containers
	cd deploy && docker compose down

docker-logs: ## Tail all container logs
	cd deploy && docker compose logs -f

# ============================
# Setup
# ============================

setup: ## First-time setup for development
	@echo "Setting up Mock Starket development environment..."
	cd backend && go mod download
	cd web && npm ci
	cp deploy/.env.example deploy/.env
	@echo "Done! Run 'make dev' to start."

clean: ## Clean build artifacts
	cd backend && rm -rf bin/ tmp/ coverage.out
	cd web && rm -rf .next/ out/ node_modules/
	cd deploy && docker compose down -v

# ============================
# Help
# ============================

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
