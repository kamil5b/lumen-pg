package transaction_repository

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (t *TransactionRepositoryImplementation) GetRowInserts(ctx context.Context, transactionID string) ([]domain.RowInsert, error) {
	return nil, errors.New("not implemented")
}
