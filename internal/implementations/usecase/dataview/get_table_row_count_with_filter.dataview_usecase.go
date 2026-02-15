package dataview

import (
	"context"
	"errors"
)

func (u *DataViewUseCaseImplementation) GetTableRowCountWithFilter(ctx context.Context, username, database, schema, table, whereClause string) (int64, error) {
	return 0, errors.New("not implemented")
}
