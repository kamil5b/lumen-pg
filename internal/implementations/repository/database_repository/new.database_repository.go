package database_repository

import (
	"database/sql"

	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
)

type DatabaseRepositoryImplementation struct {
	db *sql.DB
}

func NewDatabaseRepository(db *sql.DB) repository.DatabaseRepository {
	return &DatabaseRepositoryImplementation{
		db: db,
	}
}
