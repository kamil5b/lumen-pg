package query

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// ExecuteDML executes a DML query (INSERT, UPDATE, DELETE)
func (q *QueryRepository) ExecuteDML(ctx context.Context, query string, params ...interface{}) (*domain.QueryResult, error) {
	return nil, errors.New("not implemented yet")
}
