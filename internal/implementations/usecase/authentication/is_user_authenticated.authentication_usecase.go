package authentication

import (
	"context"
	"fmt"
)

func (u *AuthenticationUseCaseImplementation) IsUserAuthenticated(ctx context.Context, sessionID string) (bool, error) {
	session, err := u.sessionRepo.ValidateSession(ctx, sessionID)
	if err != nil {
		return false, fmt.Errorf("failed to validate session: %w", err)
	}

	if session == nil {
		return false, fmt.Errorf("session not found")
	}

	return true, nil
}
