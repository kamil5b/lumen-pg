package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
)

// SessionRepositoryConstructor is a function type that creates a SessionRepository
type SessionRepositoryConstructor func(db *sql.DB) repository.SessionRepository

// SessionRepositoryRunner runs all session repository tests against an implementation
// Maps to TEST_PLAN.md:
// - Story 2: Authentication & Identity [UC-S2-06~10, UC-S2-12, IT-S2-04~05]
// - Story 6: Isolation [UC-S6-01, UC-S6-03, IT-S6-01~03]
// - Story 7: Security & Best Practices [UC-S7-06~07, IT-S7-03, E2E-S7-04]
func SessionRepositoryRunner(t *testing.T, constructor SessionRepositoryConstructor) {
	t.Helper()

	ctx := context.Background()

	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
	)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	connStr, err := container.ConnectionString(ctx)
	require.NoError(t, err)

	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer db.Close()

	err = db.PingContext(ctx)
	require.NoError(t, err)

	repo := constructor(db)

	// UC-S2-06: Session Cookie Creation - Username
	// UC-S2-08: Session Validation - Valid Session
	// IT-S2-04: Session Persistence After Probe
	t.Run("CreateSession and GetSession", func(t *testing.T) {
		now := time.Now()
		session := &domain.Session{
			ID:        "session_123",
			Username:  "testuser",
			CreatedAt: now,
			ExpiresAt: now.Add(24 * time.Hour),
		}

		err := repo.CreateSession(ctx, session)
		require.NoError(t, err)

		retrieved, err := repo.GetSession(ctx, "session_123")
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		require.Equal(t, session.ID, retrieved.ID)
		require.Equal(t, session.Username, retrieved.Username)
	})

	// UC-S2-09: Session Validation - Expired Session
	t.Run("GetSession returns error for non-existent session", func(t *testing.T) {
		_, err := repo.GetSession(ctx, "nonexistent_session")
		require.Error(t, err)
	})

	// UC-S2-10: Session Re-authentication
	// IT-S2-04: Session Persistence After Probe
	t.Run("UpdateSession modifies existing session", func(t *testing.T) {
		now := time.Now()
		session := &domain.Session{
			ID:        "update_session",
			Username:  "testuser",
			CreatedAt: now,
			ExpiresAt: now.Add(24 * time.Hour),
		}

		err := repo.CreateSession(ctx, session)
		require.NoError(t, err)

		newExpiry := now.Add(48 * time.Hour)
		session.ExpiresAt = newExpiry

		err = repo.UpdateSession(ctx, session)
		require.NoError(t, err)

		retrieved, err := repo.GetSession(ctx, "update_session")
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		require.Equal(t, newExpiry.Unix(), retrieved.ExpiresAt.Unix())
	})

	// UC-S2-12: Logout Cookie Clearing
	// E2E-S2-04: Logout Flow
	t.Run("DeleteSession removes session", func(t *testing.T) {
		now := time.Now()
		session := &domain.Session{
			ID:        "delete_session",
			Username:  "testuser",
			CreatedAt: now,
			ExpiresAt: now.Add(24 * time.Hour),
		}

		err := repo.CreateSession(ctx, session)
		require.NoError(t, err)

		err = repo.DeleteSession(ctx, "delete_session")
		require.NoError(t, err)

		_, err = repo.GetSession(ctx, "delete_session")
		require.Error(t, err)
	})

	// UC-S2-08: Session Validation - Valid Session
	// IT-S2-04: Session Persistence After Probe
	t.Run("ValidateSession returns valid session", func(t *testing.T) {
		now := time.Now()
		session := &domain.Session{
			ID:        "valid_session",
			Username:  "testuser",
			CreatedAt: now,
			ExpiresAt: now.Add(24 * time.Hour),
		}

		err := repo.CreateSession(ctx, session)
		require.NoError(t, err)

		validated, err := repo.ValidateSession(ctx, "valid_session")
		require.NoError(t, err)
		require.NotNil(t, validated)
		require.Equal(t, session.ID, validated.ID)
	})

	// UC-S2-09: Session Validation - Expired Session
	// UC-S7-06: Session Timeout Short-Lived Cookie
	// IT-S7-03: Real Session Expiration
	t.Run("ValidateSession returns error for expired session", func(t *testing.T) {
		now := time.Now()
		session := &domain.Session{
			ID:        "expired_session",
			Username:  "testuser",
			CreatedAt: now.Add(-48 * time.Hour),
			ExpiresAt: now.Add(-1 * time.Hour),
		}

		err := repo.CreateSession(ctx, session)
		require.NoError(t, err)

		_, err = repo.ValidateSession(ctx, "expired_session")
		require.Error(t, err)
	})

	// UC-S2-08: Session Validation - Valid Session
	t.Run("ValidateSession returns error for non-existent session", func(t *testing.T) {
		_, err := repo.ValidateSession(ctx, "nonexistent_session")
		require.Error(t, err)
	})

	// UC-S2-08: Session Validation - Valid Session
	// IT-S2-05: Concurrent User Sessions with Isolated Resources
	// UC-S6-01: Session Isolation
	t.Run("GetSessionByUsername retrieves most recent session", func(t *testing.T) {
		now := time.Now()
		username := "multiuser"

		session1 := &domain.Session{
			ID:        "session_1",
			Username:  username,
			CreatedAt: now.Add(-1 * time.Hour),
			ExpiresAt: now.Add(23 * time.Hour),
		}

		session2 := &domain.Session{
			ID:        "session_2",
			Username:  username,
			CreatedAt: now,
			ExpiresAt: now.Add(24 * time.Hour),
		}

		err := repo.CreateSession(ctx, session1)
		require.NoError(t, err)

		err = repo.CreateSession(ctx, session2)
		require.NoError(t, err)

		retrieved, err := repo.GetSessionByUsername(ctx, username)
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		require.Equal(t, "session_2", retrieved.ID)
	})

	// UC-S2-08: Session Validation - Valid Session
	t.Run("GetSessionByUsername returns error for non-existent user", func(t *testing.T) {
		_, err := repo.GetSessionByUsername(ctx, "nonexistent_user")
		require.Error(t, err)
	})

	// UC-S2-12: Logout Cookie Clearing
	// UC-S6-03: Cookie Isolation
	// E2E-S6-03: One User Cannot See Another's Session
	t.Run("InvalidateUserSessions removes all user sessions", func(t *testing.T) {
		now := time.Now()
		username := "invalidate_user"

		session1 := &domain.Session{
			ID:        "inv_session_1",
			Username:  username,
			CreatedAt: now,
			ExpiresAt: now.Add(24 * time.Hour),
		}

		session2 := &domain.Session{
			ID:        "inv_session_2",
			Username:  username,
			CreatedAt: now,
			ExpiresAt: now.Add(24 * time.Hour),
		}

		err := repo.CreateSession(ctx, session1)
		require.NoError(t, err)

		err = repo.CreateSession(ctx, session2)
		require.NoError(t, err)

		err = repo.InvalidateUserSessions(ctx, username)
		require.NoError(t, err)

		_, err = repo.GetSessionByUsername(ctx, username)
		require.Error(t, err)
	})

	// UC-S2-09: Session Validation - Expired Session
	// UC-S7-06: Session Timeout Short-Lived Cookie
	// IT-S7-03: Real Session Expiration
	t.Run("InvalidateExpiredSessions removes expired sessions", func(t *testing.T) {
		now := time.Now()

		expiredSession := &domain.Session{
			ID:        "old_session",
			Username:  "old_user",
			CreatedAt: now.Add(-48 * time.Hour),
			ExpiresAt: now.Add(-1 * time.Hour),
		}

		validSession := &domain.Session{
			ID:        "current_session",
			Username:  "current_user",
			CreatedAt: now,
			ExpiresAt: now.Add(24 * time.Hour),
		}

		err := repo.CreateSession(ctx, expiredSession)
		require.NoError(t, err)

		err = repo.CreateSession(ctx, validSession)
		require.NoError(t, err)

		err = repo.InvalidateExpiredSessions(ctx)
		require.NoError(t, err)

		_, err = repo.GetSession(ctx, "old_session")
		require.Error(t, err)

		retrieved, err := repo.GetSession(ctx, "current_session")
		require.NoError(t, err)
		require.NotNil(t, retrieved)
	})

	// UC-S2-08: Session Validation - Valid Session
	// IT-S2-04: Session Persistence After Probe
	t.Run("SessionExists returns true for existing session", func(t *testing.T) {
		now := time.Now()
		session := &domain.Session{
			ID:        "exists_session",
			Username:  "testuser",
			CreatedAt: now,
			ExpiresAt: now.Add(24 * time.Hour),
		}

		err := repo.CreateSession(ctx, session)
		require.NoError(t, err)

		exists, err := repo.SessionExists(ctx, "exists_session")
		require.NoError(t, err)
		require.True(t, exists)
	})

	// UC-S2-08: Session Validation - Valid Session
	t.Run("SessionExists returns false for non-existent session", func(t *testing.T) {
		exists, err := repo.SessionExists(ctx, "nonexistent_session")
		require.NoError(t, err)
		require.False(t, exists)
	})

	// UC-S2-09: Session Validation - Expired Session
	// UC-S7-07: Session Timeout Long-Lived Cookie
	// E2E-S7-04: Session Timeout Enforcement
	t.Run("SessionExists returns false for expired session", func(t *testing.T) {
		now := time.Now()
		session := &domain.Session{
			ID:        "expired_check_session",
			Username:  "testuser",
			CreatedAt: now.Add(-48 * time.Hour),
			ExpiresAt: now.Add(-1 * time.Hour),
		}

		err := repo.CreateSession(ctx, session)
		require.NoError(t, err)

		exists, err := repo.SessionExists(ctx, "expired_check_session")
		require.NoError(t, err)
		require.False(t, exists)
	})

	// IT-S2-05: Concurrent User Sessions with Isolated Resources
	// UC-S6-01: Session Isolation
	// E2E-S6-03: One User Cannot See Another's Session
	t.Run("Multiple sessions for different users", func(t *testing.T) {
		now := time.Now()

		user1Session := &domain.Session{
			ID:        "user1_session",
			Username:  "user1",
			CreatedAt: now,
			ExpiresAt: now.Add(24 * time.Hour),
		}

		user2Session := &domain.Session{
			ID:        "user2_session",
			Username:  "user2",
			CreatedAt: now,
			ExpiresAt: now.Add(24 * time.Hour),
		}

		err := repo.CreateSession(ctx, user1Session)
		require.NoError(t, err)

		err = repo.CreateSession(ctx, user2Session)
		require.NoError(t, err)

		session1, err := repo.GetSessionByUsername(ctx, "user1")
		require.NoError(t, err)
		require.Equal(t, "user1", session1.Username)

		session2, err := repo.GetSessionByUsername(ctx, "user2")
		require.NoError(t, err)
		require.Equal(t, "user2", session2.Username)
	})

	// UC-S2-06: Session Cookie Creation - Username
	// UC-S2-07: Session Cookie Creation - Password
	t.Run("CreateSession with duplicate ID overwrites previous", func(t *testing.T) {
		now := time.Now()

		session1 := &domain.Session{
			ID:        "duplicate_id",
			Username:  "user1",
			CreatedAt: now,
			ExpiresAt: now.Add(24 * time.Hour),
		}

		err := repo.CreateSession(ctx, session1)
		require.NoError(t, err)

		session2 := &domain.Session{
			ID:        "duplicate_id",
			Username:  "user2",
			CreatedAt: now,
			ExpiresAt: now.Add(48 * time.Hour),
		}

		err = repo.CreateSession(ctx, session2)
		require.NoError(t, err)

		retrieved, err := repo.GetSession(ctx, "duplicate_id")
		require.NoError(t, err)
		require.Equal(t, "user2", retrieved.Username)
	})
}
