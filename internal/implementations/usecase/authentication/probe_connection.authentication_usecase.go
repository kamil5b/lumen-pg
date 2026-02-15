package authentication

import (
	"context"
	"errors"
)

func (u *AuthenticationUseCaseImplementation) ProbeConnection(ctx context.Context, username, password string) (bool, error) {
	return false, errors.New("not implemented")
}
