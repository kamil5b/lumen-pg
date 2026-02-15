package transaction_repository

import (
	"testing"

	testRunner "github.com/kamil5b/lumen-pg/internal/testrunners/repository"
)

func TestTransactionRepository(t *testing.T) {
	testRunner.TransactionRepositoryRunner(t, NewTransactionRepository)
}
