package usecase

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// TransactionUseCase defines operations for managing database transactions
type TransactionUseCase interface {
	// StartTransaction creates a new transaction for a user
	StartTransaction(ctx context.Context, username, database, schema, table string) (*domain.TransactionState, error)

	// GetActiveTransaction retrieves the active transaction for a user
	GetActiveTransaction(ctx context.Context, username string) (*domain.TransactionState, error)

	// CheckActiveTransaction checks if a user already has an active transaction
	CheckActiveTransaction(ctx context.Context, username string) (bool, error)

	// CommitTransaction commits all buffered changes in a transaction
	CommitTransaction(ctx context.Context, username string) error

	// RollbackTransaction rolls back all buffered changes in a transaction
	RollbackTransaction(ctx context.Context, username string) error

	// EditCell buffers an edit to a table cell
	EditCell(ctx context.Context, username, database, schema, table string, rowIndex int, columnName string, newValue interface{}) error

	// DeleteRow buffers a row deletion
	DeleteRow(ctx context.Context, username, database, schema, table string, rowIndex int) error

	// InsertRow buffers a new row insertion
	InsertRow(ctx context.Context, username, database, schema, table string, values map[string]interface{}) error

	// GetTransactionEdits retrieves all buffered edits for an active transaction
	GetTransactionEdits(ctx context.Context, username string) (map[int]domain.RowEdit, error)

	// GetTransactionDeletes retrieves all buffered deletions for an active transaction
	GetTransactionDeletes(ctx context.Context, username string) ([]int, error)

	// GetTransactionInserts retrieves all buffered insertions for an active transaction
	GetTransactionInserts(ctx context.Context, username string) ([]domain.RowInsert, error)

	// GetTransactionRemainingTime returns the remaining time for an active transaction
	GetTransactionRemainingTime(ctx context.Context, username string) (int64, error)

	// IsTransactionExpired checks if a transaction has expired
	IsTransactionExpired(ctx context.Context, username string) (bool, error)

	// CancelExpiredTransactions cancels all expired transactions
	CancelExpiredTransactions(ctx context.Context) error
}
