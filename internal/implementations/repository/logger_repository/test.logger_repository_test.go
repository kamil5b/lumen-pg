package logger_repository

import (
	"testing"

	testRunner "github.com/kamil5b/lumen-pg/internal/testrunners/repository"
)

func TestLoggerRepository(t *testing.T) {
	testRunner.LoggerRepositoryRunner(t, NewLoggerRepository)
}
