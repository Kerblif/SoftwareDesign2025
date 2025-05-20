package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"kr-02/cmd/api_gateway/server"
	"kr-02/internal/pkg/api_gateway/clients"
	pb "kr-02/internal/proto/api_gateway"
)

func main() {
	log.Println("Starting API Gateway...")

	// Initialize File Storing Service client
	fileStoringAddress := os.Getenv("FILE_STORING_SERVICE_ADDRESS")
	if fileStoringAddress == "" {
		fileStoringAddress = "file-storing-service:50051"
		log.Println("FILE_STORING_SERVICE_ADDRESS not set, using default:", fileStoringAddress)
	}

	fileStoringClient, err := clients.NewFileStoringClient(fileStoringAddress)
	if err != nil {
		log.Fatalf("Failed to initialize File Storing Service client: %v", err)
	}
	defer fileStoringClient.Close()

	// Initialize File Analysis Service client
	fileAnalysisAddress := os.Getenv("FILE_ANALYSIS_SERVICE_ADDRESS")
	if fileAnalysisAddress == "" {
		fileAnalysisAddress = "file-analysis-service:50052"
		log.Println("FILE_ANALYSIS_SERVICE_ADDRESS not set, using default:", fileAnalysisAddress)
	}

	fileAnalysisClient, err := clients.NewFileAnalysisClient(fileAnalysisAddress)
	if err != nil {
		log.Fatalf("Failed to initialize File Analysis Service client: %v", err)
	}
	defer fileAnalysisClient.Close()

	// Initialize server
	apiServer := server.NewServer(fileStoringClient, fileAnalysisClient)

	// Start gRPC server
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50050"
		log.Println("GRPC_PORT not set, using default:", grpcPort)
	}

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAPIGatewayServer(grpcServer, apiServer)

	// Start gRPC server in a goroutine
	go func() {
		log.Printf("Starting gRPC server on port %s...", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Start HTTP server with gRPC-Gateway
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
		log.Println("HTTP_PORT not set, using default:", httpPort)
	}

	// Create a client connection to the gRPC server
	conn, err := grpc.DialContext(
		context.Background(),
		"localhost:"+grpcPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to dial gRPC server: %v", err)
	}
	defer conn.Close()

	// Create a new ServeMux for the HTTP server
	gwmux := runtime.NewServeMux()

	// Register the gateway handlers
	err = pb.RegisterAPIGatewayHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	// Create a new HTTP server
	gwServer := &http.Server{
		Addr:    ":" + httpPort,
		Handler: gwmux,
	}

	// Start HTTP server
	log.Printf("Starting HTTP server on port %s...", httpPort)
	if err := gwServer.ListenAndServe(); err != nil {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}
}