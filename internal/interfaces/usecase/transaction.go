package usecase

import (
	"context"
	"github.com/kamil5b/lumen-pg/internal/domain"
)

// TransactionUseCase handles transaction management operations
type TransactionUseCase interface {
	StartTransaction(ctx context.Context, username string, tableName string) (*domain.Transaction, error)
	BufferEdit(ctx context.Context, txnID string, op domain.TransactionOperation) error
	CommitTransaction(ctx context.Context, txnID string) error
	RollbackTransaction(ctx context.Context, txnID string) error
}
