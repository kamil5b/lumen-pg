package transaction_repository

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (t *TransactionRepositoryImplementation) CreateTransaction(ctx context.Context, transaction *domain.TransactionState) error {
	return errors.New("not implemented")
}
