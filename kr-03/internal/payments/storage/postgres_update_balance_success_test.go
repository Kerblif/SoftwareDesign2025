package storage

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestUpdateBalanceSuccess(t *testing.T) {
	// Setup
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repo := &PostgresRepository{db: sqlxDB}

	ctx := context.Background()
	userID := int64(1)
	amount := 50.0
	currentBalance := 100.0
	newBalance := currentBalance + amount

	// Test success case
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT balance FROM accounts").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(currentBalance))

	mock.ExpectExec("UPDATE accounts").
		WithArgs(newBalance, sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	err = repo.UpdateBalance(ctx, userID, amount)
	assert.NoError(t, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateBalanceInsufficientFunds(t *testing.T) {
	// Setup
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repo := &PostgresRepository{db: sqlxDB}

	ctx := context.Background()
	userID := int64(1)
	currentBalance := 100.0

	// Test insufficient funds
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT balance FROM accounts").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(currentBalance))

	mock.ExpectRollback()

	err = repo.UpdateBalance(ctx, userID, -200.0) // More than current balance
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient funds")

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateBalanceErrorUpdating(t *testing.T) {
	// Setup
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repo := &PostgresRepository{db: sqlxDB}

	ctx := context.Background()
	userID := int64(1)
	amount := 50.0
	currentBalance := 100.0
	newBalance := currentBalance + amount

	// Test error updating balance
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT balance FROM accounts").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(currentBalance))

	mock.ExpectExec("UPDATE accounts").
		WithArgs(newBalance, sqlmock.AnyArg(), userID).
		WillReturnError(assert.AnError)

	mock.ExpectRollback()

	err = repo.UpdateBalance(ctx, userID, amount)
	assert.Error(t, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
