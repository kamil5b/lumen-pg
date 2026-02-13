package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
	mockRepository "github.com/kamil5b/lumen-pg/internal/testrunners/mocks/repository"
	"github.com/stretchr/testify/require"
)

// TransactionUsecaseConstructor is a function type that creates a TransactionUseCase
type TransactionUsecaseConstructor func(
	transactionRepo repository.TransactionRepository,
	databaseRepo repository.DatabaseRepository,
	rbacRepo repository.RBACRepository,
) usecase.TransactionUseCase

// TransactionUsecaseRunner runs all transaction usecase tests against an implementation
// Maps to TEST_PLAN.md:
// - Story 5: Main View & Data Interaction [UC-S5-09~16, IT-S5-04~05, E2E-S5-06~13]
// - Story 6: Isolation [UC-S6-02, IT-S6-03, E2E-S6-02]
func TransactionUsecaseRunner(t *testing.T, constructor TransactionUsecaseConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransaction := mockRepository.NewMockTransactionRepository(ctrl)
	mockDatabase := mockRepository.NewMockDatabaseRepository(ctrl)
	mockRBAC := mockRepository.NewMockRBACRepository(ctrl)

	uc := constructor(mockTransaction, mockDatabase, mockRBAC)

	ctx := context.Background()

	// UC-S5-09: Transaction Start
	// E2E-S5-06: Start Transaction Button
	t.Run("StartTransaction creates new transaction", func(t *testing.T) {
		mockTransaction.EXPECT().
			CreateTransaction(gomock.Any(), gomock.Any()).
			Return(&domain.TransactionState{
				ID:       "txn_123",
				Username: "testuser",
				Edits:    make(map[int]domain.RowEdit),
				Deletes:  []int{},
				Inserts:  []domain.RowInsert{},
			}, nil)

		txn, err := uc.StartTransaction(ctx, "testuser", "testdb", "public", "users")

		require.NoError(t, err)
		require.NotNil(t, txn)
		require.Equal(t, "testuser", txn.Username)
		require.NotEmpty(t, txn.ID)
	})

	// UC-S5-10: Transaction Already Active Error
	t.Run("StartTransaction returns error when transaction already active", func(t *testing.T) {
		mockTransaction.EXPECT().
			GetActiveTransaction(gomock.Any(), "testuser").
			Return(&domain.TransactionState{
				ID:       "txn_existing",
				Username: "testuser",
			}, nil)

		_, err := uc.StartTransaction(ctx, "testuser", "testdb", "public", "users")

		require.Error(t, err)
	})

	// UC-S5-09: Transaction Start (check if active)
	t.Run("CheckActiveTransaction returns true when transaction exists", func(t *testing.T) {
		mockTransaction.EXPECT().
			GetActiveTransaction(gomock.Any(), "testuser").
			Return(&domain.TransactionState{
				ID:       "txn_123",
				Username: "testuser",
			}, nil)

		active, err := uc.CheckActiveTransaction(ctx, "testuser")

		require.NoError(t, err)
		require.True(t, active)
	})

	t.Run("CheckActiveTransaction returns false when no transaction", func(t *testing.T) {
		mockTransaction.EXPECT().
			GetActiveTransaction(gomock.Any(), "testuser").
			Return(nil, ErrNoActiveTransaction)

		active, err := uc.CheckActiveTransaction(ctx, "testuser")

		require.NoError(t, err)
		require.False(t, active)
	})

	t.Run("GetActiveTransaction returns current transaction", func(t *testing.T) {
		mockTransaction.EXPECT().
			GetActiveTransaction(gomock.Any(), "testuser").
			Return(&domain.TransactionState{
				ID:       "txn_123",
				Username: "testuser",
				Edits:    make(map[int]domain.RowEdit),
				Deletes:  []int{},
				Inserts:  []domain.RowInsert{},
			}, nil)

		txn, err := uc.GetActiveTransaction(ctx, "testuser")

		require.NoError(t, err)
		require.NotNil(t, txn)
		require.Equal(t, "txn_123", txn.ID)
	})

	// UC-S5-11: Cell Edit Buffering
	// E2E-S5-07: Transaction Mode Cell Editing
	// E2E-S5-08: Transaction Mode Edit Buffer Display
	t.Run("EditCell buffers a cell edit", func(t *testing.T) {
		mockTransaction.EXPECT().
			UpdateTransactionEdit(gomock.Any(), "testuser", gomock.Any()).
			Return(nil)

		err := uc.EditCell(ctx, "testuser", "testdb", "public", "users", 5, "name", "NewName")

		require.NoError(t, err)
	})

	t.Run("GetTransactionEdits returns all buffered edits", func(t *testing.T) {
		edits := map[int]domain.RowEdit{
			0: {
				RowIndex:   0,
				ColumnName: "name",
				OldValue:   "OldName",
				NewValue:   "NewName",
			},
			1: {
				RowIndex:   1,
				ColumnName: "email",
				OldValue:   "old@example.com",
				NewValue:   "new@example.com",
			},
		}

		mockTransaction.EXPECT().
			GetTransactionEdits(gomock.Any(), "testuser").
			Return(edits, nil)

		result, err := uc.GetTransactionEdits(ctx, "testuser")

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, 2, len(result))
	})

	// UC-S5-15: Row Deletion Buffering
	// E2E-S5-12: Transaction Row Delete Button
	t.Run("DeleteRow buffers a row deletion", func(t *testing.T) {
		mockTransaction.EXPECT().
			UpdateTransactionDelete(gomock.Any(), "testuser", gomock.Any()).
			Return(nil)

		err := uc.DeleteRow(ctx, "testuser", "testdb", "public", "users", 3)

		require.NoError(t, err)
	})

	t.Run("GetTransactionDeletes returns all buffered deletions", func(t *testing.T) {
		deletes := []int{0, 2, 5}

		mockTransaction.EXPECT().
			GetTransactionDeletes(gomock.Any(), "testuser").
			Return(deletes, nil)

		result, err := uc.GetTransactionDeletes(ctx, "testuser")

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, 3, len(result))
	})

	// UC-S5-16: Row Insertion Buffering
	// E2E-S5-13: Transaction New Row Button
	t.Run("InsertRow buffers a new row insertion", func(t *testing.T) {
		mockTransaction.EXPECT().
			UpdateTransactionInsert(gomock.Any(), "testuser", gomock.Any()).
			Return(nil)

		values := map[string]interface{}{
			"name":  "NewUser",
			"email": "new@example.com",
		}

		err := uc.InsertRow(ctx, "testuser", "testdb", "public", "users", values)

		require.NoError(t, err)
	})

	t.Run("GetTransactionInserts returns all buffered insertions", func(t *testing.T) {
		inserts := []domain.RowInsert{
			{
				Values: map[string]interface{}{
					"name":  "User1",
					"email": "user1@example.com",
				},
			},
			{
				Values: map[string]interface{}{
					"name":  "User2",
					"email": "user2@example.com",
				},
			},
		}

		mockTransaction.EXPECT().
			GetTransactionInserts(gomock.Any(), "testuser").
			Return(inserts, nil)

		result, err := uc.GetTransactionInserts(ctx, "testuser")

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, 2, len(result))
	})

	// UC-S5-12: Transaction Commit
	// IT-S5-04: Real Transaction Commit
	// E2E-S5-09: Transaction Commit Button
	t.Run("CommitTransaction commits all buffered changes", func(t *testing.T) {
		mockTransaction.EXPECT().
			CommitTransaction(gomock.Any(), "testuser").
			Return(nil)

		mockRBAC.EXPECT().
			CheckUpdatePermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(true, nil).AnyTimes()

		mockRBAC.EXPECT().
			CheckInsertPermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(true, nil).AnyTimes()

		mockRBAC.EXPECT().
			CheckDeletePermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(true, nil).AnyTimes()

		err := uc.CommitTransaction(ctx, "testuser")

		require.NoError(t, err)
	})

	// UC-S5-13: Transaction Rollback
	// IT-S5-05: Real Transaction Rollback
	// E2E-S5-10: Transaction Rollback Button
	t.Run("RollbackTransaction rolls back all buffered changes", func(t *testing.T) {
		mockTransaction.EXPECT().
			RollbackTransaction(gomock.Any(), "testuser").
			Return(nil)

		err := uc.RollbackTransaction(ctx, "testuser")

		require.NoError(t, err)
	})

	// UC-S5-14: Transaction Timer Expiration
	// E2E-S5-11: Transaction Timer Countdown
	t.Run("GetTransactionRemainingTime returns remaining time", func(t *testing.T) {
		mockTransaction.EXPECT().
			GetTransactionRemainingTime(gomock.Any(), "testuser").
			Return(int64(45), nil)

		remaining, err := uc.GetTransactionRemainingTime(ctx, "testuser")

		require.NoError(t, err)
		require.Equal(t, int64(45), remaining)
	})

	t.Run("IsTransactionExpired returns false when transaction active", func(t *testing.T) {
		mockTransaction.EXPECT().
			GetTransactionRemainingTime(gomock.Any(), "testuser").
			Return(int64(30), nil)

		expired, err := uc.IsTransactionExpired(ctx, "testuser")

		require.NoError(t, err)
		require.False(t, expired)
	})

	t.Run("IsTransactionExpired returns true when time exceeded", func(t *testing.T) {
		mockTransaction.EXPECT().
			GetTransactionRemainingTime(gomock.Any(), "testuser").
			Return(int64(0), nil)

		expired, err := uc.IsTransactionExpired(ctx, "testuser")

		require.NoError(t, err)
		require.True(t, expired)
	})

	t.Run("CancelExpiredTransactions cancels all expired transactions", func(t *testing.T) {
		mockTransaction.EXPECT().
			GetAllActiveTransactions(gomock.Any()).
			Return([]*domain.TransactionState{
				{
					ID:       "txn_1",
					Username: "user1",
				},
				{
					ID:       "txn_2",
					Username: "user2",
				},
			}, nil)

		mockTransaction.EXPECT().
			GetTransactionRemainingTime(gomock.Any(), gomock.Any()).
			Return(int64(-10), nil).AnyTimes()

		mockTransaction.EXPECT().
			RollbackTransaction(gomock.Any(), gomock.Any()).
			Return(nil).AnyTimes()

		err := uc.CancelExpiredTransactions(ctx)

		require.NoError(t, err)
	})

	// UC-S6-02: Transaction Isolation
	// IT-S6-03: Real Transaction Isolation
	// E2E-S6-02: Simultaneous Transactions
	t.Run("Multiple users can have independent transactions", func(t *testing.T) {
		mockTransaction.EXPECT().
			CreateTransaction(gomock.Any(), gomock.Any()).
			Return(&domain.TransactionState{
				ID:       "txn_user1",
				Username: "user1",
				Edits:    make(map[int]domain.RowEdit),
			}, nil).Times(1)

		mockTransaction.EXPECT().
			CreateTransaction(gomock.Any(), gomock.Any()).
			Return(&domain.TransactionState{
				ID:       "txn_user2",
				Username: "user2",
				Edits:    make(map[int]domain.RowEdit),
			}, nil).Times(1)

		txn1, err1 := uc.StartTransaction(ctx, "user1", "testdb", "public", "users")
		txn2, err2 := uc.StartTransaction(ctx, "user2", "testdb", "public", "users")

		require.NoError(t, err1)
		require.NoError(t, err2)
		require.NotEqual(t, txn1.ID, txn2.ID)
		require.NotEqual(t, txn1.Username, txn2.Username)
	})

	t.Run("Transaction changes are isolated per user", func(t *testing.T) {
		edits1 := map[int]domain.RowEdit{
			0: {RowIndex: 0, ColumnName: "name", NewValue: "User1"},
		}

		edits2 := map[int]domain.RowEdit{
			1: {RowIndex: 1, ColumnName: "name", NewValue: "User2"},
		}

		mockTransaction.EXPECT().
			GetTransactionEdits(gomock.Any(), "user1").
			Return(edits1, nil)

		mockTransaction.EXPECT().
			GetTransactionEdits(gomock.Any(), "user2").
			Return(edits2, nil)

		result1, err1 := uc.GetTransactionEdits(ctx, "user1")
		result2, err2 := uc.GetTransactionEdits(ctx, "user2")

		require.NoError(t, err1)
		require.NoError(t, err2)
		require.NotEqual(t, result1, result2)
	})
}

var (
	ErrNoActiveTransaction = domain.ValidationError{Field: "transaction", Message: "no active transaction"}
)
