package database_repository

import (
	"testing"

	testRunner "github.com/kamil5b/lumen-pg/internal/testrunners/repository"
)

func TestDatabaseRepository(t *testing.T) {
	testRunner.DatabaseRepositoryRunner(t, NewDatabaseRepository)
}
