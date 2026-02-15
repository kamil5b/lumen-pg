package transaction

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *TransactionUseCaseImplementation) GetTransactionInserts(ctx context.Context, username string) ([]domain.RowInsert, error) {
	return nil, errors.New("not implemented")
}
