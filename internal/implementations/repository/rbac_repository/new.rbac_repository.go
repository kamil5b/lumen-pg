package rbac_repository

import (
	"database/sql"

	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
)

type RBACRepositoryImplementation struct {
	db *sql.DB
}

func NewRBACRepository(db *sql.DB) repository.RBACRepository {
	return &RBACRepositoryImplementation{
		db: db,
	}
}
