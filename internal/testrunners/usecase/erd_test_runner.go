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

// ERDUsecaseConstructor is a function type that creates an ERDUseCase
type ERDUsecaseConstructor func(
	metadataRepo repository.MetadataRepository,
	rbacRepo repository.RBACRepository,
) usecase.ERDUseCase

// ERDUsecaseRunner runs all ERD usecase tests against an implementation
// Maps to TEST_PLAN.md:
// - Story 3: ERD Viewer [UC-S3-01~04, IT-S3-01~02, E2E-S3-01~04]
func ERDUsecaseRunner(t *testing.T, constructor ERDUsecaseConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMetadata := mockRepository.NewMockMetadataRepository(ctrl)
	mockRBAC := mockRepository.NewMockRBACRepository(ctrl)

	uc := constructor(mockMetadata, mockRBAC)

	ctx := context.Background()

	// UC-S3-01: ERD Data Generation
	// IT-S3-01: ERD from Real Schema
	// E2E-S3-01: ERD Viewer Page Access
	t.Run("GenerateERD creates ERD data for schema", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetMetadata(gomock.Any(), "testdb").
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
							{
								Name: "posts",
								Columns: []domain.ColumnMetadata{
									{Name: "id", DataType: "integer", IsNullable: false, IsPrimary: true},
									{Name: "user_id", DataType: "integer", IsNullable: false, IsPrimary: false},
									{Name: "title", DataType: "text", IsNullable: true, IsPrimary: false},
								},
								PrimaryKeys: []string{"id"},
								ForeignKeys: []domain.ForeignKeyMetadata{
									{
										ColumnName:       "user_id",
										ReferencedTable:  "users",
										ReferencedSchema: "public",
									},
								},
							},
						},
					},
				},
			}, nil)

		mockRBAC.EXPECT().
			CanAccessTable(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(true, nil).
			Times(2)

		erdData, err := uc.GenerateERD(ctx, "testuser", "testdb", "public")

		require.NoError(t, err)
		require.NotNil(t, erdData)
	})

	// UC-S3-02: Table Box Representation
	// IT-S3-01: ERD from Real Schema
	t.Run("GetERDTables returns all tables in schema", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetMetadata(gomock.Any(), "testdb").
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
							{
								Name: "posts",
								Columns: []domain.ColumnMetadata{
									{Name: "id", DataType: "integer", IsNullable: false, IsPrimary: true},
									{Name: "user_id", DataType: "integer", IsNullable: false, IsPrimary: false},
								},
								PrimaryKeys: []string{"id"},
								ForeignKeys: []domain.ForeignKeyMetadata{},
							},
						},
					},
				},
			}, nil)

		mockRBAC.EXPECT().
			CanAccessTable(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(true, nil).
			Times(2)

		tables, err := uc.GetERDTables(ctx, "testuser", "testdb", "public")

		require.NoError(t, err)
		require.NotNil(t, tables)
		require.GreaterOrEqual(t, len(tables), 2)
	})

	// UC-S3-02: Table Box Representation
	t.Run("GetTableBoxData returns table metadata for ERD", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetMetadata(gomock.Any(), "testdb").
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
									{Name: "email", DataType: "varchar", IsNullable: true, IsPrimary: false},
								},
								PrimaryKeys: []string{"id"},
								ForeignKeys: []domain.ForeignKeyMetadata{},
							},
						},
					},
				},
			}, nil)

		mockRBAC.EXPECT().
			CanAccessTable(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(true, nil).
			Times(1)

		tableData, err := uc.GetTableBoxData(ctx, "testuser", "testdb", "public", "users")

		require.NoError(t, err)
		require.NotNil(t, tableData)
		require.Equal(t, "users", tableData.Name)
		require.Equal(t, 3, len(tableData.Columns))
	})

	// UC-S3-03: Relationship Lines
	// IT-S3-02: Complex Relationships
	t.Run("GetERDRelationships returns foreign key relationships", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetMetadata(gomock.Any(), "testdb").
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
									{Name: "category_id", DataType: "integer", IsNullable: true, IsPrimary: false},
								},
								PrimaryKeys: []string{"id"},
								ForeignKeys: []domain.ForeignKeyMetadata{
									{
										ColumnName:       "user_id",
										ReferencedTable:  "users",
										ReferencedSchema: "public",
									},
									{
										ColumnName:       "category_id",
										ReferencedTable:  "categories",
										ReferencedSchema: "public",
									},
								},
							},
						},
					},
				},
			}, nil)

		mockRBAC.EXPECT().
			CanAccessTable(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(true, nil).
			Times(1)

		relationships, err := uc.GetERDRelationships(ctx, "testuser", "testdb", "public")

		require.NoError(t, err)
		require.NotNil(t, relationships)
		require.GreaterOrEqual(t, len(relationships), 2)
	})

	// UC-S3-03: Relationship Lines
	t.Run("GetRelationshipLines returns relationship line data", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetMetadata(gomock.Any(), "testdb").
			Return(&domain.DatabaseMetadata{
				Name: "testdb",
				Schemas: []domain.SchemaMetadata{
					{
						Name: "public",
						Tables: []domain.TableMetadata{
							{
								Name: "users",
								Columns: []domain.ColumnMetadata{
									{Name: "id", DataType: "integer"},
								},
								PrimaryKeys: []string{"id"},
								ForeignKeys: []domain.ForeignKeyMetadata{},
							},
							{
								Name: "posts",
								ForeignKeys: []domain.ForeignKeyMetadata{
									{
										ColumnName:       "user_id",
										ReferencedTable:  "users",
										ReferencedSchema: "public",
									},
								},
							},
						},
					},
				},
			}, nil)

		mockRBAC.EXPECT().
			CanAccessTable(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(true, nil).
			Times(2)

		lines, err := uc.GetRelationshipLines(ctx, "testuser", "testdb", "public")

		require.NoError(t, err)
		require.NotNil(t, lines)
	})

	// UC-S3-04: Empty Schema ERD
	t.Run("IsSchemaEmpty returns true for empty schema", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetMetadata(gomock.Any(), "testdb").
			Return(&domain.DatabaseMetadata{
				Name: "testdb",
				Schemas: []domain.SchemaMetadata{
					{
						Name:   "empty_schema",
						Tables: []domain.TableMetadata{},
					},
				},
			}, nil)

		isEmpty, err := uc.IsSchemaEmpty(ctx, "testuser", "testdb", "empty_schema")

		require.NoError(t, err)
		require.True(t, isEmpty)
	})

	t.Run("IsSchemaEmpty returns false for non-empty schema", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetMetadata(gomock.Any(), "testdb").
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

		isEmpty, err := uc.IsSchemaEmpty(ctx, "testuser", "testdb", "public")

		require.NoError(t, err)
		require.False(t, isEmpty)
	})

	// E2E-S3-02: ERD Zoom Controls
	// E2E-S3-03: ERD Pan
	// E2E-S3-04: Table Click in ERD
	t.Run("GetAvailableSchemas returns all accessible schemas", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetMetadata(gomock.Any(), "testdb").
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
					{
						Name: "private",
						Tables: []domain.TableMetadata{
							{
								Name:        "admin_users",
								Columns:     []domain.ColumnMetadata{},
								PrimaryKeys: []string{},
								ForeignKeys: []domain.ForeignKeyMetadata{},
							},
						},
					},
				},
			}, nil)

		mockRBAC.EXPECT().
			HasSchemaUsagePermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(true, nil).
			Times(2)

		schemas, err := uc.GetAvailableSchemas(ctx, "testuser", "testdb")

		require.NoError(t, err)
		require.NotNil(t, schemas)
		require.GreaterOrEqual(t, len(schemas), 1)
	})

	// E2E-S3-01: ERD Viewer Page Access
	t.Run("GenerateERD handles complex schema with multiple relationships", func(t *testing.T) {
		mockMetadata.EXPECT().
			GetMetadata(gomock.Any(), "testdb").
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
							{
								Name: "posts",
								Columns: []domain.ColumnMetadata{
									{Name: "id", DataType: "integer", IsNullable: false, IsPrimary: true},
									{Name: "user_id", DataType: "integer", IsNullable: false, IsPrimary: false},
									{Name: "category_id", DataType: "integer", IsNullable: true, IsPrimary: false},
								},
								PrimaryKeys: []string{"id"},
								ForeignKeys: []domain.ForeignKeyMetadata{
									{
										ColumnName:       "user_id",
										ReferencedTable:  "users",
										ReferencedSchema: "public",
									},
									{
										ColumnName:       "category_id",
										ReferencedTable:  "categories",
										ReferencedSchema: "public",
									},
								},
							},
							{
								Name: "categories",
								Columns: []domain.ColumnMetadata{
									{Name: "id", DataType: "integer", IsNullable: false, IsPrimary: true},
									{Name: "name", DataType: "text", IsNullable: false, IsPrimary: false},
								},
								PrimaryKeys: []string{"id"},
								ForeignKeys: []domain.ForeignKeyMetadata{},
							},
							{
								Name: "comments",
								Columns: []domain.ColumnMetadata{
									{Name: "id", DataType: "integer", IsNullable: false, IsPrimary: true},
									{Name: "post_id", DataType: "integer", IsNullable: false, IsPrimary: false},
									{Name: "user_id", DataType: "integer", IsNullable: false, IsPrimary: false},
									{Name: "content", DataType: "text", IsNullable: true, IsPrimary: false},
								},
								PrimaryKeys: []string{"id"},
								ForeignKeys: []domain.ForeignKeyMetadata{
									{
										ColumnName:       "post_id",
										ReferencedTable:  "posts",
										ReferencedSchema: "public",
									},
									{
										ColumnName:       "user_id",
										ReferencedTable:  "users",
										ReferencedSchema: "public",
									},
								},
							},
						},
					},
				},
			}, nil)

		mockRBAC.EXPECT().
			CanAccessTable(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(true, nil).
			Times(4)

		erdData, err := uc.GenerateERD(ctx, "testuser", "testdb", "public")

		require.NoError(t, err)
		require.NotNil(t, erdData)
	})
}
