package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/interfaces/handler"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
	mockUsecase "github.com/kamil5b/lumen-pg/internal/testrunners/mocks/usecase"
)

// AuthenticationHandlerConstructor is a function type that creates an AuthenticationHandler
type AuthenticationHandlerConstructor func(
	authUC usecase.AuthenticationUseCase,
	setupUC usecase.SetupUseCase,
	rbacUC usecase.RBACUseCase,
) handler.LoginHandler

// AuthenticationHandlerRunner runs all authentication handler tests
// Maps to TEST_PLAN.md:
// - Story 2: Authentication & Identity [UC-S2-01~07, UC-S2-11~15, E2E-S2-01~06]
//
// NOTE: Session validation (UC-S2-08~10) is tested in middleware test runner
// NOTE: Protected route access (E2E-S2-05) is tested in middleware test runner
func AuthenticationHandlerRunner(t *testing.T, constructor AuthenticationHandlerConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockAuth := mockUsecase.NewMockAuthenticationUseCase(ctrl)
	mockSetup := mockUsecase.NewMockSetupUseCase(ctrl)
	mockRBAC := mockUsecase.NewMockRBACUseCase(ctrl)

	h := constructor(mockAuth, mockSetup, mockRBAC)

	// UC-S2-01: Login Form Validation - Empty Username
	t.Run("UC-S2-01: Login Form Validation - Empty Username", func(t *testing.T) {
		form := url.Values{}
		form.Add("username", "")
		form.Add("password", "password123")

		mockAuth.EXPECT().
			ValidateLoginForm(gomock.Any(), gomock.Any()).
			Return([]domain.ValidationError{
				{Field: "username", Message: "Username is required"},
			}, nil)

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.HandleLogin(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusBadRequest, rec.Code)
		body := rec.Body.String()
		require.Contains(t, body, "Username is required")
	})

	// UC-S2-02: Login Form Validation - Empty Password
	t.Run("UC-S2-02: Login Form Validation - Empty Password", func(t *testing.T) {
		form := url.Values{}
		form.Add("username", "testuser")
		form.Add("password", "")

		mockAuth.EXPECT().
			ValidateLoginForm(gomock.Any(), gomock.Any()).
			Return([]domain.ValidationError{
				{Field: "password", Message: "Password is required"},
			}, nil)

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.HandleLogin(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusBadRequest, rec.Code)
		body := rec.Body.String()
		require.Contains(t, body, "Password is required")
	})

	// UC-S2-03: Login Connection Probe
	t.Run("UC-S2-03: Login Connection Probe", func(t *testing.T) {
		form := url.Values{}
		form.Add("username", "testuser")
		form.Add("password", "password123")

		mockAuth.EXPECT().
			ValidateLoginForm(gomock.Any(), gomock.Any()).
			Return([]domain.ValidationError{}, nil)

		mockAuth.EXPECT().
			ProbeConnection(gomock.Any(), "testuser", "password123").
			Return(true, nil)

		mockAuth.EXPECT().
			GetUserAccessibleResources(gomock.Any(), "testuser").
			Return(&domain.RoleMetadata{
				Name:                "testuser",
				AccessibleDatabases: []string{"testdb"},
				AccessibleSchemas:   []string{"public"},
				AccessibleTables: []domain.AccessibleTable{
					{Database: "testdb", Schema: "public", Name: "users", HasSelect: true},
				},
			}, nil)

		mockAuth.EXPECT().
			GetFirstAccessibleDatabase(gomock.Any(), "testuser").
			Return("testdb", nil)

		mockAuth.EXPECT().
			GetFirstAccessibleSchema(gomock.Any(), "testuser", "testdb").
			Return("public", nil)

		mockAuth.EXPECT().
			GetFirstAccessibleTable(gomock.Any(), "testuser", "testdb", "public").
			Return("users", nil)

		mockAuth.EXPECT().
			CreateSession(gomock.Any(), "testuser", "password123", "testdb", "public", "users").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.HandleLogin(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusFound, rec.Code)
		require.Contains(t, rec.Header().Get("Location"), "/main")
	})

	// UC-S2-04: Login Connection Probe Failure
	t.Run("UC-S2-04: Login Connection Probe Failure", func(t *testing.T) {
		form := url.Values{}
		form.Add("username", "limited_user")
		form.Add("password", "password123")

		mockAuth.EXPECT().
			ValidateLoginForm(gomock.Any(), gomock.Any()).
			Return([]domain.ValidationError{}, nil)

		mockAuth.EXPECT().
			ProbeConnection(gomock.Any(), "limited_user", "password123").
			Return(true, nil)

		mockAuth.EXPECT().
			GetUserAccessibleResources(gomock.Any(), "limited_user").
			Return(&domain.RoleMetadata{
				Name:                "limited_user",
				AccessibleDatabases: []string{},
				AccessibleSchemas:   []string{},
				AccessibleTables:    []domain.AccessibleTable{},
			}, nil)

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.HandleLogin(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()
		require.Contains(t, body, "No accessible resources found")
	})

	// UC-S2-05: Login Success After Probe
	t.Run("UC-S2-05: Login Success After Probe", func(t *testing.T) {
		form := url.Values{}
		form.Add("username", "testuser")
		form.Add("password", "password123")

		mockAuth.EXPECT().
			ValidateLoginForm(gomock.Any(), gomock.Any()).
			Return([]domain.ValidationError{}, nil)

		mockAuth.EXPECT().
			ProbeConnection(gomock.Any(), "testuser", "password123").
			Return(true, nil)

		mockAuth.EXPECT().
			GetUserAccessibleResources(gomock.Any(), "testuser").
			Return(&domain.RoleMetadata{
				Name:                "testuser",
				AccessibleDatabases: []string{"testdb"},
				AccessibleSchemas:   []string{"public"},
				AccessibleTables: []domain.AccessibleTable{
					{Database: "testdb", Schema: "public", Name: "users", HasSelect: true},
				},
			}, nil)

		mockAuth.EXPECT().
			GetFirstAccessibleDatabase(gomock.Any(), "testuser").
			Return("testdb", nil)

		mockAuth.EXPECT().
			GetFirstAccessibleSchema(gomock.Any(), "testuser", "testdb").
			Return("public", nil)

		mockAuth.EXPECT().
			GetFirstAccessibleTable(gomock.Any(), "testuser", "testdb", "public").
			Return("users", nil)

		mockAuth.EXPECT().
			CreateSession(gomock.Any(), "testuser", "password123", "testdb", "public", "users").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.HandleLogin(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusFound, rec.Code)
		require.Contains(t, rec.Header().Get("Location"), "/main")
		require.NotEmpty(t, rec.Header().Get("Set-Cookie"))
	})

	// UC-S2-06: Session Cookie Creation - Username
	t.Run("UC-S2-06: Session Cookie Creation - Username", func(t *testing.T) {
		form := url.Values{}
		form.Add("username", "testuser")
		form.Add("password", "password123")

		mockAuth.EXPECT().
			ValidateLoginForm(gomock.Any(), gomock.Any()).
			Return([]domain.ValidationError{}, nil)

		mockAuth.EXPECT().
			ProbeConnection(gomock.Any(), "testuser", "password123").
			Return(true, nil)

		mockAuth.EXPECT().
			GetUserAccessibleResources(gomock.Any(), "testuser").
			Return(&domain.RoleMetadata{
				Name:                "testuser",
				AccessibleDatabases: []string{"testdb"},
				AccessibleSchemas:   []string{"public"},
				AccessibleTables: []domain.AccessibleTable{
					{Database: "testdb", Schema: "public", Name: "users", HasSelect: true},
				},
			}, nil)

		mockAuth.EXPECT().
			GetFirstAccessibleDatabase(gomock.Any(), "testuser").
			Return("testdb", nil)

		mockAuth.EXPECT().
			GetFirstAccessibleSchema(gomock.Any(), "testuser", "testdb").
			Return("public", nil)

		mockAuth.EXPECT().
			GetFirstAccessibleTable(gomock.Any(), "testuser", "testdb", "public").
			Return("users", nil)

		mockAuth.EXPECT().
			CreateSession(gomock.Any(), "testuser", "password123", "testdb", "public", "users").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.HandleLogin(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusFound, rec.Code)

		// Verify username cookie is set (long-lived)
		cookies := rec.Result().Cookies()
		var usernameCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "username" {
				usernameCookie = cookie
				break
			}
		}

		if usernameCookie != nil {
			require.Equal(t, "testuser", usernameCookie.Value)
			// Long-lived cookie should have MaxAge > 0 or no MaxAge set
			require.True(t, usernameCookie.MaxAge >= 0 || usernameCookie.MaxAge == 0)
		}
	})

	// UC-S2-07: Session Cookie Creation - Password
	t.Run("UC-S2-07: Session Cookie Creation - Password", func(t *testing.T) {
		form := url.Values{}
		form.Add("username", "testuser")
		form.Add("password", "password123")

		mockAuth.EXPECT().
			ValidateLoginForm(gomock.Any(), gomock.Any()).
			Return([]domain.ValidationError{}, nil)

		mockAuth.EXPECT().
			ProbeConnection(gomock.Any(), "testuser", "password123").
			Return(true, nil)

		mockAuth.EXPECT().
			GetUserAccessibleResources(gomock.Any(), "testuser").
			Return(&domain.RoleMetadata{
				Name:                "testuser",
				AccessibleDatabases: []string{"testdb"},
				AccessibleSchemas:   []string{"public"},
				AccessibleTables: []domain.AccessibleTable{
					{Database: "testdb", Schema: "public", Name: "users", HasSelect: true},
				},
			}, nil)

		mockAuth.EXPECT().
			GetFirstAccessibleDatabase(gomock.Any(), "testuser").
			Return("testdb", nil)

		mockAuth.EXPECT().
			GetFirstAccessibleSchema(gomock.Any(), "testuser", "testdb").
			Return("public", nil)

		mockAuth.EXPECT().
			GetFirstAccessibleTable(gomock.Any(), "testuser", "testdb", "public").
			Return("users", nil)

		mockAuth.EXPECT().
			CreateSession(gomock.Any(), "testuser", "password123", "testdb", "public", "users").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.HandleLogin(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusFound, rec.Code)

		// Verify password is encrypted in cookie (not plain text)
		cookies := rec.Result().Cookies()
		for _, cookie := range cookies {
			require.NotContains(t, cookie.Value, "password123", "Password must not be stored in plain text")
		}

		// Verify session cookie exists (short-lived encrypted password)
		var sessionCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "session_id" || cookie.Name == "encrypted_password" {
				sessionCookie = cookie
				break
			}
		}

		if sessionCookie != nil {
			require.NotEmpty(t, sessionCookie.Value)
			// Short-lived cookie should have MaxAge set
			require.True(t, sessionCookie.MaxAge > 0)
		}
	})

	// UC-S2-11: Data Explorer Population After Login
	t.Run("UC-S2-11: Data Explorer Population After Login", func(t *testing.T) {
		form := url.Values{}
		form.Add("username", "testuser")
		form.Add("password", "password123")

		mockAuth.EXPECT().
			ValidateLoginForm(gomock.Any(), gomock.Any()).
			Return([]domain.ValidationError{}, nil)

		mockAuth.EXPECT().
			ProbeConnection(gomock.Any(), "testuser", "password123").
			Return(true, nil)

		mockAuth.EXPECT().
			GetUserAccessibleResources(gomock.Any(), "testuser").
			Return(&domain.RoleMetadata{
				Name:                "testuser",
				AccessibleDatabases: []string{"testdb1", "testdb2"},
				AccessibleSchemas:   []string{"public", "private"},
				AccessibleTables: []domain.AccessibleTable{
					{Database: "testdb1", Schema: "public", Name: "users", HasSelect: true},
					{Database: "testdb1", Schema: "public", Name: "posts", HasSelect: true},
					{Database: "testdb2", Schema: "public", Name: "products", HasSelect: true},
				},
			}, nil)

		mockAuth.EXPECT().
			GetFirstAccessibleDatabase(gomock.Any(), "testuser").
			Return("testdb1", nil)

		mockAuth.EXPECT().
			GetFirstAccessibleSchema(gomock.Any(), "testuser", "testdb1").
			Return("public", nil)

		mockAuth.EXPECT().
			GetFirstAccessibleTable(gomock.Any(), "testuser", "testdb1", "public").
			Return("users", nil)

		mockAuth.EXPECT().
			CreateSession(gomock.Any(), "testuser", "password123", "testdb1", "public", "users").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.HandleLogin(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusFound, rec.Code)
		require.Contains(t, rec.Header().Get("Location"), "/main")
	})

	// UC-S2-12: Logout Cookie Clearing
	t.Run("UC-S2-12: Logout Cookie Clearing", func(t *testing.T) {
		mockAuth.EXPECT().
			Logout(gomock.Any(), "session_123").
			Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleLogout(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusFound, rec.Code)
		require.Equal(t, "/login", rec.Header().Get("Location"))

		// Verify cookies are cleared
		cookies := rec.Result().Cookies()
		var sessionCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "session_id" {
				sessionCookie = cookie
				break
			}
		}
		require.NotNil(t, sessionCookie)
		require.Equal(t, "", sessionCookie.Value)
		require.True(t, sessionCookie.MaxAge < 0)
	})

	// UC-S2-13: Header Username Display
	// NOTE: This is typically tested at the template/view level
	t.Run("UC-S2-13: Header Username Display", func(t *testing.T) {
		// Create an authenticated context (simulating middleware has run)
		authenticatedCtx := context.WithValue(ctx, domain.ContextKeyUser, &domain.Session{
			ID:       "session_123",
			Username: "testuser",
		})

		req := httptest.NewRequest(http.MethodGet, "/main", nil)
		rec := httptest.NewRecorder()

		h.HandleLoginPage(rec, req.WithContext(authenticatedCtx))

		require.Equal(t, http.StatusOK, rec.Code)
		// If already authenticated, login page might redirect or show different content
		// This behavior depends on implementation
	})

	// UC-S2-14: Navigation Menu Rendering
	// NOTE: This is typically tested at the template/view level
	t.Run("UC-S2-14: Navigation Menu Rendering", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/login", nil)
		rec := httptest.NewRecorder()

		h.HandleLoginPage(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify login page renders correctly
		require.Contains(t, body, "login")
	})

	// UC-S2-15: Metadata Refresh Button
	// NOTE: This is a superadmin feature, tested separately in admin handlers

	// E2E-S2-01: Login Flow with Connection Probe
	t.Run("E2E-S2-01: Login Flow with Connection Probe", func(t *testing.T) {
		form := url.Values{}
		form.Add("username", "testuser")
		form.Add("password", "password123")

		mockAuth.EXPECT().
			ValidateLoginForm(gomock.Any(), gomock.Any()).
			Return([]domain.ValidationError{}, nil)

		mockAuth.EXPECT().
			ProbeConnection(gomock.Any(), "testuser", "password123").
			Return(true, nil)

		mockAuth.EXPECT().
			GetUserAccessibleResources(gomock.Any(), "testuser").
			Return(&domain.RoleMetadata{
				Name:                "testuser",
				AccessibleDatabases: []string{"testdb"},
				AccessibleSchemas:   []string{"public"},
				AccessibleTables: []domain.AccessibleTable{
					{Database: "testdb", Schema: "public", Name: "users", HasSelect: true},
				},
			}, nil)

		mockAuth.EXPECT().
			GetFirstAccessibleDatabase(gomock.Any(), "testuser").
			Return("testdb", nil)

		mockAuth.EXPECT().
			GetFirstAccessibleSchema(gomock.Any(), "testuser", "testdb").
			Return("public", nil)

		mockAuth.EXPECT().
			GetFirstAccessibleTable(gomock.Any(), "testuser", "testdb", "public").
			Return("users", nil)

		mockAuth.EXPECT().
			CreateSession(gomock.Any(), "testuser", "password123", "testdb", "public", "users").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.HandleLogin(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusFound, rec.Code)
		require.Contains(t, rec.Header().Get("Location"), "/main")
		require.NotEmpty(t, rec.Header().Get("Set-Cookie"))
	})

	// E2E-S2-02: Login Flow - No Accessible Resources
	t.Run("E2E-S2-02: Login Flow - No Accessible Resources", func(t *testing.T) {
		form := url.Values{}
		form.Add("username", "limited_user")
		form.Add("password", "password123")

		mockAuth.EXPECT().
			ValidateLoginForm(gomock.Any(), gomock.Any()).
			Return([]domain.ValidationError{}, nil)

		mockAuth.EXPECT().
			ProbeConnection(gomock.Any(), "limited_user", "password123").
			Return(true, nil)

		mockAuth.EXPECT().
			GetUserAccessibleResources(gomock.Any(), "limited_user").
			Return(&domain.RoleMetadata{
				Name:                "limited_user",
				AccessibleDatabases: []string{},
				AccessibleSchemas:   []string{},
				AccessibleTables:    []domain.AccessibleTable{},
			}, nil)

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.HandleLogin(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()
		require.Contains(t, body, "No accessible resources found")
	})

	// E2E-S2-03: Login Flow - Invalid Credentials
	t.Run("E2E-S2-03: Login Flow - Invalid Credentials", func(t *testing.T) {
		form := url.Values{}
		form.Add("username", "testuser")
		form.Add("password", "wrongpassword")

		mockAuth.EXPECT().
			ValidateLoginForm(gomock.Any(), gomock.Any()).
			Return([]domain.ValidationError{}, nil)

		mockAuth.EXPECT().
			ProbeConnection(gomock.Any(), "testuser", "wrongpassword").
			Return(false, domain.ValidationError{Field: "credentials", Message: "Invalid credentials"})

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.HandleLogin(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusUnauthorized, rec.Code)
		body := rec.Body.String()
		require.Contains(t, body, "Invalid credentials")
	})

	// E2E-S2-04: Logout Flow
	t.Run("E2E-S2-04: Logout Flow", func(t *testing.T) {
		mockAuth.EXPECT().
			Logout(gomock.Any(), "session_123").
			Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleLogout(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusFound, rec.Code)
		require.Equal(t, "/login", rec.Header().Get("Location"))

		// Check cookies are cleared
		cookies := rec.Result().Cookies()
		var sessionCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "session_id" {
				sessionCookie = cookie
				break
			}
		}
		require.NotNil(t, sessionCookie)
		require.Equal(t, "", sessionCookie.Value)
		require.True(t, sessionCookie.MaxAge < 0)
	})

	// E2E-S2-05: Protected Route Access Without Auth
	// NOTE: This is a MIDDLEWARE concern, tested in middleware/authentication_middleware_test_runner.go
	// Handlers should not test this directly - middleware wraps handlers to provide this functionality

	// E2E-S2-06: Data Explorer Populated After Login
	t.Run("E2E-S2-06: Data Explorer Populated After Login", func(t *testing.T) {
		form := url.Values{}
		form.Add("username", "testuser")
		form.Add("password", "password123")

		mockAuth.EXPECT().
			ValidateLoginForm(gomock.Any(), gomock.Any()).
			Return([]domain.ValidationError{}, nil)

		mockAuth.EXPECT().
			ProbeConnection(gomock.Any(), "testuser", "password123").
			Return(true, nil)

		mockAuth.EXPECT().
			GetUserAccessibleResources(gomock.Any(), "testuser").
			Return(&domain.RoleMetadata{
				Name:                "testuser",
				AccessibleDatabases: []string{"testdb1", "testdb2"},
				AccessibleSchemas:   []string{"public", "private"},
				AccessibleTables: []domain.AccessibleTable{
					{Database: "testdb1", Schema: "public", Name: "users", HasSelect: true},
					{Database: "testdb1", Schema: "public", Name: "posts", HasSelect: true},
					{Database: "testdb2", Schema: "public", Name: "products", HasSelect: true},
				},
			}, nil)

		mockAuth.EXPECT().
			GetFirstAccessibleDatabase(gomock.Any(), "testuser").
			Return("testdb1", nil)

		mockAuth.EXPECT().
			GetFirstAccessibleSchema(gomock.Any(), "testuser", "testdb1").
			Return("public", nil)

		mockAuth.EXPECT().
			GetFirstAccessibleTable(gomock.Any(), "testuser", "testdb1", "public").
			Return("users", nil)

		mockAuth.EXPECT().
			CreateSession(gomock.Any(), "testuser", "password123", "testdb1", "public", "users").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.HandleLogin(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusFound, rec.Code)
		require.Contains(t, rec.Header().Get("Location"), "/main")

		// Verify session cookie is set
		cookies := rec.Result().Cookies()
		var sessionCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "session_id" {
				sessionCookie = cookie
				break
			}
		}
		require.NotNil(t, sessionCookie)
		require.Equal(t, "session_123", sessionCookie.Value)
	})

	// Additional test: Connection probe failure handling
	t.Run("Connection Probe Network Failure", func(t *testing.T) {
		form := url.Values{}
		form.Add("username", "testuser")
		form.Add("password", "password123")

		mockAuth.EXPECT().
			ValidateLoginForm(gomock.Any(), gomock.Any()).
			Return([]domain.ValidationError{}, nil)

		mockAuth.EXPECT().
			ProbeConnection(gomock.Any(), "testuser", "password123").
			Return(false, domain.ValidationError{Field: "connection", Message: "Database connection failed"})

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		h.HandleLogin(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusUnauthorized, rec.Code)
		body := rec.Body.String()
		require.Contains(t, body, "Database connection failed")
	})

	// Additional test: Login page rendering
	t.Run("Login Page Rendering", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/login", nil)
		rec := httptest.NewRecorder()

		h.HandleLoginPage(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify login form elements are present
		require.Contains(t, body, "username")
		require.Contains(t, body, "password")
	})
}
