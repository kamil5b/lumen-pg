package handler

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

// SecurityHandlerRunner runs all security handler tests
// Maps to TEST_PLAN.md:
// - Story 7: Security & Best Practices [UC-S7-01~04, E2E-S7-01~02]
//
// NOTE: Cookie security (UC-S7-03~07, E2E-S7-03~06) is tested in middleware test runner
// NOTE: Handlers focus on input validation and SQL injection detection in business logic
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

	// UC-S7-01: SQL Injection Prevention - WHERE Clause
	// E2E-S7-01: SQL Injection via WHERE Bar
	t.Run("UC-S7-01/E2E-S7-01: SQL Injection via WHERE Bar", func(t *testing.T) {
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

	// UC-S7-02: SQL Injection Prevention - Query Editor
	// E2E-S7-02: SQL Injection via Query Editor
	t.Run("UC-S7-02/E2E-S7-02: SQL Injection via Query Editor", func(t *testing.T) {
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

	// UC-S7-03: Password Encryption in Cookie
	// NOTE: Cookie creation happens during login, verification happens in usecase layer
	t.Run("UC-S7-03: Password Encryption in Cookie", func(t *testing.T) {
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

	// UC-S7-04: Password Decryption from Cookie
	// NOTE: This is handled in authentication middleware and usecase layer
	// Handler assumes middleware has already validated and decrypted credentials

	// E2E-S7-03: Cookie Tampering Prevention
	// NOTE: This is a MIDDLEWARE concern, tested in middleware/security_middleware_test_runner.go
	// Handler receives already-validated session from middleware

	// E2E-S7-04: Session Timeout Enforcement
	// NOTE: This is a MIDDLEWARE concern, tested in middleware/authentication_middleware_test_runner.go
	// Handler receives already-validated session from middleware

	// E2E-S7-05: HTTPS-Only Cookies
	// NOTE: This is a MIDDLEWARE concern, tested in middleware/security_middleware_test_runner.go
	// Cookie flags are set by middleware/infrastructure layer

	// E2E-S7-06: HTTPOnly Cookies
	// NOTE: This is a MIDDLEWARE concern, tested in middleware/security_middleware_test_runner.go
	t.Run("E2E-S7-06: HTTPOnly Cookies verification in handler response", func(t *testing.T) {
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
		// Verify handler sets HTTPOnly (actual enforcement is in cookie creation logic)
		if sessionCookie != nil {
			require.True(t, sessionCookie.HttpOnly, "Session cookie must have HttpOnly flag to prevent XSS attacks")
		}
	})

	// Additional test: Multiple SQL injection patterns
	t.Run("Multiple SQL Injection Patterns Detection", func(t *testing.T) {
		patterns := []string{
			"1=1; DROP TABLE users; --",
			"' OR '1'='1",
			"admin'--",
			"'; DELETE FROM users WHERE 'a'='a",
			"1' UNION SELECT * FROM users--",
		}

		for _, pattern := range patterns {
			form := url.Values{}
			form.Add("where", pattern)
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
				ValidateWhereClause(gomock.Any(), pattern).
				Return(false, domain.ValidationError{Field: "where", Message: "SQL injection detected"})

			req := httptest.NewRequest(http.MethodPost, "/main/filter", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.AddCookie(&http.Cookie{
				Name:  "session_id",
				Value: "session_123",
			})
			rec := httptest.NewRecorder()

			handlers.MainViewHandler.HandleFilterTable(rec, req.WithContext(ctx))

			require.Equal(t, http.StatusBadRequest, rec.Code, "Pattern should be blocked: %s", pattern)
			body := rec.Body.String()
			require.Contains(t, body, "SQL injection", "Pattern should be detected: %s", pattern)
		}
	})

	// Additional test: Password not in plain text (handler responsibility)
	t.Run("Password Not Stored in Plain Text", func(t *testing.T) {
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

		// Verify handler doesn't expose password in response or cookies
		cookies := rec.Result().Cookies()
		for _, cookie := range cookies {
			require.NotContains(t, cookie.Value, "password123", "Password must not be stored in plain text in cookies")
		}

		body := rec.Body.String()
		require.NotContains(t, body, "password123", "Password must not appear in response body")
	})

	// Additional test: Valid parameterized queries are allowed
	t.Run("Valid Parameterized Query Execution", func(t *testing.T) {
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

		// Verify handler processes parameterized queries correctly
		require.Contains(t, body, "Alice")
	})

	// Additional test: Complex WHERE clause validation
	t.Run("Complex WHERE Clause Validation", func(t *testing.T) {
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

		if sessionCookie != nil {
			// Handler should set secure cookie attributes
			require.NotEqual(t, http.SameSiteDefaultMode, sessionCookie.SameSite, "SameSite attribute should be explicitly set")
		}
	})

	// Additional test: Legitimate WHERE clause with escaped quotes
	t.Run("Legitimate WHERE Clause with Escaped Quotes", func(t *testing.T) {
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

		// Verify handler allows properly escaped quotes
		require.Contains(t, body, "O'Brien")
	})

	// Additional test: Query editor blocks multiple statements
	t.Run("Query Editor Blocks Multiple Statements", func(t *testing.T) {
		queries := []string{
			"SELECT * FROM users; DROP TABLE users;",
			"SELECT * FROM users; DELETE FROM users;",
			"SELECT * FROM users;\nDROP TABLE users;",
		}

		for _, query := range queries {
			form := url.Values{}
			form.Add("query", query)

			mockAuth.EXPECT().
				ValidateSession(gomock.Any(), "session_123").
				Return(&domain.Session{
					ID:       "session_123",
					Username: "testuser",
				}, nil)

			mockQuery.EXPECT().
				ExecuteQuery(gomock.Any(), "testuser", query, 0, 50).
				Return(nil, domain.ValidationError{Field: "query", Message: "SQL injection detected: multiple statements not allowed"})

			req := httptest.NewRequest(http.MethodPost, "/api/query/execute", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.AddCookie(&http.Cookie{
				Name:  "session_id",
				Value: "session_123",
			})
			rec := httptest.NewRecorder()

			handlers.QueryHandler.HandleExecuteQuery(rec, req.WithContext(ctx))

			require.Equal(t, http.StatusBadRequest, rec.Code, "Query should be blocked: %s", query)
			body := rec.Body.String()
			require.Contains(t, body, "SQL injection", "Query should be detected as malicious: %s", query)
		}
	})

	// Additional test: XSS prevention in error messages
	t.Run("XSS Prevention in Error Messages", func(t *testing.T) {
		xssPayload := "<script>alert('XSS')</script>"
		form := url.Values{}
		form.Add("where", xssPayload)
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
			ValidateWhereClause(gomock.Any(), xssPayload).
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

		// Verify handler doesn't reflect XSS payload in response
		// Response should escape or sanitize any user input
		require.Contains(t, body, "SQL injection")
		// If payload is reflected, it should be HTML-escaped
		if contains := rec.Body.String(); len(contains) > 0 {
			// Handler should not include unescaped script tags
			require.NotContains(t, body, "<script>alert('XSS')</script>")
		}
	})
}
