package e2e_integration

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Story2AuthE2ERunner runs end-to-end tests for Story 2: Authentication & Identity
// This tests the complete route stack with all middleware (auth, logging, security, etc.)
// Maps to TEST_PLAN.md Story 2 E2E Tests [L187-221]:
// - E2E-S2-01: Login Flow with Connection Probe
// - E2E-S2-02: Login Flow - No Accessible Resources
// - E2E-S2-03: Login Flow - Invalid Credentials
// - E2E-S2-04: Logout Flow
// - E2E-S2-05: Protected Route Access Without Auth
// - E2E-S2-06: Data Explorer Populated After Login
//
// Unlike handler test runners which mock use cases, this runner tests:
// - Real HTTP routes with full middleware stack
// - Cookie handling across requests
// - Redirects and session flow
// - Complete authentication/authorization chain
func Story2AuthE2ERunner(t *testing.T, router http.Handler) {
	t.Helper()

	// E2E-S2-01: Login Flow with Connection Probe
	t.Run("E2E-S2-01: Login Flow with Connection Probe", func(t *testing.T) {
		// Expected flow:
		// 1. GET /login -> returns login page
		// 2. POST /login with credentials -> probes connection
		// 3. If probe succeeds -> redirects to /main with session cookies
		// 4. GET /main with cookies -> shows main view with username

		// Step 1: GET login page
		req := httptest.NewRequest(http.MethodGet, "/login", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()
		assert.Contains(t, body, "login", "Login page should contain login form")

		// Step 2: POST login with valid credentials
		formData := url.Values{}
		formData.Set("username", "testuser")
		formData.Set("password", "testpass")

		req = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// Should redirect to main after successful login
		assert.Equal(t, http.StatusFound, rec.Code)
		assert.Equal(t, "/main", rec.Header().Get("Location"))

		// Should set session cookies
		cookies := rec.Result().Cookies()
		assert.NotEmpty(t, cookies, "Should set session cookies")

		var usernameCookie, passwordCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "lumen_username" {
				usernameCookie = cookie
			}
			if cookie.Name == "lumen_password" {
				passwordCookie = cookie
			}
		}

		require.NotNil(t, usernameCookie, "Should set username cookie")
		require.NotNil(t, passwordCookie, "Should set encrypted password cookie")
		assert.True(t, usernameCookie.HttpOnly, "Username cookie should be HttpOnly")
		assert.True(t, passwordCookie.HttpOnly, "Password cookie should be HttpOnly")

		// Step 3: Follow redirect to /main with session cookies
		req = httptest.NewRequest(http.MethodGet, "/main", nil)
		req.AddCookie(usernameCookie)
		req.AddCookie(passwordCookie)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		body = rec.Body.String()
		assert.Contains(t, body, "testuser", "Main page should display username")
	})

	// E2E-S2-02: Login Flow - No Accessible Resources
	t.Run("E2E-S2-02: Login Flow - No Accessible Resources", func(t *testing.T) {
		// User with valid credentials but no database permissions
		formData := url.Values{}
		formData.Set("username", "noaccessuser")
		formData.Set("password", "testpass")

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// Should NOT redirect - stays on login page with error
		assert.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()
		assert.Contains(t, body, "no accessible resources", "Should show error message")
		assert.Contains(t, body, "login", "Should still show login form")

		// Should NOT set session cookies
		cookies := rec.Result().Cookies()
		for _, cookie := range cookies {
			assert.NotEqual(t, "lumen_username", cookie.Name)
			assert.NotEqual(t, "lumen_password", cookie.Name)
		}
	})

	// E2E-S2-03: Login Flow - Invalid Credentials
	t.Run("E2E-S2-03: Login Flow - Invalid Credentials", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("username", "invaliduser")
		formData.Set("password", "wrongpassword")

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// Should return error (401 or 200 with error message)
		assert.True(t, rec.Code == http.StatusUnauthorized || rec.Code == http.StatusOK)

		body := rec.Body.String()
		assert.True(t,
			strings.Contains(body, "invalid") ||
				strings.Contains(body, "incorrect") ||
				strings.Contains(body, "failed"),
			"Should show invalid credentials error")

		// Should NOT set session cookies
		cookies := rec.Result().Cookies()
		for _, cookie := range cookies {
			assert.NotEqual(t, "lumen_username", cookie.Name)
			assert.NotEqual(t, "lumen_password", cookie.Name)
		}
	})

	// E2E-S2-04: Logout Flow
	t.Run("E2E-S2-04: Logout Flow", func(t *testing.T) {
		// Step 1: Login first to get session cookies
		formData := url.Values{}
		formData.Set("username", "testuser")
		formData.Set("password", "testpass")

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusFound, rec.Code)
		loginCookies := rec.Result().Cookies()
		require.NotEmpty(t, loginCookies)

		// Step 2: POST to /logout
		req = httptest.NewRequest(http.MethodPost, "/logout", nil)
		for _, cookie := range loginCookies {
			req.AddCookie(cookie)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// Should redirect to login page
		assert.Equal(t, http.StatusFound, rec.Code)
		assert.Equal(t, "/login", rec.Header().Get("Location"))

		// Should clear session cookies (MaxAge=-1 or empty value)
		logoutCookies := rec.Result().Cookies()
		for _, cookie := range logoutCookies {
			if cookie.Name == "lumen_username" || cookie.Name == "lumen_password" {
				assert.True(t, cookie.MaxAge < 0 || cookie.Value == "",
					"Cookie %s should be cleared", cookie.Name)
			}
		}

		// Step 3: Try accessing protected route without valid cookies
		req = httptest.NewRequest(http.MethodGet, "/main", nil)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// Should redirect to login or return unauthorized
		assert.True(t,
			rec.Code == http.StatusFound ||
				rec.Code == http.StatusUnauthorized,
			"Should require authentication")

		if rec.Code == http.StatusFound {
			assert.Equal(t, "/login", rec.Header().Get("Location"))
		}
	})

	// E2E-S2-05: Protected Route Access Without Auth
	t.Run("E2E-S2-05: Protected Route Access Without Auth", func(t *testing.T) {
		protectedRoutes := []struct {
			method string
			path   string
		}{
			{http.MethodGet, "/main"},
			{http.MethodGet, "/query-editor"},
			{http.MethodGet, "/erd-viewer"},
			{http.MethodGet, "/api/data-explorer"},
			{http.MethodGet, "/api/tables"},
			{http.MethodPost, "/api/execute-query"},
			{http.MethodPost, "/api/transaction/start"},
			{http.MethodPost, "/api/transaction/commit"},
			{http.MethodGet, "/api/metadata/refresh"},
		}

		for _, route := range protectedRoutes {
			t.Run(route.method+" "+route.path, func(t *testing.T) {
				var body io.Reader
				if route.method == http.MethodPost {
					body = strings.NewReader("{}")
				}

				req := httptest.NewRequest(route.method, route.path, body)
				if route.method == http.MethodPost {
					req.Header.Set("Content-Type", "application/json")
				}
				rec := httptest.NewRecorder()
				router.ServeHTTP(rec, req)

				// Should require authentication
				assert.True(t,
					rec.Code == http.StatusFound ||
						rec.Code == http.StatusUnauthorized ||
						rec.Code == http.StatusForbidden,
					"Route %s %s should require authentication, got %d",
					route.method, route.path, rec.Code)

				// If redirecting, should redirect to login
				if rec.Code == http.StatusFound {
					location := rec.Header().Get("Location")
					assert.Contains(t, location, "/login",
						"Should redirect to login page for %s %s", route.method, route.path)
				}
			})
		}
	})

	// E2E-S2-06: Data Explorer Populated After Login
	t.Run("E2E-S2-06: Data Explorer Populated After Login", func(t *testing.T) {
		// Step 1: Login to get session
		formData := url.Values{}
		formData.Set("username", "testuser")
		formData.Set("password", "testpass")

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusFound, rec.Code)
		cookies := rec.Result().Cookies()
		require.NotEmpty(t, cookies)

		// Step 2: Request data explorer endpoint
		req = httptest.NewRequest(http.MethodGet, "/api/data-explorer", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()
		assert.NotEmpty(t, body, "Data explorer should return content")

		// Should contain database structure (databases, schemas, tables)
		// Accept JSON array or object structure
		assert.True(t,
			strings.Contains(body, "database") ||
				strings.Contains(body, "schema") ||
				strings.Contains(body, "table") ||
				strings.Contains(body, "[") ||
				strings.Contains(body, "{"),
			"Data explorer should contain metadata structure")

		// Step 3: Verify main page includes data explorer
		req = httptest.NewRequest(http.MethodGet, "/main", nil)
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		body = rec.Body.String()

		// Main page should contain user info or data explorer elements
		assert.True(t,
			strings.Contains(body, "testuser") ||
				strings.Contains(body, "data-explorer") ||
				strings.Contains(body, "sidebar"),
			"Main page should show authenticated user content")
	})
}
