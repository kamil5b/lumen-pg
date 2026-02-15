package clock_repository

import (
	"testing"

	testRunner "github.com/kamil5b/lumen-pg/internal/testrunners/repository"
)

func TestClockRepository(t *testing.T) {
	testRunner.ClockRepositoryRunner(t, NewClockRepository)
}
