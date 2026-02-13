package usecase

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
		query := "SELECT * FROM users"
		expectedResult := &domain.QueryResult{
			Success: true,
			Columns: []string{"id", "username", "email"},
			Rows: [][]interface{}{
				{1, "user1", "user1@test.com"},
			},
			AffectedRows: 0,
		}

		mockRepo.EXPECT().ExecuteQuery(ctx, query).Return(expectedResult, nil)

		result, err := useCase.ExecuteQuery(ctx, query)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
		assert.Len(t, result.Columns, 3)
		assert.Len(t, result.Rows, 1)
	})

	t.Run("UC-S4-02: Multiple Query Execution", func(t *testing.T) {
		ctx := context.Background()
		queries := "SELECT * FROM users; SELECT * FROM posts"
		expectedResults := []*domain.QueryResult{
			{
				Success: true,
				Columns: []string{"id", "username"},
				Rows: [][]interface{}{
					{1, "user1"},
				},
			},
			{
				Success: true,
				Columns: []string{"id", "title"},
				Rows: [][]interface{}{
					{1, "post1"},
				},
			},
		}

		mockRepo.EXPECT().ExecuteMultiple(ctx, queries).Return(expectedResults, nil)

		results, err := useCase.ExecuteMultipleQueries(ctx, queries)

		require.NoError(t, err)
		assert.NotNil(t, results)
		assert.Len(t, results, 2)
		assert.True(t, results[0].Success)
		assert.True(t, results[1].Success)
	})

	t.Run("UC-S4-03: Query Result Offset Pagination", func(t *testing.T) {
		ctx := context.Background()
		query := "SELECT * FROM users LIMIT 10 OFFSET 0"
		expectedResult := &domain.QueryResult{
			Success:    true,
			Columns:    []string{"id", "username"},
			Rows:       make([][]interface{}, 10),
			TotalRows:  100, // Total rows available
			LoadedRows: 10,
		}

		mockRepo.EXPECT().ExecuteQuery(ctx, query).Return(expectedResult, nil)

		result, err := useCase.ExecuteQuery(ctx, query)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(100), result.TotalRows)
	})

	t.Run("UC-S4-03a: Query Result Actual Size Display", func(t *testing.T) {
		ctx := context.Background()
		query := "SELECT COUNT(*) as total FROM users"
		expectedResult := &domain.QueryResult{
			Success: true,
			Columns: []string{"total"},
			Rows: [][]interface{}{
				{int64(1000)},
			},
		}

		mockRepo.EXPECT().ExecuteQuery(ctx, query).Return(expectedResult, nil)

		result, err := useCase.ExecuteQuery(ctx, query)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Rows, 1)
		assert.Equal(t, int64(1000), result.Rows[0][0])
	})

	t.Run("UC-S4-03b: Query Result Limit Hard Cap", func(t *testing.T) {
		ctx := context.Background()
		query := "SELECT * FROM users"
		// Result capped at hard limit (e.g., 1000 rows max display)
		expectedResult := &domain.QueryResult{
			Success:    true,
			Columns:    []string{"id", "username"},
			Rows:       make([][]interface{}, 1000), // Hard cap
			TotalRows:  50000,                       // But actual data is larger
			LoadedRows: 1000,                        // Hard limit enforced
		}

		mockRepo.EXPECT().ExecuteQuery(ctx, query).Return(expectedResult, nil)

		result, err := useCase.ExecuteQuery(ctx, query)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Rows, 1000)
		assert.Equal(t, int64(50000), result.TotalRows)
		assert.Equal(t, 1000, result.LoadedRows)
	})

	t.Run("UC-S4-03c: Offset Pagination Next Page", func(t *testing.T) {
		ctx := context.Background()
		query := "SELECT * FROM users LIMIT 10 OFFSET 10"
		expectedResult := &domain.QueryResult{
			Success:    true,
			Columns:    []string{"id", "username"},
			Rows:       make([][]interface{}, 10),
			LoadedRows: 10,
		}

		mockRepo.EXPECT().ExecuteQuery(ctx, query).Return(expectedResult, nil)

		result, err := useCase.ExecuteQuery(ctx, query)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Rows, 10)
	})

	t.Run("UC-S4-04: DDL Query Execution", func(t *testing.T) {
		ctx := context.Background()
		query := "CREATE TABLE test (id SERIAL PRIMARY KEY, name VARCHAR(100))"

		mockRepo.EXPECT().ExecuteDDL(ctx, query).Return(nil)

		err := mockRepo.ExecuteDDL(ctx, query)

		require.NoError(t, err)
	})

	t.Run("UC-S4-05: DML Query Execution", func(t *testing.T) {
		ctx := context.Background()
		query := "INSERT INTO users (username, email) VALUES ($1, $2)"
		expectedResult := &domain.QueryResult{
			Success:      true,
			AffectedRows: 1,
			LoadedRows:   0,
		}

		mockRepo.EXPECT().ExecuteDML(ctx, query, "newuser", "new@test.com").Return(expectedResult, nil)

		result, err := mockRepo.ExecuteDML(ctx, query, "newuser", "new@test.com")

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
		assert.Equal(t, int64(1), result.AffectedRows)
	})

	t.Run("UC-S4-06: Invalid Query Error", func(t *testing.T) {
		ctx := context.Background()
		query := "SELECT * FROM nonexistent_table"
		expectedResult := &domain.QueryResult{
			Success:      false,
			ErrorMessage: "relation \"nonexistent_table\" does not exist",
		}

		mockRepo.EXPECT().ExecuteQuery(ctx, query).Return(expectedResult, nil)

		result, err := useCase.ExecuteQuery(ctx, query)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.False(t, result.Success)
		assert.NotEmpty(t, result.ErrorMessage)
	})

	t.Run("UC-S4-07: Query Splitting", func(t *testing.T) {
		ctx := context.Background()
		queries := "SELECT * FROM users; DELETE FROM posts; UPDATE comments SET approved = true"
		expectedResults := []*domain.QueryResult{
			{Success: true},
			{Success: true, AffectedRows: 5},
			{Success: true, AffectedRows: 3},
		}

		mockRepo.EXPECT().ExecuteMultiple(ctx, queries).Return(expectedResults, nil)

		results, err := useCase.ExecuteMultipleQueries(ctx, queries)

		require.NoError(t, err)
		assert.NotNil(t, results)
		assert.Len(t, results, 3)
		assert.True(t, results[0].Success)
		assert.True(t, results[1].Success)
		assert.True(t, results[2].Success)
	})

	t.Run("UC-S4-08: Parameterized Query Execution", func(t *testing.T) {
		ctx := context.Background()
		query := "SELECT * FROM users WHERE id = $1 AND username = $2"
		expectedResult := &domain.QueryResult{
			Success:    true,
			Columns:    []string{"id", "username", "email"},
			Rows:       [][]interface{}{{1, "user1", "user1@test.com"}},
			LoadedRows: 1,
		}

		mockRepo.EXPECT().ExecuteQuery(ctx, query, int64(1), "user1").Return(expectedResult, nil)

		result, err := useCase.ExecuteQuery(ctx, query, int64(1), "user1")

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
		assert.Len(t, result.Rows, 1)
		assert.Equal(t, "user1", result.Rows[0][1])
	})
}
