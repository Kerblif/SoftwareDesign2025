# Generate proto files
proto:
	bash scripts/generate_protos.sh

# Download Google proto files
proto-download:
	bash scripts/download_protos.sh

# Generate Swagger documentation
swagger:
	mkdir -p internal/pkg/api_gateway/docs
	swag init -g cmd/api_gateway/main.go -o internal/pkg/api_gateway/docs
