package cache_repository

import (
	"context"
	"errors"
)

func (c *CacheRepositoryImplementation) Set(ctx context.Context, key string, value interface{}, ttlSeconds int) error {
	return errors.New("not implemented")
}
