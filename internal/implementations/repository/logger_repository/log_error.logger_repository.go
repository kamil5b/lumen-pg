package logger_repository

import (
	"context"
	"errors"
)

func (l *LoggerRepositoryImplementation) LogError(ctx context.Context, message string, err error, fields map[string]interface{}) error {
	return errors.New("not implemented")
}
