package storage

import (
	"context"
	"os"
	"testing"
	"time"

	"cbd/internal/payments/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPostgresIntegration tests the integration with a real Postgres database
// This test will be skipped if the DATABASE_URL_PAYMENTS environment variable is not set
func TestPostgresIntegration(t *testing.T) {
	// Get the database URL from the environment
	dbURL := os.Getenv("DATABASE_URL_PAYMENTS")
	if dbURL == "" {
		t.Skip("Skipping integration test as DATABASE_URL_PAYMENTS is not set")
	}

	// Create a new repository
	repo, err := NewPostgresRepository(dbURL)
	require.NoError(t, err)
	require.NotNil(t, repo)
	defer repo.db.Close()

	// Clean up the database before and after the test
	cleanupDB(t, repo)
	defer cleanupDB(t, repo)

	// Test creating an account
	ctx := context.Background()
	userID := int64(1000) // Use a high ID to avoid conflicts with existing data

	account, err := repo.CreateAccount(ctx, userID)
	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, userID, account.UserID)
	assert.Equal(t, 0.0, account.Balance)

	// Test getting the account
	account, err = repo.GetAccountByUserID(ctx, userID)
	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, userID, account.UserID)
	assert.Equal(t, 0.0, account.Balance)

	// Test updating the balance
	err = repo.UpdateBalance(ctx, userID, 100.0)
	assert.NoError(t, err)

	account, err = repo.GetAccountByUserID(ctx, userID)
	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, userID, account.UserID)
	assert.Equal(t, 100.0, account.Balance)

	// Test inbox message operations
	inboxMessage := &models.InboxMessage{
		MessageID: "test-inbox-msg",
		Topic:     "test-topic",
		Key:       "test-key",
		Value:     []byte("test-value"),
		Processed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test saving an inbox message
	err = repo.SaveInboxMessage(ctx, nil, inboxMessage)
	assert.NoError(t, err)

	// Test getting unprocessed inbox messages
	messages, err := repo.GetUnprocessedInboxMessages(ctx, 10)
	assert.NoError(t, err)
	assert.NotEmpty(t, messages)
	found := false
	for _, msg := range messages {
		if msg.MessageID == inboxMessage.MessageID {
			found = true
			break
		}
	}
	assert.True(t, found, "Inbox message not found")

	// Test marking an inbox message as processed
	err = repo.MarkInboxMessageAsProcessed(ctx, nil, inboxMessage.MessageID)
	assert.NoError(t, err)

	// Test outbox message operations
	outboxMessage := &models.OutboxMessage{
		MessageID: "test-outbox-msg",
		Topic:     "test-topic",
		Key:       "test-key",
		Value:     []byte("test-value"),
		Sent:      false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test saving an outbox message
	err = repo.SaveOutboxMessage(ctx, nil, outboxMessage)
	assert.NoError(t, err)

	// Test getting unsent outbox messages
	outboxMessages, err := repo.GetUnsendOutboxMessages(ctx, 10)
	assert.NoError(t, err)
	assert.NotEmpty(t, outboxMessages)
	found = false
	for _, msg := range outboxMessages {
		if msg.MessageID == outboxMessage.MessageID {
			found = true
			break
		}
	}
	assert.True(t, found, "Outbox message not found")

	// Test marking an outbox message as sent
	err = repo.MarkOutboxMessageAsSent(ctx, outboxMessage.MessageID)
	assert.NoError(t, err)

	// Test transaction operations
	tx, err := repo.BeginTx(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, tx)

	err = repo.CommitTx(tx)
	assert.NoError(t, err)

	tx, err = repo.BeginTx(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, tx)

	err = repo.RollbackTx(tx)
	assert.NoError(t, err)
}

// cleanupDB cleans up the database by deleting test data
func cleanupDB(t *testing.T, repo *PostgresRepository) {
	ctx := context.Background()

	// Delete test account
	_, err := repo.db.ExecContext(ctx, "DELETE FROM accounts WHERE user_id = $1", int64(1000))
	assert.NoError(t, err)

	// Delete test inbox messages
	_, err = repo.db.ExecContext(ctx, "DELETE FROM inbox WHERE message_id = $1", "test-inbox-msg")
	assert.NoError(t, err)

	// Delete test outbox messages
	_, err = repo.db.ExecContext(ctx, "DELETE FROM outbox WHERE message_id = $1", "test-outbox-msg")
	assert.NoError(t, err)
}
