package transaction

import (
	"context"
	"time"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *TransactionUseCaseImplementation) GetTransactionRemainingTime(ctx context.Context, username string) (int64, error) {
	// Get the active transaction for the user
	txn, err := u.transactionRepo.GetUserTransaction(ctx, username)
	if err != nil {
		return 0, err
	}

	if txn == nil {
		return 0, domain.ErrNoActiveTransaction
	}

	// Calculate remaining time in seconds
	now := time.Now()
	remaining := txn.ExpiresAt.Sub(now).Seconds()

	// If transaction has already expired, return 0
	if remaining < 0 {
		return 0, nil
	}

	return int64(remaining), nil
}
