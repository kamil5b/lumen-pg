package transaction

import (
	"testing"

	testRunner "github.com/kamil5b/lumen-pg/internal/testrunners/usecase"
)

func TestTransactionUsecase(t *testing.T) {
	testRunner.TransactionUsecaseRunner(t, NewTransactionUseCaseImplementation)
}
