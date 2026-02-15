package session_repository

import (
	"context"
	"errors"
)

func (s *SessionRepositoryImplementation) InvalidateUserSessions(ctx context.Context, username string) error {
	return errors.New("not implemented")
}
