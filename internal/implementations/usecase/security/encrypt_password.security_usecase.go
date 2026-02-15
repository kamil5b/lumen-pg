package security

import (
	"context"
	"fmt"
)

func (u *SecurityUseCaseImplementation) EncryptPassword(ctx context.Context, password string) (string, error) {
	// Encrypt the password using the encryption repository
	encrypted, err := u.encryptionRepo.Encrypt(ctx, password)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt password: %w", err)
	}

	return encrypted, nil
}
