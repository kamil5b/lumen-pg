package testrunners

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
	"github.com/kamil5b/lumen-pg/internal/implementations/mocks"
)

// ConnectionUseCaseConstructor creates a connection use case with its dependencies
type ConnectionUseCaseConstructor func(repo repository.ConnectionRepository) usecase.ConnectionUseCase

// ConnectionUseCaseRunner runs test specs for connection use case (Story 1)
func ConnectionUseCaseRunner(t *testing.T, constructor ConnectionUseCaseConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockConnectionRepository(ctrl)
	useCase := constructor(mockRepo)

	t.Run("UC-S1-01: Connection String Validation - Invalid Format", func(t *testing.T) {
		invalidConnStr := "invalid connection string"

		mockRepo.EXPECT().ParseConnectionString(invalidConnStr).Return(nil, assert.AnError)

		err := useCase.ValidateConnectionString(invalidConnStr)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "validation error")
	})

	t.Run("UC-S1-02: Connection String Parsing - Valid", func(t *testing.T) {
		validConnStr := "postgres://user:pass@localhost:5432/testdb?sslmode=disable"
		expectedConn := &domain.Connection{
			Host:     "localhost",
			Port:     5432,
			Database: "testdb",
			User:     "user",
			Password: "pass",
			SSLMode:  "disable",
		}

		mockRepo.EXPECT().ParseConnectionString(validConnStr).Return(expectedConn, nil)

		err := useCase.ValidateConnectionString(validConnStr)

		require.NoError(t, err)
	})

	t.Run("UC-S1-03: Superadmin Connection Test Success", func(t *testing.T) {
		ctx := context.Background()
		connStr := "postgres://admin:pass@localhost:5432/postgres"
		conn := &domain.Connection{
			Host:     "localhost",
			Port:     5432,
			Database: "postgres",
			User:     "admin",
			Password: "pass",
		}

		mockRepo.EXPECT().ParseConnectionString(connStr).Return(conn, nil)
		mockRepo.EXPECT().ValidateConnection(ctx, conn).Return(nil)

		err := useCase.TestConnection(ctx, connStr)

		require.NoError(t, err)
	})

	t.Run("UC-S1-04: Superadmin Connection Test Failure", func(t *testing.T) {
		ctx := context.Background()
		connStr := "postgres://admin:wrongpass@localhost:5432/postgres"
		conn := &domain.Connection{
			Host:     "localhost",
			Port:     5432,
			Database: "postgres",
			User:     "admin",
			Password: "wrongpass",
		}

		mockRepo.EXPECT().ParseConnectionString(connStr).Return(conn, nil)
		mockRepo.EXPECT().ValidateConnection(ctx, conn).Return(assert.AnError)

		err := useCase.TestConnection(ctx, connStr)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "connection error")
	})
}
