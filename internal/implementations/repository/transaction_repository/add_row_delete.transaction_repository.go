package transaction_repository

import (
	"context"
	"errors"
)

func (t *TransactionRepositoryImplementation) AddRowDelete(ctx context.Context, transactionID string, rowIndex int) error {
	return errors.New("not implemented")
}
