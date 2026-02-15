package e2e

import (
	"context"
	"crypto/tls"
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

// SecurityHandlerConstructor is a struct that holds handlers for security testing
type SecurityHandlerConstructor struct {
	LoginHandler    handler.LoginHandler
	QueryHandler    handler.QueryEditorHandler
	MainViewHandler handler.MainViewHandler
}

// SecurityHandlerConstructorFunc creates all necessary handlers for security testing
type SecurityHandlerConstructorFunc func(
	authUC usecase.AuthenticationUseCase,
	securityUC usecase.SecurityUseCase,
	queryUC usecase.QueryUseCase,
	dataViewUC usecase.DataViewUseCase,
	rbacUC usecase.RBACUseCase,
	setupUC usecase.SetupUseCase,
) SecurityHandlerConstructor

// SecurityHandlerRunner runs all security E2E handler tests
// Maps to TEST_PLAN.md:
// - Story 7: Security & Best Practices [E2E-S7-01~06]
func SecurityHandlerRunner(t *testing.T, constructor SecurityHandlerConstructorFunc) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockAuth := mockUsecase.NewMockAuthenticationUseCase(ctrl)
	mockSecurity := mockUsecase.NewMockSecurityUseCase(ctrl)
	mockQuery := mockUsecase.NewMockQueryUseCase(ctrl)
	mockDataView := mockUsecase.NewMockDataViewUseCase(ctrl)
	mockRBAC := mockUsecase.NewMockRBACUseCase(ctrl)
	mockSetup := mockUsecase.NewMockSetupUseCase(ctrl)

	handlers := constructor(mockAuth, mockSecurity, mockQuery, mockDataView, mockRBAC, mockSetup)

	// E2E-S7-01: SQL Injection via WHERE Bar
	t.Run("E2E-S7-01: SQL Injection via WHERE Bar", func(t *testing.T) {
		form := url.Values{}
		form.Add("where", "1=1; DROP TABLE users; --")
		form.Add("database", "testdb")
		form.Add("schema", "public")
		form.Add("table", "users")

		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockDataView.EXPECT().
			ValidateWhereClause(gomock.Any(), "1=1; DROP TABLE users; --").
			Return(false, domain.ValidationError{Field: "where", Message: "SQL injection detected"})

		req := httptest.NewRequest(http.MethodPost, "/main/filter", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		handlers.MainViewHandler.HandleFilterTable(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusBadRequest, rec.Code)
		body := rec.Body.String()

		// Verify injection is prevented
		require.Contains(t, body, "SQL injection")
		require.NotContains(t, body, "DROP TABLE")
	})

	// E2E-S7-02: SQL Injection via Query Editor
	t.Run("E2E-S7-02: SQL Injection via Query Editor", func(t *testing.T) {
		form := url.Values{}
		form.Add("query", "SELECT * FROM users WHERE id = 1; DROP TABLE users; --")

		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockQuery.EXPECT().
			ExecuteQuery(gomock.Any(), "testuser", "SELECT * FROM users WHERE id = 1; DROP TABLE users; --", 0, 50).
			Return(nil, domain.ValidationError{Field: "query", Message: "SQL injection detected: multiple statements not allowed"})

		req := httptest.NewRequest(http.MethodPost, "/api/query/execute", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		handlers.QueryHandler.HandleExecuteQuery(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusBadRequest, rec.Code)
		body := rec.Body.String()

		// Verify injection attempt is blocked
		require.Contains(t, body, "SQL injection")
	})

	// E2E-S7-03: Cookie Tampering Prevention
	t.Run("E2E-S7-03: Cookie Tampering Prevention", func(t *testing.T) {
		// Simulate tampered cookie value
		tamperedCookie := "tampered_session_value_123"

		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), tamperedCookie).
			Return(nil, domain.ValidationError{Field: "session", Message: "Invalid session"})

		req := httptest.NewRequest(http.MethodGet, "/main", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: tamperedCookie,
		})
		rec := httptest.NewRecorder()

		handlers.MainViewHandler.HandleMainViewPage(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	// E2E-S7-04: Session Timeout Enforcement
	t.Run("E2E-S7-04: Session Timeout Enforcement", func(t *testing.T) {
		expiredSession := "expired_session_123"

		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), expiredSession).
			Return(nil, domain.ValidationError{Field: "session", Message: "Session expired"})

		req := httptest.NewRequest(http.MethodGet, "/main", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: expiredSession,
		})
		rec := httptest.NewRecorder()

		handlers.MainViewHandler.HandleMainViewPage(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	// E2E-S7-05: HTTPS-Only Cookies (if HTTPS enabled)
	t.Run("E2E-S7-05: HTTPS-Only Cookies", func(t *testing.T) {
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

		// Simulate HTTPS request
		req := httptest.NewRequest(http.MethodPost, "https://localhost/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.TLS = &tls.ConnectionState{} // Indicates HTTPS
		rec := httptest.NewRecorder()

		handlers.LoginHandler.HandleLogin(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusFound, rec.Code)

		// Verify Secure flag is set on cookies
		cookies := rec.Result().Cookies()
		var sessionCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "session_id" {
				sessionCookie = cookie
				break
			}
		}

		if sessionCookie != nil {
			require.True(t, sessionCookie.Secure, "Cookie should have Secure flag when using HTTPS")
		}
	})

	// E2E-S7-06: HTTPOnly Cookies
	t.Run("E2E-S7-06: HTTPOnly Cookies", func(t *testing.T) {
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

		handlers.LoginHandler.HandleLogin(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusFound, rec.Code)

		// Verify HTTPOnly flag is set on session cookie
		cookies := rec.Result().Cookies()
		var sessionCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "session_id" {
				sessionCookie = cookie
				break
			}
		}

		require.NotNil(t, sessionCookie)
		require.True(t, sessionCookie.HttpOnly, "Session cookie must have HttpOnly flag to prevent XSS attacks")
	})

	// Additional test: Password encryption in cookie
	t.Run("Password Encryption in Cookie", func(t *testing.T) {
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

		handlers.LoginHandler.HandleLogin(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusFound, rec.Code)

		// Verify password is not stored in plain text in cookie
		cookies := rec.Result().Cookies()
		for _, cookie := range cookies {
			// Check that no cookie contains plain password
			require.NotContains(t, cookie.Value, "password123", "Password must not be stored in plain text in cookies")
		}
	})

	// Additional test: Parameterized queries
	t.Run("Parameterized Query Execution", func(t *testing.T) {
		form := url.Values{}
		form.Add("query", "SELECT * FROM users WHERE id = $1")

		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockQuery.EXPECT().
			ExecuteQuery(gomock.Any(), "testuser", "SELECT * FROM users WHERE id = $1", 0, 50).
			Return(&domain.QueryResult{
				Columns: []string{"id", "name"},
				Rows: []map[string]interface{}{
					{"id": 1, "name": "Alice"},
				},
				RowCount:   1,
				TotalCount: 1,
			}, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/query/execute", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		handlers.QueryHandler.HandleExecuteQuery(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify parameterized query works
		require.Contains(t, body, "Alice")
	})

	// Additional test: SameSite cookie attribute
	t.Run("SameSite Cookie Attribute", func(t *testing.T) {
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

		handlers.LoginHandler.HandleLogin(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusFound, rec.Code)

		// Verify SameSite attribute is set (Strict or Lax)
		cookies := rec.Result().Cookies()
		var sessionCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "session_id" {
				sessionCookie = cookie
				break
			}
		}

		require.NotNil(t, sessionCookie)
		// SameSite should be set to prevent CSRF attacks
		require.NotEqual(t, http.SameSiteDefaultMode, sessionCookie.SameSite, "SameSite attribute should be explicitly set")
	})

	// Additional test: WHERE clause with quotes escaping
	t.Run("WHERE Clause Quote Escaping", func(t *testing.T) {
		form := url.Values{}
		form.Add("where", "name = 'O''Brien'")
		form.Add("database", "testdb")
		form.Add("schema", "public")
		form.Add("table", "users")

		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockDataView.EXPECT().
			ValidateWhereClause(gomock.Any(), "name = 'O''Brien'").
			Return(true, nil)

		mockDataView.EXPECT().
			FilterTableData(gomock.Any(), "testuser", "testdb", "public", "users", "name = 'O''Brien'", 0, 50).
			Return(&domain.QueryResult{
				Columns: []string{"id", "name"},
				Rows: []map[string]interface{}{
					{"id": 1, "name": "O'Brien"},
				},
				RowCount:   1,
				TotalCount: 1,
			}, nil)

		req := httptest.NewRequest(http.MethodPost, "/main/filter", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		handlers.MainViewHandler.HandleFilterTable(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify properly escaped quotes work
		require.Contains(t, body, "O'Brien")
	})
}
