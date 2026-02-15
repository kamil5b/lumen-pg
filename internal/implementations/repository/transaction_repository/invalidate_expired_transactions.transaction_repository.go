package transaction_repository

import (
	"context"
	"errors"
)

func (t *TransactionRepositoryImplementation) InvalidateExpiredTransactions(ctx context.Context) error {
	return errors.New("not implemented")
}
