package repository

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// SessionRepository defines operations for managing user sessions
type SessionRepository interface {
	// CreateSession creates a new user session
	CreateSession(ctx context.Context, session *domain.Session) error

	// GetSession retrieves a session by ID
	GetSession(ctx context.Context, sessionID string) (*domain.Session, error)

	// UpdateSession updates an existing session
	UpdateSession(ctx context.Context, session *domain.Session) error

	// DeleteSession removes a session
	DeleteSession(ctx context.Context, sessionID string) error

	// ValidateSession checks if a session is valid and not expired
	ValidateSession(ctx context.Context, sessionID string) (*domain.Session, error)

	// GetSessionByUsername retrieves the most recent session for a user
	GetSessionByUsername(ctx context.Context, username string) (*domain.Session, error)

	// InvalidateUserSessions invalidates all sessions for a user
	InvalidateUserSessions(ctx context.Context, username string) error

	// InvalidateExpiredSessions removes all expired sessions
	InvalidateExpiredSessions(ctx context.Context) error

	// SessionExists checks if a session exists and is valid
	SessionExists(ctx context.Context, sessionID string) (bool, error)
}
