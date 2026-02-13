package transaction

import (
	"context"
	"errors"
)

// RollbackTransaction rolls back a transaction
func (t *TransactionUseCase) RollbackTransaction(ctx context.Context, txnID string) error {
	return errors.New("not implemented yet")
}
