package transaction_repository

import (
	"context"
	"errors"
)

func (t *TransactionRepositoryImplementation) ClearRowEdits(ctx context.Context, transactionID string) error {
	return errors.New("not implemented")
}
