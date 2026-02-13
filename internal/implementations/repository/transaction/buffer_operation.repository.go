package transaction

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// BufferOperation adds an operation to the transaction buffer
func (t *TransactionRepository) BufferOperation(ctx context.Context, txnID string, op domain.TransactionOperation) error {
	return errors.New("not implemented yet")
}
