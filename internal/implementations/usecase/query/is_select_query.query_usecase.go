package query

import (
	"context"
	"errors"
)

func (u *QueryUseCaseImplementation) IsSelectQuery(ctx context.Context, query string) (bool, error) {
	return false, errors.New("not implemented")
}
