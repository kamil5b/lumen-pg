package rbac

import (
	"testing"

	testRunner "github.com/kamil5b/lumen-pg/internal/testrunners/usecase"
)

func TestRBACUsecase(t *testing.T) {
	testRunner.RBACUsecaseRunner(t, NewRBACUseCaseImplementation)
}
