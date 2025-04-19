# Zoo Management System

A web application for managing a zoo, including animals, enclosures, and feeding schedules. The application is built using Domain-Driven Design and Clean Architecture principles.

## Features

- Manage animals (add, remove, transfer between enclosures, treat)
- Manage enclosures (add, remove, clean)
- Manage feeding schedules (add, update, mark as completed)
- View zoo statistics

## Running with Docker

The easiest way to run the application is using Docker, which doesn't require Go to be installed on your machine.

### Prerequisites

- Docker

### Build and Run

```bash
# Build the Docker image and run the container
make docker-all

# Stop the container
make docker-stop

# Remove the container and image
make docker-clean
```

The application will be available at:
- REST API: http://localhost:8080
- gRPC server: localhost:9090

## Running without Docker

If you prefer to run the application without Docker, you'll need to have Go and Protocol Buffers installed.

### Prerequisites

- Go 1.24 or later
- Protocol Buffers compiler (protoc)
- Go plugins for Protocol Buffers:
  - protoc-gen-go
  - protoc-gen-go-grpc
  - protoc-gen-grpc-gateway

### Build and Run

```bash
# Generate Protocol Buffers code
make proto

# Build the application
make build

# Run the application
make run

# Or do all of the above in one command
make all
```

## API Endpoints

The application provides both REST API and gRPC endpoints.

### REST API

- **Animals**
  - GET /api/animals - Get all animals
  - GET /api/animals/{id} - Get a specific animal
  - POST /api/animals - Create a new animal
  - DELETE /api/animals/{id} - Delete an animal
  - POST /api/animals/{id}/transfer - Transfer an animal to another enclosure
  - POST /api/animals/{id}/treat - Treat a sick animal

- **Enclosures**
  - GET /api/enclosures - Get all enclosures
  - GET /api/enclosures/{id} - Get a specific enclosure
  - POST /api/enclosures - Create a new enclosure
  - DELETE /api/enclosures/{id} - Delete an enclosure
  - GET /api/enclosures/{id}/animals - Get animals in an enclosure
  - POST /api/enclosures/{id}/clean - Clean an enclosure

- **Feeding Schedules**
  - GET /api/feeding-schedules - Get all feeding schedules
  - GET /api/feeding-schedules/{id} - Get a specific feeding schedule
  - POST /api/feeding-schedules - Create a new feeding schedule
  - DELETE /api/feeding-schedules/{id} - Delete a feeding schedule
  - PUT /api/feeding-schedules/{id} - Update a feeding schedule
  - POST /api/feeding-schedules/{id}/complete - Mark a feeding schedule as completed
  - GET /api/feeding-schedules/due - Get feeding schedules that are due
  - GET /api/animals/{id}/feeding-schedules - Get feeding schedules for a specific animal

- **Statistics**
  - GET /api/statistics - Get overall zoo statistics
  - GET /api/statistics/species - Get animal count by species
  - GET /api/statistics/enclosure-utilization - Get enclosure utilization
  - GET /api/statistics/health - Get health status statistics

### gRPC

The application also provides gRPC endpoints for all the above functionality. The gRPC server is available at localhost:9090.