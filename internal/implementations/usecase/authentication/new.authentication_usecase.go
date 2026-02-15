package authentication

import (
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
)

type AuthenticationUseCaseImplementation struct {
	databaseRepo   repository.DatabaseRepository
	metadataRepo   repository.MetadataRepository
	sessionRepo    repository.SessionRepository
	rbacRepo       repository.RBACRepository
	encryptionRepo repository.EncryptionRepository
}

func NewAuthenticationUseCaseImplementation(
	databaseRepo repository.DatabaseRepository,
	metadataRepo repository.MetadataRepository,
	sessionRepo repository.SessionRepository,
	rbacRepo repository.RBACRepository,
	encryptionRepo repository.EncryptionRepository,
) usecase.AuthenticationUseCase {
	return &AuthenticationUseCaseImplementation{
		databaseRepo:   databaseRepo,
		metadataRepo:   metadataRepo,
		sessionRepo:    sessionRepo,
		rbacRepo:       rbacRepo,
		encryptionRepo: encryptionRepo,
	}
}
