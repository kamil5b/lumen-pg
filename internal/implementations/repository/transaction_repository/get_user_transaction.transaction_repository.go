package transaction_repository

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (t *TransactionRepositoryImplementation) GetUserTransaction(ctx context.Context, username string) (*domain.TransactionState, error) {
	return nil, errors.New("not implemented")
}
