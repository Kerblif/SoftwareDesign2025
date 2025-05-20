# Generate proto files
proto:
	bash scripts/generate_protos.sh

# Download Google proto files
proto-download:
	bash scripts/download_protos.sh

# Generate Swagger documentation
swagger: proto
	@echo "Swagger documentation generated in api/swagger directory"
	@echo "Run 'docker-compose up swagger-ui' to start the Swagger UI"
	@echo "Access the Swagger UI at http://localhost:8081"

# Start Swagger UI
swagger-ui:
	docker-compose up -d swagger-ui
	@echo "Swagger UI is running at http://localhost:8081"
