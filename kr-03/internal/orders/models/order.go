package models

import (
	"time"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	// OrderStatusNew represents a new order
	OrderStatusNew OrderStatus = "NEW"
	// OrderStatusFinished represents a finished order
	OrderStatusFinished OrderStatus = "FINISHED"
	// OrderStatusCancelled represents a cancelled order
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

// Order represents an order
type Order struct {
	ID          int64       `json:"id" db:"id"`
	UserID      int64       `json:"user_id" db:"user_id"`
	Amount      float64     `json:"amount" db:"amount"`
	Description string      `json:"description" db:"description"`
	Status      OrderStatus `json:"status" db:"status"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
}

// CreateOrderRequest represents the request to create a new order
type CreateOrderRequest struct {
	UserID      int64   `json:"user_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description" binding:"required"`
}

// OrderResponse represents the response with order details
type OrderResponse struct {
	ID          int64       `json:"id"`
	UserID      int64       `json:"user_id"`
	Amount      float64     `json:"amount"`
	Description string      `json:"description"`
	Status      OrderStatus `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
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