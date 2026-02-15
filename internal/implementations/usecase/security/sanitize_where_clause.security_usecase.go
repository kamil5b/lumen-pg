package security

import (
	"context"
	"errors"
)

func (u *SecurityUseCaseImplementation) SanitizeWhereClause(ctx context.Context, whereClause string) (string, error) {
	return "", errors.New("not implemented")
}
