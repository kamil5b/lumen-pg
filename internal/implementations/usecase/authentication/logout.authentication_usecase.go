package authentication

import (
	"context"
	"fmt"
)

func (u *AuthenticationUseCaseImplementation) Logout(ctx context.Context, sessionID string) error {
	err := u.sessionRepo.DeleteSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}
