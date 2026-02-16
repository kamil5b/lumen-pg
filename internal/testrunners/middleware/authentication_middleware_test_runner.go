package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/interfaces/middleware"
)

// AuthenticationMiddlewareConstructor is a function type that creates an AuthenticationMiddleware
type AuthenticationMiddlewareConstructor func() middleware.AuthenticationMiddleware

// AuthenticationMiddlewareRunner runs all authentication middleware tests
// Maps to TEST_PLAN.md:
// - Story 2: Authentication & Identity [UC-S2-08~10, E2E-S2-05]
// - Story 6: Isolation [UC-S6-01, UC-S6-03]
func AuthenticationMiddlewareRunner(t *testing.T, constructor AuthenticationMiddlewareConstructor) {
	t.Helper()

	mw := constructor()

	// UC-S2-08: Session Validation - Valid Session
	t.Run("UC-S2-08: Authenticate with valid session", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			// Verify user context was injected
			user := r.Context().Value("user")
			require.NotNil(t, user)
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.Authenticate(handler)

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "valid_session_123",
		})
		req.AddCookie(&http.Cookie{
			Name:  "username",
			Value: "testuser",
		})
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	// UC-S2-09: Session Validation - Expired Session
	t.Run("UC-S2-09: Authenticate with expired session", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.Authenticate(handler)

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "expired_session_123",
		})
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		// Handler should not be called
		require.False(t, called)
		// Should redirect to login or return 401
		require.True(t, rec.Code == http.StatusUnauthorized || rec.Code == http.StatusFound)
	})

	// UC-S2-10: Session Re-authentication
	t.Run("UC-S2-10: Authenticate requires re-authentication", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.Authenticate(handler)

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		// No session cookie - requires authentication
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.False(t, called)
		require.True(t, rec.Code == http.StatusUnauthorized || rec.Code == http.StatusFound)
	})

	// E2E-S2-05: Protected Route Access Without Auth
	t.Run("E2E-S2-05: RequireAuth blocks unauthenticated requests", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.RequireAuth(handler)

		req := httptest.NewRequest(http.MethodGet, "/main", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.False(t, called)
		require.True(t, rec.Code == http.StatusUnauthorized || rec.Code == http.StatusFound)
	})

	t.Run("RequireAuth allows authenticated requests", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			user := r.Context().Value("user")
			require.NotNil(t, user)
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.RequireAuth(handler)

		req := httptest.NewRequest(http.MethodGet, "/main", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "valid_session_123",
		})
		req.AddCookie(&http.Cookie{
			Name:  "username",
			Value: "testuser",
		})
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	// UC-S6-01: Session Isolation
	t.Run("UC-S6-01: OptionalAuth handles missing session gracefully", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			// User context should be nil for unauthenticated requests
			user := r.Context().Value("user")
			require.Nil(t, user)
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.OptionalAuth(handler)

		req := httptest.NewRequest(http.MethodGet, "/public", nil)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("OptionalAuth injects user context when authenticated", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			user := r.Context().Value("user")
			require.NotNil(t, user)
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.OptionalAuth(handler)

		req := httptest.NewRequest(http.MethodGet, "/public", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "valid_session_123",
		})
		req.AddCookie(&http.Cookie{
			Name:  "username",
			Value: "testuser",
		})
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	// UC-S6-03: Cookie Isolation
	t.Run("UC-S6-03: Authenticate validates session per request", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := r.Context().Value("user")
			require.NotNil(t, user)
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.Authenticate(handler)

		// First request with user1
		req1 := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req1.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_user1",
		})
		req1.AddCookie(&http.Cookie{
			Name:  "username",
			Value: "user1",
		})
		rec1 := httptest.NewRecorder()
		wrapped.ServeHTTP(rec1, req1)
		require.Equal(t, http.StatusOK, rec1.Code)

		// Second request with user2
		req2 := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req2.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_user2",
		})
		req2.AddCookie(&http.Cookie{
			Name:  "username",
			Value: "user2",
		})
		rec2 := httptest.NewRecorder()
		wrapped.ServeHTTP(rec2, req2)
		require.Equal(t, http.StatusOK, rec2.Code)

		// Verify sessions are independent
		require.NotEqual(t, req1.Context(), req2.Context())
	})

	t.Run("Authenticate with invalid session cookie", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.Authenticate(handler)

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "invalid_session",
		})
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.False(t, called)
		require.True(t, rec.Code == http.StatusUnauthorized || rec.Code == http.StatusFound)
	})

	t.Run("Authenticate with missing username cookie", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.Authenticate(handler)

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "valid_session_123",
		})
		// Missing username cookie
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.False(t, called)
		require.True(t, rec.Code == http.StatusUnauthorized || rec.Code == http.StatusFound)
	})

	t.Run("RequireAuth with partial session data", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.RequireAuth(handler)

		req := httptest.NewRequest(http.MethodGet, "/main", nil)
		req.AddCookie(&http.Cookie{
			Name:  "username",
			Value: "testuser",
		})
		// Missing session_id cookie
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.False(t, called)
		require.True(t, rec.Code == http.StatusUnauthorized || rec.Code == http.StatusFound)
	})

	t.Run("OptionalAuth with corrupted session cookie", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			user := r.Context().Value(domain.ContextKeyUser)
			// Should be nil for corrupted session
			require.Nil(t, user)
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.OptionalAuth(handler)

		req := httptest.NewRequest(http.MethodGet, "/public", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "corrupted!!!session",
		})
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Authenticate preserves request context", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			// Verify original context is preserved
			val := r.Context().Value("test_key")
			require.Equal(t, "test_value", val)
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.Authenticate(handler)

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "valid_session_123",
		})
		req.AddCookie(&http.Cookie{
			Name:  "username",
			Value: "testuser",
		})
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("RequireAuth chains multiple authentication checks", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		// Chain authenticate and requireAuth
		wrapped := mw.RequireAuth(mw.Authenticate(handler))

		req := httptest.NewRequest(http.MethodGet, "/main", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "valid_session_123",
		})
		req.AddCookie(&http.Cookie{
			Name:  "username",
			Value: "testuser",
		})
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("OptionalAuth does not block on authentication failure", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.OptionalAuth(handler)

		req := httptest.NewRequest(http.MethodGet, "/public", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "expired_or_invalid",
		})
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		// Should still call handler for OptionalAuth
		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})
}
