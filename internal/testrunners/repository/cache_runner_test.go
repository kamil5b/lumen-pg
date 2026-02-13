package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
)

// CacheRepositoryConstructor is a function type that creates a CacheRepository
type CacheRepositoryConstructor func() repository.CacheRepository

// CacheRepositoryRunner runs all cache repository tests against an implementation
// Maps to TEST_PLAN.md:
// - Story 1: Setup & Configuration [IT-S1-04: Cache Accessible Resources Per Role]
// - Story 7: Security & Best Practices [UC-S7-06, UC-S7-07: Session Timeout Cookies]
func CacheRepositoryRunner(t *testing.T, constructor CacheRepositoryConstructor) {
	t.Helper()

	ctx := context.Background()
	repo := constructor()

	// IT-S1-04: Cache Accessible Resources Per Role
	t.Run("Set and Get value", func(t *testing.T) {
		err := repo.Set(ctx, "key1", "value1", 3600)
		require.NoError(t, err)

		val, err := repo.Get(ctx, "key1")
		require.NoError(t, err)
		require.NotNil(t, val)
		require.Equal(t, "value1", val)
	})

	t.Run("Get non-existent key returns error", func(t *testing.T) {
		_, err := repo.Get(ctx, "nonexistent_key")
		require.Error(t, err)
	})

	t.Run("Set with zero TTL", func(t *testing.T) {
		err := repo.Set(ctx, "key_zero_ttl", "value", 0)
		require.NoError(t, err)

		val, err := repo.Get(ctx, "key_zero_ttl")
		require.NoError(t, err)
		require.Equal(t, "value", val)
	})

	t.Run("Set with large TTL", func(t *testing.T) {
		err := repo.Set(ctx, "key_large_ttl", "value", 86400*365)
		require.NoError(t, err)

		val, err := repo.Get(ctx, "key_large_ttl")
		require.NoError(t, err)
		require.Equal(t, "value", val)
	})

	t.Run("Delete removes key", func(t *testing.T) {
		err := repo.Set(ctx, "key_to_delete", "value", 3600)
		require.NoError(t, err)

		err = repo.Delete(ctx, "key_to_delete")
		require.NoError(t, err)

		_, err = repo.Get(ctx, "key_to_delete")
		require.Error(t, err)
	})

	t.Run("Delete non-existent key returns error", func(t *testing.T) {
		err := repo.Delete(ctx, "nonexistent_key_delete")
		require.Error(t, err)
	})

	t.Run("Exists returns true for existing key", func(t *testing.T) {
		err := repo.Set(ctx, "exists_key", "value", 3600)
		require.NoError(t, err)

		exists, err := repo.Exists(ctx, "exists_key")
		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("Exists returns false for non-existent key", func(t *testing.T) {
		exists, err := repo.Exists(ctx, "nonexistent_key_exists")
		require.NoError(t, err)
		require.False(t, exists)
	})

	t.Run("Clear removes all cache entries", func(t *testing.T) {
		err := repo.Set(ctx, "clear_key1", "value1", 3600)
		require.NoError(t, err)

		err = repo.Set(ctx, "clear_key2", "value2", 3600)
		require.NoError(t, err)

		err = repo.Clear(ctx)
		require.NoError(t, err)

		exists1, _ := repo.Exists(ctx, "clear_key1")
		exists2, _ := repo.Exists(ctx, "clear_key2")

		require.False(t, exists1)
		require.False(t, exists2)
	})

	t.Run("SetWithExpiration sets with Unix timestamp", func(t *testing.T) {
		expirationTime := time.Now().Add(1 * time.Hour).Unix()
		err := repo.SetWithExpiration(ctx, "expire_key", "value", expirationTime)
		require.NoError(t, err)

		val, err := repo.Get(ctx, "expire_key")
		require.NoError(t, err)
		require.Equal(t, "value", val)
	})

	t.Run("SetWithExpiration with past timestamp expires immediately", func(t *testing.T) {
		pastTime := time.Now().Add(-1 * time.Hour).Unix()
		err := repo.SetWithExpiration(ctx, "past_key", "value", pastTime)
		require.NoError(t, err)

		_, err = repo.Get(ctx, "past_key")
		require.Error(t, err)
	})

	t.Run("GetAndDelete retrieves and removes atomically", func(t *testing.T) {
		err := repo.Set(ctx, "get_delete_key", "value", 3600)
		require.NoError(t, err)

		val, err := repo.GetAndDelete(ctx, "get_delete_key")
		require.NoError(t, err)
		require.Equal(t, "value", val)

		_, err = repo.Get(ctx, "get_delete_key")
		require.Error(t, err)
	})

	t.Run("GetAndDelete non-existent key returns error", func(t *testing.T) {
		_, err := repo.GetAndDelete(ctx, "nonexistent_get_delete")
		require.Error(t, err)
	})

	t.Run("Set different value types", func(t *testing.T) {
		err := repo.Set(ctx, "string_key", "string_value", 3600)
		require.NoError(t, err)

		err = repo.Set(ctx, "int_key", 123, 3600)
		require.NoError(t, err)

		err = repo.Set(ctx, "bool_key", true, 3600)
		require.NoError(t, err)

		strVal, _ := repo.Get(ctx, "string_key")
		require.Equal(t, "string_value", strVal)

		intVal, _ := repo.Get(ctx, "int_key")
		require.Equal(t, 123, intVal)

		boolVal, _ := repo.Get(ctx, "bool_key")
		require.Equal(t, true, boolVal)
	})

	t.Run("Set overwrites existing value", func(t *testing.T) {
		err := repo.Set(ctx, "overwrite_key", "value1", 3600)
		require.NoError(t, err)

		err = repo.Set(ctx, "overwrite_key", "value2", 3600)
		require.NoError(t, err)

		val, err := repo.Get(ctx, "overwrite_key")
		require.NoError(t, err)
		require.Equal(t, "value2", val)
	})

	t.Run("Multiple keys can be stored and retrieved", func(t *testing.T) {
		keys := []string{"key1", "key2", "key3", "key4", "key5"}
		values := []string{"val1", "val2", "val3", "val4", "val5"}

		for i := 0; i < len(keys); i++ {
			err := repo.Set(ctx, keys[i], values[i], 3600)
			require.NoError(t, err)
		}

		for i := 0; i < len(keys); i++ {
			val, err := repo.Get(ctx, keys[i])
			require.NoError(t, err)
			require.Equal(t, values[i], val)
		}
	})

	t.Run("Key expiration works", func(t *testing.T) {
		err := repo.Set(ctx, "expire_quick", "value", 1)
		require.NoError(t, err)

		time.Sleep(2 * time.Second)

		_, err = repo.Get(ctx, "expire_quick")
		require.Error(t, err)
	})

	t.Run("Exists after deletion returns false", func(t *testing.T) {
		err := repo.Set(ctx, "check_exists", "value", 3600)
		require.NoError(t, err)

		exists, err := repo.Exists(ctx, "check_exists")
		require.NoError(t, err)
		require.True(t, exists)

		err = repo.Delete(ctx, "check_exists")
		require.NoError(t, err)

		exists, err = repo.Exists(ctx, "check_exists")
		require.NoError(t, err)
		require.False(t, exists)
	})

	t.Run("GetAndDelete on expired key returns error", func(t *testing.T) {
		err := repo.Set(ctx, "expired_get_delete", "value", 1)
		require.NoError(t, err)

		time.Sleep(2 * time.Second)

		_, err = repo.GetAndDelete(ctx, "expired_get_delete")
		require.Error(t, err)
	})

	t.Run("Clear idempotent", func(t *testing.T) {
		err := repo.Set(ctx, "before_clear", "value", 3600)
		require.NoError(t, err)

		err = repo.Clear(ctx)
		require.NoError(t, err)

		err = repo.Clear(ctx)
		require.NoError(t, err)

		exists, _ := repo.Exists(ctx, "before_clear")
		require.False(t, exists)
	})

	t.Run("Set empty string value", func(t *testing.T) {
		err := repo.Set(ctx, "empty_string", "", 3600)
		require.NoError(t, err)

		val, err := repo.Get(ctx, "empty_string")
		require.NoError(t, err)
		require.Equal(t, "", val)
	})

	t.Run("Complex struct can be cached", func(t *testing.T) {
		type TestStruct struct {
			Name  string
			Count int
			Items []string
		}

		data := TestStruct{
			Name:  "test",
			Count: 42,
			Items: []string{"a", "b", "c"},
		}

		err := repo.Set(ctx, "struct_key", data, 3600)
		require.NoError(t, err)

		val, err := repo.Get(ctx, "struct_key")
		require.NoError(t, err)
		require.NotNil(t, val)
	})

	t.Run("Large value can be cached", func(t *testing.T) {
		largeValue := ""
		for i := 0; i < 10000; i++ {
			largeValue += "a"
		}

		err := repo.Set(ctx, "large_value", largeValue, 3600)
		require.NoError(t, err)

		val, err := repo.Get(ctx, "large_value")
		require.NoError(t, err)
		require.Equal(t, largeValue, val)
	})

	t.Run("SetWithExpiration and regular Set work together", func(t *testing.T) {
		futureTime := time.Now().Add(1 * time.Hour).Unix()

		err := repo.Set(ctx, "regular_set", "value1", 3600)
		require.NoError(t, err)

		err = repo.SetWithExpiration(ctx, "expire_set", "value2", futureTime)
		require.NoError(t, err)

		val1, _ := repo.Get(ctx, "regular_set")
		val2, _ := repo.Get(ctx, "expire_set")

		require.Equal(t, "value1", val1)
		require.Equal(t, "value2", val2)
	})
}
