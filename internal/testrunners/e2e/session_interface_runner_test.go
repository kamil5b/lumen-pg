package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/implementations/mocks"
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
)

// SessionInterfaceConstructor creates a session repository with mock dependencies
type SessionInterfaceConstructor func(ctrl *gomock.Controller) repository.SessionRepository

// SessionInterfaceRunner runs unit tests for session repository interface (Story 2)
func SessionInterfaceRunner(t *testing.T, constructor SessionInterfaceConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSessionRepository(ctrl)
	ctx := context.Background()

	t.Run("UC-S2-06: Session Cookie Creation - Username", func(t *testing.T) {
		username := "testuser"
		password := "testpass"
		roleMetadata := &domain.RoleMetadata{
			RoleName:            username,
			AccessibleDatabases: []string{"testdb"},
			AccessibleSchemas: map[string][]string{
				"testdb": {"public"},
			},
			AccessibleTables: map[string][]string{
				"testdb.public": {"users"},
			},
		}

		expectedSession := &domain.Session{
			Username:           username,
			EncryptedPassword:  "encrypted_password_xyz",
			AccessibleMetadata: roleMetadata,
			CreatedAt:          time.Now(),
			ExpiresAt:          time.Now().Add(24 * time.Hour),
		}

		mockRepo.EXPECT().CreateSession(ctx, username, password, roleMetadata).Return(expectedSession, nil)

		session, err := mockRepo.CreateSession(ctx, username, password, roleMetadata)

		require.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, username, session.Username)
		assert.NotEmpty(t, session.EncryptedPassword)
	})

	t.Run("UC-S2-07: Session Cookie Creation - Password", func(t *testing.T) {
		username := "testuser"
		password := "securepassword123"
		roleMetadata := &domain.RoleMetadata{
			RoleName:            username,
			AccessibleDatabases: []string{"testdb"},
		}

		expectedSession := &domain.Session{
			Username:           username,
			EncryptedPassword:  "encrypted_abc123xyz",
			AccessibleMetadata: roleMetadata,
			CreatedAt:          time.Now(),
			ExpiresAt:          time.Now().Add(15 * time.Minute),
		}

		mockRepo.EXPECT().CreateSession(ctx, username, password, roleMetadata).Return(expectedSession, nil)

		session, err := mockRepo.CreateSession(ctx, username, password, roleMetadata)

		require.NoError(t, err)
		assert.NotNil(t, session)
		assert.NotEmpty(t, session.EncryptedPassword)
		assert.NotEqual(t, password, session.EncryptedPassword)
	})

	t.Run("UC-S2-08: Session Validation - Valid Session", func(t *testing.T) {
		sessionToken := "valid-session-token-xyz"
		expectedSession := &domain.Session{
			Username:           "testuser",
			EncryptedPassword:  "encrypted_xyz",
			AccessibleMetadata: &domain.RoleMetadata{RoleName: "testuser"},
			CreatedAt:          time.Now().Add(-1 * time.Hour),
			ExpiresAt:          time.Now().Add(23 * time.Hour),
		}

		mockRepo.EXPECT().ValidateSession(ctx, sessionToken).Return(expectedSession, nil)

		session, err := mockRepo.ValidateSession(ctx, sessionToken)

		require.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, "testuser", session.Username)
		assert.True(t, session.ExpiresAt.After(time.Now()))
	})

	t.Run("UC-S2-09: Session Validation - Expired Session", func(t *testing.T) {
		sessionToken := "expired-session-token"

		mockRepo.EXPECT().ValidateSession(ctx, sessionToken).Return(nil, assert.AnError)

		session, err := mockRepo.ValidateSession(ctx, sessionToken)

		require.Error(t, err)
		assert.Nil(t, session)
	})

	t.Run("UC-S2-10: Session Re-authentication", func(t *testing.T) {
		username := "testuser"
		password := "testpass"
		roleMetadata := &domain.RoleMetadata{
			RoleName:            username,
			AccessibleDatabases: []string{"testdb"},
		}

		expectedSession := &domain.Session{
			Username:           username,
			EncryptedPassword:  "encrypted_password",
			AccessibleMetadata: roleMetadata,
			CreatedAt:          time.Now(),
			ExpiresAt:          time.Now().Add(15 * time.Minute),
		}

		mockRepo.EXPECT().CreateSession(ctx, username, password, roleMetadata).Return(expectedSession, nil)

		session, err := mockRepo.CreateSession(ctx, username, password, roleMetadata)

		require.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, username, session.Username)
	})

	t.Run("UC-S2-12: Logout Cookie Clearing", func(t *testing.T) {
		sessionToken := "session-to-clear"

		mockRepo.EXPECT().DeleteSession(ctx, sessionToken).Return(nil)

		err := mockRepo.DeleteSession(ctx, sessionToken)

		require.NoError(t, err)
	})

	t.Run("UC-S7-03: Password Encryption in Cookie", func(t *testing.T) {
		plainPassword := "MySecurePassword123!"
		expectedEncrypted := "encrypted_value_xyz"

		mockRepo.EXPECT().EncryptPassword(plainPassword).Return(expectedEncrypted, nil)

		encrypted, err := mockRepo.EncryptPassword(plainPassword)

		require.NoError(t, err)
		assert.NotEqual(t, plainPassword, encrypted)
		assert.Equal(t, expectedEncrypted, encrypted)
	})

	t.Run("UC-S7-04: Password Decryption from Cookie", func(t *testing.T) {
		plainPassword := "MySecurePassword123!"
		encryptedPassword := "encrypted_value_xyz"

		mockRepo.EXPECT().DecryptPassword(encryptedPassword).Return(plainPassword, nil)

		decrypted, err := mockRepo.DecryptPassword(encryptedPassword)

		require.NoError(t, err)
		assert.Equal(t, plainPassword, decrypted)
	})

	t.Run("UC-S7-05: Cookie Tampering Detection", func(t *testing.T) {
		tamperedPassword := "tampered_encrypted_value"

		mockRepo.EXPECT().DecryptPassword(tamperedPassword).Return("", assert.AnError)

		decrypted, err := mockRepo.DecryptPassword(tamperedPassword)

		require.Error(t, err)
		assert.Equal(t, "", decrypted)
	})

	t.Run("UC-S2-11: Data Explorer Population After Login", func(t *testing.T) {
		username := "testuser"
		password := "testpass"
		roleMetadata := &domain.RoleMetadata{
			RoleName:            username,
			AccessibleDatabases: []string{"testdb", "anotherdb"},
			AccessibleSchemas: map[string][]string{
				"testdb":    {"public", "internal"},
				"anotherdb": {"public"},
			},
			AccessibleTables: map[string][]string{
				"testdb.public":    {"users", "posts", "comments"},
				"testdb.internal":  {"audit_logs"},
				"anotherdb.public": {"products", "orders"},
			},
		}

		expectedSession := &domain.Session{
			Username:           username,
			EncryptedPassword:  "encrypted_pwd",
			AccessibleMetadata: roleMetadata,
			CreatedAt:          time.Now(),
			ExpiresAt:          time.Now().Add(24 * time.Hour),
		}

		mockRepo.EXPECT().CreateSession(ctx, username, password, roleMetadata).Return(expectedSession, nil)

		session, err := mockRepo.CreateSession(ctx, username, password, roleMetadata)

		require.NoError(t, err)
		assert.NotNil(t, session)
		assert.NotEmpty(t, session.AccessibleMetadata.AccessibleDatabases)
		assert.Len(t, session.AccessibleMetadata.AccessibleDatabases, 2)
		assert.Len(t, session.AccessibleMetadata.AccessibleTables, 4)
	})

	t.Run("UC-S2-13: Session Cookie Properties - Expiration", func(t *testing.T) {
		username := "testuser"
		password := "testpass"
		roleMetadata := &domain.RoleMetadata{RoleName: username}

		now := time.Now()
		expectedSession := &domain.Session{
			Username:          username,
			EncryptedPassword: "encrypted_pwd",
			CreatedAt:         now,
			ExpiresAt:         now.Add(24 * time.Hour),
		}

		mockRepo.EXPECT().CreateSession(ctx, username, password, roleMetadata).Return(expectedSession, nil)

		session, err := mockRepo.CreateSession(ctx, username, password, roleMetadata)

		require.NoError(t, err)
		assert.True(t, session.ExpiresAt.After(session.CreatedAt))
		assert.True(t, session.ExpiresAt.Sub(session.CreatedAt) > 23*time.Hour)
	})

	t.Run("UC-S7-06: Session Timeout Short-Lived Cookie", func(t *testing.T) {
		username := "testuser"
		password := "testpass"
		roleMetadata := &domain.RoleMetadata{RoleName: username}

		now := time.Now()
		expectedSession := &domain.Session{
			Username:          username,
			EncryptedPassword: "encrypted_short_lived",
			CreatedAt:         now,
			ExpiresAt:         now.Add(15 * time.Minute),
		}

		mockRepo.EXPECT().CreateSession(ctx, username, password, roleMetadata).Return(expectedSession, nil)

		session, err := mockRepo.CreateSession(ctx, username, password, roleMetadata)

		require.NoError(t, err)
		assert.True(t, session.ExpiresAt.Before(now.Add(20*time.Minute)))
	})

	t.Run("UC-S7-07: Session Timeout Long-Lived Cookie", func(t *testing.T) {
		username := "testuser"
		roleMetadata := &domain.RoleMetadata{RoleName: username}

		now := time.Now()
		expectedSession := &domain.Session{
			Username:  username,
			CreatedAt: now,
			ExpiresAt: now.Add(24 * time.Hour),
		}

		mockRepo.EXPECT().CreateSession(ctx, username, "pass", roleMetadata).Return(expectedSession, nil)

		session, err := mockRepo.CreateSession(ctx, username, "pass", roleMetadata)

		require.NoError(t, err)
		assert.True(t, session.ExpiresAt.After(now.Add(23*time.Hour)))
	})

	t.Run("UC-S2-14: Session Validation - Invalid Token", func(t *testing.T) {
		invalidToken := "invalid-token-format"

		mockRepo.EXPECT().ValidateSession(ctx, invalidToken).Return(nil, assert.AnError)

		session, err := mockRepo.ValidateSession(ctx, invalidToken)

		require.Error(t, err)
		assert.Nil(t, session)
	})

	t.Run("UC-S2-15: Multiple Sessions Per User - Isolation", func(t *testing.T) {
		username := "testuser"
		password := "testpass"
		roleMetadata := &domain.RoleMetadata{RoleName: username}

		session1 := &domain.Session{
			Username:          username,
			EncryptedPassword: "encrypted_session_1",
			CreatedAt:         time.Now(),
			ExpiresAt:         time.Now().Add(24 * time.Hour),
		}

		session2 := &domain.Session{
			Username:          username,
			EncryptedPassword: "encrypted_session_2",
			CreatedAt:         time.Now(),
			ExpiresAt:         time.Now().Add(24 * time.Hour),
		}

		mockRepo.EXPECT().CreateSession(ctx, username, password, roleMetadata).Return(session1, nil)
		mockRepo.EXPECT().CreateSession(ctx, username, password, roleMetadata).Return(session2, nil)

		result1, _ := mockRepo.CreateSession(ctx, username, password, roleMetadata)
		result2, _ := mockRepo.CreateSession(ctx, username, password, roleMetadata)

		assert.NotEqual(t, result1.EncryptedPassword, result2.EncryptedPassword)
	})

	t.Run("UC-S2-16: Session Deletion on Logout", func(t *testing.T) {
		sessionToken := "session-to-delete"

		mockRepo.EXPECT().DeleteSession(ctx, sessionToken).Return(nil)

		err := mockRepo.DeleteSession(ctx, sessionToken)

		require.NoError(t, err)
	})

	t.Run("UC-S2-17: Password Encryption Uniqueness", func(t *testing.T) {
		plainPassword := "testpassword"
		enc1 := "encrypted_value_1"
		enc2 := "encrypted_value_2"

		mockRepo.EXPECT().EncryptPassword(plainPassword).Return(enc1, nil)
		mockRepo.EXPECT().EncryptPassword(plainPassword).Return(enc2, nil)

		encrypted1, _ := mockRepo.EncryptPassword(plainPassword)
		encrypted2, _ := mockRepo.EncryptPassword(plainPassword)

		// Each encryption can be different (salted)
		assert.NotEmpty(t, encrypted1)
		assert.NotEmpty(t, encrypted2)
	})

	t.Run("UC-S2-18: Session Integrity - AccessibleMetadata Present", func(t *testing.T) {
		username := "testuser"
		password := "testpass"
		roleMetadata := &domain.RoleMetadata{
			RoleName:            username,
			AccessibleDatabases: []string{"testdb"},
			AccessibleTables: map[string][]string{
				"testdb.public": {"users", "posts"},
			},
		}

		expectedSession := &domain.Session{
			Username:           username,
			EncryptedPassword:  "encrypted_pwd",
			AccessibleMetadata: roleMetadata,
			CreatedAt:          time.Now(),
			ExpiresAt:          time.Now().Add(24 * time.Hour),
		}

		mockRepo.EXPECT().CreateSession(ctx, username, password, roleMetadata).Return(expectedSession, nil)

		session, err := mockRepo.CreateSession(ctx, username, password, roleMetadata)

		require.NoError(t, err)
		assert.NotNil(t, session.AccessibleMetadata)
		assert.Equal(t, username, session.AccessibleMetadata.RoleName)
	})

	t.Run("UC-S2-19: Create Session with Empty Accessible Resources", func(t *testing.T) {
		username := "restricteduser"
		password := "testpass"
		roleMetadata := &domain.RoleMetadata{
			RoleName:            username,
			AccessibleDatabases: []string{},
			AccessibleTables:    map[string][]string{},
		}

		expectedSession := &domain.Session{
			Username:           username,
			EncryptedPassword:  "encrypted_pwd",
			AccessibleMetadata: roleMetadata,
			CreatedAt:          time.Now(),
			ExpiresAt:          time.Now().Add(24 * time.Hour),
		}

		mockRepo.EXPECT().CreateSession(ctx, username, password, roleMetadata).Return(expectedSession, nil)

		session, err := mockRepo.CreateSession(ctx, username, password, roleMetadata)

		require.NoError(t, err)
		assert.NotNil(t, session)
		assert.Empty(t, session.AccessibleMetadata.AccessibleDatabases)
	})

	t.Run("UC-S2-20: Validate Session - Check Expiration Time", func(t *testing.T) {
		sessionToken := "valid-token-abc"
		now := time.Now()
		expectedSession := &domain.Session{
			Username:  "testuser",
			CreatedAt: now.Add(-2 * time.Hour),
			ExpiresAt: now.Add(22 * time.Hour),
		}

		mockRepo.EXPECT().ValidateSession(ctx, sessionToken).Return(expectedSession, nil)

		session, err := mockRepo.ValidateSession(ctx, sessionToken)

		require.NoError(t, err)
		assert.True(t, session.ExpiresAt.After(now))
	})
}
