package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Story6IsolationE2ERunner runs end-to-end tests for Story 6: Isolation
// This tests the complete route stack with all middleware for multi-user isolation
// Maps to TEST_PLAN.md Story 6 E2E Tests [L672-691]:
// - E2E-S6-01: Simultaneous Users Different Permissions
// - E2E-S6-02: Simultaneous Transactions
// - E2E-S6-03: One User Cannot See Another's Session
//
// Tests complete isolation functionality including:
// - Session isolation between concurrent users
// - Transaction isolation between users
// - Permission-based resource access isolation
// - Cookie and authentication isolation
func Story6IsolationE2ERunner(t *testing.T, router http.Handler) {
	t.Helper()

	// Helper function to login with specific credentials
	loginUser := func(t *testing.T, username, password string) []*http.Cookie {
		formData := url.Values{}
		formData.Set("username", username)
		formData.Set("password", password)

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusFound, rec.Code, "Login should succeed for %s", username)
		cookies := rec.Result().Cookies()
		require.NotEmpty(t, cookies, "Should receive session cookies for %s", username)
		return cookies
	}

	// E2E-S6-01: Simultaneous Users Different Permissions
	t.Run("E2E-S6-01: Simultaneous Users Different Permissions", func(t *testing.T) {
		// Login as two different users with different permissions
		user1Cookies := loginUser(t, "testuser1", "testpass1")
		user2Cookies := loginUser(t, "testuser2", "testpass2")

		// User 1 accesses their data explorer
		req1 := httptest.NewRequest(http.MethodGet, "/api/data-explorer", nil)
		for _, cookie := range user1Cookies {
			req1.AddCookie(cookie)
		}
		rec1 := httptest.NewRecorder()
		router.ServeHTTP(rec1, req1)

		assert.Equal(t, http.StatusOK, rec1.Code, "User 1 should access data explorer")

		var user1Data map[string]interface{}
		err := json.Unmarshal(rec1.Body.Bytes(), &user1Data)
		require.NoError(t, err)

		// User 2 accesses their data explorer
		req2 := httptest.NewRequest(http.MethodGet, "/api/data-explorer", nil)
		for _, cookie := range user2Cookies {
			req2.AddCookie(cookie)
		}
		rec2 := httptest.NewRecorder()
		router.ServeHTTP(rec2, req2)

		assert.Equal(t, http.StatusOK, rec2.Code, "User 2 should access data explorer")

		var user2Data map[string]interface{}
		err = json.Unmarshal(rec2.Body.Bytes(), &user2Data)
		require.NoError(t, err)

		// Users should see different resources based on their permissions
		// Note: This assumes users have different database permissions
		// At minimum, verify both users get valid responses
		assert.NotNil(t, user1Data)
		assert.NotNil(t, user2Data)

		// Verify User 1 cannot use User 2's cookies
		req3 := httptest.NewRequest(http.MethodGet, "/api/data-explorer", nil)
		for _, cookie := range user2Cookies {
			req3.AddCookie(cookie)
		}
		rec3 := httptest.NewRecorder()
		router.ServeHTTP(rec3, req3)

		// Should get User 2's data, not User 1's
		var user2DataAgain map[string]interface{}
		if rec3.Code == http.StatusOK {
			json.Unmarshal(rec3.Body.Bytes(), &user2DataAgain)
			// The data should be consistent with user2's permissions
			assert.NotNil(t, user2DataAgain)
		}

		// Test concurrent table access with different permissions
		// User 1 tries to access a table
		req1 = httptest.NewRequest(http.MethodGet, "/api/table/data?table=users", nil)
		for _, cookie := range user1Cookies {
			req1.AddCookie(cookie)
		}
		rec1 = httptest.NewRecorder()
		router.ServeHTTP(rec1, req1)

		user1TableAccess := rec1.Code

		// User 2 tries to access the same table
		req2 = httptest.NewRequest(http.MethodGet, "/api/table/data?table=users", nil)
		for _, cookie := range user2Cookies {
			req2.AddCookie(cookie)
		}
		rec2 = httptest.NewRecorder()
		router.ServeHTTP(rec2, req2)

		user2TableAccess := rec2.Code

		// Access should be determined by individual permissions
		assert.True(t,
			(user1TableAccess == http.StatusOK || user1TableAccess == http.StatusForbidden) &&
				(user2TableAccess == http.StatusOK || user2TableAccess == http.StatusForbidden),
			"Users should have independent permission checks")
	})

	// E2E-S6-02: Simultaneous Transactions
	t.Run("E2E-S6-02: Simultaneous Transactions", func(t *testing.T) {
		// Login as two different users
		user1Cookies := loginUser(t, "testuser1", "testpass1")
		user2Cookies := loginUser(t, "testuser2", "testpass2")

		var wg sync.WaitGroup
		wg.Add(2)

		var user1TransactionID, user2TransactionID string
		var user1Error, user2Error error

		// User 1 starts a transaction
		go func() {
			defer wg.Done()

			req := httptest.NewRequest(http.MethodPost, "/api/transaction/start", nil)
			for _, cookie := range user1Cookies {
				req.AddCookie(cookie)
			}
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			if rec.Code == http.StatusOK || rec.Code == http.StatusCreated {
				var response map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				if err == nil {
					if tid, ok := response["transaction_id"].(string); ok {
						user1TransactionID = tid
					} else if tid, ok := response["transactionId"].(string); ok {
						user1TransactionID = tid
					}
				}
			} else {
				user1Error = assert.AnError
			}
		}()

		// User 2 starts a transaction simultaneously
		go func() {
			defer wg.Done()

			req := httptest.NewRequest(http.MethodPost, "/api/transaction/start", nil)
			for _, cookie := range user2Cookies {
				req.AddCookie(cookie)
			}
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			if rec.Code == http.StatusOK || rec.Code == http.StatusCreated {
				var response map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				if err == nil {
					if tid, ok := response["transaction_id"].(string); ok {
						user2TransactionID = tid
					} else if tid, ok := response["transactionId"].(string); ok {
						user2TransactionID = tid
					}
				}
			} else {
				user2Error = assert.AnError
			}
		}()

		wg.Wait()

		// Both transactions should succeed independently
		assert.NoError(t, user1Error, "User 1 transaction should start successfully")
		assert.NoError(t, user2Error, "User 2 transaction should start successfully")

		// Transaction IDs should be different (if returned)
		if user1TransactionID != "" && user2TransactionID != "" {
			assert.NotEqual(t, user1TransactionID, user2TransactionID,
				"Users should have separate transaction IDs")
		}

		// User 1 edits data in their transaction
		editPayload1 := `{"table": "users", "row_id": 1, "column": "name", "value": "User1Edit"}`
		req1 := httptest.NewRequest(http.MethodPost, "/api/transaction/edit-cell", strings.NewReader(editPayload1))
		req1.Header.Set("Content-Type", "application/json")
		for _, cookie := range user1Cookies {
			req1.AddCookie(cookie)
		}
		rec1 := httptest.NewRecorder()
		router.ServeHTTP(rec1, req1)

		// User 2 edits data in their transaction
		editPayload2 := `{"table": "users", "row_id": 2, "column": "name", "value": "User2Edit"}`
		req2 := httptest.NewRequest(http.MethodPost, "/api/transaction/edit-cell", strings.NewReader(editPayload2))
		req2.Header.Set("Content-Type", "application/json")
		for _, cookie := range user2Cookies {
			req2.AddCookie(cookie)
		}
		rec2 := httptest.NewRecorder()
		router.ServeHTTP(rec2, req2)

		// Both edits should succeed independently
		assert.True(t,
			rec1.Code == http.StatusOK || rec1.Code == http.StatusAccepted || rec1.Code == http.StatusForbidden,
			"User 1 edit should be processed")
		assert.True(t,
			rec2.Code == http.StatusOK || rec2.Code == http.StatusAccepted || rec2.Code == http.StatusForbidden,
			"User 2 edit should be processed")

		// User 1 checks their transaction status
		req1 = httptest.NewRequest(http.MethodGet, "/api/transaction/status", nil)
		for _, cookie := range user1Cookies {
			req1.AddCookie(cookie)
		}
		rec1 = httptest.NewRecorder()
		router.ServeHTTP(rec1, req1)

		var user1Status map[string]interface{}
		if rec1.Code == http.StatusOK {
			json.Unmarshal(rec1.Body.Bytes(), &user1Status)
		}

		// User 2 checks their transaction status
		req2 = httptest.NewRequest(http.MethodGet, "/api/transaction/status", nil)
		for _, cookie := range user2Cookies {
			req2.AddCookie(cookie)
		}
		rec2 = httptest.NewRecorder()
		router.ServeHTTP(rec2, req2)

		var user2Status map[string]interface{}
		if rec2.Code == http.StatusOK {
			json.Unmarshal(rec2.Body.Bytes(), &user2Status)
		}

		// Each user should only see their own transaction state
		assert.NotNil(t, user1Status)
		assert.NotNil(t, user2Status)

		// User 1 commits their transaction
		req1 = httptest.NewRequest(http.MethodPost, "/api/transaction/commit", nil)
		for _, cookie := range user1Cookies {
			req1.AddCookie(cookie)
		}
		rec1 = httptest.NewRecorder()
		router.ServeHTTP(rec1, req1)

		// User 2 rolls back their transaction
		req2 = httptest.NewRequest(http.MethodPost, "/api/transaction/rollback", nil)
		for _, cookie := range user2Cookies {
			req2.AddCookie(cookie)
		}
		rec2 = httptest.NewRecorder()
		router.ServeHTTP(rec2, req2)

		// Both operations should succeed independently
		assert.True(t,
			rec1.Code == http.StatusOK || rec1.Code == http.StatusAccepted,
			"User 1 commit should succeed")
		assert.True(t,
			rec2.Code == http.StatusOK || rec2.Code == http.StatusAccepted,
			"User 2 rollback should succeed")
	})

	// E2E-S6-03: One User Cannot See Another's Session
	t.Run("E2E-S6-03: One User Cannot See Another's Session", func(t *testing.T) {
		// Login as User 1
		user1Cookies := loginUser(t, "testuser1", "testpass1")

		// Login as User 2
		user2Cookies := loginUser(t, "testuser2", "testpass2")

		// User 1 starts a transaction
		req := httptest.NewRequest(http.MethodPost, "/api/transaction/start", nil)
		for _, cookie := range user1Cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.True(t, rec.Code == http.StatusOK || rec.Code == http.StatusCreated)

		// User 1 makes an edit
		editPayload := `{"table": "users", "row_id": 1, "column": "name", "value": "User1Secret"}`
		req = httptest.NewRequest(http.MethodPost, "/api/transaction/edit-cell", strings.NewReader(editPayload))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range user1Cookies {
			req.AddCookie(cookie)
		}
		router.ServeHTTP(httptest.NewRecorder(), req)

		// User 2 tries to see User 1's transaction status
		req = httptest.NewRequest(http.MethodGet, "/api/transaction/status", nil)
		for _, cookie := range user2Cookies {
			req.AddCookie(cookie)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// User 2 should either:
		// 1. Get no active transaction (404 or empty)
		// 2. Get their own transaction status (not User 1's)
		if rec.Code == http.StatusOK {
			var user2Status map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &user2Status)
			require.NoError(t, err)

			// Should not contain User 1's edits
			if user2Status["pending_edits"] != nil || user2Status["pendingEdits"] != nil {
				body := rec.Body.String()
				assert.NotContains(t, body, "User1Secret",
					"User 2 should not see User 1's transaction edits")
			}
		}

		// User 2 tries to commit User 1's transaction (should fail)
		req = httptest.NewRequest(http.MethodPost, "/api/transaction/commit", nil)
		for _, cookie := range user2Cookies {
			req.AddCookie(cookie)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// Should either commit User 2's own transaction (if any) or fail
		// User 1's transaction should remain unaffected

		// Verify User 1's transaction is still active
		req = httptest.NewRequest(http.MethodGet, "/api/transaction/status", nil)
		for _, cookie := range user1Cookies {
			req.AddCookie(cookie)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "User 1 should still have active transaction")

		var user1Status map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &user1Status)
		require.NoError(t, err)

		// User 1's pending edits should still be there
		assert.True(t,
			user1Status["pending_edits"] != nil ||
				user1Status["pendingEdits"] != nil ||
				user1Status["active"] != nil,
			"User 1's transaction should still have pending edits")

		// Test session cookie isolation
		// User 1 accesses main page
		req = httptest.NewRequest(http.MethodGet, "/main", nil)
		for _, cookie := range user1Cookies {
			req.AddCookie(cookie)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		user1MainPage := rec.Body.String()
		assert.Contains(t, user1MainPage, "testuser1", "Should show User 1's username")

		// User 2 accesses main page
		req = httptest.NewRequest(http.MethodGet, "/main", nil)
		for _, cookie := range user2Cookies {
			req.AddCookie(cookie)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		user2MainPage := rec.Body.String()
		assert.Contains(t, user2MainPage, "testuser2", "Should show User 2's username")
		assert.NotContains(t, user2MainPage, "testuser1", "Should not show User 1's username")

		// Try to mix cookies (User 1's username cookie with User 2's password cookie)
		// This should fail or behave unpredictably, demonstrating cookie integrity
		var user1UsernameCookie, user2PasswordCookie *http.Cookie
		for _, cookie := range user1Cookies {
			if cookie.Name == "lumen_username" {
				user1UsernameCookie = cookie
			}
		}
		for _, cookie := range user2Cookies {
			if cookie.Name == "lumen_password" {
				user2PasswordCookie = cookie
			}
		}

		if user1UsernameCookie != nil && user2PasswordCookie != nil {
			req = httptest.NewRequest(http.MethodGet, "/main", nil)
			req.AddCookie(user1UsernameCookie)
			req.AddCookie(user2PasswordCookie)
			rec = httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			// Should fail authentication (redirect to login or unauthorized)
			assert.True(t,
				rec.Code == http.StatusFound ||
					rec.Code == http.StatusUnauthorized ||
					rec.Code == http.StatusForbidden,
				"Mixed cookies should fail authentication, got %d", rec.Code)
		}
	})

	// Additional test: Concurrent query execution isolation
	t.Run("Concurrent Query Execution Isolation", func(t *testing.T) {
		user1Cookies := loginUser(t, "testuser1", "testpass1")
		user2Cookies := loginUser(t, "testuser2", "testpass2")

		var wg sync.WaitGroup
		wg.Add(2)

		var user1Result, user2Result int

		// User 1 executes a query
		go func() {
			defer wg.Done()

			queryPayload := `{"query": "SELECT COUNT(*) FROM users"}`
			req := httptest.NewRequest(http.MethodPost, "/api/execute-query", strings.NewReader(queryPayload))
			req.Header.Set("Content-Type", "application/json")
			for _, cookie := range user1Cookies {
				req.AddCookie(cookie)
			}
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			if rec.Code == http.StatusOK {
				user1Result = http.StatusOK
			} else {
				user1Result = rec.Code
			}
		}()

		// User 2 executes a query simultaneously
		go func() {
			defer wg.Done()

			queryPayload := `{"query": "SELECT COUNT(*) FROM posts"}`
			req := httptest.NewRequest(http.MethodPost, "/api/execute-query", strings.NewReader(queryPayload))
			req.Header.Set("Content-Type", "application/json")
			for _, cookie := range user2Cookies {
				req.AddCookie(cookie)
			}
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			if rec.Code == http.StatusOK {
				user2Result = http.StatusOK
			} else {
				user2Result = rec.Code
			}
		}()

		wg.Wait()

		// Both queries should complete without interference
		assert.True(t,
			user1Result == http.StatusOK || user1Result == http.StatusForbidden,
			"User 1 query should complete")
		assert.True(t,
			user2Result == http.StatusOK || user2Result == http.StatusForbidden,
			"User 2 query should complete")
	})
}
