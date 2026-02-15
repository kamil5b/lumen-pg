package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
	mockRepository "github.com/kamil5b/lumen-pg/internal/testrunners/mocks/repository"
)

// QueryUsecaseConstructor is a function type that creates a QueryUseCase
type QueryUsecaseConstructor func(
	databaseRepo repository.DatabaseRepository,
	rbacRepo repository.RBACRepository,
) usecase.QueryUseCase

// QueryUsecaseRunner runs all query usecase tests against an implementation
// Maps to TEST_PLAN.md:
// - Story 4: Manual Query Editor [UC-S4-01~08, IT-S4-01~04, E2E-S4-01~06]
func QueryUsecaseRunner(t *testing.T, constructor QueryUsecaseConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDatabase := mockRepository.NewMockDatabaseRepository(ctrl)
	mockRBAC := mockRepository.NewMockRBACRepository(ctrl)

	uc := constructor(mockDatabase, mockRBAC)

	ctx := context.Background()

	// UC-S4-01: Single Query Execution
	// IT-S4-01: Real SELECT Query
	// E2E-S4-02: Execute Single Query
	t.Run("ExecuteQuery executes single SELECT query", func(t *testing.T) {
		mockDatabase.EXPECT().
			ExecuteQuery(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(&domain.QueryResult{
				Columns:  []string{"id", "name"},
				Rows:     []map[string]interface{}{{"id": 1, "name": "test"}},
				RowCount: 1,
			}, nil)

		mockRBAC.EXPECT().
			HasSelectPermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(true, nil).AnyTimes()

		result, err := uc.ExecuteQuery(ctx, "testuser", "SELECT * FROM users LIMIT 10", 0, 10)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, 2, len(result.Columns))
		require.Equal(t, int64(1), result.RowCount)
	})

	// UC-S4-02: Multiple Query Execution
	// E2E-S4-03: Execute Multiple Queries
	t.Run("ExecuteMultipleQueries executes multiple queries", func(t *testing.T) {
		mockDatabase.EXPECT().
			ExecuteMultipleQueries(gomock.Any(), gomock.Any()).
			Return([]domain.QueryResult{
				{
					Columns:  []string{"id", "name"},
					Rows:     []map[string]interface{}{},
					RowCount: 5,
				},
				{
					Columns:  []string{"id", "title"},
					Rows:     []map[string]interface{}{},
					RowCount: 3,
				},
			}, nil)

		mockRBAC.EXPECT().
			HasSelectPermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(true, nil).AnyTimes()

		results, err := uc.ExecuteMultipleQueries(ctx, "testuser", "SELECT * FROM users; SELECT * FROM posts;")

		require.NoError(t, err)
		require.NotNil(t, results)
		require.Equal(t, 2, len(results))
	})

	// UC-S4-03: Query Result Offset Pagination
	// UC-S4-03a: Query Result Actual Size Display
	// UC-S4-03b: Query Result Limit Hard Cap
	// UC-S4-03c: Offset Pagination Next Page
	// IT-S4-01: Real SELECT Query
	// E2E-S4-05: Offset Pagination Results
	// E2E-S4-05a: Offset Pagination Navigation
	// E2E-S4-05b: Query Result Actual Size vs Display Limit
	t.Run("ExecuteQueryWithPagination returns paginated results", func(t *testing.T) {
		mockDatabase.EXPECT().
			ExecuteQueryWithPagination(gomock.Any(), gomock.Any()).
			Return(&domain.QueryResult{
				Columns:    []string{"id", "name"},
				Rows:       make([]map[string]interface{}, 50),
				RowCount:   50,
				TotalCount: 500,
			}, nil)

		mockRBAC.EXPECT().
			HasSelectPermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(true, nil).AnyTimes()

		result, err := uc.ExecuteQueryWithPagination(ctx, "testuser", domain.QueryParams{
			Query:  "SELECT * FROM users",
			Offset: 50,
			Limit:  50,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, int64(50), result.RowCount)
		require.Greater(t, result.TotalCount, int64(50))
	})

	// UC-S4-03: Query Result Offset Pagination
	t.Run("ExecuteQueryWithPagination respects hard limit", func(t *testing.T) {
		mockDatabase.EXPECT().
			ExecuteQueryWithPagination(gomock.Any(), gomock.Any()).
			Return(&domain.QueryResult{
				Columns:    []string{"id", "name"},
				Rows:       make([]map[string]interface{}, 1000),
				RowCount:   1000,
				TotalCount: 10000,
			}, nil)

		mockRBAC.EXPECT().
			HasSelectPermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(true, nil).AnyTimes()

		result, err := uc.ExecuteQueryWithPagination(ctx, "testuser", domain.QueryParams{
			Query: "SELECT * FROM large_table",
			Limit: 10000,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		require.LessOrEqual(t, result.RowCount, int64(1000))
	})

	// UC-S4-07: Query Splitting
	t.Run("SplitQueries splits multiple queries by semicolon", func(t *testing.T) {
		queries, err := uc.SplitQueries(ctx, "SELECT * FROM users; SELECT * FROM posts; SELECT COUNT(*) FROM comments;")

		require.NoError(t, err)
		require.NotNil(t, queries)
		require.GreaterOrEqual(t, len(queries), 3)
	})

	t.Run("SplitQueries handles single query", func(t *testing.T) {
		queries, err := uc.SplitQueries(ctx, "SELECT * FROM users")

		require.NoError(t, err)
		require.NotNil(t, queries)
		require.GreaterOrEqual(t, len(queries), 1)
	})

	t.Run("SplitQueries handles empty string", func(t *testing.T) {
		queries, err := uc.SplitQueries(ctx, "")

		require.NoError(t, err)
		require.NotNil(t, queries)
	})

	// UC-S4-06: Invalid Query Error
	// E2E-S4-04: Query Error Display
	t.Run("ValidateQuery accepts valid SELECT statement", func(t *testing.T) {
		valid, err := uc.ValidateQuery(ctx, "SELECT * FROM users WHERE id > 10")

		require.NoError(t, err)
		require.True(t, valid)
	})

	t.Run("ValidateQuery rejects invalid SQL", func(t *testing.T) {
		valid, err := uc.ValidateQuery(ctx, "SELECT * FORM users")

		require.NoError(t, err)
		require.False(t, valid)
	})

	// UC-S4-01: Single Query Execution
	t.Run("IsSelectQuery identifies SELECT statements", func(t *testing.T) {
		isSelect, err := uc.IsSelectQuery(ctx, "SELECT id, name FROM users")

		require.NoError(t, err)
		require.True(t, isSelect)
	})

	t.Run("IsSelectQuery rejects non-SELECT statements", func(t *testing.T) {
		isSelect, err := uc.IsSelectQuery(ctx, "INSERT INTO users VALUES (1, 'test')")

		require.NoError(t, err)
		require.False(t, isSelect)
	})

	// UC-S4-04: DDL Query Execution
	// IT-S4-02: Real DDL Query
	t.Run("IsDDLQuery identifies DDL statements", func(t *testing.T) {
		isDDL, err := uc.IsDDLQuery(ctx, "CREATE TABLE test (id INTEGER PRIMARY KEY)")

		require.NoError(t, err)
		require.True(t, isDDL)
	})

	t.Run("IsDDLQuery identifies ALTER statements", func(t *testing.T) {
		isDDL, err := uc.IsDDLQuery(ctx, "ALTER TABLE users ADD COLUMN email VARCHAR(255)")

		require.NoError(t, err)
		require.True(t, isDDL)
	})

	t.Run("IsDDLQuery identifies DROP statements", func(t *testing.T) {
		isDDL, err := uc.IsDDLQuery(ctx, "DROP TABLE test")

		require.NoError(t, err)
		require.True(t, isDDL)
	})

	t.Run("IsDDLQuery rejects SELECT statements", func(t *testing.T) {
		isDDL, err := uc.IsDDLQuery(ctx, "SELECT * FROM users")

		require.NoError(t, err)
		require.False(t, isDDL)
	})

	// UC-S4-05: DML Query Execution
	// IT-S4-03: Real DML Query
	t.Run("IsDMLQuery identifies INSERT statements", func(t *testing.T) {
		isDML, err := uc.IsDMLQuery(ctx, "INSERT INTO users VALUES (1, 'test')")

		require.NoError(t, err)
		require.True(t, isDML)
	})

	t.Run("IsDMLQuery identifies UPDATE statements", func(t *testing.T) {
		isDML, err := uc.IsDMLQuery(ctx, "UPDATE users SET name = 'test' WHERE id = 1")

		require.NoError(t, err)
		require.True(t, isDML)
	})

	t.Run("IsDMLQuery identifies DELETE statements", func(t *testing.T) {
		isDML, err := uc.IsDMLQuery(ctx, "DELETE FROM users WHERE id = 1")

		require.NoError(t, err)
		require.True(t, isDML)
	})

	t.Run("IsDMLQuery rejects SELECT statements", func(t *testing.T) {
		isDML, err := uc.IsDMLQuery(ctx, "SELECT * FROM users")

		require.NoError(t, err)
		require.False(t, isDML)
	})

	// UC-S4-08: Parameterized Query Execution
	// E2E-S4-02: Execute Single Query
	t.Run("ExecuteQuery handles parameterized queries", func(t *testing.T) {
		mockDatabase.EXPECT().
			ExecuteQuery(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(&domain.QueryResult{
				Columns:  []string{"id", "name"},
				Rows:     []map[string]interface{}{{"id": 5, "name": "specific"}},
				RowCount: 1,
			}, nil)

		mockRBAC.EXPECT().
			HasSelectPermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(true, nil).AnyTimes()

		result, err := uc.ExecuteQuery(ctx, "testuser", "SELECT * FROM users WHERE id = $1", 0, 10)

		require.NoError(t, err)
		require.NotNil(t, result)
	})

	// UC-S4-01 ~ UC-S4-08: Query Affected Row Count
	t.Run("GetQueryAffectedRowCount returns affected rows for DML", func(t *testing.T) {
		queryResult := &domain.QueryResult{
			Columns:  []string{},
			Rows:     []map[string]interface{}{},
			RowCount: 15,
		}

		count := uc.GetQueryAffectedRowCount(ctx, queryResult)

		require.GreaterOrEqual(t, count, int64(0))
	})

	// IT-S4-04: Query with Permission Denied
	t.Run("ExecuteQuery respects user permissions", func(t *testing.T) {
		mockRBAC.EXPECT().
			HasSelectPermission(gomock.Any(), "readonlyuser", gomock.Any(), gomock.Any(), gomock.Any()).
			Return(false, nil)

		_, err := uc.ExecuteQuery(ctx, "readonlyuser", "SELECT * FROM secure_table", 0, 10)

		require.Error(t, err)
	})

	// E2E-S4-06: SQL Syntax Highlighting
	t.Run("ValidateQuery for syntax highlighting", func(t *testing.T) {
		validQueries := []string{
			"SELECT * FROM users",
			"INSERT INTO users VALUES (1, 'test')",
			"UPDATE users SET name = 'new' WHERE id = 1",
			"DELETE FROM users WHERE id = 1",
			"CREATE TABLE test (id INT PRIMARY KEY)",
		}

		for _, query := range validQueries {
			valid, err := uc.ValidateQuery(ctx, query)
			require.NoError(t, err, "query: %s", query)
			require.True(t, valid, "query should be valid: %s", query)
		}
	})
}
