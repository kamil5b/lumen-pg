package transaction_repository

import (
	"context"
	"errors"
)

func (t *TransactionRepositoryImplementation) TransactionExists(ctx context.Context, transactionID string) (bool, error) {
	return false, errors.New("not implemented")
}
