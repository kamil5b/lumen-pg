package query

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *QueryUseCaseImplementation) ExecuteQuery(ctx context.Context, username, query string, offset, limit int) (*domain.QueryResult, error) {
	return nil, errors.New("not implemented")
}
