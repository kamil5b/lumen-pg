package transaction

import (
	"context"
	"errors"
)

func (u *TransactionUseCaseImplementation) GetTransactionRemainingTime(ctx context.Context, username string) (int64, error) {
	return 0, errors.New("not implemented")
}
