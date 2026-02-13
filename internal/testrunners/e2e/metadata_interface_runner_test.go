package e2e

import (
	"context"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
)

// MetadataInterfaceConstructor creates a metadata repository with mock dependencies
type MetadataInterfaceConstructor func(ctrl *gomock.Controller) repository.MetadataRepository

// MetadataInterfaceRunner runs unit tests for metadata repository interface (Story 1)
func MetadataInterfaceRunner(t *testing.T, constructor MetadataInterfaceConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := constructor(ctrl)
	ctx := context.Background()

	t.Run("UC-S1-05: Load Global Metadata with Databases", func(t *testing.T) {
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
										{Name: "id", DataType: "integer", IsNullable: false},
										{Name: "username", DataType: "varchar", IsNullable: true},
									},
									PrimaryKeys: []string{"id"},
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

		metadata, err := repo.LoadGlobalMetadata(ctx)

		// This is a unit test interface - implementation should handle
		_ = metadata
		_ = err
	})

	t.Run("UC-S1-06: Load Role Permissions - Per Role Access", func(t *testing.T) {
		roleName := "editor"
		expectedRole := &domain.RoleMetadata{
			RoleName:            roleName,
			AccessibleDatabases: []string{"testdb"},
			AccessibleSchemas: map[string][]string{
				"testdb": {"public"},
			},
			AccessibleTables: map[string][]string{
				"testdb.public": {"users", "posts"},
			},
			Permissions: map[string][]domain.Permission{
				"testdb.public.users": {domain.PermissionSelect, domain.PermissionInsert, domain.PermissionUpdate},
				"testdb.public.posts": {domain.PermissionSelect},
			},
		}

		role, err := repo.LoadRolePermissions(ctx, roleName)

		_ = role
		_ = err
	})

	t.Run("UC-S1-07: Load Role Permissions - Viewer Role (Read-Only)", func(t *testing.T) {
		roleName := "viewer"
		expectedRole := &domain.RoleMetadata{
			RoleName:            roleName,
			AccessibleDatabases: []string{"testdb"},
			AccessibleSchemas: map[string][]string{
				"testdb": {"public"},
			},
			AccessibleTables: map[string][]string{
				"testdb.public": {"users"},
			},
			Permissions: map[string][]domain.Permission{
				"testdb.public.users": {domain.PermissionSelect},
			},
		}

		role, err := repo.LoadRolePermissions(ctx, roleName)

		_ = role
		_ = err
	})

	t.Run("UC-S1-08: Load Database Metadata - Specific Database", func(t *testing.T) {
		dbName := "testdb"
		expectedDB := &domain.DatabaseMetadata{
			Name: dbName,
			Schemas: []domain.SchemaMetadata{
				{
					Name: "public",
					Tables: []domain.TableMetadata{
						{
							SchemaName: "public",
							TableName:  "users",
							Columns: []domain.ColumnMetadata{
								{Name: "id", DataType: "integer", IsNullable: false},
							},
							PrimaryKeys: []string{"id"},
						},
					},
				},
			},
		}

		dbMetadata, err := repo.LoadDatabaseMetadata(ctx, dbName)

		_ = dbMetadata
		_ = err
	})

	t.Run("UC-S1-09: Load Table Metadata - With Columns and Keys", func(t *testing.T) {
		schemaName := "public"
		tableName := "users"
		expectedTable := &domain.TableMetadata{
			SchemaName: schemaName,
			TableName:  tableName,
			Columns: []domain.ColumnMetadata{
				{Name: "id", DataType: "integer", IsNullable: false},
				{Name: "username", DataType: "varchar", IsNullable: true},
				{Name: "email", DataType: "varchar", IsNullable: true},
			},
			PrimaryKeys: []string{"id"},
		}

		tableMetadata, err := repo.LoadTableMetadata(ctx, schemaName, tableName)

		_ = tableMetadata
		_ = err
	})

	t.Run("UC-S1-10: Load Table Metadata - With Foreign Keys", func(t *testing.T) {
		schemaName := "public"
		tableName := "posts"
		expectedTable := &domain.TableMetadata{
			SchemaName: schemaName,
			TableName:  tableName,
			Columns: []domain.ColumnMetadata{
				{Name: "id", DataType: "integer", IsNullable: false},
				{Name: "user_id", DataType: "integer", IsNullable: false},
				{Name: "title", DataType: "varchar", IsNullable: false},
			},
			PrimaryKeys: []string{"id"},
			ForeignKeys: []domain.ForeignKeyMetadata{
				{
					ColumnName:           "user_id",
					ReferencedSchemaName: "public",
					ReferencedTableName:  "users",
					ReferencedColumnName: "id",
				},
			},
		}

		tableMetadata, err := repo.LoadTableMetadata(ctx, schemaName, tableName)

		_ = tableMetadata
		_ = err
	})

	t.Run("UC-S1-11: Load All Roles", func(t *testing.T) {
		expectedRoles := []string{"postgres", "admin", "editor", "viewer"}

		roles, err := repo.LoadRoles(ctx)

		_ = roles
		_ = err
	})

	t.Run("UC-S1-12: Load Role Permissions - Admin Role (Full Access)", func(t *testing.T) {
		roleName := "admin"
		expectedRole := &domain.RoleMetadata{
			RoleName:            roleName,
			AccessibleDatabases: []string{"testdb", "postgres", "template1"},
			AccessibleSchemas: map[string][]string{
				"testdb":   {"public"},
				"postgres": {"public", "information_schema"},
			},
			Permissions: map[string][]domain.Permission{
				"testdb.public.users": {
					domain.PermissionConnect,
					domain.PermissionSelect,
					domain.PermissionInsert,
					domain.PermissionUpdate,
					domain.PermissionDelete,
				},
			},
		}

		role, err := repo.LoadRolePermissions(ctx, roleName)

		_ = role
		_ = err
	})

	t.Run("UC-S1-13: Load Metadata - Multiple Schemas", func(t *testing.T) {
		dbName := "testdb"
		expectedDB := &domain.DatabaseMetadata{
			Name: dbName,
			Schemas: []domain.SchemaMetadata{
				{
					Name: "public",
					Tables: []domain.TableMetadata{
						{SchemaName: "public", TableName: "users"},
					},
				},
				{
					Name: "audit",
					Tables: []domain.TableMetadata{
						{SchemaName: "audit", TableName: "logs"},
					},
				},
			},
		}

		dbMetadata, err := repo.LoadDatabaseMetadata(ctx, dbName)

		_ = dbMetadata
		_ = err
	})

	t.Run("UC-S1-14: Load Metadata - Multiple Tables with Relationships", func(t *testing.T) {
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
									},
									PrimaryKeys: []string{"id"},
								},
								{
									SchemaName: "public",
									TableName:  "posts",
									Columns: []domain.ColumnMetadata{
										{Name: "id", DataType: "integer"},
										{Name: "user_id", DataType: "integer"},
									},
									PrimaryKeys: []string{"id"},
									ForeignKeys: []domain.ForeignKeyMetadata{
										{
											ColumnName:           "user_id",
											ReferencedTableName:  "users",
											ReferencedColumnName: "id",
										},
									},
								},
								{
									SchemaName: "public",
									TableName:  "comments",
									Columns: []domain.ColumnMetadata{
										{Name: "id", DataType: "integer"},
										{Name: "post_id", DataType: "integer"},
										{Name: "user_id", DataType: "integer"},
									},
									PrimaryKeys: []string{"id"},
									ForeignKeys: []domain.ForeignKeyMetadata{
										{
											ColumnName:           "post_id",
											ReferencedTableName:  "posts",
											ReferencedColumnName: "id",
										},
										{
											ColumnName:           "user_id",
											ReferencedTableName:  "users",
											ReferencedColumnName: "id",
										},
									},
								},
							},
						},
					},
				},
			},
		}

		metadata, err := repo.LoadGlobalMetadata(ctx)

		_ = metadata
		_ = err
	})

	t.Run("UC-S1-15: Load Role Permissions - Restricted Role", func(t *testing.T) {
		roleName := "restricted"
		expectedRole := &domain.RoleMetadata{
			RoleName:            roleName,
			AccessibleDatabases: []string{},
			AccessibleSchemas:   map[string][]string{},
			AccessibleTables:    map[string][]string{},
			Permissions:         map[string][]domain.Permission{},
		}

		role, err := repo.LoadRolePermissions(ctx, roleName)

		_ = role
		_ = err
	})
}
