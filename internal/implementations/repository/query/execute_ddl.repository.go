package query

import (
	"context"
	"errors"
)

// ExecuteDDL executes a DDL query
func (q *QueryRepository) ExecuteDDL(ctx context.Context, query string) error {
	return errors.New("not implemented yet")
}
