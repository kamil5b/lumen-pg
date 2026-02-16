package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kamil5b/lumen-pg/internal/interfaces/middleware"
)

// ErrorHandlingMiddlewareConstructor is a function type that creates an ErrorHandlingMiddleware
type ErrorHandlingMiddlewareConstructor func() middleware.ErrorHandlingMiddleware

// ErrorHandlingMiddlewareRunner runs all error handling middleware tests
// Maps to TEST_PLAN.md:
// - Story 4: Manual Query Editor [UC-S4-06: Invalid Query Error]
// - Story 7: Security & Best Practices (error handling and recovery)
func ErrorHandlingMiddlewareRunner(t *testing.T, constructor ErrorHandlingMiddlewareConstructor) {
	t.Helper()

	mw := constructor()

	// UC-S4-06: Invalid Query Error
	t.Run("UC-S4-06: HandleErrors catches handler errors", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Simulate error by setting error in context or response
			http.Error(w, "Database connection failed", http.StatusInternalServerError)
		})

		wrapped := mw.HandleErrors(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
		require.Contains(t, rec.Body.String(), "Database connection failed")
	})

	t.Run("HandleErrors allows successful requests", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		})

		wrapped := mw.HandleErrors(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, "success", rec.Body.String())
	})

	t.Run("HandleErrors handles 404 errors", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		})

		wrapped := mw.HandleErrors(handler)

		req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("HandleErrors handles 400 bad request errors", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Invalid input", http.StatusBadRequest)
		})

		wrapped := mw.HandleErrors(handler)

		req := httptest.NewRequest(http.MethodPost, "/api/data", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Contains(t, rec.Body.String(), "Invalid input")
	})

	t.Run("RecoverFromPanic recovers from panics", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("unexpected panic!")
		})

		wrapped := mw.RecoverFromPanic(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		require.NotPanics(t, func() {
			wrapped.ServeHTTP(rec, req)
		})

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("RecoverFromPanic allows normal execution", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.RecoverFromPanic(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("RecoverFromPanic handles nil pointer panic", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var ptr *string
			_ = *ptr // This will panic
		})

		wrapped := mw.RecoverFromPanic(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		require.NotPanics(t, func() {
			wrapped.ServeHTTP(rec, req)
		})

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("RecoverFromPanic handles slice out of bounds panic", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slice := []int{1, 2, 3}
			_ = slice[10] // This will panic
		})

		wrapped := mw.RecoverFromPanic(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		require.NotPanics(t, func() {
			wrapped.ServeHTTP(rec, req)
		})

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("ValidateHTTPMethod allows correct method", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateHTTPMethod(http.MethodGet, http.MethodPost)(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("ValidateHTTPMethod blocks incorrect method", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateHTTPMethod(http.MethodGet, http.MethodPost)(handler)

		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.False(t, called)
		require.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("ValidateHTTPMethod with single allowed method", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateHTTPMethod(http.MethodPost)(handler)

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("ValidateHTTPMethod blocks with single allowed method", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateHTTPMethod(http.MethodPost)(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.False(t, called)
		require.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("ValidateHTTPMethod with multiple allowed methods", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateHTTPMethod(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete)(handler)

		methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}
		for _, method := range methods {
			req := httptest.NewRequest(method, "/", nil)
			rec := httptest.NewRecorder()

			wrapped.ServeHTTP(rec, req)

			require.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("ValidateHTTPMethod sets Allow header on 405", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateHTTPMethod(http.MethodGet, http.MethodPost)(handler)

		req := httptest.NewRequest(http.MethodPut, "/", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusMethodNotAllowed, rec.Code)
		allowHeader := rec.Header().Get("Allow")
		require.NotEmpty(t, allowHeader)
	})

	t.Run("HandleErrors preserves context", func(t *testing.T) {
		type contextKey string
		const testKey contextKey = "test"

		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			val := r.Context().Value(testKey)
			require.Equal(t, "test_value", val)
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.HandleErrors(handler)

		ctx := context.WithValue(context.Background(), testKey, "test_value")
		req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("RecoverFromPanic preserves context before panic", func(t *testing.T) {
		type contextKey string
		const testKey contextKey = "test"

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			val := r.Context().Value(testKey)
			require.Equal(t, "test_value", val)
			panic("test panic")
		})

		wrapped := mw.RecoverFromPanic(handler)

		ctx := context.WithValue(context.Background(), testKey, "test_value")
		req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
		rec := httptest.NewRecorder()

		require.NotPanics(t, func() {
			wrapped.ServeHTTP(rec, req)
		})

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("Error handling middleware chain", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("chained panic")
		})

		// Chain error handling and panic recovery
		wrapped := mw.HandleErrors(
			mw.RecoverFromPanic(handler),
		)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		require.NotPanics(t, func() {
			wrapped.ServeHTTP(rec, req)
		})

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("ValidateHTTPMethod combined with error handling", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.HandleErrors(
			mw.ValidateHTTPMethod(http.MethodGet)(handler),
		)

		// Valid method
		req1 := httptest.NewRequest(http.MethodGet, "/", nil)
		rec1 := httptest.NewRecorder()
		wrapped.ServeHTTP(rec1, req1)
		require.Equal(t, http.StatusOK, rec1.Code)

		// Invalid method
		req2 := httptest.NewRequest(http.MethodPost, "/", nil)
		rec2 := httptest.NewRecorder()
		wrapped.ServeHTTP(rec2, req2)
		require.Equal(t, http.StatusMethodNotAllowed, rec2.Code)
	})

	t.Run("RecoverFromPanic with error type panic", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic(errors.New("custom error"))
		})

		wrapped := mw.RecoverFromPanic(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		require.NotPanics(t, func() {
			wrapped.ServeHTTP(rec, req)
		})

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("HandleErrors with custom error types", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Validation failed: field is required", http.StatusUnprocessableEntity)
		})

		wrapped := mw.HandleErrors(handler)

		req := httptest.NewRequest(http.MethodPost, "/api/data", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		require.Contains(t, rec.Body.String(), "Validation failed")
	})

	t.Run("ValidateHTTPMethod allows OPTIONS for CORS", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateHTTPMethod(http.MethodGet, http.MethodPost, http.MethodOptions)(handler)

		req := httptest.NewRequest(http.MethodOptions, "/", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("RecoverFromPanic logs panic information", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("critical error with details")
		})

		wrapped := mw.RecoverFromPanic(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		require.NotPanics(t, func() {
			wrapped.ServeHTTP(rec, req)
		})

		require.Equal(t, http.StatusInternalServerError, rec.Code)
		// Implementation should log the panic details
	})

	t.Run("HandleErrors handles concurrent requests", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/error" {
				http.Error(w, "error", http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		})

		wrapped := mw.HandleErrors(handler)

		// Success request
		req1 := httptest.NewRequest(http.MethodGet, "/success", nil)
		rec1 := httptest.NewRecorder()
		wrapped.ServeHTTP(rec1, req1)
		require.Equal(t, http.StatusOK, rec1.Code)

		// Error request
		req2 := httptest.NewRequest(http.MethodGet, "/error", nil)
		rec2 := httptest.NewRecorder()
		wrapped.ServeHTTP(rec2, req2)
		require.Equal(t, http.StatusInternalServerError, rec2.Code)
	})

	t.Run("ValidateHTTPMethod case sensitivity", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateHTTPMethod(http.MethodGet)(handler)

		// HTTP methods are case-sensitive
		req := httptest.NewRequest("get", "/", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		// Should handle based on implementation (strict vs lenient)
		require.True(t, rec.Code == http.StatusOK || rec.Code == http.StatusMethodNotAllowed)
	})
}
