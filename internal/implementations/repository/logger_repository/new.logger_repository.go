package logger_repository

import (
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
)

type LoggerRepositoryImplementation struct {
	// logger implementation details will be added here
}

func NewLoggerRepository() repository.LoggerRepository {
	return &LoggerRepositoryImplementation{}
}
