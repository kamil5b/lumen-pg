package transaction

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *TransactionUseCaseImplementation) GetTransactionDeletes(ctx context.Context, username string) ([]int, error) {
	// Get the active transaction for the user
	txn, err := u.transactionRepo.GetUserTransaction(ctx, username)
	if err != nil {
		return nil, err
	}

	if txn == nil {
		return nil, domain.ErrNoActiveTransaction
	}

	// Retrieve the row deletions for this transaction
	return u.transactionRepo.GetRowDeletes(ctx, username)
}
