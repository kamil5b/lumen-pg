package testrunners

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
	"github.com/kamil5b/lumen-pg/internal/implementations/mocks"
)

// AuthUseCaseConstructor creates an auth use case with its dependencies
type AuthUseCaseConstructor func(
	connRepo repository.ConnectionRepository,
	metadataRepo repository.MetadataRepository,
	sessionRepo repository.SessionRepository,
) usecase.AuthUseCase

// AuthUseCaseRunner runs test specs for auth use case (Story 2)
func AuthUseCaseRunner(t *testing.T, constructor AuthUseCaseConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnRepo := mocks.NewMockConnectionRepository(ctrl)
	mockMetadataRepo := mocks.NewMockMetadataRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	
	useCase := constructor(mockConnRepo, mockMetadataRepo, mockSessionRepo)

	t.Run("UC-S2-01: Login Form Validation - Empty Username", func(t *testing.T) {
		ctx := context.Background()
		req := domain.LoginRequest{
			Username: "",
			Password: "password",
		}

		resp, err := useCase.Login(ctx, req)

		require.NoError(t, err) // No error, but response indicates failure
		assert.NotNil(t, resp)
		assert.False(t, resp.Success)
		assert.Contains(t, resp.ErrorMessage, "username")
	})

	t.Run("UC-S2-02: Login Form Validation - Empty Password", func(t *testing.T) {
		ctx := context.Background()
		req := domain.LoginRequest{
			Username: "user",
			Password: "",
		}

		resp, err := useCase.Login(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.False(t, resp.Success)
		assert.Contains(t, resp.ErrorMessage, "password")
	})

	t.Run("UC-S2-03: Login Connection Probe Success", func(t *testing.T) {
		ctx := context.Background()
		req := domain.LoginRequest{
			Username: "testuser",
			Password: "testpass",
		}

		roleMetadata := &domain.RoleMetadata{
			RoleName:            "testuser",
			AccessibleDatabases: []string{"testdb"},
			AccessibleTables:    map[string][]string{"testdb.public": {"users"}},
		}

		mockMetadataRepo.EXPECT().LoadRolePermissions(ctx, "testuser").Return(roleMetadata, nil)
		mockConnRepo.EXPECT().ValidateConnection(ctx, gomock.Any()).Return(nil)
		mockSessionRepo.EXPECT().CreateSession(ctx, "testuser", "testpass", roleMetadata).Return(&domain.Session{
			Username: "testuser",
		}, nil)

		resp, err := useCase.Login(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.True(t, resp.Success)
		assert.Equal(t, "testdb", resp.FirstAccessibleDB)
		assert.Equal(t, "users", resp.FirstAccessibleTbl)
	})

	t.Run("UC-S2-04: Login Connection Probe Failure - No Accessible Resources", func(t *testing.T) {
		ctx := context.Background()
		req := domain.LoginRequest{
			Username: "restricteduser",
			Password: "testpass",
		}

		roleMetadata := &domain.RoleMetadata{
			RoleName:            "restricteduser",
			AccessibleDatabases: []string{},
			AccessibleTables:    map[string][]string{},
		}

		mockMetadataRepo.EXPECT().LoadRolePermissions(ctx, "restricteduser").Return(roleMetadata, nil)

		resp, err := useCase.Login(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.False(t, resp.Success)
		assert.Contains(t, resp.ErrorMessage, "no accessible resources")
	})
}
