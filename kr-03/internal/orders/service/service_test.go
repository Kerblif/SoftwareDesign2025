package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"cbd/internal/orders/models"

	"github.com/jmoiron/sqlx"
)

// MockRepository is a mock implementation of the storage.Repository interface
type MockRepository struct {
	// Mock behavior flags and return values
	BeginTxFunc                func(ctx context.Context) (*sqlx.Tx, error)
	CommitTxFunc               func(tx *sqlx.Tx) error
	RollbackTxFunc             func(tx *sqlx.Tx) error
	CreateOrderFunc            func(ctx context.Context, tx *sqlx.Tx, order *models.Order) error
	GetOrderByIDFunc           func(ctx context.Context, id int64) (*models.Order, error)
	GetOrdersByUserIDFunc      func(ctx context.Context, userID int64) ([]*models.Order, error)
	UpdateOrderStatusFunc      func(ctx context.Context, id int64, status models.OrderStatus) error
	SaveOutboxMessageFunc      func(ctx context.Context, tx *sqlx.Tx, message *models.OutboxMessage) error
	GetUnsendOutboxMessagesFunc func(ctx context.Context, limit int) ([]*models.OutboxMessage, error)
	MarkOutboxMessageAsSentFunc func(ctx context.Context, messageID string) error
}

// Implement Repository interface methods
func (m *MockRepository) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return m.BeginTxFunc(ctx)
}

func (m *MockRepository) CommitTx(tx *sqlx.Tx) error {
	return m.CommitTxFunc(tx)
}

func (m *MockRepository) RollbackTx(tx *sqlx.Tx) error {
	return m.RollbackTxFunc(tx)
}

func (m *MockRepository) CreateOrder(ctx context.Context, tx *sqlx.Tx, order *models.Order) error {
	return m.CreateOrderFunc(ctx, tx, order)
}

func (m *MockRepository) GetOrderByID(ctx context.Context, id int64) (*models.Order, error) {
	return m.GetOrderByIDFunc(ctx, id)
}

func (m *MockRepository) GetOrdersByUserID(ctx context.Context, userID int64) ([]*models.Order, error) {
	return m.GetOrdersByUserIDFunc(ctx, userID)
}

