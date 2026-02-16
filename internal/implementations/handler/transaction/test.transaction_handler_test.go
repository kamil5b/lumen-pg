package transaction_test

import (
	"testing"

	"github.com/kamil5b/lumen-pg/internal/implementations/handler/transaction"
	"github.com/kamil5b/lumen-pg/internal/interfaces/handler"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
	handlerTestRunner "github.com/kamil5b/lumen-pg/internal/testrunners/handler"
)

func TestTransactionHandler(t *testing.T) {
	constructor := func(
		txnUC usecase.TransactionUseCase,
		authUC usecase.AuthenticationUseCase,
		rbacUC usecase.RBACUseCase,
	) handler.TransactionHandler {
		return transaction.NewTransactionHandlerImplementation(txnUC, authUC, rbacUC)
	}

	handlerTestRunner.TransactionHandlerRunner(t, constructor)
}
