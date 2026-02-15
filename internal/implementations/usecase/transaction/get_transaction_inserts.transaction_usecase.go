package transaction

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *TransactionUseCaseImplementation) GetTransactionInserts(ctx context.Context, username string) ([]domain.RowInsert, error) {
	// Get the active transaction for the user
	txn, err := u.transactionRepo.GetUserTransaction(ctx, username)
	if err != nil {
		return nil, err
	}

	if txn == nil {
		return nil, domain.ErrNoActiveTransaction
	}

	// Get the row insertions for the transaction
	return u.transactionRepo.GetRowInserts(ctx, username)
}
