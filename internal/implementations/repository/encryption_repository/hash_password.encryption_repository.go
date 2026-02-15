package encryption_repository

import (
	"context"
	"errors"
)

func (e *EncryptionRepositoryImplementation) HashPassword(ctx context.Context, password string) (string, error) {
	return "", errors.New("not implemented")
}
