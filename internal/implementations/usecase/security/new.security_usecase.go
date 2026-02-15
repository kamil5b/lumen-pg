package security

import (
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
)

type SecurityUseCaseImplementation struct {
	encryptionRepo repository.EncryptionRepository
	sessionRepo    repository.SessionRepository
	clockRepo      repository.ClockRepository
}

func NewSecurityUseCaseImplementation(
	encryptionRepo repository.EncryptionRepository,
	sessionRepo repository.SessionRepository,
	clockRepo repository.ClockRepository,
) usecase.SecurityUseCase {
	return &SecurityUseCaseImplementation{
		encryptionRepo: encryptionRepo,
		sessionRepo:    sessionRepo,
		clockRepo:      clockRepo,
	}
}
