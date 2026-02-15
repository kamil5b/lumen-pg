package authentication

import (
	"context"
	"errors"
)

func (u *AuthenticationUseCaseImplementation) GetFirstAccessibleTable(ctx context.Context, username, database, schema string) (string, error) {
	return "", errors.New("not implemented")
}
