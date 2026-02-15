package transaction

import (
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
)

type TransactionUseCaseImplementation struct {
	transactionRepo repository.TransactionRepository
	databaseRepo    repository.DatabaseRepository
	rbacRepo        repository.RBACRepository
}

func NewTransactionUseCaseImplementation(
	transactionRepo repository.TransactionRepository,
	databaseRepo repository.DatabaseRepository,
	rbacRepo repository.RBACRepository,
) usecase.TransactionUseCase {
	return &TransactionUseCaseImplementation{
		transactionRepo: transactionRepo,
		databaseRepo:    databaseRepo,
		rbacRepo:        rbacRepo,
	}
}
