package cache_repository

import (
	"context"
	"errors"
)

func (c *CacheRepositoryImplementation) Delete(ctx context.Context, key string) error {
	return errors.New("not implemented")
}
