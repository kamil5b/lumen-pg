package transaction

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *TransactionUseCaseImplementation) StartTransaction(ctx context.Context, username, database, schema, table string) (*domain.TransactionState, error) {
	// Check if user already has an active transaction
	existingTxn, err := u.transactionRepo.GetUserTransaction(ctx, username)
	if err == nil && existingTxn != nil {
		// User already has an active transaction
		return nil, domain.ErrActiveTransactionExists
	}

	// Create a new transaction state
	txnID := "txn_" + uuid.New().String()
	now := time.Now()
	expiresAt := now.Add(time.Duration(domain.TransactionTimeout) * time.Second)

	txnState := &domain.TransactionState{
		ID:        txnID,
		Username:  username,
		StartedAt: now,
		ExpiresAt: expiresAt,
		Edits:     make(map[int]domain.RowEdit),
		Deletes:   make([]int, 0),
		Inserts:   make([]domain.RowInsert, 0),
	}

	// Create the transaction in the repository
	if err := u.transactionRepo.CreateTransaction(ctx, txnState); err != nil {
		return nil, err
	}

	// Retrieve the created transaction to confirm
	createdTxn, err := u.transactionRepo.GetTransaction(ctx, txnID)
	if err != nil {
		return nil, err
	}

	return createdTxn, nil
}
