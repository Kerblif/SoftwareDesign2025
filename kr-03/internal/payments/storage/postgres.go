package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"cbd/internal/payments/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Repository defines the interface for the payments storage
type Repository interface {
	// Account operations
	CreateAccount(ctx context.Context, userID int64) (*models.Account, error)
	GetAccountByUserID(ctx context.Context, userID int64) (*models.Account, error)
	UpdateBalance(ctx context.Context, userID int64, amount float64) error
	
	// Inbox operations
	SaveInboxMessage(ctx context.Context, tx *sqlx.Tx, message *models.InboxMessage) error
	GetUnprocessedInboxMessages(ctx context.Context, limit int) ([]*models.InboxMessage, error)
	MarkInboxMessageAsProcessed(ctx context.Context, tx *sqlx.Tx, messageID string) error
	
	// Outbox operations
	SaveOutboxMessage(ctx context.Context, tx *sqlx.Tx, message *models.OutboxMessage) error
	GetUnsendOutboxMessages(ctx context.Context, limit int) ([]*models.OutboxMessage, error)
	MarkOutboxMessageAsSent(ctx context.Context, messageID string) error
	
	// Transaction management
	BeginTx(ctx context.Context) (*sqlx.Tx, error)
	CommitTx(tx *sqlx.Tx) error
	RollbackTx(tx *sqlx.Tx) error
}

// PostgresRepository implements Repository interface
type PostgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository creates a new PostgresRepository
func NewPostgresRepository(connectionString string) (*PostgresRepository, error) {
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create tables if they don't exist
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &PostgresRepository{db: db}, nil
}

// createTables creates the necessary tables if they don't exist
func createTables(db *sqlx.DB) error {
	// Create accounts table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS accounts (
			id SERIAL PRIMARY KEY,
			user_id BIGINT UNIQUE NOT NULL,
			balance DECIMAL(10, 2) NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		);
	`)
	if err != nil {
		return err
	}

	// Create inbox table for transactional inbox pattern
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS inbox (
			id SERIAL PRIMARY KEY,
			message_id VARCHAR(255) UNIQUE NOT NULL,
			topic VARCHAR(255) NOT NULL,
			key VARCHAR(255) NOT NULL,
			value BYTEA NOT NULL,
			processed BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		);
	`)
	if err != nil {
		return err
	}

	// Create outbox table for transactional outbox pattern
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS outbox (
			id SERIAL PRIMARY KEY,
			message_id VARCHAR(255) UNIQUE NOT NULL,
			topic VARCHAR(255) NOT NULL,
			key VARCHAR(255) NOT NULL,
			value BYTEA NOT NULL,
			sent BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		);
	`)
	return err
}

// CreateAccount creates a new account for a user
func (r *PostgresRepository) CreateAccount(ctx context.Context, userID int64) (*models.Account, error) {
	// Check if account already exists
	var count int
	err := r.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM accounts WHERE user_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check if account exists: %w", err)
	}

	if count > 0 {
		return nil, errors.New("account already exists for this user")
	}

	// Create new account
	now := time.Now()
	account := &models.Account{
		UserID:    userID,
		Balance:   0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	query := `
		INSERT INTO accounts (user_id, balance, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err = r.db.QueryRowContext(
		ctx, query, account.UserID, account.Balance, account.CreatedAt, account.UpdatedAt,
	).Scan(&account.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return account, nil
}

// GetAccountByUserID retrieves an account by user ID
func (r *PostgresRepository) GetAccountByUserID(ctx context.Context, userID int64) (*models.Account, error) {
	account := &models.Account{}
	query := "SELECT * FROM accounts WHERE user_id = $1"
	err := r.db.GetContext(ctx, account, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("account not found")
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	return account, nil
}

// UpdateBalance updates the balance of an account
// It uses a FOR UPDATE lock to prevent concurrent updates
func (r *PostgresRepository) UpdateBalance(ctx context.Context, userID int64, amount float64) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Lock the account row for update
	var currentBalance float64
	err = tx.QueryRowContext(
		ctx,
		"SELECT balance FROM accounts WHERE user_id = $1 FOR UPDATE",
		userID,
	).Scan(&currentBalance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("account not found")
		}
		return fmt.Errorf("failed to get account balance: %w", err)
	}

	// Check if there's enough balance for debit operations
	newBalance := currentBalance + amount
	if newBalance < 0 {
		return errors.New("insufficient funds")
	}

	// Update the balance
	_, err = tx.ExecContext(
		ctx,
		"UPDATE accounts SET balance = $1, updated_at = $2 WHERE user_id = $3",
		newBalance, time.Now(), userID,
	)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	return tx.Commit()
}

