package cache_repository

import (
	"context"
	"errors"
)

func (c *CacheRepositoryImplementation) GetAndDelete(ctx context.Context, key string) (interface{}, error) {
	return nil, errors.New("not implemented")
}
