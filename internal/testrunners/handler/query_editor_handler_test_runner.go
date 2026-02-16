package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/interfaces/handler"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
	mockUsecase "github.com/kamil5b/lumen-pg/internal/testrunners/mocks/usecase"
)

// QueryEditorHandlerConstructor is a function type that creates a QueryEditorHandler
type QueryEditorHandlerConstructor func(
	queryUC usecase.QueryUseCase,
	authUC usecase.AuthenticationUseCase,
) handler.QueryEditorHandler

// QueryEditorHandlerRunner runs all query editor handler tests
// Maps to TEST_PLAN.md:
// - Story 4: Manual Query Editor [UC-S4-01~08, E2E-S4-01~06]
//
// NOTE: Authentication and session validation are MIDDLEWARE concerns
// NOTE: Handlers assume middleware has already validated session and populated user context
// NOTE: Handlers focus on business logic: query execution, validation, pagination, result formatting
func QueryEditorHandlerRunner(t *testing.T, constructor QueryEditorHandlerConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockQuery := mockUsecase.NewMockQueryUseCase(ctrl)
	mockAuth := mockUsecase.NewMockAuthenticationUseCase(ctrl)

	h := constructor(mockQuery, mockAuth)

	// E2E-S4-01: Query Editor Page Access
	t.Run("E2E-S4-01: Query Editor Page Access", func(t *testing.T) {
		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		req := httptest.NewRequest(http.MethodGet, "/query-editor", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleQueryEditorPage(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify query editor page elements
		require.Contains(t, body, "query-editor")
		require.Contains(t, body, "Execute")
		require.Contains(t, body, "textarea")
		require.Contains(t, body, "results-panel")
	})

	// E2E-S4-02: Execute Single Query
	t.Run("E2E-S4-02: Execute Single Query", func(t *testing.T) {
		form := url.Values{}
		form.Add("query", "SELECT * FROM users WHERE id > 10")

		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockQuery.EXPECT().
			ExecuteQueryWithPagination(gomock.Any(), "testuser", gomock.Any()).
			Return(&domain.QueryResult{
				Columns: []string{"id", "name", "email"},
				Rows: []map[string]interface{}{
					{"id": 11, "name": "Alice", "email": "alice@example.com"},
					{"id": 12, "name": "Bob", "email": "bob@example.com"},
				},
				RowCount:   2,
				TotalCount: 2,
			}, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/query/execute", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleExecuteQuery(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify results are displayed
		require.Contains(t, body, "Alice")
		require.Contains(t, body, "Bob")
		require.Contains(t, body, "id")
		require.Contains(t, body, "name")
		require.Contains(t, body, "email")
	})

	// E2E-S4-03: Execute Multiple Queries
	t.Run("E2E-S4-03: Execute Multiple Queries", func(t *testing.T) {
		form := url.Values{}
		form.Add("query", "SELECT * FROM users LIMIT 5; SELECT * FROM posts LIMIT 5;")

		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockQuery.EXPECT().
			ExecuteMultipleQueries(gomock.Any(), "testuser", "SELECT * FROM users LIMIT 5; SELECT * FROM posts LIMIT 5;").
			Return([]domain.QueryResult{
				{
					Columns:  []string{"id", "name"},
					Rows:     []map[string]interface{}{{"id": 1, "name": "User1"}},
					RowCount: 1,
				},
				{
					Columns:  []string{"id", "title"},
					Rows:     []map[string]interface{}{{"id": 1, "title": "Post1"}},
					RowCount: 1,
				},
			}, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/query/execute-multiple", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleExecuteMultipleQueries(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify both result sets are displayed
		require.Contains(t, body, "User1")
		require.Contains(t, body, "Post1")
		require.Contains(t, body, "Result 1")
		require.Contains(t, body, "Result 2")
	})

	// E2E-S4-04: Query Error Display
	t.Run("E2E-S4-04: Query Error Display", func(t *testing.T) {
		form := url.Values{}
		form.Add("query", "SELECT * FORM users")

		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockQuery.EXPECT().
			ExecuteQueryWithPagination(gomock.Any(), "testuser", gomock.Any()).
			Return(nil, domain.ValidationError{
				Field:   "query",
				Message: "syntax error at or near \"FORM\"",
			})

		req := httptest.NewRequest(http.MethodPost, "/api/query/execute", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleExecuteQuery(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusBadRequest, rec.Code)
		body := rec.Body.String()

		// Verify error message is displayed in red
		require.Contains(t, body, "syntax error")
		require.Contains(t, body, "error")
		require.Contains(t, body, "FORM")
	})

	// E2E-S4-05: Offset Pagination Results
	t.Run("E2E-S4-05: Offset Pagination Results", func(t *testing.T) {
		form := url.Values{}
		form.Add("query", "SELECT * FROM large_table")

		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockQuery.EXPECT().
			ExecuteQueryWithPagination(gomock.Any(), "testuser", gomock.Any()).
			Return(&domain.QueryResult{
				Columns:    []string{"id", "data"},
				Rows:       make([]map[string]interface{}, 1000),
				RowCount:   1000,
				TotalCount: 5000,
			}, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/query/execute", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleExecuteQuery(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify pagination info is displayed
		require.Contains(t, body, "Data size: 5000 rows")
		require.Contains(t, body, "pagination")
		require.Contains(t, body, "1000")
	})

	// E2E-S4-05a: Offset Pagination Navigation
	t.Run("E2E-S4-05a: Offset Pagination Navigation", func(t *testing.T) {
		form := url.Values{}
		form.Add("query", "SELECT * FROM large_table")
		form.Add("offset", "1000")

		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockQuery.EXPECT().
			ExecuteQueryWithPagination(gomock.Any(), "testuser", gomock.Any()).
			Return(&domain.QueryResult{
				Columns:    []string{"id", "data"},
				Rows:       []map[string]interface{}{},
				RowCount:   0,
				TotalCount: 10000,
			}, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/query/execute", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleExecuteQuery(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify hard limit message
		require.Contains(t, body, "hard limit of 1000 rows reached")
	})

	// E2E-S4-05b: Query Result Actual Size vs Display Limit
	t.Run("E2E-S4-05b: Query Result Actual Size vs Display Limit", func(t *testing.T) {
		form := url.Values{}
		form.Add("query", "SELECT * FROM very_large_table")

		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockQuery.EXPECT().
			ExecuteQueryWithPagination(gomock.Any(), "testuser", gomock.Any()).
			Return(&domain.QueryResult{
				Columns:    []string{"id", "data"},
				Rows:       make([]map[string]interface{}, 1000),
				RowCount:   1000,
				TotalCount: 10000,
			}, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/query/execute", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleExecuteQuery(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify actual size is shown but only first 1000 accessible
		require.Contains(t, body, "Data size: 10000 rows")
		require.Contains(t, body, "only first 1000 are accessible")
	})

	// E2E-S4-06: SQL Syntax Highlighting
	t.Run("E2E-S4-06: SQL Syntax Highlighting", func(t *testing.T) {
		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		req := httptest.NewRequest(http.MethodGet, "/query-editor", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleQueryEditorPage(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify syntax highlighting is enabled
		require.Contains(t, body, "syntax-highlight")
		require.Contains(t, body, "sql")
	})

	// Additional test: DDL query execution
	t.Run("Execute DDL Query", func(t *testing.T) {
		form := url.Values{}
		form.Add("query", "CREATE TABLE test_table (id SERIAL PRIMARY KEY, name VARCHAR(100))")

		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockQuery.EXPECT().
			ExecuteQueryWithPagination(gomock.Any(), "testuser", gomock.Any()).
			Return(&domain.QueryResult{
				Columns:  []string{},
				Rows:     []map[string]interface{}{},
				RowCount: 0,
			}, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/query/execute", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleExecuteQuery(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)

		// Verify DDL execution success
		require.Equal(t, http.StatusOK, rec.Code)
	})

	// Additional test: DML query execution
	t.Run("Execute DML Query", func(t *testing.T) {
		form := url.Values{}
		form.Add("query", "INSERT INTO users (name, email) VALUES ('John', 'john@example.com')")

		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockQuery.EXPECT().
			ExecuteQueryWithPagination(gomock.Any(), "testuser", gomock.Any()).
			Return(&domain.QueryResult{
				Columns:  []string{},
				Rows:     []map[string]interface{}{},
				RowCount: 1,
			}, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/query/execute", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleExecuteQuery(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)

		// Verify DML execution success
		require.Equal(t, http.StatusOK, rec.Code)
	})

	// Additional test: Permission denied
	t.Run("Query with Permission Denied", func(t *testing.T) {
		form := url.Values{}
		form.Add("query", "SELECT * FROM admin_only_table")

		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockQuery.EXPECT().
			ExecuteQueryWithPagination(gomock.Any(), "testuser", gomock.Any()).
			Return(nil, domain.ValidationError{
				Field:   "permission",
				Message: "permission denied for table admin_only_table",
			})

		req := httptest.NewRequest(http.MethodPost, "/api/query/execute", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleExecuteQuery(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusForbidden, rec.Code)
		body := rec.Body.String()

		// Verify permission error
		require.Contains(t, body, "permission denied")
	})

	// Additional test: Unauthorized access
	t.Run("Unauthorized Access to Query Editor", func(t *testing.T) {
		// No mock needed - handler returns early when no cookie is present
		req := httptest.NewRequest(http.MethodGet, "/query-editor", nil)
		rec := httptest.NewRecorder()

		h.HandleQueryEditorPage(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	// Additional test: Empty query
	t.Run("Execute Empty Query", func(t *testing.T) {
		form := url.Values{}
		form.Add("query", "")

		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/query/execute", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleExecuteQuery(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusBadRequest, rec.Code)
		body := rec.Body.String()

		// Verify error message
		require.Contains(t, body, "Query cannot be empty")
	})
}
