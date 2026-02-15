package authentication

import (
	"context"
	"errors"
)

func (u *AuthenticationUseCaseImplementation) IsUserAuthenticated(ctx context.Context, sessionID string) (bool, error) {
	return false, errors.New("not implemented")
}
