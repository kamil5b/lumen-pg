package transaction

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *TransactionUseCaseImplementation) CommitTransaction(ctx context.Context, username string) error {
	// Get the active transaction for the user
	txn, err := u.transactionRepo.GetUserTransaction(ctx, username)
	if err != nil {
		return err
	}

	if txn == nil {
		return domain.ErrNoActiveTransaction
	}

	// Get all buffered operations
	edits, err := u.transactionRepo.GetRowEdits(ctx, username)
	if err != nil {
		return err
	}

	inserts, err := u.transactionRepo.GetRowInserts(ctx, username)
	if err != nil {
		return err
	}

	deletes, err := u.transactionRepo.GetRowDeletes(ctx, username)
	if err != nil {
		return err
	}

	// Check permissions for updates
	if len(edits) > 0 {
		hasPermission, err := u.rbacRepo.HasUpdatePermission(ctx, username, "", "", "")
		if err != nil {
			return err
		}
		if !hasPermission {
			return domain.ErrInsufficientPermissions
		}
	}

	// Check permissions for inserts
	if len(inserts) > 0 {
		hasPermission, err := u.rbacRepo.HasInsertPermission(ctx, username, "", "", "")
		if err != nil {
			return err
		}
		if !hasPermission {
			return domain.ErrInsufficientPermissions
		}
	}

	// Check permissions for deletes
	if len(deletes) > 0 {
		hasPermission, err := u.rbacRepo.HasDeletePermission(ctx, username, "", "", "")
		if err != nil {
			return err
		}
		if !hasPermission {
			return domain.ErrInsufficientPermissions
		}
	}

	// Update the transaction to mark it as committed
	return u.transactionRepo.UpdateTransaction(ctx, txn)
}
