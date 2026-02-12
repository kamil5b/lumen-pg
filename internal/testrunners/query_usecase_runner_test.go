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

// QueryUseCaseConstructor creates a query use case with its dependencies
type QueryUseCaseConstructor func(repo repository.QueryRepository) usecase.QueryUseCase

// QueryUseCaseRunner runs test specs for query use case (Story 4)
func QueryUseCaseRunner(t *testing.T, constructor QueryUseCaseConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockQueryRepository(ctrl)
	useCase := constructor(mockRepo)

	t.Run("UC-S4-01: Single Query Execution", func(t *testing.T) {
		ctx := context.Background()
		sql := "SELECT * FROM users"
		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "username"},
			Rows:       [][]interface{}{{1, "test"}},
			TotalRows:  1,
			LoadedRows: 1,
			Success:    true,
		}

		mockRepo.EXPECT().ExecuteQuery(ctx, sql).Return(expectedResult, nil)

		result, err := useCase.ExecuteQuery(ctx, sql)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
		assert.Len(t, result.Rows, 1)
	})

	t.Run("UC-S4-02: Multiple Query Execution", func(t *testing.T) {
		ctx := context.Background()
		queries := "SELECT * FROM users; SELECT * FROM posts;"
		expectedResults := []*domain.QueryResult{
			{
				Columns:    []string{"id", "username"},
				Rows:       [][]interface{}{{1, "test"}},
				TotalRows:  1,
				LoadedRows: 1,
				Success:    true,
			},
			{
				Columns:    []string{"id", "title"},
				Rows:       [][]interface{}{{1, "post1"}},
				TotalRows:  1,
				LoadedRows: 1,
				Success:    true,
			},
		}

		mockRepo.EXPECT().ExecuteMultiple(ctx, queries).Return(expectedResults, nil)

		results, err := useCase.ExecuteMultipleQueries(ctx, queries)

		require.NoError(t, err)
		assert.Len(t, results, 2)
		assert.True(t, results[0].Success)
		assert.True(t, results[1].Success)
	})

	t.Run("UC-S4-03a: Query Result Actual Size Display", func(t *testing.T) {
		ctx := context.Background()
		sql := "SELECT * FROM large_table"
		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "data"},
			Rows:       make([][]interface{}, 1000), // First 1000 rows loaded
			TotalRows:  5000,                        // Total rows available
			LoadedRows: 1000,                        // Hard limit reached
			Success:    true,
		}

		mockRepo.EXPECT().ExecuteQuery(ctx, sql).Return(expectedResult, nil)

		result, err := useCase.ExecuteQuery(ctx, sql)

		require.NoError(t, err)
		assert.Equal(t, int64(5000), result.TotalRows) // Shows actual size
		assert.Equal(t, 1000, result.LoadedRows)       // Hard limit: only 1000 loaded
	})

	t.Run("UC-S4-06: Invalid Query Error", func(t *testing.T) {
		ctx := context.Background()
		sql := "SELECT * FROM nonexistent_table"
		expectedResult := &domain.QueryResult{
			Success:      false,
			ErrorMessage: "table does not exist",
		}

		mockRepo.EXPECT().ExecuteQuery(ctx, sql).Return(expectedResult, nil)

		result, err := useCase.ExecuteQuery(ctx, sql)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.False(t, result.Success)
		assert.Contains(t, result.ErrorMessage, "does not exist")
	})

	t.Run("UC-S4-08: Parameterized Query Execution", func(t *testing.T) {
		ctx := context.Background()
		sql := "SELECT * FROM users WHERE id = $1"
		params := []interface{}{1}
		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "username"},
			Rows:       [][]interface{}{{1, "test"}},
			TotalRows:  1,
			LoadedRows: 1,
			Success:    true,
		}

		mockRepo.EXPECT().ExecuteQuery(ctx, sql, 1).Return(expectedResult, nil)

		result, err := useCase.ExecuteQuery(ctx, sql, params...)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
	})
}
