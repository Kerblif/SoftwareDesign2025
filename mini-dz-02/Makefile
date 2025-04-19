# Makefile for Zoo Management System

# Variables
GO=go
BINARY_NAME=zoo-management
PROTO_DIR=proto
PROTO_GEN_DIR=internal/proto
MAIN_FILE=main.go
GRPC_SERVER_DIR=cmd/grpc-server
DOCKER_IMAGE=zoo-management
DOCKER_CONTAINER=zoo-management-container

# Build the application
build:
	$(GO) build -o $(BINARY_NAME) $(MAIN_FILE)

# Run the application
run:
	$(GO) run $(MAIN_FILE)

# Run the gRPC server directly
run-grpc:
	cd $(GRPC_SERVER_DIR) && $(GO) run .

# Generate proto files
proto:
	bash scripts/generate_protos.sh

# Download Google proto files
proto-download:
	bash scripts/download_protos.sh

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -rf $(PROTO_GEN_DIR)

# Run tests
test:
	$(GO) test ./...

# Run tests with coverage
test-coverage:
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out

# Run tests with coverage excluding generated files
test-all:
	@echo "Running all tests with verbose output (excluding generated files)..."
	$(GO) test -v ./pkg/...
	@echo "\nGenerating coverage report (excluding generated files)..."
	$(GO) test -coverprofile=coverage.out ./pkg/...
	@echo "\nCoverage statistics by function (excluding generated files):"
	$(GO) tool cover -func=coverage.out
	@echo "\nGenerating HTML coverage report (excluding generated files)..."
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "\nTest coverage report generated: coverage.html"

# Format code
fmt:
	$(GO) fmt ./...

# Generate documentation (excluding generated files)
doc:
	@echo "Files to be included in documentation (excluding generated files):"
	@find ./pkg -type f -name "*.go" | grep -v "_test.go" | grep -v "\.pb\.go" | grep -v "\.pb\.gw\.go"

# Vet code
vet:
	$(GO) vet ./...

# Install dependencies
deps:
	$(GO) mod download

# Build and run the application
all: proto build run

# Docker targets
docker-build:
	docker build -t $(DOCKER_IMAGE) .

docker-run:
	docker run --name $(DOCKER_CONTAINER) -p 8080:8080 -p 9090:9090 -d $(DOCKER_IMAGE)

docker-stop:
	docker stop $(DOCKER_CONTAINER) || true
	docker rm $(DOCKER_CONTAINER) || true

docker-clean: docker-stop
	docker rmi $(DOCKER_IMAGE) || true

# Build and run with Docker
docker-all: docker-build docker-run
	@echo "Docker container started. Access the API at http://localhost:8080"
	@echo "To stop the container, run: make docker-stop"

# Default target
.PHONY: build run run-grpc proto proto-download clean test test-coverage test-all test-all-no-gen doc fmt vet deps all docker-build docker-run docker-stop docker-clean docker-all
