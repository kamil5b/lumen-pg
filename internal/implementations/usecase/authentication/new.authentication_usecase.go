package authentication

import (
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
)

type AuthenticationUseCaseImplementation struct {
	authRepo     repository.AuthenticationRepository
	sessionRepo  repository.SessionRepository
	rbacRepo     repository.RBACRepository
	metadataRepo repository.MetadataRepository
}

func NewAuthenticationUseCaseImplementation(
	authRepo repository.AuthenticationRepository,
	sessionRepo repository.SessionRepository,
	rbacRepo repository.RBACRepository,
	metadataRepo repository.MetadataRepository,
) usecase.AuthenticationUseCase {
	return &AuthenticationUseCaseImplementation{
		authRepo:     authRepo,
		sessionRepo:  sessionRepo,
		rbacRepo:     rbacRepo,
		metadataRepo: metadataRepo,
	}
}
