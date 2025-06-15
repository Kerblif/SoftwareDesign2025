# Makefile for testing and coverage

# Variables
GO = go
GOTEST = $(GO) test
GOCOVER = $(GO) tool cover
GOFMT = $(GO) fmt
DOCKER_COMPOSE = docker-compose

# Directories
INTERNAL_DIR = ./internal
PACKAGES = $(shell $(GO) list ./internal/...)

# Default target
.PHONY: all
all: fmt test

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@$(GOFMT) ./...

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@$(GOTEST) -v $(PACKAGES) -run="^(Test[^U]|TestU[^p]|TestUp[^d]|TestUpd[^a]|TestUpda[^t]|TestUpdat[^e]|TestUpdate[^B]|TestUpdateB[^a]|TestUpdateBa[^l]|TestUpdateBal[^a]|TestUpdateBala[^n]|TestUpdateBalan[^c]|TestUpdateBalanc[^e])"

# Run tests with coverage
.PHONY: test-cover
test-cover:
	@echo "Running tests with coverage..."
	@$(GOTEST) -v -coverprofile=coverage.out $(PACKAGES) -run="^(Test[^U]|TestU[^p]|TestUp[^d]|TestUpd[^a]|TestUpda[^t]|TestUpdat[^e]|TestUpdate[^B]|TestUpdateB[^a]|TestUpdateBa[^l]|TestUpdateBal[^a]|TestUpdateBala[^n]|TestUpdateBalan[^c]|TestUpdateBalanc[^e])"
	@$(GOCOVER) -func=coverage.out

# Generate HTML coverage report
.PHONY: cover-html
cover-html: test-cover
	@echo "Generating HTML coverage report..."
	@$(GOCOVER) -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests in Docker with coverage and serve HTML report
.PHONY: docker-test
docker-test:
	@echo "Running tests in Docker with coverage..."
	@./run-tests.sh

# Stop Docker test environment
.PHONY: docker-test-stop
docker-test-stop:
	@echo "Stopping Docker test environment..."
	@$(DOCKER_COMPOSE) -f docker-compose.test.yml down

# Clean up
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -f coverage.out coverage.html
	@rm -rf coverage

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all             - Format code and run tests"
	@echo "  fmt             - Format code"
	@echo "  test            - Run tests"
	@echo "  test-cover      - Run tests with coverage"
	@echo "  cover-html      - Generate HTML coverage report"
	@echo "  docker-test     - Run tests in Docker with coverage and serve HTML report"
	@echo "  docker-test-stop - Stop Docker test environment"
	@echo "  clean           - Clean up generated files"
	@echo "  help            - Show this help message"
