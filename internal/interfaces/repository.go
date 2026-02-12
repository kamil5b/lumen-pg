package interfaces

import (
	"context"
	"github.com/kamil5b/lumen-pg/internal/domain"
)

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id string) (*domain.User, error)
}
