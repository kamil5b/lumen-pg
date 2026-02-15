package query

import (
	"testing"

	testRunner "github.com/kamil5b/lumen-pg/internal/testrunners/usecase"
)

func TestQueryUsecase(t *testing.T) {
	testRunner.QueryUsecaseRunner(t, NewQueryUseCaseImplementation)
}
