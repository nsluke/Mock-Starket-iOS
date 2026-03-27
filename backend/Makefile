.PHONY: run build test lint migrate seed clean

run: ## Run the server locally
	go run cmd/server/main.go

build: ## Build the server binary
	go build -o bin/server cmd/server/main.go

test: ## Run all tests
	go test ./... -race -coverprofile=coverage.out -count=1

lint: ## Run linter
	golangci-lint run ./...

migrate-up: ## Run migrations
	go run cmd/migrate/main.go up

migrate-down: ## Rollback last migration
	go run cmd/migrate/main.go down 1

seed: ## Seed the database
	go run cmd/seed/main.go

clean: ## Clean build artifacts
	rm -rf bin/ tmp/ coverage.out

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
