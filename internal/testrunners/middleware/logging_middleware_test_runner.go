package middleware

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kamil5b/lumen-pg/internal/interfaces/middleware"
)

// LoggingMiddlewareConstructor is a function type that creates a LoggingMiddleware
type LoggingMiddlewareConstructor func() middleware.LoggingMiddleware

// LoggingMiddlewareRunner runs all logging middleware tests
// Maps to TEST_PLAN.md:
// - Story 4: Manual Query Editor [UC-S4-01~08: Query execution logging]
// - Story 5: Main View & Data Interaction [UC-S5-09~14: Transaction logging]
// - Story 7: Security & Best Practices [UC-S7-01~02: Security event logging]
func LoggingMiddlewareRunner(t *testing.T, constructor LoggingMiddlewareConstructor) {
	t.Helper()

	mw := constructor()

	// UC-S4-01: Single Query Execution (logging)
	t.Run("UC-S4-01: LogRequest logs incoming requests", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.LogRequest(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data?table=users", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
		// Implementation should log: method, path, query params, remote addr
	})

	t.Run("LogRequest logs POST requests", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusCreated)
		})

		wrapped := mw.LogRequest(handler)

		body := bytes.NewBufferString(`{"query": "SELECT * FROM users"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/query", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("LogRequest logs request duration", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Simulate work
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.LogRequest(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		// Implementation should log request duration
	})

	t.Run("LogRequest logs HTTP status codes", func(t *testing.T) {
		statuses := []int{
			http.StatusOK,
			http.StatusCreated,
			http.StatusBadRequest,
			http.StatusUnauthorized,
			http.StatusForbidden,
			http.StatusNotFound,
			http.StatusInternalServerError,
		}

		for _, status := range statuses {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(status)
			})

			wrapped := mw.LogRequest(handler)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()

			wrapped.ServeHTTP(rec, req)

			require.Equal(t, status, rec.Code)
		}
	})

	// UC-S4-01~08: Query execution logging
	t.Run("UC-S4-01~08: LogQueryExecution logs SQL queries", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			// Simulate query execution context
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.LogQueryExecution(handler)

		req := httptest.NewRequest(http.MethodPost, "/api/query", nil)
		ctx := req.Context()
		ctx = context.WithValue(ctx, "query", "SELECT * FROM users WHERE id = 1")
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("LogQueryExecution logs query execution time", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Simulate query execution
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.LogQueryExecution(handler)

		req := httptest.NewRequest(http.MethodPost, "/api/query", nil)
		ctx := context.WithValue(req.Context(), "query", "SELECT COUNT(*) FROM large_table")
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		// Implementation should log query execution duration
	})

	t.Run("LogQueryExecution logs query errors", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Query execution failed", http.StatusInternalServerError)
		})

		wrapped := mw.LogQueryExecution(handler)

		req := httptest.NewRequest(http.MethodPost, "/api/query", nil)
		ctx := context.WithValue(req.Context(), "query", "SELECT * FROM nonexistent_table")
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("LogQueryExecution logs DDL queries", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.LogQueryExecution(handler)

		req := httptest.NewRequest(http.MethodPost, "/api/query", nil)
		ctx := context.WithValue(req.Context(), "query", "CREATE TABLE test (id INT)")
		ctx = context.WithValue(ctx, "query_type", "DDL")
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("LogQueryExecution logs DML queries", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.LogQueryExecution(handler)

		req := httptest.NewRequest(http.MethodPost, "/api/query", nil)
		ctx := context.WithValue(req.Context(), "query", "INSERT INTO users VALUES (1, 'test')")
		ctx = context.WithValue(ctx, "query_type", "DML")
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})

	// UC-S7-01~02: Security event logging
	t.Run("UC-S7-01~02: LogSecurityEvents logs authentication attempts", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.LogSecurityEvents(handler)

		req := httptest.NewRequest(http.MethodPost, "/login", nil)
		ctx := context.WithValue(req.Context(), "username", "testuser")
		ctx = context.WithValue(ctx, "event", "login_attempt")
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("LogSecurityEvents logs failed authentication", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		})

		wrapped := mw.LogSecurityEvents(handler)

		req := httptest.NewRequest(http.MethodPost, "/login", nil)
		ctx := context.WithValue(req.Context(), "username", "testuser")
		ctx = context.WithValue(ctx, "event", "login_failed")
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("LogSecurityEvents logs SQL injection attempts", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		})

		wrapped := mw.LogSecurityEvents(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data?where=1=1; DROP TABLE users", nil)
		ctx := context.WithValue(req.Context(), "event", "sql_injection_attempt")
		ctx = context.WithValue(ctx, "severity", "high")
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("LogSecurityEvents logs unauthorized access attempts", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Forbidden", http.StatusForbidden)
		})

		wrapped := mw.LogSecurityEvents(handler)

		req := httptest.NewRequest(http.MethodGet, "/admin/secret", nil)
		ctx := context.WithValue(req.Context(), "event", "unauthorized_access")
		ctx = context.WithValue(ctx, "resource", "/admin/secret")
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("LogSecurityEvents logs session hijacking attempts", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		})

		wrapped := mw.LogSecurityEvents(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		ctx := context.WithValue(req.Context(), "event", "session_hijack_attempt")
		ctx = context.WithValue(ctx, "session_id", "invalid_session")
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	// UC-S5-09~14: Transaction logging
	t.Run("UC-S5-09~14: LogTransactionEvents logs transaction start", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.LogTransactionEvents(handler)

		req := httptest.NewRequest(http.MethodPost, "/api/transaction/start", nil)
		ctx := context.WithValue(req.Context(), "event", "transaction_start")
		ctx = context.WithValue(ctx, "transaction_id", "txn_123")
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("LogTransactionEvents logs transaction commit", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.LogTransactionEvents(handler)

		req := httptest.NewRequest(http.MethodPost, "/api/transaction/commit", nil)
		ctx := context.WithValue(req.Context(), "event", "transaction_commit")
		ctx = context.WithValue(ctx, "transaction_id", "txn_123")
		ctx = context.WithValue(ctx, "changes_count", 5)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("LogTransactionEvents logs transaction rollback", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.LogTransactionEvents(handler)

		req := httptest.NewRequest(http.MethodPost, "/api/transaction/rollback", nil)
		ctx := context.WithValue(req.Context(), "event", "transaction_rollback")
		ctx = context.WithValue(ctx, "transaction_id", "txn_123")
		ctx = context.WithValue(ctx, "reason", "user_requested")
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("LogTransactionEvents logs transaction timeout", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Transaction timeout", http.StatusRequestTimeout)
		})

		wrapped := mw.LogTransactionEvents(handler)

		req := httptest.NewRequest(http.MethodPost, "/api/transaction/commit", nil)
		ctx := context.WithValue(req.Context(), "event", "transaction_timeout")
		ctx = context.WithValue(ctx, "transaction_id", "txn_456")
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusRequestTimeout, rec.Code)
	})

	t.Run("LogTransactionEvents logs concurrent transaction conflicts", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Conflict", http.StatusConflict)
		})

		wrapped := mw.LogTransactionEvents(handler)

		req := httptest.NewRequest(http.MethodPost, "/api/transaction/commit", nil)
		ctx := context.WithValue(req.Context(), "event", "transaction_conflict")
		ctx = context.WithValue(ctx, "transaction_id", "txn_789")
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusConflict, rec.Code)
	})

	t.Run("Logging middleware chain", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		// Chain all logging middlewares
		wrapped := mw.LogRequest(
			mw.LogQueryExecution(
				mw.LogSecurityEvents(
					mw.LogTransactionEvents(handler),
				),
			),
		)

		req := httptest.NewRequest(http.MethodPost, "/api/query", nil)
		ctx := context.WithValue(req.Context(), "query", "SELECT * FROM users")
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("LogRequest logs user agent", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.LogRequest(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("User-Agent", "Mozilla/5.0")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("LogRequest logs referer", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.LogRequest(handler)

		req := httptest.NewRequest(http.MethodGet, "/page", nil)
		req.Header.Set("Referer", "https://example.com/home")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("LogRequest handles requests without remote addr", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.LogRequest(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = ""
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("LogQueryExecution handles requests without query context", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.LogQueryExecution(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("LogSecurityEvents handles requests without security context", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.LogSecurityEvents(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("LogTransactionEvents handles requests without transaction context", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.LogTransactionEvents(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("LogRequest preserves context", func(t *testing.T) {
		type contextKey string
		const testKey contextKey = "test"

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			val := r.Context().Value(testKey)
			require.Equal(t, "test_value", val)
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.LogRequest(handler)

		ctx := context.WithValue(context.Background(), testKey, "test_value")
		req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("LogQueryExecution logs multiple queries", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.LogQueryExecution(handler)

		queries := []string{
			"SELECT * FROM users",
			"SELECT * FROM posts",
			"SELECT * FROM comments",
		}

		for _, query := range queries {
			req := httptest.NewRequest(http.MethodPost, "/api/query", nil)
			ctx := context.WithValue(req.Context(), "query", query)
			req = req.WithContext(ctx)
			rec := httptest.NewRecorder()

			wrapped.ServeHTTP(rec, req)

			require.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("LogSecurityEvents logs different severity levels", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.LogSecurityEvents(handler)

		severities := []string{"low", "medium", "high", "critical"}

		for _, severity := range severities {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			ctx := context.WithValue(req.Context(), "event", "security_event")
			ctx = context.WithValue(ctx, "severity", severity)
			req = req.WithContext(ctx)
			rec := httptest.NewRecorder()

			wrapped.ServeHTTP(rec, req)

			require.Equal(t, http.StatusOK, rec.Code)
		}
	})
}
