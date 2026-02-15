package query

import (
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
)

type QueryUseCaseImplementation struct {
	databaseRepo repository.DatabaseRepository
	rbacRepo     repository.RBACRepository
}

func NewQueryUseCaseImplementation(
	databaseRepo repository.DatabaseRepository,
	rbacRepo repository.RBACRepository,
) usecase.QueryUseCase {
	return &QueryUseCaseImplementation{
		databaseRepo: databaseRepo,
		rbacRepo:     rbacRepo,
	}
}
