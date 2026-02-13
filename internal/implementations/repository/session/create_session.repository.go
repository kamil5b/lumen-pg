package session

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// CreateSession creates a new user session
func (s *SessionRepository) CreateSession(ctx context.Context, username string, password string, roleMetadata *domain.RoleMetadata) (*domain.Session, error) {
	return nil, errors.New("not implemented yet")
}
