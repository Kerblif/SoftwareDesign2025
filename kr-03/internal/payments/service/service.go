package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"cbd/internal/payments/models"
	"cbd/internal/payments/storage"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Constants for Kafka topics
const (
	PaymentRequestTopic = "payment-requests"
	PaymentResultTopic  = "payment-results"
)

// PaymentService defines the interface for the payment service
type PaymentService interface {
	// Account operations
	CreateAccount(ctx context.Context, userID int64) (*models.Account, error)
	Deposit(ctx context.Context, userID int64, amount float64) error
	GetBalance(ctx context.Context, userID int64) (*models.BalanceResponse, error)
	
	// Payment operations
	ProcessPayment(ctx context.Context, payment *models.PaymentMessage) error
	
	// Inbox message processing
	ProcessInboxMessage(ctx context.Context, message *models.InboxMessage) error
}

// Service implements PaymentService interface
type Service struct {
	repo storage.Repository
}

// NewService creates a new payment service
func NewService(repo storage.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateAccount creates a new account for a user
func (s *Service) CreateAccount(ctx context.Context, userID int64) (*models.Account, error) {
	return s.repo.CreateAccount(ctx, userID)
}

// Deposit adds money to a user's account
func (s *Service) Deposit(ctx context.Context, userID int64, amount float64) error {
	if amount <= 0 {
		return errors.New("deposit amount must be positive")
	}
	
	return s.repo.UpdateBalance(ctx, userID, amount)
}

// GetBalance returns the current balance of a user's account
func (s *Service) GetBalance(ctx context.Context, userID int64) (*models.BalanceResponse, error) {
	account, err := s.repo.GetAccountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	return &models.BalanceResponse{
		UserID:  account.UserID,
		Balance: account.Balance,
	}, nil
}

// ProcessPayment processes a payment request
// It implements the transactional inbox pattern to ensure exactly-once processing
func (s *Service) ProcessPayment(ctx context.Context, payment *models.PaymentMessage) error {
	// Start a transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = s.repo.RollbackTx(tx)
		}
	}()
	
	// Debit the user's account
	account, err := s.repo.GetAccountByUserID(ctx, payment.UserID)
	if err != nil {
		return createPaymentResult(ctx, s.repo, tx, payment.OrderID, "fail", err.Error())
	}
	
	// Check if there's enough balance
	if account.Balance < payment.Amount {
		return createPaymentResult(ctx, s.repo, tx, payment.OrderID, "fail", "insufficient funds")
	}
	
	// Update the balance (debit)
	err = s.repo.UpdateBalance(ctx, payment.UserID, -payment.Amount)
	if err != nil {
		return createPaymentResult(ctx, s.repo, tx, payment.OrderID, "fail", err.Error())
	}
	
	// Create a success payment result
	err = createPaymentResult(ctx, s.repo, tx, payment.OrderID, "success", "")
	if err != nil {
		return err
	}
	
	// Commit the transaction
	return s.repo.CommitTx(tx)
}

// ProcessInboxMessage processes a message from the inbox
func (s *Service) ProcessInboxMessage(ctx context.Context, message *models.InboxMessage) error {
	// Start a transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = s.repo.RollbackTx(tx)
		}
	}()
	
	// Parse the payment message
	var payment models.PaymentMessage
	if err := json.Unmarshal(message.Value, &payment); err != nil {
		return fmt.Errorf("failed to unmarshal payment message: %w", err)
	}
	
	// Process the payment
	account, err := s.repo.GetAccountByUserID(ctx, payment.UserID)
	if err != nil {
		// Create a failure payment result
		if err := createPaymentResultWithTx(ctx, s.repo, tx, payment.OrderID, "fail", err.Error()); err != nil {
			return err
		}
		
		// Mark the message as processed
		if err := s.repo.MarkInboxMessageAsProcessed(ctx, tx, message.MessageID); err != nil {
			return err
		}
		
		return s.repo.CommitTx(tx)
	}
	
	// Check if there's enough balance
	if account.Balance < payment.Amount {
		// Create a failure payment result
		if err := createPaymentResultWithTx(ctx, s.repo, tx, payment.OrderID, "fail", "insufficient funds"); err != nil {
			return err
		}
		
		// Mark the message as processed
		if err := s.repo.MarkInboxMessageAsProcessed(ctx, tx, message.MessageID); err != nil {
			return err
		}
		
		return s.repo.CommitTx(tx)
	}
	
	// Update the balance (debit)
	newBalance := account.Balance - payment.Amount
	_, err = tx.ExecContext(
		ctx,
		"UPDATE accounts SET balance = $1, updated_at = $2 WHERE user_id = $3",
		newBalance, time.Now(), payment.UserID,
	)
	if err != nil {
		return err
	}
	
	// Create a success payment result
	if err := createPaymentResultWithTx(ctx, s.repo, tx, payment.OrderID, "success", ""); err != nil {
		return err
	}
	
	// Mark the message as processed
	if err := s.repo.MarkInboxMessageAsProcessed(ctx, tx, message.MessageID); err != nil {
		return err
	}
	
	// Commit the transaction
	return s.repo.CommitTx(tx)
}

// Helper function to create a payment result and save it to the outbox
func createPaymentResult(ctx context.Context, repo storage.Repository, tx *sqlx.Tx, orderID int64, status, reason string) error {
	result := models.PaymentResult{
		OrderID: orderID,
		Status:  status,
		Reason:  reason,
	}
	
	resultBytes, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal payment result: %w", err)
	}
	
	outboxMessage := &models.OutboxMessage{
		MessageID: uuid.New().String(),
		Topic:     PaymentResultTopic,
		Key:       fmt.Sprintf("%d", orderID),
		Value:     resultBytes,
		Sent:      false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	return repo.SaveOutboxMessage(ctx, tx, outboxMessage)
}

// Helper function to create a payment result with an existing transaction
func createPaymentResultWithTx(ctx context.Context, repo storage.Repository, tx *sqlx.Tx, orderID int64, status, reason string) error {
	result := models.PaymentResult{
		OrderID: orderID,
		Status:  status,
		Reason:  reason,
	}
	
	resultBytes, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal payment result: %w", err)
	}
	
	outboxMessage := &models.OutboxMessage{
		MessageID: uuid.New().String(),
		Topic:     PaymentResultTopic,
		Key:       fmt.Sprintf("%d", orderID),
		Value:     resultBytes,
		Sent:      false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	return repo.SaveOutboxMessage(ctx, tx, outboxMessage)
}