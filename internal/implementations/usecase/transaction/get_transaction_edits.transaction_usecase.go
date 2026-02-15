package transaction

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *TransactionUseCaseImplementation) GetTransactionEdits(ctx context.Context, username string) (map[int]domain.RowEdit, error) {
	// Get the active transaction for the user
	txn, err := u.transactionRepo.GetUserTransaction(ctx, username)
	if err != nil {
		return nil, err
	}

	if txn == nil {
		return nil, domain.ErrNoActiveTransaction
	}

	// Retrieve the row edits for this transaction
	return u.transactionRepo.GetRowEdits(ctx, username)
}
