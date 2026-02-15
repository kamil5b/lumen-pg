package security

import (
	"context"
	"errors"
)

func (u *SecurityUseCaseImplementation) ValidateQueryForInjection(ctx context.Context, query string) (bool, error) {
	return false, errors.New("not implemented")
}
