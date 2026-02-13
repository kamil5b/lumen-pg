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

// DataViewUsecaseConstructor is a function type that creates a DataViewUsecase
type DataViewUsecaseConstructor func(
	metadataRepo repository.MetadataRepository,
	databaseRepo repository.DatabaseRepository,
	rbacRepo repository.RBACRepository,
) usecase.DataViewUseCase

// DataViewUsecaseRunner runs all DataView usecase tests against an implementation
// Maps to TEST_PLAN.md:
// - Story 5: Main View & Data Interaction [UC-S5-01~19, IT-S5-01~07, E2E-S5-01~15a]
func DataViewUsecaseRunner(t *testing.T, constructor DataViewUsecaseConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockMetadata := mockRepository.NewMockMetadataRepository(ctrl)
	mockDatabase := mockRepository.NewMockDatabaseRepository(ctrl)
	mockRBAC := mockRepository.NewMockRBACRepository(ctrl)

	uc := constructor(mockMetadata, mockDatabase, mockRBAC)

	// UC-S5-01: Table Data Loading
	// IT-S5-01: Real Table Data Loading
	// E2E-S5-01: Main View Default Load
	t.Run("LoadTableData returns first 50 rows", func(t *testing.T) {
		mockDatabase.EXPECT().
			GetTableData(gomock.Any(), gomock.Any()).
			Return(&domain.QueryResult{
				Columns:  []string{"id", "name"},
				Rows:     make([]map[string]interface{}, 10),
				RowCount: 10,
			}, nil)

		mockRBAC.EXPECT().
			CheckTableAccess(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(true, nil)

		result, err := uc.LoadTableData(ctx, "testuser", domain.TableDataParams{
			Database: "testdb",
			Schema:   "public",
			Table:    "users",
			Limit:    50,
			Offset:   0,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, 2, len(result.Columns))
		require.Equal(t, int64(10), result.RowCount)
	})

	// UC-S5-02: Cursor Pagination Next Page
	// IT-S5-02: Real Cursor Pagination
	// E2E-S5-05: Cursor Pagination Infinite Scroll with Actual Size
	// E2E-S5-05a: Cursor Pagination Infinite Scroll Loading
	t.Run("GetTableDataWithCursorPagination returns next page with cursor", func(t *testing.T) {
		mockDatabase.EXPECT().
			GetTableData(gomock.Any(), gomock.Any()).
			Return(&domain.QueryResult{
				Columns:    []string{"id", "name"},
				Rows:       make([]map[string]interface{}, 50),
				RowCount:   50,
				TotalCount: 100,
			}, nil)

		mockRBAC.EXPECT().
			CheckTableAccess(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(true, nil)

		result, err := uc.GetTableDataWithCursorPagination(ctx, "testuser", "testdb", "public", "users", "cursor_token", 50)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, int64(50), result.RowCount)
	})

	// UC-S5-03: WHERE Clause Validation
	// IT-S5-03: Real WHERE Filter
	// E2E-S5-03: WHERE Bar Filtering
	t.Run("FilterTableData with valid WHERE clause", func(t *testing.T) {
		mockDatabase.EXPECT().
			GetTableData(gomock.Any(), gomock.Any()).
			Return(&domain.QueryResult{
				Columns:  []string{"id", "name"},
				Rows:     make([]map[string]interface{}, 5),
				RowCount: 5,
			}, nil)

		mockRBAC.EXPECT().
			CheckTableAccess(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(true, nil)

		result, err := uc.FilterTableData(ctx, "testuser", "testdb", "public", "users", "id > 10", 0, 50)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, int64(5), result.RowCount)
	})

	// UC-S5-04: WHERE Clause Injection Prevention
	t.Run("ValidateWhereClause rejects SQL injection", func(t *testing.T) {
		valid, err := uc.ValidateWhereClause(ctx, "1' OR '1'='1")

		require.NoError(t, err)
		require.False(t, valid)
	})

	t.Run("ValidateWhereClause accepts safe clause", func(t *testing.T) {
		valid, err := uc.ValidateWhereClause(ctx, "id > 10 AND status = 'active'")

		require.NoError(t, err)
		require.True(t, valid)
	})

	// UC-S5-05: Column Sorting ASC
	// UC-S5-06: Column Sorting DESC
	// IT-S5-01: Real Table Data Loading (with sort)
	// E2E-S5-04: Column Header Sorting
	t.Run("SortTableData ascending", func(t *testing.T) {
		mockDatabase.EXPECT().
			GetTableData(gomock.Any(), gomock.Any()).
			Return(&domain.QueryResult{
				Columns:  []string{"id", "name"},
				Rows:     make([]map[string]interface{}, 20),
				RowCount: 20,
			}, nil)

		mockRBAC.EXPECT().
			CheckTableAccess(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(true, nil)

		result, err := uc.SortTableData(ctx, "testuser", "testdb", "public", "users", "name", "ASC", 0, 50)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, int64(20), result.RowCount)
	})

	t.Run("SortTableData descending", func(t *testing.T) {
		mockDatabase.EXPECT().
			GetTableData(gomock.Any(), gomock.Any()).
			Return(&domain.QueryResult{
				Columns:  []string{"id", "name"},
				Rows:     make([]map[string]interface{}, 20),
				RowCount: 20,
			}, nil)

		mockRBAC.EXPECT().
			CheckTableAccess(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(true, nil)

		result, err := uc.SortTableData(ctx, "testuser", "testdb", "public", "users", "name", "DESC", 0, 50)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, int64(20), result.RowCount)
	})

	// UC-S5-07: Cursor Pagination Actual Size Display
	// E2E-S5-05b: Pagination Hard Limit Enforcement
	t.Run("GetTableRowCount returns total count", func(t *testing.T) {
		mockDatabase.EXPECT().
			GetRowCount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(int64(5000), nil)

		mockRBAC.EXPECT().
			CheckTableAccess(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(true, nil)

		count, err := uc.GetTableRowCount(ctx, "testuser", "testdb", "public", "users")

		require.NoError(t, err)
		require.Equal(t, int64(5000), count)
	})

	// UC-S5-08: Cursor Pagination Hard Limit
	t.Run("GetTableRowCountWithFilter returns filtered count", func(t *testing.T) {
		mockDatabase.EXPECT().
			GetRowCount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(int64(1500), nil)

		mockRBAC.EXPECT().
			CheckTableAccess(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(true, nil)

		count, err := uc.GetTableRowCountWithFilter(ctx, "testuser", "testdb", "public", "users", "id > 100")

		require.NoError(t, err)
		require.Equal(t, int64(1500), count)
	})

	// UC-S5-17: Foreign Key Navigation
	// IT-S5-06: Real Foreign Key Navigation
	// E2E-S5-14: FK Cell Navigation (Read-Only)
	t.Run("GetForeignKeyInfo returns FK information", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetMetadata(gomock.Any(), gomock.Any()).
			Return(&domain.DatabaseMetadata{
				Name: "testdb",
				Schemas: []domain.SchemaMetadata{
					{
						Name: "public",
						Tables: []domain.TableMetadata{
							{
								Name: "posts",
								Columns: []domain.ColumnMetadata{
									{Name: "id", DataType: "integer", IsNullable: false, IsPrimary: true},
									{Name: "user_id", DataType: "integer", IsNullable: false, IsPrimary: false},
								},
								PrimaryKeys: []string{"id"},
								ForeignKeys: []domain.ForeignKeyMetadata{
									{
										ColumnName:         "user_id",
										ReferencedTable:    "users",
										ReferencedColumn:   "id",
										ReferencedSchema:   "public",
										ReferencedDatabase: "testdb",
									},
								},
							},
						},
					},
				},
			}, nil)

		mockRBAC.EXPECT().
			CheckTableAccess(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(true, nil)

		fks, err := uc.GetForeignKeyInfo(ctx, "testuser", "testdb", "public", "posts")

		require.NoError(t, err)
		require.NotNil(t, fks)
		require.GreaterOrEqual(t, len(fks), 1)
	})

	// UC-S5-18: Primary Key Navigation
	// IT-S5-07: Real Primary Key Navigation
	// E2E-S5-15: PK Cell Navigation (Read-Only)
	// E2E-S5-15a: PK Cell Navigation - Table Click
	t.Run("GetPrimaryKeyInfo returns PK information", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetMetadata(gomock.Any(), gomock.Any()).
			Return(&domain.DatabaseMetadata{
				Name: "testdb",
				Schemas: []domain.SchemaMetadata{
					{
						Name: "public",
						Tables: []domain.TableMetadata{
							{
								Name: "users",
								Columns: []domain.ColumnMetadata{
									{Name: "id", DataType: "integer", IsNullable: false, IsPrimary: true},
									{Name: "name", DataType: "text", IsNullable: true, IsPrimary: false},
								},
								PrimaryKeys: []string{"id"},
								ForeignKeys: []domain.ForeignKeyMetadata{},
							},
						},
					},
				},
			}, nil)

		mockRBAC.EXPECT().
			CheckTableAccess(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(true, nil)

		pks, err := uc.GetPrimaryKeyInfo(ctx, "testuser", "testdb", "public", "users")

		require.NoError(t, err)
		require.NotNil(t, pks)
		require.Equal(t, 1, len(pks))
		require.Equal(t, "id", pks[0])
	})

	t.Run("GetChildTableReferences returns child table references", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetMetadata(gomock.Any(), gomock.Any()).
			Return(&domain.DatabaseMetadata{
				Name: "testdb",
				Schemas: []domain.SchemaMetadata{
					{
						Name: "public",
						Tables: []domain.TableMetadata{
							{
								Name: "posts",
								Columns: []domain.ColumnMetadata{
									{Name: "id", DataType: "integer", IsNullable: false, IsPrimary: true},
									{Name: "user_id", DataType: "integer", IsNullable: false, IsPrimary: false},
								},
								PrimaryKeys: []string{"id"},
								ForeignKeys: []domain.ForeignKeyMetadata{
									{
										ColumnName:         "user_id",
										ReferencedTable:    "users",
										ReferencedColumn:   "id",
										ReferencedSchema:   "public",
										ReferencedDatabase: "testdb",
									},
								},
							},
						},
					},
				},
			}, nil)

		mockRBAC.EXPECT().
			CheckTableAccess(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(true, nil)

		mockDatabase.EXPECT().
			GetRowCount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(int64(15), nil)

		refs, err := uc.GetChildTableReferences(ctx, "testuser", "testdb", "public", "users", map[string]interface{}{"id": 1})

		require.NoError(t, err)
		require.NotNil(t, refs)
	})

	t.Run("GetChildTableRowCount returns child row count", func(t *testing.T) {
		mockDatabase.EXPECT().
			GetRowCount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(int64(10), nil)

		mockRBAC.EXPECT().
			CheckTableAccess(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(true, nil).Times(2)

		count, err := uc.GetChildTableRowCount(ctx, "testuser", "testdb", "public", "posts", "users", "user_id", "1")

		require.NoError(t, err)
		require.Equal(t, int64(10), count)
	})

	t.Run("NavigateToParentRow returns parent row data", func(t *testing.T) {
		mockDatabase.EXPECT().
			GetTableData(gomock.Any(), gomock.Any()).
			Return(&domain.QueryResult{
				Columns:  []string{"id", "name"},
				Rows:     []map[string]interface{}{{"id": 1, "name": "John"}},
				RowCount: 1,
			}, nil)

		mockRBAC.EXPECT().
			CheckTableAccess(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(true, nil)

		mockMetadata.EXPECT().
			GetMetadata(gomock.Any(), gomock.Any()).
			Return(&domain.DatabaseMetadata{
				Name: "testdb",
				Schemas: []domain.SchemaMetadata{
					{
						Name: "public",
						Tables: []domain.TableMetadata{
							{
								Name: "posts",
								Columns: []domain.ColumnMetadata{
									{Name: "id", DataType: "integer", IsNullable: false, IsPrimary: true},
									{Name: "user_id", DataType: "integer", IsNullable: false, IsPrimary: false},
								},
								PrimaryKeys: []string{"id"},
								ForeignKeys: []domain.ForeignKeyMetadata{
									{
										ColumnName:         "user_id",
										ReferencedTable:    "users",
										ReferencedColumn:   "id",
										ReferencedSchema:   "public",
										ReferencedDatabase: "testdb",
									},
								},
							},
						},
					},
				},
			}, nil)

		result, err := uc.NavigateToParentRow(ctx, "testuser", "testdb", "public", "posts", "user_id", 1)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, int64(1), result.RowCount)
	})

	t.Run("NavigateToChildRows returns child rows data", func(t *testing.T) {
		mockDatabase.EXPECT().
			GetTableData(gomock.Any(), gomock.Any()).
			Return(&domain.QueryResult{
				Columns:  []string{"id", "user_id", "title"},
				Rows:     make([]map[string]interface{}, 5),
				RowCount: 5,
			}, nil)

		mockRBAC.EXPECT().
			CheckTableAccess(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(true, nil)

		result, err := uc.NavigateToChildRows(ctx, "testuser", "testdb", "public", "posts", "users", "user_id", "id", "1")

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, int64(5), result.RowCount)
	})

	// UC-S5-19: Read-Only Mode Enforcement
	t.Run("IsTableReadOnly returns true when no write permissions", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetMetadata(gomock.Any(), gomock.Any()).
			Return(&domain.DatabaseMetadata{
				Name: "testdb",
				Schemas: []domain.SchemaMetadata{
					{
						Name: "public",
						Tables: []domain.TableMetadata{
							{
								Name:        "users",
								Columns:     []domain.ColumnMetadata{},
								PrimaryKeys: []string{},
								ForeignKeys: []domain.ForeignKeyMetadata{},
							},
						},
					},
				},
			}, nil)

		mockRBAC.EXPECT().
			CheckSelectPermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(true, nil)

		mockRBAC.EXPECT().
			CheckInsertPermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(false, nil)

		mockRBAC.EXPECT().
			CheckUpdatePermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(false, nil)

		mockRBAC.EXPECT().
			CheckDeletePermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(false, nil)

		readonly, err := uc.IsTableReadOnly(ctx, "testuser", "testdb", "public", "users")

		require.NoError(t, err)
		require.True(t, readonly)
	})

	t.Run("IsTableReadOnly returns false when has write permissions", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetMetadata(gomock.Any(), gomock.Any()).
			Return(&domain.DatabaseMetadata{
				Name: "testdb",
				Schemas: []domain.SchemaMetadata{
					{
						Name: "public",
						Tables: []domain.TableMetadata{
							{
								Name:        "users",
								Columns:     []domain.ColumnMetadata{},
								PrimaryKeys: []string{},
								ForeignKeys: []domain.ForeignKeyMetadata{},
							},
						},
					},
				},
			}, nil)

		mockRBAC.EXPECT().
			CheckSelectPermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(true, nil)

		mockRBAC.EXPECT().
			CheckInsertPermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(true, nil)

		readonly, err := uc.IsTableReadOnly(ctx, "testuser", "testdb", "public", "users")

		require.NoError(t, err)
		require.False(t, readonly)
	})
}
