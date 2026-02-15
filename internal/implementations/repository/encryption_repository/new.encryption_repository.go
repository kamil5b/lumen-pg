package encryption_repository

import (
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
)

type EncryptionRepositoryImplementation struct {
	// encryption implementation details will be added here
}

func NewEncryptionRepository() repository.EncryptionRepository {
	return &EncryptionRepositoryImplementation{}
}
