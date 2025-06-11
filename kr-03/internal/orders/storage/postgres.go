package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"cbd/internal/orders/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Repository defines the interface for the orders storage
type Repository interface {
	// Order operations
	CreateOrder(ctx context.Context, tx *sqlx.Tx, order *models.Order) error
	GetOrderByID(ctx context.Context, id int64) (*models.Order, error)
	GetOrdersByUserID(ctx context.Context, userID int64) ([]*models.Order, error)
	UpdateOrderStatus(ctx context.Context, id int64, status models.OrderStatus) error
	
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
	// Create orders table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL,
			amount DECIMAL(10, 2) NOT NULL,
			description TEXT NOT NULL,
			status VARCHAR(20) NOT NULL,
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

// CreateOrder creates a new order
func (r *PostgresRepository) CreateOrder(ctx context.Context, tx *sqlx.Tx, order *models.Order) error {
	query := `
		INSERT INTO orders (user_id, amount, description, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	
	var err error
	if tx != nil {
		err = tx.QueryRowContext(
			ctx, query,
			order.UserID, order.Amount, order.Description, order.Status,
			order.CreatedAt, order.UpdatedAt,
		).Scan(&order.ID)
	} else {
		err = r.db.QueryRowContext(
			ctx, query,
			order.UserID, order.Amount, order.Description, order.Status,
			order.CreatedAt, order.UpdatedAt,
		).Scan(&order.ID)
	}
	
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}
	
	return nil
}

// GetOrderByID retrieves an order by ID
func (r *PostgresRepository) GetOrderByID(ctx context.Context, id int64) (*models.Order, error) {
	order := &models.Order{}
	query := "SELECT * FROM orders WHERE id = $1"
	err := r.db.GetContext(ctx, order, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	return order, nil
}

// GetOrdersByUserID retrieves orders by user ID
func (r *PostgresRepository) GetOrdersByUserID(ctx context.Context, userID int64) ([]*models.Order, error) {
	var orders []*models.Order
	query := "SELECT * FROM orders WHERE user_id = $1 ORDER BY created_at DESC"
	err := r.db.SelectContext(ctx, &orders, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	return orders, nil
}

// UpdateOrderStatus updates the status of an order
func (r *PostgresRepository) UpdateOrderStatus(ctx context.Context, id int64, status models.OrderStatus) error {
	query := "UPDATE orders SET status = $1, updated_at = $2 WHERE id = $3"
	_, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
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