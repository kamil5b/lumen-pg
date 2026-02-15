package dataview

import (
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
)

type DataViewUseCaseImplementation struct {
	metadataRepo repository.MetadataRepository
	databaseRepo repository.DatabaseRepository
	rbacRepo     repository.RBACRepository
}

func NewDataViewUseCaseImplementation(
	metadataRepo repository.MetadataRepository,
	databaseRepo repository.DatabaseRepository,
	rbacRepo repository.RBACRepository,
) usecase.DataViewUseCase {
	return &DataViewUseCaseImplementation{
		metadataRepo: metadataRepo,
		databaseRepo: databaseRepo,
		rbacRepo:     rbacRepo,
	}
}
