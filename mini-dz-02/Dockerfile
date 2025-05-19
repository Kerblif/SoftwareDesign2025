FROM golang:1.24-alpine AS builder

# Install required tools
RUN apk add --no-cache bash protobuf protobuf-dev git

# Install protoc plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Generate proto files
RUN bash scripts/generate_protos.sh

# Build the application
RUN go build -o build cmd/grpc-server/*.go

# Create a minimal runtime image
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/build .

# Expose the ports
EXPOSE 8080 9090

# Set the entrypoint
ENTRYPOINT ["./build"]