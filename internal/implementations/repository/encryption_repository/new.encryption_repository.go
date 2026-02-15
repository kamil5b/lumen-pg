package encryption_repository

import (
	"database/sql"

	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
)

type EncryptionRepositoryImplementation struct {
	db *sql.DB
}

func NewEncryptionRepository(db *sql.DB) repository.EncryptionRepository {
	return &EncryptionRepositoryImplementation{
		db: db,
	}
}
