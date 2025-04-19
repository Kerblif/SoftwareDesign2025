package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"mini-dz-02/cmd/grpc-server/servers"
	"mini-dz-02/internal/proto/zoo"
	"mini-dz-02/pkg/zoo/application"
	"mini-dz-02/pkg/zoo/infrastructure"
)

const (
	grpcPort = 9090
	httpPort = 8080
)

func main() {
	// Create repositories
	animalRepository := infrastructure.NewInMemoryAnimalRepository()
	enclosureRepository := infrastructure.NewInMemoryEnclosureRepository()
	feedingScheduleRepository := infrastructure.NewInMemoryFeedingScheduleRepository()
	eventPublisher := infrastructure.NewInMemoryEventPublisher()

	// Create application services
	animalTransferService := application.NewAnimalTransferService(
		animalRepository,
		enclosureRepository,
		eventPublisher,
	)

	feedingOrganizationService := application.NewFeedingOrganizationService(
		animalRepository,
		feedingScheduleRepository,
		eventPublisher,
	)

	zooStatisticsService := application.NewZooStatisticsService(
		animalRepository,
		enclosureRepository,
	)

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Create server implementations
	echoServer := servers.NewEchoServer()
	animalServer := servers.NewAnimalServer(animalRepository, enclosureRepository, animalTransferService)
	enclosureServer := servers.NewEnclosureServer(enclosureRepository, animalRepository)
	feedingScheduleServer := servers.NewFeedingScheduleServer(feedingScheduleRepository, animalRepository, feedingOrganizationService)
	statisticsServer := servers.NewStatisticsServer(zooStatisticsService)

	// Register services
	zoo.RegisterEchoServiceServer(grpcServer, echoServer)
	zoo.RegisterAnimalServiceServer(grpcServer, animalServer)
	zoo.RegisterEnclosureServiceServer(grpcServer, enclosureServer)
	zoo.RegisterFeedingScheduleServiceServer(grpcServer, feedingScheduleServer)
	zoo.RegisterStatisticsServiceServer(grpcServer, statisticsServer)

	// Register reflection service on gRPC server
	reflection.Register(grpcServer)

	// Start gRPC server in a goroutine
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
		log.Printf("Starting gRPC server on port %d...", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Create gRPC-Gateway mux
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	gwmux := runtime.NewServeMux()

	// Register gRPC-Gateway handlers
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := zoo.RegisterEchoServiceHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf("localhost:%d", grpcPort), opts); err != nil {
		log.Fatalf("Failed to register echo service: %v", err)
	}
	if err := zoo.RegisterAnimalServiceHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf("localhost:%d", grpcPort), opts); err != nil {
		log.Fatalf("Failed to register animal service: %v", err)
	}
	if err := zoo.RegisterEnclosureServiceHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf("localhost:%d", grpcPort), opts); err != nil {
		log.Fatalf("Failed to register enclosure service: %v", err)
	}
	if err := zoo.RegisterFeedingScheduleServiceHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf("localhost:%d", grpcPort), opts); err != nil {
		log.Fatalf("Failed to register feeding schedule service: %v", err)
	}
	if err := zoo.RegisterStatisticsServiceHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf("localhost:%d", grpcPort), opts); err != nil {
		log.Fatalf("Failed to register statistics service: %v", err)
	}

	// Create HTTP server with gRPC-Gateway mux
	mux := http.NewServeMux()
	mux.Handle("/api/", gwmux)

	// Add a simple welcome page
	mux.HandleFunc("/", ServeWelcomePage)

	// Start HTTP server
	log.Printf("Starting HTTP server on port %d...", httpPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", httpPort), mux))
}
