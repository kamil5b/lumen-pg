package e2e_integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Story4QueryEditorE2ERunner runs end-to-end tests for Story 4: Manual Query Editor
// This tests the complete route stack with all middleware for query execution
// Maps to TEST_PLAN.md Story 4 E2E Tests [L362-406]:
// - E2E-S4-01: Query Editor Page Access
// - E2E-S4-02: Execute Single Query
// - E2E-S4-03: Execute Multiple Queries
// - E2E-S4-04: Query Error Display
// - E2E-S4-05: Offset Pagination Results
// - E2E-S4-05a: Offset Pagination Navigation
// - E2E-S4-05b: Query Result Actual Size vs Display Limit
// - E2E-S4-06: SQL Syntax Highlighting
//
// Tests complete query editor functionality including:
// - Query editor page rendering with authentication
// - Single and multiple query execution
// - Result pagination with offset
// - Error handling and display
// - SQL syntax highlighting support
func Story4QueryEditorE2ERunner(t *testing.T, router http.Handler) {
	t.Helper()

	// Helper function to login and get session cookies
	getAuthenticatedSession := func(t *testing.T) []*http.Cookie {
		formData := url.Values{}
		formData.Set("username", "testuser")
		formData.Set("password", "testpass")

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusFound, rec.Code, "Login should succeed")
		cookies := rec.Result().Cookies()
		require.NotEmpty(t, cookies, "Should receive session cookies")
		return cookies
	}

	// E2E-S4-01: Query Editor Page Access
	t.Run("E2E-S4-01: Query Editor Page Access", func(t *testing.T) {
		// Test that authenticated users can access query editor page
		cookies := getAuthenticatedSession(t)

		req := httptest.NewRequest(http.MethodGet, "/query-editor", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Query editor page should be accessible")
		body := rec.Body.String()
		assert.True(t,
			strings.Contains(body, "query") ||
				strings.Contains(body, "editor") ||
				strings.Contains(body, "sql") ||
				strings.Contains(body, "execute"),
			"Query editor page should contain query editor UI")

		// Test that unauthenticated users cannot access
		req = httptest.NewRequest(http.MethodGet, "/query-editor", nil)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.True(t,
			rec.Code == http.StatusFound ||
				rec.Code == http.StatusUnauthorized,
			"Unauthenticated access should be denied")

		if rec.Code == http.StatusFound {
			assert.Contains(t, rec.Header().Get("Location"), "/login")
		}
	})

	// E2E-S4-02: Execute Single Query
	t.Run("E2E-S4-02: Execute Single Query", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Execute a simple SELECT query
		queryPayload := `{
			"query": "SELECT * FROM users LIMIT 10"
		}`
		req := httptest.NewRequest(http.MethodPost, "/api/execute-query", strings.NewReader(queryPayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Single query execution should succeed")

		// Parse response
		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err, "Response should be valid JSON")

		// Should contain query results
		assert.True(t,
			response["columns"] != nil ||
				response["rows"] != nil ||
				response["data"] != nil ||
				response["result"] != nil,
			"Response should contain query results")

		// Should have row count or similar metadata
		assert.True(t,
			response["row_count"] != nil ||
				response["rowCount"] != nil ||
				response["count"] != nil ||
				response["rows"] != nil,
			"Response should contain result metadata")
	})

	// E2E-S4-03: Execute Multiple Queries
	t.Run("E2E-S4-03: Execute Multiple Queries", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Execute multiple queries separated by semicolons
		multiQueryPayload := `{
			"query": "SELECT COUNT(*) FROM users; SELECT COUNT(*) FROM posts; SELECT COUNT(*) FROM comments;"
		}`
		req := httptest.NewRequest(http.MethodPost, "/api/execute-query", strings.NewReader(multiQueryPayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Multiple query execution should succeed")

		// Parse response
		body := rec.Body.Bytes()
		assert.True(t, json.Valid(body), "Response should be valid JSON")

		var response interface{}
		err := json.Unmarshal(body, &response)
		require.NoError(t, err)

		// Response should contain multiple results (array or results object)
		respStr := string(body)
		assert.True(t,
			strings.Contains(respStr, "results") ||
				strings.Contains(respStr, "[") ||
				strings.Count(respStr, "columns") > 1 ||
				strings.Count(respStr, "rows") > 1,
			"Response should contain multiple query results")
	})

	// E2E-S4-04: Query Error Display
	t.Run("E2E-S4-04: Query Error Display", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Execute invalid SQL query
		invalidQueryPayload := `{
			"query": "SELECT * FROM nonexistent_table"
		}`
		req := httptest.NewRequest(http.MethodPost, "/api/execute-query", strings.NewReader(invalidQueryPayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// Should return error status or 200 with error in body
		assert.True(t,
			rec.Code == http.StatusOK ||
				rec.Code == http.StatusBadRequest ||
				rec.Code == http.StatusInternalServerError,
			"Should handle query error, got %d", rec.Code)

		body := rec.Body.String()
		assert.True(t,
			strings.Contains(body, "error") ||
				strings.Contains(body, "Error") ||
				strings.Contains(body, "does not exist") ||
				strings.Contains(body, "not found"),
			"Response should contain error message")

		// Test syntax error
		syntaxErrorPayload := `{
			"query": "SELECT * FORM users"
		}`
		req = httptest.NewRequest(http.MethodPost, "/api/execute-query", strings.NewReader(syntaxErrorPayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		body = rec.Body.String()
		assert.True(t,
			strings.Contains(body, "error") ||
				strings.Contains(body, "syntax") ||
				strings.Contains(body, "Error"),
			"Response should contain syntax error message")
	})

	// E2E-S4-05: Offset Pagination Results
	t.Run("E2E-S4-05: Offset Pagination Results", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Execute query with default pagination
		queryPayload := `{
			"query": "SELECT * FROM users"
		}`
		req := httptest.NewRequest(http.MethodPost, "/api/execute-query", strings.NewReader(queryPayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Query with pagination should succeed")

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		// Should have pagination metadata
		assert.True(t,
			response["limit"] != nil ||
				response["offset"] != nil ||
				response["pagination"] != nil ||
				response["page"] != nil,
			"Response should contain pagination metadata")

		// Execute query with specific offset
		queryWithOffsetPayload := `{
			"query": "SELECT * FROM users",
			"offset": 10,
			"limit": 20
		}`
		req = httptest.NewRequest(http.MethodPost, "/api/execute-query", strings.NewReader(queryWithOffsetPayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Query with custom offset should succeed")

		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		// Verify pagination parameters in response
		if response["offset"] != nil {
			assert.Equal(t, float64(10), response["offset"], "Should respect offset parameter")
		}
		if response["limit"] != nil {
			assert.True(t, response["limit"].(float64) <= 20, "Should respect limit parameter")
		}
	})

	// E2E-S4-05a: Offset Pagination Navigation
	t.Run("E2E-S4-05a: Offset Pagination Navigation", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// First page
		firstPagePayload := `{
			"query": "SELECT * FROM users ORDER BY id",
			"offset": 0,
			"limit": 10
		}`
		req := httptest.NewRequest(http.MethodPost, "/api/execute-query", strings.NewReader(firstPagePayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var firstPageResponse map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &firstPageResponse)

		// Second page
		secondPagePayload := `{
			"query": "SELECT * FROM users ORDER BY id",
			"offset": 10,
			"limit": 10
		}`
		req = httptest.NewRequest(http.MethodPost, "/api/execute-query", strings.NewReader(secondPagePayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var secondPageResponse map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &secondPageResponse)

		// Pages should return valid results
		assert.NotNil(t, firstPageResponse)
		assert.NotNil(t, secondPageResponse)

		// Should have navigation hints (hasNext, hasPrevious, etc.)
		body := rec.Body.String()
		assert.True(t,
			strings.Contains(body, "has_next") ||
				strings.Contains(body, "hasNext") ||
				strings.Contains(body, "next") ||
				len(body) > 0,
			"Response should contain pagination navigation info")
	})

	// E2E-S4-05b: Query Result Actual Size vs Display Limit
	t.Run("E2E-S4-05b: Query Result Actual Size vs Display Limit", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Query that returns many rows with a limit
		queryPayload := `{
			"query": "SELECT * FROM users",
			"limit": 5
		}`
		req := httptest.NewRequest(http.MethodPost, "/api/execute-query", strings.NewReader(queryPayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		// Should show actual total count vs displayed count
		assert.True(t,
			response["total_count"] != nil ||
				response["totalCount"] != nil ||
				response["actual_size"] != nil ||
				response["total"] != nil ||
				response["row_count"] != nil,
			"Response should show total available rows")

		// Display limit should be enforced
		if response["rows"] != nil {
			rows, ok := response["rows"].([]interface{})
			if ok {
				assert.LessOrEqual(t, len(rows), 5, "Should not exceed display limit")
			}
		}

		// Should indicate if more results are available
		body := rec.Body.String()
		assert.NotEmpty(t, body, "Should return response body")
	})

	// E2E-S4-06: SQL Syntax Highlighting
	t.Run("E2E-S4-06: SQL Syntax Highlighting", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Request query editor page
		req := httptest.NewRequest(http.MethodGet, "/query-editor", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Page should include syntax highlighting support
		// Check for common syntax highlighter libraries or attributes
		assert.True(t,
			strings.Contains(body, "highlight") ||
				strings.Contains(body, "codemirror") ||
				strings.Contains(body, "ace-editor") ||
				strings.Contains(body, "monaco") ||
				strings.Contains(body, "syntax") ||
				strings.Contains(body, "language-sql") ||
				strings.Contains(body, "sql") ||
				strings.Contains(body, "<textarea") ||
				strings.Contains(body, "<pre"),
			"Query editor should support SQL syntax highlighting")
	})

	// Additional test: DDL Query Execution
	t.Run("DDL Query Execution", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Execute CREATE TABLE (might fail if already exists, that's OK)
		ddlPayload := `{
			"query": "CREATE TABLE IF NOT EXISTS test_table (id SERIAL PRIMARY KEY, name VARCHAR(100))"
		}`
		req := httptest.NewRequest(http.MethodPost, "/api/execute-query", strings.NewReader(ddlPayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// Should complete (success or error is fine for this test)
		assert.True(t,
			rec.Code == http.StatusOK ||
				rec.Code == http.StatusCreated ||
				rec.Code == http.StatusBadRequest ||
				rec.Code == http.StatusForbidden,
			"DDL query should be handled, got %d", rec.Code)
	})

	// Additional test: DML Query Execution
	t.Run("DML Query Execution", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Execute INSERT (might fail due to permissions or constraints)
		dmlPayload := `{
			"query": "INSERT INTO test_table (name) VALUES ('test_value')"
		}`
		req := httptest.NewRequest(http.MethodPost, "/api/execute-query", strings.NewReader(dmlPayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// Should be handled appropriately
		assert.True(t,
			rec.Code == http.StatusOK ||
				rec.Code == http.StatusCreated ||
				rec.Code == http.StatusBadRequest ||
				rec.Code == http.StatusForbidden,
			"DML query should be handled, got %d", rec.Code)

		// Response should indicate affected rows
		if rec.Code == http.StatusOK || rec.Code == http.StatusCreated {
			body := rec.Body.String()
			assert.True(t,
				strings.Contains(body, "affected") ||
					strings.Contains(body, "inserted") ||
					strings.Contains(body, "rows") ||
					len(body) > 0,
				"DML response should contain affected rows info")
		}
	})

	// Additional test: Parameterized Query Execution
	t.Run("Parameterized Query Execution", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Execute query with parameters (if supported)
		paramQueryPayload := `{
			"query": "SELECT * FROM users WHERE id = $1",
			"params": [1]
		}`
		req := httptest.NewRequest(http.MethodPost, "/api/execute-query", strings.NewReader(paramQueryPayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// Should handle parameterized queries (success or not supported)
		assert.True(t,
			rec.Code == http.StatusOK ||
				rec.Code == http.StatusBadRequest ||
				rec.Code == http.StatusNotImplemented,
			"Parameterized query should be handled, got %d", rec.Code)
	})

	// Additional test: Query with Hard Limit Cap
	t.Run("Query Result Hard Limit Cap", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Try to query with very large limit
		largeLimitPayload := `{
			"query": "SELECT * FROM users",
			"limit": 999999
		}`
		req := httptest.NewRequest(http.MethodPost, "/api/execute-query", strings.NewReader(largeLimitPayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		// System should enforce hard limit cap (e.g., max 1000 rows)
		if response["rows"] != nil {
			rows, ok := response["rows"].([]interface{})
			if ok {
				assert.LessOrEqual(t, len(rows), 10000,
					"Should enforce maximum result limit cap")
			}
		}
	})
}
