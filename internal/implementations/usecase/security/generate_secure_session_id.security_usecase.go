package security

import (
	"context"
	"errors"
)

func (u *SecurityUseCaseImplementation) GenerateSecureSessionID(ctx context.Context) (string, error) {
	return "", errors.New("not implemented")
}
