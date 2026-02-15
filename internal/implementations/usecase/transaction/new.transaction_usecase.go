package transaction

import (
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
)

type TransactionUseCaseImplementation struct {
	transactionRepo repository.TransactionRepository
	dataRepo        repository.DataRepository
	rbacUseCase     usecase.RBACUseCase
}

func NewTransactionUseCaseImplementation(
	transactionRepo repository.TransactionRepository,
	dataRepo repository.DataRepository,
	rbacUseCase usecase.RBACUseCase,
) usecase.TransactionUseCase {
	return &TransactionUseCaseImplementation{
		transactionRepo: transactionRepo,
		dataRepo:        dataRepo,
		rbacUseCase:     rbacUseCase,
	}
}
