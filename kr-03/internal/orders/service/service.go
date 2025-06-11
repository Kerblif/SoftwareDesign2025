package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cbd/internal/orders/models"
	"cbd/internal/orders/storage"

	"github.com/google/uuid"
)

// Constants for Kafka topics
const (
	PaymentRequestTopic = "payment-requests"
	PaymentResultTopic  = "payment-results"
)

// OrderService defines the interface for the order service
type OrderService interface {
	// Order operations
	CreateOrder(ctx context.Context, userID int64, amount float64, description string) (*models.Order, error)
	GetOrderByID(ctx context.Context, id int64) (*models.Order, error)
	GetOrdersByUserID(ctx context.Context, userID int64) ([]*models.Order, error)

	// Payment result processing
	ProcessPaymentResult(ctx context.Context, result *models.PaymentResult) error
}

// Service implements OrderService interface
type Service struct {
	repo storage.Repository
}

// NewService creates a new order service
func NewService(repo storage.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateOrder creates a new order and initiates the payment process
// It implements the transactional outbox pattern to ensure message delivery
func (s *Service) CreateOrder(ctx context.Context, userID int64, amount float64, description string) (*models.Order, error) {
	// Start a transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = s.repo.RollbackTx(tx)
		}
	}()

	// Create a new order
	now := time.Now()
	order := &models.Order{
		UserID:      userID,
		Amount:      amount,
		Description: description,
		Status:      models.OrderStatusNew,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Save the order to the database
	if err := s.repo.CreateOrder(ctx, tx, order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Create a payment message
	paymentMessage := &models.PaymentMessage{
		OrderID:     order.ID,
		UserID:      order.UserID,
		Amount:      order.Amount,
		Description: order.Description,
	}

	// Marshal the payment message to JSON
	paymentMessageBytes, err := json.Marshal(paymentMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payment message: %w", err)
	}

	// Create an outbox message
	outboxMessage := &models.OutboxMessage{
		MessageID: uuid.New().String(),
		Topic:     PaymentRequestTopic,
		Key:       fmt.Sprintf("%d", order.ID),
		Value:     paymentMessageBytes,
		Sent:      false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Save the outbox message to the database
	if err := s.repo.SaveOutboxMessage(ctx, tx, outboxMessage); err != nil {
		return nil, fmt.Errorf("failed to save outbox message: %w", err)
	}

	// Commit the transaction
	if err := s.repo.CommitTx(tx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return order, nil
}

// GetOrderByID retrieves an order by ID
func (s *Service) GetOrderByID(ctx context.Context, id int64) (*models.Order, error) {
	return s.repo.GetOrderByID(ctx, id)
}

// GetOrdersByUserID retrieves orders by user ID
func (s *Service) GetOrdersByUserID(ctx context.Context, userID int64) ([]*models.Order, error) {
	return s.repo.GetOrdersByUserID(ctx, userID)
}

// ProcessPaymentResult processes a payment result
func (s *Service) ProcessPaymentResult(ctx context.Context, result *models.PaymentResult) error {
	// Update the order status based on the payment result
	var status models.OrderStatus
	if result.Status == "success" {
		status = models.OrderStatusFinished
	} else {
		status = models.OrderStatusCancelled
	}

	// Update the order status in the database
	if err := s.repo.UpdateOrderStatus(ctx, result.OrderID, status); err != nil {
		return err
	}

	// Get the updated order to send to WebSocket clients
	_, err := s.repo.GetOrderByID(ctx, result.OrderID)
	if err != nil {
		return err
	}

	// Notify WebSocket clients (this will be handled by the caller)
	return nil
}
