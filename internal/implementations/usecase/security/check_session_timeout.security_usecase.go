package security

import (
	"context"
	"fmt"
	"time"
)

func (u *SecurityUseCaseImplementation) CheckSessionTimeout(ctx context.Context, sessionID string) (bool, error) {
	// Get the session from the session repository
	session, err := u.sessionRepo.GetSession(ctx, sessionID)
	if err != nil {
		return false, fmt.Errorf("failed to get session: %w", err)
	}

	if session == nil {
		return true, nil // Session doesn't exist, treat as timed out
	}

	// Check if the session has expired
	return time.Now().After(session.ExpiresAt), nil
}
