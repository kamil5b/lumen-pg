package query

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// ExecuteQuery performs query execution
func (q *QueryUseCase) ExecuteQuery(ctx context.Context, sql string, params ...interface{}) (*domain.QueryResult, error) {
	return nil, errors.New("not implemented yet")
}
