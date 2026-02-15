package transaction

import (
	"context"
	"errors"
)

func (u *TransactionUseCaseImplementation) IsTransactionExpired(ctx context.Context, username string) (bool, error) {
	return false, errors.New("not implemented")
}
