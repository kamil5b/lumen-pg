package repository

import (
	"context"
	"github.com/kamil5b/lumen-pg/internal/domain"
)

// SessionRepository handles session management
type SessionRepository interface {
	// CreateSession creates a new user session
	CreateSession(ctx context.Context, username string, password string, roleMetadata *domain.RoleMetadata) (*domain.Session, error)
	
	// ValidateSession validates an existing session
	ValidateSession(ctx context.Context, sessionToken string) (*domain.Session, error)
	
	// EncryptPassword encrypts a password for cookie storage
	EncryptPassword(password string) (string, error)
	
	// DecryptPassword decrypts a password from cookie storage
	DecryptPassword(encrypted string) (string, error)
	
	// DeleteSession removes a session
	DeleteSession(ctx context.Context, sessionToken string) error
}
