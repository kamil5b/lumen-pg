package dataview

import (
	"context"
	"errors"
)

func (u *DataViewUseCaseImplementation) GetPrimaryKeyInfo(ctx context.Context, username, database, schema, table string) ([]string, error) {
	return nil, errors.New("not implemented")
}
