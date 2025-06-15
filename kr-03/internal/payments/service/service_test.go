package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"cbd/internal/payments/models"
	"cbd/internal/payments/storage"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// MockRepository is a mock implementation of the storage.Repository interface
type MockRepository struct {
	// Mock behavior flags and return values
	CreateAccountFunc            func(ctx context.Context, userID int64) (*models.Account, error)
	GetAccountByUserIDFunc       func(ctx context.Context, userID int64) (*models.Account, error)
	UpdateBalanceFunc            func(ctx context.Context, userID int64, amount float64) error
	SaveInboxMessageFunc         func(ctx context.Context, tx *sqlx.Tx, message *models.InboxMessage) error
	GetUnprocessedInboxMessagesFunc func(ctx context.Context, limit int) ([]*models.InboxMessage, error)
	MarkInboxMessageAsProcessedFunc func(ctx context.Context, tx *sqlx.Tx, messageID string) error
	SaveOutboxMessageFunc        func(ctx context.Context, tx *sqlx.Tx, message *models.OutboxMessage) error
	GetUnsendOutboxMessagesFunc  func(ctx context.Context, limit int) ([]*models.OutboxMessage, error)
	MarkOutboxMessageAsSentFunc  func(ctx context.Context, messageID string) error
	BeginTxFunc                  func(ctx context.Context) (*sqlx.Tx, error)
	CommitTxFunc                 func(tx *sqlx.Tx) error
	RollbackTxFunc               func(tx *sqlx.Tx) error
}

// Implement Repository interface methods
func (m *MockRepository) CreateAccount(ctx context.Context, userID int64) (*models.Account, error) {
	return m.CreateAccountFunc(ctx, userID)
}

func (m *MockRepository) GetAccountByUserID(ctx context.Context, userID int64) (*models.Account, error) {
	return m.GetAccountByUserIDFunc(ctx, userID)
}

func (m *MockRepository) UpdateBalance(ctx context.Context, userID int64, amount float64) error {
	return m.UpdateBalanceFunc(ctx, userID, amount)
}

func (m *MockRepository) SaveInboxMessage(ctx context.Context, tx *sqlx.Tx, message *models.InboxMessage) error {
	return m.SaveInboxMessageFunc(ctx, tx, message)
}

func (m *MockRepository) GetUnprocessedInboxMessages(ctx context.Context, limit int) ([]*models.InboxMessage, error) {
	return m.GetUnprocessedInboxMessagesFunc(ctx, limit)
}

func (m *MockRepository) MarkInboxMessageAsProcessed(ctx context.Context, tx *sqlx.Tx, messageID string) error {
	return m.MarkInboxMessageAsProcessedFunc(ctx, tx, messageID)
}

func (m *MockRepository) SaveOutboxMessage(ctx context.Context, tx *sqlx.Tx, message *models.OutboxMessage) error {
	return m.SaveOutboxMessageFunc(ctx, tx, message)
}

func (m *MockRepository) GetUnsendOutboxMessages(ctx context.Context, limit int) ([]*models.OutboxMessage, error) {
	return m.GetUnsendOutboxMessagesFunc(ctx, limit)
}

func (m *MockRepository) MarkOutboxMessageAsSent(ctx context.Context, messageID string) error {
	return m.MarkOutboxMessageAsSentFunc(ctx, messageID)
}

func (m *MockRepository) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return m.BeginTxFunc(ctx)
}

func (m *MockRepository) CommitTx(tx *sqlx.Tx) error {
	return m.CommitTxFunc(tx)
}

func (m *MockRepository) RollbackTx(tx *sqlx.Tx) error {
	return m.RollbackTxFunc(tx)
}

// Test NewService
func TestNewService(t *testing.T) {
	mockRepo := &MockRepository{}
	service := NewService(mockRepo)

	if service == nil {
		t.Fatal("Expected service to be created, got nil")
	}
}

