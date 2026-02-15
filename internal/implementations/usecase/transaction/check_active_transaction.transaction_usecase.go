package transaction

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *TransactionUseCaseImplementation) CheckActiveTransaction(ctx context.Context, username string) (bool, error) {
	txn, err := u.transactionRepo.GetUserTransaction(ctx, username)
	if err != nil {
		// Check if it's the "no active transaction" error
		if validationErr, ok := err.(domain.ValidationError); ok {
			if validationErr.Field == "transaction" && validationErr.Message == "no active transaction" {
				return false, nil
			}
		}
		// If it's an ApplicationError
		if appErr, ok := err.(*domain.ApplicationError); ok {
			if appErr.Type == domain.ErrTypeTransaction && appErr.Message == "no active transaction" {
				return false, nil
			}
		}
		// For other errors, return the error
		return false, err
	}

	if txn == nil {
		return false, nil
	}

	return true, nil
}
