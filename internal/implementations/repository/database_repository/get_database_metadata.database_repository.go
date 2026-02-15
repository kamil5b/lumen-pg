package database_repository

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (d *DatabaseRepositoryImplementation) GetDatabaseMetadata(ctx context.Context, database string) (*domain.DatabaseMetadata, error) {
	return nil, errors.New("not implemented")
}
