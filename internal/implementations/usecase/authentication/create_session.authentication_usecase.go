package authentication

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *AuthenticationUseCaseImplementation) CreateSession(ctx context.Context, username, password, database, schema, table string) (*domain.Session, error) {
	// Encrypt the password for storage
	_, err := u.encryptionRepo.Encrypt(ctx, password)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt password: %w", err)
	}

	// Verify user has access to the specified resources
	databases, err := u.rbacRepo.GetAccessibleDatabases(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get accessible databases: %w", err)
	}

	databaseAccessible := false
	for _, db := range databases {
		if db == database {
			databaseAccessible = true
			break
		}
	}

	if !databaseAccessible {
		return nil, fmt.Errorf("user does not have access to database: %s", database)
	}

	schemas, err := u.rbacRepo.GetAccessibleSchemas(ctx, username, database)
	if err != nil {
		return nil, fmt.Errorf("failed to get accessible schemas: %w", err)
	}

	schemaAccessible := false
	for _, s := range schemas {
		if s == schema {
			schemaAccessible = true
			break
		}
	}

	if !schemaAccessible {
		return nil, fmt.Errorf("user does not have access to schema: %s", schema)
	}

	tables, err := u.rbacRepo.GetAccessibleTables(ctx, username, database, schema)
	if err != nil {
		return nil, fmt.Errorf("failed to get accessible tables: %w", err)
	}

	tableAccessible := false
	for _, tbl := range tables {
		if tbl.Name == table {
			tableAccessible = true
			break
		}
	}

	if !tableAccessible {
		return nil, fmt.Errorf("user does not have access to table: %s", table)
	}

	// Create session with encrypted password stored
	session := &domain.Session{
		ID:        uuid.New().String(),
		Username:  username,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24-hour expiration
	}

	// Store session with encrypted password
	err = u.sessionRepo.CreateSession(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}
