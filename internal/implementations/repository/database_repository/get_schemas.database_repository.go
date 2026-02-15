package database_repository

import (
	"context"
	"errors"
)

func (d *DatabaseRepositoryImplementation) GetSchemas(ctx context.Context, database string) ([]string, error) {
	return nil, errors.New("not implemented")
}
