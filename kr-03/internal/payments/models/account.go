package models

import (
	"time"
)

// Account represents a user's payment account
type Account struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Balance   float64   `json:"balance" db:"balance"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateAccountRequest represents the request to create a new account
type CreateAccountRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
}

// DepositRequest represents the request to deposit money to an account
type DepositRequest struct {
	UserID int64   `json:"user_id" binding:"required"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

// BalanceResponse represents the response with account balance
type BalanceResponse struct {
	UserID  int64   `json:"user_id"`
	Balance float64 `json:"balance"`
}

// PaymentMessage represents a payment message sent via Kafka
type PaymentMessage struct {
	OrderID     int64   `json:"order_id"`
	UserID      int64   `json:"user_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

// PaymentResult represents the result of a payment operation
type PaymentResult struct {
	OrderID int64  `json:"order_id"`
	Status  string `json:"status"` // "success" or "fail"
	Reason  string `json:"reason,omitempty"`
}

// InboxMessage represents a message in the transactional inbox
type InboxMessage struct {
	ID        int64     `db:"id"`
	MessageID string    `db:"message_id"`
	Topic     string    `db:"topic"`
	Key       string    `db:"key"`
	Value     []byte    `db:"value"`
	Processed bool      `db:"processed"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// OutboxMessage represents a message in the transactional outbox
type OutboxMessage struct {
	ID        int64     `db:"id"`
	MessageID string    `db:"message_id"`
	Topic     string    `db:"topic"`
	Key       string    `db:"key"`
	Value     []byte    `db:"value"`
	Sent      bool      `db:"sent"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
