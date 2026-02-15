package database_repository

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (d *DatabaseRepositoryImplementation) ExecuteMultipleQueries(ctx context.Context, queries string) ([]domain.QueryResult, error) {
	return nil, errors.New("not implemented")
}
