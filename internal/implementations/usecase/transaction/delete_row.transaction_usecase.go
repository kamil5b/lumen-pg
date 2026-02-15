package transaction

import (
	"context"
	"errors"
)

func (u *TransactionUseCaseImplementation) DeleteRow(ctx context.Context, username, database, schema, table string, rowIndex int) error {
	return errors.New("not implemented")
}
