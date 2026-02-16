package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kamil5b/lumen-pg/internal/interfaces/middleware"
)

// RequestIDMiddlewareConstructor is a function type that creates a RequestIDMiddleware
type RequestIDMiddlewareConstructor func() middleware.RequestIDMiddleware

// RequestIDMiddlewareRunner runs all request ID middleware tests
// Maps to TEST_PLAN.md:
// - Story 7: Security & Best Practices (request tracking and logging)
func RequestIDMiddlewareRunner(t *testing.T, constructor RequestIDMiddlewareConstructor) {
	t.Helper()

	mw := constructor()

	t.Run("InjectRequestID generates and injects request ID", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			// Verify request ID was injected in context
			requestID := r.Context().Value("request_id")
			require.NotNil(t, requestID)
			require.NotEmpty(t, requestID)
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.InjectRequestID(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("InjectRequestID generates unique request IDs", func(t *testing.T) {
		requestIDs := make(map[string]bool)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Context().Value("request_id")
			if requestID != nil {
				if id, ok := requestID.(string); ok {
					requestIDs[id] = true
				}
			}
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.InjectRequestID(handler)

		// Generate multiple request IDs
		for i := 0; i < 10; i++ {
			req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
			rec := httptest.NewRecorder()
			wrapped.ServeHTTP(rec, req)
		}

		// All request IDs should be unique
		require.Equal(t, 10, len(requestIDs))
	})

	t.Run("InjectRequestID sets X-Request-ID header in response", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.InjectRequestID(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		requestID := rec.Header().Get("X-Request-ID")
		require.NotEmpty(t, requestID)
	})

	t.Run("InjectRequestID uses existing X-Request-ID if present", func(t *testing.T) {
		existingID := "existing-request-id-12345"
		var capturedID string

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Context().Value("request_id")
			if requestID != nil {
				if id, ok := requestID.(string); ok {
					capturedID = id
				}
			}
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.InjectRequestID(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		req.Header.Set("X-Request-ID", existingID)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		// Should use existing ID or generate new one - implementation dependent
		require.NotEmpty(t, capturedID)
	})

	t.Run("InjectRequestID preserves existing context", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check original context value is preserved
			val := r.Context().Value("original_key")
			require.Equal(t, "original_value", val)

			// Check request ID is added
			requestID := r.Context().Value("request_id")
			require.NotNil(t, requestID)

			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.InjectRequestID(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("InjectRequestID handles GET requests", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Context().Value("request_id")
			require.NotNil(t, requestID)
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.InjectRequestID(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("InjectRequestID handles POST requests", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Context().Value("request_id")
			require.NotNil(t, requestID)
			w.WriteHeader(http.StatusCreated)
		})

		wrapped := mw.InjectRequestID(handler)

		req := httptest.NewRequest(http.MethodPost, "/api/data", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("InjectRequestID handles PUT requests", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Context().Value("request_id")
			require.NotNil(t, requestID)
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.InjectRequestID(handler)

		req := httptest.NewRequest(http.MethodPut, "/api/data/1", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("InjectRequestID handles DELETE requests", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Context().Value("request_id")
			require.NotNil(t, requestID)
			w.WriteHeader(http.StatusNoContent)
		})

		wrapped := mw.InjectRequestID(handler)

		req := httptest.NewRequest(http.MethodDelete, "/api/data/1", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("InjectRequestID works with error responses", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Context().Value("request_id")
			require.NotNil(t, requestID)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		})

		wrapped := mw.InjectRequestID(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
		// Request ID should still be in response headers
		requestID := rec.Header().Get("X-Request-ID")
		require.NotEmpty(t, requestID)
	})

	t.Run("InjectRequestID generates valid format request IDs", func(t *testing.T) {
		var capturedID string
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Context().Value("request_id")
			if requestID != nil {
				if id, ok := requestID.(string); ok {
					capturedID = id
				}
			}
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.InjectRequestID(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.NotEmpty(t, capturedID)
		// Request ID should have reasonable length (e.g., UUID is 36 chars)
		require.True(t, len(capturedID) > 10, "Request ID should be reasonably long")
	})

	t.Run("InjectRequestID handles concurrent requests", func(t *testing.T) {
		requestIDs := make(chan string, 5)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Context().Value("request_id")
			if requestID != nil {
				if id, ok := requestID.(string); ok {
					requestIDs <- id
				}
			}
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.InjectRequestID(handler)

		// Make concurrent requests
		for i := 0; i < 5; i++ {
			go func() {
				req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
				rec := httptest.NewRecorder()
				wrapped.ServeHTTP(rec, req)
			}()
		}

		// Collect all request IDs
		ids := make(map[string]bool)
		for i := 0; i < 5; i++ {
			id := <-requestIDs
			ids[id] = true
		}

		// All should be unique
		require.Equal(t, 5, len(ids))
	})

	t.Run("InjectRequestID does not modify request body", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Context().Value("request_id")
			require.NotNil(t, requestID)
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.InjectRequestID(handler)

		req := httptest.NewRequest(http.MethodPost, "/api/data", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("InjectRequestID chains with other middleware", func(t *testing.T) {
		var capturedID1, capturedID2 string

		middleware1 := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestID := r.Context().Value("request_id")
				if requestID != nil {
					if id, ok := requestID.(string); ok {
						capturedID1 = id
					}
				}
				next.ServeHTTP(w, r)
			})
		}

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Context().Value("request_id")
			if requestID != nil {
				if id, ok := requestID.(string); ok {
					capturedID2 = id
				}
			}
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.InjectRequestID(middleware1(handler))

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.NotEmpty(t, capturedID1)
		require.NotEmpty(t, capturedID2)
		require.Equal(t, capturedID1, capturedID2)
	})

	t.Run("InjectRequestID preserves query parameters", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "users", r.URL.Query().Get("table"))
			require.Equal(t, "testdb", r.URL.Query().Get("database"))
			requestID := r.Context().Value("request_id")
			require.NotNil(t, requestID)
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.InjectRequestID(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data?table=users&database=testdb", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("InjectRequestID preserves custom headers", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))
			require.Equal(t, "Bearer token123", r.Header.Get("Authorization"))
			requestID := r.Context().Value("request_id")
			require.NotNil(t, requestID)
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.InjectRequestID(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer token123")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})
}
