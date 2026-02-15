package erd

import (
	"context"
	"errors"
)

func (u *ERDUseCaseImplementation) IsSchemaEmpty(ctx context.Context, username, database, schema string) (bool, error) {
	return false, errors.New("not implemented")
}
