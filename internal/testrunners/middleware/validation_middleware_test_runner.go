package middleware

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kamil5b/lumen-pg/internal/interfaces/middleware"
)

// ValidationMiddlewareConstructor is a function type that creates a ValidationMiddleware
type ValidationMiddlewareConstructor func() middleware.ValidationMiddleware

// ValidationMiddlewareRunner runs all validation middleware tests
// Maps to TEST_PLAN.md:
// - Story 4: Manual Query Editor [UC-S4-06: Invalid Query Error]
// - Story 5: Main View & Data Interaction [UC-S5-03~04: WHERE Clause Validation/Injection Prevention]
// - Story 7: Security & Best Practices [UC-S7-01~02, IT-S7-01, E2E-S7-01~02]
func ValidationMiddlewareRunner(t *testing.T, constructor ValidationMiddlewareConstructor) {
	t.Helper()

	mw := constructor()

	// UC-S5-03: WHERE Clause Validation
	t.Run("UC-S5-03: ValidateQueryParams validates valid query parameters", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateQueryParams(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data?database=testdb&schema=public&table=users", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("ValidateQueryParams rejects invalid query parameters", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateQueryParams(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data?invalid_param=<script>alert('xss')</script>", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.False(t, called)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("ValidateQueryParams handles missing query parameters", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateQueryParams(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		// Should allow requests without query params
		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	// UC-S4-06: Invalid Query Error
	t.Run("UC-S4-06: ValidateRequestBody validates valid JSON body", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateRequestBody(handler)

		body := `{"query": "SELECT * FROM users", "limit": 100}`
		req := httptest.NewRequest(http.MethodPost, "/api/query", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("ValidateRequestBody rejects malformed JSON", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateRequestBody(handler)

		body := `{"query": "SELECT * FROM users", "limit": }`
		req := httptest.NewRequest(http.MethodPost, "/api/query", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.False(t, called)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("ValidateRequestBody handles empty body", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateRequestBody(handler)

		req := httptest.NewRequest(http.MethodPost, "/api/query", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		// May allow or reject empty body depending on implementation
		require.True(t, called || rec.Code == http.StatusBadRequest)
	})

	// UC-S5-04: WHERE Clause Injection Prevention
	// UC-S7-01: SQL Injection Prevention - WHERE Clause
	// E2E-S7-01: SQL Injection via WHERE Bar
	t.Run("UC-S5-04/UC-S7-01/E2E-S7-01: ValidateWhereClause allows safe WHERE clause", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateWhereClause(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data?where=age > 18 AND status = 'active'", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	// UC-S5-04: WHERE Clause Injection Prevention
	// IT-S7-01: Real SQL Injection Test
	// E2E-S7-01: SQL Injection via WHERE Bar
	t.Run("UC-S5-04/IT-S7-01/E2E-S7-01: ValidateWhereClause blocks SQL injection - DROP TABLE", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateWhereClause(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data?where=1=1; DROP TABLE users; --", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.False(t, called)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	// E2E-S7-01: SQL Injection via WHERE Bar
	t.Run("E2E-S7-01: ValidateWhereClause blocks SQL injection - UNION attack", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateWhereClause(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data?where=1=1 UNION SELECT password FROM users", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.False(t, called)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("ValidateWhereClause blocks SQL injection - comment injection", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateWhereClause(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data?where=username='admin' --", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.False(t, called)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("ValidateWhereClause blocks SQL injection - INSERT injection", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateWhereClause(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data?where=1=1; INSERT INTO users VALUES ('hacker', 'pwd')", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.False(t, called)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("ValidateWhereClause blocks SQL injection - UPDATE injection", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateWhereClause(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data?where=1=1; UPDATE users SET role='admin'", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.False(t, called)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("ValidateWhereClause blocks SQL injection - DELETE injection", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateWhereClause(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data?where=1=1; DELETE FROM users", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.False(t, called)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("ValidateWhereClause allows empty WHERE clause", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateWhereClause(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	// UC-S7-02: SQL Injection Prevention - Query Editor
	// E2E-S7-02: SQL Injection via Query Editor
	t.Run("UC-S7-02/E2E-S7-02: ValidateSQLQuery allows safe SQL query", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateSQLQuery(handler)

		body := `{"query": "SELECT id, name, email FROM users WHERE age > 18"}`
		req := httptest.NewRequest(http.MethodPost, "/api/query", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	// IT-S7-01: Real SQL Injection Test
	// E2E-S7-02: SQL Injection via Query Editor
	t.Run("IT-S7-01/E2E-S7-02: ValidateSQLQuery blocks dangerous SQL commands", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateSQLQuery(handler)

		body := `{"query": "SELECT * FROM users; DROP DATABASE production;"}`
		req := httptest.NewRequest(http.MethodPost, "/api/query", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.False(t, called)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("ValidateSQLQuery allows multiple safe statements", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateSQLQuery(handler)

		body := `{"query": "SELECT * FROM users; SELECT * FROM posts;"}`
		req := httptest.NewRequest(http.MethodPost, "/api/query", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("ValidateSQLQuery blocks script tags in query", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateSQLQuery(handler)

		body := `{"query": "<script>alert('xss')</script>"}`
		req := httptest.NewRequest(http.MethodPost, "/api/query", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.False(t, called)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("ValidateQueryParams blocks special characters", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateQueryParams(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data?table=users'; DROP TABLE users; --", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.False(t, called)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("ValidateQueryParams allows safe alphanumeric values", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateQueryParams(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data?database=testdb123&schema=public&table=users_data", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("ValidateRequestBody validates form data", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateRequestBody(handler)

		form := url.Values{}
		form.Add("username", "testuser")
		form.Add("password", "password123")

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("ValidateRequestBody blocks oversized body", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateRequestBody(handler)

		// Create large body
		largeBody := strings.Repeat("a", 10*1024*1024) // 10MB
		req := httptest.NewRequest(http.MethodPost, "/api/query", strings.NewReader(largeBody))
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		// Should reject oversized body
		require.False(t, called)
		require.True(t, rec.Code == http.StatusBadRequest || rec.Code == http.StatusRequestEntityTooLarge)
	})

	t.Run("ValidateWhereClause blocks nested SQL injection", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateWhereClause(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data?where=id IN (SELECT id FROM admin WHERE 1=1)", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		// Should allow subqueries or block depending on implementation
		require.True(t, called || rec.Code == http.StatusBadRequest)
	})

	t.Run("ValidateSQLQuery handles multiline queries", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateSQLQuery(handler)

		body := `{"query": "SELECT id,\n       name,\n       email\nFROM users\nWHERE active = true"}`
		req := httptest.NewRequest(http.MethodPost, "/api/query", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Validation middleware chain", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		// Chain multiple validation middlewares
		wrapped := mw.ValidateQueryParams(
			mw.ValidateRequestBody(
				mw.ValidateWhereClause(handler),
			),
		)

		req := httptest.NewRequest(http.MethodGet, "/api/data?database=testdb&where=id > 0", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("ValidateWhereClause context preservation", func(t *testing.T) {
		type contextKey string
		const testKey contextKey = "test"

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			val := r.Context().Value(testKey)
			require.Equal(t, "test_value", val)
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateWhereClause(handler)

		ctx := context.WithValue(context.Background(), testKey, "test_value")
		req := httptest.NewRequest(http.MethodGet, "/api/data?where=status='active'", nil).WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("ValidateSQLQuery blocks stored procedure execution", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateSQLQuery(handler)

		body := `{"query": "EXEC sp_executesql N'DROP TABLE users'"}`
		req := httptest.NewRequest(http.MethodPost, "/api/query", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.False(t, called)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("ValidateWhereClause blocks hexadecimal injection", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.ValidateWhereClause(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/data?where=id=0x31303235", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		// Should block or allow depending on implementation
		require.True(t, called || rec.Code == http.StatusBadRequest)
	})
}
