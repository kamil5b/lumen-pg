package setup

import (
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
)

type SetupUseCaseImplementation struct {
	databaseRepo repository.DatabaseRepository
	metadataRepo repository.MetadataRepository
	rbacRepo     repository.RBACRepository
}

func NewSetupUseCaseImplementation(
	databaseRepo repository.DatabaseRepository,
	metadataRepo repository.MetadataRepository,
	rbacRepo repository.RBACRepository,
) usecase.SetupUseCase {
	return &SetupUseCaseImplementation{
		databaseRepo: databaseRepo,
		metadataRepo: metadataRepo,
		rbacRepo:     rbacRepo,
	}
}
