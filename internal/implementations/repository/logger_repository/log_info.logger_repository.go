package logger_repository

import (
	"context"
	"errors"
)

func (l *LoggerRepositoryImplementation) LogInfo(ctx context.Context, message string, fields map[string]interface{}) error {
	return errors.New("not implemented")
}
