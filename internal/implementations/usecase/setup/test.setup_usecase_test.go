package setup

import (
	"testing"

	testRunner "github.com/kamil5b/lumen-pg/internal/testrunners/usecase"
)

func TestSetupUsecase(t *testing.T) {
	testRunner.SetupUsecaseRunner(t, NewSetupUseCaseImplementation)
}
