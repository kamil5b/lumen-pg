package transaction

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *TransactionUseCaseImplementation) GetActiveTransaction(ctx context.Context, username string) (*domain.TransactionState, error) {
	return nil, errors.New("not implemented")
}