// Test CreateAccount
func TestCreateAccount(t *testing.T) {
	// Setup mock repository
	expectedAccount := &models.Account{
		ID:        1,
		UserID:    123,
		Balance:   0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo := &MockRepository{
		CreateAccountFunc: func(ctx context.Context, userID int64) (*models.Account, error) {
			if userID == 123 {
				return expectedAccount, nil
			}
			return nil, errors.New("failed to create account")
		},
	}

	service := NewService(mockRepo)
	ctx := context.Background()

	// Test successful case
	account, err := service.CreateAccount(ctx, 123)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if account != expectedAccount {
		t.Errorf("Expected account to be %v, got %v", expectedAccount, account)
	}

	// Test error case
	account, err = service.CreateAccount(ctx, 456)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if account != nil {
		t.Errorf("Expected account to be nil, got %v", account)
	}
}

// Test Deposit success
func TestDeposit_Success(t *testing.T) {
	// Setup mock repository
	mockRepo := &MockRepository{
		GetAccountByUserIDFunc: func(ctx context.Context, userID int64) (*models.Account, error) {
			if userID == 123 {
				return &models.Account{
					ID:        1,
					UserID:    123,
					Balance:   100.0,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			}
			return nil, errors.New("account not found")
		},
		UpdateBalanceFunc: func(ctx context.Context, userID int64, amount float64) error {
			if userID == 123 {
				return nil
			}
			return errors.New("failed to update balance")
		},
	}

	service := NewService(mockRepo)
	ctx := context.Background()

	// Test successful case
	err := service.Deposit(ctx, 123, 50.0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Test account not found
	err = service.Deposit(ctx, 456, 50.0)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

// Test Deposit with negative amount
func TestDeposit_NegativeAmount(t *testing.T) {
	mockRepo := &MockRepository{}
	service := NewService(mockRepo)
	ctx := context.Background()

	// Test negative amount
	err := service.Deposit(ctx, 123, -50.0)
	if err == nil {
		t.Fatal("Expected error for negative amount, got nil")
	}
}

// Test GetBalance
func TestGetBalance(t *testing.T) {
	// Setup mock repository
	mockRepo := &MockRepository{
		GetAccountByUserIDFunc: func(ctx context.Context, userID int64) (*models.Account, error) {
			if userID == 123 {
				return &models.Account{
					ID:        1,
					UserID:    123,
					Balance:   150.0,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			}
			return nil, errors.New("account not found")
		},
	}

	service := NewService(mockRepo)
	ctx := context.Background()

	// Test successful case
	balance, err := service.GetBalance(ctx, 123)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if balance == nil {
		t.Fatal("Expected balance response, got nil")
	}

	if balance.Balance != 150.0 {
		t.Errorf("Expected balance to be 150.0, got %f", balance.Balance)
	}

	// Test account not found
	balance, err = service.GetBalance(ctx, 456)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if balance != nil {
		t.Errorf("Expected balance to be nil, got %v", balance)
	}
}

// Test ProcessPayment success
func TestProcessPayment_Success(t *testing.T) {
	// Setup mock repository
	mockTx := &sqlx.Tx{}
	mockRepo := &MockRepository{
		GetAccountByUserIDFunc: func(ctx context.Context, userID int64) (*models.Account, error) {
			if userID == 123 {
				return &models.Account{
					ID:        1,
					UserID:    123,
					Balance:   200.0,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			}
			return nil, errors.New("account not found")
		},
		BeginTxFunc: func(ctx context.Context) (*sqlx.Tx, error) {
			return mockTx, nil
		},
		UpdateBalanceFunc: func(ctx context.Context, userID int64, amount float64) error {
			if userID == 123 {
				return nil
			}
			return errors.New("failed to update balance")
		},
		SaveOutboxMessageFunc: func(ctx context.Context, tx *sqlx.Tx, message *models.OutboxMessage) error {
			return nil
		},
		CommitTxFunc: func(tx *sqlx.Tx) error {
			return nil
		},
		RollbackTxFunc: func(tx *sqlx.Tx) error {
			return nil
		},
	}

	service := NewService(mockRepo)
	ctx := context.Background()

	// Test successful case
	payment := &models.PaymentMessage{
		OrderID:     1,
		UserID:      123,
		Amount:      50.0,
		Description: "Test payment",
	}

	err := service.ProcessPayment(ctx, payment)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

// Test ProcessPayment with insufficient balance
func TestProcessPayment_InsufficientBalance(t *testing.T) {
	// Setup mock repository
	mockTx := &sqlx.Tx{}
	mockRepo := &MockRepository{
		GetAccountByUserIDFunc: func(ctx context.Context, userID int64) (*models.Account, error) {
			if userID == 123 {
				return &models.Account{
					ID:        1,
					UserID:    123,
					Balance:   30.0, // Less than payment amount
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			}
			return nil, errors.New("account not found")
		},
		BeginTxFunc: func(ctx context.Context) (*sqlx.Tx, error) {
			return mockTx, nil
		},
		SaveOutboxMessageFunc: func(ctx context.Context, tx *sqlx.Tx, message *models.OutboxMessage) error {
			return nil
		},
		CommitTxFunc: func(tx *sqlx.Tx) error {
			return nil
		},
		RollbackTxFunc: func(tx *sqlx.Tx) error {
			return nil
		},
	}

	service := NewService(mockRepo)
	ctx := context.Background()

	// Test insufficient balance
	payment := &models.PaymentMessage{
		OrderID:     1,
		UserID:      123,
		Amount:      50.0,
		Description: "Test payment",
	}

	err := service.ProcessPayment(ctx, payment)
	if err != nil {
		t.Fatalf("Expected no error (payment should be rejected but not error out), got %v", err)
	}
}

// Test ProcessInboxMessage
func TestProcessInboxMessage(t *testing.T) {
	// This test requires a real database connection
	// In Docker environment, we can use the environment variable to connect to the database
	dbURL := os.Getenv("DATABASE_URL_PAYMENTS")
	if dbURL == "" {
		// Use a default URL for Docker environment
		dbURL = "postgres://postgres:postgres@postgres:5432/payments?sslmode=disable"
	}

	// Create a real repository
	repo, err := storage.NewPostgresRepository(dbURL)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	// Create a service with the real repository
	service := NewService(repo)
	ctx := context.Background()

	// Create a test user account
	userID := time.Now().Unix() // Use current timestamp as a unique user ID
	account, err := service.CreateAccount(ctx, userID)
	if err != nil {
		t.Fatalf("Failed to create account: %v", err)
	}

	// Deposit some funds
	err = service.Deposit(ctx, userID, 100.0)
	if err != nil {
		t.Fatalf("Failed to deposit funds: %v", err)
	}

	// Create a payment message
	payment := &models.PaymentMessage{
		OrderID:     time.Now().Unix(), // Use current timestamp as a unique order ID
		UserID:      userID,
		Amount:      50.0,
		Description: "Test payment",
	}

	// Marshal the payment message
	paymentBytes, err := json.Marshal(payment)
	if err != nil {
		t.Fatalf("Failed to marshal payment message: %v", err)
	}

	// Create an inbox message
	inboxMessage := &models.InboxMessage{
		MessageID: uuid.New().String(),
		Topic:     PaymentRequestTopic,
		Key:       fmt.Sprintf("%d", payment.OrderID),
		Value:     paymentBytes,
		Processed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save the inbox message
	tx, err := repo.BeginTx(ctx)
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	err = repo.SaveInboxMessage(ctx, tx, inboxMessage)
	if err != nil {
		_ = repo.RollbackTx(tx)
		t.Fatalf("Failed to save inbox message: %v", err)
	}

	err = repo.CommitTx(tx)
	if err != nil {
		t.Fatalf("Failed to commit transaction: %v", err)
	}

	// Process the inbox message
	err = service.ProcessInboxMessage(ctx, inboxMessage)
	if err != nil {
		t.Fatalf("Failed to process inbox message: %v", err)
	}

	// Verify that the account balance was updated
	updatedAccount, err := repo.GetAccountByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("Failed to get account: %v", err)
	}

	expectedBalance := account.Balance + 100.0 - payment.Amount // Initial + deposit - payment
	if updatedAccount.Balance != expectedBalance {
		t.Errorf("Expected balance to be %f, got %f", expectedBalance, updatedAccount.Balance)
	}

	// Verify that the inbox message was marked as processed
	messages, err := repo.GetUnprocessedInboxMessages(ctx, 10)
	if err != nil {
		t.Fatalf("Failed to get unprocessed inbox messages: %v", err)
	}

	for _, msg := range messages {
		if msg.MessageID == inboxMessage.MessageID {
			t.Errorf("Expected inbox message to be marked as processed, but it's still unprocessed")
		}
	}
}
