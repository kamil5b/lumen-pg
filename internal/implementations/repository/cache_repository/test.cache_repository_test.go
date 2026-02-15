package cache_repository

import (
	"testing"

	testRunner "github.com/kamil5b/lumen-pg/internal/testrunners/repository"
)

func TestCacheRepository(t *testing.T) {
	testRunner.CacheRepositoryRunner(t, NewCacheRepository)
}
