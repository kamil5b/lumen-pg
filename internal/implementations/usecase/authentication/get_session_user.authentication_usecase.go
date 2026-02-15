package authentication

import (
	"context"
	"fmt"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *AuthenticationUseCaseImplementation) GetSessionUser(ctx context.Context, sessionID string) (*domain.User, error) {
	session, err := u.sessionRepo.ValidateSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate session: %w", err)
	}

	if session == nil {
		return nil, fmt.Errorf("session not found")
	}

	// Get user metadata to populate the User struct
	metadata, err := u.metadataRepo.GetRoleMetadata(ctx, session.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to get role metadata: %w", err)
	}

	if metadata == nil {
		return nil, fmt.Errorf("no metadata found for user: %s", session.Username)
	}

	// Get first accessible database, schema, and table
	database := ""
	schema := ""
	table := ""

	if len(metadata.AccessibleDatabases) > 0 {
		database = metadata.AccessibleDatabases[0]
	}

	if len(metadata.AccessibleSchemas) > 0 {
		schema = metadata.AccessibleSchemas[0]
	}

	if len(metadata.AccessibleTables) > 0 {
		table = metadata.AccessibleTables[0].Name
	}

	return &domain.User{
		Username:     session.Username,
		DatabaseName: database,
		SchemaName:   schema,
		TableName:    table,
	}, nil
}
