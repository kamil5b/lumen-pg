package authentication

import (
	"context"
	"errors"
)

func (u *AuthenticationUseCaseImplementation) GetFirstAccessibleSchema(ctx context.Context, username, database string) (string, error) {
	return "", errors.New("not implemented")
}
