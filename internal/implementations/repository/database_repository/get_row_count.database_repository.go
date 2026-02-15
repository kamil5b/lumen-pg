package database_repository

import (
	"context"
	"errors"
)

func (d *DatabaseRepositoryImplementation) GetRowCount(ctx context.Context, database, schema, table, whereClause string) (int64, error) {
	return 0, errors.New("not implemented")
}
