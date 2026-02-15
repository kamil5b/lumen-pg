package logger_repository

import (
	"context"
	"errors"
)

func (l *LoggerRepositoryImplementation) LogSecurityEvent(ctx context.Context, eventType string, username string, details map[string]interface{}) error {
	return errors.New("not implemented")
}
