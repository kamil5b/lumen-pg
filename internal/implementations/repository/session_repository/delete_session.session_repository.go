package session_repository

import (
	"context"
	"errors"
)

func (s *SessionRepositoryImplementation) DeleteSession(ctx context.Context, sessionID string) error {
	return errors.New("not implemented")
}
