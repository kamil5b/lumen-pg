package security

import (
	"context"
	"errors"
)

func (u *SecurityUseCaseImplementation) CheckSessionTimeout(ctx context.Context, sessionID string) (bool, error) {
	return false, errors.New("not implemented")
}
