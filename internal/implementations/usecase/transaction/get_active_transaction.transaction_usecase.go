package transaction

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *TransactionUseCaseImplementation) GetActiveTransaction(ctx context.Context, username string) (*domain.TransactionState, error) {
	return u.transactionRepo.GetUserTransaction(ctx, username)
}
