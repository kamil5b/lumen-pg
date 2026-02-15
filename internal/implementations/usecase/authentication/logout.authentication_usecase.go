package authentication

import (
	"context"
	"errors"
)

func (u *AuthenticationUseCaseImplementation) Logout(ctx context.Context, sessionID string) error {
	return errors.New("not implemented")
}
