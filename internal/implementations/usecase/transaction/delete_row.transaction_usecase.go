package transaction

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *TransactionUseCaseImplementation) DeleteRow(ctx context.Context, username, database, schema, table string, rowIndex int) error {
	// Get the active transaction for the user
	txn, err := u.transactionRepo.GetUserTransaction(ctx, username)
	if err != nil {
		return err
	}

	if txn == nil {
		return domain.ErrNoActiveTransaction
	}

	// Add the row deletion to the transaction
	return u.transactionRepo.AddRowDelete(ctx, username, rowIndex)
}
