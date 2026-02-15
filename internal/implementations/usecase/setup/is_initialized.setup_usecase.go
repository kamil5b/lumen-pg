package setup

import (
	"context"
	"errors"
)

func (u *SetupUseCaseImplementation) IsInitialized(ctx context.Context) (bool, error) {
	return false, errors.New("not implemented")
}
