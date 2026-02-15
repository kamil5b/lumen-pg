package authentication

import (
	"context"
	"errors"
)

func (u *AuthenticationUseCaseImplementation) ReAuthenticateWithPassword(ctx context.Context, username, encryptedPassword string) (bool, error) {
	return false, errors.New("not implemented")
}
