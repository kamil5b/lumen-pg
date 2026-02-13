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

// MetadataUseCaseConstructor creates a metadata use case with its dependencies
type MetadataUseCaseConstructor func(repo repository.MetadataRepository) usecase.MetadataUseCase

// MetadataUseCaseRunner runs test specs for metadata use case (Story 1)
func MetadataUseCaseRunner(t *testing.T, constructor MetadataUseCaseConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMetadataRepository(ctrl)
	useCase := constructor(mockRepo)

	t.Run("UC-S1-05: Metadata Initialization - Roles and Permissions", func(t *testing.T) {
		ctx := context.Background()
		expectedMetadata := &domain.GlobalMetadata{
			Databases: []domain.DatabaseMetadata{
				{
					Name: "testdb",
					Schemas: []domain.SchemaMetadata{
						{
							Name: "public",
							Tables: []domain.TableMetadata{
								{
									SchemaName: "public",
									TableName:  "users",
									Columns: []domain.ColumnMetadata{
										{Name: "id", DataType: "integer"},
										{Name: "username", DataType: "varchar"},
									},
								},
							},
						},
					},
				},
			},
			Roles: []domain.RoleMetadata{
				{
					RoleName:            "admin",
					AccessibleDatabases: []string{"testdb"},
				},
			},
		}

		mockRepo.EXPECT().LoadGlobalMetadata(ctx).Return(expectedMetadata, nil)

		result, err := useCase.LoadGlobalMetadata(ctx)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Databases, 1)
		assert.Equal(t, "testdb", result.Databases[0].Name)
		assert.Len(t, result.Roles, 1)
		assert.Equal(t, "admin", result.Roles[0].RoleName)
	})

	t.Run("UC-S1-06: In-Memory Metadata Storage - Per Role", func(t *testing.T) {
		ctx := context.Background()
		roleName := "editor"
		expectedRole := &domain.RoleMetadata{
			RoleName:            roleName,
			AccessibleDatabases: []string{"testdb"},
			AccessibleSchemas:   map[string][]string{"testdb": {"public"}},
			AccessibleTables:    map[string][]string{"testdb.public": {"users", "posts"}},
		}

		mockRepo.EXPECT().LoadRolePermissions(ctx, roleName).Return(expectedRole, nil)

		result, err := useCase.LoadRoleAccessibleResources(ctx, roleName)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, roleName, result.RoleName)
		assert.Len(t, result.AccessibleDatabases, 1)
		assert.Len(t, result.AccessibleSchemas["testdb"], 1)
		assert.Len(t, result.AccessibleTables["testdb.public"], 2)
	})

	t.Run("UC-S1-07: RBAC Initialization with User Accessibility", func(t *testing.T) {
		ctx := context.Background()
		roleName := "viewer"
		expectedRole := &domain.RoleMetadata{
			RoleName:            roleName,
			AccessibleDatabases: []string{"testdb"},
			AccessibleSchemas:   map[string][]string{"testdb": {"public"}},
			AccessibleTables:    map[string][]string{"testdb.public": {"users"}},
		}

		mockRepo.EXPECT().LoadRolePermissions(ctx, roleName).Return(expectedRole, nil)

		result, err := useCase.LoadRoleAccessibleResources(ctx, roleName)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, roleName, result.RoleName)
		assert.Contains(t, result.AccessibleDatabases, "testdb")
		assert.Contains(t, result.AccessibleTables["testdb.public"], "users")
	})
}
