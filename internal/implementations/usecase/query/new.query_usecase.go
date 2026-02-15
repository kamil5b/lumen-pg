package query

import (
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
)

type QueryUseCaseImplementation struct {
	queryRepo       repository.QueryRepository
	rbacUseCase     usecase.RBACUseCase
	securityUseCase usecase.SecurityUseCase
}

func NewQueryUseCaseImplementation(
	queryRepo repository.QueryRepository,
	rbacUseCase usecase.RBACUseCase,
	securityUseCase usecase.SecurityUseCase,
) usecase.QueryUseCase {
	return &QueryUseCaseImplementation{
		queryRepo:       queryRepo,
		rbacUseCase:     rbacUseCase,
		securityUseCase: securityUseCase,
	}
}
