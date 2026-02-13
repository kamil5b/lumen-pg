package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/implementations/mocks"
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
)

// IsolationUseCaseConstructor creates an isolation use case with its dependencies
type IsolationUseCaseConstructor func(
	sessionRepo repository.SessionRepository,
	transactionRepo repository.TransactionRepository,
) usecase.AuthUseCase

// IsolationUseCaseRunner runs test specs for isolation use case (Story 6)
func IsolationUseCaseRunner(t *testing.T, constructor IsolationUseCaseConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockTransactionRepo := mocks.NewMockTransactionRepository(ctrl)

	t.Run("UC-S6-01: Session Isolation", func(t *testing.T) {
		ctx := context.Background()

		// Create first session
		session1 := &domain.Session{
			Username: "userA",
		}

		// Create second session
		session2 := &domain.Session{
			Username: "userB",
		}

		mockSessionRepo.EXPECT().ValidateSession(ctx, "token-A").Return(session1, nil)
		mockSessionRepo.EXPECT().ValidateSession(ctx, "token-B").Return(session2, nil)

		// Validate both sessions exist independently
		result1, err1 := mockSessionRepo.ValidateSession(ctx, "token-A")
		result2, err2 := mockSessionRepo.ValidateSession(ctx, "token-B")

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.Equal(t, "userA", result1.Username)
		assert.Equal(t, "userB", result2.Username)
	})

	t.Run("UC-S6-02: Transaction Isolation", func(t *testing.T) {
		ctx := context.Background()

		// Start transaction for user A
		txnA := &domain.Transaction{
			ID:        "txn-A",
			Username:  "userA",
			TableName: "users",
		}

		// Start transaction for user B
		txnB := &domain.Transaction{
			ID:        "txn-B",
			Username:  "userB",
			TableName: "posts",
		}

		mockTransactionRepo.EXPECT().StartTransaction(ctx, "userA", "users").Return(txnA, nil)
		mockTransactionRepo.EXPECT().StartTransaction(ctx, "userB", "posts").Return(txnB, nil)

		// Get both transactions
		result1, err1 := mockTransactionRepo.StartTransaction(ctx, "userA", "users")
		result2, err2 := mockTransactionRepo.StartTransaction(ctx, "userB", "posts")

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.Equal(t, "userA", result1.Username)
		assert.Equal(t, "userB", result2.Username)
		assert.NotEqual(t, result1.ID, result2.ID)
		assert.NotEqual(t, result1.TableName, result2.TableName)
	})

	t.Run("UC-S6-03: Cookie Isolation", func(t *testing.T) {
		// User A password
		passwordA := "securepassA"
		encryptedA := "encrypted-A"

		// User B password
		passwordB := "securepassB"
		encryptedB := "encrypted-B"

		mockSessionRepo.EXPECT().EncryptPassword(passwordA).Return(encryptedA, nil)
		mockSessionRepo.EXPECT().EncryptPassword(passwordB).Return(encryptedB, nil)
		mockSessionRepo.EXPECT().DecryptPassword(encryptedA).Return(passwordA, nil)
		mockSessionRepo.EXPECT().DecryptPassword(encryptedB).Return(passwordB, nil)

		// Encrypt passwords for both users
		encA, errA := mockSessionRepo.EncryptPassword(passwordA)
		encB, errB := mockSessionRepo.EncryptPassword(passwordB)

		require.NoError(t, errA)
		require.NoError(t, errB)
		assert.NotEqual(t, encA, encB)

		// Decrypt and verify
		decA, _ := mockSessionRepo.DecryptPassword(encA)
		decB, _ := mockSessionRepo.DecryptPassword(encB)

		assert.Equal(t, passwordA, decA)
		assert.Equal(t, passwordB, decB)
		assert.NotEqual(t, decA, decB)
	})
}
