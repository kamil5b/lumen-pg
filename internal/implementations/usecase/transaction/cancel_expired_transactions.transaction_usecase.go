package transaction

import (
	"context"
)

func (u *TransactionUseCaseImplementation) CancelExpiredTransactions(ctx context.Context) error {
	return u.transactionRepo.InvalidateExpiredTransactions(ctx)
}
