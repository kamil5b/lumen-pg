package transaction

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *TransactionUseCaseImplementation) GetTransactionEdits(ctx context.Context, username string) (map[int]domain.RowEdit, error) {
	return nil, errors.New("not implemented")
}
