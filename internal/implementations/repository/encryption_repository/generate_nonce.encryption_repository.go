package encryption_repository

import (
	"context"
	"errors"
)

func (e *EncryptionRepositoryImplementation) GenerateNonce(ctx context.Context) (string, error) {
	return "", errors.New("not implemented")
}
