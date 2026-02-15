package transaction

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *TransactionUseCaseImplementation) EditCell(ctx context.Context, username, database, schema, table string, rowIndex int, columnName string, newValue interface{}) error {
	// Get the active transaction for the user
	txn, err := u.transactionRepo.GetUserTransaction(ctx, username)
	if err != nil {
		return err
	}

	if txn == nil {
		return domain.ErrNoActiveTransaction
	}

	// Create a row edit
	edit := domain.RowEdit{
		RowIndex:   rowIndex,
		ColumnName: columnName,
		NewValue:   newValue,
	}

	// Add the edit to the transaction
	return u.transactionRepo.AddRowEdit(ctx, username, edit)
}
