package metadata_repository

import (
	"database/sql"
	"sync"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
)

type MetadataRepositoryImplementation struct {
	mu            sync.RWMutex
	db            *sql.DB
	databases     map[string]*domain.DatabaseMetadata
	rolesMetadata map[string]*domain.RoleMetadata
}

func NewMetadataRepository(db *sql.DB) repository.MetadataRepository {
	return &MetadataRepositoryImplementation{
		db:            db,
		databases:     make(map[string]*domain.DatabaseMetadata),
		rolesMetadata: make(map[string]*domain.RoleMetadata),
	}
}
