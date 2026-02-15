package encryption_repository

import (
	"context"
	"errors"
)

func (e *EncryptionRepositoryImplementation) GenerateSignature(ctx context.Context, data string) (string, error) {
	return "", errors.New("not implemented")
}
