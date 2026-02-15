package encryption_repository

import (
	"context"
	"errors"
)

func (e *EncryptionRepositoryImplementation) Encrypt(ctx context.Context, plaintext string) (string, error) {
	return "", errors.New("not implemented")
}
