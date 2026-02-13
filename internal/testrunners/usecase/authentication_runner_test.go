package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
	mockRepository "github.com/kamil5b/lumen-pg/internal/testrunners/mocks/repository"
	"github.com/stretchr/testify/require"
)

// AuthenticationUsecaseConstructor is a function type that creates an AuthenticationUseCase
type AuthenticationUsecaseConstructor func(
	databaseRepo repository.DatabaseRepository,
	metadataRepo repository.MetadataRepository,
	sessionRepo repository.SessionRepository,
	rbacRepo repository.RBACRepository,
	encryptionRepo repository.EncryptionRepository,
) usecase.AuthenticationUseCase

// AuthenticationUsecaseRunner runs all authentication usecase tests against an implementation
// Maps to TEST_PLAN.md:
// - Story 2: Authentication & Identity [UC-S2-01~15, IT-S2-01~05, E2E-S2-01~06]
func AuthenticationUsecaseRunner(t *testing.T, constructor AuthenticationUsecaseConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockDatabase := mockRepository.NewMockDatabaseRepository(ctrl)
	mockMetadata := mockRepository.NewMockMetadataRepository(ctrl)
	mockSession := mockRepository.NewMockSessionRepository(ctrl)
	mockRBAC := mockRepository.NewMockRBACRepository(ctrl)
	mockEncryption := mockRepository.NewMockEncryptionRepository(ctrl)

	uc := constructor(mockDatabase, mockMetadata, mockSession, mockRBAC, mockEncryption)

	// UC-S2-01: Login Form Validation - Empty Username
	t.Run("ValidateLoginForm rejects empty username", func(t *testing.T) {
		errors, err := uc.ValidateLoginForm(ctx, domain.LoginRequest{
			Username: "",
			Password: "password123",
		})

		require.NoError(t, err)
		require.NotNil(t, errors)
		require.Greater(t, len(errors), 0)
	})

	// UC-S2-02: Login Form Validation - Empty Password
	t.Run("ValidateLoginForm rejects empty password", func(t *testing.T) {
		errors, err := uc.ValidateLoginForm(ctx, domain.LoginRequest{
			Username: "testuser",
			Password: "",
		})

		require.NoError(t, err)
		require.NotNil(t, errors)
		require.Greater(t, len(errors), 0)
	})

	t.Run("ValidateLoginForm accepts valid credentials", func(t *testing.T) {
		errors, err := uc.ValidateLoginForm(ctx, domain.LoginRequest{
			Username: "testuser",
			Password: "password123",
		})

		require.NoError(t, err)
		require.NotNil(t, errors)
		require.Equal(t, 0, len(errors))
	})

	// UC-S2-03: Login Connection Probe
	// IT-S2-01: Real PostgreSQL Connection Probe
	// E2E-S2-01: Login Flow with Connection Probe
	t.Run("ProbeConnection succeeds with valid credentials", func(t *testing.T) {
		mockDatabase.EXPECT().
			TestConnection(gomock.Any(), gomock.Any()).
			Return(nil)

		success, err := uc.ProbeConnection(ctx, "testuser", "password123")

		require.NoError(t, err)
		require.True(t, success)
	})

	// UC-S2-04: Login Connection Probe Failure
	// IT-S2-02: Real PostgreSQL Connection Probe Failure
	// E2E-S2-03: Login Flow - Invalid Credentials
	t.Run("ProbeConnection fails with invalid credentials", func(t *testing.T) {
		mockDatabase.EXPECT().
			TestConnection(gomock.Any(), gomock.Any()).
			Return(ErrInvalidCredentials)

		success, err := uc.ProbeConnection(ctx, "testuser", "wrongpassword")

		require.Error(t, err)
		require.False(t, success)
	})

	// UC-S2-05: Login Success After Probe
	// IT-S2-03: Real Role-Based Resource Access
	t.Run("Login returns LoginResponse with success flag", func(t *testing.T) {
		mockDatabase.EXPECT().
			TestConnection(gomock.Any(), gomock.Any()).
			Return(nil)

		mockRBAC.EXPECT().
			GetUserAccessibleDatabases(gomock.Any(), "testuser").
			Return([]string{"testdb"}, nil)

		response, err := uc.Login(ctx, domain.LoginRequest{
			Username: "testuser",
			Password: "password123",
		})

		require.NoError(t, err)
		require.NotNil(t, response)
		require.True(t, response.Success)
		require.Equal(t, "testuser", response.Username)
	})

	// UC-S2-06: Session Cookie Creation - Username
	// UC-S2-07: Session Cookie Creation - Password
	// IT-S2-04: Session Persistence After Probe
	// E2E-S2-01: Login Flow with Connection Probe
	t.Run("CreateSession creates new session with encrypted password", func(t *testing.T) {
		mockSession.EXPECT().
			StoreSession(gomock.Any(), gomock.Any()).
			Return(nil)

		mockEncryption.EXPECT().
			EncryptPassword(gomock.Any(), "password123").
			Return("encrypted_password", nil)

		mockRBAC.EXPECT().
			GetUserAccessibleDatabases(gomock.Any(), "testuser").
			Return([]string{"testdb"}, nil)

		mockRBAC.EXPECT().
			GetUserAccessibleSchemas(gomock.Any(), "testuser", "testdb").
			Return([]string{"public"}, nil)

		mockRBAC.EXPECT().
			GetUserAccessibleTables(gomock.Any(), "testuser", "testdb", "public").
			Return([]domain.AccessibleTable{
				{Database: "testdb", Schema: "public", Name: "users", HasSelect: true},
			}, nil)

		session, err := uc.CreateSession(ctx, "testuser", "password123", "testdb", "public", "users")

		require.NoError(t, err)
		require.NotNil(t, session)
		require.Equal(t, "testuser", session.Username)
		require.NotEmpty(t, session.ID)
	})

	// UC-S2-11: Data Explorer Population After Login
	t.Run("GetUserAccessibleResources returns role metadata", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetRoleMetadata(gomock.Any(), "testuser").
			Return(&domain.RoleMetadata{
				Name:                "testuser",
				AccessibleDatabases: []string{"testdb1", "testdb2"},
				AccessibleSchemas:   []string{"public", "private"},
				AccessibleTables: []domain.AccessibleTable{
					{Database: "testdb1", Schema: "public", Name: "users", HasSelect: true},
					{Database: "testdb1", Schema: "public", Name: "posts", HasSelect: true},
				},
			}, nil)

		resources, err := uc.GetUserAccessibleResources(ctx, "testuser")

		require.NoError(t, err)
		require.NotNil(t, resources)
		require.Equal(t, 2, len(resources.AccessibleDatabases))
	})

	// UC-S2-11: Data Explorer Population After Login
	// E2E-S2-06: Data Explorer Populated After Login
	t.Run("GetFirstAccessibleDatabase returns first database", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetRoleMetadata(gomock.Any(), "testuser").
			Return(&domain.RoleMetadata{
				Name:                "testuser",
				AccessibleDatabases: []string{"testdb1", "testdb2"},
			}, nil)

		database, err := uc.GetFirstAccessibleDatabase(ctx, "testuser")

		require.NoError(t, err)
		require.NotEmpty(t, database)
		require.Equal(t, "testdb1", database)
	})

	t.Run("GetFirstAccessibleSchema returns first schema", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetRoleMetadata(gomock.Any(), "testuser").
			Return(&domain.RoleMetadata{
				Name:              "testuser",
				AccessibleSchemas: []string{"public", "private"},
			}, nil)

		schema, err := uc.GetFirstAccessibleSchema(ctx, "testuser", "testdb")

		require.NoError(t, err)
		require.NotEmpty(t, schema)
		require.Equal(t, "public", schema)
	})

	t.Run("GetFirstAccessibleTable returns first table", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetRoleMetadata(gomock.Any(), "testuser").
			Return(&domain.RoleMetadata{
				Name: "testuser",
				AccessibleTables: []domain.AccessibleTable{
					{Database: "testdb", Schema: "public", Name: "users"},
					{Database: "testdb", Schema: "public", Name: "posts"},
				},
			}, nil)

		table, err := uc.GetFirstAccessibleTable(ctx, "testuser", "testdb", "public")

		require.NoError(t, err)
		require.NotEmpty(t, table)
		require.Equal(t, "users", table)
	})

	// UC-S2-08: Session Validation - Valid Session
	t.Run("ValidateSession returns session for valid ID", func(t *testing.T) {
		expectedSession := &domain.Session{
			ID:       "session_123",
			Username: "testuser",
		}

		mockSession.EXPECT().
			GetSession(gomock.Any(), "session_123").
			Return(expectedSession, nil)

		session, err := uc.ValidateSession(ctx, "session_123")

		require.NoError(t, err)
		require.NotNil(t, session)
		require.Equal(t, expectedSession.ID, session.ID)
	})

	// UC-S2-09: Session Validation - Expired Session
	t.Run("ValidateSession returns error for expired session", func(t *testing.T) {
		mockSession.EXPECT().
			GetSession(gomock.Any(), "expired_session").
			Return(nil, ErrSessionExpired)

		session, err := uc.ValidateSession(ctx, "expired_session")

		require.Error(t, err)
		require.Nil(t, session)
	})

	// UC-S2-10: Session Re-authentication
	t.Run("RefreshSession extends expiration time", func(t *testing.T) {
		oldSession := &domain.Session{
			ID:       "session_123",
			Username: "testuser",
		}

		mockSession.EXPECT().
			GetSession(gomock.Any(), "session_123").
			Return(oldSession, nil)

		mockSession.EXPECT().
			UpdateSession(gomock.Any(), gomock.Any()).
			Return(nil)

		session, err := uc.RefreshSession(ctx, "session_123")

		require.NoError(t, err)
		require.NotNil(t, session)
	})

	// UC-S2-12: Logout Cookie Clearing
	// E2E-S2-04: Logout Flow
	t.Run("Logout invalidates session", func(t *testing.T) {
		mockSession.EXPECT().
			DeleteSession(gomock.Any(), "session_123").
			Return(nil)

		err := uc.Logout(ctx, "session_123")

		require.NoError(t, err)
	})

	// UC-S2-13: Header Username Display
	t.Run("GetSessionUser returns user for valid session", func(t *testing.T) {
		expectedUser := &domain.User{
			Username:     "testuser",
			DatabaseName: "testdb",
			SchemaName:   "public",
			TableName:    "users",
		}

		mockSession.EXPECT().
			GetSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockMetadata.EXPECT().
			GetRoleMetadata(gomock.Any(), "testuser").
			Return(&domain.RoleMetadata{
				Name:                "testuser",
				AccessibleDatabases: []string{"testdb"},
				AccessibleSchemas:   []string{"public"},
				AccessibleTables: []domain.AccessibleTable{
					{Database: "testdb", Schema: "public", Name: "users"},
				},
			}, nil)

		user, err := uc.GetSessionUser(ctx, "session_123")

		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, "testuser", user.Username)
	})

	// UC-S2-14: Navigation Menu Rendering
	// E2E-S2-05: Protected Route Access Without Auth
	t.Run("IsUserAuthenticated returns true for valid session", func(t *testing.T) {
		mockSession.EXPECT().
			GetSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		authenticated, err := uc.IsUserAuthenticated(ctx, "session_123")

		require.NoError(t, err)
		require.True(t, authenticated)
	})

	t.Run("IsUserAuthenticated returns false for invalid session", func(t *testing.T) {
		mockSession.EXPECT().
			GetSession(gomock.Any(), "invalid_session").
			Return(nil, ErrSessionNotFound)

		authenticated, err := uc.IsUserAuthenticated(ctx, "invalid_session")

		require.Error(t, err)
		require.False(t, authenticated)
	})

	// UC-S2-15: Metadata Refresh Button
	t.Run("ReAuthenticateWithPassword verifies encrypted password", func(t *testing.T) {
		mockEncryption.EXPECT().
			DecryptPassword(gomock.Any(), "encrypted_password").
			Return("password123", nil)

		mockDatabase.EXPECT().
			TestConnection(gomock.Any(), gomock.Any()).
			Return(nil)

		verified, err := uc.ReAuthenticateWithPassword(ctx, "testuser", "encrypted_password")

		require.NoError(t, err)
		require.True(t, verified)
	})

	// IT-S2-05: Concurrent User Sessions with Isolated Resources
	t.Run("Multiple sessions remain isolated", func(t *testing.T) {
		mockSession.EXPECT().
			GetSession(gomock.Any(), "session_user1").
			Return(&domain.Session{
				ID:       "session_user1",
				Username: "user1",
			}, nil).AnyTimes()

		mockSession.EXPECT().
			GetSession(gomock.Any(), "session_user2").
			Return(&domain.Session{
				ID:       "session_user2",
				Username: "user2",
			}, nil).AnyTimes()

		user1, err1 := uc.GetSessionUser(ctx, "session_user1")
		user2, err2 := uc.GetSessionUser(ctx, "session_user2")

		require.NoError(t, err1)
		require.NoError(t, err2)
		require.NotEqual(t, user1.Username, user2.Username)
	})
}

// Error types for authentication
var (
	ErrInvalidCredentials = domain.ValidationError{Field: "credentials", Message: "invalid credentials"}
	ErrSessionExpired     = domain.ValidationError{Field: "session", Message: "session expired"}
	ErrSessionNotFound    = domain.ValidationError{Field: "session", Message: "session not found"}
)
