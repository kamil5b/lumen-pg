package query

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *QueryUseCaseImplementation) ExecuteMultipleQueries(ctx context.Context, username, queries string) ([]domain.QueryResult, error) {
	return nil, errors.New("not implemented")
}
