package rbac_repository

import (
	"testing"

	testRunner "github.com/kamil5b/lumen-pg/internal/testrunners/repository"
)

func TestRBACRepository(t *testing.T) {
	testRunner.RBACRepositoryRunner(t, NewRBACRepository)
}
