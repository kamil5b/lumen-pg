package security

import (
	"context"
	"fmt"
)

func (u *SecurityUseCaseImplementation) DecryptPassword(ctx context.Context, encryptedPassword string) (string, error) {
	// Decrypt the password using the encryption repository
	decrypted, err := u.encryptionRepo.Decrypt(ctx, encryptedPassword)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt password: %w", err)
	}

	return decrypted, nil
}
