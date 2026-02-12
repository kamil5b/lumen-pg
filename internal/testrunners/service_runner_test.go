package testrunners

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/interfaces"
	repomock "github.com/kamil5b/lumen-pg/internal/implementations/mocks"
)

type ServiceConstructor func(repo interfaces.UserRepository) interfaces.UserService

func UserServiceRunner(t *testing.T, constructor ServiceConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repomock.NewMockUserRepository(ctrl)

	svc := constructor(mockRepo)

	t.Run("CreateUser generates ID and saves", func(t *testing.T) {
		input := domain.CreateUserInput{
			Email: "service@test.com",
			Name:  "Service",
		}

		mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)

		user, err := svc.CreateUser(context.Background(), input)

		require.NoError(t, err)
		assert.NotEmpty(t, user.ID)
		assert.Equal(t, input.Email, user.Email)
	})
}
