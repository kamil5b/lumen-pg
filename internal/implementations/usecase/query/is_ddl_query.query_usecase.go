package query

import (
	"context"
	"errors"
)

func (u *QueryUseCaseImplementation) IsDDLQuery(ctx context.Context, query string) (bool, error) {
	return false, errors.New("not implemented")
}
