package storage

import (
	"context"
	"testing"
	"time"

	"cbd/internal/orders/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, *PostgresRepository) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repo := &PostgresRepository{db: sqlxDB}

	return sqlxDB, mock, repo
}

func TestNewPostgresRepository(t *testing.T) {
	// Skip this test as it requires a real database connection
	t.Skip("Skipping TestNewPostgresRepository as it requires a real database connection")
}

func TestCreateOrder(t *testing.T) {
	// Setup
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	ctx := context.Background()
	now := time.Now()
	order := &models.Order{
		UserID:      1,
		Amount:      100.0,
		Description: "Test order",
		Status:      models.OrderStatusNew,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Test with transaction
	mock.ExpectQuery("INSERT INTO orders").
		WithArgs(order.UserID, order.Amount, order.Description, order.Status, order.CreatedAt, order.UpdatedAt).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(123))

	tx, _ := sqlxDB.Beginx()
	err := repo.CreateOrder(ctx, tx, order)
	assert.NoError(t, err)
	assert.Equal(t, int64(123), order.ID)

	// Test without transaction
	mock.ExpectQuery("INSERT INTO orders").
		WithArgs(order.UserID, order.Amount, order.Description, order.Status, order.CreatedAt, order.UpdatedAt).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(456))

	err = repo.CreateOrder(ctx, nil, order)
	assert.NoError(t, err)
	assert.Equal(t, int64(456), order.ID)

	// Test error case
	mock.ExpectQuery("INSERT INTO orders").
		WithArgs(order.UserID, order.Amount, order.Description, order.Status, order.CreatedAt, order.UpdatedAt).
		WillReturnError(assert.AnError)

	err = repo.CreateOrder(ctx, nil, order)
	assert.Error(t, err)
}

