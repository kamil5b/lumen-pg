package security

import (
	"context"
	"errors"
)

func (u *SecurityUseCaseImplementation) IsHTTPSEnabled(ctx context.Context) (bool, error) {
	return false, errors.New("not implemented")
}
