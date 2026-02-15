package metadata_repository

import (
	"testing"

	testRunner "github.com/kamil5b/lumen-pg/internal/testrunners/repository"
)

func TestMetadataRepository(t *testing.T) {
	testRunner.MetadataRepositoryRunner(t, NewMetadataRepository)
}
