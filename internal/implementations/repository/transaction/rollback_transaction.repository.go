package transaction

import (
	"context"
	"errors"
)

// RollbackTransaction rolls back all buffered operations
func (t *TransactionRepository) RollbackTransaction(ctx context.Context, txnID string) error {
	return errors.New("not implemented yet")
}
