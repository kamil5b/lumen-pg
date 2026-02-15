package transaction

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *TransactionUseCaseImplementation) StartTransaction(ctx context.Context, username, database, schema, table string) (*domain.TransactionState, error) {
	return nil, errors.New("not implemented")
}
