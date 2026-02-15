package transaction_repository

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (t *TransactionRepositoryImplementation) GetRowEdits(ctx context.Context, transactionID string) (map[int]domain.RowEdit, error) {
	return nil, errors.New("not implemented")
}
