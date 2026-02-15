package transaction

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *TransactionUseCaseImplementation) RollbackTransaction(ctx context.Context, username string) error {
	// Get the active transaction for the user
	txn, err := u.transactionRepo.GetUserTransaction(ctx, username)
	if err != nil {
		return err
	}

	if txn == nil {
		return domain.ErrNoActiveTransaction
	}

	// Delete the transaction to effectively rollback all buffered changes
	return u.transactionRepo.DeleteTransaction(ctx, txn.ID)
}
