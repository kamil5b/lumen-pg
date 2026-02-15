package erd

import (
	"testing"

	testRunner "github.com/kamil5b/lumen-pg/internal/testrunners/usecase"
)

func TestERDUser(t *testing.T) {
	testRunner.ERDUsecaseRunner(t, NewERDUseCaseImplementation)
}