func (m *MockRepository) UpdateOrderStatus(ctx context.Context, id int64, status models.OrderStatus) error {
	return m.UpdateOrderStatusFunc(ctx, id, status)
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

// Test NewService
func TestNewService(t *testing.T) {
	mockRepo := &MockRepository{}
	service := NewService(mockRepo)

	if service == nil {
		t.Fatal("Expected service to be created, got nil")
	}

	if service.repo != mockRepo {
		t.Errorf("Expected repo to be %v, got %v", mockRepo, service.repo)
	}
}

// Test CreateOrder success case
func TestCreateOrder_Success(t *testing.T) {
	// Setup mock repository
	mockTx := &sqlx.Tx{}
	mockRepo := &MockRepository{
		BeginTxFunc: func(ctx context.Context) (*sqlx.Tx, error) {
			return mockTx, nil
		},
		CreateOrderFunc: func(ctx context.Context, tx *sqlx.Tx, order *models.Order) error {
			// Set ID to simulate database insert
			order.ID = 123
			return nil
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

	// Call the method
	order, err := service.CreateOrder(ctx, 1, 100.0, "Test order")

	// Verify results
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if order == nil {
		t.Fatal("Expected order to be created, got nil")
	}

	if order.ID != 123 {
		t.Errorf("Expected order ID to be 123, got %d", order.ID)
	}

	if order.UserID != 1 {
		t.Errorf("Expected user ID to be 1, got %d", order.UserID)
	}

	if order.Amount != 100.0 {
		t.Errorf("Expected amount to be 100.0, got %f", order.Amount)
	}

	if order.Description != "Test order" {
		t.Errorf("Expected description to be 'Test order', got '%s'", order.Description)
	}

	if order.Status != models.OrderStatusNew {
		t.Errorf("Expected status to be %s, got %s", models.OrderStatusNew, order.Status)
	}
}

// Test CreateOrder with BeginTx error
func TestCreateOrder_BeginTxError(t *testing.T) {
	// Setup mock repository
	mockRepo := &MockRepository{
		BeginTxFunc: func(ctx context.Context) (*sqlx.Tx, error) {
			return nil, errors.New("begin tx error")
		},
	}

	service := NewService(mockRepo)
	ctx := context.Background()

	// Call the method
	order, err := service.CreateOrder(ctx, 1, 100.0, "Test order")

	// Verify results
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if order != nil {
		t.Errorf("Expected order to be nil, got %v", order)
	}
}

// Test CreateOrder with CreateOrder error
func TestCreateOrder_CreateOrderError(t *testing.T) {
	// Setup mock repository
	mockTx := &sqlx.Tx{}
	mockRepo := &MockRepository{
		BeginTxFunc: func(ctx context.Context) (*sqlx.Tx, error) {
			return mockTx, nil
		},
		CreateOrderFunc: func(ctx context.Context, tx *sqlx.Tx, order *models.Order) error {
			return errors.New("create order error")
		},
		RollbackTxFunc: func(tx *sqlx.Tx) error {
			return nil
		},
	}

	service := NewService(mockRepo)
	ctx := context.Background()

	// Call the method
	order, err := service.CreateOrder(ctx, 1, 100.0, "Test order")

	// Verify results
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if order != nil {
		t.Errorf("Expected order to be nil, got %v", order)
	}
}

// Test GetOrderByID
func TestGetOrderByID(t *testing.T) {
	// Setup mock repository
	expectedOrder := &models.Order{
		ID:          123,
		UserID:      1,
		Amount:      100.0,
		Description: "Test order",
		Status:      models.OrderStatusNew,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRepo := &MockRepository{
		GetOrderByIDFunc: func(ctx context.Context, id int64) (*models.Order, error) {
			if id == 123 {
				return expectedOrder, nil
			}
			return nil, errors.New("order not found")
		},
	}

	service := NewService(mockRepo)
	ctx := context.Background()

	// Test successful case
	order, err := service.GetOrderByID(ctx, 123)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if order != expectedOrder {
		t.Errorf("Expected order to be %v, got %v", expectedOrder, order)
	}

	// Test error case
	order, err = service.GetOrderByID(ctx, 456)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if order != nil {
		t.Errorf("Expected order to be nil, got %v", order)
	}
}

// Test GetOrdersByUserID
func TestGetOrdersByUserID(t *testing.T) {
	// Setup mock repository
	expectedOrders := []*models.Order{
		{
			ID:          123,
			UserID:      1,
			Amount:      100.0,
			Description: "Test order 1",
			Status:      models.OrderStatusNew,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          124,
			UserID:      1,
			Amount:      200.0,
			Description: "Test order 2",
			Status:      models.OrderStatusFinished,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	mockRepo := &MockRepository{
		GetOrdersByUserIDFunc: func(ctx context.Context, userID int64) ([]*models.Order, error) {
			if userID == 1 {
				return expectedOrders, nil
			}
			return nil, errors.New("user not found")
		},
	}

	service := NewService(mockRepo)
	ctx := context.Background()

	// Test successful case
	orders, err := service.GetOrdersByUserID(ctx, 1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(orders) != len(expectedOrders) {
		t.Errorf("Expected %d orders, got %d", len(expectedOrders), len(orders))
	}

	// Test error case
	orders, err = service.GetOrdersByUserID(ctx, 2)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if orders != nil {
		t.Errorf("Expected orders to be nil, got %v", orders)
	}
}

// Test ProcessPaymentResult
func TestProcessPaymentResult(t *testing.T) {
	// Setup mock repository
	mockRepo := &MockRepository{
		UpdateOrderStatusFunc: func(ctx context.Context, id int64, status models.OrderStatus) error {
			if id == 123 {
				return nil
			}
			return errors.New("order not found")
		},
		GetOrderByIDFunc: func(ctx context.Context, id int64) (*models.Order, error) {
			if id == 123 {
				return &models.Order{
					ID:     123,
					Status: models.OrderStatusFinished,
				}, nil
			}
			return nil, errors.New("order not found")
		},
	}

	service := NewService(mockRepo)
	ctx := context.Background()

	// Test successful case with success status
	result := &models.PaymentResult{
		OrderID: 123,
		Status:  "success",
	}

	err := service.ProcessPaymentResult(ctx, result)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Test successful case with failure status
	result = &models.PaymentResult{
		OrderID: 123,
		Status:  "failure",
	}

	err = service.ProcessPaymentResult(ctx, result)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Test error case
	result = &models.PaymentResult{
		OrderID: 456,
		Status:  "success",
	}

	err = service.ProcessPaymentResult(ctx, result)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}
