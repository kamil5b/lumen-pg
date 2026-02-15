package logger_repository

import (
	"context"
	"errors"
)

func (l *LoggerRepositoryImplementation) LogTransactionEvent(ctx context.Context, username string, eventType string, details map[string]interface{}) error {
	return errors.New("not implemented")
}
