package security

import (
	"context"
	"fmt"
)

func (u *SecurityUseCaseImplementation) ValidateQueryForInjection(ctx context.Context, query string) (bool, error) {
	// Validate the query signature to detect SQL injection attempts
	// Returns true if injection detected, false if query is safe
	hasInjection, err := u.encryptionRepo.ValidateSignature(ctx, query, "")
	if err != nil {
		return false, fmt.Errorf("failed to validate query for injection: %w", err)
	}

	// ValidateSignature returns true for safe queries, so we need to invert it
	return !hasInjection, nil
}
