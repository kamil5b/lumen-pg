package setup

import (
	"context"
	"errors"
)

func (u *SetupUseCaseImplementation) TestSuperadminConnection(ctx context.Context, connString string) (bool, error) {
	return false, errors.New("not implemented")
}
