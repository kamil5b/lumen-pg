package query

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// ExecuteMultipleQueries executes multiple queries
func (q *QueryUseCase) ExecuteMultipleQueries(ctx context.Context, queries string) ([]*domain.QueryResult, error) {
	return nil, errors.New("not implemented yet")
}
