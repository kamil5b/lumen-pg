package transaction

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// BufferEdit buffers an edit operation in a transaction
func (t *TransactionUseCase) BufferEdit(ctx context.Context, txnID string, op domain.TransactionOperation) error {
	return errors.New("not implemented yet")
}
