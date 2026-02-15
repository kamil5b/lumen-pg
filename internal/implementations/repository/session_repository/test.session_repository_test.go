package session_repository

import (
	"testing"

	testRunner "github.com/kamil5b/lumen-pg/internal/testrunners/repository"
)

func TestSessionRepository(t *testing.T) {
	testRunner.SessionRepositoryRunner(t, NewSessionRepository)
}
