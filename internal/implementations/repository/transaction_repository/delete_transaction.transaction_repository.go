package transaction_repository

import (
	"context"
	"errors"
)

func (t *TransactionRepositoryImplementation) DeleteTransaction(ctx context.Context, transactionID string) error {
	return errors.New("not implemented")
}
