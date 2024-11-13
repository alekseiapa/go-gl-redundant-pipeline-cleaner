.PHONY: help audit benchmark coverage format lint test tidy

# Default settings
GO ?= go
PKG := ./...
OUT := ./bin

help: ## Display available commands
	@echo "Makefile Commands:"
	@echo "make help       - Display available commands."
	@echo "make audit      - Conduct quality checks."
	@echo "make coverage   - Generate test coverage report."
	@echo "make format     - Automatically format code."
	@echo "make lint       - Run lint checks."
	@echo "make tidy       - Tidy dependencies."

audit: ## Conduct quality checks using go vet and gosec
	$(GO) vet $(PKG)
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	gosec $(PKG)

coverage: ## Generate test coverage report
	$(GO) test -coverprofile=coverage.out $(PKG)
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

format: ## Automatically format code using gofmt and goimports
	gofmt -w .
	go install golang.org/x/tools/cmd/goimports@latest
	goimports -w .

lint: ## Run lint checks using golangci-lint
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run $(PKG)

tidy: ## Tidy dependencies
	$(GO) mod tidy

clean: ## Clean up generated files
	rm -f coverage.out coverage.html

# Run this by simply typing 'make'
.DEFAULT_GOAL := help
