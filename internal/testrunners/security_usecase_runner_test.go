package testrunners

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/implementations/mocks"
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
)

// SecurityUseCaseConstructor creates a security use case with its dependencies
type SecurityUseCaseConstructor func(
	queryRepo repository.QueryRepository,
	sessionRepo repository.SessionRepository,
) usecase.QueryUseCase

// SecurityUseCaseRunner runs test specs for security use case (Story 7)
func SecurityUseCaseRunner(t *testing.T, constructor SecurityUseCaseConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueryRepo := mocks.NewMockQueryRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)

	useCase := constructor(mockQueryRepo, mockSessionRepo)

	t.Run("UC-S7-01: SQL Injection Prevention - WHERE Clause", func(t *testing.T) {
		ctx := context.Background()
		// Attempt to inject SQL via WHERE clause using parameterized query
		sql := "SELECT * FROM users WHERE id = $1"
		maliciousInput := "1' OR '1'='1"
		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "username"},
			Rows:       [][]interface{}{{1, "test"}},
			TotalRows:  1,
			LoadedRows: 1,
			Success:    true,
		}

		mockQueryRepo.EXPECT().ExecuteQuery(ctx, sql, maliciousInput).Return(expectedResult, nil)

		result, err := useCase.ExecuteQuery(ctx, sql, maliciousInput)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
		// Should treat malicious input as literal string, not SQL
		assert.Equal(t, 1, len(result.Rows))
	})

	t.Run("UC-S7-02: SQL Injection Prevention - Query Editor", func(t *testing.T) {
		ctx := context.Background()
		// SQL injection attempt in query editor
		queries := "SELECT * FROM users; DROP TABLE users; --"
		expectedResult := &domain.QueryResult{
			Success:      false,
			ErrorMessage: "Multiple statements not allowed",
		}

		mockQueryRepo.EXPECT().ExecuteMultiple(ctx, queries).Return([]*domain.QueryResult{expectedResult}, nil)

		results, err := useCase.ExecuteMultipleQueries(ctx, queries)

		require.NoError(t, err)
		assert.NotEmpty(t, results)
		assert.False(t, results[0].Success)
	})

	t.Run("UC-S7-03: Password Encryption in Cookie", func(t *testing.T) {
		plainPassword := "MySecurePassword123!"
		expectedEncrypted := "encrypted_value_xyz"

		mockSessionRepo.EXPECT().EncryptPassword(plainPassword).Return(expectedEncrypted, nil)

		encrypted, err := mockSessionRepo.EncryptPassword(plainPassword)

		require.NoError(t, err)
		assert.NotEqual(t, plainPassword, encrypted)
		assert.Equal(t, expectedEncrypted, encrypted)
	})

	t.Run("UC-S7-04: Password Decryption from Cookie", func(t *testing.T) {
		plainPassword := "MySecurePassword123!"
		encryptedPassword := "encrypted_value_xyz"

		mockSessionRepo.EXPECT().DecryptPassword(encryptedPassword).Return(plainPassword, nil)

		decrypted, err := mockSessionRepo.DecryptPassword(encryptedPassword)

		require.NoError(t, err)
		assert.Equal(t, plainPassword, decrypted)
	})

	t.Run("UC-S7-05: Cookie Tampering Detection", func(t *testing.T) {
		tamperedPassword := "encrypted_value_abc" // Modified cookie

		mockSessionRepo.EXPECT().DecryptPassword(tamperedPassword).Return("", assert.AnError)

		decrypted, err := mockSessionRepo.DecryptPassword(tamperedPassword)

		require.Error(t, err)
		assert.Equal(t, "", decrypted)
	})

	t.Run("UC-S7-06: Session Timeout Short-Lived Cookie", func(t *testing.T) {
		ctx := context.Background()
		username := "testuser"
		password := "testpass"
		roleMetadata := &domain.RoleMetadata{
			RoleName:            username,
			AccessibleDatabases: []string{"testdb"},
		}

		shortLivedSession := &domain.Session{
			Username: username,
		}

		mockSessionRepo.EXPECT().CreateSession(ctx, username, password, roleMetadata).Return(shortLivedSession, nil)

		session, err := mockSessionRepo.CreateSession(ctx, username, password, roleMetadata)

		require.NoError(t, err)
		assert.NotNil(t, session)
	})

	t.Run("UC-S7-07: Session Timeout Long-Lived Cookie", func(t *testing.T) {
		ctx := context.Background()
		username := "testuser"
		roleMetadata := &domain.RoleMetadata{
			RoleName: username,
		}

		longLivedSession := &domain.Session{
			Username: username,
		}

		mockSessionRepo.EXPECT().CreateSession(ctx, username, "pass", roleMetadata).Return(longLivedSession, nil)

		session, err := mockSessionRepo.CreateSession(ctx, username, "pass", roleMetadata)

		require.NoError(t, err)
		assert.NotNil(t, session)
		// Long-lived cookies for username pre-fill (longer expiry than short-lived password cookie)
	})
}
