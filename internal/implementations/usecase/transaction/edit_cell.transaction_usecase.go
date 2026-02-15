package transaction

import (
	"context"
	"errors"
)

func (u *TransactionUseCaseImplementation) EditCell(ctx context.Context, username, database, schema, table string, rowIndex int, columnName string, newValue interface{}) error {
	return errors.New("not implemented")
}
