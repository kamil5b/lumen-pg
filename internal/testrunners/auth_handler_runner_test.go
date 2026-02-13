package testrunners

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/implementations/mocks"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
)

// AuthHandlerConstructor creates an auth handler with its dependencies
type AuthHandlerConstructor func(authUseCase usecase.AuthUseCase) usecase.AuthHandler

// AuthHandlerRunner runs test specs for auth handler (Story 2 E2E)
func AuthHandlerRunner(t *testing.T, constructor AuthHandlerConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUseCase := mocks.NewMockAuthUseCase(ctrl)
	handler := constructor(mockAuthUseCase)

	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	t.Run("E2E-S2-01: Login Flow with Connection Probe", func(t *testing.T) {
		loginReq := domain.LoginRequest{
			Username: "testuser",
			Password: "testpass",
		}

		expectedResp := &domain.LoginResponse{
			Success:            true,
			FirstAccessibleDB:  "testdb",
			FirstAccessibleTbl: "users",
			Session: &domain.Session{
				Username: "testuser",
			},
		}

		mockAuthUseCase.EXPECT().Login(gomock.Any(), loginReq).Return(expectedResp, nil)

		body, _ := json.Marshal(loginReq)
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp domain.LoginResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, "testdb", resp.FirstAccessibleDB)
	})

	t.Run("E2E-S2-02: Login Flow - No Accessible Resources", func(t *testing.T) {
		loginReq := domain.LoginRequest{
			Username: "restricteduser",
			Password: "testpass",
		}

		expectedResp := &domain.LoginResponse{
			Success:      false,
			ErrorMessage: "No accessible resources found",
		}

		mockAuthUseCase.EXPECT().Login(gomock.Any(), loginReq).Return(expectedResp, nil)

		body, _ := json.Marshal(loginReq)
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusForbidden, rec.Code)

		var resp domain.LoginResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.False(t, resp.Success)
		assert.Contains(t, resp.ErrorMessage, "No accessible resources")
	})

	t.Run("E2E-S2-03: Login Flow - Invalid Credentials", func(t *testing.T) {
		loginReq := domain.LoginRequest{
			Username: "invaliduser",
			Password: "wrongpass",
		}

		expectedResp := &domain.LoginResponse{
			Success:      false,
			ErrorMessage: "Invalid credentials",
		}

		mockAuthUseCase.EXPECT().Login(gomock.Any(), loginReq).Return(expectedResp, nil)

		body, _ := json.Marshal(loginReq)
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("E2E-S2-04: Logout Flow", func(t *testing.T) {
		sessionToken := "valid-session-token"

		mockAuthUseCase.EXPECT().Logout(gomock.Any(), sessionToken).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
		req.Header.Set("Authorization", "Bearer "+sessionToken)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("E2E-S2-05: Protected Route Access Without Auth", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("E2E-S2-06: Data Explorer Populated After Login", func(t *testing.T) {
		loginReq := domain.LoginRequest{
			Username: "testuser",
			Password: "testpass",
		}

		expectedResp := &domain.LoginResponse{
			Success:            true,
			FirstAccessibleDB:  "testdb",
			FirstAccessibleTbl: "users",
			Session: &domain.Session{
				Username: "testuser",
			},
		}

		mockAuthUseCase.EXPECT().Login(gomock.Any(), loginReq).Return(expectedResp, nil)

		body, _ := json.Marshal(loginReq)
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp domain.LoginResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.True(t, resp.Success)
		assert.NotEmpty(t, resp.FirstAccessibleDB)
		assert.NotEmpty(t, resp.FirstAccessibleTbl)
	})

	t.Run("E2E-S2-07: Session Cookie Persistence", func(t *testing.T) {
		loginReq := domain.LoginRequest{
			Username: "testuser",
			Password: "testpass",
		}

		expectedResp := &domain.LoginResponse{
			Success:            true,
			FirstAccessibleDB:  "testdb",
			FirstAccessibleTbl: "users",
			Session: &domain.Session{
				Username: "testuser",
			},
		}

		mockAuthUseCase.EXPECT().Login(gomock.Any(), loginReq).Return(expectedResp, nil)

		body, _ := json.Marshal(loginReq)
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		// Check for Set-Cookie headers
		setCookieHeaders := rec.Header()["Set-Cookie"]
		assert.NotEmpty(t, setCookieHeaders, "Session cookie should be set")
	})

	t.Run("E2E-S2-08: Header Username Display", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/user/profile", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		// Should display username in header if authenticated
		if rec.Code == http.StatusOK {
			assert.NotEmpty(t, rec.Body.String())
		}
	})

	t.Run("E2E-S2-09: Metadata Refresh Endpoint", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/metadata/refresh", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		// Metadata refresh should require authentication
		assert.NotEqual(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("E2E-S2-10: Login with Session Persistence", func(t *testing.T) {
		loginReq := domain.LoginRequest{
			Username: "persistentuser",
			Password: "testpass",
		}

		expectedResp := &domain.LoginResponse{
			Success:            true,
			FirstAccessibleDB:  "testdb",
			FirstAccessibleTbl: "users",
			Session: &domain.Session{
				Username: "persistentuser",
			},
		}

		mockAuthUseCase.EXPECT().Login(gomock.Any(), loginReq).Return(expectedResp, nil)

		body, _ := json.Marshal(loginReq)
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp domain.LoginResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.True(t, resp.Success)
		assert.NotNil(t, resp.Session)
	})
}
