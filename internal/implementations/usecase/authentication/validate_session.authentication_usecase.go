package authentication

import (
	"context"
	"fmt"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *AuthenticationUseCaseImplementation) ValidateSession(ctx context.Context, sessionID string) (*domain.Session, error) {
	session, err := u.sessionRepo.ValidateSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate session: %w", err)
	}

	if session == nil {
		return nil, fmt.Errorf("session not found")
	}

	return session, nil
}
