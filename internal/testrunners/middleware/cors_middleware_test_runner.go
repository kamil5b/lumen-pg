package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kamil5b/lumen-pg/internal/interfaces/middleware"
)

// CORSMiddlewareConstructor is a function type that creates a CORSMiddleware
type CORSMiddlewareConstructor func() middleware.CORSMiddleware

// CORSMiddlewareRunner runs all CORS middleware tests
// Maps to TEST_PLAN.md:
// - Story 7: Security & Best Practices (CORS handling)
func CORSMiddlewareRunner(t *testing.T, constructor CORSMiddlewareConstructor) {
	t.Helper()

	mw := constructor()

	t.Run("HandleCORS allows same-origin requests", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.HandleCORS(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		req.Header.Set("Origin", "http://localhost")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("HandleCORS handles preflight OPTIONS request", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.HandleCORS(handler)

		req := httptest.NewRequest(http.MethodOptions, "/api/data", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		req.Header.Set("Access-Control-Request-Method", "POST")
		req.Header.Set("Access-Control-Request-Headers", "Content-Type")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		// Preflight should be handled by middleware
		require.Equal(t, http.StatusOK, rec.Code)

		// Check CORS headers are set
		headers := rec.Header()
		require.NotEmpty(t, headers.Get("Access-Control-Allow-Origin"))
		require.NotEmpty(t, headers.Get("Access-Control-Allow-Methods"))
	})

	t.Run("HandleCORS sets Access-Control-Allow-Origin header", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.HandleCORS(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		allowOrigin := rec.Header().Get("Access-Control-Allow-Origin")
		require.NotEmpty(t, allowOrigin)
	})

	t.Run("HandleCORS sets Access-Control-Allow-Methods header", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.HandleCORS(handler)

		req := httptest.NewRequest(http.MethodOptions, "/api/data", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		req.Header.Set("Access-Control-Request-Method", "POST")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		allowMethods := rec.Header().Get("Access-Control-Allow-Methods")
		require.NotEmpty(t, allowMethods)
		require.Contains(t, allowMethods, "GET")
		require.Contains(t, allowMethods, "POST")
	})

	t.Run("HandleCORS sets Access-Control-Allow-Headers header", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.HandleCORS(handler)

		req := httptest.NewRequest(http.MethodOptions, "/api/data", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		req.Header.Set("Access-Control-Request-Headers", "Content-Type, Authorization")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		allowHeaders := rec.Header().Get("Access-Control-Allow-Headers")
		require.NotEmpty(t, allowHeaders)
	})

	t.Run("HandleCORS sets Access-Control-Allow-Credentials header", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.HandleCORS(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		allowCredentials := rec.Header().Get("Access-Control-Allow-Credentials")
		require.True(t, allowCredentials == "true" || allowCredentials == "")
	})

	t.Run("HandleCORS handles GET request with Origin", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.HandleCORS(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		req.Header.Set("Origin", "http://example.com")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("HandleCORS handles POST request with Origin", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusCreated)
		})

		wrapped := mw.HandleCORS(handler)

		req := httptest.NewRequest(http.MethodPost, "/api/data", nil)
		req.Header.Set("Origin", "http://example.com")
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("HandleCORS handles PUT request with Origin", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.HandleCORS(handler)

		req := httptest.NewRequest(http.MethodPut, "/api/data/1", nil)
		req.Header.Set("Origin", "http://example.com")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("HandleCORS handles DELETE request with Origin", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusNoContent)
		})

		wrapped := mw.HandleCORS(handler)

		req := httptest.NewRequest(http.MethodDelete, "/api/data/1", nil)
		req.Header.Set("Origin", "http://example.com")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("HandleCORS handles request without Origin header", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.HandleCORS(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("HandleCORS sets Access-Control-Max-Age for preflight", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.HandleCORS(handler)

		req := httptest.NewRequest(http.MethodOptions, "/api/data", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		req.Header.Set("Access-Control-Request-Method", "POST")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		maxAge := rec.Header().Get("Access-Control-Max-Age")
		require.NotEmpty(t, maxAge)
	})

	t.Run("HandleCORS handles multiple origins", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.HandleCORS(handler)

		origins := []string{
			"http://localhost:3000",
			"http://localhost:8080",
			"https://example.com",
		}

		for _, origin := range origins {
			req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
			req.Header.Set("Origin", origin)
			rec := httptest.NewRecorder()

			wrapped.ServeHTTP(rec, req)

			require.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("HandleCORS handles custom headers in preflight", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.HandleCORS(handler)

		req := httptest.NewRequest(http.MethodOptions, "/api/data", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		req.Header.Set("Access-Control-Request-Method", "POST")
		req.Header.Set("Access-Control-Request-Headers", "X-Custom-Header, X-Another-Header")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		allowHeaders := rec.Header().Get("Access-Control-Allow-Headers")
		require.NotEmpty(t, allowHeaders)
	})

	t.Run("HandleCORS preserves existing response headers", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Custom-Header", "custom-value")
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.HandleCORS(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, "custom-value", rec.Header().Get("X-Custom-Header"))
		require.NotEmpty(t, rec.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("HandleCORS handles Vary header", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.HandleCORS(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		vary := rec.Header().Get("Vary")
		require.True(t, vary == "Origin" || vary != "")
	})

	t.Run("HandleCORS does not expose credentials with wildcard origin", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.HandleCORS(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		req.Header.Set("Origin", "http://example.com")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		allowOrigin := rec.Header().Get("Access-Control-Allow-Origin")
		allowCredentials := rec.Header().Get("Access-Control-Allow-Credentials")

		if allowOrigin == "*" {
			require.NotEqual(t, "true", allowCredentials)
		}
	})

	t.Run("HandleCORS handles complex preflight with multiple methods", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.HandleCORS(handler)

		req := httptest.NewRequest(http.MethodOptions, "/api/data", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		req.Header.Set("Access-Control-Request-Method", "PUT")
		req.Header.Set("Access-Control-Request-Headers", "Content-Type, Authorization, X-CSRF-Token")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		allowMethods := rec.Header().Get("Access-Control-Allow-Methods")
		require.NotEmpty(t, allowMethods)
	})

	t.Run("HandleCORS allows requests from localhost", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.HandleCORS(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		req.Header.Set("Origin", "http://localhost")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})
}
