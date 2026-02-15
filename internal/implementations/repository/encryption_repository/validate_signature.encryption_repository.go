package encryption_repository

import (
	"context"
	"errors"
)

func (e *EncryptionRepositoryImplementation) ValidateSignature(ctx context.Context, data string, signature string) (bool, error) {
	return false, errors.New("not implemented")
}
