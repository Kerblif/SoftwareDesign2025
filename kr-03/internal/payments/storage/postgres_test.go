package storage

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"cbd/internal/payments/models"

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
	// This test requires a real database connection
	// In Docker environment, we can use the environment variable to connect to the database
	dbURL := os.Getenv("DATABASE_URL_PAYMENTS")
	if dbURL == "" {
		// Use a default URL for Docker environment
		dbURL = "postgres://postgres:postgres@postgres:5432/payments?sslmode=disable"
	}

	// Create a new repository
	repo, err := NewPostgresRepository(dbURL)
	assert.NoError(t, err)
	assert.NotNil(t, repo)
	assert.NotNil(t, repo.db)

	// Close the connection
	err = repo.db.Close()
	assert.NoError(t, err)
}

func TestCreateAccount(t *testing.T) {
	// Setup
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	ctx := context.Background()
	userID := int64(1)

	// Test account doesn't exist
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM accounts WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectQuery("INSERT INTO accounts \\(user_id, balance, created_at, updated_at\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) RETURNING id").
		WithArgs(userID, 0.0, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(123))

	account, err := repo.CreateAccount(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, int64(123), account.ID)
	assert.Equal(t, userID, account.UserID)
	assert.Equal(t, 0.0, account.Balance)

	// Test account already exists
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM accounts WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	account, err = repo.CreateAccount(ctx, userID)
	assert.Error(t, err)
	assert.Nil(t, account)
	assert.Contains(t, err.Error(), "account already exists")

	// Test error checking if account exists
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM accounts WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnError(assert.AnError)

	account, err = repo.CreateAccount(ctx, userID)
	assert.Error(t, err)
	assert.Nil(t, account)

	// Test error creating account
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM accounts WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectQuery("INSERT INTO accounts \\(user_id, balance, created_at, updated_at\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) RETURNING id").
		WithArgs(userID, 0.0, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(assert.AnError)

	account, err = repo.CreateAccount(ctx, userID)
	assert.Error(t, err)
	assert.Nil(t, account)
}

func TestGetAccountByUserID(t *testing.T) {
	// Setup
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	ctx := context.Background()
	userID := int64(1)
	now := time.Now()
	expectedAccount := &models.Account{
		ID:        123,
		UserID:    userID,
		Balance:   100.0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Test success case
	rows := sqlmock.NewRows([]string{"id", "user_id", "balance", "created_at", "updated_at"}).
		AddRow(expectedAccount.ID, expectedAccount.UserID, expectedAccount.Balance, expectedAccount.CreatedAt, expectedAccount.UpdatedAt)

	mock.ExpectQuery("SELECT \\* FROM accounts WHERE user_id = \\$1").
		WithArgs(userID).
		WillReturnRows(rows)

	account, err := repo.GetAccountByUserID(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedAccount.ID, account.ID)
	assert.Equal(t, expectedAccount.UserID, account.UserID)
	assert.Equal(t, expectedAccount.Balance, account.Balance)

	// Test account not found
	mock.ExpectQuery("SELECT \\* FROM accounts WHERE user_id = \\$1").
		WithArgs(int64(2)).
		WillReturnError(sql.ErrNoRows)

	account, err = repo.GetAccountByUserID(ctx, int64(2))
	assert.Error(t, err)
	assert.Nil(t, account)
	assert.Contains(t, err.Error(), "account not found")

	// Test other error
	mock.ExpectQuery("SELECT \\* FROM accounts WHERE user_id = \\$1").
		WithArgs(int64(3)).
		WillReturnError(assert.AnError)

	account, err = repo.GetAccountByUserID(ctx, int64(3))
	assert.Error(t, err)
	assert.Nil(t, account)
}

func TestSaveInboxMessage(t *testing.T) {
	// Setup
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	ctx := context.Background()
	now := time.Now()
	message := &models.InboxMessage{
		MessageID: "msg1",
		Topic:     "topic1",
		Key:       "key1",
		Value:     []byte("value1"),
		Processed: false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Test with transaction
	mock.ExpectExec("INSERT INTO inbox").
		WithArgs(message.MessageID, message.Topic, message.Key, message.Value, message.Processed, message.CreatedAt, message.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	tx, _ := sqlxDB.Beginx()
	err := repo.SaveInboxMessage(ctx, tx, message)
	assert.NoError(t, err)

	// Test without transaction
	mock.ExpectExec("INSERT INTO inbox").
		WithArgs(message.MessageID, message.Topic, message.Key, message.Value, message.Processed, message.CreatedAt, message.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(2, 1))

	err = repo.SaveInboxMessage(ctx, nil, message)
	assert.NoError(t, err)

	// Test error case
	mock.ExpectExec("INSERT INTO inbox").
		WithArgs(message.MessageID, message.Topic, message.Key, message.Value, message.Processed, message.CreatedAt, message.UpdatedAt).
		WillReturnError(assert.AnError)

	err = repo.SaveInboxMessage(ctx, nil, message)
	assert.Error(t, err)
}

func TestGetUnprocessedInboxMessages(t *testing.T) {
	// Setup
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	ctx := context.Background()
	now := time.Now()
	expectedMessages := []*models.InboxMessage{
		{
			ID:        1,
			MessageID: "msg1",
			Topic:     "topic1",
			Key:       "key1",
			Value:     []byte("value1"),
			Processed: false,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        2,
			MessageID: "msg2",
			Topic:     "topic2",
			Key:       "key2",
			Value:     []byte("value2"),
			Processed: false,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	// Test success case
	rows := sqlmock.NewRows([]string{"id", "message_id", "topic", "key", "value", "processed", "created_at", "updated_at"})
	for _, message := range expectedMessages {
		rows.AddRow(message.ID, message.MessageID, message.Topic, message.Key, message.Value, message.Processed, message.CreatedAt, message.UpdatedAt)
	}

	mock.ExpectQuery("SELECT \\* FROM inbox WHERE processed = false ORDER BY created_at ASC LIMIT \\$1").
		WithArgs(10).
		WillReturnRows(rows)

	messages, err := repo.GetUnprocessedInboxMessages(ctx, 10)
	assert.NoError(t, err)
	assert.Equal(t, len(expectedMessages), len(messages))
	for i, message := range messages {
		assert.Equal(t, expectedMessages[i].ID, message.ID)
		assert.Equal(t, expectedMessages[i].MessageID, message.MessageID)
		assert.Equal(t, expectedMessages[i].Topic, message.Topic)
		assert.Equal(t, expectedMessages[i].Key, message.Key)
		assert.Equal(t, expectedMessages[i].Processed, message.Processed)
	}

	// Test error case
	mock.ExpectQuery("SELECT \\* FROM inbox WHERE processed = false ORDER BY created_at ASC LIMIT \\$1").
		WithArgs(5).
		WillReturnError(assert.AnError)

	messages, err = repo.GetUnprocessedInboxMessages(ctx, 5)
	assert.Error(t, err)
	assert.Nil(t, messages)
}

func TestMarkInboxMessageAsProcessed(t *testing.T) {
	// Setup
	sqlxDB, mock, repo := setupMockDB(t)
	defer sqlxDB.Close()

	ctx := context.Background()
	messageID := "msg1"

	// Test with transaction
	mock.ExpectExec("UPDATE inbox SET processed = true, updated_at = \\$1 WHERE message_id = \\$2").
		WithArgs(sqlmock.AnyArg(), messageID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	tx, _ := sqlxDB.Beginx()
	err := repo.MarkInboxMessageAsProcessed(ctx, tx, messageID)
	assert.NoError(t, err)

	// Test without transaction
	mock.ExpectExec("UPDATE inbox SET processed = true, updated_at = \\$1 WHERE message_id = \\$2").
		WithArgs(sqlmock.AnyArg(), messageID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.MarkInboxMessageAsProcessed(ctx, nil, messageID)
	assert.NoError(t, err)

	// Test error case
	mock.ExpectExec("UPDATE inbox SET processed = true, updated_at = \\$1 WHERE message_id = \\$2").
		WithArgs(sqlmock.AnyArg(), messageID).
		WillReturnError(assert.AnError)

	err = repo.MarkInboxMessageAsProcessed(ctx, nil, messageID)
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
	messageID := "msg1"

	// Test success case
	mock.ExpectExec("UPDATE outbox SET sent = true, updated_at = \\$1 WHERE message_id = \\$2").
		WithArgs(sqlmock.AnyArg(), messageID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.MarkOutboxMessageAsSent(ctx, messageID)
	assert.NoError(t, err)

	// Test error case
	mock.ExpectExec("UPDATE outbox SET sent = true, updated_at = \\$1 WHERE message_id = \\$2").
		WithArgs(sqlmock.AnyArg(), messageID).
		WillReturnError(assert.AnError)

	err = repo.MarkOutboxMessageAsSent(ctx, messageID)
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
