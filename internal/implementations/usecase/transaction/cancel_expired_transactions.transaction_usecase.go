package transaction

import (
	"context"
	"errors"
)

func (u *TransactionUseCaseImplementation) CancelExpiredTransactions(ctx context.Context) error {
	return errors.New("not implemented")
}
