package security

import (
	"context"
	"fmt"
)

func (u *SecurityUseCaseImplementation) GenerateSecureSessionID(ctx context.Context) (string, error) {
	// Generate a secure session ID using the encryption repository
	sessionID, err := u.encryptionRepo.GenerateSecureToken(ctx, 32)
	if err != nil {
		return "", fmt.Errorf("failed to generate secure session ID: %w", err)
	}

	return sessionID, nil
}
