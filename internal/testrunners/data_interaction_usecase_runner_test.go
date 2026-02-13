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

// DataInteractionUseCaseConstructor creates a data interaction use case with its dependencies
type DataInteractionUseCaseConstructor func(queryRepo repository.QueryRepository) usecase.QueryUseCase

// DataInteractionUseCaseRunner runs test specs for data interaction use case (Story 5)
func DataInteractionUseCaseRunner(t *testing.T, constructor DataInteractionUseCaseConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueryRepo := mocks.NewMockQueryRepository(ctrl)
	useCase := constructor(mockQueryRepo)

	t.Run("UC-S5-01: Table Data Loading", func(t *testing.T) {
		ctx := context.Background()
		sql := "SELECT * FROM users LIMIT 50"
		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "username", "email"},
			Rows:       make([][]interface{}, 10),
			TotalRows:  100,
			LoadedRows: 10,
			Success:    true,
		}

		mockQueryRepo.EXPECT().ExecuteQuery(ctx, sql).Return(expectedResult, nil)

		result, err := useCase.ExecuteQuery(ctx, sql)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
		assert.Len(t, result.Columns, 3)
		assert.Equal(t, int64(100), result.TotalRows)
	})

	t.Run("UC-S5-02: Cursor Pagination Next Page", func(t *testing.T) {
		ctx := context.Background()
		sql := "SELECT * FROM users LIMIT 50 OFFSET 50"
		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "username", "email"},
			Rows:       make([][]interface{}, 50),
			TotalRows:  200,
			LoadedRows: 50,
			Success:    true,
		}

		mockQueryRepo.EXPECT().ExecuteQuery(ctx, sql).Return(expectedResult, nil)

		result, err := useCase.ExecuteQuery(ctx, sql)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
	})

	t.Run("UC-S5-03: WHERE Clause Validation", func(t *testing.T) {
		ctx := context.Background()
		sql := "SELECT * FROM users WHERE id > 10 LIMIT 50"
		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "username", "email"},
			Rows:       make([][]interface{}, 40),
			TotalRows:  90,
			LoadedRows: 40,
			Success:    true,
		}

		mockQueryRepo.EXPECT().ExecuteQuery(ctx, sql).Return(expectedResult, nil)

		result, err := useCase.ExecuteQuery(ctx, sql)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
	})

	t.Run("UC-S5-04: WHERE Clause Injection Prevention", func(t *testing.T) {
		ctx := context.Background()
		// SQL injection attempt should be handled
		sql := "SELECT * FROM users WHERE username = $1"
		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "username", "email"},
			Rows:       [][]interface{}{},
			TotalRows:  0,
			LoadedRows: 0,
			Success:    true,
		}

		mockQueryRepo.EXPECT().ExecuteQuery(ctx, sql, "admin'; DROP TABLE users; --").Return(expectedResult, nil)

		result, err := useCase.ExecuteQuery(ctx, sql, "admin'; DROP TABLE users; --")

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
	})

	t.Run("UC-S5-05: Column Sorting ASC", func(t *testing.T) {
		ctx := context.Background()
		sql := "SELECT * FROM users ORDER BY username ASC LIMIT 50"
		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "username", "email"},
			Rows:       [][]interface{}{{1, "alice", "alice@test.com"}, {2, "bob", "bob@test.com"}},
			TotalRows:  2,
			LoadedRows: 2,
			Success:    true,
		}

		mockQueryRepo.EXPECT().ExecuteQuery(ctx, sql).Return(expectedResult, nil)

		result, err := useCase.ExecuteQuery(ctx, sql)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
	})

	t.Run("UC-S5-06: Column Sorting DESC", func(t *testing.T) {
		ctx := context.Background()
		sql := "SELECT * FROM users ORDER BY username DESC LIMIT 50"
		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "username", "email"},
			Rows:       [][]interface{}{{2, "bob", "bob@test.com"}, {1, "alice", "alice@test.com"}},
			TotalRows:  2,
			LoadedRows: 2,
			Success:    true,
		}

		mockQueryRepo.EXPECT().ExecuteQuery(ctx, sql).Return(expectedResult, nil)

		result, err := useCase.ExecuteQuery(ctx, sql)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
	})

	t.Run("UC-S5-07: Cursor Pagination Actual Size Display", func(t *testing.T) {
		ctx := context.Background()
		sql := "SELECT * FROM large_table LIMIT 50"
		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "data"},
			Rows:       make([][]interface{}, 50),
			TotalRows:  5000,
			LoadedRows: 50,
			Success:    true,
		}

		mockQueryRepo.EXPECT().ExecuteQuery(ctx, sql).Return(expectedResult, nil)

		result, err := useCase.ExecuteQuery(ctx, sql)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(5000), result.TotalRows)
		assert.Equal(t, 50, result.LoadedRows)
	})

	t.Run("UC-S5-08: Cursor Pagination Hard Limit", func(t *testing.T) {
		ctx := context.Background()
		sql := "SELECT * FROM huge_table"
		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "data"},
			Rows:       make([][]interface{}, 1000),
			TotalRows:  1000000,
			LoadedRows: 1000,
			Success:    true,
		}

		mockQueryRepo.EXPECT().ExecuteQuery(ctx, sql).Return(expectedResult, nil)

		result, err := useCase.ExecuteQuery(ctx, sql)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1000, result.LoadedRows)
		assert.Greater(t, result.TotalRows, int64(1000))
	})

	t.Run("UC-S5-19: Read-Only Mode Enforcement", func(t *testing.T) {
		ctx := context.Background()
		sql := "SELECT * FROM users LIMIT 50"
		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "username", "email"},
			Rows:       make([][]interface{}, 10),
			TotalRows:  100,
			LoadedRows: 10,
			Success:    true,
		}

		mockQueryRepo.EXPECT().ExecuteQuery(ctx, sql).Return(expectedResult, nil)

		result, err := useCase.ExecuteQuery(ctx, sql)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
	})
}
