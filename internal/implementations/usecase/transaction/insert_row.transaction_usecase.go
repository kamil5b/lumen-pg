package transaction

import (
	"context"
	"errors"
)

func (u *TransactionUseCaseImplementation) InsertRow(ctx context.Context, username, database, schema, table string, values map[string]interface{}) error {
	return errors.New("not implemented")
}
