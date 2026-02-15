package security

import (
	"testing"

	testRunner "github.com/kamil5b/lumen-pg/internal/testrunners/usecase"
)

func TestSecurityUsecase(t *testing.T) {
	testRunner.SecurityUsecaseRunner(t, NewSecurityUseCaseImplementation)
}
