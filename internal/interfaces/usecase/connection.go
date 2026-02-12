package usecase

import (
	"context"
)

// ConnectionUseCase handles connection validation operations
type ConnectionUseCase interface {
	ValidateConnectionString(connStr string) error
	TestConnection(ctx context.Context, connStr string) error
}
