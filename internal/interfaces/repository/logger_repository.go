package repository

import (
	"context"
)

// LoggerRepository defines operations for logging
type LoggerRepository interface {
	// LogInfo logs an informational message
	LogInfo(ctx context.Context, message string, fields map[string]interface{}) error

	// LogWarn logs a warning message
	LogWarn(ctx context.Context, message string, fields map[string]interface{}) error

	// LogError logs an error message
	LogError(ctx context.Context, message string, err error, fields map[string]interface{}) error

	// LogDebug logs a debug message
	LogDebug(ctx context.Context, message string, fields map[string]interface{}) error

	// LogSecurityEvent logs a security-related event
	LogSecurityEvent(ctx context.Context, eventType string, username string, details map[string]interface{}) error

	// LogQueryExecution logs a query execution
	LogQueryExecution(ctx context.Context, username string, query string, executionTimeMs int64, success bool, err error) error

	// LogTransactionEvent logs a transaction event
	LogTransactionEvent(ctx context.Context, username string, eventType string, details map[string]interface{}) error
}
