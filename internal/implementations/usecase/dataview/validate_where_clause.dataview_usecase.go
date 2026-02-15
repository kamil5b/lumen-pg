package dataview

import (
	"context"
	"errors"
)

func (u *DataViewUseCaseImplementation) ValidateWhereClause(ctx context.Context, whereClause string) (bool, error) {
	return false, errors.New("not implemented")
}
