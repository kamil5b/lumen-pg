package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/interfaces/middleware"
)

// ContextMiddlewareConstructor is a function type that creates a ContextMiddleware
type ContextMiddlewareConstructor func() middleware.ContextMiddleware

// ContextMiddlewareRunner runs all context middleware tests
// Maps to TEST_PLAN.md:
// - Story 2: Authentication & Identity [UC-S2-08~10, UC-S2-11: Data Explorer Population]
// - Story 6: Isolation [UC-S6-01~03, IT-S6-01~03]
func ContextMiddlewareRunner(t *testing.T, constructor ContextMiddlewareConstructor) {
	t.Helper()

	mw := constructor()

	// UC-S2-08: Session Validation - Valid Session
	// UC-S6-01: Session Isolation
	t.Run("UC-S2-08/UC-S6-01: InjectUser injects user into context", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			user := r.Context().Value(domain.ContextKeyUser)
			require.NotNil(t, user)

			userObj, ok := user.(*domain.User)
			require.True(t, ok)
			require.Equal(t, "testuser", userObj.Username)
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.InjectUser(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  "username",
			Value: "testuser",
		})
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("InjectUser handles missing user cookie", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			user := r.Context().Value(domain.ContextKeyUser)
			require.Nil(t, user)
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.InjectUser(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	// UC-S2-09: Session Validation - Expired Session
	// UC-S6-01: Session Isolation
	t.Run("UC-S2-09/UC-S6-01: InjectSession injects session into context", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			session := r.Context().Value(domain.ContextKeySession)
			require.NotNil(t, session)

			sessionObj, ok := session.(*domain.Session)
			require.True(t, ok)
			require.Equal(t, "session_123", sessionObj.ID)
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.InjectSession(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("InjectSession handles missing session cookie", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			session := r.Context().Value(domain.ContextKeySession)
			require.Nil(t, session)
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.InjectSession(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	// UC-S6-02: Transaction Isolation
	// IT-S6-03: Real Transaction Isolation
	t.Run("UC-S6-02/IT-S6-03: InjectTransaction injects transaction into context", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			txn := r.Context().Value(domain.ContextKeyTransaction)
			require.NotNil(t, txn)

			txnObj, ok := txn.(*domain.TransactionState)
			require.True(t, ok)
			require.Equal(t, "txn_456", txnObj.ID)
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.InjectTransaction(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  "transaction_id",
			Value: "txn_456",
		})
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("InjectTransaction handles missing transaction cookie", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			txn := r.Context().Value(domain.ContextKeyTransaction)
			require.Nil(t, txn)
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.InjectTransaction(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	// UC-S2-11: Data Explorer Population After Login
	// IT-S2-03: Real Role-Based Resource Access
	t.Run("UC-S2-11/IT-S2-03: InjectUserPermissions injects permissions into context", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			perms := r.Context().Value(domain.ContextKeyPermissions)
			require.NotNil(t, perms)

			permsObj, ok := perms.(*domain.RoleMetadata)
			require.True(t, ok)
			require.Len(t, permsObj.AccessibleTables, 2)
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.InjectUserPermissions(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  "username",
			Value: "testuser",
		})
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("InjectUserPermissions handles unauthenticated user", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			perms := r.Context().Value(domain.ContextKeyPermissions)
			require.Nil(t, perms)
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.InjectUserPermissions(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	// IT-S2-04: Session Persistence After Probe
	t.Run("IT-S2-04: InjectMetadata injects cached metadata into context", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			metadata := r.Context().Value(domain.ContextKeyMetadata)
			require.NotNil(t, metadata)

			metaObj, ok := metadata.(*domain.DatabaseMetadata)
			require.True(t, ok)
			require.NotEmpty(t, metaObj.Schemas)
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.InjectMetadata(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_789",
		})
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("InjectMetadata handles missing session", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			metadata := r.Context().Value(domain.ContextKeyMetadata)
			require.Nil(t, metadata)
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.InjectMetadata(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	// UC-S6-01: Session Isolation
	// IT-S6-01: Real Multi-User Connection
	t.Run("UC-S6-01/IT-S6-01: Multiple users have isolated contexts", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := r.Context().Value(domain.ContextKeyUser)
			require.NotNil(t, user)
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.InjectUser(handler)

		// User 1
		req1 := httptest.NewRequest(http.MethodGet, "/", nil)
		req1.AddCookie(&http.Cookie{Name: "username", Value: "user1"})
		rec1 := httptest.NewRecorder()
		wrapped.ServeHTTP(rec1, req1)
		require.Equal(t, http.StatusOK, rec1.Code)

		// User 2
		req2 := httptest.NewRequest(http.MethodGet, "/", nil)
		req2.AddCookie(&http.Cookie{Name: "username", Value: "user2"})
		rec2 := httptest.NewRecorder()
		wrapped.ServeHTTP(rec2, req2)
		require.Equal(t, http.StatusOK, rec2.Code)

		// Verify contexts are different
		require.NotEqual(t, req1.Context(), req2.Context())
	})

	t.Run("Context middleware chain preserves all injected values", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true

			user := r.Context().Value(domain.ContextKeyUser)
			session := r.Context().Value(domain.ContextKeySession)
			perms := r.Context().Value(domain.ContextKeyPermissions)

			require.NotNil(t, user)
			require.NotNil(t, session)
			require.NotNil(t, perms)
			w.WriteHeader(http.StatusOK)
		})

		// Chain all context injectors
		wrapped := mw.InjectUser(
			mw.InjectSession(
				mw.InjectUserPermissions(handler),
			),
		)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: "username", Value: "testuser"})
		req.AddCookie(&http.Cookie{Name: "session_id", Value: "session_123"})
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("InjectUser preserves existing context values", func(t *testing.T) {
		type contextKey string
		const testKey contextKey = "test_key"

		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true

			// Original context value should be preserved
			val := r.Context().Value(testKey)
			require.Equal(t, "test_value", val)

			// New user should be injected
			user := r.Context().Value(domain.ContextKeyUser)
			require.NotNil(t, user)

			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.InjectUser(handler)

		ctx := context.WithValue(context.Background(), testKey, "test_value")
		req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
		req.AddCookie(&http.Cookie{Name: "username", Value: "testuser"})
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("InjectTransaction with active transaction flag", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true

			txn := r.Context().Value(domain.ContextKeyTransaction)
			require.NotNil(t, txn)

			txnObj, ok := txn.(*domain.TransactionState)
			require.True(t, ok)
			require.NotEmpty(t, txnObj.ID)

			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.InjectTransaction(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  "transaction_id",
			Value: "txn_active_123",
		})
		req.AddCookie(&http.Cookie{
			Name:  "transaction_active",
			Value: "true",
		})
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("InjectUserPermissions with multiple databases", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true

			perms := r.Context().Value(domain.ContextKeyPermissions)
			require.NotNil(t, perms)

			permsObj, ok := perms.(*domain.RoleMetadata)
			require.True(t, ok)
			require.Len(t, permsObj.AccessibleDatabases, 3)

			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.InjectUserPermissions(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: "username", Value: "poweruser"})
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("InjectMetadata with cached database schema", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true

			metadata := r.Context().Value(domain.ContextKeyMetadata)
			require.NotNil(t, metadata)

			metaObj, ok := metadata.(*domain.DatabaseMetadata)
			require.True(t, ok)
			require.NotEmpty(t, metaObj.Schemas)

			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.InjectMetadata(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: "session_id", Value: "session_with_cache"})
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	// IT-S6-02: Real Permission Isolation
	t.Run("IT-S6-02: InjectUserPermissions isolates permissions per user", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			perms := r.Context().Value(domain.ContextKeyPermissions)
			if perms != nil {
				permsObj, _ := perms.(*domain.RoleMetadata)
				// Write username to distinguish responses
				w.Write([]byte(permsObj.Name))
			}
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.InjectUserPermissions(handler)

		// User with limited permissions
		req1 := httptest.NewRequest(http.MethodGet, "/", nil)
		req1.AddCookie(&http.Cookie{Name: "username", Value: "limiteduser"})
		rec1 := httptest.NewRecorder()
		wrapped.ServeHTTP(rec1, req1)

		// User with admin permissions
		req2 := httptest.NewRequest(http.MethodGet, "/", nil)
		req2.AddCookie(&http.Cookie{Name: "username", Value: "adminuser"})
		rec2 := httptest.NewRecorder()
		wrapped.ServeHTTP(rec2, req2)

		// Responses should be different
		require.NotEqual(t, rec1.Body.String(), rec2.Body.String())
	})

	t.Run("All context injectors handle nil gracefully", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.InjectUser(
			mw.InjectSession(
				mw.InjectTransaction(
					mw.InjectUserPermissions(
						mw.InjectMetadata(handler),
					),
				),
			),
		)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		// Should not panic without cookies
		require.NotPanics(t, func() {
			wrapped.ServeHTTP(rec, req)
		})

		require.Equal(t, http.StatusOK, rec.Code)
	})
}
