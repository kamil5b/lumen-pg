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

// AuthorizationMiddlewareConstructor is a function type that creates an AuthorizationMiddleware
type AuthorizationMiddlewareConstructor func() middleware.AuthorizationMiddleware

// AuthorizationMiddlewareRunner runs all authorization middleware tests
// Maps to TEST_PLAN.md:
// - Story 2: Authentication & Identity [IT-S2-03: Real Role-Based Resource Access]
// - Story 4: Manual Query Editor [IT-S4-04: Query with Permission Denied]
// - Story 5: Main View & Data Interaction [UC-S5-19: Read-Only Mode Enforcement]
// - Story 6: Isolation [UC-S6-02, IT-S6-02: Real Permission Isolation]
func AuthorizationMiddlewareRunner(t *testing.T, constructor AuthorizationMiddlewareConstructor) {
	t.Helper()

	mw := constructor()

	// IT-S2-03: Real Role-Based Resource Access
	t.Run("IT-S2-03: RequireTableAccess grants access to permitted table", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.RequireTableAccess(handler)

		req := httptest.NewRequest(http.MethodGet, "/table/users?database=testdb&schema=public&table=users", nil)
		// Inject role metadata with table access
		ctx := req.Context()
		ctx = context.WithValue(ctx, "permissions", &domain.RoleMetadata{
			Name: "testuser",
			AccessibleTables: []domain.AccessibleTable{
				{Database: "testdb", Schema: "public", Name: "users", HasSelect: true},
			},
		})
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("RequireTableAccess denies access to unpermitted table", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.RequireTableAccess(handler)

		req := httptest.NewRequest(http.MethodGet, "/table/secrets?database=testdb&schema=public&table=secrets", nil)
		ctx := req.Context()
		ctx = context.WithValue(ctx, "permissions", &domain.RoleMetadata{
			Name: "testuser",
			AccessibleTables: []domain.AccessibleTable{
				{Database: "testdb", Schema: "public", Name: "users", HasSelect: true},
			},
		})
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.False(t, called)
		require.Equal(t, http.StatusForbidden, rec.Code)
	})

	// UC-S5-19: Read-Only Mode Enforcement
	t.Run("UC-S5-19: RequireSelectPermission grants access with SELECT permission", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.RequireSelectPermission(handler)

		req := httptest.NewRequest(http.MethodGet, "/table/users?database=testdb&schema=public&table=users", nil)
		ctx := req.Context()
		ctx = context.WithValue(ctx, "permissions", &domain.RoleMetadata{
			Name: "testuser",
			AccessibleTables: []domain.AccessibleTable{
				{Database: "testdb", Schema: "public", Name: "users", HasSelect: true},
			},
		})
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("RequireSelectPermission denies access without SELECT permission", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.RequireSelectPermission(handler)

		req := httptest.NewRequest(http.MethodGet, "/table/users?database=testdb&schema=public&table=users", nil)
		ctx := req.Context()
		ctx = context.WithValue(ctx, "permissions", &domain.RoleMetadata{
			Name: "testuser",
			AccessibleTables: []domain.AccessibleTable{
				{Database: "testdb", Schema: "public", Name: "users", HasSelect: false},
			},
		})
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.False(t, called)
		require.Equal(t, http.StatusForbidden, rec.Code)
	})

	// UC-S5-19: Read-Only Mode Enforcement
	t.Run("UC-S5-19: RequireInsertPermission grants access with INSERT permission", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.RequireInsertPermission(handler)

		req := httptest.NewRequest(http.MethodPost, "/table/users?database=testdb&schema=public&table=users", nil)
		ctx := req.Context()
		ctx = context.WithValue(ctx, "permissions", &domain.RoleMetadata{
			Name: "testuser",
			AccessibleTables: []domain.AccessibleTable{
				{Database: "testdb", Schema: "public", Name: "users", HasInsert: true},
			},
		})
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("RequireInsertPermission denies access without INSERT permission", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.RequireInsertPermission(handler)

		req := httptest.NewRequest(http.MethodPost, "/table/users?database=testdb&schema=public&table=users", nil)
		ctx := req.Context()
		ctx = context.WithValue(ctx, "permissions", &domain.RoleMetadata{
			Name: "testuser",
			AccessibleTables: []domain.AccessibleTable{
				{Database: "testdb", Schema: "public", Name: "users", HasInsert: false},
			},
		})
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.False(t, called)
		require.Equal(t, http.StatusForbidden, rec.Code)
	})

	// UC-S5-19: Read-Only Mode Enforcement
	t.Run("UC-S5-19: RequireUpdatePermission grants access with UPDATE permission", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.RequireUpdatePermission(handler)

		req := httptest.NewRequest(http.MethodPut, "/table/users?database=testdb&schema=public&table=users", nil)
		ctx := req.Context()
		ctx = context.WithValue(ctx, "permissions", &domain.RoleMetadata{
			Name: "testuser",
			AccessibleTables: []domain.AccessibleTable{
				{Database: "testdb", Schema: "public", Name: "users", HasUpdate: true},
			},
		})
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("RequireUpdatePermission denies access without UPDATE permission", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.RequireUpdatePermission(handler)

		req := httptest.NewRequest(http.MethodPut, "/table/users?database=testdb&schema=public&table=users", nil)
		ctx := req.Context()
		ctx = context.WithValue(ctx, "permissions", &domain.RoleMetadata{
			Name: "testuser",
			AccessibleTables: []domain.AccessibleTable{
				{Database: "testdb", Schema: "public", Name: "users", HasUpdate: false},
			},
		})
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.False(t, called)
		require.Equal(t, http.StatusForbidden, rec.Code)
	})

	// UC-S5-19: Read-Only Mode Enforcement
	t.Run("UC-S5-19: RequireDeletePermission grants access with DELETE permission", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.RequireDeletePermission(handler)

		req := httptest.NewRequest(http.MethodDelete, "/table/users?database=testdb&schema=public&table=users", nil)
		ctx := req.Context()
		ctx = context.WithValue(ctx, "permissions", &domain.RoleMetadata{
			Name: "testuser",
			AccessibleTables: []domain.AccessibleTable{
				{Database: "testdb", Schema: "public", Name: "users", HasDelete: true},
			},
		})
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("RequireDeletePermission denies access without DELETE permission", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.RequireDeletePermission(handler)

		req := httptest.NewRequest(http.MethodDelete, "/table/users?database=testdb&schema=public&table=users", nil)
		ctx := req.Context()
		ctx = context.WithValue(ctx, "permissions", &domain.RoleMetadata{
			Name: "testuser",
			AccessibleTables: []domain.AccessibleTable{
				{Database: "testdb", Schema: "public", Name: "users", HasDelete: false},
			},
		})
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.False(t, called)
		require.Equal(t, http.StatusForbidden, rec.Code)
	})

	// IT-S2-03: Real Role-Based Resource Access
	t.Run("IT-S2-03: RequireDatabaseAccess grants access to permitted database", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.RequireDatabaseAccess(handler)

		req := httptest.NewRequest(http.MethodGet, "/database?database=testdb", nil)
		ctx := req.Context()
		ctx = context.WithValue(ctx, "permissions", &domain.RoleMetadata{
			Name:                "testuser",
			AccessibleDatabases: []string{"testdb", "testdb2"},
		})
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("RequireDatabaseAccess denies access to unpermitted database", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.RequireDatabaseAccess(handler)

		req := httptest.NewRequest(http.MethodGet, "/database?database=secretdb", nil)
		ctx := req.Context()
		ctx = context.WithValue(ctx, "permissions", &domain.RoleMetadata{
			Name:                "testuser",
			AccessibleDatabases: []string{"testdb"},
		})
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.False(t, called)
		require.Equal(t, http.StatusForbidden, rec.Code)
	})

	// UC-S6-02: Transaction Isolation
	// IT-S6-02: Real Permission Isolation
	t.Run("UC-S6-02/IT-S6-02: Different users have different permissions", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.RequireUpdatePermission(handler)

		// User1 with UPDATE permission
		req1 := httptest.NewRequest(http.MethodPut, "/table/users?database=testdb&schema=public&table=users", nil)
		ctx1 := req1.Context()
		ctx1 = context.WithValue(ctx1, "permissions", &domain.RoleMetadata{
			Name: "user1",
			AccessibleTables: []domain.AccessibleTable{
				{Database: "testdb", Schema: "public", Name: "users", HasUpdate: true},
			},
		})
		req1 = req1.WithContext(ctx1)
		rec1 := httptest.NewRecorder()
		wrapped.ServeHTTP(rec1, req1)
		require.Equal(t, http.StatusOK, rec1.Code)

		// User2 without UPDATE permission
		req2 := httptest.NewRequest(http.MethodPut, "/table/users?database=testdb&schema=public&table=users", nil)
		ctx2 := req2.Context()
		ctx2 = context.WithValue(ctx2, "permissions", &domain.RoleMetadata{
			Name: "user2",
			AccessibleTables: []domain.AccessibleTable{
				{Database: "testdb", Schema: "public", Name: "users", HasUpdate: false},
			},
		})
		req2 = req2.WithContext(ctx2)
		rec2 := httptest.NewRecorder()
		wrapped.ServeHTTP(rec2, req2)
		require.Equal(t, http.StatusForbidden, rec2.Code)
	})

	t.Run("Authorization fails with missing user context", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.RequireTableAccess(handler)

		req := httptest.NewRequest(http.MethodGet, "/table/users?database=testdb&schema=public&table=users", nil)
		// No user context injected
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.False(t, called)
		require.True(t, rec.Code == http.StatusUnauthorized || rec.Code == http.StatusForbidden)
	})

	t.Run("Authorization checks multiple tables independently", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.RequireSelectPermission(handler)

		ctx := context.WithValue(context.Background(), "permissions", &domain.RoleMetadata{
			Name: "testuser",
			AccessibleTables: []domain.AccessibleTable{
				{Database: "testdb", Schema: "public", Name: "users", HasSelect: true},
				{Database: "testdb", Schema: "public", Name: "posts", HasSelect: false},
			},
		})

		// Access to users table - allowed
		req1 := httptest.NewRequest(http.MethodGet, "/table/users?database=testdb&schema=public&table=users", nil)
		req1 = req1.WithContext(ctx)
		rec1 := httptest.NewRecorder()
		wrapped.ServeHTTP(rec1, req1)
		require.Equal(t, http.StatusOK, rec1.Code)

		// Access to posts table - denied
		req2 := httptest.NewRequest(http.MethodGet, "/table/posts?database=testdb&schema=public&table=posts", nil)
		req2 = req2.WithContext(ctx)
		rec2 := httptest.NewRecorder()
		wrapped.ServeHTTP(rec2, req2)
		require.Equal(t, http.StatusForbidden, rec2.Code)
	})

	t.Run("Authorization chain multiple permission checks", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		// Chain database and table access checks
		wrapped := mw.RequireTableAccess(mw.RequireDatabaseAccess(handler))

		req := httptest.NewRequest(http.MethodGet, "/table/users?database=testdb&schema=public&table=users", nil)
		ctx := req.Context()
		ctx = context.WithValue(ctx, "permissions", &domain.RoleMetadata{
			Name:                "testuser",
			AccessibleDatabases: []string{"testdb"},
			AccessibleTables: []domain.AccessibleTable{
				{Database: "testdb", Schema: "public", Name: "users", HasSelect: true},
			},
		})
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.True(t, called)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Authorization with complex permission scenarios", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.RequireSelectPermission(handler)

		// User with multiple tables, different permissions
		ctx := context.WithValue(context.Background(), "permissions", &domain.RoleMetadata{
			Name: "poweruser",
			AccessibleTables: []domain.AccessibleTable{
				{Database: "db1", Schema: "public", Name: "table1", HasSelect: true, HasInsert: true, HasUpdate: true, HasDelete: true},
				{Database: "db1", Schema: "public", Name: "table2", HasSelect: true, HasInsert: false, HasUpdate: false, HasDelete: false},
				{Database: "db2", Schema: "public", Name: "table3", HasSelect: false, HasInsert: false, HasUpdate: false, HasDelete: false},
			},
		})

		// table1 - full permissions
		req1 := httptest.NewRequest(http.MethodGet, "/table?database=db1&schema=public&table=table1", nil)
		req1 = req1.WithContext(ctx)
		rec1 := httptest.NewRecorder()
		wrapped.ServeHTTP(rec1, req1)
		require.Equal(t, http.StatusOK, rec1.Code)

		// table2 - select only
		req2 := httptest.NewRequest(http.MethodGet, "/table?database=db1&schema=public&table=table2", nil)
		req2 = req2.WithContext(ctx)
		rec2 := httptest.NewRecorder()
		wrapped.ServeHTTP(rec2, req2)
		require.Equal(t, http.StatusOK, rec2.Code)

		// table3 - no select permission
		req3 := httptest.NewRequest(http.MethodGet, "/table?database=db2&schema=public&table=table3", nil)
		req3 = req3.WithContext(ctx)
		rec3 := httptest.NewRecorder()
		wrapped.ServeHTTP(rec3, req3)
		require.Equal(t, http.StatusForbidden, rec3.Code)
	})

	t.Run("Authorization handles empty permission lists", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		wrapped := mw.RequireTableAccess(handler)

		req := httptest.NewRequest(http.MethodGet, "/table/users?database=testdb&schema=public&table=users", nil)
		ctx := req.Context()
		ctx = context.WithValue(ctx, "permissions", &domain.RoleMetadata{
			Name:             "limiteduser",
			AccessibleTables: []domain.AccessibleTable{}, // Empty list
		})
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		require.False(t, called)
		require.Equal(t, http.StatusForbidden, rec.Code)
	})
}
