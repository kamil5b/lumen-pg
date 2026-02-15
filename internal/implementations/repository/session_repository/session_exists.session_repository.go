package session_repository

import (
	"context"
	"errors"
)

func (s *SessionRepositoryImplementation) SessionExists(ctx context.Context, sessionID string) (bool, error) {
	return false, errors.New("not implemented")
}
