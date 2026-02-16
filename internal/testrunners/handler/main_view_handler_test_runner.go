package handler

import (
	"context"
	"io"
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

// MainViewHandlerConstructor is a function type that creates a MainViewHandler
type MainViewHandlerConstructor func(
	dataViewUC usecase.DataViewUseCase,
	authUC usecase.AuthenticationUseCase,
	rbacUC usecase.RBACUseCase,
) handler.MainViewHandler

// MainViewHandlerRunner runs all main view handler tests
// Maps to TEST_PLAN.md:
// - Story 5: Main View & Data Interaction [UC-S5-01~19, E2E-S5-01~15a]
//
// NOTE: Authentication and session validation are MIDDLEWARE concerns
// NOTE: Handlers assume middleware has already validated session and populated user context
// NOTE: Handlers focus on business logic: data retrieval, filtering, transactions, navigation
func MainViewHandlerRunner(t *testing.T, constructor MainViewHandlerConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockDataView := mockUsecase.NewMockDataViewUseCase(ctrl)
	mockAuth := mockUsecase.NewMockAuthenticationUseCase(ctrl)
	mockRBAC := mockUsecase.NewMockRBACUseCase(ctrl)

	h := constructor(mockDataView, mockAuth, mockRBAC)

	// E2E-S5-01: Main View Default Load
	t.Run("E2E-S5-01: Main View Default Load", func(t *testing.T) {
		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockAuth.EXPECT().
			GetUserAccessibleResources(gomock.Any(), "testuser").
			Return(&domain.RoleMetadata{
				Name:                "testuser",
				AccessibleDatabases: []string{"testdb"},
				AccessibleSchemas:   []string{"public"},
				AccessibleTables: []domain.AccessibleTable{
					{Database: "testdb", Schema: "public", Name: "users", HasSelect: true},
				},
			}, nil)

		mockDataView.EXPECT().
			LoadTableData(gomock.Any(), "testuser", gomock.Any()).
			Return(&domain.QueryResult{
				Columns: []string{"id", "name", "email"},
				Rows: []map[string]interface{}{
					{"id": 1, "name": "Alice", "email": "alice@example.com"},
					{"id": 2, "name": "Bob", "email": "bob@example.com"},
				},
				RowCount:   2,
				TotalCount: 2,
			}, nil)

		req := httptest.NewRequest(http.MethodGet, "/main", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleMainViewPage(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify first accessible table is loaded
		require.Contains(t, body, "users")
		require.Contains(t, body, "Alice")
		require.Contains(t, body, "Bob")
	})

	// E2E-S5-02: Table Selection from Sidebar
	t.Run("E2E-S5-02: Table Selection from Sidebar", func(t *testing.T) {
		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockRBAC.EXPECT().
			CheckTableAccess(gomock.Any(), "testuser", "testdb", "public", "posts").
			Return(true, nil)

		mockDataView.EXPECT().
			LoadTableData(gomock.Any(), "testuser", gomock.Any()).
			Return(&domain.QueryResult{
				Columns: []string{"id", "user_id", "title", "content"},
				Rows: []map[string]interface{}{
					{"id": 1, "user_id": 1, "title": "First Post", "content": "Content here"},
				},
				RowCount:   1,
				TotalCount: 1,
			}, nil)

		req := httptest.NewRequest(http.MethodPost, "/main/select-table", nil)
		form := url.Values{}
		form.Add("database", "testdb")
		form.Add("schema", "public")
		form.Add("table", "posts")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Body = io.NopCloser(strings.NewReader(form.Encode()))
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleTableSelect(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify posts table data is loaded
		require.Contains(t, body, "First Post")
		require.Contains(t, body, "user_id")
	})

	// E2E-S5-03: WHERE Bar Filtering
	t.Run("E2E-S5-03: WHERE Bar Filtering", func(t *testing.T) {
		form := url.Values{}
		form.Add("where", "id > 10 AND status = 'active'")
		form.Add("database", "testdb")
		form.Add("schema", "public")
		form.Add("table", "users")

		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockDataView.EXPECT().
			ValidateWhereClause(gomock.Any(), "id > 10 AND status = 'active'").
			Return(true, nil)

		mockDataView.EXPECT().
			FilterTableData(gomock.Any(), "testuser", "testdb", "public", "users", "id > 10 AND status = 'active'", 0, 50).
			Return(&domain.QueryResult{
				Columns: []string{"id", "name", "status"},
				Rows: []map[string]interface{}{
					{"id": 11, "name": "Charlie", "status": "active"},
				},
				RowCount:   1,
				TotalCount: 1,
			}, nil)

		req := httptest.NewRequest(http.MethodPost, "/main/filter", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleFilterTable(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify filtered results
		require.Contains(t, body, "Charlie")
		require.Contains(t, body, "active")
	})

	// E2E-S5-04: Column Header Sorting
	t.Run("E2E-S5-04: Column Header Sorting", func(t *testing.T) {
		form := url.Values{}
		form.Add("database", "testdb")
		form.Add("schema", "public")
		form.Add("table", "users")
		form.Add("column", "name")
		form.Add("direction", "ASC")

		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockDataView.EXPECT().
			SortTableData(gomock.Any(), "testuser", "testdb", "public", "users", "name", "ASC", 0, 50).
			Return(&domain.QueryResult{
				Columns: []string{"id", "name"},
				Rows: []map[string]interface{}{
					{"id": 1, "name": "Alice"},
					{"id": 2, "name": "Bob"},
					{"id": 3, "name": "Charlie"},
				},
				RowCount:   3,
				TotalCount: 3,
			}, nil)

		req := httptest.NewRequest(http.MethodPost, "/main/sort", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleSortTable(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify sorted data (alphabetically)
		require.Contains(t, body, "Alice")
		require.Contains(t, body, "Bob")
		require.Contains(t, body, "Charlie")
	})

	// E2E-S5-05: Cursor Pagination Infinite Scroll with Actual Size
	t.Run("E2E-S5-05: Cursor Pagination Infinite Scroll with Actual Size", func(t *testing.T) {
		form := url.Values{}
		form.Add("database", "testdb")
		form.Add("schema", "public")
		form.Add("table", "large_table")

		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockDataView.EXPECT().
			LoadTableData(gomock.Any(), "testuser", gomock.Any()).
			Return(&domain.QueryResult{
				Columns:    []string{"id", "data"},
				Rows:       make([]map[string]interface{}, 50),
				RowCount:   50,
				TotalCount: 5000,
			}, nil)

		req := httptest.NewRequest(http.MethodPost, "/main/load-data", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleLoadTableData(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify data size indicator and first 50 rows loaded
		require.Contains(t, body, "5000")
		require.Contains(t, body, "50")
	})

	// E2E-S5-05a: Cursor Pagination Infinite Scroll Loading
	t.Run("E2E-S5-05a: Cursor Pagination Infinite Scroll Loading", func(t *testing.T) {
		form := url.Values{}
		form.Add("database", "testdb")
		form.Add("schema", "public")
		form.Add("table", "large_table")
		form.Add("cursor", "cursor_token_page2")

		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockDataView.EXPECT().
			GetTableDataWithCursorPagination(gomock.Any(), "testuser", "testdb", "public", "large_table", "cursor_token_page2", 50).
			Return(&domain.QueryResult{
				Columns:    []string{"id", "data"},
				Rows:       make([]map[string]interface{}, 50),
				RowCount:   50,
				TotalCount: 5000,
			}, nil)

		req := httptest.NewRequest(http.MethodPost, "/main/pagination-next", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandlePaginationNext(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify next 50 rows loaded
		require.NotEmpty(t, body)
	})

	// E2E-S5-05b: Pagination Hard Limit Enforcement
	t.Run("E2E-S5-05b: Pagination Hard Limit Enforcement", func(t *testing.T) {
		form := url.Values{}
		form.Add("database", "testdb")
		form.Add("schema", "public")
		form.Add("table", "large_table")

		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockDataView.EXPECT().
			LoadTableData(gomock.Any(), "testuser", gomock.Any()).
			Return(&domain.QueryResult{
				Columns:    []string{"id", "data"},
				Rows:       make([]map[string]interface{}, 1000),
				RowCount:   1000,
				TotalCount: 5000,
			}, nil)

		req := httptest.NewRequest(http.MethodPost, "/main/load-data", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleLoadTableData(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify hard limit message
		require.Contains(t, body, "1000")
		require.Contains(t, body, "5000")
	})

	// E2E-S5-14: FK Cell Navigation (Read-Only)
	// NOTE: This uses the NavigationHandler, not MainViewHandler
	// Skipping this test in MainViewHandler runner as it belongs to NavigationHandler

	// E2E-S5-15: PK Cell Navigation (Read-Only)
	// NOTE: This uses the NavigationHandler, not MainViewHandler
	// Skipping this test in MainViewHandler runner as it belongs to NavigationHandler

	// E2E-S5-15a: PK Cell Navigation - Table Click
	// NOTE: This uses the NavigationHandler, not MainViewHandler
	// Skipping this test in MainViewHandler runner as it belongs to NavigationHandler

	// Additional test: Read-only mode enforcement
	t.Run("Read-Only Mode Enforcement", func(t *testing.T) {
		form := url.Values{}
		form.Add("database", "testdb")
		form.Add("schema", "public")
		form.Add("table", "users")

		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "readonly_user",
			}, nil)

		// No RBAC check needed - loading data only requires read permission (already checked by middleware)
		mockDataView.EXPECT().
			LoadTableData(gomock.Any(), "readonly_user", gomock.Any()).
			Return(&domain.QueryResult{
				Columns: []string{"id", "name"},
				Rows: []map[string]interface{}{
					{"id": 1, "name": "Alice"},
				},
				RowCount:   1,
				TotalCount: 1,
			}, nil)

		req := httptest.NewRequest(http.MethodPost, "/main/load-data", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleLoadTableData(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify read-only indicators (no transaction button in UI)
		require.Contains(t, body, "Alice")
		// The absence of "Start Transaction" button would be in the main page
		// which is tested via HandleMainViewPage
	})

	// Additional test: Unauthorized access
	t.Run("Unauthorized Access to Main View", func(t *testing.T) {
		// No mock needed - handler returns early when no cookie is present
		req := httptest.NewRequest(http.MethodGet, "/main", nil)
		rec := httptest.NewRecorder()

		h.HandleMainViewPage(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusFound, rec.Code) // Redirects to /login
	})

	// Additional test: Table not accessible
	t.Run("Table Not Accessible", func(t *testing.T) {
		form := url.Values{}
		form.Add("database", "testdb")
		form.Add("schema", "public")
		form.Add("table", "admin_only_table")

		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockRBAC.EXPECT().
			CheckTableAccess(gomock.Any(), "testuser", "testdb", "public", "admin_only_table").
			Return(false, nil)

		req := httptest.NewRequest(http.MethodPost, "/main/select-table", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleTableSelect(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusForbidden, rec.Code)
	})

	// Additional test: Invalid WHERE clause
	t.Run("Invalid WHERE Clause", func(t *testing.T) {
		form := url.Values{}
		form.Add("where", "id = 1; DROP TABLE users;")
		form.Add("database", "testdb")
		form.Add("schema", "public")
		form.Add("table", "users")

		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockDataView.EXPECT().
			ValidateWhereClause(gomock.Any(), "id = 1; DROP TABLE users;").
			Return(false, domain.ValidationError{Field: "where", Message: "SQL injection detected"})

		req := httptest.NewRequest(http.MethodPost, "/main/filter", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleFilterTable(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusBadRequest, rec.Code)
		body := rec.Body.String()
		require.Contains(t, body, "SQL injection")
	})

	// Additional test: Pagination previous
	t.Run("Pagination Previous", func(t *testing.T) {
		form := url.Values{}
		form.Add("database", "testdb")
		form.Add("schema", "public")
		form.Add("table", "large_table")
		form.Add("cursor", "cursor_token_page1")

		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockDataView.EXPECT().
			LoadTableData(gomock.Any(), "testuser", gomock.Any()).
			Return(&domain.QueryResult{
				Columns:    []string{"id", "data"},
				Rows:       make([]map[string]interface{}, 50),
				RowCount:   50,
				TotalCount: 5000,
			}, nil)

		req := httptest.NewRequest(http.MethodPost, "/main/pagination-previous", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandlePaginationPrevious(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify previous page loaded
		require.NotEmpty(t, body)
	})
}
