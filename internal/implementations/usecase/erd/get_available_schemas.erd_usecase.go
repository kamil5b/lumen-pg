package erd

import (
	"context"
	"errors"
)

func (u *ERDUseCaseImplementation) GetAvailableSchemas(ctx context.Context, username, database string) ([]string, error) {
	return nil, errors.New("not implemented")
}
