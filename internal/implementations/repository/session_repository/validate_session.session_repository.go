package session_repository

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (s *SessionRepositoryImplementation) ValidateSession(ctx context.Context, sessionID string) (*domain.Session, error) {
	return nil, errors.New("not implemented")
}