// SaveInboxMessage saves a message to the inbox table
func (r *PostgresRepository) SaveInboxMessage(ctx context.Context, tx *sqlx.Tx, message *models.InboxMessage) error {
	query := `
		INSERT INTO inbox (message_id, topic, key, value, processed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (message_id) DO NOTHING
	`
	
	var err error
	if tx != nil {
		_, err = tx.ExecContext(
			ctx, query,
			message.MessageID, message.Topic, message.Key, message.Value,
			message.Processed, message.CreatedAt, message.UpdatedAt,
		)
	} else {
		_, err = r.db.ExecContext(
			ctx, query,
			message.MessageID, message.Topic, message.Key, message.Value,
			message.Processed, message.CreatedAt, message.UpdatedAt,
		)
	}
	
	if err != nil {
		return fmt.Errorf("failed to save inbox message: %w", err)
	}
	
	return nil
}

// GetUnprocessedInboxMessages retrieves unprocessed messages from the inbox
func (r *PostgresRepository) GetUnprocessedInboxMessages(ctx context.Context, limit int) ([]*models.InboxMessage, error) {
	query := `
		SELECT * FROM inbox
		WHERE processed = false
		ORDER BY created_at ASC
		LIMIT $1
	`
	
	var messages []*models.InboxMessage
	err := r.db.SelectContext(ctx, &messages, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get unprocessed inbox messages: %w", err)
	}
	
	return messages, nil
}

// MarkInboxMessageAsProcessed marks an inbox message as processed
func (r *PostgresRepository) MarkInboxMessageAsProcessed(ctx context.Context, tx *sqlx.Tx, messageID string) error {
	query := `
		UPDATE inbox
		SET processed = true, updated_at = $1
		WHERE message_id = $2
	`
	
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, time.Now(), messageID)
	} else {
		_, err = r.db.ExecContext(ctx, query, time.Now(), messageID)
	}
	
	if err != nil {
		return fmt.Errorf("failed to mark inbox message as processed: %w", err)
	}
	
	return nil
}

// SaveOutboxMessage saves a message to the outbox table
func (r *PostgresRepository) SaveOutboxMessage(ctx context.Context, tx *sqlx.Tx, message *models.OutboxMessage) error {
	query := `
		INSERT INTO outbox (message_id, topic, key, value, sent, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	
	var err error
	if tx != nil {
		_, err = tx.ExecContext(
			ctx, query,
			message.MessageID, message.Topic, message.Key, message.Value,
			message.Sent, message.CreatedAt, message.UpdatedAt,
		)
	} else {
		_, err = r.db.ExecContext(
			ctx, query,
			message.MessageID, message.Topic, message.Key, message.Value,
			message.Sent, message.CreatedAt, message.UpdatedAt,
		)
	}
	
	if err != nil {
		return fmt.Errorf("failed to save outbox message: %w", err)
	}
	
	return nil
}

// GetUnsendOutboxMessages retrieves unsent messages from the outbox
func (r *PostgresRepository) GetUnsendOutboxMessages(ctx context.Context, limit int) ([]*models.OutboxMessage, error) {
	query := `
		SELECT * FROM outbox
		WHERE sent = false
		ORDER BY created_at ASC
		LIMIT $1
	`
	
	var messages []*models.OutboxMessage
	err := r.db.SelectContext(ctx, &messages, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get unsent outbox messages: %w", err)
	}
	
	return messages, nil
}

// MarkOutboxMessageAsSent marks an outbox message as sent
func (r *PostgresRepository) MarkOutboxMessageAsSent(ctx context.Context, messageID string) error {
	query := `
		UPDATE outbox
		SET sent = true, updated_at = $1
		WHERE message_id = $2
	`
	
	_, err := r.db.ExecContext(ctx, query, time.Now(), messageID)
	if err != nil {
		return fmt.Errorf("failed to mark outbox message as sent: %w", err)
	}
	
	return nil
}

// BeginTx begins a new transaction
func (r *PostgresRepository) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return r.db.BeginTxx(ctx, nil)
}

// CommitTx commits a transaction
func (r *PostgresRepository) CommitTx(tx *sqlx.Tx) error {
	return tx.Commit()
}

// RollbackTx rolls back a transaction
func (r *PostgresRepository) RollbackTx(tx *sqlx.Tx) error {
	return tx.Rollback()
}