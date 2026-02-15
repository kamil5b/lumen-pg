package erd

import (
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
)

type ERDUseCaseImplementation struct {
	metadataRepo repository.MetadataRepository
	rbacRepo     repository.RBACRepository
}

func NewERDUseCaseImplementation(
	metadataRepo repository.MetadataRepository,
	rbacRepo repository.RBACRepository,
) usecase.ERDUseCase {
	return &ERDUseCaseImplementation{
		metadataRepo: metadataRepo,
		rbacRepo:     rbacRepo,
	}
}
