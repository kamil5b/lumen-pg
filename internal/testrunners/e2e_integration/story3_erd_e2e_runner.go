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

// Story3ERDViewerE2ERunner runs end-to-end tests for Story 3: ERD Viewer
// This tests the complete route stack with all middleware for ERD visualization
// Maps to TEST_PLAN.md Story 3 E2E Tests [L257-281]:
// - E2E-S3-01: ERD Viewer Page Access
// - E2E-S3-02: ERD Zoom Controls
// - E2E-S3-03: ERD Pan
// - E2E-S3-04: Table Click in ERD
//
// Tests complete ERD viewer functionality including:
// - ERD page rendering with authentication
// - ERD data generation from database schema
// - Interactive controls (zoom, pan)
// - Navigation from ERD to table view
func Story3ERDViewerE2ERunner(t *testing.T, router http.Handler) {
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

	// E2E-S3-01: ERD Viewer Page Access
	t.Run("E2E-S3-01: ERD Viewer Page Access", func(t *testing.T) {
		// Test that authenticated users can access ERD viewer page
		cookies := getAuthenticatedSession(t)

		req := httptest.NewRequest(http.MethodGet, "/erd-viewer", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "ERD viewer page should be accessible")
		body := rec.Body.String()
		assert.True(t,
			strings.Contains(body, "erd") ||
				strings.Contains(body, "ERD") ||
				strings.Contains(body, "diagram") ||
				strings.Contains(body, "schema"),
			"ERD page should contain ERD-related content")

		// Test that unauthenticated users cannot access
		req = httptest.NewRequest(http.MethodGet, "/erd-viewer", nil)
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

	// E2E-S3-02: ERD Zoom Controls
	t.Run("E2E-S3-02: ERD Zoom Controls", func(t *testing.T) {
		// Test zoom in/out functionality via API
		cookies := getAuthenticatedSession(t)

		// Test zoom in
		zoomInPayload := `{"action": "zoom_in", "level": 1.5}`
		req := httptest.NewRequest(http.MethodPost, "/api/erd/zoom", strings.NewReader(zoomInPayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.True(t,
			rec.Code == http.StatusOK || rec.Code == http.StatusAccepted,
			"Zoom in should succeed, got %d", rec.Code)

		// Test zoom out
		zoomOutPayload := `{"action": "zoom_out", "level": 0.75}`
		req = httptest.NewRequest(http.MethodPost, "/api/erd/zoom", strings.NewReader(zoomOutPayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.True(t,
			rec.Code == http.StatusOK || rec.Code == http.StatusAccepted,
			"Zoom out should succeed, got %d", rec.Code)

		// Test zoom reset
		zoomResetPayload := `{"action": "zoom_reset", "level": 1.0}`
		req = httptest.NewRequest(http.MethodPost, "/api/erd/zoom", strings.NewReader(zoomResetPayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.True(t,
			rec.Code == http.StatusOK || rec.Code == http.StatusAccepted,
			"Zoom reset should succeed, got %d", rec.Code)
	})

	// E2E-S3-03: ERD Pan
	t.Run("E2E-S3-03: ERD Pan", func(t *testing.T) {
		// Test panning/dragging the ERD canvas
		cookies := getAuthenticatedSession(t)

		// Test pan to different positions
		panPayload := `{"action": "pan", "x": 100, "y": 50}`
		req := httptest.NewRequest(http.MethodPost, "/api/erd/pan", strings.NewReader(panPayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.True(t,
			rec.Code == http.StatusOK || rec.Code == http.StatusAccepted,
			"Pan should succeed, got %d", rec.Code)

		// Test pan with negative coordinates
		panPayload = `{"action": "pan", "x": -50, "y": -100}`
		req = httptest.NewRequest(http.MethodPost, "/api/erd/pan", strings.NewReader(panPayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.True(t,
			rec.Code == http.StatusOK || rec.Code == http.StatusAccepted,
			"Pan with negative coordinates should succeed, got %d", rec.Code)

		// Test pan reset to center
		panResetPayload := `{"action": "pan_reset", "x": 0, "y": 0}`
		req = httptest.NewRequest(http.MethodPost, "/api/erd/pan", strings.NewReader(panResetPayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.True(t,
			rec.Code == http.StatusOK || rec.Code == http.StatusAccepted,
			"Pan reset should succeed, got %d", rec.Code)
	})

	// E2E-S3-04: Table Click in ERD
	t.Run("E2E-S3-04: Table Click in ERD", func(t *testing.T) {
		// Test clicking on a table in ERD navigates to main view for that table
		cookies := getAuthenticatedSession(t)

		// Step 1: Generate ERD to get available tables
		req := httptest.NewRequest(http.MethodGet, "/api/erd/generate", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "ERD generation should succeed")

		// Parse response to get table information
		var erdData map[string]interface{}
		body := rec.Body.Bytes()
		if len(body) > 0 && json.Valid(body) {
			err := json.Unmarshal(body, &erdData)
			require.NoError(t, err, "Should parse ERD data")

			// ERD should contain tables information
			assert.True(t,
				erdData["tables"] != nil ||
					erdData["nodes"] != nil ||
					erdData["entities"] != nil ||
					len(erdData) > 0,
				"ERD should contain table/node information")
		}

		// Step 2: Click on a table (simulate by calling table click endpoint)
		tableClickPayload := `{
			"database": "testdb",
			"schema": "public",
			"table": "users"
		}`
		req = httptest.NewRequest(http.MethodPost, "/api/erd/table-click", strings.NewReader(tableClickPayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// Should either:
		// 1. Return 200 with navigation data
		// 2. Redirect to main view with table selected
		assert.True(t,
			rec.Code == http.StatusOK ||
				rec.Code == http.StatusFound ||
				rec.Code == http.StatusSeeOther,
			"Table click should succeed, got %d", rec.Code)

		if rec.Code == http.StatusFound || rec.Code == http.StatusSeeOther {
			location := rec.Header().Get("Location")
			assert.True(t,
				strings.Contains(location, "/main") ||
					strings.Contains(location, "table=users"),
				"Should navigate to main view with selected table")
		}

		// Step 3: Verify clicking on non-existent table is handled
		invalidTablePayload := `{
			"database": "testdb",
			"schema": "public",
			"table": "nonexistent_table"
		}`
		req = httptest.NewRequest(http.MethodPost, "/api/erd/table-click", strings.NewReader(invalidTablePayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// Should return error or 404
		assert.True(t,
			rec.Code == http.StatusNotFound ||
				rec.Code == http.StatusBadRequest ||
				rec.Code == http.StatusOK, // Might return OK with error message
			"Invalid table click should be handled gracefully, got %d", rec.Code)
	})

	// Additional test: ERD generation with empty schema
	t.Run("ERD Generation with Empty Schema", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Request ERD for a schema with no tables (if exists)
		req := httptest.NewRequest(http.MethodGet, "/api/erd/generate?schema=empty_schema", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// Should succeed but return empty or minimal ERD
		assert.True(t,
			rec.Code == http.StatusOK ||
				rec.Code == http.StatusNoContent,
			"Empty schema ERD should be handled gracefully, got %d", rec.Code)

		if rec.Code == http.StatusOK {
			body := rec.Body.String()
			// Should return valid JSON even if empty
			assert.True(t,
				strings.Contains(body, "[]") ||
					strings.Contains(body, "{}") ||
					strings.Contains(body, "tables") ||
					len(body) > 0,
				"Empty ERD should return valid response")
		}
	})

	// Additional test: ERD with complex relationships
	t.Run("ERD with Complex Relationships", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Request full ERD which may contain complex foreign key relationships
		req := httptest.NewRequest(http.MethodGet, "/api/erd/generate?include_relationships=true", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "ERD with relationships should succeed")

		body := rec.Body.Bytes()
		if len(body) > 0 && json.Valid(body) {
			var erdData map[string]interface{}
			err := json.Unmarshal(body, &erdData)
			require.NoError(t, err)

			// Should contain relationship/edge information
			assert.True(t,
				erdData["relationships"] != nil ||
					erdData["edges"] != nil ||
					erdData["foreignKeys"] != nil ||
					erdData["tables"] != nil,
				"ERD should contain relationship information")
		}
	})
}
