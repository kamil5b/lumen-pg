package cache_repository

import (
	"context"
	"errors"
)

func (c *CacheRepositoryImplementation) Get(ctx context.Context, key string) (interface{}, error) {
	return nil, errors.New("not implemented")
}
