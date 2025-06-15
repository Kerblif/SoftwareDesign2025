package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"cbd/internal/orders/models"
	"cbd/internal/orders/service"
	"cbd/internal/orders/storage"

	"github.com/segmentio/kafka-go"
)

// Producer represents a Kafka producer for payment requests
type Producer struct {
	writer       *kafka.Writer
	reader       *kafka.Reader
	repo         storage.Repository
	orderSvc     service.OrderService
	wsHandler    WebSocketNotifier
	outboxPoller *OutboxPoller
	stopCh       chan struct{}
	readerStopCh chan struct{}
}

// WebSocketNotifier defines the interface for notifying WebSocket clients
type WebSocketNotifier interface {
	NotifyOrderStatusChange(ctx context.Context, order *models.Order)
}

// NewProducer creates a new Kafka producer
func NewProducer(brokers []string, groupID string, repo storage.Repository, orderSvc service.OrderService, wsHandler WebSocketNotifier) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Balancer: &kafka.LeastBytes{},
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		GroupID:     groupID,
		Topic:       service.PaymentResultTopic,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		MaxWait:     1 * time.Second,
		StartOffset: kafka.FirstOffset,
	})

	return &Producer{
		writer:       writer,
		reader:       reader,
		repo:         repo,
		orderSvc:     orderSvc,
		wsHandler:    wsHandler,
		stopCh:       make(chan struct{}),
		readerStopCh: make(chan struct{}),
	}
}

// Start starts the producer and outbox poller
func (p *Producer) Start(ctx context.Context) error {
	// Start the outbox poller
	p.outboxPoller = NewOutboxPoller(p.repo, p.writer)
	go p.outboxPoller.Start(ctx, p.stopCh)

	// Start consuming payment results
	go p.consumePaymentResults(ctx)

	return nil
}

// Stop stops the producer and outbox poller
func (p *Producer) Stop() {
	close(p.stopCh)
	close(p.readerStopCh)
	p.writer.Close()
	p.reader.Close()
}

// consumePaymentResults reads payment results from Kafka and processes them
func (p *Producer) consumePaymentResults(ctx context.Context) {
	for {
		select {
		case <-p.readerStopCh:
			return
		default:
			message, err := p.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading payment result: %v", err)
				continue
			}

			// Process the payment result
			if err := p.processPaymentResult(ctx, message); err != nil {
				log.Printf("Error processing payment result: %v", err)
			}
		}
	}
}

// processPaymentResult processes a payment result message
func (p *Producer) processPaymentResult(ctx context.Context, message kafka.Message) error {
	// Parse the payment result
	var result models.PaymentResult
	if err := json.Unmarshal(message.Value, &result); err != nil {
		return fmt.Errorf("failed to unmarshal payment result: %w", err)
	}

	// Process the payment result
	if err := p.orderSvc.ProcessPaymentResult(ctx, &result); err != nil {
		return err
	}

	// Get the updated order to send to WebSocket clients
	order, err := p.repo.GetOrderByID(ctx, result.OrderID)
	if err != nil {
		return err
	}

	// Notify WebSocket clients
	if p.wsHandler != nil {
		p.wsHandler.NotifyOrderStatusChange(ctx, order)
	}

	return nil
}

// OutboxPoller polls the outbox table for unsent messages
type OutboxPoller struct {
	repo   storage.Repository
	writer *kafka.Writer
}

// NewOutboxPoller creates a new outbox poller
func NewOutboxPoller(repo storage.Repository, writer *kafka.Writer) *OutboxPoller {
	return &OutboxPoller{
		repo:   repo,
		writer: writer,
	}
}

// Start starts the outbox poller
func (p *OutboxPoller) Start(ctx context.Context, stopCh <-chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stopCh:
			return
		case <-ticker.C:
			if err := p.pollOutbox(ctx); err != nil {
				log.Printf("Error polling outbox: %v", err)
			}
		}
	}
}

// pollOutbox polls the outbox table for unsent messages
func (p *OutboxPoller) pollOutbox(ctx context.Context) error {
	// Get unsent messages from the outbox
	messages, err := p.repo.GetUnsendOutboxMessages(ctx, 10)
	if err != nil {
		return fmt.Errorf("failed to get unsent outbox messages: %w", err)
	}

	// Send each message to Kafka
	for _, message := range messages {
		// Send the message to Kafka
		kafkaMsg := kafka.Message{
			Key:   []byte(message.Key),
			Value: message.Value,
		}

		// Only set the topic if it's not already set in the writer
		if p.writer.Topic == "" {
			kafkaMsg.Topic = message.Topic
		}

		err := p.writer.WriteMessages(ctx, kafkaMsg)
		if err != nil {
			log.Printf("Error sending message to Kafka: %v", err)
			continue
		}

		// Mark the message as sent
		if err := p.repo.MarkOutboxMessageAsSent(ctx, message.MessageID); err != nil {
			log.Printf("Error marking outbox message as sent: %v", err)
			continue
		}
	}

	return nil
}
