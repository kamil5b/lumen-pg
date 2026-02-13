package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/implementations/mocks"
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
)

// TransactionInterfaceConstructor creates a transaction repository with mock dependencies
type TransactionInterfaceConstructor func(ctrl *gomock.Controller) repository.TransactionRepository

// TransactionInterfaceRunner runs unit tests for transaction repository interface (Story 5)
func TransactionInterfaceRunner(t *testing.T, constructor TransactionInterfaceConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)
	ctx := context.Background()

	t.Run("UC-S5-09: Transaction Start", func(t *testing.T) {
		username := "testuser"
		tableName := "users"
		expectedTxn := &domain.Transaction{
			ID:        "txn-123",
			Username:  username,
			TableName: tableName,
			StartedAt: time.Now(),
			ExpiresAt: time.Now().Add(5 * time.Minute),
		}

		mockRepo.EXPECT().StartTransaction(ctx, username, tableName).Return(expectedTxn, nil)

		txn, err := mockRepo.StartTransaction(ctx, username, tableName)

		require.NoError(t, err)
		assert.NotNil(t, txn)
		assert.Equal(t, username, txn.Username)
		assert.Equal(t, tableName, txn.TableName)
		assert.True(t, txn.ExpiresAt.After(txn.StartedAt))
	})

	t.Run("UC-S5-10: Transaction Already Active Error", func(t *testing.T) {
		username := "testuser"
		tableName := "users"

		mockRepo.EXPECT().StartTransaction(ctx, username, tableName).Return(nil, assert.AnError)

		txn, err := mockRepo.StartTransaction(ctx, username, tableName)

		require.Error(t, err)
		assert.Nil(t, txn)
	})

	t.Run("UC-S5-11: Cell Edit Buffering", func(t *testing.T) {
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

		err := mockRepo.BufferOperation(ctx, txnID, op)

		require.NoError(t, err)
	})

	t.Run("UC-S5-12: Transaction Commit", func(t *testing.T) {
		txnID := "txn-123"

		mockRepo.EXPECT().CommitTransaction(ctx, txnID).Return(nil)

		err := mockRepo.CommitTransaction(ctx, txnID)

		require.NoError(t, err)
	})

	t.Run("UC-S5-13: Transaction Rollback", func(t *testing.T) {
		txnID := "txn-123"

		mockRepo.EXPECT().RollbackTransaction(ctx, txnID).Return(nil)

		err := mockRepo.RollbackTransaction(ctx, txnID)

		require.NoError(t, err)
	})

	t.Run("UC-S5-14: Transaction Timer Expiration", func(t *testing.T) {
		username := "testuser"
		tableName := "users"
		expiredTxn := &domain.Transaction{
			ID:        "txn-123",
			Username:  username,
			TableName: tableName,
			StartedAt: time.Now().Add(-6 * time.Minute),
			ExpiresAt: time.Now().Add(-1 * time.Minute),
		}

		mockRepo.EXPECT().GetTransaction(ctx, "txn-123").Return(expiredTxn, nil)

		txn, err := mockRepo.GetTransaction(ctx, "txn-123")

		require.NoError(t, err)
		assert.NotNil(t, txn)
		assert.True(t, time.Now().After(txn.ExpiresAt))
	})

	t.Run("UC-S5-15: Row Deletion Buffering", func(t *testing.T) {
		txnID := "txn-123"
		op := domain.TransactionOperation{
			Type:       domain.OperationDelete,
			TableName:  "users",
			PrimaryKey: 1,
		}

		mockRepo.EXPECT().BufferOperation(ctx, txnID, op).Return(nil)

		err := mockRepo.BufferOperation(ctx, txnID, op)

		require.NoError(t, err)
	})

	t.Run("UC-S5-16: Row Insertion Buffering", func(t *testing.T) {
		txnID := "txn-123"
		op := domain.TransactionOperation{
			Type:      domain.OperationInsert,
			TableName: "users",
			NewValue: map[string]interface{}{
				"username": "newuser",
				"email":    "new@test.com",
			},
		}

		mockRepo.EXPECT().BufferOperation(ctx, txnID, op).Return(nil)

		err := mockRepo.BufferOperation(ctx, txnID, op)

		require.NoError(t, err)
	})

	t.Run("UC-S5-17: Get Active Transaction", func(t *testing.T) {
		txnID := "txn-456"
		expectedTxn := &domain.Transaction{
			ID:        txnID,
			Username:  "testuser",
			TableName: "posts",
			StartedAt: time.Now(),
			ExpiresAt: time.Now().Add(5 * time.Minute),
			Operations: []domain.TransactionOperation{
				{
					Type:       domain.OperationUpdate,
					TableName:  "posts",
					PrimaryKey: 1,
					Column:     "title",
					OldValue:   "old title",
					NewValue:   "new title",
				},
			},
		}

		mockRepo.EXPECT().GetTransaction(ctx, txnID).Return(expectedTxn, nil)

		txn, err := mockRepo.GetTransaction(ctx, txnID)

		require.NoError(t, err)
		assert.NotNil(t, txn)
		assert.Len(t, txn.Operations, 1)
	})

	t.Run("UC-S5-18: Buffer Multiple Operations", func(t *testing.T) {
		txnID := "txn-123"
		ops := []domain.TransactionOperation{
			{
				Type:       domain.OperationUpdate,
				TableName:  "users",
				PrimaryKey: 1,
				Column:     "username",
				OldValue:   "oldname",
				NewValue:   "newname",
			},
			{
				Type:       domain.OperationUpdate,
				TableName:  "users",
				PrimaryKey: 1,
				Column:     "email",
				OldValue:   "old@test.com",
				NewValue:   "new@test.com",
			},
		}

		for _, op := range ops {
			mockRepo.EXPECT().BufferOperation(ctx, txnID, op).Return(nil)
		}

		for _, op := range ops {
			err := mockRepo.BufferOperation(ctx, txnID, op)
			require.NoError(t, err)
		}
	})

	t.Run("UC-S5-19: Transaction Isolation - Different Transactions", func(t *testing.T) {
		txn1 := &domain.Transaction{
			ID:        "txn-1",
			Username:  "user1",
			TableName: "users",
			StartedAt: time.Now(),
		}

		txn2 := &domain.Transaction{
			ID:        "txn-2",
			Username:  "user2",
			TableName: "posts",
			StartedAt: time.Now(),
		}

		mockRepo.EXPECT().GetTransaction(ctx, "txn-1").Return(txn1, nil)
		mockRepo.EXPECT().GetTransaction(ctx, "txn-2").Return(txn2, nil)

		result1, _ := mockRepo.GetTransaction(ctx, "txn-1")
		result2, _ := mockRepo.GetTransaction(ctx, "txn-2")

		assert.NotEqual(t, result1.ID, result2.ID)
		assert.NotEqual(t, result1.Username, result2.Username)
		assert.NotEqual(t, result1.TableName, result2.TableName)
	})

	t.Run("UC-S5-20: Commit Transaction with Multiple Operations", func(t *testing.T) {
		txnID := "txn-123"

		mockRepo.EXPECT().CommitTransaction(ctx, txnID).Return(nil)

		err := mockRepo.CommitTransaction(ctx, txnID)

		require.NoError(t, err)
	})

	t.Run("UC-S5-21: Rollback Transaction with Multiple Operations", func(t *testing.T) {
		txnID := "txn-123"

		mockRepo.EXPECT().RollbackTransaction(ctx, txnID).Return(nil)

		err := mockRepo.RollbackTransaction(ctx, txnID)

		require.NoError(t, err)
	})

	t.Run("UC-S5-22: Get Non-Existent Transaction", func(t *testing.T) {
		txnID := "nonexistent"

		mockRepo.EXPECT().GetTransaction(ctx, txnID).Return(nil, assert.AnError)

		txn, err := mockRepo.GetTransaction(ctx, txnID)

		require.Error(t, err)
		assert.Nil(t, txn)
	})

	t.Run("UC-S5-23: Transaction Expiration Check", func(t *testing.T) {
		txnID := "txn-expired"
		now := time.Now()
		expiredTxn := &domain.Transaction{
			ID:        txnID,
			Username:  "testuser",
			TableName: "users",
			StartedAt: now.Add(-10 * time.Minute),
			ExpiresAt: now.Add(-1 * time.Minute),
		}

		mockRepo.EXPECT().GetTransaction(ctx, txnID).Return(expiredTxn, nil)

		txn, err := mockRepo.GetTransaction(ctx, txnID)

		require.NoError(t, err)
		assert.True(t, txn.ExpiresAt.Before(time.Now()))
	})

	t.Run("UC-S5-24: Buffer INSERT Operation with All Fields", func(t *testing.T) {
		txnID := "txn-123"
		op := domain.TransactionOperation{
			Type:      domain.OperationInsert,
			TableName: "users",
			NewValue: map[string]interface{}{
				"id":       5,
				"username": "alice",
				"email":    "alice@test.com",
				"active":   true,
			},
		}

		mockRepo.EXPECT().BufferOperation(ctx, txnID, op).Return(nil)

		err := mockRepo.BufferOperation(ctx, txnID, op)

		require.NoError(t, err)
	})

	t.Run("UC-S5-25: Buffer DELETE Operation by Primary Key", func(t *testing.T) {
		txnID := "txn-123"
		op := domain.TransactionOperation{
			Type:       domain.OperationDelete,
			TableName:  "users",
			PrimaryKey: 5,
		}

		mockRepo.EXPECT().BufferOperation(ctx, txnID, op).Return(nil)

		err := mockRepo.BufferOperation(ctx, txnID, op)

		require.NoError(t, err)
	})
}
