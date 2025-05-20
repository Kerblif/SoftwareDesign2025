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
- Protocol Buffers compiler (protoc)
- jq (for merging Swagger files)

## Getting Started

### Clone the Repository

```bash
git clone <repository-url>
cd kr-02
```

### Download Google Proto Files

```bash
make proto-download
```

### Generate Proto Files

```bash
make proto
```

### Build and Run the Services

```bash
docker-compose up -d
```

This will start all the services:
- PostgreSQL database on port 5432
- File Storing Service on port 50051
- File Analysis Service on port 50052
- API Gateway on ports 50050 (gRPC) and 8080 (HTTP)
- Swagger UI on port 8081

## API Documentation

### Generate Swagger Documentation

```bash
make swagger
```

### Access Swagger UI

```bash
make swagger-ui
```

Then open http://localhost:8081 in your browser to view the API documentation.

## API Endpoints

### Upload a File

```
POST /api/v1/files
```

Request body:
```json
{
  "file_name": "example.txt",
  "content": "base64-encoded-content"
}
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

Response:
```json
{
  "file_name": "example.txt",
  "content": "base64-encoded-content"
}
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

Response: Word cloud image (binary data)

## Development

### Project Structure

```
kr-02/
├── api/                      # API documentation
│   └── swagger/              # Swagger files
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
│   └── proto/                # Generated proto files
├── proto/                    # Proto definitions
├── scripts/                  # Utility scripts
└── tests/                    # Tests
```

### Running Tests

```bash
go test ./tests/...
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.