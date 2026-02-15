package dataview

import (
	"context"
	"errors"
)

func (u *DataViewUseCaseImplementation) IsTableReadOnly(ctx context.Context, username, database, schema, table string) (bool, error) {
	return false, errors.New("not implemented")
}
