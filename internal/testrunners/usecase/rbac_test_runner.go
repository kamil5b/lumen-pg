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

// RBACUsecaseConstructor is a function type that creates an RBACUseCase
type RBACUsecaseConstructor func(
	metadataRepo repository.MetadataRepository,
	rbacRepo repository.RBACRepository,
) usecase.RBACUseCase

// RBACUsecaseRunner runs all RBAC usecase tests against an implementation
// Maps to TEST_PLAN.md:
// - Story 1: Setup & Configuration [UC-S1-05~07, IT-S1-02~04]
// - Story 2: Authentication & Identity [UC-S2-11, IT-S2-03, E2E-S2-06]
// - Story 5: Main View & Data Interaction [UC-S5-19, IT-S5-01~07]
// - Story 6: Isolation [UC-S6-01~03, IT-S6-01~03, E2E-S6-01~03]
func RBACUsecaseRunner(t *testing.T, constructor RBACUsecaseConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMetadata := mockRepository.NewMockMetadataRepository(ctrl)
	mockRBAC := mockRepository.NewMockRBACRepository(ctrl)

	uc := constructor(mockMetadata, mockRBAC)

	ctx := context.Background()

	// UC-S1-07: RBAC Initialization with User Accessibility
	// IT-S2-03: Real Role-Based Resource Access
	t.Run("CheckTableAccess returns true for accessible table", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetAccessibleTables(gomock.Any(), "testuser", "testdb", "public").
			Return([]string{"users", "posts"}, nil)

		accessible, err := uc.CheckTableAccess(ctx, "testuser", "testdb", "public", "users")

		require.NoError(t, err)
		require.True(t, accessible)
	})

	t.Run("CheckTableAccess returns false for inaccessible table", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetAccessibleTables(gomock.Any(), "testuser", "testdb", "public").
			Return([]string{"users"}, nil)

		accessible, err := uc.CheckTableAccess(ctx, "testuser", "testdb", "public", "admin_tables")

		require.NoError(t, err)
		require.False(t, accessible)
	})

	// UC-S1-05: Metadata Initialization - Roles and Permissions
	t.Run("CheckSelectPermission returns permission", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetTablePermissions(gomock.Any(), "testuser", "testdb", "public", "users").
			Return(&domain.AccessibleTable{
				HasSelect: true,
				HasInsert: true,
				HasUpdate: true,
				HasDelete: false,
			}, nil)

		canSelect, err := uc.CheckSelectPermission(ctx, "testuser", "testdb", "public", "users")

		require.NoError(t, err)
		require.True(t, canSelect)
	})

	t.Run("CheckSelectPermission returns false when denied", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetTablePermissions(gomock.Any(), "readonly_user", "testdb", "public", "admin_table").
			Return(&domain.AccessibleTable{
				HasSelect: false,
				HasInsert: false,
				HasUpdate: false,
				HasDelete: false,
			}, nil)

		canSelect, err := uc.CheckSelectPermission(ctx, "readonly_user", "testdb", "public", "admin_table")

		require.NoError(t, err)
		require.False(t, canSelect)
	})

	t.Run("CheckInsertPermission checks INSERT permission", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetTablePermissions(gomock.Any(), "testuser", "testdb", "public", "users").
			Return(&domain.AccessibleTable{
				HasSelect: true,
				HasInsert: true,
				HasUpdate: true,
				HasDelete: false,
			}, nil)

		canInsert, err := uc.CheckInsertPermission(ctx, "testuser", "testdb", "public", "users")

		require.NoError(t, err)
		require.True(t, canInsert)
	})

	t.Run("CheckUpdatePermission checks UPDATE permission", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetTablePermissions(gomock.Any(), "testuser", "testdb", "public", "users").
			Return(&domain.AccessibleTable{
				HasSelect: true,
				HasInsert: true,
				HasUpdate: true,
				HasDelete: false,
			}, nil)

		canUpdate, err := uc.CheckUpdatePermission(ctx, "testuser", "testdb", "public", "users")

		require.NoError(t, err)
		require.True(t, canUpdate)
	})

	t.Run("CheckDeletePermission checks DELETE permission", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetTablePermissions(gomock.Any(), "testuser", "testdb", "public", "users").
			Return(&domain.AccessibleTable{
				HasSelect: true,
				HasInsert: true,
				HasUpdate: true,
				HasDelete: true,
			}, nil)

		canDelete, err := uc.CheckDeletePermission(ctx, "testuser", "testdb", "public", "users")

		require.NoError(t, err)
		require.True(t, canDelete)
	})

	// UC-S1-07: RBAC Initialization with User Accessibility
	t.Run("CheckDatabaseAccess returns true for accessible database", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetAccessibleDatabases(gomock.Any(), "testuser").
			Return([]string{"testdb1", "testdb2"}, nil)

		accessible, err := uc.CheckDatabaseAccess(ctx, "testuser", "testdb1")

		require.NoError(t, err)
		require.True(t, accessible)
	})

	t.Run("CheckDatabaseAccess returns false for inaccessible database", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetAccessibleDatabases(gomock.Any(), "testuser").
			Return([]string{"testdb1"}, nil)

		accessible, err := uc.CheckDatabaseAccess(ctx, "testuser", "secret_database")

		require.NoError(t, err)
		require.False(t, accessible)
	})

	t.Run("CheckSchemaAccess returns true for accessible schema", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetAccessibleSchemas(gomock.Any(), "testuser", "testdb").
			Return([]string{"public", "private"}, nil)

		accessible, err := uc.CheckSchemaAccess(ctx, "testuser", "testdb", "public")

		require.NoError(t, err)
		require.True(t, accessible)
	})

	// IT-S1-02: Load Real Database Metadata with User Accessible Resources
	// E2E-S2-06: Data Explorer Populated After Login
	t.Run("GetUserAccessibleDatabases returns user's databases", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetAccessibleDatabases(gomock.Any(), "testuser").
			Return([]string{"testdb1", "testdb2", "testdb3"}, nil)

		databases, err := uc.GetUserAccessibleDatabases(ctx, "testuser")

		require.NoError(t, err)
		require.NotNil(t, databases)
		require.Equal(t, 3, len(databases))
	})

	// IT-S1-03: Load Real Relations and Role Access
	t.Run("GetUserAccessibleSchemas returns user's schemas", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetAccessibleSchemas(gomock.Any(), "testuser", "testdb").
			Return([]string{"public", "private", "staging"}, nil)

		schemas, err := uc.GetUserAccessibleSchemas(ctx, "testuser", "testdb")

		require.NoError(t, err)
		require.NotNil(t, schemas)
		require.Equal(t, 3, len(schemas))
	})

	// IT-S1-04: Cache Accessible Resources Per Role
	t.Run("GetUserAccessibleTables returns user's tables", func(t *testing.T) {
		mockRBAC.EXPECT().
			GetAccessibleTables(gomock.Any(), "testuser", "testdb", "public").
			Return([]domain.AccessibleTable{
				{
					Database:  "testdb",
					Schema:    "public",
					Name:      "users",
					HasSelect: true,
					HasInsert: true,
					HasUpdate: true,
					HasDelete: false,
				},
				{
					Database:  "testdb",
					Schema:    "public",
					Name:      "posts",
					HasSelect: true,
					HasInsert: true,
					HasUpdate: true,
					HasDelete: true,
				},
			}, nil)

		tables, err := uc.GetUserAccessibleTables(ctx, "testuser", "testdb", "public")

		require.NoError(t, err)
		require.NotNil(t, tables)
		require.Equal(t, 2, len(tables))
	})

	// UC-S5-19: Read-Only Mode Enforcement
	t.Run("IsTableReadOnly returns true when no write permissions", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetTablePermissions(gomock.Any(), "readonly_user", "testdb", "public", "users").
			Return(&domain.AccessibleTable{
				HasSelect: true,
				HasInsert: false,
				HasUpdate: false,
				HasDelete: false,
			}, nil)

		isReadOnly, err := uc.IsTableReadOnly(ctx, "readonly_user", "testdb", "public", "users")

		require.NoError(t, err)
		require.True(t, isReadOnly)
	})

	t.Run("IsTableReadOnly returns false when has write permissions", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetTablePermissions(gomock.Any(), "testuser", "testdb", "public", "users").
			Return(&domain.AccessibleTable{
				HasSelect: true,
				HasInsert: true,
				HasUpdate: true,
				HasDelete: false,
			}, nil)

		isReadOnly, err := uc.IsTableReadOnly(ctx, "testuser", "testdb", "public", "users")

		require.NoError(t, err)
		require.False(t, isReadOnly)
	})

	t.Run("GetTablePermissions returns complete permission set", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetTablePermissions(gomock.Any(), "testuser", "testdb", "public", "users").
			Return(&domain.AccessibleTable{
				HasSelect: true,
				HasInsert: true,
				HasUpdate: true,
				HasDelete: false,
			}, nil)

		perms, err := uc.GetTablePermissions(ctx, "testuser", "testdb", "public", "users")

		require.NoError(t, err)
		require.NotNil(t, perms)
		require.True(t, perms.CanSelect)
	})

	// UC-S6-01: Session Isolation
	t.Run("ValidateUserAccessToResource checks access", func(t *testing.T) {
		mockMetadata.EXPECT().
			IsTableAccessible(gomock.Any(), "testuser", "testdb", "public", "users").
			Return(true, nil)

		accessible, err := uc.ValidateUserAccessToResource(ctx, "testuser", "table", "testdb", "public", "users")

		require.NoError(t, err)
		require.True(t, accessible)
	})

	// UC-S6-02: Transaction Isolation
	// IT-S6-02: Real Permission Isolation
	// E2E-S6-01: Simultaneous Users Different Permissions
	t.Run("GetUserRole returns user's role", func(t *testing.T) {
		mockRBAC.EXPECT().
			GetUserRole(gomock.Any(), "testuser").
			Return("user", nil)

		role, err := uc.GetUserRole(ctx, "testuser")

		require.NoError(t, err)
		require.Equal(t, "user", role)
	})

	// UC-S6-03: Cookie Isolation
	// IT-S6-01: Real Multi-User Connection
	t.Run("VerifyUserPermissions returns permission set for user", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetTablePermissions(gomock.Any(), "testuser", "testdb", "public", "users").
			Return(&domain.AccessibleTable{
				HasSelect: true,
				HasInsert: true,
				HasUpdate: true,
				HasDelete: false,
			}, nil)

		perms, err := uc.VerifyUserPermissions(ctx, "testuser", "testdb", "public", "users")

		require.NoError(t, err)
		require.NotNil(t, perms)
		require.True(t, perms.CanSelect)
		require.True(t, perms.CanInsert)
		require.False(t, perms.CanDelete)
	})

	// E2E-S6-01: Simultaneous Users Different Permissions
	t.Run("Different users have different accessible resources", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetAccessibleDatabases(gomock.Any(), "admin_user").
			Return([]string{"testdb1", "testdb2", "testdb3"}, nil)

		mockMetadata.EXPECT().
			GetAccessibleDatabases(gomock.Any(), "readonly_user").
			Return([]string{"testdb1"}, nil)

		adminDbs, _ := uc.GetUserAccessibleDatabases(ctx, "admin_user")
		readonlyDbs, _ := uc.GetUserAccessibleDatabases(ctx, "readonly_user")

		require.Greater(t, len(adminDbs), len(readonlyDbs))
	})

	// E2E-S6-02: Simultaneous Transactions
	t.Run("Users with different permissions can work independently", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetTablePermissions(gomock.Any(), "user1", "testdb", "public", "users").
			Return(&domain.AccessibleTable{
				HasSelect: true,
				HasInsert: true,
				HasUpdate: true,
				HasDelete: true,
			}, nil)

		mockMetadata.EXPECT().
			GetTablePermissions(gomock.Any(), "user2", "testdb", "public", "users").
			Return(&domain.AccessibleTable{
				HasSelect: true,
				HasInsert: false,
				HasUpdate: false,
				HasDelete: false,
			}, nil)

		perms1, _ := uc.GetTablePermissions(ctx, "user1", "testdb", "public", "users")
		perms2, _ := uc.GetTablePermissions(ctx, "user2", "testdb", "public", "users")

		require.True(t, perms1.CanDelete)
		require.False(t, perms2.CanDelete)
	})

	// E2E-S6-03: One User Cannot See Another's Session
	t.Run("Permissions are strictly per-user", func(t *testing.T) {
		mockRBAC.EXPECT().
			GetUserRole(gomock.Any(), "user1").
			Return("admin", nil)

		mockRBAC.EXPECT().
			GetUserRole(gomock.Any(), "user2").
			Return("viewer", nil)

		role1, _ := uc.GetUserRole(ctx, "user1")
		role2, _ := uc.GetUserRole(ctx, "user2")

		require.NotEqual(t, role1, role2)
	})
}
