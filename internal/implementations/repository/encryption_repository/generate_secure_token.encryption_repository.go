package encryption_repository

import (
	"context"
	"errors"
)

func (e *EncryptionRepositoryImplementation) GenerateSecureToken(ctx context.Context, length int) (string, error) {
	return "", errors.New("not implemented")
}
