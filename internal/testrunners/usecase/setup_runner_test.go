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

// SetupUsecaseConstructor is a function type that creates a SetupUseCase
type SetupUsecaseConstructor func(
	databaseRepo repository.DatabaseRepository,
	metadataRepo repository.MetadataRepository,
	rbacRepo repository.RBACRepository,
) usecase.SetupUseCase

// SetupUsecaseRunner runs all setup usecase tests against an implementation
// Maps to TEST_PLAN.md:
// - Story 1: Setup & Configuration [UC-S1-01~07, IT-S1-01~04]
func SetupUsecaseRunner(t *testing.T, constructor SetupUsecaseConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDatabase := mockRepository.NewMockDatabaseRepository(ctrl)
	mockMetadata := mockRepository.NewMockMetadataRepository(ctrl)
	mockRBAC := mockRepository.NewMockRBACRepository(ctrl)

	uc := constructor(mockDatabase, mockMetadata, mockRBAC)

	ctx := context.Background()

	// UC-S1-01: Connection String Validation
	t.Run("ValidateConnectionString accepts valid string", func(t *testing.T) {
		valid, err := uc.ValidateConnectionString(ctx, "postgres://user:password@localhost:5432/testdb?sslmode=disable")

		require.NoError(t, err)
		require.True(t, valid)
	})

	t.Run("ValidateConnectionString rejects invalid string", func(t *testing.T) {
		valid, err := uc.ValidateConnectionString(ctx, "invalid connection string")

		require.NoError(t, err)
		require.False(t, valid)
	})

	// UC-S1-02: Connection String Parsing
	t.Run("ParseConnectionString returns parsed components", func(t *testing.T) {
		connStr := "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"

		parsed, err := uc.ParseConnectionString(ctx, connStr)

		require.NoError(t, err)
		require.NotNil(t, parsed)
		require.Equal(t, "testuser", parsed.Username)
		require.Equal(t, "testdb", parsed.Database)
		require.Equal(t, "localhost", parsed.Host)
		require.Equal(t, "5432", parsed.Port)
	})

	t.Run("ParseConnectionString handles missing components", func(t *testing.T) {
		connStr := "postgres://testuser@localhost/testdb"

		parsed, err := uc.ParseConnectionString(ctx, connStr)

		require.NoError(t, err)
		require.NotNil(t, parsed)
		require.Equal(t, "testuser", parsed.Username)
		require.Equal(t, "testdb", parsed.Database)
	})

	// UC-S1-03: Superadmin Connection Test Success
	// IT-S1-01: Connect to Real PostgreSQL
	t.Run("TestSuperadminConnection succeeds with valid credentials", func(t *testing.T) {
		mockDatabase.EXPECT().
			TestConnection(gomock.Any(), gomock.Any()).
			Return(nil)

		success, err := uc.TestSuperadminConnection(ctx, "postgres://user:pass@localhost/postgres?sslmode=disable")

		require.NoError(t, err)
		require.True(t, success)
	})

	// UC-S1-04: Superadmin Connection Test Failure
	t.Run("TestSuperadminConnection fails with invalid credentials", func(t *testing.T) {
		mockDatabase.EXPECT().
			TestConnection(gomock.Any(), gomock.Any()).
			Return(ErrConnectionFailed)

		success, err := uc.TestSuperadminConnection(ctx, "postgres://invalid:invalid@localhost/postgres?sslmode=disable")

		require.Error(t, err)
		require.False(t, success)
	})

	// UC-S1-05: Metadata Initialization - Roles and Permissions
	// UC-S1-06: In-Memory Metadata Storage - Per Role
	// UC-S1-07: RBAC Initialization with User Accessibility
	// IT-S1-02: Load Real Database Metadata with User Accessible Resources
	// IT-S1-03: Load Real Relations and Role Access
	// IT-S1-04: Cache Accessible Resources Per Role
	t.Run("InitializeMetadata fetches and caches all metadata", func(t *testing.T) {
		mockDatabase.EXPECT().
			TestConnection(gomock.Any(), gomock.Any()).
			Return(nil)

		mockDatabase.EXPECT().
			GetDatabaseMetadata(gomock.Any(), gomock.Any()).
			Return(&domain.DatabaseMetadata{
				Name: "testdb",
				Schemas: []domain.SchemaMetadata{
					{
						Name: "public",
						Tables: []domain.TableMetadata{
							{
								Name: "users",
								Columns: []domain.ColumnMetadata{
									{
										Name:       "id",
										DataType:   "integer",
										IsNullable: false,
										IsPrimary:  true,
									},
									{
										Name:       "name",
										DataType:   "text",
										IsNullable: true,
										IsPrimary:  false,
									},
								},
								PrimaryKeys: []string{"id"},
								ForeignKeys: []domain.ForeignKeyMetadata{},
							},
						},
					},
				},
			}, nil)

		mockMetadata.EXPECT().
			StoreMetadata(gomock.Any(), gomock.Any()).
			Return(nil)

		err := uc.InitializeMetadata(ctx, "postgres://user:pass@localhost/testdb?sslmode=disable")

		require.NoError(t, err)
	})

	t.Run("InitializeRBAC fetches and caches role metadata", func(t *testing.T) {
		mockRBAC.EXPECT().
			GetAllRoles(gomock.Any()).
			Return([]string{"postgres", "testuser", "readonly"}, nil)

		mockRBAC.EXPECT().
			GetRoleAccessibility(gomock.Any(), gomock.Any()).
			Return(&domain.RoleMetadata{
				Name:                "testuser",
				AccessibleDatabases: []string{"testdb"},
				AccessibleSchemas:   []string{"public"},
				AccessibleTables: []domain.AccessibleTable{
					{
						Database:  "testdb",
						Schema:    "public",
						Name:      "users",
						HasSelect: true,
						HasInsert: true,
						HasUpdate: true,
						HasDelete: false,
					},
				},
			}, nil).Times(3)

		mockMetadata.EXPECT().
			StoreAllRolesMetadata(gomock.Any(), gomock.Any()).
			Return(nil)

		err := uc.InitializeRBAC(ctx, "postgres://user:pass@localhost/testdb?sslmode=disable")

		require.NoError(t, err)
	})

	t.Run("GetAllRoles returns list of database roles", func(t *testing.T) {
		mockRBAC.EXPECT().
			GetAllRoles(gomock.Any()).
			Return([]string{"postgres", "testuser", "readonly"}, nil)

		roles, err := uc.GetAllRoles(ctx)

		require.NoError(t, err)
		require.NotNil(t, roles)
		require.GreaterOrEqual(t, len(roles), 3)
	})

	t.Run("GetRoleAccessibility returns role metadata", func(t *testing.T) {
		mockRBAC.EXPECT().
			GetRoleAccessibility(gomock.Any(), "testuser").
			Return(&domain.RoleMetadata{
				Name:                "testuser",
				AccessibleDatabases: []string{"testdb1", "testdb2"},
				AccessibleSchemas:   []string{"public", "private"},
				AccessibleTables: []domain.AccessibleTable{
					{
						Database:  "testdb1",
						Schema:    "public",
						Name:      "users",
						HasSelect: true,
						HasInsert: true,
						HasUpdate: true,
						HasDelete: true,
					},
					{
						Database:  "testdb1",
						Schema:    "public",
						Name:      "posts",
						HasSelect: true,
						HasInsert: true,
						HasUpdate: false,
						HasDelete: false,
					},
				},
			}, nil)

		metadata, err := uc.GetRoleAccessibility(ctx, "testuser")

		require.NoError(t, err)
		require.NotNil(t, metadata)
		require.Equal(t, "testuser", metadata.Name)
		require.GreaterOrEqual(t, len(metadata.AccessibleDatabases), 1)
		require.GreaterOrEqual(t, len(metadata.AccessibleTables), 1)
	})

	// UC-S2-15: Metadata Refresh Button
	t.Run("RefreshMetadata reloads all cached metadata", func(t *testing.T) {
		mockMetadata.EXPECT().
			InvalidateAllMetadata(gomock.Any()).
			Return(nil)

		mockDatabase.EXPECT().
			GetDatabaseMetadata(gomock.Any(), gomock.Any()).
			Return(&domain.DatabaseMetadata{
				Name: "testdb",
				Schemas: []domain.SchemaMetadata{
					{
						Name:   "public",
						Tables: []domain.TableMetadata{},
					},
				},
			}, nil)

		mockMetadata.EXPECT().
			StoreMetadata(gomock.Any(), gomock.Any()).
			Return(nil)

		err := uc.RefreshMetadata(ctx)

		require.NoError(t, err)
	})

	t.Run("RefreshRBACMetadata reloads all cached RBAC metadata", func(t *testing.T) {
		mockMetadata.EXPECT().
			InvalidateAllMetadata(gomock.Any()).
			Return(nil)

		mockRBAC.EXPECT().
			GetAllRoles(gomock.Any()).
			Return([]string{"testuser"}, nil)

		mockRBAC.EXPECT().
			GetRoleAccessibility(gomock.Any(), gomock.Any()).
			Return(&domain.RoleMetadata{
				Name:                "testuser",
				AccessibleDatabases: []string{"testdb"},
				AccessibleSchemas:   []string{"public"},
				AccessibleTables:    []domain.AccessibleTable{},
			}, nil)

		mockMetadata.EXPECT().
			StoreAllRolesMetadata(gomock.Any(), gomock.Any()).
			Return(nil)

		err := uc.RefreshRBACMetadata(ctx)

		require.NoError(t, err)
	})

	t.Run("IsInitialized returns true after successful initialization", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetMetadata(gomock.Any(), gomock.Any()).
			Return(&domain.DatabaseMetadata{
				Name: "testdb",
				Schemas: []domain.SchemaMetadata{
					{
						Name:   "public",
						Tables: []domain.TableMetadata{},
					},
				},
			}, nil)

		initialized, err := uc.IsInitialized(ctx)

		require.NoError(t, err)
		require.True(t, initialized)
	})

	t.Run("IsInitialized returns false when not initialized", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetMetadata(gomock.Any(), gomock.Any()).
			Return(nil, ErrMetadataNotFound)

		initialized, err := uc.IsInitialized(ctx)

		require.NoError(t, err)
		require.False(t, initialized)
	})
}

// Error types for setup
var (
	ErrConnectionFailed = domain.ValidationError{Field: "connection", Message: "connection failed"}
	ErrMetadataNotFound = domain.ValidationError{Field: "metadata", Message: "metadata not found"}
)
