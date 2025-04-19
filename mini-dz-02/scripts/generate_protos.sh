#!/bin/bash

# Create the output directory if it doesn't exist
mkdir -p ./internal/proto

# Generate Go code from proto files
protoc -I ./proto \
   --go_out ./internal/proto --go_opt paths=source_relative \
   --go-grpc_out ./internal/proto --go-grpc_opt paths=source_relative \
   --grpc-gateway_out ./internal/proto --grpc-gateway_opt paths=source_relative \
   proto/zoo/*.proto
