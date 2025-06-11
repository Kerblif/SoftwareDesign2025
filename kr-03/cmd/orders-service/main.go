package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"cbd/internal/orders/handler"
	"cbd/internal/orders/messaging"
	"cbd/internal/orders/service"
	"cbd/internal/orders/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	// Get configuration from environment variables
	listenAddr := getEnv("LISTEN_ADDR", ":8001")
	databaseURL := getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/orders_db?sslmode=disable")
	kafkaBrokers := strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ",")

	// Create repository
	repo, err := storage.NewPostgresRepository(databaseURL)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	// Create service
	orderSvc := service.NewService(repo)

	// Create WebSocket handler
	wsHandler := handler.NewWebSocketHandler()

	// Create Kafka producer
	producer := messaging.NewProducer(kafkaBrokers, "orders-service", repo, orderSvc, wsHandler)

	// Create HTTP handler
	h := handler.NewHandler(orderSvc)

	// Create router
	router := gin.Default()
	h.RegisterRoutes(router)
	wsHandler.RegisterRoutes(router)

	// Start Kafka producer
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := producer.Start(ctx); err != nil {
		log.Fatalf("Failed to start producer: %v", err)
	}
	defer producer.Stop()

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := router.Run(listenAddr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")
	cancel()
	log.Println("Server exited")
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
