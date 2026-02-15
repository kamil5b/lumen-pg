package authentication

import (
	"context"
	"fmt"
	"time"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *AuthenticationUseCaseImplementation) RefreshSession(ctx context.Context, sessionID string) (*domain.Session, error) {
	// Validate the session first
	session, err := u.sessionRepo.ValidateSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate session: %w", err)
	}

	if session == nil {
		return nil, fmt.Errorf("session not found")
	}

	// Extend the expiration time
	session.ExpiresAt = time.Now().Add(24 * time.Hour)

	// Update the session with new expiration time
	err = u.sessionRepo.UpdateSession(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	return session, nil
}
