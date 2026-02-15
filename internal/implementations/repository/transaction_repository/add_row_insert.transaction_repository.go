package transaction_repository

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (t *TransactionRepositoryImplementation) AddRowInsert(ctx context.Context, transactionID string, insert domain.RowInsert) error {
	return errors.New("not implemented")
}
