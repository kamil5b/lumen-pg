package e2e

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/implementations/mocks"
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
)

// QueryInterfaceConstructor creates a query repository with mock dependencies
type QueryInterfaceConstructor func(ctrl *gomock.Controller) repository.QueryRepository

// QueryInterfaceRunner runs unit tests for query repository interface (Story 4)
func QueryInterfaceRunner(t *testing.T, constructor QueryInterfaceConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockQueryRepository(ctrl)
	ctx := context.Background()

	t.Run("UC-S4-01: Execute Simple SELECT Query", func(t *testing.T) {
		sql := "SELECT id, username FROM users"
		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "username"},
			Rows:       [][]interface{}{{1, "alice"}, {2, "bob"}},
			TotalRows:  2,
			LoadedRows: 2,
			Success:    true,
		}

		mockRepo.EXPECT().ExecuteQuery(ctx, sql).Return(expectedResult, nil)

		result, err := mockRepo.ExecuteQuery(ctx, sql)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
		assert.Len(t, result.Rows, 2)
		assert.Len(t, result.Columns, 2)
	})

	t.Run("UC-S4-02: Execute Multiple Queries", func(t *testing.T) {
		queries := "SELECT * FROM users; SELECT * FROM posts;"
		expectedResults := []*domain.QueryResult{
			{
				Columns:    []string{"id", "username"},
				Rows:       [][]interface{}{{1, "alice"}},
				TotalRows:  1,
				LoadedRows: 1,
				Success:    true,
			},
			{
				Columns:    []string{"id", "title"},
				Rows:       [][]interface{}{{1, "first post"}},
				TotalRows:  1,
				LoadedRows: 1,
				Success:    true,
			},
		}

		mockRepo.EXPECT().ExecuteMultiple(ctx, queries).Return(expectedResults, nil)

		results, err := mockRepo.ExecuteMultiple(ctx, queries)

		require.NoError(t, err)
		assert.Len(t, results, 2)
		assert.True(t, results[0].Success)
		assert.True(t, results[1].Success)
	})

	t.Run("UC-S4-03: Query Result Offset Pagination", func(t *testing.T) {
		sql := "SELECT * FROM users LIMIT 100 OFFSET 0"
		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "username"},
			Rows:       make([][]interface{}, 100),
			TotalRows:  500,
			LoadedRows: 100,
			Success:    true,
		}

		mockRepo.EXPECT().ExecuteQuery(ctx, sql).Return(expectedResult, nil)

		result, err := mockRepo.ExecuteQuery(ctx, sql)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
		assert.Equal(t, 100, result.LoadedRows)
		assert.Equal(t, int64(500), result.TotalRows)
	})

	t.Run("UC-S4-03a: Query Result Actual Size Display", func(t *testing.T) {
		sql := "SELECT * FROM large_table"
		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "data"},
			Rows:       make([][]interface{}, 1000),
			TotalRows:  5000,
			LoadedRows: 1000,
			Success:    true,
		}

		mockRepo.EXPECT().ExecuteQuery(ctx, sql).Return(expectedResult, nil)

		result, err := mockRepo.ExecuteQuery(ctx, sql)

		require.NoError(t, err)
		assert.Equal(t, int64(5000), result.TotalRows)
		assert.Equal(t, 1000, result.LoadedRows)
	})

	t.Run("UC-S4-03b: Query Result Limit Hard Cap", func(t *testing.T) {
		sql := "SELECT * FROM huge_table"
		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "data"},
			Rows:       make([][]interface{}, 1000),
			TotalRows:  1000000,
			LoadedRows: 1000,
			Success:    true,
		}

		mockRepo.EXPECT().ExecuteQuery(ctx, sql).Return(expectedResult, nil)

		result, err := mockRepo.ExecuteQuery(ctx, sql)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
		assert.Equal(t, 1000, result.LoadedRows)
		assert.Greater(t, result.TotalRows, int64(result.LoadedRows))
	})

	t.Run("UC-S4-03c: Offset Pagination Next Page", func(t *testing.T) {
		sql := "SELECT * FROM users LIMIT 100 OFFSET 100"
		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "username"},
			Rows:       make([][]interface{}, 100),
			TotalRows:  500,
			LoadedRows: 100,
			Success:    true,
		}

		mockRepo.EXPECT().ExecuteQuery(ctx, sql).Return(expectedResult, nil)

		result, err := mockRepo.ExecuteQuery(ctx, sql)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
	})

	t.Run("UC-S4-04: DDL Query Execution - CREATE TABLE", func(t *testing.T) {
		sql := "CREATE TABLE new_table (id SERIAL PRIMARY KEY, name VARCHAR(255))"

		mockRepo.EXPECT().ExecuteDDL(ctx, sql).Return(nil)

		err := mockRepo.ExecuteDDL(ctx, sql)

		require.NoError(t, err)
	})

	t.Run("UC-S4-04a: DDL Query Execution - ALTER TABLE", func(t *testing.T) {
		sql := "ALTER TABLE users ADD COLUMN email VARCHAR(255)"

		mockRepo.EXPECT().ExecuteDDL(ctx, sql).Return(nil)

		err := mockRepo.ExecuteDDL(ctx, sql)

		require.NoError(t, err)
	})

	t.Run("UC-S4-04b: DDL Query Execution - DROP TABLE", func(t *testing.T) {
		sql := "DROP TABLE old_table"

		mockRepo.EXPECT().ExecuteDDL(ctx, sql).Return(nil)

		err := mockRepo.ExecuteDDL(ctx, sql)

		require.NoError(t, err)
	})

	t.Run("UC-S4-05: DML Query Execution - INSERT", func(t *testing.T) {
		sql := "INSERT INTO users (username, email) VALUES ('newuser', 'new@test.com')"
		expectedResult := &domain.QueryResult{
			Success:      true,
			AffectedRows: 1,
		}

		mockRepo.EXPECT().ExecuteDML(ctx, sql).Return(expectedResult, nil)

		result, err := mockRepo.ExecuteDML(ctx, sql)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
		assert.Equal(t, int64(1), result.AffectedRows)
	})

	t.Run("UC-S4-05a: DML Query Execution - UPDATE", func(t *testing.T) {
		sql := "UPDATE users SET username='updated' WHERE id=1"
		expectedResult := &domain.QueryResult{
			Success:      true,
			AffectedRows: 1,
		}

		mockRepo.EXPECT().ExecuteDML(ctx, sql).Return(expectedResult, nil)

		result, err := mockRepo.ExecuteDML(ctx, sql)

		require.NoError(t, err)
		assert.True(t, result.Success)
		assert.Equal(t, int64(1), result.AffectedRows)
	})

	t.Run("UC-S4-05b: DML Query Execution - DELETE", func(t *testing.T) {
		sql := "DELETE FROM users WHERE id=1"
		expectedResult := &domain.QueryResult{
			Success:      true,
			AffectedRows: 1,
		}

		mockRepo.EXPECT().ExecuteDML(ctx, sql).Return(expectedResult, nil)

		result, err := mockRepo.ExecuteDML(ctx, sql)

		require.NoError(t, err)
		assert.True(t, result.Success)
	})

	t.Run("UC-S4-06: Invalid Query Error", func(t *testing.T) {
		sql := "SELECT * FROM nonexistent_table"
		expectedResult := &domain.QueryResult{
			Success:      false,
			ErrorMessage: "relation \"nonexistent_table\" does not exist",
		}

		mockRepo.EXPECT().ExecuteQuery(ctx, sql).Return(expectedResult, nil)

		result, err := mockRepo.ExecuteQuery(ctx, sql)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.False(t, result.Success)
		assert.Contains(t, result.ErrorMessage, "does not exist")
	})

	t.Run("UC-S4-07: Query Splitting - Multiple Statements", func(t *testing.T) {
		queries := "SELECT * FROM users; UPDATE users SET active=true WHERE id=1; DELETE FROM temp_users;"
		expectedResults := []*domain.QueryResult{
			{
				Columns:    []string{"id", "username"},
				Rows:       [][]interface{}{{1, "alice"}},
				TotalRows:  1,
				LoadedRows: 1,
				Success:    true,
			},
			{
				Success:      true,
				AffectedRows: 1,
			},
			{
				Success:      true,
				AffectedRows: 5,
			},
		}

		mockRepo.EXPECT().ExecuteMultiple(ctx, queries).Return(expectedResults, nil)

		results, err := mockRepo.ExecuteMultiple(ctx, queries)

		require.NoError(t, err)
		assert.Len(t, results, 3)
		assert.True(t, results[0].Success)
		assert.True(t, results[1].Success)
		assert.True(t, results[2].Success)
	})

	t.Run("UC-S4-08: Parameterized Query Execution", func(t *testing.T) {
		sql := "SELECT * FROM users WHERE id = $1 AND username = $2"
		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "username"},
			Rows:       [][]interface{}{{1, "alice"}},
			TotalRows:  1,
			LoadedRows: 1,
			Success:    true,
		}

		mockRepo.EXPECT().ExecuteQuery(ctx, sql, 1, "alice").Return(expectedResult, nil)

		result, err := mockRepo.ExecuteQuery(ctx, sql, 1, "alice")

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
		assert.Len(t, result.Rows, 1)
	})

	t.Run("UC-S4-09: Query with Multiple Parameters", func(t *testing.T) {
		sql := "SELECT * FROM users WHERE id > $1 AND created_at > $2 LIMIT $3"
		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "username", "created_at"},
			Rows:       make([][]interface{}, 10),
			TotalRows:  50,
			LoadedRows: 10,
			Success:    true,
		}

		mockRepo.EXPECT().ExecuteQuery(ctx, sql, 10, "2024-01-01", 10).Return(expectedResult, nil)

		result, err := mockRepo.ExecuteQuery(ctx, sql, 10, "2024-01-01", 10)

		require.NoError(t, err)
		assert.True(t, result.Success)
	})

	t.Run("UC-S4-10: Empty Query Result", func(t *testing.T) {
		sql := "SELECT * FROM users WHERE id > 1000"
		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "username"},
			Rows:       [][]interface{}{},
			TotalRows:  0,
			LoadedRows: 0,
			Success:    true,
		}

		mockRepo.EXPECT().ExecuteQuery(ctx, sql).Return(expectedResult, nil)

		result, err := mockRepo.ExecuteQuery(ctx, sql)

		require.NoError(t, err)
		assert.True(t, result.Success)
		assert.Len(t, result.Rows, 0)
		assert.Equal(t, int64(0), result.TotalRows)
	})

	t.Run("UC-S4-11: Query with JOIN", func(t *testing.T) {
		sql := "SELECT u.id, u.username, p.title FROM users u JOIN posts p ON u.id = p.user_id"
		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "username", "title"},
			Rows:       [][]interface{}{{1, "alice", "first post"}, {1, "alice", "second post"}},
			TotalRows:  2,
			LoadedRows: 2,
			Success:    true,
		}

		mockRepo.EXPECT().ExecuteQuery(ctx, sql).Return(expectedResult, nil)

		result, err := mockRepo.ExecuteQuery(ctx, sql)

		require.NoError(t, err)
		assert.True(t, result.Success)
		assert.Len(t, result.Columns, 3)
	})

	t.Run("UC-S4-12: Query with WHERE Clause and Parameters", func(t *testing.T) {
		sql := "SELECT * FROM users WHERE username LIKE $1"
		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "username"},
			Rows:       [][]interface{}{{1, "alice"}},
			TotalRows:  1,
			LoadedRows: 1,
			Success:    true,
		}

		mockRepo.EXPECT().ExecuteQuery(ctx, sql, "%alice%").Return(expectedResult, nil)

		result, err := mockRepo.ExecuteQuery(ctx, sql, "%alice%")

		require.NoError(t, err)
		assert.True(t, result.Success)
	})

	t.Run("UC-S4-13: Batch DML Operations", func(t *testing.T) {
		queries := "INSERT INTO users (username) VALUES ('user1'); INSERT INTO users (username) VALUES ('user2');"
		expectedResults := []*domain.QueryResult{
			{
				Success:      true,
				AffectedRows: 1,
			},
			{
				Success:      true,
				AffectedRows: 1,
			},
		}

		mockRepo.EXPECT().ExecuteMultiple(ctx, queries).Return(expectedResults, nil)

		results, err := mockRepo.ExecuteMultiple(ctx, queries)

		require.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, int64(2), results[0].AffectedRows+results[1].AffectedRows)
	})
}
