package transaction

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// GetTransaction retrieves an active transaction
func (t *TransactionRepository) GetTransaction(ctx context.Context, txnID string) (*domain.Transaction, error) {
	return nil, errors.New("not implemented yet")
}
