package dataview

import (
	"context"
	"errors"
)

func (u *DataViewUseCaseImplementation) GetChildTableRowCount(ctx context.Context, username, database, schema, childTable, parentTable string, fkColumn, pkValue string) (int64, error) {
	return 0, errors.New("not implemented")
}
