package messaging

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"cbd/internal/orders/models"

	"github.com/jmoiron/sqlx"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of the storage.Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(*sqlx.Tx), args.Error(1)
}

func (m *MockRepository) CommitTx(tx *sqlx.Tx) error {
	args := m.Called(tx)
	return args.Error(0)
}

func (m *MockRepository) RollbackTx(tx *sqlx.Tx) error {
	args := m.Called(tx)
	return args.Error(0)
}

func (m *MockRepository) CreateOrder(ctx context.Context, tx *sqlx.Tx, order *models.Order) error {
	args := m.Called(ctx, tx, order)
	return args.Error(0)
}

func (m *MockRepository) GetOrderByID(ctx context.Context, id int64) (*models.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockRepository) GetOrdersByUserID(ctx context.Context, userID int64) ([]*models.Order, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*models.Order), args.Error(1)
}

func (m *MockRepository) UpdateOrderStatus(ctx context.Context, id int64, status models.OrderStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockRepository) SaveOutboxMessage(ctx context.Context, tx *sqlx.Tx, message *models.OutboxMessage) error {
	args := m.Called(ctx, tx, message)
	return args.Error(0)
}

func (m *MockRepository) GetUnsendOutboxMessages(ctx context.Context, limit int) ([]*models.OutboxMessage, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]*models.OutboxMessage), args.Error(1)
}

func (m *MockRepository) MarkOutboxMessageAsSent(ctx context.Context, messageID string) error {
	args := m.Called(ctx, messageID)
	return args.Error(0)
}

// MockOrderService is a mock implementation of the service.OrderService interface
type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) CreateOrder(ctx context.Context, userID int64, amount float64, description string) (*models.Order, error) {
	args := m.Called(ctx, userID, amount, description)
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderService) GetOrderByID(ctx context.Context, id int64) (*models.Order, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderService) GetOrdersByUserID(ctx context.Context, userID int64) ([]*models.Order, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*models.Order), args.Error(1)
}

func (m *MockOrderService) ProcessPaymentResult(ctx context.Context, result *models.PaymentResult) error {
	args := m.Called(ctx, result)
	return args.Error(0)
}

// MockWebSocketNotifier is a mock implementation of the WebSocketNotifier interface
type MockWebSocketNotifier struct {
	mock.Mock
}

func (m *MockWebSocketNotifier) NotifyOrderStatusChange(ctx context.Context, order *models.Order) {
	m.Called(ctx, order)
}

// MockKafkaWriter is a mock implementation of the kafka.Writer
type MockKafkaWriter struct {
	mock.Mock
}

func (m *MockKafkaWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	args := m.Called(ctx, msgs)
	return args.Error(0)
}

func (m *MockKafkaWriter) Close() error {
	args := m.Called()
	return args.Error(0)
}

// MockKafkaReader is a mock implementation of the kafka.Reader
type MockKafkaReader struct {
	mock.Mock
}

func (m *MockKafkaReader) ReadMessage(ctx context.Context) (kafka.Message, error) {
	args := m.Called(ctx)
	return args.Get(0).(kafka.Message), args.Error(1)
}

func (m *MockKafkaReader) Close() error {
	args := m.Called()
	return args.Error(0)
}

// TestNewProducer tests the NewProducer function
func TestNewProducer(t *testing.T) {
	// Setup
	mockRepo := new(MockRepository)
	mockOrderSvc := new(MockOrderService)
	mockWsHandler := new(MockWebSocketNotifier)

	// Test
	producer := NewProducer([]string{"localhost:9092"}, "test-group", mockRepo, mockOrderSvc, mockWsHandler)

	// Verify
	assert.NotNil(t, producer)
	assert.NotNil(t, producer.writer)
	assert.NotNil(t, producer.reader)
	assert.Equal(t, mockRepo, producer.repo)
	assert.Equal(t, mockOrderSvc, producer.orderSvc)
	assert.Equal(t, mockWsHandler, producer.wsHandler)
	assert.NotNil(t, producer.stopCh)
	assert.NotNil(t, producer.readerStopCh)
}

// TestNewOutboxPoller tests the NewOutboxPoller function
func TestNewOutboxPoller(t *testing.T) {
	// Setup
	mockRepo := new(MockRepository)
	mockWriter := &kafka.Writer{}

	// Test
	poller := NewOutboxPoller(mockRepo, mockWriter)

	// Verify
	assert.NotNil(t, poller)
	assert.Equal(t, mockRepo, poller.repo)
	assert.Equal(t, mockWriter, poller.writer)
}

// TestProcessPaymentResult tests the processPaymentResult function
func TestProcessPaymentResult(t *testing.T) {
	// Setup
	mockRepo := new(MockRepository)
	mockOrderSvc := new(MockOrderService)
	mockWsHandler := new(MockWebSocketNotifier)

	producer := &Producer{
		repo:      mockRepo,
		orderSvc:  mockOrderSvc,
		wsHandler: mockWsHandler,
	}

	// Create a payment result
	result := &models.PaymentResult{
		OrderID: 123,
		Status:  "success",
	}

	// Marshal the result to JSON
	resultBytes, _ := json.Marshal(result)

	// Create a Kafka message
	message := kafka.Message{
		Value: resultBytes,
	}

	// Create an order
	order := &models.Order{
		ID:     123,
		Status: models.OrderStatusFinished,
	}

	// Setup expectations
	mockOrderSvc.On("ProcessPaymentResult", mock.Anything, result).Return(nil)
	mockRepo.On("GetOrderByID", mock.Anything, int64(123)).Return(order, nil)
	mockWsHandler.On("NotifyOrderStatusChange", mock.Anything, order).Return()

	// Test
	err := producer.processPaymentResult(context.Background(), message)

	// Verify
	assert.NoError(t, err)
	mockOrderSvc.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockWsHandler.AssertExpectations(t)
}

// TestPollOutbox tests the pollOutbox function with a real Kafka writer
func TestPollOutbox(t *testing.T) {
	// Skip this test in Docker environment as it requires a real Kafka broker with the topic created
	if os.Getenv("KAFKA_BROKERS") != "" {
		t.Skip("Skipping TestPollOutbox in Docker environment as it requires a real Kafka broker with the topic created")
	}

	// Use localhost for local testing
	kafkaHost := "localhost:9092"
	if host := getEnv("KAFKA_HOST", ""); host != "" {
		kafkaHost = host
	}

	// Setup a real Kafka writer
	writer := &kafka.Writer{
		Addr:     kafka.TCP(kafkaHost),
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	// Setup a mock repository
	mockRepo := new(MockRepository)

	// Create an outbox poller with the real writer and mock repo
	poller := NewOutboxPoller(mockRepo, writer)

	// Create some test messages
	messages := []*models.OutboxMessage{
		{
			MessageID: "msg1",
			Topic:     "test-topic",
			Key:       "key1",
			Value:     []byte("value1"),
			Sent:      false,
		},
		{
			MessageID: "msg2",
			Topic:     "test-topic",
			Key:       "key2",
			Value:     []byte("value2"),
			Sent:      false,
		},
	}

	// Setup expectations
	mockRepo.On("GetUnsendOutboxMessages", mock.Anything, 10).Return(messages, nil)
	mockRepo.On("MarkOutboxMessageAsSent", mock.Anything, "msg1").Return(nil)
	mockRepo.On("MarkOutboxMessageAsSent", mock.Anything, "msg2").Return(nil)

	// Test
	err := poller.pollOutbox(context.Background())

	// Verify
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Helper function to get environment variables with default values
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
