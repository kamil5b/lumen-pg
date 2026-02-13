package transaction

import (
	"context"
	"errors"
)

// CommitTransaction commits a transaction
func (t *TransactionUseCase) CommitTransaction(ctx context.Context, txnID string) error {
	return errors.New("not implemented yet")
}
