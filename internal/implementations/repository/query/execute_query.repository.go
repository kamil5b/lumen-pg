package query

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// ExecuteQuery executes a SQL query with parameters
func (q *QueryRepository) ExecuteQuery(ctx context.Context, query string, params ...interface{}) (*domain.QueryResult, error) {
	return nil, errors.New("not implemented yet")
}
