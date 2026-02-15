package database_repository

import (
	"context"
	"errors"
)

func (d *DatabaseRepositoryImplementation) InsertRow(ctx context.Context, database, schema, table string, values map[string]interface{}) error {
	return errors.New("not implemented")
}
