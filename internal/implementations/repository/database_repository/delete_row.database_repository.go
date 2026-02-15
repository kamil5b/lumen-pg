package database_repository

import (
	"context"
	"errors"
)

func (d *DatabaseRepositoryImplementation) DeleteRow(ctx context.Context, database, schema, table string, pkValues map[string]interface{}) error {
	return errors.New("not implemented")
}
