package encryption_repository

import (
	"context"
	"errors"
)

func (e *EncryptionRepositoryImplementation) ComparePasswordHash(ctx context.Context, password, hash string) (bool, error) {
	return false, errors.New("not implemented")
}
