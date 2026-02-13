package repository

import (
	"context"
)

// CacheRepository defines operations for caching
type CacheRepository interface {
	// Set stores a value in the cache with a key
	Set(ctx context.Context, key string, value interface{}, ttlSeconds int) error

	// Get retrieves a value from the cache by key
	Get(ctx context.Context, key string) (interface{}, error)

	// Delete removes a value from the cache
	Delete(ctx context.Context, key string) error

	// Exists checks if a key exists in the cache
	Exists(ctx context.Context, key string) (bool, error)

	// Clear clears all values from the cache
	Clear(ctx context.Context) error

	// SetWithExpiration sets a value with an expiration time
	SetWithExpiration(ctx context.Context, key string, value interface{}, expirationTime int64) error

	// GetAndDelete retrieves a value and deletes it atomically
	GetAndDelete(ctx context.Context, key string) (interface{}, error)
}
