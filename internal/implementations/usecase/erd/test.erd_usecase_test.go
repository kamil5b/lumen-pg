package erd

import (
	"testing"

	testRunner "github.com/kamil5b/lumen-pg/internal/testrunners/usecase"
)

func TestERDUsecase(t *testing.T) {
	testRunner.ERDUsecaseRunner(t, NewERDUseCaseImplementation)
}
