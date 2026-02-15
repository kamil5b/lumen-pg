package transaction_repository

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (t *TransactionRepositoryImplementation) GetTransaction(ctx context.Context, transactionID string) (*domain.TransactionState, error) {
	return nil, errors.New("not implemented")
}
