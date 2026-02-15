package setup

import (
	"context"
	"errors"
)

func (u *SetupUseCaseImplementation) ValidateConnectionString(ctx context.Context, connString string) (bool, error) {
	return false, errors.New("not implemented")
}
