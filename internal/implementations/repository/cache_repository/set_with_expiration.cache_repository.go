package cache_repository

import (
	"context"
	"errors"
)

func (c *CacheRepositoryImplementation) SetWithExpiration(ctx context.Context, key string, value interface{}, expirationTime int64) error {
	return errors.New("not implemented")
}
