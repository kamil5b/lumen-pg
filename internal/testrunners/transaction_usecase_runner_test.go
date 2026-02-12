package testrunners

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/interfaces"
	"github.com/kamil5b/lumen-pg/internal/implementations/mocks"
)

// TransactionUseCaseConstructor creates a transaction use case with its dependencies
type TransactionUseCaseConstructor func(repo interfaces.TransactionRepository) TransactionUseCase

// TransactionUseCase represents the transaction management use case
type TransactionUseCase interface {
	StartTransaction(ctx context.Context, username string, tableName string) (*domain.Transaction, error)
	BufferEdit(ctx context.Context, txnID string, op domain.TransactionOperation) error
	CommitTransaction(ctx context.Context, txnID string) error
	RollbackTransaction(ctx context.Context, txnID string) error
}

// TransactionUseCaseRunner runs test specs for transaction use case (Story 5)
func TransactionUseCaseRunner(t *testing.T, constructor TransactionUseCaseConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)
	useCase := constructor(mockRepo)

	t.Run("UC-S5-09: Transaction Start", func(t *testing.T) {
		ctx := context.Background()
		username := "testuser"
		tableName := "users"
		expectedTxn := &domain.Transaction{
			ID:        "txn-123",
			Username:  username,
			TableName: tableName,
			StartedAt: time.Now(),
			ExpiresAt: time.Now().Add(1 * time.Minute),
		}

		mockRepo.EXPECT().StartTransaction(ctx, username, tableName).Return(expectedTxn, nil)

		txn, err := useCase.StartTransaction(ctx, username, tableName)

		require.NoError(t, err)
		assert.NotNil(t, txn)
		assert.Equal(t, username, txn.Username)
		assert.Equal(t, tableName, txn.TableName)
		assert.True(t, txn.ExpiresAt.After(txn.StartedAt))
	})

	t.Run("UC-S5-11: Cell Edit Buffering", func(t *testing.T) {
		ctx := context.Background()
		txnID := "txn-123"
		op := domain.TransactionOperation{
			Type:       domain.OperationUpdate,
			TableName:  "users",
			PrimaryKey: 1,
			Column:     "username",
			OldValue:   "oldname",
			NewValue:   "newname",
		}

		mockRepo.EXPECT().BufferOperation(ctx, txnID, op).Return(nil)

		err := useCase.BufferEdit(ctx, txnID, op)

		require.NoError(t, err)
	})

	t.Run("UC-S5-12: Transaction Commit", func(t *testing.T) {
		ctx := context.Background()
		txnID := "txn-123"

		mockRepo.EXPECT().CommitTransaction(ctx, txnID).Return(nil)

		err := useCase.CommitTransaction(ctx, txnID)

		require.NoError(t, err)
	})

	t.Run("UC-S5-13: Transaction Rollback", func(t *testing.T) {
		ctx := context.Background()
		txnID := "txn-123"

		mockRepo.EXPECT().RollbackTransaction(ctx, txnID).Return(nil)

		err := useCase.RollbackTransaction(ctx, txnID)

		require.NoError(t, err)
	})

	t.Run("UC-S5-15: Row Deletion Buffering", func(t *testing.T) {
		ctx := context.Background()
		txnID := "txn-123"
		op := domain.TransactionOperation{
			Type:       domain.OperationDelete,
			TableName:  "users",
			PrimaryKey: 1,
		}

		mockRepo.EXPECT().BufferOperation(ctx, txnID, op).Return(nil)

		err := useCase.BufferEdit(ctx, txnID, op)

		require.NoError(t, err)
	})

	t.Run("UC-S5-16: Row Insertion Buffering", func(t *testing.T) {
		ctx := context.Background()
		txnID := "txn-123"
		op := domain.TransactionOperation{
			Type:      domain.OperationInsert,
			TableName: "users",
			NewValue:  map[string]interface{}{"username": "newuser", "email": "new@test.com"},
		}

		mockRepo.EXPECT().BufferOperation(ctx, txnID, op).Return(nil)

		err := useCase.BufferEdit(ctx, txnID, op)

		require.NoError(t, err)
	})
}
