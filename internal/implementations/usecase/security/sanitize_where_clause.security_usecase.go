package security

import (
	"context"
	"fmt"
)

func (u *SecurityUseCaseImplementation) SanitizeWhereClause(ctx context.Context, whereClause string) (string, error) {
	// Sanitize the WHERE clause by encrypting it to prevent SQL injection
	sanitized, err := u.encryptionRepo.Encrypt(ctx, whereClause)
	if err != nil {
		return "", fmt.Errorf("failed to sanitize where clause: %w", err)
	}

	return sanitized, nil
}
