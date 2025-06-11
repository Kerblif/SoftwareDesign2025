package messaging

import (
	"context"
	"fmt"
	"log"
	"time"

	"cbd/internal/payments/models"
	"cbd/internal/payments/service"
	"cbd/internal/payments/storage"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

// Consumer represents a Kafka consumer for payment requests
type Consumer struct {
	reader        *kafka.Reader
	paymentSvc    service.PaymentService
	repo          storage.Repository
	inboxPoller   *InboxPoller
	outboxPoller  *OutboxPoller
	stopCh        chan struct{}
	inboxStopCh   chan struct{}
	outboxStopCh  chan struct{}
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(brokers []string, groupID string, paymentSvc service.PaymentService, repo storage.Repository) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		GroupID:     groupID,
		Topic:       service.PaymentRequestTopic,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		MaxWait:     1 * time.Second,
		StartOffset: kafka.FirstOffset,
	})

	return &Consumer{
		reader:       reader,
		paymentSvc:   paymentSvc,
		repo:         repo,
		stopCh:       make(chan struct{}),
		inboxStopCh:  make(chan struct{}),
		outboxStopCh: make(chan struct{}),
	}
}

// Start starts the consumer and pollers
func (c *Consumer) Start(ctx context.Context) error {
	// Start the inbox poller
	c.inboxPoller = NewInboxPoller(c.repo, c.paymentSvc)
	go c.inboxPoller.Start(ctx, c.inboxStopCh)

	// Start the outbox poller
	c.outboxPoller = NewOutboxPoller(c.repo, []string{c.reader.Config().Brokers[0]})
	go c.outboxPoller.Start(ctx, c.outboxStopCh)

	// Start consuming messages
	go c.consume(ctx)

	return nil
}

// Stop stops the consumer and pollers
func (c *Consumer) Stop() {
	close(c.stopCh)
	close(c.inboxStopCh)
	close(c.outboxStopCh)
	c.reader.Close()
}

// consume reads messages from Kafka and processes them
func (c *Consumer) consume(ctx context.Context) {
	for {
		select {
		case <-c.stopCh:
			return
		default:
			message, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}

			// Process the message
			if err := c.processMessage(ctx, message); err != nil {
				log.Printf("Error processing message: %v", err)
			}
		}
	}
}

// processMessage processes a Kafka message by saving it to the inbox
func (c *Consumer) processMessage(ctx context.Context, message kafka.Message) error {
	// Create an inbox message
	inboxMessage := &models.InboxMessage{
		MessageID: uuid.New().String(),
		Topic:     message.Topic,
		Key:       string(message.Key),
		Value:     message.Value,
		Processed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save the message to the inbox
	return c.repo.SaveInboxMessage(ctx, nil, inboxMessage)
}

// InboxPoller polls the inbox table for unprocessed messages
type InboxPoller struct {
	repo       storage.Repository
	paymentSvc service.PaymentService
}

// NewInboxPoller creates a new inbox poller
func NewInboxPoller(repo storage.Repository, paymentSvc service.PaymentService) *InboxPoller {
	return &InboxPoller{
		repo:       repo,
		paymentSvc: paymentSvc,
	}
}

// Start starts the inbox poller
func (p *InboxPoller) Start(ctx context.Context, stopCh <-chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stopCh:
			return
		case <-ticker.C:
			if err := p.pollInbox(ctx); err != nil {
				log.Printf("Error polling inbox: %v", err)
			}
		}
	}
}

// pollInbox polls the inbox table for unprocessed messages
func (p *InboxPoller) pollInbox(ctx context.Context) error {
	// Get unprocessed messages from the inbox
	messages, err := p.repo.GetUnprocessedInboxMessages(ctx, 10)
	if err != nil {
		return fmt.Errorf("failed to get unprocessed inbox messages: %w", err)
	}

	// Process each message
	for _, message := range messages {
		if err := p.paymentSvc.ProcessInboxMessage(ctx, message); err != nil {
			log.Printf("Error processing inbox message: %v", err)
			continue
		}
	}

	return nil
}

// OutboxPoller polls the outbox table for unsent messages
type OutboxPoller struct {
	repo    storage.Repository
	writer  *kafka.Writer
}

// NewOutboxPoller creates a new outbox poller
func NewOutboxPoller(repo storage.Repository, brokers []string) *OutboxPoller {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Balancer: &kafka.LeastBytes{},
	}

	return &OutboxPoller{
		repo:   repo,
		writer: writer,
	}
}

// Start starts the outbox poller
func (p *OutboxPoller) Start(ctx context.Context, stopCh <-chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	defer p.writer.Close()

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
		err := p.writer.WriteMessages(ctx, kafka.Message{
			Topic: message.Topic,
			Key:   []byte(message.Key),
			Value: message.Value,
		})
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
