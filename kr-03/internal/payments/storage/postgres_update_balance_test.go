package storage

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestUpdateBalanceAccountNotFound(t *testing.T) {
	// Setup
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repo := &PostgresRepository{db: sqlxDB}

	ctx := context.Background()
	userID := int64(2)
	amount := 50.0

	// Test account not found
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT balance FROM accounts").
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()

	err = repo.UpdateBalance(ctx, userID, amount)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "account not found")

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
