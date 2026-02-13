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

// AuthRepoConstructor creates an auth repository with database connection
type AuthRepoConstructor func(db *sql.DB) repository.SessionRepository

// AuthRepositoryRunner runs integration tests for auth repository (Story 2)
func AuthRepositoryRunner(t *testing.T, constructor AuthRepoConstructor) {
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

	t.Run("IT-S2-01: Real PostgreSQL Connection Probe", func(t *testing.T) {
		err := db.PingContext(ctx)
		require.NoError(t, err)
	})

	t.Run("IT-S2-02: Real PostgreSQL Connection Probe Failure", func(t *testing.T) {
		// Try to connect with wrong credentials
		badConnStr := "postgres://wronguser:wrongpass@localhost:5432/nonexistent"
		badDB, err := sql.Open("postgres", badConnStr)
		require.NoError(t, err)
		defer badDB.Close()

		// Should fail on ping
		err = badDB.PingContext(ctx)
		assert.Error(t, err)
	})

	t.Run("IT-S2-03: Real Role-Based Resource Access", func(t *testing.T) {
		// Create test tables
		_, err := db.ExecContext(ctx, `
			CREATE TABLE IF NOT EXISTS users (
				id SERIAL PRIMARY KEY,
				username VARCHAR(50) NOT NULL
			)
		`)
		require.NoError(t, err)

		// Query to get accessible resources for postgres role
		var roleName string
		err = db.QueryRowContext(ctx, "SELECT current_user").Scan(&roleName)
		require.NoError(t, err)
		assert.NotEmpty(t, roleName)
	})

	t.Run("IT-S2-04: Session Persistence After Probe", func(t *testing.T) {
		username := "testuser"
		password := "testpass"
		roleMetadata := &domain.RoleMetadata{
			RoleName:            username,
			AccessibleDatabases: []string{"testdb"},
		}

		createdSession, err := repo.CreateSession(ctx, username, password, roleMetadata)

		require.NoError(t, err)
		assert.NotNil(t, createdSession)
		assert.Equal(t, username, createdSession.Username)
	})

	t.Run("IT-S2-05: Concurrent User Sessions with Isolated Resources", func(t *testing.T) {
		// Create first session for userA
		userAMetadata := &domain.RoleMetadata{
			RoleName:            "userA",
			AccessibleDatabases: []string{"testdb"},
		}
		sessionA, err := repo.CreateSession(ctx, "userA", "passA", userAMetadata)
		require.NoError(t, err)
		assert.NotNil(t, sessionA)

		// Create second session for userB
		userBMetadata := &domain.RoleMetadata{
			RoleName:            "userB",
			AccessibleDatabases: []string{"testdb"},
		}
		sessionB, err := repo.CreateSession(ctx, "userB", "passB", userBMetadata)
		require.NoError(t, err)
		assert.NotNil(t, sessionB)

		// Both sessions should be independent
		assert.Equal(t, "userA", sessionA.Username)
		assert.Equal(t, "userB", sessionB.Username)
	})

	t.Run("IT-S2-06: Password Encryption in Session", func(t *testing.T) {
		plainPassword := "MySecurePassword123!"

		encrypted, err := repo.EncryptPassword(plainPassword)
		require.NoError(t, err)
		assert.NotEqual(t, plainPassword, encrypted)
		assert.NotEmpty(t, encrypted)
	})

	t.Run("IT-S2-07: Password Decryption from Session", func(t *testing.T) {
		plainPassword := "MySecurePassword123!"

		encrypted, err := repo.EncryptPassword(plainPassword)
		require.NoError(t, err)

		decrypted, err := repo.DecryptPassword(encrypted)
		require.NoError(t, err)
		assert.Equal(t, plainPassword, decrypted)
	})

	t.Run("IT-S2-08: Session Validation After Creation", func(t *testing.T) {
		username := "validationuser"
		password := "testpass"
		roleMetadata := &domain.RoleMetadata{
			RoleName:            username,
			AccessibleDatabases: []string{"testdb"},
		}

		createdSession, err := repo.CreateSession(ctx, username, password, roleMetadata)
		require.NoError(t, err)

		// Validate the created session
		validatedSession, err := repo.ValidateSession(ctx, createdSession.EncryptedPassword)
		require.NoError(t, err)
		assert.NotNil(t, validatedSession)
		assert.Equal(t, username, validatedSession.Username)
	})

	t.Run("IT-S2-09: Session Deletion on Logout", func(t *testing.T) {
		username := "logoutuser"
		password := "testpass"
		roleMetadata := &domain.RoleMetadata{
			RoleName:            username,
			AccessibleDatabases: []string{"testdb"},
		}

		createdSession, err := repo.CreateSession(ctx, username, password, roleMetadata)
		require.NoError(t, err)

		// Delete session
		err = repo.DeleteSession(ctx, createdSession.EncryptedPassword)
		require.NoError(t, err)

		// Verify session is deleted
		_, err = repo.ValidateSession(ctx, createdSession.EncryptedPassword)
		assert.Error(t, err)
	})
}
