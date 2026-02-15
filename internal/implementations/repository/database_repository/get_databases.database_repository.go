package database_repository

import (
	"context"
	"errors"
)

func (d *DatabaseRepositoryImplementation) GetDatabases(ctx context.Context) ([]string, error) {
	return nil, errors.New("not implemented")
}
