package encryption_repository

import (
	"context"
	"errors"
)

func (e *EncryptionRepositoryImplementation) Decrypt(ctx context.Context, ciphertext string) (string, error) {
	return "", errors.New("not implemented")
}
