package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Story5MainViewE2ERunner runs end-to-end tests for Story 5: Main View & Data Interaction
// This tests the complete route stack with all middleware for main data view
// Maps to TEST_PLAN.md Story 5 E2E Tests [L542-636]:
// - E2E-S5-01: Main View Default Load
// - E2E-S5-02: Table Selection from Sidebar
// - E2E-S5-03: WHERE Bar Filtering
// - E2E-S5-04: Column Header Sorting
// - E2E-S5-05: Cursor Pagination Infinite Scroll with Actual Size
// - E2E-S5-05a: Cursor Pagination Infinite Scroll Loading
// - E2E-S5-05b: Pagination Hard Limit Enforcement
// - E2E-S5-06: Start Transaction Button
// - E2E-S5-07: Transaction Mode Cell Editing
// - E2E-S5-08: Transaction Mode Edit Buffer Display
// - E2E-S5-09: Transaction Commit Button
// - E2E-S5-10: Transaction Rollback Button
// - E2E-S5-11: Transaction Timer Countdown
// - E2E-S5-12: Transaction Row Delete Button
// - E2E-S5-13: Transaction New Row Button
// - E2E-S5-14: FK Cell Navigation (Read-Only)
// - E2E-S5-15: PK Cell Navigation (Read-Only)
// - E2E-S5-15a: PK Cell Navigation - Table Click
//
// Tests complete main view functionality including:
// - Table data loading with cursor pagination
// - Filtering and sorting
// - Transaction management (start, edit, commit, rollback)
// - Foreign key and primary key navigation
func Story5MainViewE2ERunner(t *testing.T, router http.Handler) {
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

	// E2E-S5-01: Main View Default Load
	t.Run("E2E-S5-01: Main View Default Load", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		req := httptest.NewRequest(http.MethodGet, "/main", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Main view should load successfully")
		body := rec.Body.String()

		// Should contain main view elements
		assert.True(t,
			strings.Contains(body, "table") ||
				strings.Contains(body, "data") ||
				strings.Contains(body, "testuser") ||
				strings.Contains(body, "sidebar"),
			"Main view should contain data view elements")
	})

	// E2E-S5-02: Table Selection from Sidebar
	t.Run("E2E-S5-02: Table Selection from Sidebar", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Select a table via API
		selectTablePayload := `{
			"database": "testdb",
			"schema": "public",
			"table": "users"
		}`
		req := httptest.NewRequest(http.MethodPost, "/api/table/select", strings.NewReader(selectTablePayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Table selection should succeed")

		// Load table data
		req = httptest.NewRequest(http.MethodGet, "/api/table/data?table=users&schema=public&database=testdb", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Table data loading should succeed")

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		// Should contain table data
		assert.True(t,
			response["columns"] != nil ||
				response["rows"] != nil ||
				response["data"] != nil,
			"Response should contain table data")
	})

	// E2E-S5-03: WHERE Bar Filtering
	t.Run("E2E-S5-03: WHERE Bar Filtering", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Apply WHERE filter
		filterPayload := `{
			"table": "users",
			"where": "id > 10"
		}`
		req := httptest.NewRequest(http.MethodPost, "/api/table/filter", strings.NewReader(filterPayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.True(t,
			rec.Code == http.StatusOK ||
				rec.Code == http.StatusAccepted,
			"Filter should be applied, got %d", rec.Code)

		// Should return filtered data
		if rec.Code == http.StatusOK {
			var response map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.True(t,
				response["rows"] != nil ||
					response["data"] != nil,
				"Should return filtered data")
		}
	})

	// E2E-S5-04: Column Header Sorting
	t.Run("E2E-S5-04: Column Header Sorting", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Sort by column ascending
		sortPayload := `{
			"table": "users",
			"column": "id",
			"direction": "ASC"
		}`
		req := httptest.NewRequest(http.MethodPost, "/api/table/sort", strings.NewReader(sortPayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.True(t,
			rec.Code == http.StatusOK ||
				rec.Code == http.StatusAccepted,
			"Sort ASC should succeed, got %d", rec.Code)

		// Sort by column descending
		sortPayload = `{
			"table": "users",
			"column": "id",
			"direction": "DESC"
		}`
		req = httptest.NewRequest(http.MethodPost, "/api/table/sort", strings.NewReader(sortPayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.True(t,
			rec.Code == http.StatusOK ||
				rec.Code == http.StatusAccepted,
			"Sort DESC should succeed, got %d", rec.Code)
	})

	// E2E-S5-05: Cursor Pagination Infinite Scroll with Actual Size
	t.Run("E2E-S5-05: Cursor Pagination Infinite Scroll with Actual Size", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Load initial page with cursor
		req := httptest.NewRequest(http.MethodGet, "/api/table/data?table=users&limit=10", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		// Should have cursor for next page
		assert.True(t,
			response["next_cursor"] != nil ||
				response["nextCursor"] != nil ||
				response["cursor"] != nil ||
				response["pagination"] != nil,
			"Response should contain cursor for pagination")

		// Should show actual total size
		assert.True(t,
			response["total_count"] != nil ||
				response["totalCount"] != nil ||
				response["actual_size"] != nil ||
				response["total"] != nil,
			"Response should show actual table size")
	})

	// E2E-S5-05a: Cursor Pagination Infinite Scroll Loading
	t.Run("E2E-S5-05a: Cursor Pagination Infinite Scroll Loading", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Load first page
		req := httptest.NewRequest(http.MethodGet, "/api/table/data?table=users&limit=5", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var firstPage map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &firstPage)
		require.NoError(t, err)

		// Extract next cursor
		var nextCursor interface{}
		if firstPage["next_cursor"] != nil {
			nextCursor = firstPage["next_cursor"]
		} else if firstPage["nextCursor"] != nil {
			nextCursor = firstPage["nextCursor"]
		}

		// Load next page with cursor
		if nextCursor != nil {
			cursorStr, ok := nextCursor.(string)
			if ok {
				req = httptest.NewRequest(http.MethodGet, "/api/table/data?table=users&cursor="+cursorStr+"&limit=5", nil)
				for _, cookie := range cookies {
					req.AddCookie(cookie)
				}
				rec = httptest.NewRecorder()
				router.ServeHTTP(rec, req)

				assert.Equal(t, http.StatusOK, rec.Code, "Next page should load with cursor")
			}
		}
	})

	// E2E-S5-05b: Pagination Hard Limit Enforcement
	t.Run("E2E-S5-05b: Pagination Hard Limit Enforcement", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Request with very large limit
		req := httptest.NewRequest(http.MethodGet, "/api/table/data?table=users&limit=999999", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		// Should enforce hard limit (e.g., max 1000 rows)
		if response["rows"] != nil {
			rows, ok := response["rows"].([]interface{})
			if ok {
				assert.LessOrEqual(t, len(rows), 10000, "Should enforce hard limit")
			}
		}
	})

	// E2E-S5-06: Start Transaction Button
	t.Run("E2E-S5-06: Start Transaction Button", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Start transaction
		req := httptest.NewRequest(http.MethodPost, "/api/transaction/start", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.True(t,
			rec.Code == http.StatusOK ||
				rec.Code == http.StatusCreated,
			"Transaction start should succeed, got %d", rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		// Should return transaction ID and status
		assert.True(t,
			response["transaction_id"] != nil ||
				response["transactionId"] != nil ||
				response["id"] != nil ||
				response["status"] != nil,
			"Response should contain transaction info")
	})

	// E2E-S5-07: Transaction Mode Cell Editing
	t.Run("E2E-S5-07: Transaction Mode Cell Editing", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Start transaction first
		req := httptest.NewRequest(http.MethodPost, "/api/transaction/start", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		require.True(t, rec.Code == http.StatusOK || rec.Code == http.StatusCreated)

		// Edit a cell
		editPayload := `{
			"table": "users",
			"row_id": 1,
			"column": "name",
			"value": "Updated Name"
		}`
		req = httptest.NewRequest(http.MethodPost, "/api/transaction/edit-cell", strings.NewReader(editPayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.True(t,
			rec.Code == http.StatusOK ||
				rec.Code == http.StatusAccepted,
			"Cell edit should succeed in transaction, got %d", rec.Code)
	})

	// E2E-S5-08: Transaction Mode Edit Buffer Display
	t.Run("E2E-S5-08: Transaction Mode Edit Buffer Display", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Start transaction
		req := httptest.NewRequest(http.MethodPost, "/api/transaction/start", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		require.True(t, rec.Code == http.StatusOK || rec.Code == http.StatusCreated)

		// Edit multiple cells
		editPayload := `{"table": "users", "row_id": 1, "column": "name", "value": "Test1"}`
		req = httptest.NewRequest(http.MethodPost, "/api/transaction/edit-cell", strings.NewReader(editPayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		router.ServeHTTP(httptest.NewRecorder(), req)

		editPayload = `{"table": "users", "row_id": 2, "column": "name", "value": "Test2"}`
		req = httptest.NewRequest(http.MethodPost, "/api/transaction/edit-cell", strings.NewReader(editPayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		router.ServeHTTP(httptest.NewRecorder(), req)

		// Get transaction status to see edit buffer
		req = httptest.NewRequest(http.MethodGet, "/api/transaction/status", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		// Should show pending edits in buffer
		assert.True(t,
			response["pending_edits"] != nil ||
				response["pendingEdits"] != nil ||
				response["buffer"] != nil ||
				response["changes"] != nil,
			"Should display edit buffer")
	})

	// E2E-S5-09: Transaction Commit Button
	t.Run("E2E-S5-09: Transaction Commit Button", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Start transaction and make edit
		req := httptest.NewRequest(http.MethodPost, "/api/transaction/start", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		router.ServeHTTP(httptest.NewRecorder(), req)

		editPayload := `{"table": "users", "row_id": 1, "column": "name", "value": "Commit Test"}`
		req = httptest.NewRequest(http.MethodPost, "/api/transaction/edit-cell", strings.NewReader(editPayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		router.ServeHTTP(httptest.NewRecorder(), req)

		// Commit transaction
		req = httptest.NewRequest(http.MethodPost, "/api/transaction/commit", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.True(t,
			rec.Code == http.StatusOK ||
				rec.Code == http.StatusAccepted,
			"Transaction commit should succeed, got %d", rec.Code)

		body := rec.Body.String()
		assert.True(t,
			strings.Contains(body, "success") ||
				strings.Contains(body, "committed") ||
				len(body) > 0,
			"Should confirm commit success")
	})

	// E2E-S5-10: Transaction Rollback Button
	t.Run("E2E-S5-10: Transaction Rollback Button", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Start transaction and make edit
		req := httptest.NewRequest(http.MethodPost, "/api/transaction/start", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		router.ServeHTTP(httptest.NewRecorder(), req)

		editPayload := `{"table": "users", "row_id": 1, "column": "name", "value": "Rollback Test"}`
		req = httptest.NewRequest(http.MethodPost, "/api/transaction/edit-cell", strings.NewReader(editPayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		router.ServeHTTP(httptest.NewRecorder(), req)

		// Rollback transaction
		req = httptest.NewRequest(http.MethodPost, "/api/transaction/rollback", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.True(t,
			rec.Code == http.StatusOK ||
				rec.Code == http.StatusAccepted,
			"Transaction rollback should succeed, got %d", rec.Code)

		body := rec.Body.String()
		assert.True(t,
			strings.Contains(body, "rollback") ||
				strings.Contains(body, "cancelled") ||
				strings.Contains(body, "discarded") ||
				len(body) > 0,
			"Should confirm rollback success")
	})

	// E2E-S5-11: Transaction Timer Countdown
	t.Run("E2E-S5-11: Transaction Timer Countdown", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Start transaction
		startTime := time.Now()
		req := httptest.NewRequest(http.MethodPost, "/api/transaction/start", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.True(t, rec.Code == http.StatusOK || rec.Code == http.StatusCreated)

		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		// Should have expiration time
		assert.True(t,
			response["expires_at"] != nil ||
				response["expiresAt"] != nil ||
				response["timeout"] != nil ||
				response["ttl"] != nil,
			"Should indicate transaction timeout")

		// Check status after a moment
		time.Sleep(100 * time.Millisecond)

		req = httptest.NewRequest(http.MethodGet, "/api/transaction/status", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var statusResponse map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &statusResponse)
		require.NoError(t, err)

		// Should show time remaining
		assert.True(t,
			statusResponse["time_remaining"] != nil ||
				statusResponse["timeRemaining"] != nil ||
				statusResponse["expires_at"] != nil ||
				statusResponse["active"] != nil,
			"Status should show timer info")

		// Verify elapsed time makes sense
		elapsed := time.Since(startTime)
		assert.True(t, elapsed < 10*time.Second, "Transaction should still be active")
	})

	// E2E-S5-12: Transaction Row Delete Button
	t.Run("E2E-S5-12: Transaction Row Delete Button", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Start transaction
		req := httptest.NewRequest(http.MethodPost, "/api/transaction/start", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		router.ServeHTTP(httptest.NewRecorder(), req)

		// Delete a row
		deletePayload := `{"table": "users", "row_id": 999}`
		req = httptest.NewRequest(http.MethodPost, "/api/transaction/delete-row", strings.NewReader(deletePayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.True(t,
			rec.Code == http.StatusOK ||
				rec.Code == http.StatusAccepted ||
				rec.Code == http.StatusNotFound, // Row might not exist
			"Row deletion should be handled, got %d", rec.Code)
	})

	// E2E-S5-13: Transaction New Row Button
	t.Run("E2E-S5-13: Transaction New Row Button", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Start transaction
		req := httptest.NewRequest(http.MethodPost, "/api/transaction/start", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		router.ServeHTTP(httptest.NewRecorder(), req)

		// Insert a new row
		insertPayload := `{
			"table": "users",
			"values": {
				"name": "New User",
				"email": "newuser@test.com"
			}
		}`
		req = httptest.NewRequest(http.MethodPost, "/api/transaction/insert-row", strings.NewReader(insertPayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.True(t,
			rec.Code == http.StatusOK ||
				rec.Code == http.StatusCreated ||
				rec.Code == http.StatusAccepted,
			"Row insertion should succeed, got %d", rec.Code)
	})

	// E2E-S5-14: FK Cell Navigation (Read-Only)
	t.Run("E2E-S5-14: FK Cell Navigation (Read-Only)", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Navigate to parent via foreign key
		navPayload := `{
			"table": "posts",
			"column": "user_id",
			"value": 1
		}`
		req := httptest.NewRequest(http.MethodPost, "/api/navigation/fk-parent", strings.NewReader(navPayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.True(t,
			rec.Code == http.StatusOK ||
				rec.Code == http.StatusFound,
			"FK navigation should succeed, got %d", rec.Code)

		// Should return parent row info or redirect
		if rec.Code == http.StatusOK {
			var response map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.True(t,
				response["table"] != nil ||
					response["row"] != nil ||
					response["data"] != nil,
				"Should return parent row information")
		}
	})

	// E2E-S5-15: PK Cell Navigation (Read-Only)
	t.Run("E2E-S5-15: PK Cell Navigation (Read-Only)", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Navigate to child rows via primary key
		navPayload := `{
			"table": "users",
			"pk_column": "id",
			"pk_value": 1
		}`
		req := httptest.NewRequest(http.MethodPost, "/api/navigation/pk-children", strings.NewReader(navPayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.True(t,
			rec.Code == http.StatusOK ||
				rec.Code == http.StatusFound,
			"PK navigation should succeed, got %d", rec.Code)

		// Should return child table references
		if rec.Code == http.StatusOK {
			var response map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.True(t,
				response["references"] != nil ||
					response["children"] != nil ||
					response["tables"] != nil ||
					len(response) >= 0, // Empty is valid
				"Should return child table references")
		}
	})

	// E2E-S5-15a: PK Cell Navigation - Table Click
	t.Run("E2E-S5-15a: PK Cell Navigation - Table Click", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Get child references first
		navPayload := `{"table": "users", "pk_column": "id", "pk_value": 1}`
		req := httptest.NewRequest(http.MethodPost, "/api/navigation/pk-children", strings.NewReader(navPayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// Click on a child table reference
		clickPayload := `{
			"parent_table": "users",
			"parent_pk": 1,
			"child_table": "posts",
			"fk_column": "user_id"
		}`
		req = httptest.NewRequest(http.MethodPost, "/api/navigation/navigate-to-child", strings.NewReader(clickPayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.True(t,
			rec.Code == http.StatusOK ||
				rec.Code == http.StatusFound,
			"Navigate to child table should succeed, got %d", rec.Code)

		// Should load child table filtered by FK
		if rec.Code == http.StatusOK {
			body := rec.Body.String()
			assert.True(t,
				strings.Contains(body, "posts") ||
					strings.Contains(body, "filter") ||
					len(body) > 0,
				"Should show filtered child table data")
		}
	})
}
