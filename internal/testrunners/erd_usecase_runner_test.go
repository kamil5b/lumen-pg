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

// ERDUseCaseConstructor creates an ERD use case with its dependencies
type ERDUseCaseConstructor func(metadataRepo repository.MetadataRepository) usecase.MetadataUseCase

// ERDUseCaseRunner runs test specs for ERD use case (Story 3)
func ERDUseCaseRunner(t *testing.T, constructor ERDUseCaseConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMetadataRepo := mocks.NewMockMetadataRepository(ctrl)
	useCase := constructor(mockMetadataRepo)

	t.Run("UC-S3-01: ERD Data Generation", func(t *testing.T) {
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
									PrimaryKeys: []string{"id"},
								},
								{
									SchemaName: "public",
									TableName:  "posts",
									Columns: []domain.ColumnMetadata{
										{Name: "id", DataType: "integer"},
										{Name: "user_id", DataType: "integer"},
										{Name: "title", DataType: "varchar"},
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
							},
						},
					},
				},
			},
		}

		mockMetadataRepo.EXPECT().LoadGlobalMetadata(ctx).Return(expectedMetadata, nil)

		result, err := useCase.LoadGlobalMetadata(ctx)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.Databases)
		assert.NotEmpty(t, result.Databases[0].Schemas[0].Tables)
	})

	t.Run("UC-S3-02: Table Box Representation", func(t *testing.T) {
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
										{Name: "id", DataType: "integer", IsNullable: false},
										{Name: "username", DataType: "varchar", IsNullable: true},
										{Name: "email", DataType: "varchar", IsNullable: true},
									},
									PrimaryKeys: []string{"id"},
								},
							},
						},
					},
				},
			},
		}

		mockMetadataRepo.EXPECT().LoadGlobalMetadata(ctx).Return(expectedMetadata, nil)

		result, err := useCase.LoadGlobalMetadata(ctx)

		require.NoError(t, err)
		assert.NotNil(t, result)
		table := result.Databases[0].Schemas[0].Tables[0]
		assert.Equal(t, "users", table.TableName)
		assert.Len(t, table.Columns, 3)
		assert.Contains(t, table.PrimaryKeys, "id")
	})

	t.Run("UC-S3-03: Relationship Lines", func(t *testing.T) {
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

		mockMetadataRepo.EXPECT().LoadGlobalMetadata(ctx).Return(expectedMetadata, nil)

		result, err := useCase.LoadGlobalMetadata(ctx)

		require.NoError(t, err)
		assert.NotNil(t, result)

		// Verify relationships exist
		postsTable := result.Databases[0].Schemas[0].Tables[1]
		assert.NotEmpty(t, postsTable.ForeignKeys)
		assert.Equal(t, "users", postsTable.ForeignKeys[0].ReferencedTableName)

		commentsTable := result.Databases[0].Schemas[0].Tables[2]
		assert.Len(t, commentsTable.ForeignKeys, 2)
	})

	t.Run("UC-S3-04: Empty Schema ERD", func(t *testing.T) {
		ctx := context.Background()
		expectedMetadata := &domain.GlobalMetadata{
			Databases: []domain.DatabaseMetadata{
				{
					Name: "emptydb",
					Schemas: []domain.SchemaMetadata{
						{
							Name:   "public",
							Tables: []domain.TableMetadata{},
						},
					},
				},
			},
		}

		mockMetadataRepo.EXPECT().LoadGlobalMetadata(ctx).Return(expectedMetadata, nil)

		result, err := useCase.LoadGlobalMetadata(ctx)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result.Databases[0].Schemas[0].Tables)
	})
}
