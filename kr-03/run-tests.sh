#!/bin/bash
set -e

# Create coverage directory if it doesn't exist
mkdir -p coverage

# Check if docker-compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "docker-compose is not installed. Please install it first."
    exit 1
fi

# Build and start the test environment
echo "Building and starting the test environment..."
docker-compose -f docker-compose.test.yml up -d

# Wait for services to be ready
echo "Waiting for services to be ready..."
sleep 10

# Run tests with coverage
echo "Running tests with coverage..."
docker-compose -f docker-compose.test.yml exec test-runner go test -v -coverprofile=/app/coverage/coverage.out ./internal/...

# Generate HTML coverage report
echo "Generating HTML coverage report..."
docker-compose -f docker-compose.test.yml exec test-runner go tool cover -html=/app/coverage/coverage.out -o /app/coverage/coverage.html

# Copy the coverage report to the host
echo "Copying coverage report to host..."
docker cp test-runner:/app/coverage/coverage.html ./coverage/

# Set permissions for the coverage directory
echo "Setting permissions for coverage directory..."
chmod -R 755 ./coverage

echo "Coverage report generated at ./coverage/coverage.html"
echo "You can view the coverage report at http://localhost:8080/coverage.html"

# Keep the services running to serve the coverage report
echo "Services are still running to serve the coverage report."
echo "To stop the services, run: docker-compose -f docker-compose.test.yml down"
