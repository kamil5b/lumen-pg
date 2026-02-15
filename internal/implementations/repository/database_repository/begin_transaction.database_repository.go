package database_repository

import (
	"context"
	"database/sql"
	"errors"
)

func (d *DatabaseRepositoryImplementation) BeginTransaction(ctx context.Context) (*sql.Tx, error) {
	return nil, errors.New("not implemented")
}
