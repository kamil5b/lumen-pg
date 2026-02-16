package transaction

import (
	"github.com/kamil5b/lumen-pg/internal/interfaces/handler"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
)

type TransactionHandlerImplementation struct {
	transactionUC usecase.TransactionUseCase
	authUC        usecase.AuthenticationUseCase
	rbacUC        usecase.RBACUseCase
}

func NewTransactionHandlerImplementation(
	transactionUC usecase.TransactionUseCase,
	authUC usecase.AuthenticationUseCase,
	rbacUC usecase.RBACUseCase,
) handler.TransactionHandler {
	return &TransactionHandlerImplementation{
		transactionUC: transactionUC,
		authUC:        authUC,
		rbacUC:        rbacUC,
	}
}
