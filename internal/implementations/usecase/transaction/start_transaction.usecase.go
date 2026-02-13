package transaction

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// StartTransaction starts a new transaction
func (t *TransactionUseCase) StartTransaction(ctx context.Context, username string, tableName string) (*domain.Transaction, error) {
	return nil, errors.New("not implemented yet")
}
