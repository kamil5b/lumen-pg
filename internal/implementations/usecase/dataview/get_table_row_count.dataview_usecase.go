package dataview

import (
	"context"
	"errors"
)

func (u *DataViewUseCaseImplementation) GetTableRowCount(ctx context.Context, username, database, schema, table string) (int64, error) {
	return 0, errors.New("not implemented")
}
