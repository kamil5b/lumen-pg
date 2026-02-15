package transaction

import (
	"context"
	"errors"
)

func (u *TransactionUseCaseImplementation) RollbackTransaction(ctx context.Context, username string) error {
	return errors.New("not implemented")
}
