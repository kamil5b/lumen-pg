package metadata_repository

import (
	"sync"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
)

type MetadataRepositoryImplementation struct {
	mu            sync.RWMutex
	databases     map[string]*domain.DatabaseMetadata
	rolesMetadata map[string]*domain.RoleMetadata
}

func NewMetadataRepository() repository.MetadataRepository {
	return &MetadataRepositoryImplementation{
		databases:     make(map[string]*domain.DatabaseMetadata),
		rolesMetadata: make(map[string]*domain.RoleMetadata),
	}
}
