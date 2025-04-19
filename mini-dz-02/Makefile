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

# Format code
fmt:
	$(GO) fmt ./...

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
.PHONY: build run run-grpc proto proto-download clean test test-coverage fmt vet deps all docker-build docker-run docker-stop docker-clean docker-all
