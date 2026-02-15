package cache_repository

import (
	"context"
	"errors"
)

func (c *CacheRepositoryImplementation) Exists(ctx context.Context, key string) (bool, error) {
	return false, errors.New("not implemented")
}
