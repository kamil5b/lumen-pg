package erd

import "github.com/kamil5b/lumen-pg/internal/interfaces/repository"

type ERDUseCaseImplementation struct {
	metadataRepo repository.MetadataRepository
	rbacRepo     repository.RBACRepository
}

func NewERDUseCaseImplementation(
	metadataRepo repository.MetadataRepository,
	rbacRepo repository.RBACRepository,
) *ERDUseCaseImplementation {
	return &ERDUseCaseImplementation{
		metadataRepo: metadataRepo,
		rbacRepo:     rbacRepo,
	}
}
