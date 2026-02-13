package repository

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// TransactionRepository defines operations for managing user transactions
type TransactionRepository interface {
	// CreateTransaction creates a new transaction state
	CreateTransaction(ctx context.Context, transaction *domain.TransactionState) error

	// GetTransaction retrieves a transaction by ID
	GetTransaction(ctx context.Context, transactionID string) (*domain.TransactionState, error)

	// GetUserTransaction retrieves the active transaction for a user
	GetUserTransaction(ctx context.Context, username string) (*domain.TransactionState, error)

	// UpdateTransaction updates an existing transaction
	UpdateTransaction(ctx context.Context, transaction *domain.TransactionState) error

	// DeleteTransaction removes a transaction
	DeleteTransaction(ctx context.Context, transactionID string) error

	// AddRowEdit buffers a cell edit in a transaction
	AddRowEdit(ctx context.Context, transactionID string, edit domain.RowEdit) error

	// AddRowDelete buffers a row deletion in a transaction
	AddRowDelete(ctx context.Context, transactionID string, rowIndex int) error

	// AddRowInsert buffers a new row insertion in a transaction
	AddRowInsert(ctx context.Context, transactionID string, insert domain.RowInsert) error

	// GetRowEdits retrieves all buffered edits for a transaction
	GetRowEdits(ctx context.Context, transactionID string) (map[int]domain.RowEdit, error)

	// GetRowDeletes retrieves all buffered deletions for a transaction
	GetRowDeletes(ctx context.Context, transactionID string) ([]int, error)

	// GetRowInserts retrieves all buffered insertions for a transaction
	GetRowInserts(ctx context.Context, transactionID string) ([]domain.RowInsert, error)

	// ClearRowEdits clears all buffered edits for a transaction
	ClearRowEdits(ctx context.Context, transactionID string) error

	// ClearRowDeletes clears all buffered deletions for a transaction
	ClearRowDeletes(ctx context.Context, transactionID string) error

	// ClearRowInserts clears all buffered insertions for a transaction
	ClearRowInserts(ctx context.Context, transactionID string) error

	// TransactionExists checks if a transaction exists and is active
	TransactionExists(ctx context.Context, transactionID string) (bool, error)

	// InvalidateExpiredTransactions removes all expired transactions
	InvalidateExpiredTransactions(ctx context.Context) error
}