func TestGetOrderByID(t *testing.T) {
	// Setup
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	ctx := context.Background()
	now := time.Now()
	expectedOrder := &models.Order{
		ID:          123,
		UserID:      1,
		Amount:      100.0,
		Description: "Test order",
		Status:      models.OrderStatusNew,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Test success case
	rows := sqlmock.NewRows([]string{"id", "user_id", "amount", "description", "status", "created_at", "updated_at"}).
		AddRow(expectedOrder.ID, expectedOrder.UserID, expectedOrder.Amount, expectedOrder.Description, expectedOrder.Status, expectedOrder.CreatedAt, expectedOrder.UpdatedAt)

	mock.ExpectQuery("SELECT \\* FROM orders WHERE id = \\$1").
		WithArgs(123).
		WillReturnRows(rows)

	order, err := repo.GetOrderByID(ctx, 123)
	assert.NoError(t, err)
	assert.Equal(t, expectedOrder.ID, order.ID)
	assert.Equal(t, expectedOrder.UserID, order.UserID)
	assert.Equal(t, expectedOrder.Amount, order.Amount)
	assert.Equal(t, expectedOrder.Description, order.Description)
	assert.Equal(t, expectedOrder.Status, order.Status)

	// Test error case
	mock.ExpectQuery("SELECT \\* FROM orders WHERE id = \\$1").
		WithArgs(456).
		WillReturnError(assert.AnError)

	order, err = repo.GetOrderByID(ctx, 456)
	assert.Error(t, err)
	assert.Nil(t, order)
}

func TestGetOrdersByUserID(t *testing.T) {
	// Setup
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	ctx := context.Background()
	now := time.Now()
	expectedOrders := []*models.Order{
		{
			ID:          123,
			UserID:      1,
			Amount:      100.0,
			Description: "Test order 1",
			Status:      models.OrderStatusNew,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          124,
			UserID:      1,
			Amount:      200.0,
			Description: "Test order 2",
			Status:      models.OrderStatusFinished,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	// Test success case
	rows := sqlmock.NewRows([]string{"id", "user_id", "amount", "description", "status", "created_at", "updated_at"})
	for _, order := range expectedOrders {
		rows.AddRow(order.ID, order.UserID, order.Amount, order.Description, order.Status, order.CreatedAt, order.UpdatedAt)
	}

	mock.ExpectQuery("SELECT \\* FROM orders WHERE user_id = \\$1 ORDER BY created_at DESC").
		WithArgs(1).
		WillReturnRows(rows)

	orders, err := repo.GetOrdersByUserID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, len(expectedOrders), len(orders))
	for i, order := range orders {
		assert.Equal(t, expectedOrders[i].ID, order.ID)
		assert.Equal(t, expectedOrders[i].UserID, order.UserID)
		assert.Equal(t, expectedOrders[i].Amount, order.Amount)
		assert.Equal(t, expectedOrders[i].Description, order.Description)
		assert.Equal(t, expectedOrders[i].Status, order.Status)
	}

	// Test error case
	mock.ExpectQuery("SELECT \\* FROM orders WHERE user_id = \\$1 ORDER BY created_at DESC").
		WithArgs(2).
		WillReturnError(assert.AnError)

	orders, err = repo.GetOrdersByUserID(ctx, 2)
	assert.Error(t, err)
	assert.Nil(t, orders)
}

func TestUpdateOrderStatus(t *testing.T) {
	// Setup
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	ctx := context.Background()

	// Test success case
	mock.ExpectExec("UPDATE orders SET status = \\$1, updated_at = \\$2 WHERE id = \\$3").
		WithArgs(models.OrderStatusFinished, sqlmock.AnyArg(), 123).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateOrderStatus(ctx, 123, models.OrderStatusFinished)
	assert.NoError(t, err)

	// Test error case
	mock.ExpectExec("UPDATE orders SET status = \\$1, updated_at = \\$2 WHERE id = \\$3").
		WithArgs(models.OrderStatusCancelled, sqlmock.AnyArg(), 456).
		WillReturnError(assert.AnError)

	err = repo.UpdateOrderStatus(ctx, 456, models.OrderStatusCancelled)
	assert.Error(t, err)
}

func TestSaveOutboxMessage(t *testing.T) {
	// Setup
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	ctx := context.Background()
	now := time.Now()
	message := &models.OutboxMessage{
		MessageID: "msg1",
		Topic:     "topic1",
		Key:       "key1",
		Value:     []byte("value1"),
		Sent:      false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Test with transaction
	mock.ExpectExec("INSERT INTO outbox").
		WithArgs(message.MessageID, message.Topic, message.Key, message.Value, message.Sent, message.CreatedAt, message.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	tx, _ := sqlxDB.Beginx()
	err := repo.SaveOutboxMessage(ctx, tx, message)
	assert.NoError(t, err)

	// Test without transaction
	mock.ExpectExec("INSERT INTO outbox").
		WithArgs(message.MessageID, message.Topic, message.Key, message.Value, message.Sent, message.CreatedAt, message.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(2, 1))

	err = repo.SaveOutboxMessage(ctx, nil, message)
	assert.NoError(t, err)

	// Test error case
	mock.ExpectExec("INSERT INTO outbox").
		WithArgs(message.MessageID, message.Topic, message.Key, message.Value, message.Sent, message.CreatedAt, message.UpdatedAt).
		WillReturnError(assert.AnError)

	err = repo.SaveOutboxMessage(ctx, nil, message)
	assert.Error(t, err)
}

func TestGetUnsendOutboxMessages(t *testing.T) {
	// Setup
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	ctx := context.Background()
	now := time.Now()
	expectedMessages := []*models.OutboxMessage{
		{
			ID:        1,
			MessageID: "msg1",
			Topic:     "topic1",
			Key:       "key1",
			Value:     []byte("value1"),
			Sent:      false,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        2,
			MessageID: "msg2",
			Topic:     "topic2",
			Key:       "key2",
			Value:     []byte("value2"),
			Sent:      false,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	// Test success case
	rows := sqlmock.NewRows([]string{"id", "message_id", "topic", "key", "value", "sent", "created_at", "updated_at"})
	for _, message := range expectedMessages {
		rows.AddRow(message.ID, message.MessageID, message.Topic, message.Key, message.Value, message.Sent, message.CreatedAt, message.UpdatedAt)
	}

	mock.ExpectQuery("SELECT \\* FROM outbox WHERE sent = false ORDER BY created_at ASC LIMIT \\$1").
		WithArgs(10).
		WillReturnRows(rows)

	messages, err := repo.GetUnsendOutboxMessages(ctx, 10)
	assert.NoError(t, err)
	assert.Equal(t, len(expectedMessages), len(messages))
	for i, message := range messages {
		assert.Equal(t, expectedMessages[i].ID, message.ID)
		assert.Equal(t, expectedMessages[i].MessageID, message.MessageID)
		assert.Equal(t, expectedMessages[i].Topic, message.Topic)
		assert.Equal(t, expectedMessages[i].Key, message.Key)
		assert.Equal(t, expectedMessages[i].Value, message.Value)
		assert.Equal(t, expectedMessages[i].Sent, message.Sent)
	}

	// Test error case
	mock.ExpectQuery("SELECT \\* FROM outbox WHERE sent = false ORDER BY created_at ASC LIMIT \\$1").
		WithArgs(5).
		WillReturnError(assert.AnError)

	messages, err = repo.GetUnsendOutboxMessages(ctx, 5)
	assert.Error(t, err)
	assert.Nil(t, messages)
}

func TestMarkOutboxMessageAsSent(t *testing.T) {
	// Setup
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	ctx := context.Background()

	// Test success case
	mock.ExpectExec("UPDATE outbox SET sent = true, updated_at = \\$1 WHERE message_id = \\$2").
		WithArgs(sqlmock.AnyArg(), "msg1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.MarkOutboxMessageAsSent(ctx, "msg1")
	assert.NoError(t, err)

	// Test error case
	mock.ExpectExec("UPDATE outbox SET sent = true, updated_at = \\$1 WHERE message_id = \\$2").
		WithArgs(sqlmock.AnyArg(), "msg2").
		WillReturnError(assert.AnError)

	err = repo.MarkOutboxMessageAsSent(ctx, "msg2")
	assert.Error(t, err)
}

func TestBeginTx(t *testing.T) {
	// Setup
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	ctx := context.Background()

	// Test success case
	mock.ExpectBegin()

	tx, err := repo.BeginTx(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, tx)

	// Test error case
	mock.ExpectBegin().WillReturnError(assert.AnError)

	tx, err = repo.BeginTx(ctx)
	assert.Error(t, err)
	assert.Nil(t, tx)
}

func TestCommitTx(t *testing.T) {
	// Setup
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	// Test success case
	mock.ExpectBegin()
	mock.ExpectCommit()

	tx, _ := sqlxDB.Beginx()
	err := repo.CommitTx(tx)
	assert.NoError(t, err)

	// Test error case
	mock.ExpectBegin()
	mock.ExpectCommit().WillReturnError(assert.AnError)

	tx, _ = sqlxDB.Beginx()
	err = repo.CommitTx(tx)
	assert.Error(t, err)
}

func TestRollbackTx(t *testing.T) {
	// Setup
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	// Test success case
	mock.ExpectBegin()
	mock.ExpectRollback()

	tx, _ := sqlxDB.Beginx()
	err := repo.RollbackTx(tx)
	assert.NoError(t, err)

	// Test error case
	mock.ExpectBegin()
	mock.ExpectRollback().WillReturnError(assert.AnError)

	tx, _ = sqlxDB.Beginx()
	err = repo.RollbackTx(tx)
	assert.Error(t, err)
}
