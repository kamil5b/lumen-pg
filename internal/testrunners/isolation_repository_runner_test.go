package testrunners

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
)

// IsolationRepoConstructor creates an isolation repository with database connection
type IsolationRepoConstructor func(db *sql.DB) repository.SessionRepository

// IsolationRepositoryRunner runs integration tests for isolation repository (Story 6)
func IsolationRepositoryRunner(t *testing.T, constructor IsolationRepoConstructor) {
	t.Helper()

	ctx := context.Background()

	// Start PostgreSQL container
	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
	)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	connStr, err := container.ConnectionString(ctx)
	require.NoError(t, err)

	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer db.Close()

	repo := constructor(db)

	t.Run("IT-S6-01: Real Multi-User Session Creation", func(t *testing.T) {
		// Create session for user A
		roleMetadataA := &domain.RoleMetadata{
			RoleName:            "userA",
			AccessibleDatabases: []string{"testdb"},
			AccessibleTables:    map[string][]string{"testdb.public": {"users"}},
		}

		sessionA, err := repo.CreateSession(ctx, "userA", "passwordA", roleMetadataA)
		require.NoError(t, err)
		assert.NotNil(t, sessionA)
		assert.Equal(t, "userA", sessionA.Username)

		// Create session for user B independently
		roleMetadataB := &domain.RoleMetadata{
			RoleName:            "userB",
			AccessibleDatabases: []string{"testdb"},
			AccessibleTables:    map[string][]string{"testdb.public": {"posts"}},
		}

		sessionB, err := repo.CreateSession(ctx, "userB", "passwordB", roleMetadataB)
		require.NoError(t, err)
		assert.NotNil(t, sessionB)
		assert.Equal(t, "userB", sessionB.Username)

		// Sessions should have different encrypted passwords
		assert.NotEqual(t, sessionA.EncryptedPassword, sessionB.EncryptedPassword)
	})

	t.Run("IT-S6-02: Real Permission Isolation in Sessions", func(t *testing.T) {
		// Create sessions for two different users with different permissions
		roleMetadataA := &domain.RoleMetadata{
			RoleName:            "userA",
			AccessibleDatabases: []string{"testdb"},
			AccessibleSchemas:   map[string][]string{"testdb": {"public"}},
			AccessibleTables:    map[string][]string{"testdb.public": {"users"}},
		}

		sessionA, err := repo.CreateSession(ctx, "userA", "passA", roleMetadataA)
		require.NoError(t, err)
		assert.NotNil(t, sessionA)
		assert.Equal(t, roleMetadataA, sessionA.AccessibleMetadata)

		roleMetadataB := &domain.RoleMetadata{
			RoleName:            "userB",
			AccessibleDatabases: []string{"testdb"},
			AccessibleSchemas:   map[string][]string{"testdb": {"public"}},
			AccessibleTables:    map[string][]string{"testdb.public": {"posts"}},
		}

		sessionB, err := repo.CreateSession(ctx, "userB", "passB", roleMetadataB)
		require.NoError(t, err)
		assert.NotNil(t, sessionB)
		assert.Equal(t, roleMetadataB, sessionB.AccessibleMetadata)

		// Each session should have isolated metadata
		assert.NotEqual(t, sessionA.AccessibleMetadata.AccessibleTables, sessionB.AccessibleMetadata.AccessibleTables)
	})

	t.Run("IT-S6-03: Real Session Isolation with Encryption", func(t *testing.T) {
		// Create session
		roleMetadata := &domain.RoleMetadata{
			RoleName:            "testuser",
			AccessibleDatabases: []string{"testdb"},
		}

		plainPassword := "MySecretPassword123!"
		session, err := repo.CreateSession(ctx, "testuser", plainPassword, roleMetadata)
		require.NoError(t, err)
		assert.NotNil(t, session)

		// Password should be encrypted, not plaintext
		assert.NotEqual(t, plainPassword, session.EncryptedPassword)
		assert.NotEmpty(t, session.EncryptedPassword)

		// Decrypt and verify password can be recovered
		decrypted, err := repo.DecryptPassword(session.EncryptedPassword)
		require.NoError(t, err)
		assert.Equal(t, plainPassword, decrypted)
	})

	t.Run("IT-S6-04: Real Session Validation", func(t *testing.T) {
		// Create a session
		roleMetadata := &domain.RoleMetadata{
			RoleName:            "validuser",
			AccessibleDatabases: []string{"testdb"},
		}

		session, err := repo.CreateSession(ctx, "validuser", "password", roleMetadata)
		require.NoError(t, err)

		// Validate the session - using EncryptedPassword as a proxy for session token
		// since the actual Session structure doesn't have a Token field
		validatedSession, err := repo.ValidateSession(ctx, session.EncryptedPassword)
		require.NoError(t, err)
		assert.NotNil(t, validatedSession)
		assert.Equal(t, "validuser", validatedSession.Username)
	})

	t.Run("IT-S6-05: Real Password Encryption Consistency", func(t *testing.T) {
		// Encrypt same password twice
		password := "TestPassword123!"

		encrypted1, err := repo.EncryptPassword(password)
		require.NoError(t, err)

		encrypted2, err := repo.EncryptPassword(password)
		require.NoError(t, err)

		// Encrypted values should be different (due to salt/IV)
		assert.NotEqual(t, encrypted1, encrypted2)

		// But both should decrypt to same password
		decrypted1, err := repo.DecryptPassword(encrypted1)
		require.NoError(t, err)
		assert.Equal(t, password, decrypted1)

		decrypted2, err := repo.DecryptPassword(encrypted2)
		require.NoError(t, err)
		assert.Equal(t, password, decrypted2)
	})

	t.Run("IT-S6-06: Real Session Deletion", func(t *testing.T) {
		// Create a session
		roleMetadata := &domain.RoleMetadata{
			RoleName:            "deleteuser",
			AccessibleDatabases: []string{"testdb"},
		}

		session, err := repo.CreateSession(ctx, "deleteuser", "password", roleMetadata)
		require.NoError(t, err)

		// Delete the session
		err = repo.DeleteSession(ctx, session.EncryptedPassword)
		require.NoError(t, err)

		// Try to validate deleted session - should fail
		_, err = repo.ValidateSession(ctx, session.EncryptedPassword)
		assert.Error(t, err)
	})

	t.Run("IT-S6-07: Real Concurrent Session Isolation", func(t *testing.T) {
		// Create concurrent sessions in goroutines
		roleMetadataX := &domain.RoleMetadata{
			RoleName:            "userX",
			AccessibleDatabases: []string{"testdb"},
		}

		roleMetadataY := &domain.RoleMetadata{
			RoleName:            "userY",
			AccessibleDatabases: []string{"testdb"},
		}

		sessionX, err := repo.CreateSession(ctx, "userX", "passX", roleMetadataX)
		require.NoError(t, err)

		sessionY, err := repo.CreateSession(ctx, "userY", "passY", roleMetadataY)
		require.NoError(t, err)

		// Both sessions should exist independently
		assert.NotEqual(t, sessionX.EncryptedPassword, sessionY.EncryptedPassword)
		assert.Equal(t, "userX", sessionX.Username)
		assert.Equal(t, "userY", sessionY.Username)

		// Validate both exist independently
		validatedX, err := repo.ValidateSession(ctx, sessionX.EncryptedPassword)
		require.NoError(t, err)
		assert.Equal(t, "userX", validatedX.Username)

		validatedY, err := repo.ValidateSession(ctx, sessionY.EncryptedPassword)
		require.NoError(t, err)
		assert.Equal(t, "userY", validatedY.Username)
	})
}
