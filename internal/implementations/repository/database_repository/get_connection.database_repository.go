package database_repository

import (
	"database/sql"
)

func (d *DatabaseRepositoryImplementation) GetConnection() *sql.DB {
	return d.db
}
