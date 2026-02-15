package rbac

import (
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
)

type RBACUseCaseImplementation struct {
	rbacRepo     repository.RBACRepository
	metadataRepo repository.MetadataRepository
}

func NewRBACUseCaseImplementation(
	rbacRepo repository.RBACRepository,
	metadataRepo repository.MetadataRepository,
) usecase.RBACUseCase {
	return &RBACUseCaseImplementation{
		rbacRepo:     rbacRepo,
		metadataRepo: metadataRepo,
	}
}
