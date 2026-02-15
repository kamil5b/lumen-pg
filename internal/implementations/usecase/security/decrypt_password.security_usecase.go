package security

import (
	"context"
	"errors"
)

func (u *SecurityUseCaseImplementation) DecryptPassword(ctx context.Context, encryptedPassword string) (string, error) {
	return "", errors.New("not implemented")
}
