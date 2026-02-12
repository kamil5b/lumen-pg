package repository

import (
	"context"
	"github.com/kamil5b/lumen-pg/internal/domain"
)

// TransactionRepository handles database transactions
type TransactionRepository interface {
	// StartTransaction starts a new transaction
	StartTransaction(ctx context.Context, username string, tableName string) (*domain.Transaction, error)
	
	// BufferOperation adds an operation to the transaction buffer
	BufferOperation(ctx context.Context, txnID string, op domain.TransactionOperation) error
	
	// CommitTransaction commits all buffered operations atomically
	CommitTransaction(ctx context.Context, txnID string) error
	
	// RollbackTransaction rolls back all buffered operations
	RollbackTransaction(ctx context.Context, txnID string) error
	
	// GetTransaction retrieves an active transaction
	GetTransaction(ctx context.Context, txnID string) (*domain.Transaction, error)
}
