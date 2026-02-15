package security

import (
	"context"
	"errors"
)

func (u *SecurityUseCaseImplementation) EncryptPassword(ctx context.Context, password string) (string, error) {
	return "", errors.New("not implemented")
}
