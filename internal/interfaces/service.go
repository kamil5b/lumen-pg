package interfaces

import (
	"context"
	"github.com/kamil5b/lumen-pg/internal/domain"
)

type UserService interface {
	CreateUser(ctx context.Context, input domain.CreateUserInput) (*domain.User, error)
	GetUser(ctx context.Context, id string) (*domain.User, error)
}
