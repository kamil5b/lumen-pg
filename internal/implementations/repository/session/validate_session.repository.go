package session

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// ValidateSession validates an existing session
func (s *SessionRepository) ValidateSession(ctx context.Context, sessionToken string) (*domain.Session, error) {
	return nil, errors.New("not implemented yet")
}
