package database_repository

import (
	"context"
	"database/sql"
	"errors"
)

func (d *DatabaseRepositoryImplementation) RollbackTransaction(ctx context.Context, tx *sql.Tx) error {
	if tx == nil {
		return errors.New("transaction is nil")
	}
	return tx.Rollback()
}
