package database_repository

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (d *DatabaseRepositoryImplementation) GetTableMetadata(ctx context.Context, database, schema, table string) (*domain.TableMetadata, error) {
	return nil, errors.New("not implemented")
}
