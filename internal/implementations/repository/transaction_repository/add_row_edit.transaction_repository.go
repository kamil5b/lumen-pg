package transaction_repository

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (t *TransactionRepositoryImplementation) AddRowEdit(ctx context.Context, transactionID string, edit domain.RowEdit) error {
	return errors.New("not implemented")
}
