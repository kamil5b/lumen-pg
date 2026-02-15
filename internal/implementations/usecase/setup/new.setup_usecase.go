package setup

import (
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
)

type SetupUseCaseImplementation struct {
	metadataRepo repository.MetadataRepository
	rbacRepo     repository.RBACRepository
	dataRepo     repository.DataRepository
}

func NewSetupUseCaseImplementation(
	metadataRepo repository.MetadataRepository,
	rbacRepo repository.RBACRepository,
	dataRepo repository.DataRepository,
) usecase.SetupUseCase {
	return &SetupUseCaseImplementation{
		metadataRepo: metadataRepo,
		rbacRepo:     rbacRepo,
		dataRepo:     dataRepo,
	}
}
