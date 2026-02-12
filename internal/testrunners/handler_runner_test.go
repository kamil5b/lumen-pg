package testrunners

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/kamil5b/lumen-pg/internal/domain"
	svcmock "github.com/kamil5b/lumen-pg/internal/implementations/mocks"
	"github.com/kamil5b/lumen-pg/internal/interfaces"
)

type HandlerConstructor func(svc interfaces.UserService) interfaces.UserHandler

func UserHandlerRunner(t *testing.T, constructor HandlerConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := svcmock.NewMockUserService(ctrl)

	h := constructor(mockSvc)

	r := chi.NewRouter()
	h.RegisterRoutes(r)

	t.Run("POST /users returns 201", func(t *testing.T) {
		input := domain.CreateUserInput{
			Email: "http@test.com",
			Name:  "HTTP",
		}

		mockSvc.EXPECT().CreateUser(gomock.Any(), input).Return(&domain.User{
			ID:    "123",
			Email: input.Email,
			Name:  input.Name,
		}, nil)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		require.Equal(t, http.StatusCreated, rec.Code)
	})
}
