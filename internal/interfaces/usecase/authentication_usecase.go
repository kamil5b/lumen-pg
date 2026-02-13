package usecase

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// AuthenticationUseCase defines operations for user authentication and session management
type AuthenticationUseCase interface {
	// ValidateLoginForm validates login form input (username and password)
	ValidateLoginForm(ctx context.Context, req domain.LoginRequest) ([]domain.ValidationError, error)

	// ProbeConnection tests if a user can connect to their first accessible database
	ProbeConnection(ctx context.Context, username, password string) (bool, error)

	// Login authenticates a user with their PostgreSQL credentials
	Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error)

	// CreateSession creates a new session after successful authentication
	CreateSession(ctx context.Context, username, password, database, schema, table string) (*domain.Session, error)

	// GetUserAccessibleResources returns databases, schemas, and tables accessible by a user
	GetUserAccessibleResources(ctx context.Context, username string) (*domain.RoleMetadata, error)

	// GetFirstAccessibleDatabase returns the first database accessible by a user
	GetFirstAccessibleDatabase(ctx context.Context, username string) (string, error)

	// GetFirstAccessibleSchema returns the first schema accessible by a user in a database
	GetFirstAccessibleSchema(ctx context.Context, username, database string) (string, error)

	// GetFirstAccessibleTable returns the first table accessible by a user in a schema
	GetFirstAccessibleTable(ctx context.Context, username, database, schema string) (string, error)

	// ValidateSession validates if a session is still active and not expired
	ValidateSession(ctx context.Context, sessionID string) (*domain.Session, error)

	// RefreshSession extends a session's expiration time
	RefreshSession(ctx context.Context, sessionID string) (*domain.Session, error)

	// Logout invalidates a user's session
	Logout(ctx context.Context, sessionID string) error

	// GetSessionUser returns the user associated with a session
	GetSessionUser(ctx context.Context, sessionID string) (*domain.User, error)

	// IsUserAuthenticated checks if a user has a valid session
	IsUserAuthenticated(ctx context.Context, sessionID string) (bool, error)

	// ReAuthenticateWithPassword re-authenticates a user using their stored password cookie
	ReAuthenticateWithPassword(ctx context.Context, username, encryptedPassword string) (bool, error)
}
