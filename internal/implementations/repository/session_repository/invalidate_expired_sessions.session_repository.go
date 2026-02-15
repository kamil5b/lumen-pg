package session_repository

import (
	"context"
	"errors"
)

func (s *SessionRepositoryImplementation) InvalidateExpiredSessions(ctx context.Context) error {
	return errors.New("not implemented")
}
