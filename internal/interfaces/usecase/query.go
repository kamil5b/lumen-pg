package usecase

import (
	"context"
	"github.com/kamil5b/lumen-pg/internal/domain"
)

// QueryUseCase handles query execution operations
type QueryUseCase interface {
	ExecuteQuery(ctx context.Context, sql string, params ...interface{}) (*domain.QueryResult, error)
	ExecuteMultipleQueries(ctx context.Context, queries string) ([]*domain.QueryResult, error)
}
