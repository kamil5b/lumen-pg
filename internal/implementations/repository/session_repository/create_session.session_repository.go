package session_repository

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (s *SessionRepositoryImplementation) CreateSession(ctx context.Context, session *domain.Session) error {
	return errors.New("not implemented")
}
