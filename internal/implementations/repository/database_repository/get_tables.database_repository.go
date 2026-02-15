package database_repository

import (
	"context"
	"errors"
)

func (d *DatabaseRepositoryImplementation) GetTables(ctx context.Context, database, schema string) ([]string, error) {
	return nil, errors.New("not implemented")
}
