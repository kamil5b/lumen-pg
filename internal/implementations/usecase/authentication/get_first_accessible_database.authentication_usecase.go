package authentication

import (
	"context"
	"errors"
)

func (u *AuthenticationUseCaseImplementation) GetFirstAccessibleDatabase(ctx context.Context, username string) (string, error) {
	return "", errors.New("not implemented")
}
