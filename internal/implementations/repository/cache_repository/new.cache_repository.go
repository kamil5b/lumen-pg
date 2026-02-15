package cache_repository

import (
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
)

type CacheRepositoryImplementation struct {
	// cache implementation details will be added here
}

func NewCacheRepository() repository.CacheRepository {
	return &CacheRepositoryImplementation{}
}
