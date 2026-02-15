package authentication

import (
	"context"
	"fmt"
)

func (u *AuthenticationUseCaseImplementation) ReAuthenticateWithPassword(ctx context.Context, username, encryptedPassword string) (bool, error) {
	// Decrypt the password
	decryptedPassword, err := u.encryptionRepo.Decrypt(ctx, encryptedPassword)
	if err != nil {
		return false, fmt.Errorf("failed to decrypt password: %w", err)
	}

	// Test the connection with the decrypted password
	success, err := u.ProbeConnection(ctx, username, decryptedPassword)
	if err != nil || !success {
		return false, err
	}

	return true, nil
}
