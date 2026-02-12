package usecase

import (
	"context"
	"github.com/kamil5b/lumen-pg/internal/domain"
)

// AuthUseCase handles authentication operations
type AuthUseCase interface {
	Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error)
	Logout(ctx context.Context, sessionToken string) error
}
