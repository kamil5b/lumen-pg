package database_repository

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (d *DatabaseRepositoryImplementation) ExecuteQueryWithPagination(ctx context.Context, params domain.QueryParams) (*domain.QueryResult, error) {
	return nil, errors.New("not implemented")
}
