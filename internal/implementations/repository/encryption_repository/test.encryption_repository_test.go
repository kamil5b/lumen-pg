package encryption_repository

import (
	"testing"

	testRunner "github.com/kamil5b/lumen-pg/internal/testrunners/repository"
)

func TestEncryptionRepository(t *testing.T) {
	testRunner.EncryptionRepositoryRunner(t, NewEncryptionRepository)
}
