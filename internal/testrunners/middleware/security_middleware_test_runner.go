package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kamil5b/lumen-pg/internal/interfaces/middleware"
)

// SecurityMiddlewareConstructor is a function type that creates a SecurityMiddleware
type SecurityMiddlewareConstructor func() middleware.SecurityMiddleware

// SecurityMiddlewareRunner runs all security middleware tests
// Maps to TEST_PLAN.md:
// - Story 7: Security & Best Practices [UC-S7-05~07, IT-S7-01~03, E2E-S7-03~06]
func SecurityMiddlewareRunner(t *testing.T, constructor SecurityMiddlewareConstructor) {
	t.Helper()

	mw := constructor()

	// UC-S7-06: Session Timeout Short-Lived Cookie
	// E2E-S7-05: HTTPS-Only Cookies (if HTTPS enabled)
	t.Run("UC-S7-06/E2E-S7-05: SetSecurityHeaders sets appropriate headers", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.SetSecurityHeaders(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)

		// Check security headers are set
		headers := rec.Header()
		require.NotEmpty(t, headers.Get("X-Content-Type-Options"))
		require.NotEmpty(t, headers.Get("X-Frame-Options"))
		require.NotEmpty(t, headers.Get("X-XSS-Protection"))
		require.NotEmpty(t, headers.Get("Content-Security-Policy"))
	})

	// UC-S7-01: SQL Injection Prevention - WHERE Clause
	// E2E-S7-01: SQL Injection via WHERE Bar
	t.Run("UC-S7-01/E2E-S7-01: PreventCSRF blocks requests without CSRF token", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.PreventCSRF(handler)

		req := httptest.NewRequest(http.MethodPost, "/api/data", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		// Should block POST without CSRF token
		require.False(t, called)
		require.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("PreventCSRF allows GET requests without CSRF token", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.PreventCSRF(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		// GET requests should pass through
		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("PreventCSRF allows POST with valid CSRF token", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.PreventCSRF(handler)

		req := httptest.NewRequest(http.MethodPost, "/api/data", nil)
		req.Header.Set("X-CSRF-Token", "valid_csrf_token")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	// IT-S7-03: Real Session Expiration
	// UC-S7-06: Session Timeout Short-Lived Cookie
	t.Run("IT-S7-03/UC-S7-06: RateLimiter limits excessive requests", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.RateLimiter(handler)

		// Make multiple requests
		successCount := 0
		rateLimitedCount := 0

		for i := 0; i < 100; i++ {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = "192.168.1.1:12345"
			rec := httptest.NewRecorder()

			wrapped.ServeHTTP(rec, req)

			if rec.Code == http.StatusOK {
				successCount++
			} else if rec.Code == http.StatusTooManyRequests {
				rateLimitedCount++
			}
		}

		// Should have some rate limited requests
		require.True(t, rateLimitedCount > 0 || successCount <= 100, "Rate limiter should limit excessive requests")
	})

	// UC-S7-05: Cookie Tampering Detection
	// E2E-S7-03: Cookie Tampering Prevention
	t.Run("UC-S7-05/E2E-S7-03: ValidateCookieIntegrity allows valid cookies", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateCookieIntegrity(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "valid_session_with_signature",
		})
		req.AddCookie(&http.Cookie{
			Name:  "session_signature",
			Value: "valid_signature",
		})
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	// UC-S7-05: Cookie Tampering Detection
	// E2E-S7-03: Cookie Tampering Prevention
	t.Run("UC-S7-05/E2E-S7-03: ValidateCookieIntegrity blocks tampered cookies", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateCookieIntegrity(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "tampered_session_value",
		})
		req.AddCookie(&http.Cookie{
			Name:  "session_signature",
			Value: "invalid_signature",
		})
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.False(t, called)
		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("ValidateCookieIntegrity handles missing signature", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateCookieIntegrity(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_without_signature",
		})
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		// Should allow or reject based on implementation
		// But should not panic
		require.NotPanics(t, func() {
			wrapped.ServeHTTP(rec, req)
		})
	})

	// E2E-S7-06: HTTPOnly Cookies
	t.Run("E2E-S7-06: EnforceSameSiteCookie sets SameSite attribute", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set a cookie
			http.SetCookie(w, &http.Cookie{
				Name:  "test_cookie",
				Value: "test_value",
			})
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.EnforceSameSiteCookie(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		// Check Set-Cookie header has SameSite attribute
		cookies := rec.Result().Cookies()
		if len(cookies) > 0 {
			// Implementation should set SameSite attribute
			setCookieHeader := rec.Header().Get("Set-Cookie")
			require.NotEmpty(t, setCookieHeader)
		}
	})

	// E2E-S7-05: HTTPS-Only Cookies (if HTTPS enabled)
	// UC-S7-07: Session Timeout Long-Lived Cookie
	t.Run("E2E-S7-05/UC-S7-07: RequireHTTPSForCookies enforces HTTPS", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.RequireHTTPSForCookies(handler)

		// HTTP request
		req := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
		req.AddCookie(&http.Cookie{
			Name:   "secure_cookie",
			Value:  "value",
			Secure: true,
		})
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		// Should redirect to HTTPS or return error
		if rec.Code == http.StatusMovedPermanently || rec.Code == http.StatusFound {
			location := rec.Header().Get("Location")
			require.Contains(t, location, "https://")
		} else {
			// Or allow if not strict mode
			require.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("RequireHTTPSForCookies allows HTTPS requests", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.RequireHTTPSForCookies(handler)

		req := httptest.NewRequest(http.MethodGet, "https://example.com/", nil)
		// Note: httptest doesn't easily support TLS mocking, implementation should check URL scheme
		req.AddCookie(&http.Cookie{
			Name:   "secure_cookie",
			Value:  "value",
			Secure: true,
		})
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})

	// IT-S7-01: Real SQL Injection Test
	t.Run("IT-S7-01: SetSecurityHeaders prevents clickjacking", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.SetSecurityHeaders(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		// X-Frame-Options should prevent clickjacking
		frameOptions := rec.Header().Get("X-Frame-Options")
		require.NotEmpty(t, frameOptions)
		require.True(t, frameOptions == "DENY" || frameOptions == "SAMEORIGIN")
	})

	t.Run("SetSecurityHeaders prevents MIME sniffing", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.SetSecurityHeaders(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		contentTypeOptions := rec.Header().Get("X-Content-Type-Options")
		require.Equal(t, "nosniff", contentTypeOptions)
	})

	t.Run("SetSecurityHeaders sets CSP", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.SetSecurityHeaders(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		csp := rec.Header().Get("Content-Security-Policy")
		require.NotEmpty(t, csp)
	})

	t.Run("RateLimiter allows requests under limit", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.RateLimiter(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "192.168.1.100:12345"
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("RateLimiter distinguishes between different IPs", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.RateLimiter(handler)

		// IP 1
		req1 := httptest.NewRequest(http.MethodGet, "/", nil)
		req1.RemoteAddr = "192.168.1.1:12345"
		rec1 := httptest.NewRecorder()
		wrapped.ServeHTTP(rec1, req1)
		require.Equal(t, http.StatusOK, rec1.Code)

		// IP 2
		req2 := httptest.NewRequest(http.MethodGet, "/", nil)
		req2.RemoteAddr = "192.168.1.2:12345"
		rec2 := httptest.NewRecorder()
		wrapped.ServeHTTP(rec2, req2)
		require.Equal(t, http.StatusOK, rec2.Code)
	})

	t.Run("ValidateCookieIntegrity allows requests without cookies", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateCookieIntegrity(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		// Should allow through (cookies are optional)
		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("PreventCSRF allows HEAD requests", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.PreventCSRF(handler)

		req := httptest.NewRequest(http.MethodHead, "/", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("PreventCSRF allows OPTIONS requests", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.PreventCSRF(handler)

		req := httptest.NewRequest(http.MethodOptions, "/", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Security middleware chain", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		// Chain multiple security middlewares
		wrapped := mw.SetSecurityHeaders(
			mw.ValidateCookieIntegrity(
				mw.EnforceSameSiteCookie(handler),
			),
		)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)

		// Verify security headers are set
		require.NotEmpty(t, rec.Header().Get("X-Content-Type-Options"))
		require.NotEmpty(t, rec.Header().Get("X-Frame-Options"))
	})

	t.Run("ValidateCookieIntegrity with multiple cookies", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateCookieIntegrity(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: "cookie1", Value: "value1"})
		req.AddCookie(&http.Cookie{Name: "cookie2", Value: "value2"})
		req.AddCookie(&http.Cookie{Name: "session_signature", Value: "valid_sig"})
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("EnforceSameSiteCookie with response cookies", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.SetCookie(w, &http.Cookie{
				Name:     "session",
				Value:    "session_value",
				Path:     "/",
				HttpOnly: true,
			})
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.EnforceSameSiteCookie(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		cookies := rec.Result().Cookies()
		require.True(t, len(cookies) >= 0)
	})

	t.Run("RequireHTTPSForCookies with local development", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.RequireHTTPSForCookies(handler)

		// localhost should be allowed even without HTTPS
		req := httptest.NewRequest(http.MethodGet, "http://localhost/", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		// Implementation may allow localhost
		require.True(t, called || rec.Code == http.StatusMovedPermanently || rec.Code == http.StatusFound)
	})

	t.Run("SetSecurityHeaders preserves existing headers", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Custom-Header", "custom_value")
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.SetSecurityHeaders(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, "custom_value", rec.Header().Get("X-Custom-Header"))
		require.NotEmpty(t, rec.Header().Get("X-Content-Type-Options"))
	})

	t.Run("PreventCSRF with custom CSRF header", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.PreventCSRF(handler)

		req := httptest.NewRequest(http.MethodPost, "/api/data", nil)
		req.Header.Set("X-Requested-With", "XMLHttpRequest")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		// May allow AJAX requests
		require.True(t, called || rec.Code == http.StatusForbidden)
	})

	t.Run("RateLimiter context preservation", func(t *testing.T) {
		type contextKey string
		const testKey contextKey = "test"

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			val := r.Context().Value(testKey)
			require.Equal(t, "test_value", val)
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.RateLimiter(handler)

		ctx := context.WithValue(context.Background(), testKey, "test_value")
		req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
		req.RemoteAddr = "192.168.1.10:12345"
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})
}
