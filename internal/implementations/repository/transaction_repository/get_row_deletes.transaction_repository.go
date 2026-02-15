package transaction_repository

import (
	"context"
	"errors"
)

func (t *TransactionRepositoryImplementation) GetRowDeletes(ctx context.Context, transactionID string) ([]int, error) {
	return nil, errors.New("not implemented")
}
