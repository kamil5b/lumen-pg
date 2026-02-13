package query

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// ExecuteMultiple executes multiple queries separated by semicolons
func (q *QueryRepository) ExecuteMultiple(ctx context.Context, queries string) ([]*domain.QueryResult, error) {
	return nil, errors.New("not implemented yet")
}
