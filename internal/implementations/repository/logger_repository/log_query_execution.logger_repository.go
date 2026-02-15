package logger_repository

import (
	"context"
	"errors"
)

func (l *LoggerRepositoryImplementation) LogQueryExecution(ctx context.Context, username string, query string, executionTimeMs int64, success bool, err error) error {
	return errors.New("not implemented")
}
