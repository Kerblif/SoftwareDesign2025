package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"cbd/internal/payments/handler"
	"cbd/internal/payments/messaging"
	"cbd/internal/payments/service"
	"cbd/internal/payments/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	// Get configuration from environment variables
	listenAddr := getEnv("LISTEN_ADDR", ":8002")
	databaseURL := getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/payments_db?sslmode=disable")
	kafkaBrokers := strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ",")

	// Create repository
	repo, err := storage.NewPostgresRepository(databaseURL)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	// Create service
	paymentSvc := service.NewService(repo)

	// Create Kafka consumer
	consumer := messaging.NewConsumer(kafkaBrokers, "payments-service", paymentSvc, repo)

	// Create HTTP handler
	h := handler.NewHandler(paymentSvc)

	// Create router
	router := gin.Default()
	h.RegisterRoutes(router)

	// Start Kafka consumer
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := consumer.Start(ctx); err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}
	defer consumer.Stop()

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