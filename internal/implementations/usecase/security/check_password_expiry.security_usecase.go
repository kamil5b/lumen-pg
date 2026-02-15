package security

import (
	"context"
	"errors"
)

func (u *SecurityUseCaseImplementation) CheckPasswordExpiry(ctx context.Context, username string, encryptedPassword string) (bool, error) {
	return false, errors.New("not implemented")
}
