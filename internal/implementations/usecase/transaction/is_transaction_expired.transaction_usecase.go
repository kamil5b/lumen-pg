package transaction

import (
	"context"
	"time"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *TransactionUseCaseImplementation) IsTransactionExpired(ctx context.Context, username string) (bool, error) {
	// Get the active transaction for the user
	txn, err := u.transactionRepo.GetUserTransaction(ctx, username)
	if err != nil {
		// If there's no active transaction, consider it expired
		if validationErr, ok := err.(domain.ValidationError); ok {
			if validationErr.Field == "transaction" && validationErr.Message == "no active transaction" {
				return true, nil
			}
		}
		if appErr, ok := err.(*domain.ApplicationError); ok {
			if appErr.Type == domain.ErrTypeTransaction && appErr.Message == "no active transaction" {
				return true, nil
			}
		}
		return false, err
	}

	if txn == nil {
		return true, nil
	}

	// Check if the transaction has expired
	now := time.Now()
	return now.After(txn.ExpiresAt), nil
}
