package transaction

import (
	"context"
	"errors"
)

func (u *TransactionUseCaseImplementation) GetTransactionDeletes(ctx context.Context, username string) ([]int, error) {
	return nil, errors.New("not implemented")
}
