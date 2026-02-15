package database_repository

import (
	"context"
	"errors"
)

func (d *DatabaseRepositoryImplementation) UpdateRow(ctx context.Context, database, schema, table string, pkValues map[string]interface{}, values map[string]interface{}) error {
	return errors.New("not implemented")
}
