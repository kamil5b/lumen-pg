package transaction

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *TransactionUseCaseImplementation) InsertRow(ctx context.Context, username, database, schema, table string, values map[string]interface{}) error {
	// Get the active transaction for the user
	txn, err := u.transactionRepo.GetUserTransaction(ctx, username)
	if err != nil {
		return err
	}

	if txn == nil {
		return domain.ErrNoActiveTransaction
	}

	// Create the row insert
	insert := domain.RowInsert{
		Values: values,
	}

	// Add the row insertion to the transaction
	return u.transactionRepo.AddRowInsert(ctx, username, insert)
}
