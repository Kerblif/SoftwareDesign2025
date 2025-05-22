# Text Analyzer Microservice Architecture

This project implements a microservice architecture for analyzing text files, including statistics calculation, plagiarism detection, and word cloud generation.

## Architecture

The system consists of three microservices:

1. **API Gateway** - Responsible for routing requests to the appropriate services
2. **File Storing Service** - Responsible for storing and retrieving files
3. **File Analysis Service** - Responsible for analyzing files and storing results

## Prerequisites

- Docker and Docker Compose
- Go 1.20 or later
- Swag (for generating Swagger documentation)
  ```bash
  go install github.com/swaggo/swag/cmd/swag@latest
  ```

## Getting Started

### Clone the Repository

```bash
git clone <repository-url>
cd kr-02
```

### Generate Swagger Documentation

```bash
make swagger
```

### Build and Run the Services

```bash
docker-compose up -d
```

This will start all the services:
- PostgreSQL database on port 5432
- File Storing Service on port 50051
- File Analysis Service on port 50052
- API Gateway on port 8080 (HTTP)
- Swagger UI on port 8081

## API Documentation

### Generate Swagger Documentation

```bash
make swagger
```

This will generate Swagger documentation using swag based on the annotations in the code.

### Access Swagger UI

```bash
make swagger-ui
```

Then open http://localhost:8081 in your browser to view the API documentation.

You can also access the Swagger UI directly from the API Gateway at http://localhost:8080/swagger/index.html.

## API Endpoints

### Upload a File

```
POST /api/v1/files
```

Request: multipart/form-data with a file field named "file"

Example using curl:
```bash
curl -X POST -F "file=@example.txt" http://localhost:8080/api/v1/files
```

Response:
```json
{
  "file_id": "unique-file-id"
}
```

### Get a File

```
GET /api/v1/files/{file_id}
```

Response: Binary file content with appropriate Content-Disposition header for download

Example using curl:
```bash
curl -OJ http://localhost:8080/api/v1/files/{file_id}
```

### Analyze a File

```
POST /api/v1/analysis
```

Request body:
```json
{
  "file_id": "unique-file-id",
  "generate_word_cloud": true
}
```

Response:
```json
{
  "paragraph_count": 5,
  "word_count": 100,
  "character_count": 500,
  "is_plagiarism": false,
  "similar_file_ids": [],
  "word_cloud_location": "word-cloud-location"
}
```

### Get a Word Cloud

```
GET /api/v1/wordcloud/{location}
```

Response: Word cloud image (binary data) with Content-Type: image/png

Example using curl:
```bash
curl -o wordcloud.png http://localhost:8080/api/v1/wordcloud/{location}
```

## Development

### Project Structure

```
kr-02/
├── build/                    # Dockerfiles
│   ├── api_gateway/
│   ├── file_analysis_service/
│   └── file_storing_service/
├── cmd/                      # Entry points for each service
│   ├── api_gateway/
│   ├── file_analysis_service/
│   └── file_storing_service/
├── configs/                  # Configuration files
├── internal/                 # Internal packages
│   ├── pkg/                  # Shared packages
│   │   ├── api_gateway/      # API Gateway implementation
│   │   │   ├── clients/      # Service clients
│   │   │   ├── docs/         # Generated Swagger docs
│   │   │   └── handlers/     # HTTP handlers
│   │   ├── file_analysis/    # File Analysis Service implementation
│   │   └── file_storing/     # File Storing Service implementation
│   └── proto/                # Generated proto files for gRPC services
├── proto/                    # Proto definitions for gRPC services
├── scripts/                  # Utility scripts
└── tests/                    # Tests
```

### Running Tests

```bash
go test ./tests/...
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
