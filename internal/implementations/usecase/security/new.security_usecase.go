package security

import (
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
)

type SecurityUseCaseImplementation struct {
	securityRepo repository.SecurityRepository
}

func NewSecurityUseCaseImplementation(
	securityRepo repository.SecurityRepository,
) usecase.SecurityUseCase {
	return &SecurityUseCaseImplementation{
		securityRepo: securityRepo,
	}
}
