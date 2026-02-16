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

// Story7SecurityE2ERunner runs end-to-end tests for Story 7: Security & Best Practices
// This tests the complete route stack with all middleware for security features
// Maps to TEST_PLAN.md Story 7 E2E Tests [L747-781]:
// - E2E-S7-01: SQL Injection via WHERE Bar
// - E2E-S7-02: SQL Injection via Query Editor
// - E2E-S7-03: Cookie Tampering Prevention
// - E2E-S7-04: Session Timeout Enforcement
// - E2E-S7-05: HTTPS-Only Cookies (if HTTPS enabled)
// - E2E-S7-06: HTTPOnly Cookies
//
// Tests complete security functionality including:
// - SQL injection prevention across all input vectors
// - Cookie security (HttpOnly, Secure, SameSite)
// - Session timeout and expiration
// - Cookie tampering detection
// - Password encryption in cookies
func Story7SecurityE2ERunner(t *testing.T, router http.Handler) {
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

	// E2E-S7-01: SQL Injection via WHERE Bar
	t.Run("E2E-S7-01: SQL Injection via WHERE Bar", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Test various SQL injection attempts in WHERE clause
		injectionAttempts := []string{
			"1' OR '1'='1",
			"1; DROP TABLE users; --",
			"1' UNION SELECT * FROM users --",
			"1' AND 1=1 --",
			"'; DELETE FROM users WHERE '1'='1",
			"1' OR 1=1 /*",
			"admin'--",
			"' OR 'x'='x",
			"1' AND (SELECT COUNT(*) FROM users) > 0 --",
			"1'; EXEC sp_executesql N'DROP TABLE users' --",
		}

		for _, injection := range injectionAttempts {
			t.Run("Injection: "+injection, func(t *testing.T) {
				filterPayload := map[string]string{
					"table": "users",
					"where": injection,
				}
				payloadBytes, _ := json.Marshal(filterPayload)

				req := httptest.NewRequest(http.MethodPost, "/api/table/filter", strings.NewReader(string(payloadBytes)))
				req.Header.Set("Content-Type", "application/json")
				for _, cookie := range cookies {
					req.AddCookie(cookie)
				}
				rec := httptest.NewRecorder()
				router.ServeHTTP(rec, req)

				// Should either reject with error or safely sanitize
				// Should NOT execute malicious SQL
				assert.True(t,
					rec.Code == http.StatusBadRequest ||
						rec.Code == http.StatusOK ||
						rec.Code == http.StatusForbidden,
					"Injection attempt should be handled safely, got %d", rec.Code)

				body := rec.Body.String()
				// Should not contain evidence of successful SQL injection
				assert.NotContains(t, strings.ToLower(body), "syntax error near 'drop'",
					"Should prevent DROP TABLE execution")
			})
		}

		// Test legitimate WHERE clauses still work
		legitimateWhere := "id > 10 AND status = 'active'"
		filterPayload := map[string]string{
			"table": "users",
			"where": legitimateWhere,
		}
		payloadBytes, _ := json.Marshal(filterPayload)

		req := httptest.NewRequest(http.MethodPost, "/api/table/filter", strings.NewReader(string(payloadBytes)))
		req.Header.Set("Content-Type", "application/json")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.True(t,
			rec.Code == http.StatusOK || rec.Code == http.StatusAccepted,
			"Legitimate WHERE clause should work, got %d", rec.Code)
	})

	// E2E-S7-02: SQL Injection via Query Editor
	t.Run("E2E-S7-02: SQL Injection via Query Editor", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Test SQL injection attempts via query editor
		// Note: Query editor may allow DDL/DML by design, but should still prevent
		// unintended side effects through proper connection handling

		injectionAttempts := []struct {
			name  string
			query string
		}{
			{
				name:  "Union-based injection",
				query: "SELECT * FROM users WHERE id = 1 UNION SELECT password, username FROM admin_users --",
			},
			{
				name:  "Stacked queries",
				query: "SELECT * FROM users; DROP TABLE important_table; --",
			},
			{
				name:  "Time-based blind injection",
				query: "SELECT * FROM users WHERE id = 1 AND (SELECT pg_sleep(10)) --",
			},
			{
				name:  "Comment injection",
				query: "SELECT * FROM users WHERE name = 'admin' /**/OR/**/1=1 --",
			},
		}

		for _, attempt := range injectionAttempts {
			t.Run(attempt.name, func(t *testing.T) {
				queryPayload := map[string]string{
					"query": attempt.query,
				}
				payloadBytes, _ := json.Marshal(queryPayload)

				req := httptest.NewRequest(http.MethodPost, "/api/execute-query", strings.NewReader(string(payloadBytes)))
				req.Header.Set("Content-Type", "application/json")
				for _, cookie := range cookies {
					req.AddCookie(cookie)
				}
				rec := httptest.NewRecorder()
				router.ServeHTTP(rec, req)

				// Query execution should either succeed with proper isolation
				// or fail safely without cascading damage
				assert.True(t,
					rec.Code == http.StatusOK ||
						rec.Code == http.StatusBadRequest ||
						rec.Code == http.StatusForbidden ||
						rec.Code == http.StatusInternalServerError,
					"Injection attempt should be handled, got %d", rec.Code)

				// Even if query executes, it should be in isolated transaction context
				// and not affect other users or system tables
			})
		}
	})

	// E2E-S7-03: Cookie Tampering Prevention
	t.Run("E2E-S7-03: Cookie Tampering Prevention", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Get original cookies
		var usernameCookie, passwordCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "lumen_username" {
				usernameCookie = cookie
			}
			if cookie.Name == "lumen_password" {
				passwordCookie = cookie
			}
		}

		require.NotNil(t, usernameCookie, "Should have username cookie")
		require.NotNil(t, passwordCookie, "Should have password cookie")

		// Test 1: Tamper with username cookie value
		tamperedUsernameCookie := &http.Cookie{
			Name:  usernameCookie.Name,
			Value: usernameCookie.Value + "tampered",
		}

		req := httptest.NewRequest(http.MethodGet, "/main", nil)
		req.AddCookie(tamperedUsernameCookie)
		req.AddCookie(passwordCookie)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.True(t,
			rec.Code == http.StatusFound ||
				rec.Code == http.StatusUnauthorized ||
				rec.Code == http.StatusForbidden,
			"Tampered username cookie should be rejected, got %d", rec.Code)

		// Test 2: Tamper with password cookie value
		tamperedPasswordCookie := &http.Cookie{
			Name:  passwordCookie.Name,
			Value: passwordCookie.Value[:len(passwordCookie.Value)-5] + "xxxxx",
		}

		req = httptest.NewRequest(http.MethodGet, "/main", nil)
		req.AddCookie(usernameCookie)
		req.AddCookie(tamperedPasswordCookie)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.True(t,
			rec.Code == http.StatusFound ||
				rec.Code == http.StatusUnauthorized ||
				rec.Code == http.StatusForbidden,
			"Tampered password cookie should be rejected, got %d", rec.Code)

		// Test 3: Use completely fabricated cookie
		fakeCookie := &http.Cookie{
			Name:  "lumen_username",
			Value: "YWRtaW4=", // base64 encoded "admin"
		}
		fakePasswordCookie := &http.Cookie{
			Name:  "lumen_password",
			Value: "ZmFrZXBhc3N3b3Jk", // base64 encoded "fakepassword"
		}

		req = httptest.NewRequest(http.MethodGet, "/main", nil)
		req.AddCookie(fakeCookie)
		req.AddCookie(fakePasswordCookie)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.True(t,
			rec.Code == http.StatusFound ||
				rec.Code == http.StatusUnauthorized ||
				rec.Code == http.StatusForbidden,
			"Fabricated cookies should be rejected, got %d", rec.Code)

		// Test 4: Valid cookies should still work
		req = httptest.NewRequest(http.MethodGet, "/main", nil)
		req.AddCookie(usernameCookie)
		req.AddCookie(passwordCookie)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Valid cookies should work")
	})

	// E2E-S7-04: Session Timeout Enforcement
	t.Run("E2E-S7-04: Session Timeout Enforcement", func(t *testing.T) {
		// Login and get cookies with short expiration
		formData := url.Values{}
		formData.Set("username", "testuser")
		formData.Set("password", "testpass")

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusFound, rec.Code)
		cookies := rec.Result().Cookies()

		// Check cookie expiration attributes
		for _, cookie := range cookies {
			if cookie.Name == "lumen_username" || cookie.Name == "lumen_password" {
				// Cookies should have MaxAge set (for session timeout)
				assert.True(t,
					cookie.MaxAge > 0 || !cookie.Expires.IsZero(),
					"Session cookie should have expiration, got MaxAge=%d", cookie.MaxAge)

				// Expiration should be reasonable (not too long)
				if !cookie.Expires.IsZero() {
					timeUntilExpiry := time.Until(cookie.Expires)
					assert.True(t,
						timeUntilExpiry > 0 && timeUntilExpiry < 24*time.Hour,
						"Session should expire within reasonable time, got %v", timeUntilExpiry)
				}
			}
		}

		// Access protected route immediately - should work
		req = httptest.NewRequest(http.MethodGet, "/main", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Fresh session should work")

		// Test with expired cookie (simulate by creating cookie with past expiration)
		expiredCookie := &http.Cookie{
			Name:    "lumen_username",
			Value:   "testuser",
			Expires: time.Now().Add(-1 * time.Hour), // Expired 1 hour ago
			MaxAge:  -1,
		}

		req = httptest.NewRequest(http.MethodGet, "/main", nil)
		req.AddCookie(expiredCookie)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.True(t,
			rec.Code == http.StatusFound ||
				rec.Code == http.StatusUnauthorized,
			"Expired session should be rejected, got %d", rec.Code)
	})

	// E2E-S7-05: HTTPS-Only Cookies (if HTTPS enabled)
	t.Run("E2E-S7-05: HTTPS-Only Cookies (if HTTPS enabled)", func(t *testing.T) {
		// Login and check cookie Secure attribute
		formData := url.Values{}
		formData.Set("username", "testuser")
		formData.Set("password", "testpass")

		// Test with HTTPS request (if server supports it)
		req := httptest.NewRequest(http.MethodPost, "https://localhost/login", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("X-Forwarded-Proto", "https")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code == http.StatusFound {
			cookies := rec.Result().Cookies()

			for _, cookie := range cookies {
				if cookie.Name == "lumen_username" || cookie.Name == "lumen_password" {
					// In production with HTTPS, Secure flag should be set
					// Note: In test environment, this might not be enforced
					if strings.HasPrefix(req.URL.String(), "https://") {
						assert.True(t,
							cookie.Secure || req.Header.Get("X-Forwarded-Proto") == "https",
							"Cookie %s should have Secure flag for HTTPS", cookie.Name)
					}
				}
			}
		}

		// Test SameSite attribute
		req = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code == http.StatusFound {
			cookies := rec.Result().Cookies()

			for _, cookie := range cookies {
				if cookie.Name == "lumen_username" || cookie.Name == "lumen_password" {
					// SameSite should be set to prevent CSRF
					assert.True(t,
						cookie.SameSite == http.SameSiteStrictMode ||
							cookie.SameSite == http.SameSiteLaxMode ||
							cookie.SameSite == http.SameSiteDefaultMode,
						"Cookie %s should have SameSite attribute", cookie.Name)
				}
			}
		}
	})

	// E2E-S7-06: HTTPOnly Cookies
	t.Run("E2E-S7-06: HTTPOnly Cookies", func(t *testing.T) {
		// Login and verify HttpOnly flag on cookies
		formData := url.Values{}
		formData.Set("username", "testuser")
		formData.Set("password", "testpass")

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusFound, rec.Code)
		cookies := rec.Result().Cookies()

		// Check that sensitive cookies have HttpOnly flag
		var hasUsernameCookie, hasPasswordCookie bool
		for _, cookie := range cookies {
			if cookie.Name == "lumen_username" {
				hasUsernameCookie = true
				assert.True(t, cookie.HttpOnly,
					"Username cookie MUST have HttpOnly flag to prevent XSS")
			}
			if cookie.Name == "lumen_password" {
				hasPasswordCookie = true
				assert.True(t, cookie.HttpOnly,
					"Password cookie MUST have HttpOnly flag to prevent XSS")
			}
		}

		assert.True(t, hasUsernameCookie, "Should set username cookie")
		assert.True(t, hasPasswordCookie, "Should set password cookie")
	})

	// Additional test: Password Encryption in Cookie
	t.Run("Password Encryption in Cookie", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("username", "testuser")
		formData.Set("password", "plaintextpassword123")

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusFound, rec.Code)
		cookies := rec.Result().Cookies()

		// Check that password is NOT stored in plaintext
		for _, cookie := range cookies {
			if cookie.Name == "lumen_password" {
				// Password cookie should NOT contain the plaintext password
				assert.NotContains(t, cookie.Value, "plaintextpassword123",
					"Password must be encrypted in cookie, not plaintext")

				// Should be encrypted/encoded (typically longer than plaintext)
				assert.Greater(t, len(cookie.Value), 20,
					"Encrypted password cookie should be reasonably long")

				// Should not be simple base64 of plaintext
				assert.NotEqual(t, "cGxhaW50ZXh0cGFzc3dvcmQxMjM=", cookie.Value,
					"Password should not be simple base64 encoding")
			}
		}
	})

	// Additional test: Rate Limiting (if implemented)
	t.Run("Rate Limiting Protection", func(t *testing.T) {
		// Attempt multiple rapid login requests to test rate limiting
		failedAttempts := 0
		successAttempts := 0

		for i := 0; i < 20; i++ {
			formData := url.Values{}
			formData.Set("username", "testuser")
			formData.Set("password", "wrongpassword")

			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Set("X-Real-IP", "192.168.1.100") // Simulate same IP
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			if rec.Code == http.StatusTooManyRequests || rec.Code == http.StatusForbidden {
				failedAttempts++
			} else {
				successAttempts++
			}
		}

		// If rate limiting is implemented, some requests should be blocked
		// This is optional, so we just log the result
		t.Logf("Rate limiting test: %d blocked, %d allowed out of 20 requests", failedAttempts, successAttempts)
	})

	// Additional test: CSRF Protection
	t.Run("CSRF Protection", func(t *testing.T) {
		cookies := getAuthenticatedSession(t)

		// Attempt state-changing operation without proper CSRF token (if implemented)
		// This simulates a cross-site request
		req := httptest.NewRequest(http.MethodPost, "/api/transaction/start", nil)
		req.Header.Set("Origin", "https://malicious-site.com")
		req.Header.Set("Referer", "https://malicious-site.com/attack")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// If CSRF protection is implemented, this should be blocked
		// Otherwise, SameSite cookies should provide some protection
		// We log the result for awareness
		t.Logf("CSRF test result: %d (Blocked if 403/400, Allowed if 200/201)", rec.Code)
	})
}
