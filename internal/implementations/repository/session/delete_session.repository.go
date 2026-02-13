package session

import (
	"context"
	"errors"
)

// DeleteSession removes a session
func (s *SessionRepository) DeleteSession(ctx context.Context, sessionToken string) error {
	return errors.New("not implemented yet")
}
